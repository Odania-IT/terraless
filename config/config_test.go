package config

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func baseDir() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "..")
}

func TestTerralessConfig_NewTerralessConfig(t *testing.T) {
	// given
	args := schema.Arguments{
		Config:       filepath.Join(baseDir(), "examples", "terraless-project.yml"),
		Environment:  "develop",
		GlobalConfig: filepath.Join(baseDir(), "examples", "terraless.yml"),
	}

	// when
	terralessData := NewTerralessData(args)

	// then
	expected := schema.TerralessData{
		Arguments: args,
		Config: schema.TerralessConfig{
			Backend: schema.TerralessBackend{
				Name: "myBackend",
				Type: "s3",
				Data: map[string]string{
					"bucket":         "my-bucket-name",
					"encrypt":        "true",
					"region":         "eu-central-1",
					"key":            "myProjectKey",
					"dynamodb_table": "terraform-state-lock",
					"profile":        "my-aws-infrastrucutre-profile-developer",
				},
				Workspaces: map[string]string{},
			},
			Certificates: map[string]schema.TerralessCertificate{
				"MyCert": {
					Aliases: []string{
						"*.my-domain.com",
					},
					Domain: "my-domain.com",
					Type:   "aws",
					Providers: []string{
						"aws.us-east",
					},
					ZoneId:        "${aws_route53_zone.zone.id}",
					ProjectName:   "examples",
					TerraformName: "terraless-certificate-my-domain-com",
				},
			},
			Functions: map[string]schema.TerralessFunction{
				"MyTestLambda": {
					Description: "My Test Lambda Description",
					Events: []schema.TerralessFunctionEvent{
						{
							Type:   "http",
							Method: "ANY",
						},
					},
					Handler:    "test.Handler",
					MemorySize: 512,
					Runtime:    "ruby2.5",
					Timeout:    60,
					Type:       "aws",
				},
				"MyTestLambda2": {
					Description: "My Test Lambda Description",
					Events: []schema.TerralessFunctionEvent{
						{
							Type: "sqs",
							Arn:  "arn:aws:sqs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${self:provider.stage}-my-queue",
						},
					},
					Handler:    "test.Handler2",
					MemorySize: 512,
					Runtime:    "ruby2.5",
					Timeout:    60,
					Type:       "aws",
				},
			},
			Package: schema.TerralessPackage{
				SourceDir: "src",
			},
			ProjectName: "examples",
			Providers: map[string]schema.TerralessProvider{
				"aws-default": {
					Type: "aws",
					Name: "aws-default",
					Data: map[string]string{
						"accountId": "01234556678",
						"profile":   "my-aws-profile",
						"region":    "eu-central-1",
					},
				},
				"aws-develop-developer": {
					Type: "aws",
					Name: "aws-develop-developer",
					Data: map[string]string{
						"accountId": "01234556678",
						"profile":   "my-aws-profile",
						"region":    "eu-central-1",
					},
				},
				"custom-aws": {
					Type: "aws",
					Name: "custom-aws",
					Data: map[string]string{
						"alias":   "eu-west-1",
						"profile": "my-custom-aws-profile",
						"region":  "eu-west-1",
					},
				},
			},
			Settings: schema.TerralessSettings{
				AutoSignIn: true,
				Runtime:    "ruby2.5",
			},
			SourcePath: filepath.Join(baseDir(), "examples"),
			TargetPath: filepath.Join(baseDir(), "examples", ".terraless"),
			Uploads: []schema.TerralessUpload{
				{
					Type:   "s3",
					Bucket: "example-develop-public",
					Cloudfront: schema.TerralessCloudfront{
						Certificate: "MyCert",
						Domain:      "admin.my-domain.com",
						PriceClass:  "PriceClass_100",
					},
					Provider: "aws-develop-developer",
					Region:   "eu-central-1",
					Source:   "public",
					Target:   "admin",
				},
			},
			HasProvider: map[string]bool{
				"aws": true,
			},
		},
		GlobalConfig: schema.TerralessGlobalConfig{
			Backends: []schema.TerralessBackend{
				{
					Data: map[string]string{
						"bucket":         "my-bucket-name",
						"encrypt":        "true",
						"region":         "eu-central-1",
						"dynamodb_table": "terraform-state-lock",
						"profile":        "my-aws-infrastrucutre-profile-developer",
					},
					Name:     "myBackend",
					Provider: "",
					Type:     "s3",
				},
			},
			Teams: []schema.TerralessTeam{
				{
					Name: "Team A",
					Providers: []schema.TerralessProvider{
						{
							Data: map[string]string{
								"accountId": "01234556678",
								"profile":   "my-aws-profile",
								"region":    "eu-central-1",
							},
							Name: "aws-default",
							Roles: []string{
								"developer",
							},
							Type: "aws",
						},
						{
							Data: map[string]string{
								"accountId": "01234556678",
								"profile":   "my-aws-profile",
								"region":    "eu-central-1",
							},
							Name: "aws-develop",
							Roles: []string{
								"developer",
							},
							Type: "aws",
						},
						{
							Data: map[string]string{
								"accountId": "01234556678",
								"profile":   "my-aws-profile",
								"region":    "eu-central-1",
							},
							Name: "my-aws-infrastrucutre-profile",
							Roles: []string{
								"developer",
							},
							Type: "aws",
						},
					},
				},
			},
		},
	}
	assert.Equal(t, &expected, terralessData)
}
