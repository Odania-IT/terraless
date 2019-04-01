package schema

type Arguments struct {
	Config               string
	Environment          string
	ForceDeploy          bool
	GlobalConfig         string
	LogLevel             string
	NoDeploy             bool
	NoProviderGeneration bool
	NoUpload             bool
	TerraformCommand     string
}
