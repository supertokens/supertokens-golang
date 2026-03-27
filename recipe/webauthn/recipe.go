/* Copyright (c) 2025, VRAI Labs and/or its affiliates. All rights reserved.
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

package webauthn

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/webauthn/api"
	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "webauthn"

type Recipe struct {
	RecipeModule  supertokens.RecipeModule
	Config        webauthnmodels.TypeNormalisedInput
	RecipeImpl    webauthnmodels.RecipeInterface
	APIImpl       webauthnmodels.APIInterface
	EmailDelivery emaildelivery.Ingredient
}

var singletonInstance *Recipe

func MakeRecipe(
	recipeId string,
	appInfo supertokens.NormalisedAppinfo,
	config *webauthnmodels.TypeInput,
	emailDeliveryIngredient *emaildelivery.Ingredient,
	onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter),
) (Recipe, error) {
	r := &Recipe{}
	verifiedConfig := validateAndNormaliseUserInput(appInfo, config)
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	if emailDeliveryIngredient != nil {
		r.EmailDelivery = *emailDeliveryIngredient
	} else {
		r.EmailDelivery = emaildelivery.MakeIngredient(verifiedConfig.GetEmailDeliveryConfig())
	}

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}

	r.RecipeImpl = verifiedConfig.Override.Functions(MakeRecipeImplementation(*querierInstance, r.EmailDelivery))

	r.RecipeModule = supertokens.MakeRecipeModule(
		recipeId, appInfo,
		r.handleAPIRequest,
		r.getAllCORSHeaders,
		r.getAPIsHandled,
		nil,
		r.handleError,
		onSuperTokensAPIError,
	)

	r.RecipeModule.ResetForTest = resetForTest

	return *r, nil
}

func Init(config *webauthnmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, nil, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("webauthn recipe has already been initialised. Please check your code for bugs.")
	}
}

func GetRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, errors.New("initialisation not done. Did you forget to call the Init function?")
}

func GetRecipeInstance() *Recipe {
	return singletonInstance
}

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	registerOptionsNorm, err := supertokens.NewNormalisedURLPath(registerOptionsAPI)
	if err != nil {
		return nil, err
	}
	signInOptionsNorm, err := supertokens.NewNormalisedURLPath(signInOptionsAPI)
	if err != nil {
		return nil, err
	}
	signUpNorm, err := supertokens.NewNormalisedURLPath(signUpAPI)
	if err != nil {
		return nil, err
	}
	signInNorm, err := supertokens.NewNormalisedURLPath(signInAPI)
	if err != nil {
		return nil, err
	}
	generateRecoverAccountTokenNorm, err := supertokens.NewNormalisedURLPath(generateRecoverAccountTokenAPI)
	if err != nil {
		return nil, err
	}
	recoverAccountNorm, err := supertokens.NewNormalisedURLPath(recoverAccountAPI)
	if err != nil {
		return nil, err
	}
	emailExistsNorm, err := supertokens.NewNormalisedURLPath(doesEmailExistAPI)
	if err != nil {
		return nil, err
	}
	registerCredentialNorm, err := supertokens.NewNormalisedURLPath(registerCredentialAPI)
	if err != nil {
		return nil, err
	}
	listCredentialsNorm, err := supertokens.NewNormalisedURLPath(listCredentialsAPI)
	if err != nil {
		return nil, err
	}
	removeCredentialNorm, err := supertokens.NewNormalisedURLPath(removeCredentialAPI)
	if err != nil {
		return nil, err
	}

	return []supertokens.APIHandled{
		{
			Method:                 http.MethodPost,
			PathWithoutAPIBasePath: registerOptionsNorm,
			ID:                     registerOptionsAPI,
			Disabled:               r.APIImpl.RegisterOptionsPOST == nil,
		},
		{
			Method:                 http.MethodPost,
			PathWithoutAPIBasePath: signInOptionsNorm,
			ID:                     signInOptionsAPI,
			Disabled:               r.APIImpl.SignInOptionsPOST == nil,
		},
		{
			Method:                 http.MethodPost,
			PathWithoutAPIBasePath: signUpNorm,
			ID:                     signUpAPI,
			Disabled:               r.APIImpl.SignUpPOST == nil,
		},
		{
			Method:                 http.MethodPost,
			PathWithoutAPIBasePath: signInNorm,
			ID:                     signInAPI,
			Disabled:               r.APIImpl.SignInPOST == nil,
		},
		{
			Method:                 http.MethodPost,
			PathWithoutAPIBasePath: generateRecoverAccountTokenNorm,
			ID:                     generateRecoverAccountTokenAPI,
			Disabled:               r.APIImpl.GenerateRecoverAccountTokenPOST == nil,
		},
		{
			Method:                 http.MethodPost,
			PathWithoutAPIBasePath: recoverAccountNorm,
			ID:                     recoverAccountAPI,
			Disabled:               r.APIImpl.RecoverAccountPOST == nil,
		},
		{
			Method:                 http.MethodGet,
			PathWithoutAPIBasePath: emailExistsNorm,
			ID:                     doesEmailExistAPI,
			Disabled:               r.APIImpl.EmailExistsGET == nil,
		},
		{
			Method:                 http.MethodPost,
			PathWithoutAPIBasePath: registerCredentialNorm,
			ID:                     registerCredentialAPI,
			Disabled:               r.APIImpl.RegisterCredentialPOST == nil,
		},
		{
			Method:                 http.MethodGet,
			PathWithoutAPIBasePath: listCredentialsNorm,
			ID:                     listCredentialsAPI,
			Disabled:               r.APIImpl.ListCredentialsGET == nil,
		},
		{
			Method:                 http.MethodPost,
			PathWithoutAPIBasePath: removeCredentialNorm,
			ID:                     removeCredentialAPI,
			Disabled:               r.APIImpl.RemoveCredentialPOST == nil,
		},
	}, nil
}

func (r *Recipe) handleAPIRequest(id string, tenantId string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string, userContext supertokens.UserContext) error {
	options := webauthnmodels.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		AppInfo:              r.RecipeModule.GetAppInfo(),
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirHandler,
		EmailDelivery:        r.EmailDelivery,
	}
	switch id {
	case registerOptionsAPI:
		return api.RegisterOptions(r.APIImpl, tenantId, options, userContext)
	case signInOptionsAPI:
		return api.SignInOptions(r.APIImpl, tenantId, options, userContext)
	case signUpAPI:
		return api.SignUp(r.APIImpl, tenantId, options, userContext)
	case signInAPI:
		return api.SignIn(r.APIImpl, tenantId, options, userContext)
	case generateRecoverAccountTokenAPI:
		return api.GenerateRecoverAccountToken(r.APIImpl, tenantId, options, userContext)
	case recoverAccountAPI:
		return api.RecoverAccount(r.APIImpl, tenantId, options, userContext)
	case doesEmailExistAPI:
		return api.EmailExists(r.APIImpl, tenantId, options, userContext)
	case registerCredentialAPI:
		return api.RegisterCredential(r.APIImpl, tenantId, options, userContext)
	case listCredentialsAPI:
		return api.ListCredentials(r.APIImpl, tenantId, options, userContext)
	case removeCredentialAPI:
		return api.RemoveCredential(r.APIImpl, tenantId, options, userContext)
	}
	return errors.New("should never come here: unknown API id " + id)
}

func (r *Recipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *Recipe) handleError(err error, _ *http.Request, _ http.ResponseWriter, _ supertokens.UserContext) (bool, error) {
	return false, nil
}

func resetForTest() {
	singletonInstance = nil
}
