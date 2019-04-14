package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_Provider(t *testing.T) {
	// given

	// when
	provider := Provider()

	// then
	assert.Equal(t, true, provider.CanHandle("aws"))
	assert.Equal(t, false, provider.CanHandle("aws2"))
	assert.Equal(t, false, provider.CanHandle("dummy"))
	assert.Equal(t, "terraless-provider-aws", provider.Name())
}

func TestTemplatesFunctions_CanHandle(t *testing.T) {
	// given

	// when

	// then
	assert.Equal(t, true, canHandle("aws"))
	assert.Equal(t, false, canHandle("aws2"))
	assert.Equal(t, false, canHandle("dummy"))
}

func TestTemplatesFunctions_ProviderName(t *testing.T) {
	// given

	// when

	// then
	assert.Equal(t, "terraless-provider-aws", providerName())
}

func TestTemplatesFunctions_AwsTemplates(t *testing.T) {
	// given

	// when
	template := awsTemplates("iam.tf.tmpl")

	// then
	assert.Contains(t, template, `resource "aws_iam_role" "terraless-lambda-iam-role"`)
}

func TestTemplatesFunctions_FinalizeTemplates(t *testing.T) {
	// given
	addTerralessLambdaRole = true
	terralessData := schema.TerralessData{
		Arguments: schema.Arguments{
			Environment: "DummyEnvironment",
		},
		Config: schema.TerralessConfig{
			ProjectName: "DummyProjectName",
		},
	}

	// when
	buffer := finalizeTemplates(terralessData, bytes.Buffer{})

	// then
	assert.Contains(t, buffer.String(), `resource "aws_iam_role" "terraless-lambda-iam-role"`)
}
