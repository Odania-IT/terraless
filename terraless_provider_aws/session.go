package terraless_provider_aws

import (
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gofrs/flock"
	"github.com/sirupsen/logrus"
)

type AwsProfileWriter struct {
	awsConfigFile   string
	credentialsFile string
	lock            *flock.Flock
}

var intermediateProfilesProcessed = map[string]string{}
var execAssumeRoleFunc = execAssumeRole
var retrieveCallerIdentityFunc = retrieveCallerIdentity
var writeSessionProfileFunc = writeSessionProfile

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
			logrus.Infof("Trying auto login for provider %s [intermediate profile: %s]\n", provider.Name, intermediateProfile)
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

	baseProfile := provider.Data["base-profile"]
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

	writeSessionProfileFunc(*awsCredentials, intermediateProfile, region)
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
	svc := sts.New(sessionForProvider(signInProvider))

	logrus.Debugf("Trying to assume role %s\n", arn)
	output, err := execAssumeRoleFunc(svc, sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(getDurationFromData(provider.Data, "session-duration", TargetSessionTokenDuration)),
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(support.SanitizeSessionName(provider.Name)),
	})
	if err != nil {
		logrus.Debugln(provider.Data)
		logrus.Fatalf("[Provider: %s] Failed retrieving session token! Error: %s\n", provider.Name, err)
	}

	profileName := provider.Name
	if provider.Data["profile"] != "" {
		logrus.Debugf("Using profile name from data %s [Provider: %s]\n", provider.Data["profile"], provider.Name)
		profileName = provider.Data["profile"]
	}

	region := provider.Data["region"]
	if region == "" {
		region = "eu-central-1"
	}

	writeSessionProfileFunc(*output.Credentials, profileName, region)
}

func execAssumeRole(svc *sts.STS, input sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	return svc.AssumeRole(&input)
}

func sessionValid(provider schema.TerralessProvider) (bool, error) {
	logrus.Debugf("Checking validity of AWS Provider: %s", provider)
	identity, err := retrieveCallerIdentityFunc(provider)

	if err != nil {
		logrus.Debugf("Invalid AWS Session for provider: %s Error: %s\n", provider.Name, err)
		return false, err
	}

	logrus.Debugf("Valid AWS Session for provider: %s User: %s\n", provider.Name, identity)
	return true, nil
}

func retrieveCallerIdentity(provider schema.TerralessProvider) (*sts.GetCallerIdentityOutput, error) {
	svc := sts.New(sessionForProvider(provider))
	return svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
}

func sessionForProvider(provider schema.TerralessProvider) *session.Session {
	profileName := provider.Data["profile"]
	if profileName == "" {
		profileName = provider.Name
	}

	currentCredentials := credentials.NewSharedCredentials("", profileName)
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

func simpleSession(provider schema.TerralessProvider) *session.Session {
	config := aws.Config{
		Region: aws.String(provider.Data["region"]),
	}

	sess, err := session.NewSession(&config)

	if err != nil {
		logrus.Fatalf("Failed creating AWS Session for provider: %s Error: %s\n", provider, err)
	}

	return sess
}
