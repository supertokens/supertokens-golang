package thirdparty

import (
	defaultErrors "errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "thirdparty"

type Recipe struct {
	RecipeModule            supertokens.RecipeModule
	Config                  models.TypeNormalisedInput
	RecipeImpl              models.RecipeImplementation
	APIImpl                 models.APIImplementation
	EmailVerificationRecipe emailverification.Recipe
	Providers               []models.TypeProvider
}

var r *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) (Recipe, error) {
	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, handleAPIRequest, getAllCORSHeaders, getAPIsHandled, handleError)
	recipeImplementation := makeRecipeImplementation(*querierInstance)
	verifiedConfig := validateAndNormaliseUserInput(recipeImplementation, appInfo, config)
	emailVerificationRecipe, err := emailverification.MakeRecipe(recipeId, appInfo, &verifiedConfig.EmailVerificationFeature)
	if err != nil {
		return Recipe{}, err
	}

	return Recipe{
		RecipeModule:            recipeModuleInstance,
		Config:                  verifiedConfig,
		RecipeImpl:              verifiedConfig.Override.Functions(recipeImplementation),
		APIImpl:                 verifiedConfig.Override.APIs(api.MakeAPIImplementation()),
		EmailVerificationRecipe: emailVerificationRecipe,
	}, nil
}

func RecipeInit(config *models.TypeInput) supertokens.RecipeListFunction {
	return func(appInfo supertokens.NormalisedAppinfo) (*supertokens.RecipeModule, error) {
		if r == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config)
			if err != nil {
				return nil, err
			}
			r = &recipe
			return &r.RecipeModule, nil
		}
		return nil, defaultErrors.New("Emailpassword recipe has already been initialised. Please check your code for bugs.")
	}
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if r != nil {
		return r, nil
	}
	return nil, defaultErrors.New("Initialisation not done. Did you forget to call the init function?")
}

// implement RecipeModule

func getAPIsHandled() ([]supertokens.APIHandled, error) {
	signInUpAPI, err := supertokens.NewNormalisedURLPath(SignInUpAPI)
	if err != nil {
		return nil, err
	}
	authorisationAPI, err := supertokens.NewNormalisedURLPath(AuthorisationAPI)
	if err != nil {
		return nil, err
	}
	emailverificationAPIhandled, err := r.EmailVerificationRecipe.RecipeModule.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	return append([]supertokens.APIHandled{{
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: *signInUpAPI,
		ID:                     SignInUpAPI,
		Disabled:               r.APIImpl.SignInUpPOST == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: *authorisationAPI,
		ID:                     AuthorisationAPI,
		Disabled:               r.APIImpl.AuthorisationUrlGET == nil,
	}}, emailverificationAPIhandled...), nil
}

func handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path supertokens.NormalisedURLPath, method string) error {
	options := models.APIOptions{
		Config:               r.Config,
		OtherHandler:         theirHandler,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
	}
	if id == SignInUpAPI {
		return api.SignInUpAPI(r.APIImpl, options)
	} else if id == AuthorisationAPI {
		return api.AuthorisationUrlAPI(r.APIImpl, options)
	}
	return r.EmailVerificationRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirHandler, path, method)
}

func getAllCORSHeaders() []string {
	return r.EmailVerificationRecipe.RecipeModule.GetAllCORSHeaders()
}

func handleError(err error, req *http.Request, res http.ResponseWriter) bool {
	return false
}
