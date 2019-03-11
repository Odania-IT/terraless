package terraless_provider_aws

import (
	"archive/zip"
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/templates"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

func renderCloudfrontTemplates(currentConfig schema.TerralessConfig, targetFileName string) {
	buffer := bytes.Buffer{}

	for _, upload := range currentConfig.Uploads {
		targetFile := lambdaAtEdgeZip(currentConfig)
		data := map[string]string {
			"FileName": targetFile,
			"ProjectName": currentConfig.ProjectName,
		}
		buffer = templates.RenderTemplateToBuffer(data, buffer, awsTemplates("lambda-at-edge.tf.tmpl"))

		if upload.Cloudfront.Domain != "" {
			logrus.Debugf("Generating cloudfront template for %s\n", upload.Cloudfront.Domain)
			upload.Cloudfront.Aliases = append(upload.Cloudfront.Aliases, upload.Cloudfront.Domain)
			upload.ProjectName = currentConfig.ProjectName
			upload.Certificate = currentConfig.Certificates[upload.Cloudfront.Certificate]
			upload.OwnCertificate = upload.Certificate.Domain != ""
			upload.LambdaAtEdgeFile = targetFile

			buffer = templates.RenderTemplateToBuffer(upload, buffer, awsTemplates("cloudfront.tf.tmpl"))
			for _, alias := range upload.Cloudfront.Aliases {
				buffer = Route53AliasRecordFor(alias, upload.Certificate.ZoneId, buffer)
			}
		}
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

func lambdaAtEdgeZip(config schema.TerralessConfig) string {
	targetFile := filepath.Join(config.TargetPath, "lambda-at-edge.zip")

	info, _ := os.Stat(targetFile)
	if info.Size() > 0 {
		err := os.Remove(targetFile)

		if err != nil {
			logrus.Fatalf("Failed deleting old zip %s\n", err)
		}
	}

	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)
	writer, err := zipWriter.Create("lambda.js")

	if err != nil {
		logrus.Fatalf("Failed creating zip %s\n", err)
	}

	_, err = writer.Write([]byte(awsTemplates("lambda-at-edge.js")))

	if err != nil {
		logrus.Fatalf("Failed creating zip %s\n", err)
	}

	err = zipWriter.Close()
	if err != nil {
		logrus.Fatalf("Failed creating zip %s\n", err)
	}

	err = ioutil.WriteFile(targetFile, buffer.Bytes(), 0640)

	if err != nil {
		logrus.Fatalf("Failed creating zip %s\n", err)
	}

	return targetFile
}
