package support

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestTerralessSupport_DetectContentType_CommonTypes(t *testing.T) {
	properties := map[string]string {
		"/tmp/test/dummy.jpg": "image/jpeg",
		"/tmp/test/asd.html": "text/html; charset=utf-8",
		"/tmp/test/archive.zip": "application/zip",
		"/tmp/test/dummy.txt": "text/plain",
		"/tmp/test/dummy.zzz": "application/octet-stream",
	}

	for fileName, expectedContentType := range properties {
		contentType := DetectContentType(fileName)
		assert.Equal(t, contentType, expectedContentType)
	}
}
