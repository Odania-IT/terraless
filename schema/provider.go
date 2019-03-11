package schema

type Provider struct {
	CanHandle      CanHandleFunc
	PrepareSession PrepareSessionFunc
	ProcessUpload  ProcessUploadFunc
}

type CanHandleFunc func(resourceType string) bool
type PrepareSessionFunc func(terralessConfig TerralessConfig)
type ProcessUploadFunc func(config TerralessConfig, upload TerralessUpload)
