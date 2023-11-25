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

package jwt

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/jwt/api"
	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "jwt"

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       jwtmodels.TypeNormalisedInput
	RecipeImpl   jwtmodels.RecipeInterface
	APIImpl      jwtmodels.APIInterface
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *jwtmodels.TypeInput, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}
	verifiedConfig := validateAndNormaliseUserInput(appInfo, config)
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance, verifiedConfig, appInfo)
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, nil, r.handleError, onSuperTokensAPIError)
	r.RecipeModule = recipeModuleInstance

	return *r, nil
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the init function?")
}

func recipeInit(config *jwtmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("JWT recipe has already been initialised. Please check your code for bugs.")
	}
}

// implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	getJWKSAPINormalised, err := supertokens.NewNormalisedURLPath(GetJWKSAPI)
	if err != nil {
		return nil, err
	}

	return []supertokens.APIHandled{{
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: getJWKSAPINormalised,
		ID:                     GetJWKSAPI,
		Disabled:               r.APIImpl.GetJWKSGET == nil,
	}}, nil
}

func (r *Recipe) handleAPIRequest(id string, tenantId string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string, userContext supertokens.UserContext) error {
	options := jwtmodels.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirHandler,
	}
	return api.GetJWKS(r.APIImpl, options, userContext)
}

func (r *Recipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter, userContext supertokens.UserContext) (bool, error) {
	return false, nil
}

func ResetForTest() {
	singletonInstance = nil
}
