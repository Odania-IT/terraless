package plugin

import (
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

type PluginData struct {
	Name      string
	Version   string
	Type      PluginType
	Extension schema.Extension
	Provider  schema.Provider
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
var pluginsData []PluginData

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
	existingPlugins := ExistingPlugins(pluginDirectory)
	var pluginInfos = []string{}
	for _, configPlugin := range configPlugins {
		version := configPlugin.Version
		if configPlugin.Version == "" {
			version = "~any"
		}

		currentVersion := verifyPluginInstalled(configPlugin, existingPlugins)

		pluginInfo := fmt.Sprintf("  - %s (Current Version: %s Wanted Version: %s)", configPlugin.Name, currentVersion, version)
		pluginInfos = append(pluginInfos, pluginInfo)
	}

	logrus.Infof("Terraless Plugins:\n%s\n", strings.Join(pluginInfos, "\n"))

	return true
}

func verifyPluginInstalled(plugin schema.TerralessPlugin, existingPlugins map[string]schema.TerralessPlugin) string {
	var existingVersions []string
	for _, existingPlugin := range existingPlugins {
		if plugin.Name == existingPlugin.Name {
			if plugin.Version == existingPlugin.Version || plugin.Version == "~any" {
				logrus.Debugf("Plugin is present %s Version: %s\n", plugin.Name, plugin.Version)
				return plugin.Version
			}

			existingVersions = append(existingVersions, existingPlugin.Version)
		}
	}

	if len(existingVersions) > 0 {
		return strings.Join(existingVersions, ", ")
	}

	return installPlugin(plugin)
}

func ExistingPlugins(pluginDirectory string) map[string]schema.TerralessPlugin {
	result := map[string]schema.TerralessPlugin{}

	logrus.Debugf("Listing plugin directory %s\n", pluginDirectory)
	files, err := ioutil.ReadDir(pluginDirectory)
	if err != nil {
		logrus.Debugf("Failed reading plugin directory: %s\n", pluginDirectory)

		return result
	}

	for _, file := range files {
		fileName := file.Name()
		logrus.Debug(fileName)
		detectPluginAndLoad(filepath.Join(pluginDirectory, fileName))
		logrus.Debug(fileName)
	}

	return result
}

func installPlugin(plugin schema.TerralessPlugin) string {

	return plugin.Version
}

func detectPluginAndLoad(file string) PluginData {
	fileName := filepath.Base(file)

	for _, pluginType := range pluginTypes {
		if strings.HasPrefix(fileName, pluginType.Prefix) {
			logrus.Debugf("Detected plugin of type '%s' - %s\n", pluginType.Type, fileName)
			return loadPlugin(pluginType, file)
		}
	}

	logrus.Warnf("Found unknown plugin %s\n", fileName)
	return PluginData{}
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

func loadPlugin(pluginType PluginType, file string) PluginData {
	fileName := filepath.Base(file)

	pluginMap := map[string]plugin.Plugin{}
	pluginMap[pluginType.Type] = pluginType.pluginMapValue()

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: pluginType.HandshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(file),
		Logger:          logger,
	})
	clients = append(clients, client)

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		logrus.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pluginType.Type)
	if err != nil {
		logrus.Fatal(err)
	}

	pluginData := PluginData{
		Type:      pluginType,
	}

	var pluginInfo schema.PluginInfo
	if pluginType.Type == ExtensionPluginType {
		extension := raw.(schema.Extension)
		pluginInfo = extension.Info()
		pluginData.Extension = extension
	} else if pluginType.Type == ProviderPluginType {
		provider := raw.(schema.Provider)
		pluginInfo = provider.Info()
		pluginData.Provider = provider
	} else {
		logrus.Warnf("Could not detect plugin type for %s\n", fileName)
	}

	logrus.Debugf("Loaded Plugin %s Version: %s\n", pluginInfo.Name, pluginInfo.Version)
	pluginData.Name = pluginInfo.Name
	pluginData.Version = pluginInfo.Version
	pluginsData = append(pluginsData, pluginData)
	return pluginData
}
