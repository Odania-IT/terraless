package support

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/user"
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


func WriteToFileIfNotExists(file string, data string) {
	_, err := os.Stat(file)

	if err != nil {
		buffer := bytes.Buffer{}
		buffer.WriteString(data)
		WriteToFile(file, buffer)
	}
}

func ReadFile(file string) string {
	content, err := ioutil.ReadFile(file)

	if err != nil {
		logrus.Fatalf("Could not read file %s\n", file)
	}

	return string(content)
}

func HomeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		logrus.Warnf("Could not detect user home folder")
		return ""
	}

	return usr.HomeDir
}
