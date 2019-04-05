package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerralessSchema_EnrichWithData_NoOverride(t *testing.T) {
	// given
	data := map[string]string{
		"a": "1",
		"b": "123",
	}
	override := map[string]string{}

	// when
	result := EnrichWithData(data, override)

	// then
	assert.Equal(t, 2, len(result))
	assert.Equal(t, data, result)
}

func TestTerralessSchema_EnrichWithData_OverrideNoDuplicate(t *testing.T) {
	// given
	data := map[string]string{
		"a": "1",
		"b": "123",
	}
	override := map[string]string{
		"c": "456",
	}

	// when
	result := EnrichWithData(data, override)

	// then
	assert.Equal(t, 3, len(result))
}

func TestTerralessSchema_EnrichWithData_OverrideDuplicate(t *testing.T) {
	// given
	data := map[string]string{
		"a": "1",
		"b": "123",
	}
	override := map[string]string{
		"b": "456",
	}

	// when
	result := EnrichWithData(data, override)

	// then
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "456", result["b"])
}

func TestTerralessSchema_EnrichWithData_ProcessString(t *testing.T) {
	// given
	arguments := Arguments{
		Environment: "develop",
		Variables: map[string]string{
			"command": "123",
		},
	}
	settings := TerralessSettings{
		Variables: map[string]string{
			"test1": "asd",
			"test2": "bse",
		},
	}

	// when

	// then
	assert.Equal(t, "asd-123", ProcessString("asd-123", arguments, settings))
	assert.Equal(t, "asd-develop", ProcessString("asd-${environment}", arguments, settings))
	assert.Equal(t, "asd-w", ProcessString("${test1}-w", arguments, settings))
	assert.Equal(t, "bse-develop", ProcessString("${test2}-${environment}", arguments, settings))
	assert.Equal(t, "asd-bse", ProcessString("${test1}-${test2}", arguments, settings))
	assert.Equal(t, "123-bse", ProcessString("${command}-${test2}", arguments, settings))
}

func TestTerralessSchema_EnrichWithData_ProcessData(t *testing.T) {
	// given
	arguments := Arguments{
		Environment: "develop",
	}
	data := map[string]string{
		"field1": "asd-123",
		"field2": "asd-${environment}",
		"field3": "asd-${environment}-${test1}",
	}
	settings := TerralessSettings{
		Variables: map[string]string{
			"test1": "dse",
		},
	}

	// when
	result := ProcessData(data, arguments, settings)

	// then
	expected := map[string]string{
		"field1": "asd-123",
		"field2": "asd-develop",
		"field3": "asd-develop-dse",
	}
	assert.Equal(t, expected, result)
}
