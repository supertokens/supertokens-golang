package recipe

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/api"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/recipeimplementation"
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
	return &Recipe{
		RecipeModule:        recipeModuleInstance,
		Config:              emailverification.ValidateAndNormaliseUserInput(appInfo, config), // doubt about `this`
		RecipeInterfaceImpl: config.Override.Functions(*recipeimplementation.NewRecipeImplementation(*instance)),
		APIImpl:             config.Override.APIs(api.APIImplementation{}),
	}
}

func GetInstanceOrThrowError() (*Recipe, error) {
	if r.instance != nil {
		return r.instance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the SuperTokens.init function?")
}

// discussion required
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
	generateEmailVerifyTokenAPI, _ := supertokens.NewNormalisedURLPath(emailverification.GenerateEmailVerifyTokenAPI)
	emailVerifyAPI, _ := supertokens.NewNormalisedURLPath(emailverification.EmailVerifyAPI)
	return []supertokens.APIHandled{{
		Method:                 "post",
		PathWithoutAPIBasePath: *generateEmailVerifyTokenAPI,
		ID:                     emailverification.GenerateEmailVerifyTokenAPI,
		Disabled:               false, // doubt
	}, {
		Method:                 "post",
		PathWithoutAPIBasePath: *emailVerifyAPI,
		ID:                     emailverification.EmailVerifyAPI,
		Disabled:               false, // doubt
	}, {
		Method:                 "get",
		PathWithoutAPIBasePath: *emailVerifyAPI,
		ID:                     emailverification.EmailVerifyAPI,
		Disabled:               false, // doubt
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
	if id == emailverification.GenerateEmailVerifyTokenAPI {
		api.GenerateEmailVerifyToken(r.APIImpl, options)
	}
}
