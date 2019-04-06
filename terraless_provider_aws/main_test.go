package terraless_provider_aws

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
