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

package emailverification

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/api"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "emailverification"

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       evmodels.TypeNormalisedInput
	RecipeImpl   evmodels.RecipeInterface
	APIImpl      evmodels.APIInterface
}

var r = &Recipe{}

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *evmodels.TypeInput, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	verifiedConfig := validateAndNormaliseUserInput(appInfo, *config)
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance)
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, handleAPIRequest, getAllCORSHeaders, getAPIsHandled, handleError, onGeneralError)
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

func recipeInit(config *evmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if r == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, onGeneralError)
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
		PathWithoutAPIBasePath: generateEmailVerifyTokenAPINormalised,
		ID:                     generateEmailVerifyTokenAPI,
		Disabled:               r.APIImpl.GenerateEmailVerifyTokenPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: emailVerifyAPINormalised,
		ID:                     emailVerifyAPI,
		Disabled:               r.APIImpl.VerifyEmailPOST == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: emailVerifyAPINormalised,
		ID:                     emailVerifyAPI,
		Disabled:               r.APIImpl.IsEmailVerifiedGET == nil,
	}}, nil
}

func handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string) error {
	options := evmodels.APIOptions{
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

func handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	return false, nil
}
