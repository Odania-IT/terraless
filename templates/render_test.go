package templates

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestTerralessTemplates_RenderTemplateToBuffer_Render(t *testing.T) {
	// given
	terralessData := schema.TerralessData{
		Config: schema.TerralessConfig{
			Authorizers: map[string]schema.TerralessAuthorizer{
				"auth": {
					Type: "dummy",
				},
			},
			Certificates: map[string]schema.TerralessCertificate{
				"cert": {
					Type: "dummy",
				},
			},
			Endpoints: []schema.TerralessEndpoint{
				{
					Type: "dummy",
				},
			},
			Functions: map[string]schema.TerralessFunction{
				"func": {
					Type: "dummy",
				},
			},
			Package: schema.TerralessPackage{
				SourceDir: "dummy-source",
			},
			Uploads: []schema.TerralessUpload{
				{
					Type: "dummy",
				},
			},
		},
	}

	// when
	buffer := Render(&terralessData, bytes.Buffer{})

	// then
	expected := `## Terraless: Lambda Package

data "archive_file" "lambda-archive" {
  source_dir = "${path.root}/dummy-source"

  output_path = "lambda.zip"
  type = "zip"
}

`
	assert.Equal(t, buffer.String(), expected)
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
