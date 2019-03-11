package main

import (
	"github.com/Odania-IT/terraless/config"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/terraless_provider_aws"
	"github.com/Odania-IT/terraless/uploads"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func detectTerralessProviders() []schema.Provider {
	var terralessProviders []schema.Provider
	terralessProviders = append(terralessProviders, terraless_provider_aws.Provider())

	return terralessProviders
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
	})

	var terralessProviders []schema.Provider
	terralessProviders = append(terralessProviders, terraless_provider_aws.Provider())

	logrus.Info("Running terraless")
	arguments := config.ParseArguments()
	terralessData := config.NewTerralessData(arguments, detectTerralessProviders())
	currentConfig := terralessData.Config

	logrus.Debugf("Active Providers in Config: %d\n", len(currentConfig.Providers))

	logrus.Debug("Config", currentConfig)

	logrus.Debug("Terraless target folder: ", currentConfig.TargetPath)
	err := os.MkdirAll(currentConfig.TargetPath, os.ModePerm)

	if err != nil {
		logrus.Fatalf("Error creating target directory: %s\n", err)
	}

	for _, terralessProvider := range terralessProviders {
		terralessProvider.PrepareSession(currentConfig)
	}

	terraless_provider_aws.RenderTemplates(*terralessData)

	if arguments.NoDeploy {
		logrus.Debug("Not deploying due to arguments....")
	} else {
		if currentConfig.Package.SourceDir != "" {
			logrus.Debugf("Executing before package hooks depending on runtime")
			if support.ContainsStartsWith(currentConfig.Runtimes, "ruby") {
				executeCommand(filepath.Join(currentConfig.SourcePath, currentConfig.Package.SourceDir), "bundle",
					[]string{
						"install", "--deployment", "--without", "test", "development",
					}, false)
			}
		}

		deploy(currentConfig, arguments.Environment, arguments.ForceDeploy, arguments.TerraformCommand)

		if currentConfig.Package.SourceDir != "" {
			logrus.Debugf("Executing after package hooks depending on runtime")
			if support.ContainsStartsWith(currentConfig.Runtimes, "ruby") {
				executeCommand(filepath.Join(currentConfig.SourcePath, currentConfig.Package.SourceDir), "bundle",
					[]string{
						"install", "--quiet", "--no-deployment", "--with", "test", "development",
					}, false)
			}
		}
	}

	if arguments.NoUpload {
		logrus.Debug("Not uploading due to arguments....")
	} else {
		uploads.ProcessUploads(*terralessData)
	}
}
