package templates

import (
	"bytes"
	"github.com/Odania-IT/terraless/dummy"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunctions_ProcessFunctions(t *testing.T) {
	// given
	buffer := bytes.Buffer{}
	terralessData := schema.TerralessData{
		Config: schema.TerralessConfig{
			Functions: map[string]schema.TerralessFunction{
				"DummyFunction": {
					Type:        "aws",
					Handler:     "dummy.handler",
					Description: "My description",
					Events: []schema.TerralessFunctionEvent{
						{
							Type: "http",
							Path: "dummy",
						},
					},
				},
				"OtherType": {
					Type:        "gcloud",
					Handler:     "dummy.handler",
					Description: "My description",
					Events: []schema.TerralessFunctionEvent{
						{
							Type: "http",
							Path: "dummy",
						},
					},
				},
				"SpecificFunction": {
					Type:       "aws",
					Handler:    "specific.Handler",
					Runtime:    "dummyRuntime",
					MemorySize: 512,
					Timeout:    60,
					Events: []schema.TerralessFunctionEvent{
						{
							Type:   "http",
							Method: "GET",
							Path:   "dummy/specific",
						},
						{
							Type: "sqs",
							Arn:  "arn:aws::::sqs",
						},
					},
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
	response := processFunctions(&terralessData, providers, buffer)

	// then
	testProcessed := provider.TestProcessed()
	expected := ""
	assert.Equal(t, expected, response.String())
	assert.Equal(t, true, testProcessed["RenderFunctionTemplates"])
}
