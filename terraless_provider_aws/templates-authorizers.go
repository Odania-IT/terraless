package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
	"github.com/sirupsen/logrus"
)

func RenderAuthorizerTemplates(config schema.TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
	for _, authorizer := range config.Authorizers {
		if authorizer.Type == "aws" {
			logrus.Debugf("Generating authorizer template for %s\n", authorizer.Name)
			authorizer.TerraformName = "terraless-authorizer-" + support.SanitizeString(authorizer.Name)

			buffer = templates.RenderTemplateToBuffer(authorizer, buffer, awsTemplates("authorizer.tf.tmpl"), "terraless-authorizer")
		}
	}

	return buffer
}
