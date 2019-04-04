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
		"/tmp/asd.gz": "application/gzip",
		"/tmp/asd.7z": "application/x-7z-compressed",
		"/tmp/asd.pdf": "application/pdf",
		"/tmp/asd.epub": "application/epub+zip",
		"/tmp/asd.apk": "application/vnd.android.package-archive",
		"/tmp/asd.doc": "application/msword",
		"/tmp/asd.json": "application/json",
		"/tmp/asd.php": "text/x-php; charset=utf-8",
		"/tmp/asd.lua": "text/x-lua",
		"/tmp/asd.pl": "text/x-perl",
		"/tmp/asd.py": "application/x-python",
		"/tmp/asd.tcl": "text/x-tcl",
		"/tmp/asd.svg": "image/svg+xml",
		"/tmp/asd.x3d": "model/x3d+xml",
		"/tmp/asd.kml": "application/vnd.google-earth.kml+xml",
		"/tmp/asd.png": "image/png",
		"/tmp/asd.gif": "image/gif",
		"/tmp/asd.webp": "image/webp",
		"/tmp/asd.tiff": "image/tiff",
		"/tmp/asd.bmp": "image/bmp",
		"/tmp/asd.ico": "image/x-icon",
		"/tmp/asd.mp3": "audio/mpeg",
		"/tmp/asd.flac": "audio/flac",
		"/tmp/asd.midi": "audio/midi",
		"/tmp/asd.ape": "audio/ape",
		"/tmp/asd.mpc": "audio/musepack",
		"/tmp/asd.wav": "audio/wav",
		"/tmp/asd.aiff": "audio/aiff",
		"/tmp/asd.au": "audio/basic",
		"/tmp/asd.amr": "audio/amr",
		"/tmp/asd.mp4": "video/mp4",
		"/tmp/asd.webm": "video/webm",
		"/tmp/asd.mpeg": "video/mpeg",
		"/tmp/asd.mov": "video/quicktime",
		"/tmp/asd.3gp": "video/3gp",
		"/tmp/asd.avi": "video/x-msvideo",
		"/tmp/asd.flv": "video/x-flv",
		"/tmp/asd.mkv": "video/x-matroska",
		"/tmp/asd.jar": "application/jar",
		"/tmp/asd.swf": "application/x-shockwave-flash",
		"/tmp/asd.crx": "application/x-chrome-extension",
		"/tmp/asd.css": "text/css",
	}

	for fileName, expectedContentType := range properties {
		contentType := DetectContentType(fileName)
		assert.Equal(t, contentType, expectedContentType)
	}
}
