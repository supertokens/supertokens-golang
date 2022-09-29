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

package dashboard

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/api"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "dashboard"

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       dashboardmodels.TypeNormalisedInput
	RecipeImpl   dashboardmodels.RecipeInterface
	APIImpl      dashboardmodels.APIInterface
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config dashboardmodels.TypeInput, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}
	verifiedConfig := validateAndNormaliseUserInput(appInfo, config)
	r.Config = verifiedConfig

	recipeImplementation := makeRecipeImplementation()
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, r.getAPIIdIfCanHandleRequest, r.handleError, onSuperTokensAPIError)
	r.RecipeModule = recipeModuleInstance

	return *r, nil
}

func recipeInit(config dashboardmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("Dashboard recipe has already been initialised. Please check your code for bugs.")
	}
}

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	return []supertokens.APIHandled{}, nil
}

func (r *Recipe) getAPIIdIfCanHandleRequest(path supertokens.NormalisedURLPath, method string) (*string, error) {
	ok, err := isApiPath(path, r.RecipeModule.GetAppInfo())
	if err != nil {
		return nil, err
	}
	if ok {
		return getApiIdIfMatched(path, method)
	}

	dashboardAPIPath, err := supertokens.NewNormalisedURLPath(dashboardAPI)
	if err != nil {
		return nil, err
	}
	dashboardBundlePath := r.RecipeModule.GetAppInfo().APIBasePath.AppendPath(dashboardAPIPath)

	if path.StartsWith(dashboardBundlePath) {
		val := dashboardAPI
		return &val, nil
	}

	return nil, nil
}

func (r *Recipe) handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string) error {
	options := dashboardmodels.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		AppInfo:              r.RecipeModule.GetAppInfo(),
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirHandler,
	}
	if id == dashboardAPI {
		return api.Dashboard(r.APIImpl, options)
	} else if id == validateKeyAPI {
		return api.ValidateKey(r.APIImpl, options)
	}

	// Do API key validation for the remaining APIs
	if id == usersListGetAPI || id == usersCountAPI {
		userContext := supertokens.MakeDefaultUserContextFromAPI(req)
		return apiKeyProtector(r.APIImpl, options, userContext, func() error {
			if id == usersListGetAPI {
				return api.UsersGet(r.APIImpl, options)
			} else if id == usersCountAPI {
				return api.UsersCountGet(r.APIImpl, options)
			}
			return errors.New("should never come here")
		})
	}
	return errors.New("should never come here")
}

func (r *Recipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	return false, nil
}

func ResetForTest() {
	singletonInstance = nil
}
