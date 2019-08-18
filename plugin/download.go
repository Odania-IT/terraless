package plugin

import (
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"runtime"
)

func DownloadPlugin(plugin schema.TerralessPlugin, pluginDirectory string) {
	version := plugin.Version
	if plugin.Version == "" {
		version = "latest"
	}

	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}

	err := os.MkdirAll(pluginDirectory, 0755)
	if err != nil {
		logrus.Fatalf("Could not create plugin directory: %s\n", pluginDirectory)
	}

	logrus.Infof("Trying to download plugin %s [Version: %s]\n", plugin.Name, version)
	url := "https://terraless-plugins.s3.eu-central-1.amazonaws.com/" +
		plugin.Name +
		"/" +
		version +
		"/" +
		plugin.Name +
		"_" +
		runtime.GOOS +
		"_amd64" +
		extension

	logrus.Debugf("Downloading plugin from url: %s\n", url)
	fileName := pluginDirectory + "/" + plugin.Name + "-" + version
	err = downloadFile(fileName, url)

	if err != nil {
		logrus.Fatalf("Failed downloading plugin %s [Version: %s]\n", plugin.Name, version)
	}

	if runtime.GOOS != "windows" {
		err := os.Chmod(fileName, 0755)
		if err != nil {
			logrus.Fatalf("Can not make plugin executable!", err)
		}
	}
	logrus.Infof("Downloaded plugin to %s\n", fileName)
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Failed downloading from url %s\n", url))
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
