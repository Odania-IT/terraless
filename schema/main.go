package schema

type TerralessData struct {
	ActiveProviders    map[string]TerralessProvider
	Arguments          Arguments
	Config             TerralessConfig
	TerralessProviders []Provider
}
