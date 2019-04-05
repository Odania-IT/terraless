package templates

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestTerralessTemplates_RenderTemplateToBuffer_Render(t *testing.T) {
	// given
	terralessData := schema.TerralessData{}

	// when
	buffer := Render(&terralessData, bytes.Buffer{})

	// then
	assert.Equal(t, buffer.String(), "")
}

func TestTerralessTemplates_RenderTemplateToBuffer_Simple(t *testing.T) {
	// given
	data := map[string]string{}
	tpl := `Template`

	// when
	buffer := RenderTemplateToBuffer(data, bytes.Buffer{}, tpl, "test")

	// then
	assert.Equal(t, buffer.String(), tpl)
}

func TestTerralessTemplates_RenderTemplateToBuffer_WithData(t *testing.T) {
	// given
	data := map[string]string{
		"Content": "TemplateRendered",
	}
	tpl := `{{ .Content }}`

	// when
	buffer := RenderTemplateToBuffer(data, bytes.Buffer{}, tpl, "test")

	// then
	assert.Equal(t, buffer.String(), data["Content"])
}
