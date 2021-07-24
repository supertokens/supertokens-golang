package thirdpartyemailpassword

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	evm "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
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

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput, emailVerificationInstance *emailverification.Recipe, thirdPartyInstance *thirdparty.Recipe, emailPasswordInstance *emailpassword.Recipe) (Recipe, error) {
	r = &Recipe{}
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
	r.Config = verifiedConfig
	APIImpl := verifiedConfig.Override.APIs(api.MakeAPIImplementation())
	r.APIImpl = APIImpl
	RecipeImpl := verifiedConfig.Override.Functions(recipeImplementation)
	r.RecipeImpl = RecipeImpl

	var emailVerificationRecipe emailverification.Recipe
	if emailVerificationInstance == nil {
		emailVerificationRecipe, err = emailverification.MakeRecipe(recipeId, appInfo, &verifiedConfig.EmailVerificationFeature)
		if err != nil {
			return Recipe{}, err
		}
	} else {
		emailVerificationRecipe = *emailVerificationInstance
	}
	r.EmailVerificationRecipe = emailVerificationRecipe

	var emailPasswordRecipe emailpassword.Recipe
	if emailPasswordInstance == nil {
		emailPasswordConfig := &epm.TypeInput{
			SessionFeature: &epm.TypeNormalisedInputSessionFeature{
				SetJwtPayload: func(user epm.User, formFields []epm.TypeFormField, action string) map[string]interface{} {
					return verifiedConfig.SessionFeature.SetJwtPayload(models.User{
						ID:         user.ID,
						Email:      user.Email,
						TimeJoined: user.TimeJoined,
					}, models.TypeContext{
						FormFields:                 formFields,
						ThirdPartyAuthCodeResponse: nil,
					}, action)
				},
				SetSessionData: func(user epm.User, formFields []epm.TypeFormField, action string) map[string]interface{} {
					return verifiedConfig.SessionFeature.SetSessionData(models.User{
						ID:         user.ID,
						Email:      user.Email,
						TimeJoined: user.TimeJoined,
					}, models.TypeContext{
						FormFields:                 formFields,
						ThirdPartyAuthCodeResponse: nil,
					}, action)
				},
			},
			SignUpFeature: &epm.TypeInputSignUp{
				FormFields: normalisedToType(verifiedConfig.SignUpFeature.FormFields),
			},
			ResetPasswordUsingTokenFeature: verifiedConfig.ResetPasswordUsingTokenFeature,
			Override: &struct {
				Functions                func(originalImplementation epm.RecipeImplementation) epm.RecipeImplementation
				APIs                     func(originalImplementation epm.APIImplementation) epm.APIImplementation
				EmailVerificationFeature *struct {
					Functions func(originalImplementation evm.RecipeImplementation) evm.RecipeImplementation
					APIs      func(originalImplementation evm.APIImplementation) evm.APIImplementation
				}
			}{
				Functions: func(_ epm.RecipeImplementation) epm.RecipeImplementation {
					return recipeimplementation.MakeEmailPasswordRecipeImplementation(recipeImplementation)
				},
				APIs: func(_ epm.APIImplementation) epm.APIImplementation {
					return api.GetEmailPasswordIterfaceImpl(APIImpl)
				},
				EmailVerificationFeature: nil,
			},
		}
		emailPasswordRecipe, err = emailpassword.MakeRecipe(recipeId, appInfo, emailPasswordConfig, &emailVerificationRecipe)
		if err != nil {
			return Recipe{}, err
		}
	} else {
		emailPasswordRecipe = *emailPasswordInstance
	}
	r.emailPasswordRecipe = &emailPasswordRecipe

	var thirdPartyRecipe *thirdparty.Recipe
	if len(verifiedConfig.Providers) > 0 {
		if thirdPartyInstance == nil {
			thirdPartyConfig := &tpm.TypeInput{
				SessionFeature: &tpm.TypeNormalisedInputSessionFeature{
					SetJwtPayload: func(user tpm.User, thirdPartyAuthCodeResponse interface{}, action string) map[string]interface{} {
						return verifiedConfig.SessionFeature.SetJwtPayload(models.User{
							ID:         user.ID,
							Email:      user.Email,
							TimeJoined: user.TimeJoined,
						}, models.TypeContext{
							FormFields:                 nil,
							ThirdPartyAuthCodeResponse: thirdPartyAuthCodeResponse,
						}, action)
					},
					SetSessionData: func(user tpm.User, thirdPartyAuthCodeResponse interface{}, action string) map[string]interface{} {
						return verifiedConfig.SessionFeature.SetSessionData(models.User{
							ID:         user.ID,
							Email:      user.Email,
							TimeJoined: user.TimeJoined,
						}, models.TypeContext{
							FormFields:                 nil,
							ThirdPartyAuthCodeResponse: thirdPartyAuthCodeResponse,
						}, action)
					},
				},
				SignInAndUpFeature: tpm.TypeInputSignInAndUp{
					Providers: verifiedConfig.Providers,
				},
				Override: &struct {
					Functions                func(originalImplementation tpm.RecipeImplementation) tpm.RecipeImplementation
					APIs                     func(originalImplementation tpm.APIImplementation) tpm.APIImplementation
					EmailVerificationFeature *struct {
						Functions func(originalImplementation evm.RecipeImplementation) evm.RecipeImplementation
						APIs      func(originalImplementation evm.APIImplementation) evm.APIImplementation
					}
				}{
					Functions: func(_ tpm.RecipeImplementation) tpm.RecipeImplementation {
						return recipeimplementation.MakeThirdPartyRecipeImplementation(recipeImplementation)
					},
					APIs: func(_ tpm.APIImplementation) tpm.APIImplementation {
						return api.GetThirdPartyIterfaceImpl(APIImpl)
					},
					EmailVerificationFeature: nil,
				},
			}
			thirdPartyRecipeinstance, err := thirdparty.MakeRecipe(recipeId, appInfo, thirdPartyConfig, &emailVerificationRecipe)
			if err != nil {
				return Recipe{}, err
			}
			thirdPartyRecipe = &thirdPartyRecipeinstance
		} else {
			thirdPartyRecipe = thirdPartyInstance
		}
		r.thirdPartyRecipe = thirdPartyRecipe
	}

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, handleAPIRequest, getAllCORSHeaders, getAPIsHandled, handleError)

	return Recipe{
		RecipeModule:            recipeModuleInstance,
		Config:                  verifiedConfig,
		RecipeImpl:              RecipeImpl,
		APIImpl:                 APIImpl,
		EmailVerificationRecipe: emailVerificationRecipe,
		emailPasswordRecipe:     &emailPasswordRecipe,
		thirdPartyRecipe:        thirdPartyRecipe,
	}, nil
}

func RecipeInit(config *models.TypeInput) supertokens.RecipeListFunction {
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

func handleError(err error, req *http.Request, res http.ResponseWriter) bool {
	handleError := r.emailPasswordRecipe.RecipeModule.HandleError(err, req, res)
	if !handleError {
		if r.thirdPartyRecipe != nil {
			handleError = r.thirdPartyRecipe.RecipeModule.HandleError(err, req, res)
		}
		if !handleError {
			handleError = r.emailPasswordRecipe.RecipeModule.HandleError(err, req, res)
		}
	}
	return handleError
}
