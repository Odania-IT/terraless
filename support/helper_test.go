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
