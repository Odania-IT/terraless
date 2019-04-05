package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func baseDir() string {
	dir, _ := os.Getwd()
	return dir
}

func TestArguments_DetectGlobalConfig(t *testing.T) {
	examplesDir := filepath.Join(baseDir(), "examples")
	_ = os.Chdir(examplesDir)
	configFolders := []string{
		"dummy1",
	}

	configFile := detectGlobalConfig(configFolders)

	assert.Equal(t, path.Join(examplesDir, "terraless.yml"), configFile)
}

func TestArguments_ParseArguments(t *testing.T) {
	arguments := []string{
		"-c",
		"/dummy/test.yml",
		"-e",
		"test",
		"info",
	}

	args, kingpinResult := parseArguments(arguments)

	assert.Equal(t, "info", kingpinResult)
	assert.Equal(t, "/dummy/test.yml", args.Config)
	assert.Equal(t, "test", args.Environment)
}
