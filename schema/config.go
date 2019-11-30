package schema

import (
	"strings"
)

type TerralessBackend struct {
	Data       map[string]string `yaml:"Data"`
	Name       string            `yaml:"Name"`
	Provider   string            `yaml:"Provider"`
	Type       string            `yaml:"Type"`
	Workspaces map[string]string `yaml:"Workspaces"`
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
	Aliases                   []string                   `yaml:"Aliases"`
	Caching                   TerralessCloudfrontCaching `yaml:"Caching"`
	Certificate               string                     `yaml:"Certificate"`
	Handler                   string                     `yaml:"Handler"`
	LambdaFunctionAssociation map[string]string          `yaml:"LambdaFunctionAssociation"`
	NoCreateBucket            bool                       `yaml:"NoCreateBucket"`
	Domain                    string                     `yaml:"Domain"`
	LoggingEnabled            bool                       `yaml:"LoggingEnabled"`
	PriceClass                string                     `yaml:"PriceClass"`
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

var HttpMethods = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "ANY"}

type TerralessFunctionEvent struct {
	Authorizer string            `yaml:"Authorizer"`
	Arn        string            `yaml:"Arn"`
	Cors       TerralessCors     `yaml:"Cors"`
	Event      map[string]string `yaml:"Event"`
	Method     string            `yaml:"Method"`
	Path       string            `yaml:"Path"`
	Rate       string            `yaml:"Rate"`
	Route      string            `yaml:"Route"`
	Type       string            `yaml:"Type"`

	// only for rendering template
	Authorization       string
	AuthorizerId        string
	FunctionName        string
	Idx                 string
	ProjectName         string
	PathsRendered       map[string]string
	ResourceNameForPath string
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
	Type        string                   `yaml:"Type"`

	// only for rendering template
	AddApiGatewayPermission bool
	FunctionName            string
	ProjectName             string
	RenderEnvironment       bool
}

type TerralessPackage struct {
	BuildCommand      string `yaml:"BuildCommand"`
	LambdaArchiveFile string `yaml:"LambdaArchiveFile"`
	SourceDir         string `yaml:"SourceDir"`
	OutputPath        string
}

type TerralessProvider struct {
	Data  map[string]string `yaml:"Data"`
	Name  string            `yaml:"Name"`
	Roles []string          `yaml:"Roles"`
	Type  string            `yaml:"Type"`
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
	Environment      string
	LambdaAtEdgeFile string
	OwnCertificate   bool
	ProjectName      string
}

type TerralessTeam struct {
	Data      map[string]string   `yaml:"Data"`
	Name      string              `yaml:"Name"`
	Providers []TerralessProvider `yaml:"Providers"`
}

func getRoleAndNameFromProvider(name string) (string, string) {
	parts := strings.Split(name, "-")
	role := parts[len(parts)-1]
	globalName := strings.Join(parts[:len(parts)-1], "-")

	return globalName, role
}

func (profile TerralessProvider) is(name string) bool {
	if profile.Name == name {
		return true
	}

	globalName, role := getRoleAndNameFromProvider(name)
	if profile.Name == globalName {
		for _, profileRole := range profile.Roles {
			if profileRole == role {
				return true
			}
		}
	}

	return false
}
