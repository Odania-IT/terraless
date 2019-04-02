package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_RenderEndpointTemplates(t *testing.T) {
	// given
	buffer := bytes.Buffer{}
	config := schema.TerralessConfig{
			ProjectName: "DummyProjectName",
			Endpoints: []schema.TerralessEndpoint{
				{
					Type: "dummy",
					Domain: "my-secret-dummy-domain.org",
				},
				{
					Type: "apigateway",
					Domain: "my-secret-domain.org",
				},
			},
		}

	// when
	buffer = RenderEndpointTemplates(config, buffer)

	// then
	assert.Contains(t, buffer.String(), `domain_name     = "my-secret-domain.org"`)
	assert.Contains(t, buffer.String(), `resource "aws_api_gateway_base_path_mapping" "terraless-endpoint-my-secret-domain-org"`)
	assert.Contains(t, buffer.String(), `resource "aws_route53_record" "terraless-endpoint-my-secret-domain-org"`)
}
