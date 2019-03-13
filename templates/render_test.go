package templates

import (
	"bytes"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestTerralessTemplates_RenderTemplateToBuffer_Simple(t *testing.T) {
	// given
	data := map[string]string{}
	tpl := `Template`

	// when
	buffer := RenderTemplateToBuffer(data, bytes.Buffer{}, tpl)

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
	buffer := RenderTemplateToBuffer(data, bytes.Buffer{}, tpl)

	// then
	assert.Equal(t, buffer.String(), data["Content"])
}
