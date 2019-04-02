package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_RenderUploadTemplates(t *testing.T) {
	// given
	buffer := bytes.Buffer{}
	config := schema.TerralessConfig{
		ProjectName: "DummyProjectName",
		Uploads: []schema.TerralessUpload{
			{
				Type: "s3",
				Cloudfront: schema.TerralessCloudfront{

				},
			},
		},
	}

	// when
	buffer = RenderUploadTemplates(config, buffer)

	// then
	assert.Contains(t, buffer.String(), `## Terraless Lambda@Edge`)
	assert.Contains(t, buffer.String(), `resource "aws_cloudwatch_log_group" "lambda-log-terraless-lambda-cloudfront"`)
	assert.Contains(t, buffer.String(), `resource "aws_lambda_function" "terraless-lambda-cloudfront"`)
}
