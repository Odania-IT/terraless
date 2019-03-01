package main

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type TerralessFunction struct {
	Description string `yaml:"Description"`
	Environment map[string]string `yaml:"Environment"`
	Handler string `yaml:"Handler"`
	MemorySize int `yaml:"MemorySize"`
	RoleArn string `yaml:"RoleArn"`
	Runtime string `yaml:"Runtime"`
	Timeout int `yaml:"Timeout"`

	// only for rendering template
	FunctionName string
	RenderEnvironment bool
}

type TerralessPackage struct {
	SourceDir string `yaml:"SourceDir"`
	OutputPath string
}

type TerralessType struct {
	Data map[string]string `yaml:"Data"`
	Name string `yaml:"Name"`
	Type string `yaml:"Type"`
}

type TerralessConfig struct {
	Backend TerralessType `yaml:"Backend"`
	Functions map[string]TerralessFunction `yaml:"Functions"`
	HasAwsProvider bool
	Package TerralessPackage `yaml:"Package"`
	ProjectName string `yaml:"ProjectName"`
	Providers []TerralessType `yaml:"Providers"`
	SourcePath string `yaml:"SourcePath"`
	Swamp []string `yaml:"Swamp"`
	TargetPath string `yaml:"TargetPath"`
}

func processString(check string, arguments Arguments) string {
	r, _ := regexp.Compile("\\$\\{[a-z]+\\}")
	matches := r.FindStringSubmatch(check)

	if len(matches) > 0 {
		for _, match := range matches {
			if "${environment}" == match {
				check = strings.Replace(check, match, arguments.Environment, -1)
			} else {
				logrus.Fatal("Not implemented! Currently only environment can be replaced!")
			}
		}
	}

	return check
}

func processData(data map[string]string, arguments Arguments) map[string]string {
	for key, val := range data {
		data[key] = processString(val, arguments)
	}

	return data
}

func processConfig(config TerralessConfig, arguments Arguments) TerralessConfig {
	processedConfig := &TerralessConfig{
		Backend: config.Backend,
		Functions: config.Functions,
		Package: config.Package,
		ProjectName: processString(config.ProjectName, arguments),
	}

	processedConfig.Backend.Name = processString(processedConfig.Backend.Name, arguments)
	processedConfig.Backend.Data = processData(processedConfig.Backend.Data, arguments)

	for _, provider := range config.Providers {
		provider.Name = processString(provider.Name, arguments)
		provider.Data = processData(provider.Data, arguments)

		processedConfig.Providers = append(processedConfig.Providers, provider)
	}

	return *processedConfig
}

func readYamlConfig(file string, arguments Arguments) TerralessConfig {
	bytes, err := ioutil.ReadFile(file)
	config := &TerralessConfig{}

	if err != nil {
		logrus.Fatal("Could not read config", err)
	}

	if err := yaml.Unmarshal(bytes, config); err != nil {
		logrus.Fatal("Could not parse config", err)
	}

	return processConfig(*config, arguments)
}

func newTerralessConfig(arguments Arguments) TerralessConfig {
	config := readYamlConfig(arguments.Config, arguments)

	// Set ProjectName if none is in the config
	if config.ProjectName == "" {
		config.ProjectName = filepath.Base(filepath.Dir(arguments.Config))
	}

	// Set SourcePath if none is in the config
	if config.SourcePath == "" {
		config.SourcePath = filepath.Join(filepath.Dir(arguments.Config))
	}

	// Set TargetPath if none is in the config
	if config.TargetPath == "" {
		config.TargetPath = filepath.Join(filepath.Dir(arguments.Config), ".terraless")
	}

	if arguments.GlobalConfig == "" {
		validateConfig(&config)
		return config
	}

	logrus.Debug("Building configuration for project")
	globalConfig := readYamlConfig(arguments.GlobalConfig, arguments)
	result := buildFinalConfig(config, globalConfig)
	validateConfig(&result)

	return result
}

func findAndEnrichByProviderByName(config TerralessConfig, src TerralessType) TerralessType {
	for _, provider := range config.Providers {
		if provider.Name == src.Name {
			provider.Data = enrichWithData(provider.Data, src.Data)
			return provider
		}
	}

	logrus.Fatal("Could not find provider:", src.Name)

	// Why is this required??
	return TerralessType{}
}

func enrichWithData(data map[string]string, override map[string]string) map[string]string {
	for key, val := range override {
		data[key] = val
	}

	return data
}

func buildFinalConfig(config TerralessConfig, globalConfig TerralessConfig) TerralessConfig {
	result  := &TerralessConfig{
		Functions: config.Functions,
		Package: config.Package,
		ProjectName: config.ProjectName,
		SourcePath: config.SourcePath,
		Swamp: config.Swamp,
		TargetPath: config.TargetPath,
	}

	for _, provider := range config.Providers {
		if provider.Type == "global" {
			result.Providers = append(result.Providers, findAndEnrichByProviderByName(globalConfig, provider))
		} else {
			result.Providers = append(result.Providers, provider)
		}
	}

	if config.Backend.Type == "global" {
		result.Backend = globalConfig.Backend
		globalConfig.Backend.Data = enrichWithData(globalConfig.Backend.Data, config.Backend.Data)
	} else {
		result.Backend = config.Backend
	}

	return *result
}

func validateConfig(config *TerralessConfig) {
	logrus.Debug("Verifying config", config)
	for _, provider := range config.Providers {
		if provider.Type == "global" {
			logrus.Fatal("Unresolved global in provider found!", provider)
		}

		if provider.Type == "aws" {
			config.HasAwsProvider = true
		}
	}

	if config.Backend.Type == "global" {
		logrus.Fatal("Unresolved global in backend found!", config.Backend.Type)
	}
}
