package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
	"github.com/sirupsen/logrus"
	"strconv"
)

var lambdaFunctionsTemplate = `
# Lambda Function {{.FunctionName}}

resource "aws_cloudwatch_log_group" "lambda-log-{{.FunctionName}}" {
  name = "/aws/lambda/{{ .ProjectName }}-{{.FunctionName}}"
  retention_in_days = 14
}

resource "aws_lambda_function" "lambda-{{.FunctionName}}" {
  filename = "${data.archive_file.lambda-archive.output_path}"
  function_name = "{{ .ProjectName }}-{{.FunctionName}}"
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

var addTerralessLambdaRole bool
func renderBaseFunction(functionConfig schema.TerralessFunction, functionName string, config schema.TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
	logrus.Debug("Rendering Template for Lambda Function: ", functionName)
	functionConfig.RenderEnvironment = len(functionConfig.Environment) > 0
	functionConfig.FunctionName = functionName
	functionConfig.ProjectName = config.ProjectName

	// Set default runtime if none is specified for the function
	if functionConfig.Runtime == "" {
		functionConfig.Runtime = config.Settings.Runtime
	}

	if functionConfig.RoleArn == "" {
		functionConfig.RoleArn = "${aws_iam_role.terraless-lambda-iam-role.arn}"
		addTerralessLambdaRole = true
	}

	return templates.RenderTemplateToBuffer(functionConfig, buffer, lambdaFunctionsTemplate, "aws-lambda-function")
}

func RenderFunctionTemplates(resourceType string, functionEvents schema.FunctionEvents, terralessData *schema.TerralessData, buffer bytes.Buffer) bytes.Buffer {
	if !canHandle(resourceType) {
		return buffer
	}

	buffer.WriteString("## Terraless Functions AWS\n\n")
	functionsToRender := map[string]bool{}
	for eventType, functionEventArray := range functionEvents.Events {
		baseTemplate := eventTemplates[eventType]

		if baseTemplate == "" {
			baseTemplate = "## Terraless " + eventType + "\n"
		}

		buffer = templates.RenderTemplateToBuffer(terralessData, buffer, baseTemplate, "function-event-" + eventType)

		// Events
		pathsRendered := map[string]string{}
		for key, event := range functionEventArray {
			logrus.Debugf("[EventType %s][AWS %s] Rendering Event %s\n", eventType, event.FunctionName, event)
			functionsToRender[event.FunctionName] = true

			// Render function template
			functionEvent := event.FunctionEvent
			functionEvent.FunctionName = event.FunctionName
			functionEvent.Idx = strconv.FormatInt(int64(key), 10)
			functionEvent.ProjectName = terralessData.Config.ProjectName
			functionEvent.PathsRendered = pathsRendered
			functionEvent.ResourceNameForPath = support.SanitizeString(functionEvent.Path)
			integrationTemplate := functionIntegrationTemplates[functionEvent.Type]

			if integrationTemplate == "" {
				logrus.Fatal("Event Type ", functionEvent.Type, " unknown! Function: ", event.FunctionName)
			}

			buffer = templates.RenderTemplateToBuffer(functionEvent, buffer, integrationTemplate, "function-event-" + functionEvent.Type + "-" + functionEvent.Idx)
			pathsRendered[functionEvent.Path] = support.SanitizeString(functionEvent.Path)
		}
	}

	// Render function base
	for functionName := range functionsToRender {
		functionConfig := terralessData.Config.Functions[functionName]
		buffer = renderBaseFunction(functionConfig, functionName, terralessData.Config, buffer)

		terralessData.Config.Runtimes = append(terralessData.Config.Runtimes, functionConfig.Runtime)
	}

	// AWS Lambda
	if addTerralessLambdaRole {
		buffer = templates.RenderTemplateToBuffer(terralessData, buffer, awsTemplates("iam.tf.tmpl"), "aws-lambda-iam")
	}

	return buffer
}
