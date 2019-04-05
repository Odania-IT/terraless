package templates

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

var functionRendered bool
func dummyTerralessProviderFunctions() schema.Provider {
	return schema.Provider{
		CanHandle: func(resourceType string) bool {
			return resourceType == "dummy"
		},
		PrepareSession: func(terralessConfig schema.TerralessConfig) {
		},
		RenderFunctionTemplates: func(resourceType string, functionEvents schema.FunctionEvents, terralessData *schema.TerralessData, buffer bytes.Buffer) bytes.Buffer {
			functionRendered = true

			return buffer
		},
	}
}

func TestFunctions_ProcessFunctions(t *testing.T) {
	// given
	buffer := bytes.Buffer{}
	terralessData := schema.TerralessData{
		TerralessProviders: []schema.Provider{
			dummyTerralessProviderFunctions(),
		},
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

	// when
	response := processFunctions(&terralessData, buffer)

	// then
	expected := ""
	assert.Equal(t, expected, response.String())
	assert.Equal(t, true, functionRendered)
}
