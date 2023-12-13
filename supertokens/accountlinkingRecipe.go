/* Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package supertokens

import (
	"errors"
	"net/http"
)

const RECIPE_ID = "accountlinking"

type AccountLinkingRecipe struct {
	RecipeModule RecipeModule
	Config       AccountLinkingTypeNormalisedInput
	RecipeImpl   AccountLinkingRecipeInterface
}

var singletonInstance *AccountLinkingRecipe

func makeAccountLinkingRecipe(recipeId string, appInfo NormalisedAppinfo, config *AccountLinkingTypeInput, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (AccountLinkingRecipe, error) {
	r := &AccountLinkingRecipe{}
	verifiedConfig := validateAndNormaliseAccountLinkingUserInput(appInfo, config)
	r.Config = verifiedConfig

	querierInstance, err := GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return AccountLinkingRecipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance, verifiedConfig)
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	recipeModuleInstance := MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, nil, r.handleError, onSuperTokensAPIError)
	r.RecipeModule = recipeModuleInstance

	return *r, nil
}

func getAccountLinkingRecipeInstanceOrThrowError() (*AccountLinkingRecipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the init function?")
}

func accountLinkingRecipeInit(config *AccountLinkingTypeInput) Recipe {
	return func(appInfo NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := makeAccountLinkingRecipe(RECIPE_ID, appInfo, config, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe

			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("Account linking recipe has already been initialised. Please check your code for bugs.")
	}
}

// implement RecipeModule

func (r *AccountLinkingRecipe) getAPIsHandled() ([]APIHandled, error) {
	return []APIHandled{}, nil
}

func (r *AccountLinkingRecipe) handleAPIRequest(id string, tenantId string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ NormalisedURLPath, _ string, userContext UserContext) error {
	return errors.New("should never come here")
}

func (r *AccountLinkingRecipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *AccountLinkingRecipe) handleError(err error, req *http.Request, res http.ResponseWriter, userContext UserContext) (bool, error) {
	return false, nil
}
