package schema

import (
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"net/rpc"
)

type Extension interface {
	Exec(globalConfig TerralessGlobalConfig, data TerralessData) error
	Info() PluginInfo
}

type ExtensionExecArgs struct {
	GlobalConfig TerralessGlobalConfig
	TerralessData TerralessData
}

// RPC
type ExtensionRPC struct {
	client *rpc.Client
}

func (g *ExtensionRPC) Exec(globalConfig TerralessGlobalConfig, terralessData TerralessData) error {
	args := &ExtensionExecArgs{globalConfig, terralessData}
	err := g.client.Call("Plugin.Exec", args, new(interface{}))
	if err != nil {
		logrus.Fatal("Error executing Extension:Exec()", err)
	}

	return err
}

func (g *ExtensionRPC) Info() PluginInfo {
	var resp PluginInfo
	err := g.client.Call("Plugin.Info", new(interface{}), &resp)
	if err != nil {
		logrus.Fatal("Error executing Extension:Info()", err)
	}

	return resp
}

// RPC Server
type ExtensionRPCServer struct {
	Impl Extension
}

func (server *ExtensionRPCServer) Exec(args ExtensionExecArgs, resp *error) error {
	*resp = server.Impl.Exec(args.GlobalConfig, args.TerralessData)
	return nil
}

func (server *ExtensionRPCServer) Info(args interface{}, resp *PluginInfo) error {
	*resp = server.Impl.Info()
	return nil
}

// Implementation of plugin.Plugin
type ExtensionPlugin struct {
	Impl Extension
}

func (plugin *ExtensionPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &ExtensionRPCServer{Impl: plugin.Impl}, nil
}

func (ExtensionPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ExtensionRPC{client: c}, nil
}
