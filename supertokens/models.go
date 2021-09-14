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
	APIDomain       string
	WebsiteBasePath *string
	APIBasePath     *string
	APIGatewayPath  *string
}

// TODO: change this to a better name.
type Recipe func(appInfo NormalisedAppinfo) (*RecipeModule, error)

type TypeInput struct {
	Supertokens    *ConnectionInfo
	AppInfo        AppInfo
	RecipeList     []Recipe
	Telemetry      *bool
	OnGeneralError func(err error, req *http.Request, res http.ResponseWriter)
}

type ConnectionInfo struct {
	ConnectionURI string
	APIKey        string
}

type APIHandled struct {
	PathWithoutAPIBasePath NormalisedURLPath
	Method                 string
	ID                     string
	Disabled               bool
}
