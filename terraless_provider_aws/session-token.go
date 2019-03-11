package terraless_provider_aws

import (
	"bufio"
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

const (
	IntermediateSessionTokenDuration = int64(12 * 60 * 60)
	TargetSessionTokenDuration       = int64(60 * 60)
)

func askForTokenCode(tokenSerialNumber string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter mfa token for %s: ", tokenSerialNumber)

	tokenCode, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatalf("Error reading MFA input! Error: %s\n", err)
	}

	return tokenCode
}

func getTokenCode(mfaArn string) string {
	tokenCode := askForTokenCode(mfaArn)

	return strings.Trim(tokenCode, " \r\n")
}

func getIntermediateSessionToken(provider schema.TerralessProvider) *sts.Credentials {
	logrus.Debugf("Retrieving session for AWS Provider: %s", provider)
	svc := sts.New(sessionForProfile(provider))

	mfaDevice := provider.Data["mfa-device"]
	getSessionTokenInput := sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(getDurationFromData(provider.Data, "intermediate-session-duration", IntermediateSessionTokenDuration)),
	}

	if mfaDevice != "" {
		tokenCode := getTokenCode(mfaDevice)
		getSessionTokenInput.SerialNumber = aws.String(mfaDevice)
		getSessionTokenInput.TokenCode = aws.String(tokenCode)
	}

	logrus.Debugln(getSessionTokenInput)
	output, err := svc.GetSessionToken(&getSessionTokenInput)
	if err != nil {
		logrus.Fatalf("[getIntermediateSessionToken] Failed retrieving session token! Error: %s\n", err)
	}

	return output.Credentials
}

func getDurationFromData(data map[string]string, key string, defaultValue int64) int64 {
	val := data[key]

	if val == "" {
		return defaultValue
	}

	parsedInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		logrus.Errorf("Invalid value for %s specified! Please specify a int! Using default value now... Error: %s\n", key, err)
		return defaultValue
	}

	return parsedInt
}
