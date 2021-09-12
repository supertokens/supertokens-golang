package thirdparty

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "thirdparty"

type Recipe struct {
	RecipeModule            supertokens.RecipeModule
	Config                  tpmodels.TypeNormalisedInput
	RecipeImpl              tpmodels.RecipeInterface
	APIImpl                 tpmodels.APIInterface
	EmailVerificationRecipe emailverification.Recipe
	Providers               []tpmodels.TypeProvider
}

var r *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *tpmodels.TypeInput, emailVerificationInstance *emailverification.Recipe) (Recipe, error) {
	r = &Recipe{}

	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, handleAPIRequest, getAllCORSHeaders, getAPIsHandled, handleError)

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	verifiedConfig, err := validateAndNormaliseUserInput(r, appInfo, config)
	if err != nil {
		return Recipe{}, err
	}
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())
	r.RecipeImpl = verifiedConfig.Override.Functions(MakeRecipeImplementation(*querierInstance))
	r.Providers = config.SignInAndUpFeature.Providers

	if emailVerificationInstance == nil {
		emailVerificationRecipe, err := emailverification.MakeRecipe(recipeId, appInfo, &verifiedConfig.EmailVerificationFeature)
		if err != nil {
			return Recipe{}, err
		}
		r.EmailVerificationRecipe = emailVerificationRecipe

	} else {
		r.EmailVerificationRecipe = *emailVerificationInstance
	}

	return *r, nil
}

func recipeInit(config *tpmodels.TypeInput) supertokens.RecipeListFunction {
	return func(appInfo supertokens.NormalisedAppinfo) (*supertokens.RecipeModule, error) {
		if r == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, nil)
			if err != nil {
				return nil, err
			}
			r = &recipe
			return &r.RecipeModule, nil
		}
		return nil, errors.New("ThirdParty recipe has already been initialised. Please check your code for bugs.")
	}
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if r != nil {
		return r, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the init function?")
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
		PathWithoutAPIBasePath: signInUpAPI,
		ID:                     SignInUpAPI,
		Disabled:               r.APIImpl.SignInUpPOST == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: authorisationAPI,
		ID:                     AuthorisationAPI,
		Disabled:               r.APIImpl.AuthorisationUrlGET == nil,
	}}, emailverificationAPIhandled...), nil
}

func handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path supertokens.NormalisedURLPath, method string) error {
	options := tpmodels.APIOptions{
		Config:                                r.Config,
		OtherHandler:                          theirHandler,
		RecipeID:                              r.RecipeModule.GetRecipeID(),
		RecipeImplementation:                  r.RecipeImpl,
		EmailVerificationRecipeImplementation: r.EmailVerificationRecipe.RecipeImpl,
		Providers:                             r.Providers,
		Req:                                   req,
		Res:                                   res,
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

func handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	return r.EmailVerificationRecipe.RecipeModule.HandleError(err, req, res)
}

func (r *Recipe) getEmailForUserId(userID string) (string, error) {
	userInfo, err := r.RecipeImpl.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	if userInfo == nil {
		return "", errors.New("unknown User ID provided")
	}
	return userInfo.Email, nil
}
