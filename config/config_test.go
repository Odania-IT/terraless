package config

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func baseDir() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "..")
}

func TestTerralessConfig_NewTerralessConfig(t *testing.T) {
	args := schema.Arguments{
		Config: filepath.Join(baseDir(), "examples", "terraless-project.yml"),
	}

	terralessData := NewTerralessData(args, []schema.Provider{})

	assert.Equal(t, &schema.TerralessData{}, terralessData)
}
