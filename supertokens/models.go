package supertokens

import "net/http"

type NormalisedAppinfo struct {
	AppName         string
	WebsiteDomain   NormalisedURLDomain
	APIDomain       NormalisedURLDomain
	APIBasePath     NormalisedURLPath
	APIGatewayPath  NormalisedURLPath
	WebsiteBasePath NormalisedURLPath
}

type AppInfo struct {
	AppName         string
	WebsiteDomain   string
	WebsiteBasePath *string
	APIDomain       string
	APIBasePath     *string
	APIGatewayPath  *string
}

type RecipeListFunction func(appInfo NormalisedAppinfo) (*RecipeModule, error)

type TypeInput struct {
	Supertokens     *SupertokenTypeInput
	AppInfo        AppInfo
	RecipeList     []RecipeListFunction
	Telemetry      *bool
	OnGeneralError func(err error, req *http.Request, res http.ResponseWriter)
}

type SupertokenTypeInput struct {
	ConnectionURI string
	APIKey        *string
}

type APIHandled struct {
	PathWithoutAPIBasePath NormalisedURLPath
	Method                 string
	ID                     string
	Disabled               bool
}
