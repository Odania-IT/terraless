package terraless_provider_aws

var eventTemplates = map[string]string{
	"http": `## Terraless Function Event HTTP

resource "aws_api_gateway_rest_api" "terraless-api-gateway" {
  name        = "{{ .Config.ProjectName }}-{{ .Arguments.Environment }}"
  description = "Terraless Api Gateway for {{ .Config.ProjectName }}-{{ .Arguments.Environment }}"
}

resource "aws_api_gateway_deployment" "terraless-api-gateway-v1" {
  rest_api_id = "${aws_api_gateway_rest_api.terraless-api-gateway.id}"
  stage_name = "v1"
  stage_description = "Deployed at {{ currentTime }}"
}

output "api-gateway-invoke-url" {
  value = "${aws_api_gateway_deployment.terraless-api-gateway-v1.invoke_url}"
}

`,
}

var functionIntegrationTemplates = map[string]string{
	"http": `
# Function {{ .ProjectName }} {{ .FunctionName }} EventKey: {{.Idx}}

{{ if resourceForPathRendered .PathsRendered .Path }}
resource "aws_api_gateway_resource" "terraless-lambda-{{.FunctionName}}-{{.ResourceNameForPath}}" {
  rest_api_id = "${aws_api_gateway_rest_api.terraless-api-gateway.id}"
  parent_id   = "${aws_api_gateway_rest_api.terraless-api-gateway.root_resource_id}"
  path_part   = "{{ .Path }}"
}
{{ end }}

resource "aws_api_gateway_method" "terraless-lambda-{{.FunctionName}}-{{.Idx}}" {
  rest_api_id   = "${aws_api_gateway_rest_api.terraless-api-gateway.id}"
  {{ if stringNotEmpty .Path }}
  resource_id   = "${aws_api_gateway_resource.terraless-lambda-{{.FunctionName}}-{{.ResourceNameForPath}}.id}"
  {{ else }}
  resource_id   = "${aws_api_gateway_rest_api.terraless-api-gateway.root_resource_id}"
  {{ end }}
  http_method   = "{{ .Method }}"
  authorization = "{{ .Authorization }}"
  authorizer_id = "{{ .AuthorizerId }}"
}

resource "aws_api_gateway_integration" "terraless-lambda-{{.FunctionName}}-{{.Idx}}" {
  {{ if stringNotEmpty .Path }}
  depends_on = ["aws_api_gateway_resource.terraless-lambda-{{.FunctionName}}-{{.ResourceNameForPath}}"]
  {{ end }}

  rest_api_id = "${aws_api_gateway_rest_api.terraless-api-gateway.id}"
  resource_id = "${aws_api_gateway_method.terraless-lambda-{{.FunctionName}}-{{.Idx}}.resource_id}"
  http_method = "${aws_api_gateway_method.terraless-lambda-{{.FunctionName}}-{{.Idx}}.http_method}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.lambda-{{.FunctionName}}.invoke_arn}"
}

`,
	"sqs": `
# Function {{ .ProjectName }} {{ .FunctionName }} EventKey: {{.Idx}}

`,
}
