package support

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func WriteToFile(targetFileName string, buffer bytes.Buffer) {
	targetFile, err := os.Create(targetFileName)
	if err != nil {
		logrus.Fatal("Failed creating file: ", filepath.Base(targetFileName), err)
	}

	_, err = targetFile.Write(buffer.Bytes())
	if err != nil {
		logrus.Fatal("Failed writing to file: ", filepath.Base(targetFileName), err)
	}

	_ = targetFile.Close()
}
