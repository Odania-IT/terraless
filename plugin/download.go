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
	err := downloadFile(pluginDirectory+"/"+plugin.Name+"_"+version, url)

	if err != nil {
		logrus.Fatalf("Failed downloading plugin %s [Version: %s]\n", plugin.Name, version)
	}
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
