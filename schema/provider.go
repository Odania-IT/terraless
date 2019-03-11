package schema

type Provider struct {
	CanHandle             CanHandleFunc
	// FilterActiveProviders FilterActiveProfilesFunc
	PrepareSession        PrepareSessionFunc
}

type CanHandleFunc func(resourceType string) bool
// type FilterActiveProfilesFunc func(terralessConfig TerralessConfig) map[string]TerralessProvider
type PrepareSessionFunc func(terralessConfig TerralessConfig)
