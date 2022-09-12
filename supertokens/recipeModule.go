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

package supertokens

import "net/http"

type RecipeModule struct {
	recipeID                   string
	appInfo                    NormalisedAppinfo
	HandleAPIRequest           func(ID string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path NormalisedURLPath, method string) error
	GetAllCORSHeaders          func() []string
	GetAPIsHandled             func() ([]APIHandled, error)
	GetAPIIdIfCanHandleRequest func(path NormalisedURLPath, method string) (*string, error)
	HandleError                func(err error, req *http.Request, res http.ResponseWriter) (bool, error)
	OnSuperTokensAPIError      func(err error, req *http.Request, res http.ResponseWriter)
}

func MakeRecipeModule(
	recipeId string,
	appInfo NormalisedAppinfo,
	handleAPIRequest func(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path NormalisedURLPath, method string) error,
	getAllCORSHeaders func() []string,
	getAPIsHandled func() ([]APIHandled, error),
	getAPIIdIfCanHandleRequest func(path NormalisedURLPath, method string) (*string, error),
	handleError func(err error, req *http.Request, res http.ResponseWriter) (bool, error),
	onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) RecipeModule {
	if handleError == nil {
		// Execution will come here only if there is a bug in the code
		panic("nil passed for handleError in recipe")
	}
	if onSuperTokensAPIError == nil {
		// Execution will come here only if there is a bug in the code
		panic("nil passed for OnSuperTokensAPIError in recipe")
	}
	return RecipeModule{
		recipeID:                   recipeId,
		appInfo:                    appInfo,
		HandleAPIRequest:           handleAPIRequest,
		GetAllCORSHeaders:          getAllCORSHeaders,
		GetAPIsHandled:             getAPIsHandled,
		GetAPIIdIfCanHandleRequest: getAPIIdIfCanHandleRequest,
		HandleError:                handleError,
		OnSuperTokensAPIError:      onSuperTokensAPIError,
	}
}

func (r RecipeModule) GetRecipeID() string {
	return r.recipeID
}

func (r RecipeModule) GetAppInfo() NormalisedAppinfo {
	return r.appInfo
}

func (r *RecipeModule) ReturnAPIIdIfCanHandleRequest(path NormalisedURLPath, method string) (*string, error) {
	apisHandled, err := r.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	for _, APIshandled := range apisHandled {
		pathAppend := r.appInfo.APIBasePath.AppendPath(APIshandled.PathWithoutAPIBasePath)
		if !APIshandled.Disabled && APIshandled.Method == method && pathAppend.Equals(path) {
			return &APIshandled.ID, nil
		}
	}

	if r.GetAPIIdIfCanHandleRequest != nil {
		return r.GetAPIIdIfCanHandleRequest(path, method)
	}

	return nil, nil
}
