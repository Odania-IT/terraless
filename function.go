package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
)

const (
	currentEnvironmentTemplate = `#!/usr/bin/env bash
{{ range $key, $val := . }}export {{ $key }}={{ $val }}
{{ end }}
`

	helperTemplate = `#!/usr/bin/env bash
terraless $@
source {{.TerralessDirectory}}/currentEnvironment.sh
`
)

func initHelperFunction(arguments schema.Arguments) {
	logrus.Debug("Verifying helper bash function exists")
	if runtime.GOOS == "windows" {
		logrus.Info("Can currently not create helper function for Windows")
		return
	}

	helperFile := filepath.Join(arguments.TerralessDirectory, "terraless-env")
	buffer := bytes.Buffer{}
	buffer = templates.RenderTemplateToBuffer(arguments, buffer, helperTemplate, "helper-template")

	err := os.MkdirAll(arguments.TerralessDirectory, 0755)
	if err != nil {
		logrus.Fatalf("Could not create plugin directory: %s\n", arguments.TerralessDirectory)
	}

	logrus.Infof("Creating or updating helper file %s\n", helperFile)
	support.WriteToFile(helperFile, buffer)

	os.Chmod(helperFile, 0755)

	writeEnvironmentFile(arguments, map[string]string{})
}

func writeEnvironmentFile(arguments schema.Arguments, environmentVariables map[string]string) {
	buffer := bytes.Buffer{}
	buffer = templates.RenderTemplateToBuffer(environmentVariables, buffer, currentEnvironmentTemplate, "current-environment-template")

	currentEnvironmentFile := filepath.Join(arguments.TerralessDirectory, "currentEnvironment.sh")
	logrus.Debugf("Writting environment to %s file\n", currentEnvironmentFile)
	support.WriteToFile(currentEnvironmentFile, buffer)
}
