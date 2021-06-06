package emailverification

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/receipe/emailverification/api"
	"github.com/supertokens/supertokens-golang/receipe/emailverification/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "emailverification"

type Recipe struct {
	Instance            supertokens.RecipeModule
	RecipeID            string
	Config              schema.TypeNormalisedInput
	RecipeInterfaceImpl schema.RecipeInterface
	APIImpl             schema.APIInterface
}

func NewRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config schema.TypeInput) *Recipe {
	return &Recipe{
		// config: config,
		// recipeInterfaceImpl: config.override.functions(NewRecipeImplementation()),
	}
}

func (r *Recipe) GetAPIsHandled() []supertokens.APIHandled {
	generateEmailVerifyTokenAPI, _ := supertokens.NewNormalisedURLPath(GenerateEmailVerifyTokenAPI)
	emailVerifyAPI, _ := supertokens.NewNormalisedURLPath(EmailVerifyAPI)
	return []supertokens.APIHandled{{
		Method:                 "post",
		PathWithoutAPIBasePath: *generateEmailVerifyTokenAPI,
		ID:                     GenerateEmailVerifyTokenAPI,
		Disabled:               false, // doubt
	}, {
		Method:                 "post",
		PathWithoutAPIBasePath: *emailVerifyAPI,
		ID:                     EmailVerifyAPI,
		Disabled:               false, // doubt
	}, {
		Method:                 "get",
		PathWithoutAPIBasePath: *emailVerifyAPI,
		ID:                     EmailVerifyAPI,
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
		RecipeID:             r.Instance.GetRecipeID(),
		RecipeImplementation: r.RecipeInterfaceImpl,
		Req:                  req,
		Res:                  w,
	}
	if id == GenerateEmailVerifyTokenAPI {
		api.GenerateEmailVerifyToken(r.APIImpl, options)
	}
}
