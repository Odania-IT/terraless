package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/templates"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var lambdaFunctionsTemplate = `
# Lambda Function {{.FunctionName}}

resource "aws_cloudwatch_log_group" "lambda-log-{{.FunctionName}}" {
  name = "/aws/lambda/{{.FunctionName}}"
  retention_in_days = 14
}

resource "aws_lambda_function" "lambda-{{.FunctionName}}" {
  filename = "${data.archive_file.lambda-archive.output_path}"
  function_name = "{{.FunctionName}}"
  role = "{{.RoleArn}}"
  handler = "{{.Handler}}"
  source_code_hash = "${data.archive_file.lambda-archive.output_base64sha256}"
  runtime = "{{.Runtime}}"

  {{ if .RenderEnvironment }}
  environment {
    variables = {
      {{ range $key, $val := .Environment }}{{ $key }} = "{{ $val }}"
      {{ end }}
    }
  }
  {{ end }}
}
`

func renderFunctionTemplate(terralessData schema.TerralessData, targetFileName string, tpl string) {
	config := terralessData.Config
	targetFile, err := os.Create(targetFileName)
	if err != nil {
		logrus.Fatal("Failed creating file: ", filepath.Base(targetFileName), err)
	}

	_, err = targetFile.WriteString("# This file is generated by Terraless")
	if err != nil {
		logrus.Fatal("Failed writing File: ", filepath.Base(targetFileName), err)
	}

	var addTerralessLambdaRole bool
	eventConfig := map[string]bytes.Buffer{}
	for functionName, functionConfig := range config.Functions {
		logrus.Debug("Rendering Template for Lambda Function: ", functionName)
		functionConfig.RenderEnvironment = len(functionConfig.Environment) > 0
		functionConfig.FunctionName = functionName

		// Set default runtime if none is specified for the function
		if functionConfig.Runtime == "" {
			functionConfig.Runtime = terralessData.Config.Settings.Runtime
		}

		if functionConfig.RoleArn == "" {
			functionConfig.RoleArn = "${aws_iam_role.terraless-lambda-iam-role.arn}"
			addTerralessLambdaRole = true
		}

		tmpl := template.Must(template.New("lambda-function.tf").Parse(tpl))
		err = tmpl.Execute(targetFile, functionConfig)

		if err != nil {
			logrus.Fatal("Failed writing File: ", filepath.Base(targetFileName), err)
		}

		for _, event := range functionConfig.Events {
			event.FunctionName = functionName
			integrationTemplate := IntegrationTemplates[event.Type]

			if integrationTemplate == "" {
				logrus.Fatal("Event Type ", event.Type, " unknown! Function: ", functionName)
			}

			eventConfig[event.Type] = templates.RenderTemplateToBuffer(event, eventConfig[event.Type], integrationTemplate)
		}

		config.Runtimes = append(config.Runtimes, functionConfig.Runtime)
	}

	if addTerralessLambdaRole {
		renderTemplate(terralessData, filepath.Join(config.SourcePath, "terraless-lambda-iam.tf"), awsTemplates("iam.tf.tmpl"))
	}

	// Generate event specific files
	writeEventConfigs(terralessData, eventConfig)
}

func writeEventConfigs(terralessData schema.TerralessData, eventConfig map[string]bytes.Buffer) {
	config := terralessData.Config

	for eventType, buffer := range eventConfig {
		if buffer.Len() > 0 {
			baseTemplate := eventTemplates[eventType]

			if baseTemplate == "" {
				baseTemplate = "# This file is generated by Terraless\n"
			}

			targetFileName := filepath.Join(config.SourcePath, "terraless-function-"+eventType+".tf")
			targetFile, err := os.Create(targetFileName)
			if err != nil {
				logrus.Fatal("Failed creating file: ", filepath.Base(targetFileName), err)
			}

			tmpl := template.Must(template.New("terraless-function-"+eventType+".tf").Parse(baseTemplate))
			err = tmpl.Execute(targetFile, terralessData)
			if err != nil {
				logrus.Fatal("Failed writing function base template to file: ", filepath.Base(targetFileName), err)
			}

			_, err = targetFile.Write(buffer.Bytes())
			if err != nil {
				logrus.Fatal("Failed writing to file: ", filepath.Base(targetFileName), err)
			}

			_ = targetFile.Close()
		}
	}
}
