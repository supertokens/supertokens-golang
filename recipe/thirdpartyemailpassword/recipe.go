package thirdpartyemailpassword

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
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
	RecipeImpl              models.RecipeInterface
	APIImpl                 models.APIInterface
}

var r *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput, emailVerificationInstance *emailverification.Recipe, thirdPartyInstance *thirdparty.Recipe, emailPasswordInstance *emailpassword.Recipe) (Recipe, error) {
	r = &Recipe{}
	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, handleAPIRequest, getAllCORSHeaders, getAPIsHandled, handleError)

	verifiedConfig, err := validateAndNormaliseUserInput(r, appInfo, config)
	if err != nil {
		return Recipe{}, err
	}
	r.Config = verifiedConfig
	{
		emailpasswordquerierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(emailpassword.RECIPE_ID)
		if err != nil {
			return Recipe{}, err
		}
		thirdpartyquerierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(thirdparty.RECIPE_ID)
		if err != nil {
			return Recipe{}, err
		}

		r.RecipeImpl = verifiedConfig.Override.Functions(recipeimplementation.MakeRecipeImplementation(*emailpasswordquerierInstance, thirdpartyquerierInstance))
	}
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	if emailVerificationInstance == nil {
		emailVerificationRecipe, err := emailverification.MakeRecipe(recipeId, appInfo, &verifiedConfig.EmailVerificationFeature)
		if err != nil {
			return Recipe{}, err
		}
		r.EmailVerificationRecipe = emailVerificationRecipe

	} else {
		r.EmailVerificationRecipe = *emailVerificationInstance
	}

	var emailPasswordRecipe emailpassword.Recipe
	if emailPasswordInstance == nil {
		emailPasswordConfig := &epm.TypeInput{
			SignUpFeature:                  verifiedConfig.SignUpFeature,
			ResetPasswordUsingTokenFeature: verifiedConfig.ResetPasswordUsingTokenFeature,
			Override: &epm.OverrideStruct{
				Functions: func(_ epm.RecipeInterface) epm.RecipeInterface {
					return recipeimplementation.MakeEmailPasswordRecipeImplementation(r.RecipeImpl)
				},
				APIs: func(_ epm.APIInterface) epm.APIInterface {
					return api.GetEmailPasswordIterfaceImpl(r.APIImpl)
				},
				EmailVerificationFeature: nil,
			},
		}
		emailPasswordRecipe, err = emailpassword.MakeRecipe(recipeId, appInfo, emailPasswordConfig, &r.EmailVerificationRecipe)
		if err != nil {
			return Recipe{}, err
		}
		r.emailPasswordRecipe = &emailPasswordRecipe
	} else {
		r.emailPasswordRecipe = emailPasswordInstance
	}

	if len(verifiedConfig.Providers) > 0 {
		if thirdPartyInstance == nil {
			thirdPartyConfig := &tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: verifiedConfig.Providers,
				},
				Override: &tpmodels.OverrideStruct{
					Functions: func(_ tpmodels.RecipeInterface) tpmodels.RecipeInterface {
						return recipeimplementation.MakeThirdPartyRecipeImplementation(r.RecipeImpl)
					},
					APIs: func(_ tpmodels.APIInterface) tpmodels.APIInterface {
						return api.GetThirdPartyIterfaceImpl(r.APIImpl)
					},
					EmailVerificationFeature: nil,
				},
			}
			thirdPartyRecipeinstance, err := thirdparty.MakeRecipe(recipeId, appInfo, thirdPartyConfig, &r.EmailVerificationRecipe)
			if err != nil {
				return Recipe{}, err
			}
			r.thirdPartyRecipe = &thirdPartyRecipeinstance
		} else {
			r.thirdPartyRecipe = thirdPartyInstance
		}
	}

	return *r, nil
}

func recipeInit(config *models.TypeInput) supertokens.RecipeListFunction {
	return func(appInfo supertokens.NormalisedAppinfo) (*supertokens.RecipeModule, error) {
		if r == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, nil, nil, nil)
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

func handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	handleError, err := r.emailPasswordRecipe.RecipeModule.HandleError(err, req, res)
	if err != nil || handleError {
		return handleError, err
	}
	if r.thirdPartyRecipe != nil {
		handleError, err = r.thirdPartyRecipe.RecipeModule.HandleError(err, req, res)
		if err != nil || handleError {
			return handleError, err
		}
	}
	return r.EmailVerificationRecipe.RecipeModule.HandleError(err, req, res)
}

func (r *Recipe) getEmailForUserId(userID string) (string, error) {
	userInfo, err := r.RecipeImpl.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	if userInfo == nil {
		return "", errors.New("Unknown User ID provided")
	}
	return userInfo.Email, nil
}
