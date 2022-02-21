/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package emailpassword

import (
	defaultErrors "errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/api"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/constants"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/errors"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "emailpassword"

type Recipe struct {
	RecipeModule            supertokens.RecipeModule
	Config                  epmodels.TypeNormalisedInput
	RecipeImpl              epmodels.RecipeInterface
	APIImpl                 epmodels.APIInterface
	EmailVerificationRecipe emailverification.Recipe
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *epmodels.TypeInput, emailVerificationInstance *emailverification.Recipe, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}
	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, r.handleError, onGeneralError)

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	verifiedConfig := validateAndNormaliseUserInput(r, appInfo, config)
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())
	r.RecipeImpl = verifiedConfig.Override.Functions(MakeRecipeImplementation(*querierInstance))

	if emailVerificationInstance == nil {
		emailVerificationRecipe, err := emailverification.MakeRecipe(recipeId, appInfo, verifiedConfig.EmailVerificationFeature, onGeneralError)
		if err != nil {
			return Recipe{}, err
		}
		r.EmailVerificationRecipe = emailVerificationRecipe

	} else {
		r.EmailVerificationRecipe = *emailVerificationInstance
	}

	return *r, nil
}

func recipeInit(config *epmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, nil, onGeneralError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, defaultErrors.New("emailpassword recipe has already been initialised. Please check your code for bugs.")
	}
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, defaultErrors.New("initialisation not done. Did you forget to call the init function?")
}

// implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
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
		PathWithoutAPIBasePath: signUpAPI,
		ID:                     constants.SignUpAPI,
		Disabled:               r.APIImpl.SignUpPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: signInAPI,
		ID:                     constants.SignInAPI,
		Disabled:               r.APIImpl.SignInPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: generatePasswordResetTokenAPI,
		ID:                     constants.GeneratePasswordResetTokenAPI,
		Disabled:               r.APIImpl.GeneratePasswordResetTokenPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: passwordResetAPI,
		ID:                     constants.PasswordResetAPI,
		Disabled:               r.APIImpl.PasswordResetPOST == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: signupEmailExistsAPI,
		ID:                     constants.SignupEmailExistsAPI,
		Disabled:               r.APIImpl.EmailExistsGET == nil,
	}}, emailverificationAPIhandled...), nil
}

func (r *Recipe) handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path supertokens.NormalisedURLPath, method string) error {
	options := epmodels.APIOptions{
		Config:                                r.Config,
		OtherHandler:                          theirHandler,
		RecipeID:                              r.RecipeModule.GetRecipeID(),
		RecipeImplementation:                  r.RecipeImpl,
		EmailVerificationRecipeImplementation: r.EmailVerificationRecipe.RecipeImpl,
		Req:                                   req,
		Res:                                   res,
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

func (r *Recipe) getAllCORSHeaders() []string {
	return r.EmailVerificationRecipe.RecipeModule.GetAllCORSHeaders()
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	if defaultErrors.As(err, &errors.FieldError{}) {
		errs := err.(errors.FieldError)
		return true, supertokens.Send200Response(res, map[string]interface{}{
			"status":     "FIELD_ERROR",
			"formFields": errs.Payload,
		})
	}
	return r.EmailVerificationRecipe.RecipeModule.HandleError(err, req, res)
}

func (r *Recipe) getEmailForUserId(userID string, userContext supertokens.UserContext) (string, error) {
	userInfo, err := (*r.RecipeImpl.GetUserByID)(userID, userContext)
	if err != nil {
		return "", err
	}
	if userInfo == nil {
		return "", defaultErrors.New("unknown User ID provided")
	}
	return userInfo.Email, nil
}

func ResetForTest() {
	singletonInstance = nil
}
