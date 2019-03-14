package templates

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerralessFunctions_ConsolidateEventData(t *testing.T) {
	// given
	terralessData := schema.TerralessData{
		TerralessProviders: []schema.Provider{},
		Config: schema.TerralessConfig{
			Functions: map[string]schema.TerralessFunction{
				"DummyFunction": {
					Type: "aws",
					Handler: "dummy.handler",
					Description: "My description",
					Events: []schema.TerralessFunctionEvent{
						{
							Type: "http",
							Path: "dummy",
						},
					},
				},
				"OtherType": {
					Type: "gcloud",
					Handler: "dummy.handler",
					Description: "My description",
					Events: []schema.TerralessFunctionEvent{
						{
							Type: "http",
							Path: "dummy",
						},
					},
				},
				"SpecificFunction": {
					Type: "aws",
					Handler: "specific.Handler",
					Runtime: "dummyRuntime",
					MemorySize: 512,
					Timeout: 60,
					Events: []schema.TerralessFunctionEvent{
						{
							Type: "http",
							Method: "GET",
							Path: "dummy/specific",
						},
						{
							Type: "sqs",
							Arn: "arn:aws::::sqs",
						},
					},
				},
			},
		},
	}

	// when
	functionEvents := consolidateEventData(terralessData)

	// then
	expected := map[string]schema.FunctionEvents{
		"aws": {
			Events: map[string][]schema.FunctionEvent{
				"http": {
					{
						FunctionName: "DummyFunction",
						FunctionEvent: schema.TerralessFunctionEvent{
							Type: "http",
							Path: "dummy",
						},
					},
					{
						FunctionName: "SpecificFunction",
						FunctionEvent: schema.TerralessFunctionEvent{
							Type: "http",
							Method: "GET",
							Path: "dummy/specific",
						},
					},
				},
				"sqs": {
					{
						FunctionName: "SpecificFunction",
						FunctionEvent: schema.TerralessFunctionEvent{
							Type: "sqs",
							Arn: "arn:aws::::sqs",
						},
					},
				},
			},
		},
		"gcloud": {
			Events: map[string][]schema.FunctionEvent{
				"http": {
					{
						FunctionName: "OtherType",
						FunctionEvent: schema.TerralessFunctionEvent{
							Type: "http",
							Path: "dummy",
						},
					},
				},
			},
		},
	}
	assert.Equal(t, expected, functionEvents)
}
