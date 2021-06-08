package emailverification

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/api"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "emailverification"

type Recipe struct {
	RecipeModule        *supertokens.RecipeModule
	instance            *Recipe
	Config              schema.TypeNormalisedInput
	RecipeInterfaceImpl schema.RecipeInterface
	APIImpl             schema.APIInterface
}

var r Recipe

func NewRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config schema.TypeInput) *Recipe {
	q := supertokens.Querier{}
	instance, _ := q.GetNewInstanceOrThrowError(recipeId)
	recipeModuleInstance := supertokens.NewRecipeModule(recipeId, appInfo)
	verifiedConfig := ValidateAndNormaliseUserInput(appInfo, config)
	recipeInterface := NewRecipeImplementation(*instance)
	return &Recipe{
		RecipeModule:        recipeModuleInstance,
		Config:              verifiedConfig,
		RecipeInterfaceImpl: verifiedConfig.Override.Functions(recipeInterface),
		APIImpl:             verifiedConfig.Override.APIs(api.NewAPIImplementation()),
	}
}

func GetInstanceOrThrowError() (*Recipe, error) {
	if r.instance != nil {
		return r.instance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the SuperTokens.init function?")
}

func RecipeInit(config schema.TypeInput) supertokens.RecipeListFunction {
	return func(appInfo supertokens.NormalisedAppinfo) *supertokens.RecipeModule {
		recipe := NewRecipe(RECIPE_ID, appInfo, config)
		if r.instance == nil {
			r.instance = recipe.instance
			r.Config = recipe.Config
			r.RecipeInterfaceImpl = recipe.RecipeInterfaceImpl
			r.APIImpl = recipe.APIImpl
			return recipe.RecipeModule
		}
		// handle errors.New("Emailverification recipe has already been initialised. Please check your code for bugs.")
		return nil
	}
}

func (r *Recipe) GetAPIsHandled() []supertokens.APIHandled {
	generateEmailVerifyTokenAPI, _ := supertokens.NewNormalisedURLPath(GenerateEmailVerifyTokenAPI)
	emailVerifyAPI, _ := supertokens.NewNormalisedURLPath(EmailVerifyAPI)
	return []supertokens.APIHandled{{
		Method:                 "post",
		PathWithoutAPIBasePath: *generateEmailVerifyTokenAPI,
		ID:                     GenerateEmailVerifyTokenAPI,
		Disabled:               r.APIImpl.GenerateEmailVerifyTokenPOST == nil,
	}, {
		Method:                 "post",
		PathWithoutAPIBasePath: *emailVerifyAPI,
		ID:                     EmailVerifyAPI,
		Disabled:               r.APIImpl.VerifyEmailPOST == nil,
	}, {
		Method:                 "get",
		PathWithoutAPIBasePath: *emailVerifyAPI,
		ID:                     EmailVerifyAPI,
		Disabled:               r.APIImpl.IsEmailVerifiedGET == nil,
	}}
}

func (r *Recipe) HandleAPIRequest(
	id string,
	req *http.Request,
	w http.ResponseWriter,
	path supertokens.NormalisedURLPath,
	method string) {
	options := schema.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeInterfaceImpl,
		Req:                  req,
		Res:                  w,
	}
	if id == GenerateEmailVerifyTokenAPI {
		api.GenerateEmailVerifyToken(r.APIImpl, options)
	}
}
