package uploads

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/sirupsen/logrus"
)

func ProcessUploads(terralessData schema.TerralessData) {
	terralessConfig := terralessData.Config
	logrus.Debug("Processing uploads")
	if len(terralessConfig.Uploads) == 0 {
		logrus.Debug("... no uploads")
		return
	}

	for _, upload := range terralessConfig.Uploads {
		logrus.Debugf("Processing upload: %#v\n", upload)

		for _, terralessProvider := range terralessData.TerralessProviders {
			terralessProvider.ProcessUpload(terralessConfig, upload)
		}
	}
}
