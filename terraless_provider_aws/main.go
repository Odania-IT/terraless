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
		// FilterActiveProviders: filterActiveProviders,
		PrepareSession: prepareSession,
	}
}

func canHandle(resourceType string) bool {
	if resourceType == "aws" {
		return true
	}

	return false
}

// func filterActiveProviders(terralessConfig schema.TerralessConfig) map[string]schema.TerralessProvider {
// 	var result = map[string]schema.TerralessProvider{}
//
// 	var backendProfile = "none-none"
// 	if terralessConfig.Backend.Type == "s3" {
// 		logrus.Debug("Checking profile for terraform backend s3")
// 		profile := terralessConfig.Backend.Data["profile"]
//
// 		if profile != "" && result[profile].Name == "" {
// 			backendProfile = profile
// 		}
// 	}
//
// 	parts := strings.Split(backendProfile, "-")
// 	role := parts[len(parts)-1]
// 	globalName := strings.Join(parts[:len(parts)-1], "-")
//
// 	for _, team := range terralessConfig.Teams {
// 		for _, provider := range team.Providers {
// 			result[provider.Name] = provider
//
//
// 		}
// 	}
// 	logrus.Fatal(result)
//
// 	return result
// }
