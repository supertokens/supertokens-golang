package supertokens

type RecipeModule struct {
	RecipeID string
	AppInfo  NormalisedAppinfo
}

func (r RecipeModule) GetRecipeID() string {
	return r.RecipeID
}

func (r RecipeModule) GetAppInfo() NormalisedAppinfo {
	return r.AppInfo
}
