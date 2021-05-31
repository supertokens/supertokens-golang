package supertokens

type RecipeModule struct {
	RecipeID string
	AppInfo  NormalisedAppinfo
}

func NewRecipeModule(recipeId string, appInfo NormalisedAppinfo) *RecipeModule {
	return &RecipeModule{
		RecipeID: recipeId,
		AppInfo:  appInfo,
	}
}

func (r *RecipeModule) GetRecipeID() string {
	return r.RecipeID
}

func (r *RecipeModule) GetAppInfo() NormalisedAppinfo {
	return r.AppInfo
}

func (r *RecipeModule) ReturnAPIIdIfCanHandleRequest(path NormalisedURLPath, method string) string {
	// apisHandled
	return ""
}


// func (r *RecipeModule) getAPIsHandled() []APIHandled {

// }
