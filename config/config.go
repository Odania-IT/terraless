package config

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)


func readProjectYamlConfig(arguments schema.Arguments) *schema.TerralessProjectConfig {
	if arguments.Config == "" {
		logrus.Fatal("Can not load YAML no project config file provided!")
	}

	bytes, err := ioutil.ReadFile(arguments.Config)
	config := &schema.TerralessProjectConfig{}

	if err != nil {
		logrus.Warnf("Could not read project config! Error: %s\n", err)
		return config.Validate()
	}

	if err := yaml.Unmarshal(bytes, config); err != nil {
		logrus.Fatalf("Could not parse project config! Error: %s\n", err)
	}

	logrus.Debugf("Read project YAML config %s\n", arguments.Config)

	return config.Validate()
}

func readGlobalYamlConfig(arguments schema.Arguments) *schema.TerralessGlobalConfig {
	config := &schema.TerralessGlobalConfig{}
	if arguments.GlobalConfig == "" {
		return config
	}

	logrus.Debugf("Loading global config %s\n", arguments.GlobalConfig)
	bytes, err := ioutil.ReadFile(arguments.GlobalConfig)

	if err != nil {
		logrus.Fatalf("Could not read global config! Error: %s\n", err)
	}

	if err := yaml.Unmarshal(bytes, config); err != nil {
		logrus.Fatalf("Could not parse global config! Error: %s\n", err)
	}

	logrus.Debugf("Read global YAML config %s\n", arguments.GlobalConfig)

	return config // config.ProcessConfig(arguments)
}
