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

package session

import (
	defaultErrors "errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/api"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "session"

var r *sessmodels.SessionRecipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *sessmodels.TypeInput, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (sessmodels.SessionRecipe, error) {
	r = &sessmodels.SessionRecipe{}
	verifiedConfig, configError := validateAndNormaliseUserInput(appInfo, config)
	if configError != nil {
		return sessmodels.SessionRecipe{}, configError
	}
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return sessmodels.SessionRecipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance, verifiedConfig)
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, handleAPIRequest, getAllCORSHeaders, getAPIsHandled, handleError, onGeneralError)
	r.RecipeModule = recipeModuleInstance

	return sessmodels.SessionRecipe{
		RecipeModule: recipeModuleInstance,
		Config:       verifiedConfig,
		RecipeImpl:   verifiedConfig.Override.Functions(recipeImplementation),
		APIImpl:      verifiedConfig.Override.APIs(api.MakeAPIImplementation()),
	}, nil
}

func getRecipeInstanceOrThrowError() (*sessmodels.SessionRecipe, error) {
	if r != nil {
		return r, nil
	}
	return nil, defaultErrors.New("Initialisation not done. Did you forget to call the init function?")
}

func recipeInit(config *sessmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if r == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, onGeneralError)
			if err != nil {
				return nil, err
			}
			r = &recipe
			return &r.RecipeModule, nil
		}
		return nil, defaultErrors.New("Session recipe has already been initialised. Please check your code for bugs.")
	}
}

// Implement RecipeModule

func getAPIsHandled() ([]supertokens.APIHandled, error) {
	refreshAPIPathNormalised, err := supertokens.NewNormalisedURLPath(refreshAPIPath)
	if err != nil {
		return nil, err
	}
	signoutAPIPathNormalised, err := supertokens.NewNormalisedURLPath(signoutAPIPath)
	if err != nil {
		return nil, err
	}
	return []supertokens.APIHandled{{
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: refreshAPIPathNormalised,
		ID:                     refreshAPIPath,
		Disabled:               r.APIImpl.RefreshPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: signoutAPIPathNormalised,
		ID:                     signoutAPIPath,
		Disabled:               r.APIImpl.SignOutPOST == nil,
	}}, nil
}

func handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirhandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string) error {
	options := sessmodels.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirhandler,
	}
	if id == refreshAPIPath {
		return api.HandleRefreshAPI(r.APIImpl, options)
	} else {
		return api.SignOutAPI(r.APIImpl, options)
	}
}

func getAllCORSHeaders() []string {
	return getCORSAllowedHeaders()
}

func handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	if defaultErrors.As(err, &errors.UnauthorizedError{}) {
		return true, r.Config.ErrorHandlers.OnUnauthorised(err.Error(), req, res)
	} else if defaultErrors.As(err, &errors.TryRefreshTokenError{}) {
		return true, r.Config.ErrorHandlers.OnTryRefreshToken(err.Error(), req, res)
	} else if defaultErrors.As(err, &errors.TokenTheftDetectedError{}) {
		errs := err.(errors.TokenTheftDetectedError)
		return true, r.Config.ErrorHandlers.OnTokenTheftDetected(errs.Payload.SessionHandle, errs.Payload.UserID, req, res)
	}
	return false, nil
}

func ResetForTest() {
	r = nil
}
