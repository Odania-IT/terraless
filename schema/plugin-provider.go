package schema

import (
	"bytes"
	"net/rpc"
)

type Provider interface {
	CanHandle(resourceType string) bool
	FinalizeTemplates(terralessData TerralessData, buffer bytes.Buffer) bytes.Buffer
	Info() PluginInfo
	Name() string
	PrepareSession(terralessConfig TerralessConfig)
	ProcessUpload(terralessData TerralessData, upload TerralessUpload) []string
	RenderAuthorizerTemplates(config TerralessConfig, buffer bytes.Buffer) bytes.Buffer
	RenderCertificateTemplates(config TerralessConfig, buffer bytes.Buffer) bytes.Buffer
	RenderEndpointTemplates(config TerralessConfig, buffer bytes.Buffer) bytes.Buffer
	RenderFunctionTemplates(resourceType string, functionEvents FunctionEvents, terralessData *TerralessData, buffer bytes.Buffer) bytes.Buffer
	RenderUploadTemplates(terralessData TerralessData, buffer bytes.Buffer) bytes.Buffer
}

// RPC
type ProviderRPC struct {
	client *rpc.Client
}

// RPC Server
type ProviderRPCServer struct {
	Impl Provider
}

// Implementation of plugin.Plugin

type ProviderPlugin struct {
	Impl Provider
}
