package main

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const (
	CODENAME = "Flying Eagle"
	VERSION  = "0.1.18"
)

var dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
var (
	app        = kingpin.New("terraless", "Terraless cloud army swiss knife")
	configFlag = app.Flag("config", "Project Config file").
		Short('c').
		Default(path.Join(dir, "terraless-project.yml")).
		String()
	environment      = app.Flag("environment", "Environment").Short('e').Default("develop").String()
	forceDeploy      = app.Flag("force-deploy", "Force deployment (Do not ask the user)").Bool()
	globalConfig     = app.Flag("global-config", "Global Config file").Short('g').String()
	logLevel         = app.Flag("log-level", "Log level").Default("info").String()
	noDeploy         = app.Flag("no-deploy", "Do not execute deploy").Bool()
	noUpload         = app.Flag("no-upload", "Do not upload").Bool()
	terraformCommand = app.Flag("terraform-command", "Terraform Command").Default("terraform").String()

	// Commands
	deployCommand  = app.Command("deploy", "Deploy")
	initCommand    = app.Command("init", "Initialize Templates")
	infoCommand    = app.Command("info", "Display information")
	sessionCommand = app.Command("session", "Handle Provider sessions")
	versionCommand = app.Command("version", "Version")
	uploadCommand  = app.Command("upload", "Upload")

	// Deploy Command Options
	deployNoProviderGeneration = deployCommand.Flag("no-provider-generation", "Do not generate terraform provider").Default("false").Bool()
)

func detectGlobalConfig() string {
	logrus.Info("Trying to detect global config")
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
			logrus.Infof("Found global config: %s\n", fullPath)
			return fullPath
		}
	}

	return ""
}

func parseArguments() (schema.Arguments, string) {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	arguments := &schema.Arguments{
		Config:               *configFlag,
		Environment:          *environment,
		ForceDeploy:          *forceDeploy,
		GlobalConfig:         *globalConfig,
		LogLevel:             *logLevel,
		NoDeploy:             *noDeploy,
		NoProviderGeneration: *deployNoProviderGeneration,
		NoUpload:             *noUpload,
		TerraformCommand:     *terraformCommand,
	}

	if arguments.GlobalConfig == "" {
		arguments.GlobalConfig = detectGlobalConfig()
	}

	// Set log level
	level, _ := logrus.ParseLevel(arguments.LogLevel)
	logrus.SetLevel(level)

	kingpinResult := kingpin.MustParse(app.Parse(os.Args[1:]))

	return *arguments, kingpinResult
}
