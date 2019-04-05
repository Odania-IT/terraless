package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func captureOutputProcessCommand(terralessData schema.TerralessData, kingpinResult string) string {
	oldStdout := os.Stdout
	readFile, writeFile, _ := os.Pipe()
	os.Stdout = writeFile

	print()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, readFile)
		outC <- buf.String()
	}()

	processCommands(&terralessData, kingpinResult)

	_ = writeFile.Close()
	os.Stdout = oldStdout
	out := <-outC

	return out
}

func TestMain_InfoCommand(t *testing.T) {
	// given
	terralessData := schema.TerralessData{}
	kingpinResult := infoCommand.FullCommand()

	// when
	output := captureOutputProcessCommand(terralessData, kingpinResult)

	// then
	assert.Contains(t, output, "Terraless Version: " + VERSION + " [Codename: " + CODENAME + "]")
}
