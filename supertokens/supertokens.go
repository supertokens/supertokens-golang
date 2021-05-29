package supertokens

type SuperTokens struct {
	AppInfo           NormalisedAppinfo
	IsInServerlessEnv bool
	RecipeModules     []RecipeModule
}

func (s SuperTokens) init(config TypeInput) {

}
