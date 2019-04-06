package terraless_provider_aws

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSessionToken_GetDurationFromData(t *testing.T) {
	// given
	data := map[string]string{
		"key1": "123",
		"notNumeric": "a",
	}

	// when

	// then
	assert.Equal(t, int64(2), getDurationFromData(data, "key", 2))
	assert.Equal(t, int64(123), getDurationFromData(data, "key1", 2))
	assert.Equal(t, int64(2), getDurationFromData(data, "notNumeric", 2))
}

func TestSessionToken_GetTokenCode(t *testing.T) {
	// given
	reader1 := strings.NewReader("serial1\n")

	// when
	code := getTokenCode("myMfa", reader1)

	// then
	assert.Equal(t, "serial1", code)
}
