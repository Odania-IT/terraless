package terraless_provider_aws

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
