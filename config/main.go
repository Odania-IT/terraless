package config

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

func NewTerralessData(arguments schema.Arguments, terralessProviders []schema.Provider) *schema.TerralessData {
	terralessData := &schema.TerralessData{
		Arguments:          arguments,
		TerralessProviders: terralessProviders,
	}

	projectConfig := readProjectYamlConfig(terralessData.Arguments)
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

	// for _, terralessProvider := range terralessProviders {
	// 	activeProviders := terralessProvider.FilterActiveProviders(*projectConfig)
	// 	for _, activeProvider := range activeProviders {
	// 		terralessData.ActiveProviders[activeProvider.Name] = activeProvider
	// 	}
	// }

	globalConfig := readGlobalYamlConfig(terralessData.Arguments)
	// globalConfig.BuildFinalConfig(projectConfig, arguments)

	logrus.Debugln(projectConfig.ActiveProviders)

	terralessData.Config = schema.BuildTerralessConfig(*globalConfig, *projectConfig, arguments)
	terralessData.Config.Validate()

	// for _, terralessProvider := range terralessProviders {
	// 	activeProviders := terralessProvider.FilterActiveProviders(*projectConfig)
	// 	for _, activeProvider := range activeProviders {
	// 		terralessData.ActiveProviders[activeProvider.Name] = activeProvider
	// 	}
	// }

	// logrus.Debug(terralessData.Config.Backend.Name)
	// logrus.Debug(terralessData.Config.Backend.Type)
	// logrus.Debug(terralessData.Config.Backend.Data)
	// logrus.Fatal(terralessData.Config)

	return terralessData
}
