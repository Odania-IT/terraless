package templates

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

var uploadProcessed bool

func dummyTerralessProvider() schema.Provider {
	return schema.Provider{
		CanHandle: func(resourceType string) bool {
			return resourceType == "dummy"
		},
		PrepareSession: func(terralessConfig schema.TerralessConfig) {
		},
		ProcessUpload: func(terralessData schema.TerralessData, upload schema.TerralessUpload) []string {
			uploadProcessed = true

			return []string{}
		},
	}
}

func TestTerralessUploads_ProcessUploads(t *testing.T) {
	// given
	terralessData := schema.TerralessData{
		Config: schema.TerralessConfig{
			Uploads: []schema.TerralessUpload{
				{
					Bucket:      "myBucket",
					Provider:    "myProvider",
					ProjectName: "myProject",
					Region:      "myRegion",
					Source:      "mySource",
					Target:      "myTarget",
					Type:        "dummy",
				},
			},
		},
	}
	providers := []schema.Provider{
		dummyTerralessProvider(),
	}

	// when
	ProcessUploads(terralessData, providers)

	// then
	assert.Equal(t, true, uploadProcessed)
}
