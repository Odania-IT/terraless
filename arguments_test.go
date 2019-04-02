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
