package schema

import (
	"fmt"
	"github.com/Odania-IT/terraless/support"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type TerralessConfig struct {
	Authorizers  map[string]TerralessAuthorizer
	Backend      TerralessBackend
	Certificates map[string]TerralessCertificate
	Endpoints    []TerralessEndpoint
	Functions    map[string]TerralessFunction
	Package      TerralessPackage
	ProjectName  string
	Providers    map[string]TerralessProvider
	Settings     TerralessSettings
	SourcePath   string
	TargetPath   string
	Uploads      []TerralessUpload

	// only for rendering template
	HasProvider map[string]bool
	Runtimes    []string
}

func (cfg *TerralessConfig) applyFunctionDefaults() {
	for functionKey, function := range cfg.Functions {
		for eventKey, event := range function.Events {
			if event.Type == "http" && event.Method == "" {
				cfg.Functions[functionKey].Events[eventKey].Method = "ANY"
			}
		}
	}
}

func BuildTerralessConfig(globalCfg TerralessGlobalConfig, projectCfg TerralessProjectConfig, arguments Arguments) TerralessConfig {
	result := TerralessConfig{
		Authorizers:  projectCfg.Authorizers,
		Certificates: map[string]TerralessCertificate{},
		Endpoints:    projectCfg.Endpoints,
		Functions:    projectCfg.Functions,
		HasProvider:  map[string]bool{},
		Package:      projectCfg.Package,
		ProjectName:  ProcessString(projectCfg.ProjectName, arguments, projectCfg.Settings),
		Providers:    map[string]TerralessProvider{},
		Settings:     projectCfg.Settings,
		SourcePath:   projectCfg.SourcePath,
		TargetPath:   projectCfg.TargetPath,
	}

	result.buildCertificates(globalCfg, projectCfg, arguments)
	result.buildBackend(globalCfg, projectCfg, arguments)
	result.buildProviders(globalCfg, projectCfg, arguments)
	result.buildUploads(globalCfg, projectCfg, arguments)
	result.setProviderForBackend(globalCfg)
	result.applyFunctionDefaults()

	for _, provider := range result.Providers {
		logrus.Debugf("Provider: %s Provider-Type: %s Data: %s\n", provider.Name, provider.Type, provider.Data)
		result.HasProvider[provider.Type] = true
	}

	// Write joined config to target dir
	yamlData, err := yaml.Marshal(result)
	if err != nil {
		logrus.Fatalf("Failed serializing yaml! Error: %s\n", err)
	}

	err = ioutil.WriteFile(filepath.Join(projectCfg.TargetPath, "config.yml"), yamlData, 0644)
	if err != nil {
		logrus.Fatalf("Failed writing joined config to target dir! Error: %s\n", err)
	}

	return result
}

func (cfg *TerralessConfig) buildBackend(globalCfg TerralessGlobalConfig, projectCfg TerralessProjectConfig, arguments Arguments) {
	if projectCfg.Backend.Type == "" {
		logrus.Debugln("No Backend defined....")
		return
	}

	logrus.Debug("Building terraless backend")
	cfg.Backend = TerralessBackend{
		Data:       ProcessData(projectCfg.Backend.Data, arguments, projectCfg.Settings),
		Name:       ProcessString(projectCfg.Backend.Name, arguments, projectCfg.Settings),
		Provider:   ProcessString(projectCfg.Backend.Provider, arguments, projectCfg.Settings),
		Type:       ProcessString(projectCfg.Backend.Type, arguments, projectCfg.Settings),
		Workspaces: ProcessData(projectCfg.Backend.Workspaces, arguments, projectCfg.Settings),
	}

	if cfg.Backend.Type == "global" {
		logrus.Debugf("Processing global backend %s\n", cfg.Backend.Name)
		for _, globalBackend := range globalCfg.Backends {
			if globalBackend.Name == cfg.Backend.Name {
				cfg.Backend.Type = globalBackend.Type
				cfg.Backend.Data = ProcessData(EnrichWithData(cfg.Backend.Data, globalBackend.Data), arguments, projectCfg.Settings)
				cfg.Backend.Provider = ProcessString(globalBackend.Provider, arguments, projectCfg.Settings)
				cfg.Backend.Workspaces = ProcessData(EnrichWithData(cfg.Backend.Workspaces, globalBackend.Workspaces), arguments, projectCfg.Settings)

				return
			}
		}

		logrus.Fatalf("Global Backend '%s' not found\n", cfg.Backend.Name)
	}
}

func (cfg *TerralessConfig) buildCertificates(globalCfg TerralessGlobalConfig, projectCfg TerralessProjectConfig, arguments Arguments) {
	for key, certificate := range projectCfg.Certificates {
		certificate.ProjectName = projectCfg.ProjectName
		certificate.TerraformName = "terraless-certificate-" + support.SanitizeString(certificate.Domain)

		cfg.Certificates[key] = certificate
	}
}

func getProviderRole(provider TerralessProvider) string {
	if provider.Data["role"] != "" {
		return provider.Data["role"]
	}

	for _, role := range provider.Roles {
		return role
	}

	return ""
}

func (cfg *TerralessConfig) buildProviders(globalCfg TerralessGlobalConfig, projectCfg TerralessProjectConfig, arguments Arguments) {
	for _, activeProvider := range projectCfg.ActiveProviders {
		for _, provider := range activeProvider.Providers {
			team := globalCfg.findTeamByName(activeProvider.Team)

			newProvider := TerralessProvider{
				Data: ProcessData(EnrichWithData(dataWithoutProfile(team.Data), provider.Data), arguments, projectCfg.Settings),
				Name: ProcessString(provider.Name, arguments, projectCfg.Settings),
				Type: ProcessString(provider.Type, arguments, projectCfg.Settings),
			}

			if newProvider.Type == "global" {
				logrus.Debugf("Processing global provider %s\n", provider.Name)
				globalProvider := findGlobalProvider(activeProvider, provider, globalCfg, arguments, projectCfg.Settings)

				// Make sure the profile name includes the role
				role := getProviderRole(provider)
				if !strings.HasSuffix(newProvider.Name, role) {
					logrus.Debugf("Adding role suffix to provider name: %s [Role: %s]\n", newProvider.Name, role)
					newProvider.Name += "-" + role
				}

				newProvider.Data = ProcessData(EnrichWithData(globalProvider.Data, newProvider.Data), arguments, projectCfg.Settings)
				newProvider.Type = globalProvider.Type
			}

			cfg.Providers[newProvider.Name] = newProvider
		}
	}
}

func (cfg *TerralessConfig) buildUploads(globalCfg TerralessGlobalConfig, projectCfg TerralessProjectConfig, arguments Arguments) {
	for _, upload := range projectCfg.Uploads {
		logrus.Debugf("Processing upload %s\n", upload.Type)
		newUpload := TerralessUpload{
			Bucket:           ProcessString(upload.Bucket, arguments, projectCfg.Settings),
			Certificate:      upload.Certificate,
			Cloudfront:       upload.Cloudfront,
			LambdaAtEdgeFile: ProcessString(upload.LambdaAtEdgeFile, arguments, projectCfg.Settings),
			OwnCertificate:   upload.OwnCertificate,
			Provider:         ProcessString(upload.Provider, arguments, projectCfg.Settings),
			ProjectName:      ProcessString(upload.ProjectName, arguments, projectCfg.Settings),
			Region:           ProcessString(upload.Region, arguments, projectCfg.Settings),
			Source:           ProcessString(upload.Source, arguments, projectCfg.Settings),
			Target:           ProcessString(upload.Target, arguments, projectCfg.Settings),
			Type:             ProcessString(upload.Type, arguments, projectCfg.Settings),
		}

		cfg.Uploads = append(cfg.Uploads, newUpload)
	}
}

func (cfg *TerralessConfig) setProviderForBackend(globalCfg TerralessGlobalConfig) {
	if cfg.Backend.Type == "" || cfg.Backend.Provider == "" {
		return
	}

	if cfg.Providers[cfg.Backend.Provider].Name != "" {
		logrus.Debugln("Provider for Backend already defined....")
		return
	}

	for _, team := range globalCfg.Teams {
		provider := team.findProviderByName(cfg.Backend.Provider)

		if provider.Name != "" {
			parts := strings.Split(cfg.Backend.Provider, "-")

			provider.Data["alias"] = "backend"
			provider.Data["role"] = parts[len(parts)-1]
			cfg.Providers[provider.Name] = provider
			return
		}
	}

	logrus.Fatalf("Could not set provider for Backend '%s' [Provider: %s]\n", cfg.Backend, cfg.Backend.Provider)
}

func findGlobalProvider(activeProvider TerralessActiveProvider, provider TerralessProvider, globalCfg TerralessGlobalConfig, arguments Arguments, settings TerralessSettings) TerralessProvider {
	team := globalCfg.findTeamByName(activeProvider.Team)

	if team.Name == "" {
		logrus.Fatalf("[findGlobalProvider] Team '%s' not found in global config\n", activeProvider.Team)
	}

	providerName := ProcessString(provider.Name, arguments, settings)
	providerByName := team.findProviderByName(providerName)

	if providerByName.Name == "" {
		logrus.Fatalf("[Team: %s] Provider '%s' not found\n", team.Name, providerName)
	}

	return TerralessProvider{
		Data:  EnrichWithData(dataWithoutProfile(team.Data), providerByName.Data),
		Name:  providerByName.Name,
		Type:  providerByName.Type,
		Roles: providerByName.Roles,
	}
}

func (cfg TerralessConfig) Validate() []string {
	logrus.Debug("Verifying config", cfg)
	providerNames := map[string]bool{}
	result := []string{}
	for _, provider := range cfg.Providers {
		if provider.Type == "global" {
			result = append(result, fmt.Sprintf("Unresolved global in provider found! %v\n", provider))
		}

		if providerNames[provider.Name] {
			logrus.Warnf("Provider Name %s is duplicate!! [Provider: %s]", provider.Name, provider)
		}

		providerNames[provider.Name] = true
	}

	if cfg.Backend.Type == "global" {
		result = append(result, fmt.Sprintf("Unresolved global in backend found! %s\n", cfg.Backend.Name))
	}

	for functionName, functionConfig := range cfg.Functions {
		for _, event := range functionConfig.Events {
			if event.Type == "" {
				result = append(result, fmt.Sprintf("Function %s does have event without Type! %s\n", functionName, event))
			}

			if event.Type == "http" && !support.Contains(HttpMethods, event.Method) {
				result = append(result, fmt.Sprintf("Invalid Method in HTTP-Event Function: %s. Method: %s\n", functionName, event.Method))
			}

			if strings.HasPrefix(event.Path, "/") {
				result = append(result, fmt.Sprintf("[ERROR] Path in HTTP-Event starts with '/'. Function: %s. Method: %s\n", functionName, event.Method))
			}
		}
	}

	return result
}
