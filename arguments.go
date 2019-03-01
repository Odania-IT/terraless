package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const (
	VERSION = "0.1.0"
)

type Arguments struct {
	Config string
	Environment string
	ForceDeploy bool
	GlobalConfig string
	LogLevel string
	TerraformCommand string
}

func detectGlobalConfig() string {
	configFolders := []string{
		".terraless",
		filepath.Join(".config", ".terraless"),
		".aws",
		filepath.Join(".config", "gcloud"),
	}
	usr, err := user.Current()

	if err != nil {
		logrus.Fatal("Could not detect user home folder")
	}

	for _, folder := range configFolders {
		fullPath := filepath.Join(usr.HomeDir, folder, "terraless.yml")

		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	return ""
}

func parseArguments() Arguments {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	arguments := &Arguments{
		Config: path.Join(dir, "terraless.yml"),
		Environment: "develop",
		TerraformCommand: "terraform",
	}
	flag.StringVar(&arguments.Config, "config", arguments.Config, "Config file")
	flag.StringVar(&arguments.Environment, "environment", arguments.Environment, "Environment")
	flag.BoolVar(&arguments.ForceDeploy, "force-deploy", arguments.ForceDeploy, "Force deployment (Do not ask the user)")
	flag.StringVar(&arguments.GlobalConfig, "global-config", arguments.Config, "Global config file")
	flag.StringVar(&arguments.LogLevel, "log-level", arguments.LogLevel, "Log level")
	flag.Usage = flagUsage
	flag.Parse()

	return setArgumentDefaults(*arguments)
}

func flagUsage() {
	_, _ = fmt.Fprintf(os.Stderr, "Version of %s: %s\n", os.Args[0], VERSION)
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func setArgumentDefaults(arguments Arguments) Arguments {
	if arguments.LogLevel == "" {
		arguments.LogLevel = "INFO"
	}
	level, _ := logrus.ParseLevel(arguments.LogLevel)
	logrus.SetLevel(level)

	if arguments.GlobalConfig == "" {
		logrus.Debug("Global config not specified! Trying to detect it...")
		arguments.GlobalConfig = detectGlobalConfig()
	}

	return arguments
}
