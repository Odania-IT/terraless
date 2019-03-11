package uploads

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/aws/aws-sdk-go/aws"
	credentials2 "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

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
			err = AddFileToS3(svc, bucketName, filename, targetFile)
			result = append(result, targetFile)
		}
	}

	return result
}

func AddFileToS3(svc *s3manager.Uploader, bucket string, filename string, targetFile string) error {
	// Open the file for use
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	_, err = file.Read(buffer)

	if err != nil {
		logrus.Fatalf("Can not read file %s! Error: %s\n", filename, err)
	}

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	contentType := support.DetectContentType(filename)
	uploadResult, err := svc.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(targetFile),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(buffer),
		ContentType:        aws.String(contentType),
	})

	if err != nil {
		logrus.Fatalf("Can not read file %s! Error: %s\n", filename, err)
	}

	logrus.Debugf("Successfully uploaded %s to %s [Content-Type: %s]\n", filename, uploadResult.Location, contentType)
	return err
}

func ProcessUploads(terralessConfig schema.TerralessConfig) {
	logrus.Debug("Processing uploads")
	if len(terralessConfig.Uploads) == 0 {
		logrus.Debug("... no uploads")
		return
	}

	for _, upload := range terralessConfig.Uploads {
		logrus.Debug("Processing upload: ", upload)
		provider, _ := terralessConfig.Providers[upload.Provider]
		credentials := credentials2.NewSharedCredentials("", provider.Data["profile"])

		sess, err := session.NewSession(&aws.Config{
			Credentials: credentials,
			Region: aws.String(upload.Region),
		})

		if err != nil {
			logrus.Fatal("Error creating aws session: ", err)
		}

		svc := s3manager.NewUploader(sess)

		uploadedFiles := recursiveUpload(filepath.Join(terralessConfig.SourcePath, upload.Source), upload.Target, upload.Bucket, svc)
		logrus.Debugf("Uploaded files: %s\n", uploadedFiles)
	}
}
