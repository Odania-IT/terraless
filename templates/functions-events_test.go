package templates

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunctionsEvents_ConsolidateEventData(t *testing.T) {
	// given
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

	// when
	functionEvents := consolidateEventData(terralessData)

	// then
	expectedGcloud := map[string][]schema.FunctionEvent{
		"http": {
			{
				FunctionName: "OtherType",
				FunctionEvent: schema.TerralessFunctionEvent{
					Type: "http",
					Path: "dummy",
				},
			},
		},
	}
	assert.Equal(t, 2, len(functionEvents))
	assert.Equal(t, expectedGcloud, functionEvents["gcloud"].Events)

	expectedAwsSqs := []schema.FunctionEvent{
		{
			FunctionName: "SpecificFunction",
			FunctionEvent: schema.TerralessFunctionEvent{
				Type: "sqs",
				Arn:  "arn:aws::::sqs",
			},
		},
	}
	assert.Equal(t, expectedAwsSqs, functionEvents["aws"].Events["sqs"])
	expectedAwsHttp1 := schema.FunctionEvent{
		FunctionName: "DummyFunction",
		FunctionEvent: schema.TerralessFunctionEvent{
			Type: "http",
			Path: "dummy",
		},
	}
	expectedAwsHttp2 := schema.FunctionEvent{
		FunctionName: "SpecificFunction",
		FunctionEvent: schema.TerralessFunctionEvent{
			Type:   "http",
			Method: "GET",
			Path:   "dummy/specific",
		},
	}

	awsHttpEvents := functionEvents["aws"].Events["http"]
	if awsHttpEvents[0].FunctionName == "DummyFunction" {
		assert.Equal(t, expectedAwsHttp1, awsHttpEvents[0])
		assert.Equal(t, expectedAwsHttp2, awsHttpEvents[1])
	} else {
		assert.Equal(t, expectedAwsHttp2, awsHttpEvents[0])
		assert.Equal(t, expectedAwsHttp1, awsHttpEvents[1])
	}
}
