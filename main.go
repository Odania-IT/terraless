package main

import (
	"github.com/Sirupsen/logrus"
	"os"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
	})

	logrus.Info("Running terraless")
	arguments := parseArguments()
	config := newTerralessConfig(arguments)

	logrus.Debug("Config", config)

	logrus.Debug("Terraless target folder: ", config.TargetPath)
	err := os.MkdirAll(config.TargetPath, os.ModePerm)

	if err != nil {
		logrus.Fatal("Error creating target directory: ", err)
	}

	renderTemplates(config)
	// findAndCopyFiles(filepath.Join(config.SourcePath, "*.tf"), config.TargetPath)

	execSwamp(config)
	deploy(config, arguments.Environment, arguments.ForceDeploy, arguments.TerraformCommand)
}
