package terraless_provider_aws

import (
	"bytes"
	"github.com/Odania-IT/terraless/support"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gofrs/flock"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var (
	testAwsConfigFile = filepath.Join(os.TempDir(), "terraless-provider-aws-test-aws-config-file.cfg")
	testAwsCredentialsFile = filepath.Join(os.TempDir(), "terraless-provider-aws-test-credentials-file.cfg")
	testLockFile = filepath.Join(os.TempDir(), "terraless-provider-aws-test.lock")
)

func TestProfileWriter_ProfileWriterConfigsNotExists(t *testing.T) {
	// given
	_ = os.Remove(testAwsConfigFile)
	_ = os.Remove(testAwsCredentialsFile)
	_ = os.Remove(testLockFile)
	awsProfileWriter := AwsProfileWriter{
		awsConfigFile:   testAwsConfigFile,
		credentialsFile: testAwsCredentialsFile,
		lock:            flock.New(testLockFile),
	}
	accessKeyId := "dummy-access-key-must-be-long"
	secretAccessKey := "dummy-secret-access-key"
	sessionToken := "dummy-session-token"
	credentials := sts.Credentials{
		AccessKeyId: &accessKeyId,
		SecretAccessKey: &secretAccessKey,
		SessionToken: &sessionToken,
	}
	targetProfile := "dummy-profile"
	region := "dummy-region"

	// when
	awsProfileWriter.lockAndWriteAwsCredentials(credentials, targetProfile, region)

	// then
	configFileData := support.ReadFile(testAwsConfigFile)
	assert.Contains(t, configFileData, `[dummy-profile]`)
	assert.Contains(t, configFileData, `region = dummy-region`)

	credentialsFileData := support.ReadFile(testAwsCredentialsFile)
	assert.Contains(t, credentialsFileData, `[dummy-profile]`)
	assert.Contains(t, credentialsFileData, `aws_access_key_id     = dummy-access-key-must-be-long`)
	assert.Contains(t, credentialsFileData, `aws_secret_access_key = dummy-secret-access-key`)
	assert.Contains(t, credentialsFileData, `aws_session_token     = dummy-session-token`)

	data := support.ReadFile(testLockFile)
	assert.Equal(t, data, "")
}

func strToBuffer(data string) bytes.Buffer {
	buffer := bytes.Buffer{}
	buffer.WriteString(data)

	return buffer
}

func TestProfileWriter_ProfileWriterConfigsExists(t *testing.T) {
	// given
	cfgData := strToBuffer(`[default]
region = eu-central-1
`)
	support.WriteToFile(testAwsConfigFile, cfgData)

	credentialsData := strToBuffer(`[asd]`)
	support.WriteToFile(testAwsCredentialsFile, credentialsData)
	_ = os.Remove(testLockFile)
	awsProfileWriter := AwsProfileWriter{
		awsConfigFile:   testAwsConfigFile,
		credentialsFile: testAwsCredentialsFile,
		lock:            flock.New(testLockFile),
	}
	accessKeyId := "dummy-access-key-must-be-long"
	secretAccessKey := "dummy-secret-access-key"
	sessionToken := "dummy-session-token"
	credentials := sts.Credentials{
		AccessKeyId: &accessKeyId,
		SecretAccessKey: &secretAccessKey,
		SessionToken: &sessionToken,
	}
	targetProfile := "dummy-profile"
	region := "dummy-region"

	// when
	awsProfileWriter.lockAndWriteAwsCredentials(credentials, targetProfile, region)

	// then
	configFileData := support.ReadFile(testAwsConfigFile)
	assert.Contains(t, configFileData, `[dummy-profile]`)
	assert.Contains(t, configFileData, `region = dummy-region`)

	credentialsFileData := support.ReadFile(testAwsCredentialsFile)
	assert.Contains(t, credentialsFileData, `[dummy-profile]`)
	assert.Contains(t, credentialsFileData, `aws_access_key_id     = dummy-access-key-must-be-long`)
	assert.Contains(t, credentialsFileData, `aws_secret_access_key = dummy-secret-access-key`)
	assert.Contains(t, credentialsFileData, `aws_session_token     = dummy-session-token`)

	data := support.ReadFile(testLockFile)
	assert.Equal(t, data, "")
}

func TestProfileWriter_GetCredentialsFile(t *testing.T) {
	assert.Contains(t, getCredentialsFile(), "/.aws/credentials")
}

func TestProfileWriter_GetAwsConfigFile(t *testing.T) {
	assert.Contains(t, getAwsConfigFile(), "/.aws/config")
}
