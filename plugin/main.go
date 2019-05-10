package plugin

import (
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

func HandlePlugins(plugins []schema.TerralessPlugin, pluginDirectory string) bool {
	logrus.Info("Processing Plugins")
	existingPlugins := existingPlugins(pluginDirectory)
	var pluginInfos = []string{}
	for _, plugin := range plugins {
		version := plugin.Version
		if plugin.Version == "" {
			version = "~any"
		}

		currentVersion := verifyPluginInstalled(plugin, existingPlugins)

		pluginInfo := fmt.Sprintf("  - %s (Current Version: %s Wanted Version: %s)", plugin.Name, currentVersion, version)
		pluginInfos = append(pluginInfos, pluginInfo)
	}

	logrus.Infof("Terraless Plugins:\n%s\n", strings.Join(pluginInfos, "\n"))

	return true
}

func verifyPluginInstalled(plugin schema.TerralessPlugin, existingPlugins map[string]schema.TerralessPlugin) string {
	var existingVersions []string
	for _, existingPlugin := range existingPlugins {
		if plugin.Name == existingPlugin.Name {
			if plugin.Version == existingPlugin.Version ||plugin.Version == "~any" {
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

func existingPlugins(pluginDirectory string) map[string]schema.TerralessPlugin {
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
	}

	return result
}

func installPlugin(plugin schema.TerralessPlugin) string {


	return plugin.Version
}
