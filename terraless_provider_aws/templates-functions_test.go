package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_RenderFunctionTemplates_DoesNotHandleWrongType(t *testing.T) {
	// given
	buffer := bytes.Buffer{}

	// when
	buffer = RenderFunctionTemplates("dummy", schema.FunctionEvents{}, &schema.TerralessData{}, buffer)

	// then
	expected := ``
	assert.Equal(t, expected, buffer.String())
}

func TestTemplatesFunctions_RenderFunctionTemplates_HttpEvent(t *testing.T) {
	// given
	buffer := bytes.Buffer{}
	functionEvents := schema.FunctionEvents{
		Events: map[string][]schema.FunctionEvent{
			"http": {
				{
					FunctionName: "DummyFunction",
					FunctionEvent: schema.TerralessFunctionEvent{
						Type: "http",
						Path: "dummy",
					},
				},
			},
		},
	}
	terralessData := schema.TerralessData{
		Arguments: schema.Arguments{
			Environment: "DummyEnvironment",
		},
		Config: schema.TerralessConfig{
			ProjectName: "DummyProjectName",
		},
	}

	// when
	buffer = RenderFunctionTemplates("aws", functionEvents, &terralessData, buffer)

	// then
	expected := `resource "aws_api_gateway_rest_api" "terraless-api-gateway" {
  name        = "DummyProjectName-DummyEnvironment"
  description = "Terraless Api Gateway for DummyProjectName-DummyEnvironment"
}`
	assert.Contains(t, buffer.String(), expected)
	assert.Contains(t, buffer.String(), `output "api-gateway-invoke-url"`)
	assert.Contains(t, buffer.String(), `resource "aws_api_gateway_resource" "terraless-lambda-DummyFunction-dummy"`)
	assert.Contains(t, buffer.String(), `resource "aws_iam_role" "terraless-lambda-iam-role"`)
	assert.Contains(t, buffer.String(), `resource "aws_cloudwatch_log_group" "lambda-log-DummyFunction"`)
}

func TestTemplatesFunctions_RenderFunctionTemplates_SqsEvent(t *testing.T) {
	// given
	buffer := bytes.Buffer{}
	functionEvents := schema.FunctionEvents{
		Events: map[string][]schema.FunctionEvent{
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
	}
	terralessData := schema.TerralessData{
		Arguments: schema.Arguments{
			Environment: "DummyEnvironment",
		},
		Config: schema.TerralessConfig{
			ProjectName: "DummyProjectName",
		},
	}

	// when
	buffer = RenderFunctionTemplates("aws", functionEvents, &terralessData, buffer)

	// then
	assert.Contains(t, buffer.String(), `## Terraless sqs`)
}
