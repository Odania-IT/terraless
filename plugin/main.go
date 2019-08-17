package plugin

import (
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	ExtensionPluginType = "Extension"
	ProviderPluginType  = "Provider"
)

type PluginType struct {
	Type            string
	Prefix          string
	HandshakeConfig plugin.HandshakeConfig
}

var pluginTypes = []PluginType{
	{
		Type:   ExtensionPluginType,
		Prefix: "terraless-extension",
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "extension-plugin",
			MagicCookieValue: "terraless",
		},
	},
	{
		Type:   ProviderPluginType,
		Prefix: "terraless-provider",
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "provider-plugin",
			MagicCookieValue: "terraless",
		},
	},
}

var clients []*plugin.Client

func Extensions() []schema.Extension {
	var result []schema.Extension

	for _, pluginData := range pluginsData {
		if pluginData.Type.Type == ExtensionPluginType {
			result = append(result, pluginData.Extension)
		}
	}

	return result
}

func Providers() []schema.Provider {
	var result []schema.Provider

	for _, pluginData := range pluginsData {
		if pluginData.Type.Type == ProviderPluginType {
			result = append(result, pluginData.Provider)
		}
	}

	return result
}

func HandlePlugins(configPlugins []schema.TerralessPlugin, pluginDirectory string) bool {
	logrus.Info("Processing Plugins")
	var pluginInfos []string
	for _, configPlugin := range configPlugins {
		version := configPlugin.Version
		if configPlugin.Version == "" {
			version = "~any"
		}

		currentVersion := verifyPluginInstalled(configPlugin)

		pluginInfo := fmt.Sprintf("  - %s (Current Version: %s Wanted Version: %s)", configPlugin.Name, currentVersion, version)
		pluginInfos = append(pluginInfos, pluginInfo)
	}

	logrus.Infof("Terraless Plugins:\n%s\n", strings.Join(pluginInfos, "\n"))

	return true
}

func verifyPluginInstalled(plugin schema.TerralessPlugin) string {
	var existingVersions []string
	for _, pluginData := range pluginsData {
		if plugin.Name == pluginData.Name {
			if plugin.Version == pluginData.Version || plugin.Version == "~any" {
				logrus.Debugf("Plugin is present %s Version: %s\n", plugin.Name, plugin.Version)
				return plugin.Version
			}

			existingVersions = append(existingVersions, pluginData.Version)
		}
	}

	if len(existingVersions) > 0 {
		return strings.Join(existingVersions, ", ")
	}

	return installPlugin(plugin)
}

func installPlugin(plugin schema.TerralessPlugin) string {

	return plugin.Version
}

func (pluginType PluginType) pluginMapValue() plugin.Plugin {
	if pluginType.Type == ExtensionPluginType {
		return &schema.ExtensionPlugin{}
	}

	if pluginType.Type == ProviderPluginType {
		return &schema.ProviderPlugin{}
	}

	logrus.Warnf("Missing plugin data for plugin type %s\n", pluginType.Type)
	return nil
}

func ClosePlugins() {
	logrus.Debug("Closing all plugins")
	for _, client := range clients {
		client.Kill()
	}
}
