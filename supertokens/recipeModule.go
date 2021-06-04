package supertokens

import "net/http"

type RecipeModule struct {
	recipeID string
	appInfo  NormalisedAppinfo
}

func NewRecipeModule(recipeId string, appInfo NormalisedAppinfo) *RecipeModule {
	return &RecipeModule{
		recipeID: recipeId,
		appInfo:  appInfo,
	}
}

func (r *RecipeModule) GetRecipeID() string {
	return r.recipeID
}

func (r *RecipeModule) GetAppInfo() NormalisedAppinfo {
	return r.appInfo
}

func (r *RecipeModule) ReturnAPIIdIfCanHandleRequest(path NormalisedURLPath, method string) string {
	// todo
	return ""
}

func (r *RecipeModule) handleAPIRequest(
	id string,
	req *http.Request,
	w http.ResponseWriter,
	path NormalisedURLPath,
	method http.HandlerFunc) {
}

// func (r *RecipeModule) getAPIsHandled() []APIHandled {

// }
