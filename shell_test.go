package main

import (
	"os"
	"testing"
)

func TestTerraless_ExecuteCommand(t *testing.T) {
	// given
	args := []string{
		"HelloTerraless",
	}

	// when
	executeCommand(os.TempDir(), "echo", args, false)

	// then
}
