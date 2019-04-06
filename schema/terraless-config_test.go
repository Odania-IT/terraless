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
		Certificates: map[string]TerralessCertificate{},
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
		Certificates: map[string]TerralessCertificate{},
		HasProvider: map[string]bool{},
		Providers: map[string]TerralessProvider{},
		Settings: TerralessSettings{
			AutoSignIn: true,
		},
	}
	assert.Equal(t, expected, config)
}

func TestTerralessConfig_BuildTerralessConfig_Uploads(t *testing.T) {
	// given
	arguments := Arguments{}
	globalConfig := TerralessGlobalConfig{}
	projectConfig := TerralessProjectConfig{
		Settings: TerralessSettings{
			AutoSignIn: true,
		},
		Uploads: []TerralessUpload{
			{
				Type: "s3",
				Bucket: "dummyBucket",
				Source: "my-src",
				Target: "my-target",
			},
		},
	}

	// when
	config := BuildTerralessConfig(globalConfig, projectConfig, arguments)

	// then
	expected := TerralessConfig{
		Certificates: map[string]TerralessCertificate{},
		HasProvider: map[string]bool{},
		Providers: map[string]TerralessProvider{},
		Settings: TerralessSettings{
			AutoSignIn: true,
		},
		Uploads: []TerralessUpload{
			{
				Type: "s3",
				Bucket: "dummyBucket",
				Source: "my-src",
				Target: "my-target",
			},
		},
	}
	assert.Equal(t, expected, config)
}

func TestTerralessConfig_Validate_Nothing(t *testing.T) {
	// given
	cfg := TerralessConfig{}

	// when
	validate := cfg.Validate()

	// then
	assert.Equal(t, 0, len(validate))
}

func TestTerralessConfig_Validate_Errors(t *testing.T) {
	// given
	cfg := TerralessConfig{
		Backend: TerralessBackend{
			Type: "global",
			Name: "DummyBackend",
		},
		Functions: map[string]TerralessFunction{
			"func1": {
				Type: "aws",
				Events: []TerralessFunctionEvent{
					{
						Type: "http",
						Path: "/ohoh",
						Method: "GET",
					},
				},
			},
		},
		Providers: map[string]TerralessProvider{
			"provider1": {
				Type: "global",
				Name: "provider1",
			},
			"provider2": {
				Type: "dummy",
				Name: "provider1",
			},
		},
	}

	// when
	validate := cfg.Validate()

	// then
	expected := []string{
		"Unresolved global in provider found! {map[] provider1 [] global}\n",
		"Unresolved global in backend found! DummyBackend\n",
		"[ERROR] Path in HTTP-Event starts with '/'. Function: func1. Method: GET\n",
	}

	assert.Equal(t, expected, validate)
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
		Certificates: map[string]TerralessCertificate{},
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
		Certificates: map[string]TerralessCertificate{},
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
		Certificates: map[string]TerralessCertificate{},
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
		Certificates: map[string]TerralessCertificate{},
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
						Data: map[string]string{},
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
		Certificates: map[string]TerralessCertificate{},
		HasProvider: map[string]bool{
			"dummy": true,
		},
		Providers: map[string]TerralessProvider{
			"dummy-provider": {
				Name: "dummy-provider",
				Type: "dummy",
				Data: map[string]string{
					"alias": "backend",
					"role": "provider",
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
		Certificates: map[string]TerralessCertificate{},
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
						Data: map[string]string{},
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
		Certificates: map[string]TerralessCertificate{},
		HasProvider: map[string]bool{
			"dummy": true,
		},
		Providers: map[string]TerralessProvider{
			"dummy-provider": {
				Type: "dummy",
				Name: "dummy-provider",
				Data: map[string]string{
					"alias": "backend",
					"role": "provider",
					"teamKey": "teamValue",
				},
			},
		},
	}
	assert.Equal(t, expected, config)
}
