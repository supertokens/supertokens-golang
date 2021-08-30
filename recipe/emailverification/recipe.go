package emailverification

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/api"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "emailverification"

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       models.TypeNormalisedInput
	RecipeImpl   models.RecipeInterface
	APIImpl      models.APIInterface
}

var r = &Recipe{}

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) (Recipe, error) {
	verifiedConfig := validateAndNormaliseUserInput(appInfo, *config)
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance)
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, handleAPIRequest, getAllCORSHeaders, getAPIsHandled, nil)
	r.RecipeModule = recipeModuleInstance

	return Recipe{
		RecipeModule: recipeModuleInstance,
		Config:       verifiedConfig,
		RecipeImpl:   verifiedConfig.Override.Functions(recipeImplementation),
		APIImpl:      verifiedConfig.Override.APIs(api.MakeAPIImplementation()),
	}, nil
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if r != nil {
		return r, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the init function?")
}

func recipeInit(config *models.TypeInput) supertokens.RecipeListFunction {
	return func(appInfo supertokens.NormalisedAppinfo) (*supertokens.RecipeModule, error) {
		if r == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config)
			if err != nil {
				return nil, err
			}
			r = &recipe
			return &r.RecipeModule, nil
		}
		return nil, errors.New("Emailverification recipe has already been initialised. Please check your code for bugs.")
	}
}

// implement RecipeModule

func getAPIsHandled() ([]supertokens.APIHandled, error) {
	generateEmailVerifyTokenAPINormalised, err := supertokens.NewNormalisedURLPath(generateEmailVerifyTokenAPI)
	if err != nil {
		return nil, err
	}
	emailVerifyAPINormalised, err := supertokens.NewNormalisedURLPath(emailVerifyAPI)
	if err != nil {
		return nil, err
	}

	return []supertokens.APIHandled{{
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: *generateEmailVerifyTokenAPINormalised,
		ID:                     generateEmailVerifyTokenAPI,
		Disabled:               r.APIImpl.GenerateEmailVerifyTokenPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: *emailVerifyAPINormalised,
		ID:                     emailVerifyAPI,
		Disabled:               r.APIImpl.VerifyEmailPOST == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: *emailVerifyAPINormalised,
		ID:                     emailVerifyAPI,
		Disabled:               r.APIImpl.IsEmailVerifiedGET == nil,
	}}, nil
}

func handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string) error {
	options := models.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirHandler,
	}
	if id == generateEmailVerifyTokenAPI {
		return api.GenerateEmailVerifyToken(r.APIImpl, options)
	} else {
		return api.EmailVerify(r.APIImpl, options)
	}
}

func getAllCORSHeaders() []string {
	return []string{}
}
