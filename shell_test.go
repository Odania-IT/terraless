package main

import (
	"github.com/magiconair/properties/assert"
	"os"
	"strings"
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

func TestTerraless_CheckApprove_Fail(t *testing.T) {
	// given
	reader := strings.NewReader("n")

	// when
	result := checkApprove(reader)

	// then
	assert.Equal(t, false, result)
}

func TestTerraless_CheckApprove_Ok(t *testing.T) {
	// given
	reader := strings.NewReader("y")

	// when
	result := checkApprove(reader)

	// then
	assert.Equal(t, true, result)
}
