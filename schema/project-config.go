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

// func (cfg TerralessProjectConfig) ProcessConfig(arguments Arguments) *TerralessProjectConfig {
// 	processedConfig := &TerralessProjectConfig{
// 		ActiveProviders: []TerralessActiveProvider{},
// 		Backend:         cfg.Backend,
// 		Certificates:    cfg.Certificates,
// 		Functions:       cfg.Functions,
// 		HasProvider:     map[string]bool{},
// 		Package:         cfg.Package,
// 		ProjectName:     ProcessString(cfg.ProjectName, arguments),
// 		Settings:        cfg.Settings,
// 		TeamData:        map[string]map[string]string{},
// 	}
//
// 	processedConfig.Backend.Name = ProcessString(processedConfig.Backend.Name, arguments)
// 	processedConfig.Backend.Data = ProcessData(processedConfig.Backend.Data, arguments)
//
// 	for _, activeProvider := range cfg.ActiveProviders {
// 		newActiveProvider := TerralessActiveProvider{
// 			Team: activeProvider.Team,
// 		}
//
// 		for _, provider := range activeProvider.Providers {
// 			provider.Name = ProcessString(provider.Name, arguments)
// 			provider.Data = ProcessData(provider.Data, arguments)
// 			newActiveProvider.Providers = append(newActiveProvider.Providers, provider)
// 		}
//
// 		processedConfig.ActiveProviders = append(processedConfig.ActiveProviders, newActiveProvider)
// 	}
//
// 	for _, upload := range cfg.Uploads {
// 		upload.Bucket = ProcessString(upload.Bucket, arguments)
// 		upload.Provider = ProcessString(upload.Provider, arguments)
// 		processedConfig.Uploads = append(processedConfig.Uploads, upload)
// 	}
//
// 	return processedConfig
// }

// func (cfg TerralessProjectConfig) FindTeamByName(teamName string) TerralessTeam {
// 	for _, team := range cfg.Teams {
// 		if team.Name == teamName {
// 			return team
// 		}
// 	}
//
// 	return TerralessTeam{}
// }

// func (cfg TerralessProjectConfig) FindProviderByName(name string) (TerralessProvider, TerralessActiveProvider) {
// 	for _, activeProvider := range cfg.ActiveProviders {
// 		for _, provider := range activeProvider.Providers {
// 			if provider.Name == name {
// 				return provider, activeProvider
// 			}
// 		}
// 	}
//
// 	logrus.Fatalf("Could not find provider %s\n", name)
//
// 	// Why is this required??
// 	return TerralessProvider{}, TerralessActiveProvider{}
// }
