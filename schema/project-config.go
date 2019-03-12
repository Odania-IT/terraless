package schema

type TerralessActiveProvider struct {
	Team      string              `yaml:"Team"`
	Providers []TerralessProvider `yaml:"Providers"`
}

type TerralessProjectConfig struct {
	ActiveProviders []TerralessActiveProvider       `yaml:"ActiveProviders"`
	Backend         TerralessBackend                `yaml:"Backend"`
	Certificates    map[string]TerralessCertificate `yaml:"Certificates"`
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
