package schema

import "strings"

type TerralessAuthorizer struct {
	Type string `yaml:"Type"`
}

type TerralessBackend struct {
	Data     map[string]string `yaml:"Data"`
	Name     string            `yaml:"Name"`
	Provider string            `yaml:"Provider"`
	Type     string            `yaml:"Type"`
}

type TerralessCertificate struct {
	Aliases   []string `yaml:"Aliases"`
	Domain    string   `yaml:"Domain"`
	Type      string   `yaml:"Type"`
	Providers []string `yaml:"Providers"`
	ZoneId    string   `yaml:"ZoneId"`

	// only for rendering template
	ProjectName   string
	TerraformName string
}

type TerralessCloudfront struct {
	Aliases        []string                   `yaml:"Aliases"`
	Caching        TerralessCloudfrontCaching `yaml:"Caching"`
	Certificate    string                     `yaml:"Certificate"`
	Domain         string                     `yaml:"Domain"`
	LoggingEnabled bool                       `yaml:"LoggingEnabled"`
	PriceClass     string                     `yaml:"PriceClass"`
}

type TerralessCloudfrontCaching struct {
	DefaultTTL int `yaml:"DefaultTTL"`
	MinTTL     int `yaml:"MinTTL"`
	MaxTTL     int `yaml:"MaxTTL"`
}

type TerralessCors struct {
	Headers []string `yaml:"Headers"`
	Origin  string   `yaml:"Origin"`
}

type TerralessFunctionEvent struct {
	Authorizer TerralessAuthorizer `yaml:"Authorizer"`
	Arn        string              `yaml:"Arn"`
	Cors       TerralessCors       `yaml:"Cors"`
	Method     string              `yaml:"Method"`
	Path       string              `yaml:"Path"`
	Rate       string              `yaml:"Rate"`
	Route      string              `yaml:"Route"`
	Type       string              `yaml:"Type"`

	// only for rendering template
	FunctionName string
}

type TerralessFunction struct {
	Description string                   `yaml:"Description"`
	Environment map[string]string        `yaml:"Environment"`
	Events      []TerralessFunctionEvent `yaml:"Events"`
	Handler     string                   `yaml:"Handler"`
	MemorySize  int                      `yaml:"MemorySize"`
	RoleArn     string                   `yaml:"RoleArn"`
	Runtime     string                   `yaml:"Runtime"`
	Timeout     int                      `yaml:"Timeout"`

	// only for rendering template
	FunctionName      string
	RenderEnvironment bool
}

type TerralessPackage struct {
	SourceDir  string `yaml:"SourceDir"`
	OutputPath string
}

type TerralessProvider struct {
	Data  map[string]string `yaml:"Data"`
	Name  string            `yaml:"Name"`
	Roles []string          `yaml:"Roles"`
	Type  string            `yaml:"Type"`
}

type TerralessSettings struct {
	AutoSignIn bool `yaml:"AutoSignIn"`
}

type TerralessType struct {
	Data map[string]string `yaml:"Data"`
	Name string            `yaml:"Name"`
	Type string            `yaml:"Type"`
}

type TerralessUpload struct {
	Bucket     string              `yaml:"Bucket"`
	Cloudfront TerralessCloudfront `yaml:"Cloudfront"`
	Provider   string              `yaml:"Provider"`
	Region     string              `yaml:"Region"`
	Source     string              `yaml:"Source"`
	Target     string              `yaml:"Target"`
	Type       string              `yaml:"Type"`

	// only for rendering template
	Certificate      TerralessCertificate
	LambdaAtEdgeFile string
	OwnCertificate   bool
	ProjectName      string
}

type TerralessTeam struct {
	Data      map[string]string   `yaml:"Data"`
	Name      string              `yaml:"Name"`
	Providers []TerralessProvider `yaml:"Providers"`
}

func (profile TerralessProvider) is(name string) bool {
	if profile.Name == name {
		return true
	}

	parts := strings.Split(name, "-")
	role := parts[len(parts)-1]
	globalName := strings.Join(parts[:len(parts)-1], "-")

	if profile.Name == globalName {
		for _, profileRole := range profile.Roles {
			if profileRole == role {
				return true
			}
		}
	}

	return false
}
