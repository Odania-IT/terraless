package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_RenderAuthorizerTemplates(t *testing.T) {
	// given
	buffer := bytes.Buffer{}
	config := schema.TerralessConfig{
		ProjectName: "DummyProjectName",
		Authorizers: map[string]schema.TerralessAuthorizer{
			"unsed": {
				Type: "dummy",
				Name: "UnusupportedAuthorizer",
			},
			"dummy": {
				Type: "aws",
				Name: "SupportedAuthorizer",
				ProviderArns: []string{
					"arn1",
					"arn2",
				},
			},
		},
	}

	// when
	buffer = RenderAuthorizerTemplates(config, buffer)

	// then
	assert.Contains(t, buffer.String(), `## Terraless Authorizer`)
	assert.Contains(t, buffer.String(), `resource "aws_api_gateway_authorizer" "terraless-authorizer-SupportedAuthorizer"`)
	assert.Contains(t, buffer.String(), `"arn1"`)
	assert.Contains(t, buffer.String(), `"arn2"`)
}
