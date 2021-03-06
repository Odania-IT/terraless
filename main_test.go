package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/dummy"
	"github.com/Odania-IT/terraless/plugin"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
)

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
	arguments := schema.Arguments{}

	// when
	plugin.ExistingPlugins(arguments)

	// then
	assert.Equal(t, 0, len(providers))
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
	}
	kingpinResult := deployCommand.FullCommand()
	provider := dummy.TerralessProvider{}
	provider.Reset()
	providers = []schema.Provider{
		provider,
	}

	_ = os.Mkdir(filepath.Join(os.TempDir(), "terraless-provider-aws-test"), 0755)

	// when
	output := captureOutputProcessCommand(terralessData, kingpinResult)

	// then
	testProcessed := provider.TestProcessed()
	assert.Contains(t, output, "apply -input=false terraform.plan")
	assert.Equal(t, true, testProcessed["PrepareSession"])
	assert.Equal(t, true, testProcessed["ProcessUpload"])
	assert.Equal(t, true, testProcessed["RenderAuthorizerTemplates"])
	assert.Equal(t, true, testProcessed["RenderCertificateTemplates"])
	assert.Equal(t, true, testProcessed["RenderEndpointTemplates"])
	assert.Equal(t, true, testProcessed["RenderFunctionTemplates"])
	assert.Equal(t, true, testProcessed["RenderUploadTemplates"])
}

func TestMain_Init(t *testing.T) {
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
			ProjectName: "DummyProject",
			Providers: map[string]schema.TerralessProvider{
				"DummyProvider1": {
					Name: "DummyProvider1",
					Type: "dummy",
					Roles: []string{
						"role1",
					},
				},
			},
		},
	}
	provider := dummy.TerralessProvider{}
	provider.Reset()
	kingpinResult := initTemplatesCommand.FullCommand()

	_ = os.Mkdir(filepath.Join(os.TempDir(), "terraless-provider-aws-test"), 0755)

	// when
	_ = captureOutputProcessCommand(terralessData, kingpinResult)

	// then
	testProcessed := provider.TestProcessed()
	assert.Equal(t, false, testProcessed["PrepareSession"])
	assert.Equal(t, false, testProcessed["ProcessUpload"])
	assert.Equal(t, false, testProcessed["RenderAuthorizerTemplates"])
	assert.Equal(t, false, testProcessed["RenderCertificateTemplates"])
	assert.Equal(t, false, testProcessed["RenderEndpointTemplates"])
	assert.Equal(t, false, testProcessed["RenderFunctionTemplates"])
	assert.Equal(t, false, testProcessed["RenderUploadTemplates"])
}

func TestMain_Init_RemovesFileIfNoContent(t *testing.T) {
	// given
	tmpFolder := filepath.Join(os.TempDir(), "terraless-provider-aws-test")
	terralessData := schema.TerralessData{
		ActiveProviders: map[string]schema.TerralessProvider{},
		Arguments: schema.Arguments{
			Config:           filepath.Join(tmpFolder, "my-project-config.yml"),
			Environment:      "test",
			ForceDeploy:      true,
			GlobalConfig:     filepath.Join(tmpFolder, "my-global-config.yml"),
			TerraformCommand: "echo",
		},
		Config: schema.TerralessConfig{
			ProjectName: "DummyProject",
			Providers: map[string]schema.TerralessProvider{
				"DummyProvider1": {
					Name: "DummyProvider1",
					Type: "dummy",
					Roles: []string{
						"role1",
					},
				},
			},
			SourcePath: tmpFolder,
		},
	}
	provider := dummy.TerralessProvider{}
	provider.Reset()
	kingpinResult := initTemplatesCommand.FullCommand()

	_ = os.Mkdir(tmpFolder, 0755)
	ressourcesFile := filepath.Join(tmpFolder, "terraless-resources.tf")
	support.WriteToFileIfNotExists(ressourcesFile, "# dummy")

	// when
	_ = captureOutputProcessCommand(terralessData, kingpinResult)

	// then
	testProcessed := provider.TestProcessed()
	assert.Equal(t, false, testProcessed["PrepareSession"])
	assert.Equal(t, false, testProcessed["ProcessUpload"])
	assert.Equal(t, false, testProcessed["RenderAuthorizerTemplates"])
	assert.Equal(t, false, testProcessed["RenderCertificateTemplates"])
	assert.Equal(t, false, testProcessed["RenderEndpointTemplates"])
	assert.Equal(t, false, testProcessed["RenderFunctionTemplates"])
	assert.Equal(t, false, testProcessed["RenderUploadTemplates"])

	info, err := os.Stat(ressourcesFile)
	assert.Equal(t, nil, info)
	assert.NotEqual(t, nil, err)
}

func TestMain_Upload(t *testing.T) {
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
			ProjectName: "DummyProject",
			Providers: map[string]schema.TerralessProvider{
				"DummyProvider1": {
					Name: "DummyProvider1",
					Type: "dummy",
					Roles: []string{
						"role1",
					},
				},
			},
			Uploads: []schema.TerralessUpload{
				{
					Type: "dummyUpload",
				},
			},
		},
	}
	kingpinResult := uploadCommand.FullCommand()
	provider := dummy.TerralessProvider{}
	provider.Reset()

	_ = os.Mkdir(filepath.Join(os.TempDir(), "terraless-provider-aws-test"), 0755)

	// when
	_ = captureOutputProcessCommand(terralessData, kingpinResult)

	// then
	testProcessed := provider.TestProcessed()
	assert.Equal(t, true, testProcessed["PrepareSession"])
	assert.Equal(t, true, testProcessed["ProcessUpload"])
	assert.Equal(t, false, testProcessed["RenderAuthorizerTemplates"])
	assert.Equal(t, false, testProcessed["RenderCertificateTemplates"])
	assert.Equal(t, false, testProcessed["RenderEndpointTemplates"])
	assert.Equal(t, false, testProcessed["RenderFunctionTemplates"])
	assert.Equal(t, true, testProcessed["RenderUploadTemplates"])
}

func TestMain_Auth(t *testing.T) {
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
			ProjectName: "DummyProject",
			Providers: map[string]schema.TerralessProvider{
				"DummyProvider1": {
					Name: "DummyProvider1",
					Type: "dummy",
					Roles: []string{
						"role1",
					},
				},
			},
			Uploads: []schema.TerralessUpload{
				{
					Type: "dummyUpload",
				},
			},
		},
	}
	kingpinResult := authCommand.FullCommand()
	provider := dummy.TerralessProvider{}
	provider.Reset()

	_ = os.Mkdir(filepath.Join(os.TempDir(), "terraless-provider-aws-test"), 0755)

	// when
	_ = captureOutputProcessCommand(terralessData, kingpinResult)

	// then
	testProcessed := provider.TestProcessed()
	assert.Equal(t, true, testProcessed["PrepareSession"])
	assert.Equal(t, false, testProcessed["ProcessUpload"])
	assert.Equal(t, false, testProcessed["RenderAuthorizerTemplates"])
	assert.Equal(t, false, testProcessed["RenderCertificateTemplates"])
	assert.Equal(t, false, testProcessed["RenderEndpointTemplates"])
	assert.Equal(t, false, testProcessed["RenderFunctionTemplates"])
	assert.Equal(t, false, testProcessed["RenderUploadTemplates"])
}
