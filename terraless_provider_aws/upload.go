package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var uploadFileFunc = s3Upload
func processUpload(terralessData schema.TerralessData, upload schema.TerralessUpload) []string {
	config := terralessData.Config
	if upload.Type != "s3" {
		logrus.Debugf("AWS-Provider can not handle upload %s\n", upload.Type)
		return []string{}
	}

	logrus.Infof("Processing Upload Source: %s Bucket: %s\n", upload.Source, upload.Bucket)
	provider := config.Providers[schema.ProcessString(upload.Provider, terralessData.Arguments, terralessData.Config.Settings)]
	sess := sessionForProvider(provider)

	svc := s3manager.NewUploader(sess)

	uploadedFiles := recursiveUpload(filepath.Join(config.SourcePath, upload.Source), upload.Target, upload.Bucket, svc)
	logrus.Debugf("Uploaded files: %s\n", uploadedFiles)

	return uploadedFiles
}

func recursiveUpload(sourceDir string, targetPrefix string, bucketName string, svc *s3manager.Uploader) []string {
	var result []string
	matches, err := filepath.Glob(filepath.Join(sourceDir, "**"))

	if err != nil {
		logrus.Fatal("Failed locating upload files: ", filepath.Base(sourceDir), " Error: ", err)
	}

	logrus.Debugf("%d Objects found to upload to %s\n", len(matches), bucketName)
	for _, filename := range matches {
		info, err := os.Stat(filename)
		targetFile := filepath.Join(targetPrefix, filepath.Base(filename))

		if err != nil {
			logrus.Fatalf("Can not stat %s! Error: %s\n", filename, err)
		}

		if info.IsDir() {
			logrus.Debugf("Processing directory %s", targetFile)
			result = append(result, recursiveUpload(filename, targetFile, bucketName, svc)...)
		} else {
			err = addFileToS3(svc, bucketName, filename, targetFile)
			if err != nil {
				logrus.Fatalf("Failed uploading file %s to s3 bucket %s\n", targetFile, bucketName)
			}

			result = append(result, targetFile)
		}
	}

	return result
}

func addFileToS3(svc *s3manager.Uploader, bucket string, filename string, targetFile string) error {
	// Open the file for use
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	_, err = file.Read(buffer)

	if err != nil {
		logrus.Fatalf("Can not read file %s! Error: %s\n", filename, err)
	}

	err = file.Close()
	if err != nil {
		logrus.Fatalf("Can close file %s! Error: %s\n", filename, err)
	}

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	contentType := support.DetectContentType(filename)
	uploadResult, err := uploadFileFunc(svc, s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(targetFile),
		ACL:         aws.String("private"),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		logrus.Fatalf("Can not read file %s! Error: %s\n", filename, err)
	}

	logrus.Debugf("Successfully uploaded %s to %s [Content-Type: %s]\n", filename, uploadResult.Location, contentType)
	return err
}

func s3Upload(svc *s3manager.Uploader, uploadInput s3manager.UploadInput) (*s3manager.UploadOutput, error) {
	return svc.Upload(&uploadInput)
}
