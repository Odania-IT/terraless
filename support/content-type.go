package support

import (
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

// unfortunately mimetype golang module can not handle all files correclty
// js -> file is identified as text/plain
// so we make a simple test on file extension
func DetectContentType(fileName string) string {
	extension := strings.ToLower(filepath.Ext(fileName))
	logrus.Debugf("File %s Extension: %s\n",  fileName, extension)

	if extension == ".gz" {
		return "application/gzip"
	}

	if extension == ".7z" {
		return "application/x-7z-compressed"
	}

	if extension == ".pdf" {
		return "application/pdf"
	}

	if extension == ".xlsx" {
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	if extension == ".docx" {
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}

	if extension == ".pptx" {
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	}

	if extension == ".epub" {
		return "application/epub+zip"
	}

	if extension == ".apk" {
		return "application/vnd.android.package-archive"
	}

	if extension == ".doc" {
		return "application/msword"
	}

	if extension == ".ppt" {
		return "application/vnd.ms-powerpoint"
	}

	if extension == ".xls" {
		return "application/vnd.ms-excel"
	}

	if extension == ".ps" {
		return "application/postscript"
	}

	if extension == ".psd" {
		return "application/x-photoshop"
	}

	if extension == ".ogg" {
		return "application/ogg"
	}

	if extension == ".json" {
		return "application/json"
	}

	if extension == ".html" {
		return "text/html; charset=utf-8"
	}

	if extension == ".php" {
		return "text/x-php; charset=utf-8"
	}

	if extension == ".rtf" {
		return "text/rtf"
	}

	if extension == ".js" {
		return "application/javascript"
	}

	if extension == ".lua" {
		return "text/x-lua"
	}

	if extension == ".pl" {
		return "text/x-perl"
	}

	if extension == ".py" {
		return "application/x-python"
	}

	if extension == ".tcl" {
		return "text/x-tcl"
	}

	if extension == ".svg" {
		return "image/svg+xml"
	}

	if extension == ".x3d" {
		return "model/x3d+xml"
	}

	if extension == ".kml" {
		return "application/vnd.google-earth.kml+xml"
	}

	if extension == ".dae" {
		return "model/vnd.collada+xml"
	}

	if extension == ".gml" {
		return "application/gml+xml"
	}

	if extension == ".gpx" {
		return "application/gpx+xml"
	}

	if extension == ".png" {
		return "image/png"
	}

	if extension == ".jpg" {
		return "image/jpeg"
	}

	if extension == ".gif" {
		return "image/gif"
	}

	if extension == ".webp" {
		return "image/webp"
	}

	if extension == ".tiff" {
		return "image/tiff"
	}

	if extension == ".bmp" {
		return "image/bmp"
	}

	if extension == ".ico" {
		return "image/x-icon"
	}

	if extension == ".mp3" {
		return "audio/mpeg"
	}

	if extension == ".flac" {
		return "audio/flac"
	}

	if extension == ".midi" {
		return "audio/midi"
	}

	if extension == ".ape" {
		return "audio/ape"
	}

	if extension == ".mpc" {
		return "audio/musepack"
	}

	if extension == ".wav" {
		return "audio/wav"
	}

	if extension == ".aiff" {
		return "audio/aiff"
	}

	if extension == ".au" {
		return "audio/basic"
	}

	if extension == ".amr" {
		return "audio/amr"
	}

	if extension == ".mp4" {
		return "video/mp4"
	}

	if extension == ".webm" {
		return "video/webm"
	}

	if extension == ".mpeg" {
		return "video/mpeg"
	}

	if extension == ".mov" {
		return "video/quicktime"
	}

	if extension == ".3gp" {
		return "video/3gp"
	}

	if extension == ".avi" {
		return "video/x-msvideo"
	}

	if extension == ".flv" {
		return "video/x-flv"
	}

	if extension == ".mkv" {
		return "video/x-matroska"
	}

	if extension == ".jar" || extension == ".apk" {
		return "application/jar"
	}

	if extension == ".swf" {
		return "application/x-shockwave-flash"
	}

	if extension == ".crx" {
		return "application/x-chrome-extension"
	}

	if extension == ".css" {
		return "text/css"
	}

	zips := []string{
		".zip", ".xlsx", ".docx", ".pptx", ".pub",
	}
	for _, zipType := range zips {
		if extension == zipType {
			return "application/zip"
		}
	}

	xmls := []string{
		".xml", ".svg", ".x3d", ".kml", ".collada", ".gml", ".gpx",
	}
	for _, xmlType := range xmls {
		if extension == xmlType {
			return "text/xml; charset=utf-8"
		}
	}

	txts := []string{
		".txt", ".php", ".rb", ".lua", ".perl", ".python", ".py", ".rtf", ".tcl",
	}
	for _, txtType := range txts {
		if extension == txtType {
			return "text/plain"
		}
	}

	return "application/octet-stream"
}
