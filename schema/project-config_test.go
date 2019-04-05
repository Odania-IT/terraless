package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerralessConfig_ProjectConfig_Validate(t *testing.T) {
	// given
	projectConfig := TerralessProjectConfig{
		Functions: map[string]TerralessFunction{
			"f1": {
				FunctionName: "f1",
				Runtime: "dummy-runtime",
			},
			"f2": {
				FunctionName: "f2",
			},
		},
	}

	// when
	validatedConfig := projectConfig.Validate()

	// then
	assert.Equal(t, 2, len(validatedConfig.Functions))
	assert.Equal(t, "dummy-runtime", validatedConfig.Functions["f1"].Runtime)
	assert.Equal(t, "", validatedConfig.Functions["f2"].Runtime)
	assert.Equal(t, "ruby2.5", validatedConfig.Settings.Runtime)
}

func TestTerralessConfig_ProjectConfig_ValidateKeepsDefaultRuntime(t *testing.T) {
	// given
	projectConfig := TerralessProjectConfig{
		Functions: map[string]TerralessFunction{
			"f2": {
				FunctionName: "f2",
			},
		},
		Settings: TerralessSettings{
			Runtime: "dummy-runtime",
		},
	}

	// when
	validatedConfig := projectConfig.Validate()

	// then
	assert.Equal(t, 1, len(validatedConfig.Functions))
	assert.Equal(t, "", validatedConfig.Functions["f2"].Runtime)
	assert.Equal(t, "dummy-runtime", validatedConfig.Settings.Runtime)
}
