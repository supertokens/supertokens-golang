package emailpassword

import (
	defaultErrors "errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/api"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/constants"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/errors"
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

var r *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) (Recipe, error) {
	r = &Recipe{}
	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance)
	verifiedConfig := validateAndNormaliseUserInput(recipeImplementation, appInfo, config)
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	emailVerificationRecipe, err := emailverification.MakeRecipe(recipeId, appInfo, &verifiedConfig.EmailVerificationFeature)
	if err != nil {
		return Recipe{}, err
	}
	r.EmailVerificationRecipe = emailVerificationRecipe

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, handleAPIRequest, getAllCORSHeaders, getAPIsHandled, handleError)

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
	signUpAPI, err := supertokens.NewNormalisedURLPath(constants.SignUpAPI)
	if err != nil {
		return nil, err
	}
	signInAPI, err := supertokens.NewNormalisedURLPath(constants.SignInAPI)
	if err != nil {
		return nil, err
	}
	generatePasswordResetTokenAPI, err := supertokens.NewNormalisedURLPath(constants.GeneratePasswordResetTokenAPI)
	if err != nil {
		return nil, err
	}
	passwordResetAPI, err := supertokens.NewNormalisedURLPath(constants.PasswordResetAPI)
	if err != nil {
		return nil, err
	}
	signupEmailExistsAPI, err := supertokens.NewNormalisedURLPath(constants.SignupEmailExistsAPI)
	if err != nil {
		return nil, err
	}
	emailverificationAPIhandled, err := r.EmailVerificationRecipe.RecipeModule.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	return append([]supertokens.APIHandled{{
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: *signUpAPI,
		ID:                     constants.SignUpAPI,
		Disabled:               r.APIImpl.SignUpPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: *signInAPI,
		ID:                     constants.SignInAPI,
		Disabled:               r.APIImpl.SignInPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: *generatePasswordResetTokenAPI,
		ID:                     constants.GeneratePasswordResetTokenAPI,
		Disabled:               r.APIImpl.GeneratePasswordResetTokenPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: *passwordResetAPI,
		ID:                     constants.PasswordResetAPI,
		Disabled:               r.APIImpl.PasswordResetPOST == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: *signupEmailExistsAPI,
		ID:                     constants.SignupEmailExistsAPI,
		Disabled:               r.APIImpl.EmailExistsGET == nil,
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
	if id == constants.SignUpAPI {
		return api.SignUpAPI(r.APIImpl, options)
	} else if id == constants.SignInAPI {
		return api.SignInAPI(r.APIImpl, options)
	} else if id == constants.GeneratePasswordResetTokenAPI {
		return api.GeneratePasswordResetToken(r.APIImpl, options)
	} else if id == constants.PasswordResetAPI {
		return api.PasswordReset(r.APIImpl, options)
	} else if id == constants.SignupEmailExistsAPI {
		return api.EmailExists(r.APIImpl, options)
	}
	return r.EmailVerificationRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirHandler, path, method)
}

func getAllCORSHeaders() []string {
	return r.EmailVerificationRecipe.RecipeModule.GetAllCORSHeaders()
}

func handleError(err error, req *http.Request, res http.ResponseWriter) bool {
	if defaultErrors.As(err, &errors.FieldError{}) {
		errs := err.(errors.FieldError)
		supertokens.Send200Response(res, map[string]interface{}{
			"status":     "FIELD_ERROR",
			"formFields": errs.Payload,
		})
		return true
	}
	return false
}
