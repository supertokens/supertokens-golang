package supertokens

type NormalisedAppinfo struct {
	AppName         string
	WebsiteDomain   NormalisedURLDomain
	APIDomain       NormalisedURLDomain
	APIBasePath     NormalisedURLPath
	APIGatewayPath  NormalisedURLPath
	WebsiteBasePath NormalisedURLPath
}

type AppInfo struct {
	appName         string
	websiteDomain   string
	websiteBasePath string
	apiDomain       string
	apiBasePath     string
	apiGatewayPath  string
}

type RecipeListFunction func(appInfo NormalisedAppinfo) *RecipeModule

type TypeInput struct {
	Supertoken SupertokenTypeInput
	AppInfo    AppInfo
	RecipeList []RecipeListFunction
	Telemetry  bool
}

type SupertokenTypeInput struct {
	ConnectionURI string
	APIKey        string
}

type APIHandled struct {
	PathWithoutAPIBasePath NormalisedURLPath
	Method                 string
	ID                     string
	Disabled               bool
}
