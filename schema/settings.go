package schema

type TerralessSettings struct {
	AutoSignIn           bool              `yaml:"AutoSignIn"`
	AutoSignInInCloud    bool              `yaml:"AutoSignInInCloud"`
	NoProviderGeneration bool              `yaml:"NoProviderGeneration"`
	Runtime              string            `yaml:"Runtime"`
	TerraformPluginDir   string            `yaml:"TerraformPluginDir"`
	Variables            map[string]string `yaml:"Variables"`
}
