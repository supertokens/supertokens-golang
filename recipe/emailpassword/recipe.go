package emailpassword

import (
	defaultErrors "errors"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "emailpassword"

type Recipe struct {
	RecipeModule            supertokens.RecipeModule
	Config                  models.TypeNormalisedInput
	RecipeImpl              models.RecipeImplementation
	APIImpl                 models.APIImplementation
	EmailVerificationRecipe emailverification.Recipe
}

var r *Recipe = nil

// func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config models.TypeInput, recipes struct {
// 	EmailVerificationInstance *emailverification.Recipe
// }) (Recipe, error) {
// 	verifiedConfig := validateAndNormaliseUserInput(r.RecipeImpl, appInfo, config)
// 	emailVerificationRecipe, err := emailverification.MakeRecipe(recipeId, appInfo, verifiedConfig.EmailVerificationFeature)
// 	if recipes.EmailVerificationInstance != nil {
// 		emailVerificationRecipe = *recipes.EmailVerificationInstance
// 	}
// 	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
// 	if err != nil {
// 		return Recipe{}, err
// 	}
// 	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, HandleAPIRequest, GetAllCORSHeaders, GetAPIsHandled)
// 	recipeImplementation := MakeRecipeImplementation(*querierInstance)

// 	return Recipe{
// 		RecipeModule:            recipeModuleInstance,
// 		Config:                  verifiedConfig,
// 		RecipeImpl:              verifiedConfig.Override.Functions(recipeImplementation),
// 		APIImpl:                 verifiedConfig.Override.APIs(api.MakeAPIImplementation()),
// 		EmailVerificationRecipe: emailVerificationRecipe,
// 	}, nil
// }

func GetRecipeInstanceOrThrowError() (*Recipe, error) {
	if r != nil {
		return r, nil
	}
	return nil, defaultErrors.New("Initialisation not done. Did you forget to call the init function?")
}
