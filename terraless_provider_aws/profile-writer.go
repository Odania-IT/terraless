package terraless_provider_aws

import (
	"github.com/Odania-IT/terraless/support"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"github.com/gofrs/flock"
	"github.com/sirupsen/logrus"
	"os"
	"os/user"
	"path/filepath"
)


func writeSessionProfile(credentials sts.Credentials, targetProfile string, region string) {
	awsProfileWriter := AwsProfileWriter{
		awsConfigFile:   getAwsConfigFile(),
		credentialsFile: getCredentialsFile(),
		lock:            flock.New(filepath.Join(os.TempDir(), "terraless-provider-aws.lock")),
	}

	awsProfileWriter.lockAndWriteAwsCredentials(credentials, targetProfile, region)
}

func (pw AwsProfileWriter) lockAndWriteAwsCredentials(credentials sts.Credentials, targetProfile string, region string) {
	defer func() {
		cerr := pw.lock.Unlock()
		if cerr != nil {
			logrus.Fatalf("[Provider: AWS] Failed to unlock Error %s\n", cerr)
		}
	}()

	locked, err := pw.lock.TryLock()

	if err != nil {
		logrus.Fatalf("Failed aquiring lock for updating AWS credentials! %s\n", err)
	}

	if locked {
		pw.writeAwsCredentials(credentials, targetProfile)
		pw.writeAwsConfig(region, targetProfile)

		return
	}

	logrus.Fatal("AWS credentials lock already locked!")
}

func (pw AwsProfileWriter) writeAwsCredentials(credentials sts.Credentials, targetProfile string) {
	logrus.Debugf("Loading credentials file %s\n", pw.credentialsFile)
	support.WriteToFileIfEmpty(pw.credentialsFile, "[default]")
	cfg, err := ini.Load(pw.credentialsFile)

	if err != nil {
		logrus.Fatalf("Error loading aws credentials! %s\n", err)
	}

	section := cfg.Section(targetProfile)

	if section == nil {
		section, err = cfg.NewSection(targetProfile)

		if err != nil {
			logrus.Fatalf("Failed creating section in aws credentials file! Error: %s\n", err)
		}
	} else {
		section.DeleteKey("aws_access_key_id")
		section.DeleteKey("aws_secret_access_key")
		section.DeleteKey("aws_session_token")
	}

	writeKeyToSection(section, "aws_access_key_id", *credentials.AccessKeyId)
	writeKeyToSection(section, "aws_secret_access_key", *credentials.SecretAccessKey)
	writeKeyToSection(section, "aws_session_token", *credentials.SessionToken)

	err = cfg.SaveTo(pw.credentialsFile)
	if err != nil {
		logrus.Fatalf("Failed writing config file %s! Error: %s\n", pw.credentialsFile, err)
	}

	logrus.Debugf("Wrote session token for profile %s\n", targetProfile)
	logrus.Debugf("Token is valid until: %v\n", credentials.Expiration)
}

func (pw AwsProfileWriter) writeAwsConfig(region string, targetProfile string) {
	logrus.Debugf("Loading config file %s\n", pw.awsConfigFile)
	support.WriteToFileIfEmpty(pw.awsConfigFile, "[default]")
	cfg, err := ini.Load(pw.awsConfigFile)

	if err != nil {
		logrus.Fatalf("Error loading aws config! %s\n", err)
	}

	section := cfg.Section(targetProfile)

	if section == nil {
		section, err = cfg.NewSection(targetProfile)

		if err != nil {
			logrus.Fatalf("Failed creating section in aws config file! Error: %s\n", err)
		}
	} else {
		section.DeleteKey("region")
	}

	writeKeyToSection(section, "region", region)

	err = cfg.SaveTo(pw.awsConfigFile)
	if err != nil {
		logrus.Fatalf("Failed writing config file %s! Error: %s\n", pw.credentialsFile, err)
	}

	logrus.Debugf("Wrote aws config section for profile %s\n", targetProfile)
}

func writeKeyToSection(section *ini.Section, key string, val string) {
	_, err := section.NewKey(key, val)

	if err != nil {
		logrus.Fatalf("Failed writting key %s to aws profile section\n", err)
	}
}

func getCredentialsFile() string {
	credentialsPath := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")

	if credentialsPath != "" {
		return credentialsPath
	}

	usr, err := user.Current()
	if err != nil {
		logrus.Fatalf("Error fetching home dir: %s", err)
	}

	return filepath.Join(usr.HomeDir, ".aws", "credentials")
}

func getAwsConfigFile() string {
	usr, err := user.Current()
	if err != nil {
		logrus.Fatalf("Error fetching home dir: %s", err)
	}

	return filepath.Join(usr.HomeDir, ".aws", "config")
}
