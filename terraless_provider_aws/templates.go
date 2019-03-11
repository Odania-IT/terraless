package terraless_provider_aws

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var mainTfTemplate = `# This file is generated by Terraless

{{ range .Config.Providers }}
# {{ .Name }}
provider "{{.Type}}" {
  {{ range $key, $value := awsProviderKeys .Data .Name }}{{$key}} = "{{$value}}"
  {{ end }}
}
{{ end }}

{{ if .Config.Backend }}
terraform {
	backend "{{ .Config.Backend.Type }}" {
		{{ range $key, $value := .Config.Backend.Data }}{{$key}} = "{{$value}}"
		{{ end }}
	}
}
{{ end }}

{{if providerPresent .Config.HasProvider "aws" }}
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
{{ end }}
`

var lambdaPackageTemplate = `# This file is generated by Terraless

data "archive_file" "lambda-archive" {
  source_dir = "{{.Config.Package.SourceDir}}"

  output_path = "{{.Config.Package.OutputPath}}"
  type = "zip"
}

`

func renderTemplate(terralessData schema.TerralessData, targetFileName string, tpl string) {
	targetFile, err := os.Create(targetFileName)
	if err != nil {
		logrus.Fatal("Failed creating file: ", filepath.Base(targetFileName), err)
	}

	tmpl := template.Must(template.New(filepath.Base(targetFileName)).Funcs(template.FuncMap{
		"awsProviderKeys": func(data map[string]string, profileName string) map[string]string {
			result := map[string]string{}

			if data["alias"] != "" {
				result["alias"] = data["alias"]
			}

			if data["profile"] != "" {
				result["profile"] = data["profile"]
			}

			if data["region"] != "" {
				result["region"] = data["region"]
			}

			return result
		},
		"providerPresent": func(hasProvider map[string]bool, key string) bool {
			return hasProvider[key]
		},
	}).Parse(tpl))
	err = tmpl.Execute(targetFile, terralessData)

	if err != nil {
		logrus.Fatal("Failed writing File: ", filepath.Base(targetFileName), err)
	}
}

func RenderTemplates(terralessData schema.TerralessData) {
	config := terralessData.Config
	renderTemplate(terralessData, filepath.Join(config.SourcePath, "terraless-main.tf"), mainTfTemplate)

	if config.Package.SourceDir != "" {
		logrus.Debug("Creating package template")
		terralessData.Config.Package.OutputPath = filepath.Join(config.TargetPath, "lambda.zip")
		renderTemplate(terralessData, filepath.Join(config.SourcePath, "terraless-package.tf"), lambdaPackageTemplate)
	}

	if len(config.Functions) > 0 {
		logrus.Debug("Creating function templates")
		renderFunctionTemplate(terralessData, filepath.Join(config.SourcePath, "terraless-functions.tf"), lambdaFunctionsTemplate)
	}

	if len(config.Certificates) > 0 {
		logrus.Debug("Creating certificate templates")
		renderCertificateTemplates(config, filepath.Join(config.SourcePath, "terraless-certificate.tf"))
	}

	if len(config.Uploads) > 0 {
		logrus.Debug("Creating cloudfront templates")
		renderCloudfrontTemplates(config, filepath.Join(config.SourcePath, "terraless-cloudfront.tf"))
	}
}
