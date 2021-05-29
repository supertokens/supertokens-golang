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

func RecipeListFunction(appInfo NormalisedAppinfo, isInServerlessEnv bool) RecipeModule {
	return RecipeModule{}
}

type TypeInput struct {
	Supertoken        SupertokenTypeInput
	AppInfo           AppInfo
	RecipeList        []func()
	Telemetry         bool
	IsInServerlessEnv bool
}

type SupertokenTypeInput struct {
	ConnectionURI string
	APIKey        string
}
