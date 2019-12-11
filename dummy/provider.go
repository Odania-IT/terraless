package dummy

import (
	"github.com/Odania-IT/terraless/schema"
)

var testProcessed map[string]bool

type TerralessProvider struct {}

func (provider TerralessProvider) Info() schema.PluginInfo {
	return schema.PluginInfo{
		Name:    "dummy",
		Version: "0.4.42",
	}
}

func (provider TerralessProvider) CanHandle(resourceType string) bool {
	return resourceType == "dummy"
}

func (provider TerralessProvider) PrepareSession(terralessConfig schema.TerralessConfig) map[string]string {
	testProcessed["PrepareSession"] = true

	return map[string]string{}
}

func (provider TerralessProvider) ProcessUpload(terralessData schema.TerralessData, upload schema.TerralessUpload) []string {
	testProcessed["ProcessUpload"] = true

	return []string{}
}

func (provider TerralessProvider) FinalizeTemplates(terralessData schema.TerralessData) string {
	testProcessed["FinalizeTemplates"] = true

	return ""
}

func (provider TerralessProvider) GenerateHelperFunctionCommand(teamName string, providerName string, roleName string) string {
	testProcessed["GenerateHelperFunctionCommand"] = true

	return "generated-command-" + teamName + "-" + providerName + "-" + roleName
}

func (provider TerralessProvider) RenderAuthorizerTemplates(config schema.TerralessConfig) string {
	testProcessed["RenderAuthorizerTemplates"] = true

	return ""
}

func (provider TerralessProvider) RenderCertificateTemplates(config schema.TerralessConfig) string {
	testProcessed["RenderCertificateTemplates"] = true

	return ""
}

func (provider TerralessProvider) RenderEndpointTemplates(config schema.TerralessConfig) string {
	testProcessed["RenderEndpointTemplates"] = true

	return ""
}

func (provider TerralessProvider) RenderFunctionTemplates(resourceType string, functionEvents schema.FunctionEvents, terralessData *schema.TerralessData) string {
	testProcessed["RenderFunctionTemplates"] = true

	return ""
}

func (provider TerralessProvider) RenderUploadTemplates(terralessData schema.TerralessData) string {
	testProcessed["RenderUploadTemplates"] = true

	return ""
}

func (provider TerralessProvider) Reset() {
	testProcessed = map[string]bool{}
}

func (provider TerralessProvider) TestProcessed() map[string]bool {
	return testProcessed
}
