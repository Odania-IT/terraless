package templates

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/sirupsen/logrus"
)

func ProcessUploads(terralessData schema.TerralessData, providers []schema.Provider) {
	terralessConfig := terralessData.Config
	logrus.Debug("Processing uploads")
	if len(terralessConfig.Uploads) == 0 {
		logrus.Debug("... no uploads")
		return
	}

	for _, upload := range terralessConfig.Uploads {
		logrus.Debugf("Processing upload: %#v\n", upload)

		for _, terralessProvider := range providers {
			_ = terralessProvider.ProcessUpload(terralessData, upload)
		}
	}
}
