package plugin

import (
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

type PluginData struct {
	Name      string
	Version   string
	Type      PluginType
	Extension schema.Extension
	Provider  schema.Provider
}

var pluginsData []PluginData

func ExistingPlugins(arguments schema.Arguments) map[string]schema.TerralessPlugin {
	result := map[string]schema.TerralessPlugin{}

	logrus.Debugf("Listing plugin directory %s\n", arguments.PluginDirectory)
	files, err := ioutil.ReadDir(arguments.PluginDirectory)
	if err != nil {
		logrus.Debugf("Failed reading plugin directory: %s\n", arguments.PluginDirectory)

		return result
	}

	plugins := []string{}
	for _, file := range files {
		fileName := file.Name()
		pluginData := detectPluginAndLoad(filepath.Join(arguments.PluginDirectory, fileName))
		plugins = append(plugins, pluginData.Name)
	}

	return result
}

func PluginsData() []PluginData {
	return pluginsData
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
		Type: pluginType,
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
