package main

import (
	"bytes"
	"fmt"
	"github.com/Odania-IT/terraless/config"
	"github.com/Odania-IT/terraless/plugin"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defer closePlugins()
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
	})

	logrus.Info("Running terraless")
	arguments, kingpinResult := parseArguments(os.Args[1:])

	if kingpinResult == versionCommand.FullCommand() {
		versionInfo()
		return
	}

	plugin.ExistingPlugins(arguments)
	terralessData := config.NewTerralessData(arguments, plugin.Providers())
	processCommands(terralessData, kingpinResult)
}

func closePlugins() {
	plugin.ClosePlugins()
}

func processCommands(terralessData *schema.TerralessData, kingpinResult string) {
	arguments := terralessData.Arguments
	currentConfig := terralessData.Config
	currentConfig.Settings.Variables = schema.EnrichWithData(currentConfig.Settings.Variables, arguments.Variables)

	logrus.Debugf("Active Providers in Config: %d\n", len(currentConfig.Providers))

	logrus.Debug("Config", currentConfig)

	switch kingpinResult {
	case authCommand.FullCommand():
		logrus.Debug("Handling Auth Command")

		if authProviderName == nil || *authProviderName == "" {
			stepPrepareSesssion(terralessData)
		}

	case deployCommand.FullCommand():
		logrus.Debug("Handling Deploy Command")
		stepDeploy(terralessData)
	case initCommand.FullCommand():
		logrus.Debug("Handling Init Command")
		plugin.HandlePlugins(terralessData.Plugins, terralessData.Arguments.PluginDirectory)

	case initTemplatesCommand.FullCommand():
		logrus.Debug("Handling Init-Templates Command")
		stepInitialize(terralessData)
	case uploadCommand.FullCommand():
		logrus.Debug("Handling Upload Command")
		stepUpload(terralessData)
	case infoCommand.FullCommand():
		logrus.Debug("Handling Info Command")
		versionInfo()
		fmt.Println()
		fmt.Printf("Global Config: %s\n", arguments.GlobalConfig)
		fmt.Printf("Project Config: %s\n", arguments.Config)

		var allProviders []string
		for _, provider := range terralessData.Config.Providers {
			allProviders = append(allProviders, provider.Name)
		}
		fmt.Printf("Providers: %s\n", strings.Join(allProviders, ", "))

		var allTerralessExtensions []string
		var allTerralessProviders []string
		for _, pluginData := range plugin.PluginsData() {
			info := fmt.Sprintf("%s-%s", pluginData.Name, pluginData.Version)
			if pluginData.Type.Type == plugin.ExtensionPluginType {
				allTerralessExtensions = append(allTerralessExtensions, info)
			} else if pluginData.Type.Type == plugin.ProviderPluginType {
				allTerralessProviders = append(allTerralessProviders, info)
			}
		}
		fmt.Printf("Terraless Extensions: %s\n", strings.Join(allTerralessExtensions, ", "))
		fmt.Printf("Terraless Providers: %s\n", strings.Join(allTerralessProviders, ", "))

		fmt.Println("Variables:")
		for key, val := range terralessData.Config.Settings.Variables {
			fmt.Printf("- %s: %s\n", key, val)
		}
		if len(terralessData.Config.Settings.Variables) == 0 {
			fmt.Println("  none")
		}
	default:
		logrus.Debug("Invalid step")
		kingpin.Usage()
	}
}

func versionInfo() {
	fmt.Printf("Terraless Version: %s [Codename: %s]\n", VERSION, CODENAME)
}

func stepDeploy(terralessData *schema.TerralessData) {

	stepInitialize(terralessData)
	stepPrepareSesssion(terralessData)

	arguments := terralessData.Arguments
	currentConfig := terralessData.Config
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

		deployTerraform(currentConfig, arguments.Environment, arguments.ForceDeploy, arguments.TerraformCommand)

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
		templates.ProcessUploads(*terralessData)
	}
}

func stepInitialize(terralessData *schema.TerralessData) {
	buffer := bytes.Buffer{}
	buffer = templates.Render(terralessData, buffer)
	targetFileName := filepath.Join(terralessData.Config.SourcePath, "terraless-resources.tf")

	if buffer.Len() == 0 {
		logrus.Debug("Nothing to write to terraless-resources.tf")

		// Remove terraless-resources.tf if it exists
		if _, err := os.Stat(targetFileName); err == nil {
			err := os.Remove(targetFileName)

			if err != nil {
				logrus.Fatal("Failed to remove terraless-resources.tf")
			}
		}

		return
	}

	// Writing buffer to file
	logrus.Debugf("Writing file %s\n", targetFileName)

	finalBuffer := bytes.Buffer{}
	finalBuffer.WriteString("# This file is generated by Terraless\n\n")
	finalBuffer.Write(buffer.Bytes())
	support.WriteToFile(targetFileName, finalBuffer)
}

func stepPrepareSesssion(terralessData *schema.TerralessData) {
	if support.RunningInAws() && !terralessData.Config.Settings.AutoSignInInCloud {
		logrus.Info("Not executing prepare session! Settings: AutoSignInInCloud is false")
		return
	}

	for _, terralessProvider := range terralessData.TerralessProviders {
		terralessProvider.PrepareSession(terralessData.Config)
	}
}

func stepUpload(terralessData *schema.TerralessData) {
	stepInitialize(terralessData)
	stepPrepareSesssion(terralessData)
	templates.ProcessUploads(*terralessData)
}

func deployTerraform(config schema.TerralessConfig, environment string, forceDeploy bool, terraformCommand string) {
	logrus.Info("Executing terraform init")
	initArguments := []string{"init"}

	if config.Settings.TerraformPluginDir != "" {
		initArguments = append(initArguments, "-plugin-dir", config.Settings.TerraformPluginDir)
	}

	executeCommand(config.SourcePath, terraformCommand, initArguments, false)
	logrus.Info("Creating new terraform workspace")
	executeCommand(config.SourcePath, terraformCommand, []string{"workspace", "new", config.ProjectName}, true)
	logrus.Info("Selecting new terraform workspace")
	executeCommand(config.SourcePath, terraformCommand, []string{"workspace", "select", config.ProjectName}, false)

	logrus.Info("Executing terraform plan")
	planArgs := []string{
		"plan",
		"-out",
		"terraform.plan",
		"-input=false",
		"-var",
		"environment=" + environment,
	}
	executeCommand(config.SourcePath, terraformCommand, planArgs, false)

	if forceDeploy || checkApprove(os.Stdin) {
		logrus.Info("Deploying terraform plan")
		executeCommand(config.SourcePath, terraformCommand, []string{"apply", "-input=false", "terraform.plan"}, false)
	} else {
		logrus.Info("Not deploying...")
	}
}
