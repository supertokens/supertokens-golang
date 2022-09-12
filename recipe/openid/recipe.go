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

package openid

import (
	defaultErrors "errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/jwt"
	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/recipe/openid/api"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       openidmodels.TypeNormalisedInput
	RecipeImpl   openidmodels.RecipeInterface
	JwtRecipe    jwt.Recipe
	APIImpl      openidmodels.APIInterface
}

const RECIPE_ID = "openid"

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *openidmodels.TypeInput, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}

	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, nil, r.handleError, onSuperTokensAPIError)

	verifiedConfig, configError := validateAndNormaliseUserInput(appInfo, config)
	if configError != nil {
		return Recipe{}, configError
	}
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	jwtRecipe, err := jwt.MakeRecipe(recipeId, appInfo, &jwtmodels.TypeInput{
		JwtValiditySeconds: verifiedConfig.JwtValiditySeconds,
		Override:           verifiedConfig.Override.JwtFeature,
	}, onSuperTokensAPIError)
	if err != nil {
		return Recipe{}, err
	}
	r.RecipeImpl = verifiedConfig.Override.Functions(makeRecipeImplementation(verifiedConfig, jwtRecipe.RecipeImpl))
	r.JwtRecipe = jwtRecipe

	return *r, nil
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, defaultErrors.New("Initialisation not done. Did you forget to call the init function?")
}

func recipeInit(config *openidmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, defaultErrors.New("OpenID recipe has already been initialised. Please check your code for bugs.")
	}
}

// Implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	normalisedPath, err := supertokens.NewNormalisedURLPath(GetDiscoveryConfigUrl)
	if err != nil {
		return nil, err
	}
	resp := []supertokens.APIHandled{{
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: normalisedPath,
		ID:                     GetDiscoveryConfigUrl,
		Disabled:               r.APIImpl.GetOpenIdDiscoveryConfigurationGET == nil,
	}}

	jwtAPIs, err := r.JwtRecipe.RecipeModule.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	resp = append(resp, jwtAPIs...)

	return resp, nil
}

func (r *Recipe) handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirhandler http.HandlerFunc, path supertokens.NormalisedURLPath, method string) error {
	options := openidmodels.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirhandler,
	}
	if id == GetDiscoveryConfigUrl {
		return api.GetOpenIdDiscoveryConfiguration(r.APIImpl, options)
	} else {
		return r.JwtRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirhandler, path, method)
	}
}

func (r *Recipe) getAllCORSHeaders() []string {
	return r.JwtRecipe.RecipeModule.GetAllCORSHeaders()
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	return r.JwtRecipe.RecipeModule.HandleError(err, req, res)
}

func ResetForTest() {
	singletonInstance = nil
}
