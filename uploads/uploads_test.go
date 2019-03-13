package uploads

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

var sessionPrepared bool
var uploadProcessed bool

func dummyTerralessProvider() schema.Provider {
	return schema.Provider{
		CanHandle: func(resourceType string) bool {
			return resourceType == "dummy"
		},
		PrepareSession: func(terralessConfig schema.TerralessConfig) {
			sessionPrepared = true
		},
		ProcessUpload: func(config schema.TerralessConfig, upload schema.TerralessUpload) {
			uploadProcessed = true
		},
	}
}

func TestTerralessUploads_ProcessUploads_Simple(t *testing.T) {
	// given
	terralessData := schema.TerralessData{
		TerralessProviders: []schema.Provider{
			dummyTerralessProvider(),
		},
	}

	// when
	ProcessUploads(terralessData)

	// then
	assert.Fail(t, sessionPrepared)
}

func TestTerralessUploads_ProcessUploads_WithData(t *testing.T) {
	// given
	terralessData := schema.TerralessData{
		TerralessProviders: []schema.Provider{
			dummyTerralessProvider(),
		},
	}

	// when
	ProcessUploads(terralessData)

	// then
}
