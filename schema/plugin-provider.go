package schema

import (
	"bytes"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"net/rpc"
)

type Provider interface {
	CanHandle(resourceType string) bool
	FinalizeTemplates(terralessData TerralessData, buffer bytes.Buffer) bytes.Buffer
	Info() PluginInfo
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

func (g *ProviderRPC) CanHandle(resourceType string) bool {
	var resp bool
	err := g.client.Call("Plugin.CanHandle", resourceType, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:CanHandle()", err)
	}

	return resp
}

type finalizeTemplatesArgs struct {
	TerralessData TerralessData
	Buffer        bytes.Buffer
}

func (g *ProviderRPC) FinalizeTemplates(terralessData TerralessData, buffer bytes.Buffer) bytes.Buffer {
	var resp bytes.Buffer
	args := &finalizeTemplatesArgs{terralessData, buffer}
	err := g.client.Call("Plugin.FinalizeTemplates", args, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:FinalizeTemplates()", err)
	}

	return resp
}

func (g *ProviderRPC) Info() PluginInfo {
	var resp PluginInfo
	err := g.client.Call("Plugin.Info", new(interface{}), &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:Info()", err)
	}

	return resp
}

func (g *ProviderRPC) PrepareSession(terralessConfig TerralessConfig) {
	var resp string
	err := g.client.Call("Plugin.PrepareSession", terralessConfig, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:PrepareSession()", err)
	}
}

type processUploadArgs struct {
	TerralessData TerralessData
	Upload        TerralessUpload
}

func (g *ProviderRPC) ProcessUpload(terralessData TerralessData, upload TerralessUpload) []string {
	var resp []string
	args := &processUploadArgs{terralessData, upload}
	err := g.client.Call("Plugin.ProcessUpload", args, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:ProcessUpload()", err)
	}

	return resp
}

type renderWithConfigArgs struct {
	Config TerralessConfig
	Buffer bytes.Buffer
}

func (g *ProviderRPC) RenderAuthorizerTemplates(config TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
	var resp bytes.Buffer
	args := &renderWithConfigArgs{Config: config, Buffer: buffer}
	err := g.client.Call("Plugin.RenderAuthorizerTemplates", args, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:RenderAuthorizerTemplates()", err)
	}

	return resp
}

func (g *ProviderRPC) RenderCertificateTemplates(config TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
	var resp bytes.Buffer
	args := &renderWithConfigArgs{Config: config, Buffer: buffer}
	err := g.client.Call("Plugin.RenderCertificateTemplates", args, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:RenderCertificateTemplates()", err)
	}

	return resp
}

func (g *ProviderRPC) RenderEndpointTemplates(config TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
	var resp bytes.Buffer
	args := &renderWithConfigArgs{Config: config, Buffer: buffer}
	err := g.client.Call("Plugin.RenderEndpointTemplates", args, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:RenderEndpointTemplates()", err)
	}

	return resp
}

type renderFunctionTemplatesArgs struct {
	ResourceType string
	FunctionEvents FunctionEvents
	TerralessData *TerralessData
	Buffer bytes.Buffer
}

func (g *ProviderRPC) RenderFunctionTemplates(resourceType string, functionEvents FunctionEvents, terralessData *TerralessData, buffer bytes.Buffer) bytes.Buffer {
	var resp bytes.Buffer
	args := &renderFunctionTemplatesArgs{
		ResourceType: resourceType,
		FunctionEvents: functionEvents,
		TerralessData: terralessData,
		Buffer: buffer,
	}
	err := g.client.Call("Plugin.RenderFunctionTemplates", args, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:RenderFunctionTemplates()", err)
	}

	return resp
}

type renderUploadTemplatesArgs struct {
	TerralessData TerralessData
	Buffer        bytes.Buffer
}

func (g *ProviderRPC) RenderUploadTemplates(terralessData TerralessData, buffer bytes.Buffer) bytes.Buffer {
	var resp bytes.Buffer
	args := &renderUploadTemplatesArgs{TerralessData: terralessData, Buffer: buffer}
	err := g.client.Call("Plugin.RenderUploadTemplates", args, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:RenderUploadTemplates()", err)
	}

	return resp
}

// RPC Server
type ProviderRPCServer struct {
	Impl Provider
}

func (server *ProviderRPCServer) CanHandle(resourceType string, resp *bool) error {
	*resp = server.Impl.CanHandle(resourceType)
	return nil
}

func (server *ProviderRPCServer) FinalizeTemplates(terralessData TerralessData, buffer bytes.Buffer, resp *bytes.Buffer) error {
	*resp = server.Impl.FinalizeTemplates(terralessData, buffer)
	return nil
}

func (server *ProviderRPCServer) Info(resp *PluginInfo) error {
	*resp = server.Impl.Info()
	return nil
}

func (server *ProviderRPCServer) PrepareSession(terralessConfig TerralessConfig) error {
	server.Impl.PrepareSession(terralessConfig)
	return nil
}

func (server *ProviderRPCServer) ProcessUpload(terralessData TerralessData, upload TerralessUpload, resp *[]string) error {
	*resp = server.Impl.ProcessUpload(terralessData, upload)
	return nil
}

func (server *ProviderRPCServer) RenderAuthorizerTemplates(config TerralessConfig, buffer bytes.Buffer, resp *bytes.Buffer) error {
	*resp = server.Impl.RenderAuthorizerTemplates(config, buffer)
	return nil
}

func (server *ProviderRPCServer) RenderCertificateTemplates(config TerralessConfig, buffer bytes.Buffer, resp *bytes.Buffer) error {
	*resp = server.Impl.RenderCertificateTemplates(config, buffer)
	return nil
}

func (server *ProviderRPCServer) RenderEndpointTemplates(config TerralessConfig, buffer bytes.Buffer, resp *bytes.Buffer) error {
	*resp = server.Impl.RenderEndpointTemplates(config, buffer)
	return nil
}

func (server *ProviderRPCServer) RenderFunctionTemplates(resourceType string, functionEvents FunctionEvents, terralessData *TerralessData, buffer bytes.Buffer, resp *bytes.Buffer) error {
	*resp = server.Impl.RenderFunctionTemplates(resourceType, functionEvents, terralessData, buffer)
	return nil
}

func (server *ProviderRPCServer) RenderUploadTemplates(terralessData TerralessData, buffer bytes.Buffer, resp *bytes.Buffer) error {
	*resp = server.Impl.RenderUploadTemplates(terralessData, buffer)
	return nil
}

// Implementation of plugin.Plugin

type ProviderPlugin struct {
	Impl Provider
}

func (plugin *ProviderPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &ProviderRPCServer{Impl: plugin.Impl}, nil
}

func (ProviderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ProviderRPC{client: c}, nil
}
