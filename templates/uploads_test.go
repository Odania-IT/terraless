package templates

import (
	"github.com/Odania-IT/terraless/dummy"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	provider := dummy.TerralessProvider{}
	provider.Reset()
	providers := []schema.Provider{
		provider,
	}

	// when
	ProcessUploads(terralessData, providers)

	// then
	testProcessed := provider.TestProcessed()
	assert.Equal(t, true, testProcessed["RenderUploadTemplates"])
}
