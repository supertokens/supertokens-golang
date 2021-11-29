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

	"github.com/supertokens/supertokens-golang/recipe/jwt"
	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/api"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessionwithjwt"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       sessmodels.TypeNormalisedInput
	RecipeImpl   sessmodels.RecipeInterface
	JwtRecipe    *jwt.Recipe
	APIImpl      sessmodels.APIInterface
}

const RECIPE_ID = "session"

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *sessmodels.TypeInput, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}

	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, r.handleError, onGeneralError)

	verifiedConfig, configError := validateAndNormaliseUserInput(appInfo, config)
	if configError != nil {
		return Recipe{}, configError
	}
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance, verifiedConfig)

	if verifiedConfig.Jwt.Enable {
		jwtRecipe, err := jwt.MakeRecipe(recipeId, appInfo, &jwtmodels.TypeInput{
			Override: verifiedConfig.Override.JwtFeature,
		}, onGeneralError)
		if err != nil {
			return Recipe{}, err
		}
		r.RecipeImpl = verifiedConfig.Override.Functions(sessionwithjwt.MakeRecipeImplementation(recipeImplementation, jwtRecipe.RecipeImpl, verifiedConfig, appInfo))
		r.JwtRecipe = &jwtRecipe
	} else {
		r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)
	}

	return *r, nil
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, defaultErrors.New("Initialisation not done. Did you forget to call the init function?")
}

func recipeInit(config *sessmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, onGeneralError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, defaultErrors.New("Session recipe has already been initialised. Please check your code for bugs.")
	}
}

// Implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	refreshAPIPathNormalised, err := supertokens.NewNormalisedURLPath(refreshAPIPath)
	if err != nil {
		return nil, err
	}
	signoutAPIPathNormalised, err := supertokens.NewNormalisedURLPath(signoutAPIPath)
	if err != nil {
		return nil, err
	}
	resp := []supertokens.APIHandled{{
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: refreshAPIPathNormalised,
		ID:                     refreshAPIPath,
		Disabled:               r.APIImpl.RefreshPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: signoutAPIPathNormalised,
		ID:                     signoutAPIPath,
		Disabled:               r.APIImpl.SignOutPOST == nil,
	}}

	if r.JwtRecipe != nil {
		jwtAPIs, err := r.JwtRecipe.RecipeModule.GetAPIsHandled()
		if err != nil {
			return nil, err
		}
		resp = append(resp, jwtAPIs...)
	}

	return resp, nil
}

func (r *Recipe) handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirhandler http.HandlerFunc, path supertokens.NormalisedURLPath, method string) error {
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
	} else if id == signoutAPIPath {
		return api.SignOutAPI(r.APIImpl, options)
	} else if r.JwtRecipe != nil {
		return r.JwtRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirhandler, path, method)
	}
	return nil
}

func (r *Recipe) getAllCORSHeaders() []string {
	resp := getCORSAllowedHeaders()
	if r.JwtRecipe != nil {
		resp = append(resp, r.JwtRecipe.RecipeModule.GetAllCORSHeaders()...)
	}
	return resp
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	if defaultErrors.As(err, &errors.UnauthorizedError{}) {
		return true, r.Config.ErrorHandlers.OnUnauthorised(err.Error(), req, res)
	} else if defaultErrors.As(err, &errors.TryRefreshTokenError{}) {
		return true, r.Config.ErrorHandlers.OnTryRefreshToken(err.Error(), req, res)
	} else if defaultErrors.As(err, &errors.TokenTheftDetectedError{}) {
		errs := err.(errors.TokenTheftDetectedError)
		return true, r.Config.ErrorHandlers.OnTokenTheftDetected(errs.Payload.SessionHandle, errs.Payload.UserID, req, res)
	} else if r.JwtRecipe != nil {
		return r.JwtRecipe.RecipeModule.HandleError(err, req, res)
	}
	return false, nil
}

func ResetForTest() {
	singletonInstance = nil
}
