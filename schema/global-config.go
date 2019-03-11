package schema

import (
	"github.com/sirupsen/logrus"
)

type TerralessGlobalConfig struct {
	Backends []TerralessBackend `yaml:"Backends"`
	Teams    []TerralessTeam    `yaml:"Teams"`
	Uploads  []TerralessUpload  `yaml:"Uploads"`
}

// func (globalCfg *TerralessGlobalConfig) BuildFinalConfig(projectConfig *TerralessProjectConfig, arguments Arguments) {
// 	if projectConfig.Backend.Type == "global" {
// 		globalBackend := globalCfg.findBackendByName(projectConfig.Backend.Name)
//
// 		if globalBackend.Type == "" {
// 			fmt.Println("Available Backends:")
// 			for _, backend := range globalCfg.Backends {
// 				fmt.Printf("* %s\n", backend.Name)
// 			}
//
// 			logrus.Fatalf("Global Backend '%s' not defined!", projectConfig.Backend.Name)
// 		}
//
// 		projectConfig.Backend.Name = globalBackend.Name
// 		projectConfig.Backend.Type = globalBackend.Type
// 		projectConfig.Backend.Data = EnrichWithData(projectConfig.Backend.Data, globalBackend.Data)
// 	}
//
// 	for activeProviderKey, activeProvider := range projectConfig.ActiveProviders {
// 		teamName := ProcessString(activeProvider.Team, arguments)
// 		team := globalCfg.findTeamByName(teamName)
//
// 		if team.Name == "" {
// 			logrus.Fatalf("Team '%s' not found in global config\n", teamName)
// 		}
//
// 		projectConfig.TeamData[teamName] = team.Data
//
// 		for providerKey, provider := range activeProvider.Providers {
// 			if provider.Type == "global" {
// 				teamProvider := team.findProviderByName(provider.Name)
// 				provider.Name = ProcessString(teamProvider.Name, arguments)
// 				provider.Type = ProcessString(teamProvider.Type, arguments)
// 				provider.Data = ProcessData(teamProvider.Data, arguments)
// 			}
//
// 			projectConfig.ActiveProviders[activeProviderKey].Providers[providerKey] = provider
// 		}
// 	}
// }

func (globalCfg TerralessGlobalConfig) findTeamByName(teamName string) TerralessTeam {
	for _, team := range globalCfg.Teams {
		if team.Name == teamName {
			return team
		}
	}

	logrus.Fatalf("Team '%s' not found in global config\n", teamName)
	return TerralessTeam{}
}

func (globalCfg TerralessGlobalConfig) findBackendByName(backendName string) TerralessBackend {
	for _, backend := range globalCfg.Backends {
		if backend.Name == backendName {
			return backend
		}
	}

	logrus.Fatalf("Could not find backend '%s'\n", backendName)
	return TerralessBackend{}
}
