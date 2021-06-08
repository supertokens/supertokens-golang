package session

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "session"

// WIP

type Recipe struct {
	RecipeModule        *supertokens.RecipeModule
	instance            *Recipe
	Config              schema.TypeNormalisedInput
	RecipeInterfaceImpl schema.RecipeInterface
	APIImpl             schema.APIInterface
}

var r Recipe

func NewRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config schema.TypeInput) *Recipe {
	// q := supertokens.Querier{}
	// instance, _ := q.GetNewInstanceOrThrowError(recipeId)
	recipeModuleInstance := supertokens.NewRecipeModule(recipeId, appInfo)
	recipeModuleInstance.GetAPIsHandled = func() []supertokens.APIHandled {
		return GetAPIsHandled()
	}
	recipeModuleInstance.HandleAPIRequest = func(id string, req *http.Request, w http.ResponseWriter, path supertokens.NormalisedURLPath, method string) {
		HandleAPIRequest(id, req, w, path, method)
	}
	recipeModuleInstance.GetAllCORSHeaders = func() []string {
		return GetAllCORSHeaders()
	}
	// verifiedConfig := ValidateAndNormaliseUserInput(appInfo, config)
	// recipeInterface := NewRecipeImplementation(*instance)
	// return &Recipe{
	// 	RecipeModule:        recipeModuleInstance,
	// 	Config:              verifiedConfig,
	// 	RecipeInterfaceImpl: verifiedConfig.Override.Functions(recipeInterface),
	// 	APIImpl:             verifiedConfig.Override.APIs(api.NewAPIImplementation()),
	// }
	return nil
}

// implement RecipeModule

func GetAPIsHandled() []supertokens.APIHandled {
	refreshAPIPath, _ := supertokens.NewNormalisedURLPath(RefreshAPIPath)
	signoutAPIPath, _ := supertokens.NewNormalisedURLPath(SignoutAPIPath)
	return []supertokens.APIHandled{{
		Method:                 "post",
		PathWithoutAPIBasePath: *refreshAPIPath,
		ID:                     RefreshAPIPath,
		Disabled:               r.APIImpl.RefreshPOST == nil,
	}, {
		Method:                 "post",
		PathWithoutAPIBasePath: *signoutAPIPath,
		ID:                     SignoutAPIPath,
		// Disabled:               r.APIImpl.signOutPOST == nil,
	}}
}

func HandleAPIRequest(id string, req *http.Request, w http.ResponseWriter, path supertokens.NormalisedURLPath, method string) {
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
}

func GetAllCORSHeaders() []string {
	return []string{}
}
