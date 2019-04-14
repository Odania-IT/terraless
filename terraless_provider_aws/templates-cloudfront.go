package terraless_provider_aws

import (
	"archive/zip"
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

func RenderUploadTemplates(terralessData schema.TerralessData, buffer bytes.Buffer) bytes.Buffer {
	currentConfig := terralessData.Config
	for _, upload := range currentConfig.Uploads {
		if upload.Type == "s3" && upload.Cloudfront.Domain != "" {
			upload.Cloudfront.Handler = lambdaAtEdgeHandler(upload.Cloudfront.Handler)

			if upload.Cloudfront.Handler != "" {
				buffer = renderLambdaAtEdge(currentConfig, upload, buffer)
			}

			logrus.Debugf("Generating cloudfront template for %s\n", upload.Cloudfront.Domain)
			upload.Cloudfront.Aliases = append(upload.Cloudfront.Aliases, upload.Cloudfront.Domain)
			upload.Environment = terralessData.Arguments.Environment
			upload.ProjectName = currentConfig.ProjectName
			upload.Certificate = currentConfig.Certificates[upload.Cloudfront.Certificate]
			upload.OwnCertificate = upload.Certificate.Domain != ""

			buffer = templates.RenderTemplateToBuffer(upload, buffer, awsTemplates("cloudfront.tf.tmpl"), "terraless-upload-cloudfront")
			for _, alias := range upload.Cloudfront.Aliases {
				buffer = Route53AliasRecordFor(alias, upload.Certificate.ZoneId, buffer)
			}

			addTerralessLambdaRole = true
		}
	}

	return buffer
}

func renderLambdaAtEdge(currentConfig schema.TerralessConfig, upload schema.TerralessUpload, buffer bytes.Buffer) bytes.Buffer {
	targetFile := lambdaAtEdgeZip(currentConfig)
	data := map[string]string{
		"FileName":    targetFile,
		"Handler": upload.Cloudfront.Handler,
		"ProjectName": currentConfig.ProjectName,
	}
	return templates.RenderTemplateToBuffer(data, buffer, awsTemplates("lambda-at-edge.tf.tmpl"), "lambda-at-edge.tf")
}

func lambdaAtEdgeHandler(handler string) string {
	availableHandlers := []string {
		"redirectToWww",
		"singleEntryPointHandler",
		"staticHandler",
	}

	if support.Contains(availableHandlers, handler) {
		return handler
	}

	return ""
}

func lambdaAtEdgeZip(config schema.TerralessConfig) string {
	targetFile := filepath.Join(config.TargetPath, "lambda-at-edge.zip")

	info, _ := os.Stat(targetFile)
	if info != nil && info.Size() > 0 {
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
