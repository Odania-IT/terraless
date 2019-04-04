package support

import (
	"bytes"
	"github.com/magiconair/properties/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestTerralessSupport_WriteToFile(t *testing.T) {
	// given
	testContent := "This is our test content"
	testFile := filepath.Join(os.TempDir(), "terraless-aws-provider-test-file.txt")
	buffer := bytes.Buffer{}
	buffer.WriteString(testContent)

	// when
	WriteToFile(testFile, buffer)

	// then
	fileContent := ReadFile(testFile)
	assert.Equal(t, testContent, fileContent)
}

func TestTerralessSupport_WriteToFileIfNotExists(t *testing.T) {
	// given
	testContent := "This is our test content"
	testFile := filepath.Join(os.TempDir(), "terraless-aws-provider-test-file.txt")

	// when
	WriteToFileIfNotExists(testFile, testContent)

	// then
	fileContent := ReadFile(testFile)
	assert.Equal(t, testContent, fileContent)
}
