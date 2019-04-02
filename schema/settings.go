package schema

type TerralessSettings struct {
	AutoSignIn           bool   `yaml:"AutoSignIn"`
	NoProviderGeneration bool   `yaml:"NoProviderGeneration"`
	Runtime              string `yaml:"Runtime"`
}
