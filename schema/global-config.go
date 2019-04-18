package schema

import (
	"github.com/sirupsen/logrus"
)

type TerralessGlobalConfig struct {
	Backends []TerralessBackend `yaml:"Backends"`
	Teams    []TerralessTeam    `yaml:"Teams"`
	Uploads  []TerralessUpload  `yaml:"Uploads"`
}

func (globalCfg TerralessGlobalConfig) findTeamByName(teamName string) TerralessTeam {
	if teamName == "" {
		return TerralessTeam{}
	}

	for _, team := range globalCfg.Teams {
		if team.Name == teamName {
			return team
		}
	}

	logrus.Warnf("Team '%s' not found in global config\n", teamName)
	return TerralessTeam{}
}
