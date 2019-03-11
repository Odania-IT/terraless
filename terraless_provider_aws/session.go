package terraless_provider_aws

import (
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"github.com/gofrs/flock"
	"github.com/sirupsen/logrus"
	"os"
	"os/user"
	"path/filepath"
)

type AwsProfileWriter struct {
	awsConfigFile   string
	credentialsFile string
	lock            *flock.Flock
}

var intermediateProfilesProcessed = map[string]string{}

func prepareSession(terralessConfig schema.TerralessConfig) {
	for _, provider := range terralessConfig.Providers {
		if canHandle(provider.Type) {
			logrus.Debugf("Found AWS Provider: %s\n", provider.Name)

			intermediateProfile := processIntermediateProfile(provider, terralessConfig.Settings.AutoSignIn)

			verifyOrUpdateSession(provider, intermediateProfile, terralessConfig.Settings.AutoSignIn)
		}
	}
}

func processIntermediateProfile(provider schema.TerralessProvider, autoSignIn bool) string {
	intermediateProfile := provider.Data["intermediate-profile"]

	if intermediateProfilesProcessed[provider.Name] == "" {
		if intermediateProfile == "" {
			logrus.Debug("No intermediate profile! Using default....")
			intermediateProfile = "terraless-session"
		}

		if autoSignIn {
			validateOrRefreshIntermediateSession(provider, intermediateProfile)
		}

		intermediateProfilesProcessed[provider.Name] = intermediateProfile
	}

	return intermediateProfilesProcessed[provider.Name]
}

func verifyOrUpdateSession(provider schema.TerralessProvider, intermediateProfile string, autoSignIn bool) {
	logrus.Debugf("Checking provider %s\n", provider)
	validSession, err := sessionValid(provider)
	if !validSession {
		if autoSignIn {
			logrus.Infof("Trying auto login for provider %s\n", provider.Name)
			assumeRole(intermediateProfile, provider)
			validSession, err = sessionValid(provider)
		}

		if !validSession {
			logrus.Fatalf("No AWS Session for provider: %s [Error: %s]\n", provider.Name, err)
		}
	}
}

func validateOrRefreshIntermediateSession(provider schema.TerralessProvider, intermediateProfile string) {
	mfaDevice := provider.Data["mfa-device"]

	if mfaDevice == "" {
		logrus.Debug("No mfa-device! Nothing to do....")
		return
	}

	region := provider.Data["region"]
	if region == "" {
		region = "eu-central-1"
	}

	baseProfile := provider.Data["profile"]
	if baseProfile == "" {
		baseProfile = "default"
	}
	logrus.Debugf("Creating intermediate profile session. Region: %s IntermediateProfile: %s BaseProfile: %s\n",
		region, intermediateProfile, baseProfile)

	intermediateProvider := schema.TerralessProvider{
		Name: intermediateProfile,
		Data: map[string]string{
			"mfa-device":  mfaDevice,
			"region":  region,
			"profile": intermediateProfile,
		},
	}
	validSession, err := sessionValid(intermediateProvider)
	if err == nil && validSession {
		logrus.Debug("Intermediate session still valid....")
		return
	}

	// Retrieve session token for base profile in order to store it as intermediate profile
	intermediateProvider.Data["profile"] = baseProfile
	awsCredentials := getIntermediateSessionToken(intermediateProvider)
	logrus.Debug(awsCredentials)

	writeSessionProfile(*awsCredentials, intermediateProfile, region)
}

func writeSessionProfile(credentials sts.Credentials, targetProfile string, region string) {
	awsProfileWriter := AwsProfileWriter{
		awsConfigFile:   getAwsConfigFile(),
		credentialsFile: getCredentialsFile(),
		lock:            flock.New(filepath.Join(os.TempDir(), "terraless-provider-aws.lock")),
	}

	awsProfileWriter.lockAndWriteAwsCredentials(credentials, targetProfile, region)
}

func assumeRole(intermediateProfile string, provider schema.TerralessProvider) {
	accountId := provider.Data["accountId"]
	role := provider.Data["role"]

	if accountId == "" || role == "" {
		logrus.Fatalf("Can not assume role without accountId and role! Provider: %s Data: %s\n", provider.Name, provider.Data)
	}

	arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, role)
	signInProvider := schema.TerralessProvider{
		Name: intermediateProfile,
		Data: map[string]string{
			"profile": intermediateProfile,
		},
	}
	svc := sts.New(sessionForProfile(signInProvider))

	logrus.Debugf("Trying to assume role %s\n", arn)
	output, err := svc.AssumeRole(&sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(getDurationFromData(provider.Data, "session-duration", TargetSessionTokenDuration)),
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(provider.Name),
	})
	if err != nil {
		logrus.Debugln(provider.Data)
		logrus.Fatalf("[Provider: %s] Failed retrieving session token! Error: %s\n", provider.Name, err)
	}

	profileName := provider.Name
	if provider.Data["profile"] != "" {
		profileName = provider.Data["profile"]
	}

	region := provider.Data["region"]
	if region == "" {
		region = "eu-central-1"
	}

	writeSessionProfile(*output.Credentials, profileName, region)
}

func (pw AwsProfileWriter) lockAndWriteAwsCredentials(credentials sts.Credentials, targetProfile string, region string) {
	defer pw.lock.Unlock()

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

func writeKeyToSection(section *ini.Section, key string, val string) {
	_, err := section.NewKey(key, val)

	if err != nil {
		logrus.Fatalf("Failed writting key %s to aws profile section\n", err)
	}
}

func (pw AwsProfileWriter) writeAwsConfig(region string, targetProfile string) {
	logrus.Debugf("Loading config file %s\n", pw.awsConfigFile)
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

func sessionValid(provider schema.TerralessProvider) (bool, error) {
	logrus.Debugf("Checking validity of AWS Provider: %s", provider)
	svc := sts.New(sessionForProfile(provider))
	identity, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		logrus.Debugf("Invalid AWS Session for provider: %s Error: %s\n", provider.Name, err)
		return false, err
	}

	logrus.Debugf("Valid AWS Session for provider: %s User: %s\n", provider.Name, identity)
	return true, nil
}

func sessionForProfile(provider schema.TerralessProvider) *session.Session {
	currentCredentials := credentials.NewSharedCredentials("", provider.Data["profile"])
	config := aws.Config{
		Credentials: currentCredentials,
		Region:      aws.String(provider.Data["region"]),
	}

	logrus.Debugf("AWS Session Profile for config %s\n", provider.Data)
	sess, err := session.NewSession(&config)

	if err != nil {
		logrus.Fatalf("Failed creating AWS Session for provider: %s Error: %s\n", provider, err)
	}

	return sess
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
