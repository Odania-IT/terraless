package schema

import "bytes"

type Provider struct {
	CanHandle                  CanHandleFunc
	PrepareSession             PrepareSessionFunc
	ProcessUpload              ProcessUploadFunc
	RenderAuthorizerTemplates  RenderAuthorizerTemplatesFunc
	RenderCertificateTemplates RenderCertificateTemplatesFunc
	RenderEndpointTemplates    RenderEndpointTemplatesFunc
	RenderFunctionTemplates    RenderFunctionTemplatesFunc
	RenderUploadTemplates      RenderUploadTemplatesFunc
}

type CanHandleFunc func(resourceType string) bool
type PrepareSessionFunc func(terralessConfig TerralessConfig)
type ProcessUploadFunc func(terralessData TerralessData, upload TerralessUpload)
type RenderAuthorizerTemplatesFunc func(config TerralessConfig, buffer bytes.Buffer) bytes.Buffer
type RenderCertificateTemplatesFunc func(config TerralessConfig, buffer bytes.Buffer) bytes.Buffer
type RenderEndpointTemplatesFunc func(config TerralessConfig, buffer bytes.Buffer) bytes.Buffer
type RenderFunctionTemplatesFunc func(resourceType string, functionEvents FunctionEvents, terralessData *TerralessData, buffer bytes.Buffer) bytes.Buffer
type RenderUploadTemplatesFunc func(currentConfig TerralessConfig, buffer bytes.Buffer) bytes.Buffer
