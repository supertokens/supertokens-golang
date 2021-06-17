package session

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/api"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "session"

var r *models.SessionRecipe = nil

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) models.SessionRecipe {
	querierInstance, _ := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, HandleAPIRequest, GetAllCORSHeaders, GetAPIsHandled)
	verifiedConfig, _ := validateAndNormaliseUserInput(r, appInfo, config)
	recipeImplementation := MakeRecipeImplementation(*querierInstance, verifiedConfig)

	return models.SessionRecipe{
		RecipeModule: recipeModuleInstance,
		Config:       verifiedConfig,
		RecipeImpl:   verifiedConfig.Override.Functions(recipeImplementation),
		APIImpl:      verifiedConfig.Override.APIs(api.MakeAPIImplementation()),
	}
}

func GetInstanceOrThrowError() (*models.SessionRecipe, error) {
	if r != nil {
		return r, nil
	}
	return nil, errors.BadInputError{Msg: "Initialisation not done. Did you forget to call the SuperTokens.init function?"}
}

func RecipeInit(config models.TypeInput) supertokens.RecipeListFunction {
	return func(appInfo supertokens.NormalisedAppinfo) (*supertokens.RecipeModule, error) {
		if r == nil {
			recipe := MakeRecipe(RECIPE_ID, appInfo, &config)
			r = &recipe
			return &r.RecipeModule, nil
		}
		return nil, errors.BadInputError{Msg: "Emailverification recipe has already been initialised. Please check your code for bugs."}
	}
}

// Implement RecipeModule

func GetAPIsHandled() ([]supertokens.APIHandled, error) {
	refreshAPIPath, err := supertokens.NewNormalisedURLPath(RefreshAPIPath)
	if err != nil {
		return nil, err
	}
	signoutAPIPath, err := supertokens.NewNormalisedURLPath(SignoutAPIPath)
	if err != nil {
		return nil, err
	}
	return []supertokens.APIHandled{{
		Method:                 "post",
		PathWithoutAPIBasePath: *refreshAPIPath,
		ID:                     RefreshAPIPath,
		Disabled:               r.APIImpl.RefreshPOST == nil,
	}, {
		Method:                 "post",
		PathWithoutAPIBasePath: *signoutAPIPath,
		ID:                     SignoutAPIPath,
		Disabled:               r.APIImpl.SignOutPOST == nil,
	}}, nil
}

func HandleAPIRequest(id string, req *http.Request, res http.ResponseWriter, thierhandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string) error {
	options := models.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
		OtherHandler:         thierhandler,
	}
	if id == RefreshAPIPath {
		api.HandleRefreshAPI(r.APIImpl, options)
	} else {
		return api.SignOutAPI(r.APIImpl, options)
	}
	return nil
}

func GetAllCORSHeaders() []string {
	return []string{antiCsrfHeaderKey, ridHeaderKey}
}
