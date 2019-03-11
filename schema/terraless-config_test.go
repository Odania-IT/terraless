package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerralessConfig_BuildTerralessConfig_NoData(t *testing.T) {
	// given
	arguments := Arguments{}
	globalConfig := TerralessGlobalConfig{}
	projectConfig := TerralessProjectConfig{}

	// when
	config := BuildTerralessConfig(globalConfig, projectConfig, arguments)

	// then
	expected := TerralessConfig{
		HasProvider: map[string]bool{},
		Providers: map[string]TerralessProvider{},
	}
	assert.Equal(t, expected, config)
}

func TestTerralessConfig_BuildTerralessConfig_SettingsIsInConfig(t *testing.T) {
	// given
	arguments := Arguments{}
	globalConfig := TerralessGlobalConfig{}
	projectConfig := TerralessProjectConfig{
		Settings: TerralessSettings{
			AutoSignIn: true,
		},
	}

	// when
	config := BuildTerralessConfig(globalConfig, projectConfig, arguments)

	// then
	expected := TerralessConfig{
		HasProvider: map[string]bool{},
		Providers: map[string]TerralessProvider{},
		Settings: TerralessSettings{
			AutoSignIn: true,
		},
	}
	assert.Equal(t, expected, config)
}

func TestTerralessConfig_BuildTerralessConfig_ProvidersIsInGlobalConfig(t *testing.T) {
	// given
	arguments := Arguments{}
	globalConfig := TerralessGlobalConfig{
		Teams: []TerralessTeam{
			{
				Providers: []TerralessProvider{
					{
						Type: "aws",
					},
				},
			},
		},
	}
	projectConfig := TerralessProjectConfig{}

	// when
	config := BuildTerralessConfig(globalConfig, projectConfig, arguments)

	// then
	expected := TerralessConfig{
		HasProvider: map[string]bool{},
		Providers: map[string]TerralessProvider{},
	}
	assert.Equal(t, expected, config)
}

func TestTerralessConfig_BuildTerralessConfig_ProvidersIsInProjectConfig(t *testing.T) {
	// given
	arguments := Arguments{}
	globalConfig := TerralessGlobalConfig{}
	projectConfig := TerralessProjectConfig{
		ActiveProviders: []TerralessActiveProvider{
			{
				Team: "myTeam",
				Providers: []TerralessProvider{
					{
						Type: "aws",
						Data: map[string]string{
							"dummyData": "dummyValue",
						},
						Name: "myTeam-myEnvironment-myRole",
						Roles: []string{
							"myRole",
						},
					},
					{
						Type: "aws",
						Data: map[string]string{
							"dummySecondData": "example",
						},
						Name: "myTeam-myEnvironment-mySecondRole",
						Roles: []string{
							"mySecondRole",
						},
					},
				},
			},
		},
	}

	// when
	config := BuildTerralessConfig(globalConfig, projectConfig, arguments)

	// then
	assert.Equal(t, 2, len(config.Providers))

	expected := TerralessConfig{
		HasProvider: map[string]bool{
			"aws": true,
		},
		Providers: map[string]TerralessProvider{
			"myTeam-myEnvironment-myRole": {
				Name: "myTeam-myEnvironment-myRole",
				Type: "aws",
				Data: map[string]string{
					"dummyData": "dummyValue",
				},
			},
			"myTeam-myEnvironment-mySecondRole": {
				Name: "myTeam-myEnvironment-mySecondRole",
				Type: "aws",
				Data: map[string]string{
					"dummySecondData": "example",
				},
			},
		},
	}
	assert.Equal(t, expected, config)
}

func TestTerralessConfig_BuildTerralessConfig_ProjectConfigWithGlobalProvider(t *testing.T) {
	// given
	arguments := Arguments{}
	globalConfig := TerralessGlobalConfig{
		Teams: []TerralessTeam{
			{
				Name: "myTeam",
				Providers: []TerralessProvider{
					{
						Type: "aws",
						Name: "myTeam-myEnvironment",
						Data: map[string]string{
							"dummyData": "dummyValue",
						},
						Roles: []string{
							"myRole",
						},
					},
				},
			},
		},
	}
	projectConfig := TerralessProjectConfig{
		ActiveProviders: []TerralessActiveProvider{
			{
				Team: "myTeam",
				Providers: []TerralessProvider{
					{
						Type: "global",
						Name: "myTeam-myEnvironment-myRole",
						Roles: []string{
							"myRole",
						},
					},
					{
						Type: "aws",
						Data: map[string]string{
							"dummySecondData": "example",
						},
						Name: "myTeam-myEnvironment-mySecondRole",
						Roles: []string{
							"mySecondRole",
						},
					},
				},
			},
		},
	}

	// when
	config := BuildTerralessConfig(globalConfig, projectConfig, arguments)

	// then
	assert.Equal(t, 2, len(config.Providers))

	expected := TerralessConfig{
		HasProvider: map[string]bool{
			"aws": true,
		},
		Providers: map[string]TerralessProvider{
			"myTeam-myEnvironment-myRole": {
				Name: "myTeam-myEnvironment-myRole",
				Type: "aws",
				Data: map[string]string{
					"dummyData": "dummyValue",
				},
			},
			"myTeam-myEnvironment-mySecondRole": {
				Name: "myTeam-myEnvironment-mySecondRole",
				Type: "aws",
				Data: map[string]string{
					"dummySecondData": "example",
				},
			},
		},
	}
	assert.Equal(t, expected, config)
}

func TestTerralessConfig_BuildTerralessConfig_Backend(t *testing.T) {
	// given
	arguments := Arguments{}
	globalConfig := TerralessGlobalConfig{}
	projectConfig := TerralessProjectConfig{
		Backend: TerralessBackend{
			Type: "dummy",
			Name: "dummy-backend",
			Data: map[string]string{
				"key": "val",
			},
		},
	}

	// when
	config := BuildTerralessConfig(globalConfig, projectConfig, arguments)

	// then
	expected := TerralessConfig{
		Backend: TerralessBackend{
			Type: "dummy",
			Name: "dummy-backend",
			Data: map[string]string{
				"key": "val",
			},
		},
		HasProvider: map[string]bool{},
		Providers: map[string]TerralessProvider{},
	}
	assert.Equal(t, expected, config)
}

func TestTerralessConfig_BuildTerralessConfig_BackendWithProvider(t *testing.T) {
	// given
	arguments := Arguments{}
	globalConfig := TerralessGlobalConfig{
		Teams: []TerralessTeam{
			{
				Name: "myTeam",
				Data: map[string]string{
					"teamData": "teamValue",
				},
				Providers: []TerralessProvider{
					{
						Name: "dummy-provider",
						Type: "dummy",
					},
				},
			},
		},
	}
	projectConfig := TerralessProjectConfig{
		Backend: TerralessBackend{
			Type: "dummy",
			Name: "dummy-backend",
			Provider: "dummy-provider",
			Data: map[string]string{
				"key": "val",
			},
		},
	}

	// when
	config := BuildTerralessConfig(globalConfig, projectConfig, arguments)

	// then
	expected := TerralessConfig{
		Backend: TerralessBackend{
			Type: "dummy",
			Name: "dummy-backend",
			Provider: "dummy-provider",
			Data: map[string]string{
				"key": "val",
			},
		},
		HasProvider: map[string]bool{
			"dummy": true,
		},
		Providers: map[string]TerralessProvider{
			"dummy-provider": {
				Name: "dummy-provider",
				Type: "dummy",
				Data: map[string]string{
					"teamData": "teamValue",
				},
			},
		},
	}
	assert.Equal(t, expected, config)
}

func TestTerralessConfig_BuildTerralessConfig_GlobalBackend(t *testing.T) {
	// given
	arguments := Arguments{}
	globalConfig := TerralessGlobalConfig{
		Backends: []TerralessBackend{
			{
				Type: "dummy",
				Name: "dummy-backend",
				Data: map[string]string{
					"secondKey": "val",
				},
			},
		},
	}
	projectConfig := TerralessProjectConfig{
		Backend: TerralessBackend{
			Type: "global",
			Name: "dummy-backend",
			Data: map[string]string{
				"key": "val",
			},
		},
	}

	// when
	config := BuildTerralessConfig(globalConfig, projectConfig, arguments)

	// then
	expected := TerralessConfig{
		Backend: TerralessBackend{
			Type: "dummy",
			Name: "dummy-backend",
			Data: map[string]string{
				"key": "val",
				"secondKey": "val",
			},
		},
		HasProvider: map[string]bool{},
		Providers: map[string]TerralessProvider{},
	}
	assert.Equal(t, expected, config)
}

func TestTerralessConfig_BuildTerralessConfig_GlobalBackendWithProvider(t *testing.T) {
	// given
	arguments := Arguments{}
	globalConfig := TerralessGlobalConfig{
		Backends: []TerralessBackend{
			{
				Type: "dummy",
				Name: "dummy-backend",
				Provider: "dummy-provider",
				Data: map[string]string{
					"secondKey": "val",
				},
			},
		},
		Teams: []TerralessTeam{
			{
				Name: "myTeam",
				Data: map[string]string{
					"teamKey": "teamValue",
				},
				Providers: []TerralessProvider{
					{
						Type: "dummy",
						Name: "dummy-provider",
					},
				},
			},
		},
	}
	projectConfig := TerralessProjectConfig{
		Backend: TerralessBackend{
			Type: "global",
			Name: "dummy-backend",
			Data: map[string]string{
				"key": "val",
			},
		},
	}

	// when
	config := BuildTerralessConfig(globalConfig, projectConfig, arguments)

	// then
	expected := TerralessConfig{
		Backend: TerralessBackend{
			Type: "dummy",
			Name: "dummy-backend",
			Provider: "dummy-provider",
			Data: map[string]string{
				"key": "val",
				"secondKey": "val",
			},
		},
		HasProvider: map[string]bool{
			"dummy": true,
		},
		Providers: map[string]TerralessProvider{
			"dummy-provider": {
				Type: "dummy",
				Name: "dummy-provider",
				Data: map[string]string{
					"teamKey": "teamValue",
				},
			},
		},
	}
	assert.Equal(t, expected, config)
}
