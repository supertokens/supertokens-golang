package supertokens

import "net/http"

type RecipeModule struct {
	recipeID          string
	appInfo           NormalisedAppinfo
	HandleAPIRequest  func(ID string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path NormalisedURLPath, method string) error
	GetAllCORSHeaders func() []string
	GetAPIsHandled    func() ([]APIHandled, error)
	HandleError       func(err error, req *http.Request, res http.ResponseWriter) bool
}

func MakeRecipeModule(
	recipeId string,
	appInfo NormalisedAppinfo,
	handleAPIRequest func(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path NormalisedURLPath, method string) error,
	getAllCORSHeaders func() []string,
	getAPIsHandled func() ([]APIHandled, error),
	handleError func(err error, req *http.Request, res http.ResponseWriter) bool) RecipeModule {
	return RecipeModule{
		recipeID:          recipeId,
		appInfo:           appInfo,
		HandleAPIRequest:  handleAPIRequest,
		GetAllCORSHeaders: getAllCORSHeaders,
		GetAPIsHandled:    getAPIsHandled,
		HandleError:       handleError,
	}
}

func (r RecipeModule) GetRecipeID() string {
	return r.recipeID
}

func (r RecipeModule) GetAppInfo() NormalisedAppinfo {
	return r.appInfo
}

func (r *RecipeModule) ReturnAPIIdIfCanHandleRequest(path NormalisedURLPath, method string) (*string, error) {
	apisHandled, err := r.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	for _, APIshandled := range apisHandled {
		pathAppend := r.appInfo.APIBasePath.AppendPath(APIshandled.PathWithoutAPIBasePath)
		if !APIshandled.Disabled && APIshandled.Method == method && pathAppend.Equals(path) {
			return &APIshandled.ID, nil
		}
	}
	return nil, nil
}
