package schema

import (
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"net/rpc"
)

type Extension interface {
	Exec(data TerralessData) error
	Info(logLevel string) PluginInfo
}

// RPC
type ExtensionRPC struct {
	client *rpc.Client
}

func (g *ExtensionRPC) Exec(data TerralessData) error {
	err := g.client.Call("Plugin.Exec", data, new(interface{}))
	if err != nil {
		logrus.Fatal("Error executing Extension:Exec()", err)
	}

	return err
}

func (g *ExtensionRPC) Info(logLevel string) PluginInfo {
	var resp PluginInfo
	err := g.client.Call("Plugin.Info", logLevel, &resp)
	if err != nil {
		logrus.Fatal("Error executing Extension:Info()", err)
	}

	return resp
}

// RPC Server
type ExtensionRPCServer struct {
	Impl Extension
}

func (server *ExtensionRPCServer) Exec(data TerralessData, resp *error) error {
	*resp = server.Impl.Exec(data)
	return nil
}

func (server *ExtensionRPCServer) Info(logLevel string, resp *PluginInfo) error {
	*resp = server.Impl.Info(logLevel)
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
