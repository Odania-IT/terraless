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
	reader1 := strings.NewReader("n")
	reader2 := strings.NewReader("yn")
	reader3 := strings.NewReader("sed")

	// when
	result1 := checkApprove(reader1)
	result2 := checkApprove(reader2)
	result3 := checkApprove(reader3)

	// then
	assert.Equal(t, false, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, false, result3)
}

func TestTerraless_CheckApprove_Ok(t *testing.T) {
	// given
	reader1 := strings.NewReader("y")
	reader2 := strings.NewReader("Y")

	// when
	result1 := checkApprove(reader1)
	result2 := checkApprove(reader2)

	// then
	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
}
