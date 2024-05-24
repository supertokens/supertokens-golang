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

import (
	"fmt"
	"net/http"
	"regexp"
)

type RecipeModule struct {
	recipeID                      string
	appInfo                       NormalisedAppinfo
	HandleAPIRequest              func(ID string, tenantId string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path NormalisedURLPath, method string, userContext UserContext) error
	GetAllCORSHeaders             func() []string
	GetAPIsHandled                func() ([]APIHandled, error)
	ReturnAPIIdIfCanHandleRequest func(path NormalisedURLPath, method string, userContext UserContext) (*string, string, error)
	HandleError                   func(err error, req *http.Request, res http.ResponseWriter, userContext UserContext) (bool, error)
	OnSuperTokensAPIError         func(err error, req *http.Request, res http.ResponseWriter)
	ResetForTest                  func()
}

func MakeRecipeModule(
	recipeId string,
	appInfo NormalisedAppinfo,
	handleAPIRequest func(id string, tenantId string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path NormalisedURLPath, method string, userContext UserContext) error,
	getAllCORSHeaders func() []string,
	getAPIsHandled func() ([]APIHandled, error),
	returnAPIIdIfCanHandleRequest func(path NormalisedURLPath, method string, userContext UserContext) (*string, string, error),
	handleError func(err error, req *http.Request, res http.ResponseWriter, userContext UserContext) (bool, error),
	onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) RecipeModule {
	if handleError == nil {
		// Execution will come here only if there is a bug in the code
		panic("nil passed for handleError in recipe")
	}
	if onSuperTokensAPIError == nil {
		// Execution will come here only if there is a bug in the code
		panic("nil passed for OnSuperTokensAPIError in recipe")
	}

	if returnAPIIdIfCanHandleRequest == nil {
		returnAPIIdIfCanHandleRequest = func(path NormalisedURLPath, method string, userContext UserContext) (*string, string, error) {
			apisHandled, err := getAPIsHandled()
			if err != nil {
				return nil, "", err
			}

			basePathStr := appInfo.APIBasePath.GetAsStringDangerous()
			pathStr := path.GetAsStringDangerous()
			regexStr := fmt.Sprintf(`^%s(?:/([a-zA-Z0-9-]+))?(/.*)$`, regexp.QuoteMeta(basePathStr))
			regex := regexp.MustCompile(regexStr)
			var remainingPath *NormalisedURLPath
			var tenantId string = DefaultTenantId

			if regex.MatchString(pathStr) {
				matches := regex.FindStringSubmatch(pathStr)
				tenantId = matches[1]
				remainingPath2, err := NewNormalisedURLPath(matches[2])
				if err != nil {
					return nil, "", err
				}
				remainingPath = &remainingPath2
			}

			for _, APIshandled := range apisHandled {
				pathAppend := appInfo.APIBasePath.AppendPath(APIshandled.PathWithoutAPIBasePath)
				if !APIshandled.Disabled && APIshandled.Method == method {
					if pathAppend.Equals(path) {
						return &APIshandled.ID, DefaultTenantId, nil
					} else if remainingPath != nil && appInfo.APIBasePath.AppendPath(APIshandled.PathWithoutAPIBasePath).Equals(appInfo.APIBasePath.AppendPath(*remainingPath)) {
						return &APIshandled.ID, tenantId, nil
					}
				}
			}

			return nil, "", nil
		}
	}
	return RecipeModule{
		recipeID:                      recipeId,
		appInfo:                       appInfo,
		HandleAPIRequest:              handleAPIRequest,
		GetAllCORSHeaders:             getAllCORSHeaders,
		GetAPIsHandled:                getAPIsHandled,
		ReturnAPIIdIfCanHandleRequest: returnAPIIdIfCanHandleRequest,
		HandleError:                   handleError,
		OnSuperTokensAPIError:         onSuperTokensAPIError,
	}
}

func (r RecipeModule) GetRecipeID() string {
	return r.recipeID
}

func (r RecipeModule) GetAppInfo() NormalisedAppinfo {
	return r.appInfo
}
