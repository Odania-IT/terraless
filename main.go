package main

import (
	"fmt"
	"github.com/Odania-IT/terraless/config"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/terraless_provider_aws"
	"github.com/Odania-IT/terraless/uploads"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"path/filepath"
)

var terralessProviders []schema.Provider

func detectTerralessProviders() []schema.Provider {
	var terralessProviders []schema.Provider
	terralessProviders = append(terralessProviders, terraless_provider_aws.Provider())

	return terralessProviders
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
	})

	terralessProviders = append(terralessProviders, terraless_provider_aws.Provider())

	logrus.Info("Running terraless")
	arguments, kingpinResult := parseArguments()
	terralessData := config.NewTerralessData(arguments, detectTerralessProviders())
	currentConfig := terralessData.Config

	logrus.Debugf("Active Providers in Config: %d\n", len(currentConfig.Providers))

	logrus.Debug("Config", currentConfig)

	switch kingpinResult {
	case deployCommand.FullCommand():
		logrus.Debug("Handling Deploy Command")
		stepDeploy(terralessData)
	case sessionCommand.FullCommand():
		logrus.Debug("Handling Session Command")
		stepPrepareSesssion(terralessData)
	case versionCommand.FullCommand():
		fmt.Printf("Terraless Version: %s [Codename: %s]\n", VERSION, CODENAME)
	default:
		logrus.Debug("Invalid step")
		kingpin.Usage()
	}
}

func stepDeploy(terralessData *schema.TerralessData) {
	arguments := terralessData.Arguments
	currentConfig := terralessData.Config

	stepPrepareSesssion(terralessData)

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

func stepPrepareSesssion(terralessData *schema.TerralessData) {
	for _, terralessProvider := range terralessProviders {
		terralessProvider.PrepareSession(terralessData.Config)
	}
}
