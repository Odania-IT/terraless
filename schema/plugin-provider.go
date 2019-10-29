package schema

import (
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"net/rpc"
)

type Provider interface {
	CanHandle(resourceType string) bool
	FinalizeTemplates(terralessData TerralessData) string
	Info() PluginInfo
	PrepareSession(terralessConfig TerralessConfig) map[string]string
	ProcessUpload(terralessData TerralessData, upload TerralessUpload) []string
	RenderAuthorizerTemplates(config TerralessConfig) string
	RenderCertificateTemplates(config TerralessConfig) string
	RenderEndpointTemplates(config TerralessConfig) string
	RenderFunctionTemplates(resourceType string, functionEvents FunctionEvents, terralessData *TerralessData) string
	RenderUploadTemplates(terralessData TerralessData) string
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

func (g *ProviderRPC) FinalizeTemplates(terralessData TerralessData) string {
	var resp string
	err := g.client.Call("Plugin.FinalizeTemplates", terralessData, &resp)
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

func (g *ProviderRPC) PrepareSession(terralessConfig TerralessConfig) map[string]string {
	var resp map[string]string
	err := g.client.Call("Plugin.PrepareSession", terralessConfig, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:PrepareSession()", err)
	}

	return resp
}

type ProcessUploadArgs struct {
	TerralessData TerralessData
	Upload        TerralessUpload
}

func (g *ProviderRPC) ProcessUpload(terralessData TerralessData, upload TerralessUpload) []string {
	var resp []string
	args := &ProcessUploadArgs{terralessData, upload}
	err := g.client.Call("Plugin.ProcessUpload", args, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:ProcessUpload()", err)
	}

	return resp
}

func (g *ProviderRPC) RenderAuthorizerTemplates(config TerralessConfig) string {
	var resp string
	err := g.client.Call("Plugin.RenderAuthorizerTemplates", config, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:RenderAuthorizerTemplates()", err)
	}

	return resp
}

func (g *ProviderRPC) RenderCertificateTemplates(config TerralessConfig) string {
	var resp string
	err := g.client.Call("Plugin.RenderCertificateTemplates", config, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:RenderCertificateTemplates()", err)
	}

	return resp
}

func (g *ProviderRPC) RenderEndpointTemplates(config TerralessConfig) string {
	var resp string
	err := g.client.Call("Plugin.RenderEndpointTemplates", config, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:RenderEndpointTemplates()", err)
	}

	return resp
}

type RenderFunctionTemplatesArgs struct {
	ResourceType string
	FunctionEvents FunctionEvents
	TerralessData *TerralessData
}

func (g *ProviderRPC) RenderFunctionTemplates(resourceType string, functionEvents FunctionEvents, terralessData *TerralessData) string {
	var resp string
	args := &RenderFunctionTemplatesArgs{
		ResourceType: resourceType,
		FunctionEvents: functionEvents,
		TerralessData: terralessData,
	}
	err := g.client.Call("Plugin.RenderFunctionTemplates", args, &resp)
	if err != nil {
		logrus.Fatal("Error executing Provider:RenderFunctionTemplates()", err)
	}

	return resp
}

func (g *ProviderRPC) RenderUploadTemplates(terralessData TerralessData) string {
	var resp string
	err := g.client.Call("Plugin.RenderUploadTemplates", terralessData, &resp)
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

func (server *ProviderRPCServer) FinalizeTemplates(terralessData TerralessData, resp *string) error {
	*resp = server.Impl.FinalizeTemplates(terralessData)
	return nil
}

func (server *ProviderRPCServer) Info(args interface{}, resp *PluginInfo) error {
	*resp = server.Impl.Info()
	return nil
}

func (server *ProviderRPCServer) PrepareSession(terralessConfig TerralessConfig, resp *map[string]string) error {
	*resp = server.Impl.PrepareSession(terralessConfig)
	return nil
}

func (server *ProviderRPCServer) ProcessUpload(args ProcessUploadArgs, resp *[]string) error {
	*resp = server.Impl.ProcessUpload(args.TerralessData, args.Upload)
	return nil
}

func (server *ProviderRPCServer) RenderAuthorizerTemplates(config TerralessConfig, resp *string) error {
	*resp = server.Impl.RenderAuthorizerTemplates(config)
	return nil
}

func (server *ProviderRPCServer) RenderCertificateTemplates(config TerralessConfig, resp *string) error {
	*resp = server.Impl.RenderCertificateTemplates(config)
	return nil
}

func (server *ProviderRPCServer) RenderEndpointTemplates(config TerralessConfig, resp *string) error {
	*resp = server.Impl.RenderEndpointTemplates(config)
	return nil
}

func (server *ProviderRPCServer) RenderFunctionTemplates(args RenderFunctionTemplatesArgs, resp *string) error {
	*resp = server.Impl.RenderFunctionTemplates(args.ResourceType, args.FunctionEvents, args.TerralessData)
	return nil
}

func (server *ProviderRPCServer) RenderUploadTemplates(terralessData TerralessData, resp *string) error {
	*resp = server.Impl.RenderUploadTemplates(terralessData)
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
