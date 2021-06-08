package supertokens

import "net/http"

type RecipeModule struct {
	recipeID          string
	appInfo           NormalisedAppinfo
	HandleAPIRequest  func(id string, req *http.Request, w http.ResponseWriter, path NormalisedURLPath, method string)
	GetAllCORSHeaders func() []string
	GetAPIsHandled    func() []APIHandled
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
	apisHandled := r.GetAPIsHandled()
	for _, APIshandled := range apisHandled {
		pathAppend := r.appInfo.APIBasePath.AppendPath(APIshandled.PathWithoutAPIBasePath)
		if !APIshandled.Disabled && APIshandled.Method == method && pathAppend.Equals(path) {
			return APIshandled.ID
		}
	}
	return ""
}
