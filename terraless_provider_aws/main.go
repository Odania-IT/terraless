package terraless_provider_aws

import (
	"github.com/Odania-IT/terraless/schema"
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
		CanHandle: canHandle,
		PrepareSession: prepareSession,
		ProcessUpload: processUpload,
	}
}

func canHandle(resourceType string) bool {
	if resourceType == "aws" {
		return true
	}

	return false
}
