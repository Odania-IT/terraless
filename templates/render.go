package templates

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

var mainTfTemplate = `# This file is generated by Terraless

{{ range .Config.Providers }}
# {{ .Name }}
provider "{{.Type}}" {
  {{ range $key, $value := awsProviderKeys .Data .Name }}{{$key}} = "{{$value}}"
  {{ end }}
}
{{ end }}

{{ if .Config.Backend.Type }}
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

var lambdaPackageTemplate = `## Terraless: Lambda Packaga

data "archive_file" "lambda-archive" {
  source_dir = "{{.Config.Package.SourceDir}}"

  output_path = "{{.Config.Package.OutputPath}}"
  type = "zip"
}

`

func Render(terralessData schema.TerralessData, buffer bytes.Buffer) bytes.Buffer {
	config := terralessData.Config
	renderTemplate(terralessData, filepath.Join(config.SourcePath, "terraless-main.tf"), mainTfTemplate)

	if len(config.Authorizers) > 0 {
		logrus.Debug("Creating authorizer templates")

		for _, terralessProvider := range terralessData.TerralessProviders {
			buffer = terralessProvider.RenderAuthorizerTemplates(config, buffer)
		}
	}

	if len(terralessData.Config.Functions) > 0 {
		logrus.Debug("Creating function templates")
		buffer = processFunctions(terralessData, buffer)
	}

	if config.Package.SourceDir != "" {
		logrus.Debug("Creating package template")
		terralessData.Config.Package.OutputPath = filepath.Join(config.TargetPath, "lambda.zip")

		buffer = RenderTemplateToBuffer(terralessData, buffer, lambdaPackageTemplate, "terraless-package")
	}

	if len(config.Certificates) > 0 {
		logrus.Debug("Creating certificate templates")

		for _, terralessProvider := range terralessData.TerralessProviders {
			buffer = terralessProvider.RenderCertificateTemplates(config, buffer)
		}
	}

	if len(config.Endpoints) > 0 {
		logrus.Debug("Creating endpoint templates")

		for _, terralessProvider := range terralessData.TerralessProviders {
			buffer = terralessProvider.RenderEndpointTemplates(config, buffer)
		}
	}

	if len(config.Uploads) > 0 {
		logrus.Debug("Creating cloudfront templates")
		for _, terralessProvider := range terralessData.TerralessProviders {
			buffer = terralessProvider.RenderUploadTemplates(config, buffer)
		}
	}

	return buffer
}

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

func RenderTemplateToBuffer(config interface{}, buffer bytes.Buffer, tpl string, name string) bytes.Buffer {
	tmpl := template.Must(template.New(name).Funcs(template.FuncMap{
		"resourceForPathRendered": func(pathRendered map[string]string, key string) bool {
			return pathRendered[key] == ""
		},
		"terralessResourceName": func(pathRendered map[string]string, key string) string {
			return pathRendered[key]
		},
		"currentTime": func() string {
			currentTime := time.Now()

			return currentTime.Format("2006-01-02 15:04:05")
		},
	}).Parse(tpl))
	err := tmpl.Execute(&buffer, config)

	if err != nil {
		logrus.Fatal("Failed writing to Buffer: ", err)
	}

	return buffer
}
