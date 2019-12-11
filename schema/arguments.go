package schema

type Arguments struct {
	AuthProvider         string
	Config               string
	Environment          string
	ForceDeploy          bool
	GlobalConfig         string
	LogLevel             string
	NoDeploy             bool
	NoProviderGeneration bool
	NoUpload             bool
	PluginDirectory      string
	TerralessDirectory   string
	TerraformCommand     string
	Variables            map[string]string
}
