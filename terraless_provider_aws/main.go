package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/templates"
	"github.com/gobuffalo/packr"
	"github.com/sirupsen/logrus"
)

func awsTemplates(name string) string {
	pckr := packr.NewBox("./templates")

	tpl, err := pckr.FindString(name)
	if err != nil {
		logrus.Fatal("Failed retrieving template: ", err)
	}

	return tpl
}

func Provider() schema.Provider {
	return schema.Provider{
		CanHandle:                  canHandle,
		FinalizeTemplates:          finalizeTemplates,
		Name:                       providerName,
		PrepareSession:             prepareSession,
		ProcessUpload:              processUpload,
		RenderAuthorizerTemplates:  RenderAuthorizerTemplates,
		RenderCertificateTemplates: RenderCertificateTemplates,
		RenderEndpointTemplates:    RenderEndpointTemplates,
		RenderFunctionTemplates:    RenderFunctionTemplates,
		RenderUploadTemplates:      RenderUploadTemplates,
	}
}

func canHandle(resourceType string) bool {
	return resourceType == "aws"
}

func finalizeTemplates(terralessData schema.TerralessData, buffer bytes.Buffer) bytes.Buffer {
	if addTerralessLambdaRole {
		buffer = templates.RenderTemplateToBuffer(terralessData, buffer, awsTemplates("iam.tf.tmpl"), "aws-lambda-iam")
	}

	return buffer
}

func providerName() string {
	return "terraless-provider-aws"
}
