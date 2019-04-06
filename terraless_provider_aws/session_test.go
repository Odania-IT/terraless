package terraless_provider_aws

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/assert"
	"testing"
)

var preparedSessions []string
var callerIdentityFor = map[string]bool{}
var writeSessionProfileCalls int
func execAssumeRoleFuncMock(svc *sts.STS, input sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	preparedSessions = append(preparedSessions, *input.RoleSessionName)

	response := sts.AssumeRoleOutput{
		Credentials: &sts.Credentials{},
	}
	return &response, nil
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func retrieveCallerIdentityMock(provider schema.TerralessProvider) (*sts.GetCallerIdentityOutput, error) {
	if !callerIdentityFor[provider.Name] {
		callerIdentityFor[provider.Name] = true
		return nil, &errorString{s: "mock-error"}
	}

	response := sts.GetCallerIdentityOutput{
		Account: &provider.Name,
	}
	return &response, nil
}

func writeSessionProfileMock(credentials sts.Credentials, targetProfile string, region string) {
	writeSessionProfileCalls += 1
}

func TestSession_SessionForProfile(t *testing.T) {
	// given
	provider := schema.TerralessProvider{
		Name: "DummyProvider",
		Data: map[string]string{
			"region": "dummy-region",
		},
	}

	// when
	session := sessionForProvider(provider)

	// then
	assert.Equal(t, `dummy-region`, *session.Config.Region)
}

func TestSession_SessionValid(t *testing.T) {
	// given
	provider := schema.TerralessProvider{
		Name: "DummyProvider",
		Data: map[string]string{
			"region": "dummy-region",
		},
	}

	// when
	valid, err := sessionValid(provider)

	// then
	assert.False(t, valid)
	assert.Contains(t, err.Error(), "SharedCredsLoad: failed to")
}

func TestSession_PrepareSession(t *testing.T) {
	// given
	execAssumeRoleFunc = execAssumeRoleFuncMock
	retrieveCallerIdentityFunc = retrieveCallerIdentityMock
	writeSessionProfileFunc = writeSessionProfileMock
	terralessConfig := schema.TerralessConfig{
		Providers: map[string]schema.TerralessProvider{
			"provider1": {
				Type: "aws",
				Name: "provider1-profile",
				Data: map[string]string{
					"accountId": "123",
					"session-duration": "10",
					"role": "developer",
				},
			},
			"provider2": {
				Type: "aws",
				Name: "provider2-profile",
				Data: map[string]string{
					"accountId": "456",
					"profile": "my-second-profile",
					"role": "admin",
				},
			},
		},
		Settings: schema.TerralessSettings{
			AutoSignIn: true,
		},
	}

	// when
	prepareSession(terralessConfig)

	// then
	assert.Equal(t, 2, len(preparedSessions))
	assert.Equal(t, 2, writeSessionProfileCalls)
}
