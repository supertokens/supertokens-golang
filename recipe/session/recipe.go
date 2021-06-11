package session

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/api"
	"github.com/supertokens/supertokens-golang/recipe/session/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "session"

type SessionRecipe struct {
	RecipeModule supertokens.RecipeModule
	Config       schema.TypeNormalisedInput
	RecipeImpl   schema.RecipeImplementation
	APIImpl      schema.APIImplementation
}

var r *SessionRecipe = nil

func NewRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *schema.TypeInput) SessionRecipe {
	querierInstance, _ := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, HandleAPIRequest, GetAllCORSHeaders, GetAPIsHandled)
	verifiedConfig, _ := ValidateAndNormaliseUserInput(r, appInfo, config)

	recipeImplementation := MakeRecipeImplementation(*querierInstance, verifiedConfig)

	return SessionRecipe{
		RecipeModule: recipeModuleInstance,
		Config:       verifiedConfig,
		RecipeImpl:   verifiedConfig.Override.Functions(recipeImplementation),
		APIImpl:      verifiedConfig.Override.APIs(api.MakeAPIImplementation()),
	}
}

// implement RecipeModule

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

func HandleAPIRequest(id string, req *http.Request, w http.ResponseWriter, path supertokens.NormalisedURLPath, method string) error {
	// options := schema.APIOptions{
	// 	Config:               r.Config,
	// 	RecipeID:             r.RecipeModule.GetRecipeID(),
	// 	RecipeImplementation: r.RecipeInterfaceImpl,
	// 	Req:                  req,
	// 	Res:                  w,
	// }
	if id == RefreshAPIPath {
		// api.handleRefreshAPI(r.APIImpl, options)
	} else {
		// api.signOutAPI(r.APIImpl, options)
	}
	return nil
}

func GetAllCORSHeaders() []string {
	return []string{}
}
