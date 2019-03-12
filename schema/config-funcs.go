package schema

import (
	"github.com/sirupsen/logrus"
)

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
		if key == "profile" {
			key = "base-profile"
		}
		result[key] = val
	}

	return result
}
