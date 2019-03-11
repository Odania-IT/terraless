package schema

import (
	"github.com/sirupsen/logrus"
)

// func FindAndEnrichByProviderByName(globalCfg TerralessConfig, teamName string, src TerralessProvider) TerralessProvider {
// 	team := globalCfg.FindTeamByName(teamName)
//
// 	if team.Name == "" {
// 		logrus.Fatalf("Could not find team %s in global config!\n", teamName)
// 	}
//
// 	parts := strings.Split(src.Name, "-")
// 	role := parts[len(parts)-1]
// 	globalName := strings.Join(parts[:len(parts)-1], "-")
//
// 	for _, provider := range team.Providers {
// 		if provider.Name == globalName {
// 			result := provider
// 			result.Name = src.Name
//
// 			provider.Data = EnrichWithData(provider.Data, src.Data)
//
// 			return result
// 		}
// 	}
//
// 	logrus.Fatalf("Could not find provider with name: %s [Global Name: %s, Role: %s] \n", src.Name, globalName, role)
//
// 	// Why is this required??
// 	return TerralessProvider{}
// }
//
// func joinCertificates(dest map[string]TerralessCertificate, override map[string]TerralessCertificate) {
// 	for key, cert := range override {
// 		curCert := dest[key]
//
// 		if curCert.Domain == "" {
// 			dest[key] = cert
// 		} else {
// 			if len(cert.Aliases) > 0 {
// 				curCert.Aliases = append(curCert.Aliases, cert.Aliases...)
// 			}
//
// 			if cert.Domain != "" {
// 				curCert.Domain = cert.Domain
// 			}
//
// 			if cert.Type != "" {
// 				curCert.Type = cert.Type
// 			}
//
// 			if len(cert.Providers) > 0 {
// 				curCert.Providers = append(curCert.Providers, cert.Providers...)
// 			}
//
// 			if cert.ZoneId != "" {
// 				curCert.ZoneId = cert.ZoneId
// 			}
// 		}
// 	}
// }
//
// func joinProviders(dest []TerralessProvider, override []TerralessProvider) {
// 	for key, provider := range override {
// 		curProvider := dest[key]
//
// 		if curProvider.Name == "" {
// 			dest[key] = provider
// 		} else {
// 			if len(provider.Data) > 0 {
// 				curProvider.Data = EnrichWithData(curProvider.Data, provider.Data)
// 			}
//
// 			if provider.Name != "" {
// 				curProvider.Name = provider.Name
// 			}
//
// 			if len(provider.Roles) > 0 {
// 				curProvider.Roles = append(curProvider.Roles, provider.Roles...)
// 			}
// 		}
// 	}
// }
//
// func joinTeams(dest map[string]TerralessTeam, override map[string]TerralessTeam) {
// 	for key, team := range override {
// 		curTeam := dest[key]
//
// 		if curTeam.Name == "" {
// 			dest[key] = team
// 		} else {
// 			if len(team.Data) > 0 {
// 				curTeam.Data = EnrichWithData(curTeam.Data, team.Data)
// 			}
//
// 			if team.Name != "" {
// 				curTeam.Name = team.Name
// 			}
//
// 			if len(team.Providers) > 0 {
// 				joinProviders(curTeam.Providers, team.Providers)
// 			}
// 		}
// 	}
// }

func (team TerralessTeam) findProviderByName(providerName string) TerralessProvider {
	for _, provider := range team.Providers {
		if provider.is(providerName) {
			provider.Data["profile"] = providerName

			newProvider := TerralessProvider{
				Data:  EnrichWithData(dataWithoutProfile(team.Data), provider.Data),
				Name:  providerName,
				Roles: provider.Roles,
				Type:  provider.Type,
			}

			return newProvider
		}
	}

	logrus.Fatalf("[Team: %s] Provider '%s' not found\n", team.Name, providerName)

	return TerralessProvider{}
}

func dataWithoutProfile(data map[string]string) map[string]string {
	result := map[string]string{}

	for key, val := range data {
		if key != "profile" {
			result[key] = val
		}
	}

	return result
}
