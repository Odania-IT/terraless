package templates

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func dummyTerralessRenderProvider() schema.Provider {
	return schema.Provider{
		CanHandle: func(resourceType string) bool {
			return resourceType == "dummy"
		},
		FinalizeTemplates: func(terralessData schema.TerralessData) string {
			return ""
		},
		PrepareSession: func(terralessConfig schema.TerralessConfig) {
		},
		ProcessUpload: func(terralessData schema.TerralessData, upload schema.TerralessUpload) []string {
			return []string{}
		},
		RenderAuthorizerTemplates: func(config schema.TerralessConfig) string {
			return ""
		},
		RenderCertificateTemplates: func(config schema.TerralessConfig) string {
			return ""
		},
		RenderEndpointTemplates: func(config schema.TerralessConfig) string {
			return ""
		},
		RenderUploadTemplates: func(terralessData schema.TerralessData) string {
			return ""
		},
	}
}

func TestTerralessTemplates_Render(t *testing.T) {
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
	providers:= []schema.Provider{
		dummyTerralessRenderProvider(),
	}

	// when
	buffer := Render(&terralessData, providers, bytes.Buffer{})

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

func TestTerralessTemplates_AwsProviderKeys_WithProfileInData(t *testing.T) {
	// given
	data := map[string]string{
		"alias": "myAlias",
		"profile": "myProfile",
		"region": "myRegion",
	}

	// when
	result := awsProviderKeys(data, "profile1")

	// then
	assert.Equal(t, "myAlias", result["alias"])
	assert.Equal(t, "myProfile", result["profile"])
	assert.Equal(t, "myRegion", result["region"])
}

func TestTerralessTemplates_AwsProviderKeys_WithoutProfileInData(t *testing.T) {
	// given
	data := map[string]string{
		"alias": "myAlias",
		"region": "myRegion",
	}

	// when
	result := awsProviderKeys(data, "profile1")

	// then
	assert.Equal(t, "myAlias", result["alias"])
	assert.Equal(t, "profile1", result["profile"])
	assert.Equal(t, "myRegion", result["region"])
}
