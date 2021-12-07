package openidmodels

import "net/http"

type APIOptions struct {
	RecipeImplementation RecipeInterface
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
}

type APIInterface struct {
	GetOpenIdDiscoveryConfigurationGET *func(options APIOptions) (GetOpenIdDiscoveryConfigurationResponse, error)
}
