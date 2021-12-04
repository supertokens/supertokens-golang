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

package passwordless

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/passwordless/api"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "passwordless"

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       plessmodels.TypeNormalisedInput
	RecipeImpl   plessmodels.RecipeInterface
	APIImpl      plessmodels.APIInterface
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config plessmodels.TypeInput, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}
	verifiedConfig := validateAndNormaliseUserInput(appInfo, config)
	r.Config = verifiedConfig

	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance)
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, r.handleError, onGeneralError)
	r.RecipeModule = recipeModuleInstance

	return *r, nil
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, errors.New("initialisation not done. Did you forget to call the init function?")
}

func recipeInit(config plessmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, onGeneralError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("passwordless recipe has already been initialised. Please check your code for bugs")
	}
}

// implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	consumeCodeAPINormalised, err := supertokens.NewNormalisedURLPath(consumeCodeAPI)
	if err != nil {
		return nil, err
	}
	createCodeAPINormalised, err := supertokens.NewNormalisedURLPath(createCodeAPI)
	if err != nil {
		return nil, err
	}
	doesEmailExistsAPINormalised, err := supertokens.NewNormalisedURLPath(doesEmailExistAPI)
	if err != nil {
		return nil, err
	}
	doesPhoneNumberExistsAPINormalised, err := supertokens.NewNormalisedURLPath(doesPhoneNumberExistAPI)
	if err != nil {
		return nil, err
	}
	resendCodeAPINormalised, err := supertokens.NewNormalisedURLPath(resendCodeAPI)
	if err != nil {
		return nil, err
	}

	return []supertokens.APIHandled{{
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: consumeCodeAPINormalised,
		ID:                     consumeCodeAPI,
		Disabled:               r.APIImpl.ConsumeCodePOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: createCodeAPINormalised,
		ID:                     createCodeAPI,
		Disabled:               r.APIImpl.CreateCodePOST == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: doesEmailExistsAPINormalised,
		ID:                     doesEmailExistAPI,
		Disabled:               r.APIImpl.EmailExistsGET == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: doesPhoneNumberExistsAPINormalised,
		ID:                     doesPhoneNumberExistAPI,
		Disabled:               r.APIImpl.PhoneNumberExistsGET == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: resendCodeAPINormalised,
		ID:                     resendCodeAPI,
		Disabled:               r.APIImpl.ResendCodePOST == nil,
	}}, nil
}

func (r *Recipe) handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string) error {
	options := plessmodels.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirHandler,
	}
	if id == consumeCodeAPI {
		return api.ConsumeCode(r.APIImpl, options)
	} else if id == createCodeAPI {
		return api.CreateCode(r.APIImpl, options)
	} else if id == doesEmailExistAPI {
		return api.DoesEmailExist(r.APIImpl, options)
	} else if id == doesPhoneNumberExistAPI {
		return api.DoesPhoneNumberExist(r.APIImpl, options)
	} else {
		return api.ResendCode(r.APIImpl, options)
	}
}

func (r *Recipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	return false, nil
}
