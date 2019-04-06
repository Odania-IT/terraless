package terraless_provider_aws

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var uploadedFiles int
func uploadFileMock(svc *s3manager.Uploader, bucket string, filename string, targetFile string) error {
	uploadedFiles += 1
	return nil
}

func TestTemplatesFunctions_RecursiveUpload(t *testing.T) {
	// given
	dir, _ := os.Getwd()
	uploadFileFunc = uploadFileMock
	terralessData := schema.TerralessData{
		Config: schema.TerralessConfig{
			SourcePath: dir,
		},
	}
	upload := schema.TerralessUpload{
		Type: "s3",
		Source: "templates",
	}

	// when
	uploadedFilenames := processUpload(terralessData, upload)

	// then
	expected := []string{
		"authorizer.tf.tmpl",
		"certificate.tf.tmpl",
		"cloudfront.tf.tmpl",
		"endpoint.tf.tmpl",
		"iam.tf.tmpl",
		"lambda-at-edge.js",
		"lambda-at-edge.tf.tmpl",
	}
	assert.Equal(t, 7, uploadedFiles)
	assert.Equal(t, expected, uploadedFilenames)
}
