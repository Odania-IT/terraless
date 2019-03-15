package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func handleExecuteCommandErr(err error, allowFailure bool) {
	if err != nil {
		if allowFailure {
			logrus.Debug("Failure executing command, but allowFailure set!", err)
			return
		}

		logrus.Fatal("Could not execute command: ", err)
	}
}

func executeCommand(folder string, command string, arguments []string, allowFailure bool) {
	logrus.Debug("Executing command ", command, " ", arguments, " in: ", folder)
	cmd := exec.Command(command, arguments...)
	cmd.Dir = folder
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logrus.Fatal("Error executing command: ", err)
	}

	err := cmd.Wait()
	handleExecuteCommandErr(err, allowFailure)
}

func checkApprove() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Deploy? (y)")
	input, _ := reader.ReadString('\n')
	input = strings.Trim(input, "\n")

	if input == "y" || input == "Y" {
		return true
	}

	return false
}
