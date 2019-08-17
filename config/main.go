package config

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

var globalConfig *schema.TerralessGlobalConfig

func NewTerralessData(arguments schema.Arguments) *schema.TerralessData {
	terralessData := &schema.TerralessData{
		Arguments:          arguments,
	}

	projectConfig := readProjectYamlConfig(terralessData.Arguments)

	if arguments.AuthProvider != "" {
		parts := strings.Split(arguments.AuthProvider, ":")

		if len(parts) < 2 {
			logrus.Fatal("Invalid format for auth provider!")
		}

		projectConfig.ActiveProviders = []schema.TerralessActiveProvider{
			{
				Team: parts[0],
				Providers: []schema.TerralessProvider{
					{
						Type: "global",
						Name: parts[1],
						Data: dataFromParts(parts),
					},
				},
			},
		}
		projectConfig.Settings.AutoSignIn = true
	}

	// Set ProjectName if none is in the config
	if projectConfig.ProjectName == "" {
		projectConfig.ProjectName = filepath.Base(filepath.Dir(terralessData.Arguments.Config))
	}

	// Set SourcePath if none is in the projectConfig
	if projectConfig.SourcePath == "" {
		projectConfig.SourcePath = filepath.Join(filepath.Dir(terralessData.Arguments.Config))
	}

	// Set TargetPath if none is in the projectConfig
	if projectConfig.TargetPath == "" {
		projectConfig.TargetPath = filepath.Join(filepath.Dir(terralessData.Arguments.Config), ".terraless")
	}

	logrus.Debug("Terraless target folder: ", projectConfig.TargetPath)
	err := os.MkdirAll(projectConfig.TargetPath, os.ModePerm)

	if err != nil {
		logrus.Fatalf("Error creating target directory: %s\n", err)
	}

	globalConfig = readGlobalYamlConfig(terralessData.Arguments)

	logrus.Debugln(projectConfig.ActiveProviders)

	terralessData.Config = schema.BuildTerralessConfig(*globalConfig, *projectConfig, arguments)
	terralessData.Plugins = globalConfig.Plugins
	validate := terralessData.Config.Validate()

	if len(validate) > 0 {
		logrus.Fatalf("Failed to verify config! %s\n", strings.Join(validate, ", "))
	}

	return terralessData
}

func GlobalConfig() schema.TerralessGlobalConfig {
	return *globalConfig
}

func dataFromParts(data []string) map[string]string {
	result := map[string]string{}

	for idx, keyVal := range data {
		if idx > 1 {
			parts := strings.Split(keyVal, "=")
			result[parts[0]] = parts[1]
		}
	}

	return result
}
