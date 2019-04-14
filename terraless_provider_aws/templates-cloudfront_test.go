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
	terralessData := schema.TerralessData{
		Config: schema.TerralessConfig{
			ProjectName: "DummyProjectName",
			Uploads: []schema.TerralessUpload{
				{
					Type: "s3",
					Cloudfront: schema.TerralessCloudfront{
						Handler: "singleEntryPointHandler",
					},
				},
			},
		},
	}

	// when
	buffer = RenderUploadTemplates(terralessData, buffer)

	// then
	assert.NotContains(t, buffer.String(), `## Terraless Lambda@Edge`)
	assert.NotContains(t, buffer.String(), `resource "aws_cloudwatch_log_group" "lambda-log-terraless-lambda-cloudfront"`)
	assert.NotContains(t, buffer.String(), `resource "aws_lambda_function" "terraless-lambda-cloudfront"`)
}

func TestTemplatesFunctions_RenderUploadTemplates_WithDomain(t *testing.T) {
	// given
	buffer := bytes.Buffer{}
	terralessData := schema.TerralessData{
		Config: schema.TerralessConfig{
			ProjectName: "DummyProjectName",
			Uploads: []schema.TerralessUpload{
				{
					Type: "s3",
					Cloudfront: schema.TerralessCloudfront{
						Domain: "my-dummy-domain.org",
						Handler: "redirectToWww",
					},
				},
			},
		},
	}

	// when
	buffer = RenderUploadTemplates(terralessData, buffer)

	// then
	assert.Contains(t, buffer.String(), `## Terraless Lambda@Edge`)
	assert.Contains(t, buffer.String(), `resource "aws_cloudwatch_log_group" "lambda-log-terraless-lambda-cloudfront"`)
	assert.Contains(t, buffer.String(), `resource "aws_lambda_function" "terraless-lambda-cloudfront"`)
	assert.Contains(t, buffer.String(), `resource "aws_cloudfront_origin_access_identity" "terraless-default"`)
	assert.Contains(t, buffer.String(), `resource "aws_cloudfront_distribution" "terraless-default"`)
}
