package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
)

var testProcessed map[string]bool

func dummyTerralessProvider() schema.Provider {
	return schema.Provider{
		CanHandle: func(resourceType string) bool {
			return resourceType == "dummy"
		},
		Name: func() string {
			return "terraless-provider-dummy"
		},
		PrepareSession: func(terralessConfig schema.TerralessConfig) {
			testProcessed["PrepareSession"] = true
		},
		ProcessUpload: func(terralessData schema.TerralessData, upload schema.TerralessUpload) {
			testProcessed["ProcessUpload"] = true
		},
		RenderAuthorizerTemplates: func(config schema.TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
			testProcessed["RenderAuthorizerTemplates"] = true
			return bytes.Buffer{}
		},
		RenderCertificateTemplates: func(config schema.TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
			testProcessed["RenderCertificateTemplates"] = true
			return bytes.Buffer{}
		},
		RenderEndpointTemplates: func(config schema.TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
			testProcessed["RenderEndpointTemplates"] = true
			return bytes.Buffer{}
		},
		RenderFunctionTemplates: func(resourceType string, functionEvents schema.FunctionEvents, terralessData *schema.TerralessData, buffer bytes.Buffer) bytes.Buffer {
			testProcessed["RenderFunctionTemplates"] = true
			return bytes.Buffer{}
		},
		RenderUploadTemplates: func(currentConfig schema.TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
			testProcessed["RenderUploadTemplates"] = true
			return bytes.Buffer{}
		},
	}
}

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
	assert.Contains(t, output, "Terraless Version: "+VERSION+" [Codename: "+CODENAME+"]")
}

func TestMain_DetectTerralessProviders(t *testing.T) {
	// given

	// when
	providers := detectTerralessProviders()

	// then
	assert.Equal(t, 1, len(providers))
	assert.Equal(t, "terraless-provider-aws", providers[0].Name())
}

func TestMain_Deploy(t *testing.T) {
	// given
	terralessData := schema.TerralessData{
		ActiveProviders: map[string]schema.TerralessProvider{},
		Arguments: schema.Arguments{
			Config:           filepath.Join(os.TempDir(), "terraless-provider-aws-test", "my-project-config.yml"),
			Environment:      "test",
			ForceDeploy:      true,
			GlobalConfig:     filepath.Join(os.TempDir(), "terraless-provider-aws-test", "my-global-config.yml"),
			TerraformCommand: "echo",
		},
		Config: schema.TerralessConfig{
			Authorizers: map[string]schema.TerralessAuthorizer{
				"dummyAuthorizer": {
					Name: "dummyAuthorizer",
					Type: "dummy",
				},
			},
			Backend: schema.TerralessBackend{
				Name: "dummyBackend",
				Type: "dummy",
			},
			Certificates: map[string]schema.TerralessCertificate{
				"dummyCertificate": {
					Type:   "dummy",
					Domain: "dummy-domain.local",
				},
			},
			Endpoints: []schema.TerralessEndpoint{
				{
					Type: "dummy",
				},
			},
			Functions: map[string]schema.TerralessFunction{
				"DummyFunc": schema.TerralessFunction{
					Type: "dummy",
					Events: []schema.TerralessFunctionEvent{
						{
							Type:       "dummy",
							Authorizer: "dummyAuthorizer",
						},
					},
				},
			},
			ProjectName: "DummyProject",
			Providers: map[string]schema.TerralessProvider{
				"DummyProvider1": {
					Name: "DummyProvider1",
					Type: "dummy",
					Roles: []string{
						"role1",
					},
				},
				"DummyProvider2": {
					Name: "DummyProvider2",
					Type: "dummy",
					Roles: []string{
						"role2",
					},
				},
			},
			Package: schema.TerralessPackage{
				SourceDir: filepath.Join(os.TempDir(), "terraless-provider-aws-test", "src"),
			},
			Settings: schema.TerralessSettings{
				AutoSignIn: true,
			},
			Uploads: []schema.TerralessUpload{
				{
					Type: "dummyUpload",
				},
			},
		},
		TerralessProviders: []schema.Provider{
			dummyTerralessProvider(),
		},
	}
	kingpinResult := deployCommand.FullCommand()
	testProcessed = map[string]bool{}

	_ = os.Mkdir(filepath.Join(os.TempDir(), "terraless-provider-aws-test"), 0755)

	// when
	output := captureOutputProcessCommand(terralessData, kingpinResult)

	// then
	assert.Contains(t, output, "apply -input=false terraform.plan")
	assert.Equal(t, true, testProcessed["PrepareSession"])
	assert.Equal(t, true, testProcessed["ProcessUpload"])
	assert.Equal(t, true, testProcessed["RenderAuthorizerTemplates"])
	assert.Equal(t, true, testProcessed["RenderCertificateTemplates"])
	assert.Equal(t, true, testProcessed["RenderEndpointTemplates"])
	assert.Equal(t, true, testProcessed["RenderFunctionTemplates"])
	assert.Equal(t, true, testProcessed["RenderUploadTemplates"])
}