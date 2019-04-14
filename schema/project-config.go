package schema

type TerralessActiveProvider struct {
	Team      string              `yaml:"Team"`
	Providers []TerralessProvider `yaml:"Providers"`
}

type TerralessAuthorizer struct {
	AuthorizerType string   `yaml:"AuthorizerType"`
	Type           string   `yaml:"Type"`
	Name           string   `yaml:"Name"`
	ProviderArns   []string `yaml:"ProviderArns"`

	// only for rendering template
	TerraformName string
}

type TerralessEndpoint struct {
	Aliases     []string `yaml:"Aliases"`
	BasePath    string   `yaml:"BasePath"`
	Certificate string   `yaml:"Certificate"`
	Domain      string   `yaml:"Domain"`
	Type        string   `yaml:"Type"`

	// only for rendering template
	TerralessCertificate TerralessCertificate
	TerraformName        string
}

type TerralessProjectConfig struct {
	Authorizers     map[string]TerralessAuthorizer  `yaml:"Authorizers"`
	ActiveProviders []TerralessActiveProvider       `yaml:"ActiveProviders"`
	Backend         TerralessBackend                `yaml:"Backend"`
	Certificates    map[string]TerralessCertificate `yaml:"Certificates"`
	Endpoints       []TerralessEndpoint             `yaml:"Endpoints"`
	Functions       map[string]TerralessFunction    `yaml:"Functions"`
	Package         TerralessPackage                `yaml:"Package"`
	ProjectName     string                          `yaml:"ProjectName"`
	Settings        TerralessSettings               `yaml:"Settings"`
	SourcePath      string                          `yaml:"SourcePath"`
	TargetPath      string                          `yaml:"TargetPath"`
	Uploads         []TerralessUpload               `yaml:"Uploads"`

	// only for rendering template
	HasProvider map[string]bool
	Runtimes    []string
	TeamData    map[string]map[string]string
}

func (cfg *TerralessProjectConfig) Validate() *TerralessProjectConfig {
	if cfg.Settings.Runtime == "" {
		cfg.Settings.Runtime = "ruby2.5"
	}

	return cfg
}
