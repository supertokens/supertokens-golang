package thirdpartyemailpassword

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/recipeimplementation"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "thirdpartyemailpassword"

type Recipe struct {
	RecipeModule            supertokens.RecipeModule
	Config                  models.TypeNormalisedInput
	EmailVerificationRecipe emailverification.Recipe
	emailPasswordRecipe     *emailpassword.Recipe
	thirdPartyRecipe        *thirdparty.Recipe
	RecipeImpl              models.RecipeImplementation
	APIImpl                 models.APIImplementation
}

var r *Recipe

// TODO
func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) (Recipe, error) {
	emailpasswordquerierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(emailpassword.RECIPE_ID)
	if err != nil {
		return Recipe{}, err
	}
	thirdpartyquerierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(thirdparty.RECIPE_ID)
	if err != nil {
		return Recipe{}, err
	}
	recipeImplementation := recipeimplementation.MakeRecipeImplementation(*emailpasswordquerierInstance, thirdpartyquerierInstance)
	verifiedConfig, err := validateAndNormaliseUserInput(recipeImplementation, appInfo, config)
	if err != nil {
		return Recipe{}, err
	}

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, handleAPIRequest, getAllCORSHeaders, getAPIsHandled, handleError)

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
		return nil, errors.New("ThirdPartyEmailPassword recipe has already been initialised. Please check your code for bugs.")
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
	emailpasswordAPIhandled, err := r.emailPasswordRecipe.RecipeModule.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	emailverificationAPIhandled, err := r.EmailVerificationRecipe.RecipeModule.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	apisHandled := append(emailpasswordAPIhandled, emailverificationAPIhandled...)
	if r.thirdPartyRecipe != nil {
		thirdpartyAPIhandled, err := r.thirdPartyRecipe.RecipeModule.GetAPIsHandled()
		if err != nil {
			return nil, err
		}
		apisHandled = append(apisHandled, thirdpartyAPIhandled...)
	}
	return apisHandled, nil
}

func handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path supertokens.NormalisedURLPath, method string) error {
	ok, err := r.emailPasswordRecipe.RecipeModule.ReturnAPIIdIfCanHandleRequest(path, method)
	if err != nil {
		return err
	}
	if ok != nil {
		return r.emailPasswordRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirHandler, path, method)
	}
	if r.thirdPartyRecipe != nil {
		ok, err := r.thirdPartyRecipe.RecipeModule.ReturnAPIIdIfCanHandleRequest(path, method)
		if err != nil {
			return err
		}
		if ok != nil {
			return r.thirdPartyRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirHandler, path, method)
		}
	}
	return r.EmailVerificationRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirHandler, path, method)
}

func getAllCORSHeaders() []string {
	corsHeaders := append(r.EmailVerificationRecipe.RecipeModule.GetAllCORSHeaders(), r.emailPasswordRecipe.RecipeModule.GetAllCORSHeaders()...)
	if r.thirdPartyRecipe != nil {
		corsHeaders = append(corsHeaders, r.thirdPartyRecipe.RecipeModule.GetAllCORSHeaders()...)
	}
	return corsHeaders
}

func handleError(err error, req *http.Request, res http.ResponseWriter) bool {
	handleError := r.emailPasswordRecipe.RecipeModule.HandleError(err, req, res)
	if !handleError {
		if r.thirdPartyRecipe != nil {
			handleError = r.thirdPartyRecipe.RecipeModule.HandleError(err, req, res)
		}
		if !handleError {
			handleError = r.EmailVerificationRecipe.RecipeModule.HandleError(err, req, res)
		}
	}
	return handleError
}
