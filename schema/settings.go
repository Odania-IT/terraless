package schema

type PermissionData struct {
	Actions   []string `yaml:"Actions"`
	Effect    string   `yaml:"Effect"`
	Resources []string `yaml:"Resources"`
}

type Permissions struct {
	Type string           `yaml:"Type"`
	Data []PermissionData `yaml:"Data"`
}

type TerralessSettings struct {
	AutoSignIn           bool              `yaml:"AutoSignIn"`
	AutoSignInInCloud    bool              `yaml:"AutoSignInInCloud"`
	NoProviderGeneration bool              `yaml:"NoProviderGeneration"`
	Permissions          []Permissions     `yaml:"Permissions"`
	Runtime              string            `yaml:"Runtime"`
	TerraformPluginDir   string            `yaml:"TerraformPluginDir"`
	Variables            map[string]string `yaml:"Variables"`
}
