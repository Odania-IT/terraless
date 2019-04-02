package support

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestTerralessSupport_ContainsStartsWith(t *testing.T) {
	// given
	haystack := []string {
		"Dummy",
		"ExampleCom",
	}

	// when

	// then
	assert.Equal(t, ContainsStartsWith(haystack, "ExampleCom"), true)
	assert.Equal(t, ContainsStartsWith(haystack, "Example"), true)
	assert.Equal(t, ContainsStartsWith(haystack, "A"), false)
	assert.Equal(t, ContainsStartsWith(haystack, "ample"), false)
	assert.Equal(t, ContainsStartsWith(haystack, "Com"), false)
}

func TestTerralessSupport_SanitizeString(t *testing.T) {
	// given

	// when

	// then
	assert.Equal(t, SanitizeString("ExampleCom"), "ExampleCom")
	assert.Equal(t, SanitizeString("Example.Com"), "Example-Com")
	assert.Equal(t, SanitizeString("Example-Com"), "Example-Com")
}

func TestTerralessSupport_SanitizeSessionName(t *testing.T) {
	// given

	// when

	// then
	assert.Equal(t, SanitizeSessionName("ExampleCom"), "ExampleCom")
	assert.Equal(t, SanitizeSessionName("Example.Com"), "Example.Com")
	assert.Equal(t, SanitizeSessionName("Example/Com"), "Example-Com")
}

func TestTerralessSupport_Contains(t *testing.T) {
	// given
	haystack := []string {
		"Dummy",
		"ExampleCom",
	}

	// when

	// then
	assert.Equal(t, Contains(haystack, "ExampleCom"), true)
	assert.Equal(t, Contains(haystack, "Example"), false)
	assert.Equal(t, Contains(haystack, "A"), false)
	assert.Equal(t, Contains(haystack, "ample"), false)
	assert.Equal(t, Contains(haystack, "Com"), false)
	assert.Equal(t, Contains(haystack, "Dummy"), true)
}
