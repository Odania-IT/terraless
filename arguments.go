package main

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path"
	"path/filepath"
)

const (
	CODENAME = "Flying Eagle"
	VERSION  = "0.1.36"
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
	pluginDir        = app.Flag("plugin-directory", "Plugin Directory").String()
	terraformCommand = app.Flag("terraform-command", "Terraform Command").Default("terraform").String()
	variables        = app.Flag("var", "Variable").StringMap()

	// Commands
	authCommand          = app.Command("auth", "Authenticate with providers")
	deployCommand        = app.Command("deploy", "Deploy")
	initCommand          = app.Command("init", "Initialize Terraless")
	initTemplatesCommand = app.Command("init-templates", "Initialize Templates")
	infoCommand          = app.Command("info", "Display information")
	uploadCommand        = app.Command("upload", "Upload")
	versionCommand       = app.Command("version", "Version")

	// Auth Command Options
	authProviderName = authCommand.Flag("provider", "Provider to authenticate with").String()

	// Deploy Command Options
	deployNoProviderGeneration = deployCommand.Flag("no-provider-generation", "Do not generate terraform provider").Default("false").Bool()
)

func detectGlobalConfig(configFolders []string) string {
	logrus.Info("Trying to detect global config")

	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		logrus.Fatal("Could not detect current directory")
	}

	fullPath := filepath.Join(currentWorkingDirectory, "terraless.yml")
	if _, err := os.Stat(fullPath); err == nil {
		logrus.Infof("Found global config: %s\n", fullPath)
		return fullPath
	}

	homeDirectory := support.HomeDirectory()
	for _, folder := range configFolders {
		fullPath := filepath.Join(homeDirectory, folder, "terraless.yml")

		if _, err := os.Stat(fullPath); err == nil {
			logrus.Infof("Found global config: %s\n", fullPath)
			return fullPath
		}
	}

	return ""
}

func parseArguments(args []string) (schema.Arguments, string) {
	kingpin.MustParse(app.Parse(args))

	arguments := &schema.Arguments{
		Config:               *configFlag,
		Environment:          *environment,
		ForceDeploy:          *forceDeploy,
		GlobalConfig:         *globalConfig,
		LogLevel:             *logLevel,
		PluginDirectory:      *pluginDir,
		NoDeploy:             *noDeploy,
		NoProviderGeneration: *deployNoProviderGeneration,
		NoUpload:             *noUpload,
		TerraformCommand:     *terraformCommand,
		Variables:            *variables,
	}

	if arguments.GlobalConfig == "" {
		configFolders := []string{
			".terraless",
			filepath.Join(".config", ".terraless"),
			".aws",
			filepath.Join(".config", "gcloud"),
		}
		arguments.GlobalConfig = detectGlobalConfig(configFolders)
	}

	// Set log level
	level, _ := logrus.ParseLevel(arguments.LogLevel)
	logrus.SetLevel(level)

	// Set plugin directory
	if arguments.PluginDirectory == "" {
		homeDirectory := support.HomeDirectory()
		arguments.PluginDirectory = filepath.Join(homeDirectory, ".terraless", "plugins")
	}

	kingpinResult := kingpin.MustParse(app.Parse(args))

	return *arguments, kingpinResult
}
