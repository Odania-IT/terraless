package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func renderCertificateTemplates(config schema.TerralessConfig, targetFileName string) {
	buffer := bytes.Buffer{}

	for _, certificate := range config.Certificates {
		logrus.Debugf("Generating certificate template for %s\n", certificate.Domain)
		certificate.ProjectName = config.ProjectName
		certificate.TerraformName = "terraless-certificate-" + support.SanitizeString(certificate.Domain)

		buffer = templates.RenderTemplateToBuffer(certificate, buffer, awsTemplates("certificate.tf.tmpl"))
	}

	// Write to file
	targetFile, err := os.Create(targetFileName)
	if err != nil {
		logrus.Fatal("Failed creating file: ", filepath.Base(targetFileName), err)
	}

	_, err = targetFile.Write(buffer.Bytes())
	if err != nil {
		logrus.Fatal("Failed wrtiting to file: ", filepath.Base(targetFileName), err)
	}
}
