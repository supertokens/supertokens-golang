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
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// This function is required to be here because calling multitenancy recipe from this module causes cyclic dependency
// this function is initialized by the init function in multitenancy recipe
var GetTenantIdFuncFromUsingMultitenancyRecipe func(tenantIdFromFrontend string, userContext UserContext) (string, error)

type superTokens struct {
	AppInfo               NormalisedAppinfo
	SuperTokens           ConnectionInfo
	RecipeModules         []RecipeModule
	OnSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)
	Telemetry             *bool
}

// this will be set to true if this is used in a test app environment
var IsTestFlag = false

var superTokensInstance *superTokens

func supertokensInit(config TypeInput) error {
	if superTokensInstance != nil {
		return nil
	}

	superTokens := &superTokens{}

	superTokens.OnSuperTokensAPIError = defaultOnSuperTokensAPIError
	if config.OnSuperTokensAPIError != nil {
		superTokens.OnSuperTokensAPIError = config.OnSuperTokensAPIError
	}

	DebugEnabled = config.Debug

	LogDebugMessage("Started SuperTokens with debug logging (supertokens.Init called)")

	appInfoJsonString, _ := json.Marshal(config.AppInfo)
	LogDebugMessage("AppInfo: " + string(appInfoJsonString))

	var err error
	superTokens.AppInfo, err = NormaliseInputAppInfoOrThrowError(config.AppInfo)
	if err != nil {
		return err
	}

	if config.Supertokens != nil {
		if len(config.Supertokens.ConnectionURI) != 0 {
			hostList := strings.Split(config.Supertokens.ConnectionURI, ";")
			hosts := []QuerierHost{}
			for _, h := range hostList {
				domain, err := NewNormalisedURLDomain(h)
				if err != nil {
					return err
				}
				basePath, err := NewNormalisedURLPath(h)
				if err != nil {
					return err
				}
				hosts = append(hosts, QuerierHost{
					Domain:   domain,
					BasePath: basePath,
				})
			}
			initQuerier(hosts, config.Supertokens.APIKey)
			superTokens.SuperTokens = *config.Supertokens
		} else {
			return errors.New("please provide 'ConnectionURI' value. If you do not want to provide a connection URI, then set config.Supertokens to nil")
		}
	} else {
		// TODO: Add tests for init without supertokens core.
	}

	if config.RecipeList == nil || len(config.RecipeList) == 0 {
		return errors.New("please provide at least one recipe to the supertokens.init function call")
	}

	multitenancyFound := false

	for _, elem := range config.RecipeList {
		recipeModule, err := elem(superTokens.AppInfo, superTokens.OnSuperTokensAPIError)
		if err != nil {
			return err
		}
		superTokens.RecipeModules = append(superTokens.RecipeModules, *recipeModule)

		if recipeModule.GetRecipeID() == "multitenancy" {
			multitenancyFound = true
		}
	}

	if !multitenancyFound && DefaultMultitenancyRecipe != nil {
		recipeModule, err := DefaultMultitenancyRecipe(superTokens.AppInfo, superTokens.OnSuperTokensAPIError)
		if err != nil {
			return err
		}
		superTokens.RecipeModules = append(superTokens.RecipeModules, *recipeModule)
	}

	superTokens.Telemetry = config.Telemetry
	superTokensInstance = superTokens

	return nil
}

func defaultOnSuperTokensAPIError(err error, req *http.Request, res http.ResponseWriter) {
	http.Error(res, err.Error(), 500)
}

func GetInstanceOrThrowError() (*superTokens, error) {
	if superTokensInstance != nil {
		return superTokensInstance, nil
	}
	return nil, errors.New("initialisation not done. Did you forget to call the SuperTokens.init function?")
}

func (s *superTokens) middleware(theirHandler http.Handler) http.Handler {
	LogDebugMessage("middleware: Started")
	if theirHandler == nil {
		theirHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dw := MakeDoneWriter(w)
		userContext := MakeDefaultUserContextFromAPI(r)
		reqURL, err := NewNormalisedURLPath(r.URL.Path)
		if err != nil {
			err = s.errorHandler(err, r, dw)
			if err != nil && !dw.IsDone() {
				s.OnSuperTokensAPIError(err, r, dw)
			}
			return
		}
		path := s.AppInfo.APIGatewayPath.AppendPath(reqURL)
		method := r.Method

		if !strings.HasPrefix(path.GetAsStringDangerous(), s.AppInfo.APIBasePath.GetAsStringDangerous()) {
			LogDebugMessage("middleware: Not handling because request path did not start with config path. Request path: " + path.GetAsStringDangerous())
			theirHandler.ServeHTTP(dw, r)
			return
		}
		requestRID := getRIDFromRequest(r)
		LogDebugMessage("middleware: requestRID is: " + requestRID)
		if requestRID == "anti-csrf" {
			// See https://github.com/supertokens/supertokens-node/issues/202
			requestRID = ""
		}
		if requestRID != "" {
			var matchedRecipe *RecipeModule
			for _, recipeModule := range s.RecipeModules {
				LogDebugMessage("middleware: Checking recipe ID for match: " + recipeModule.GetRecipeID())
				if recipeModule.GetRecipeID() == requestRID {
					matchedRecipe = &recipeModule
					break
				}
			}
			if matchedRecipe == nil {
				LogDebugMessage("middleware: Not handling because no recipe matched")
				theirHandler.ServeHTTP(dw, r)
				return
			}

			LogDebugMessage("middleware: Matched with recipe ID: " + matchedRecipe.GetRecipeID())

			id, tenantId, err := matchedRecipe.ReturnAPIIdIfCanHandleRequest(path, method, userContext)

			if err != nil {
				err = s.errorHandler(err, r, dw)
				if err != nil && !dw.IsDone() {
					s.OnSuperTokensAPIError(err, r, dw)
				}
				return
			}

			if id == nil {
				LogDebugMessage("middleware: Not handling because recipe doesn't handle request path or method. Request path: " + path.GetAsStringDangerous() + ", request method: " + method)
				theirHandler.ServeHTTP(dw, r)
				return
			}

			LogDebugMessage("middleware: Request being handled by recipe. ID is: " + *id)

			tenantId, err = GetTenantIdFuncFromUsingMultitenancyRecipe(tenantId, userContext)
			if err != nil {
				err = s.errorHandler(err, r, dw)
				if err != nil && !dw.IsDone() {
					s.OnSuperTokensAPIError(err, r, dw)
				}
				return
			}

			apiErr := matchedRecipe.HandleAPIRequest(*id, tenantId, r, dw, theirHandler.ServeHTTP, path, method, userContext)
			if apiErr != nil {
				apiErr = s.errorHandler(apiErr, r, dw)
				if apiErr != nil && !dw.IsDone() {
					s.OnSuperTokensAPIError(apiErr, r, dw)
				}
				return
			}
			LogDebugMessage("middleware: Ended")
		} else {
			for _, recipeModule := range s.RecipeModules {
				id, tenantId, err := recipeModule.ReturnAPIIdIfCanHandleRequest(path, method, userContext)
				LogDebugMessage("middleware: Checking recipe ID for match: " + recipeModule.GetRecipeID())
				if err != nil {
					err = s.errorHandler(err, r, dw)
					if err != nil && !dw.IsDone() {
						s.OnSuperTokensAPIError(err, r, dw)
					}
					return
				}

				if id != nil {
					LogDebugMessage("middleware: Request being handled by recipe. ID is: " + *id)
					err := recipeModule.HandleAPIRequest(*id, tenantId, r, dw, theirHandler.ServeHTTP, path, method, userContext)
					if err != nil {
						err = s.errorHandler(err, r, dw)
						if err != nil && !dw.IsDone() {
							s.OnSuperTokensAPIError(err, r, dw)
						}
					} else {
						LogDebugMessage("middleware: Ended")
					}
					return
				}
			}

			LogDebugMessage("middleware: Not handling because no recipe matched")
			theirHandler.ServeHTTP(dw, r)
		}
	})
}

func (s *superTokens) getAllCORSHeaders() []string {
	headerMap := map[string]bool{HeaderRID: true, HeaderFDI: true}
	for _, recipe := range s.RecipeModules {
		headers := recipe.GetAllCORSHeaders()
		for _, header := range headers {
			headerMap[header] = true
		}
	}
	var headers []string
	for header := range headerMap {
		headers = append(headers, header)
	}
	return headers
}

func (s *superTokens) errorHandler(originalError error, req *http.Request, res http.ResponseWriter) error {
	LogDebugMessage("errorHandler: Started")
	if errors.As(originalError, &BadInputError{}) {
		LogDebugMessage("errorHandler: Sending 400 status code response")
		err := SendNon200ResponseWithMessage(res, originalError.Error(), 400)
		if err != nil {
			// this function can return an error, so we should return
			// the error here. Once returned, either the user will handle
			// the error themselves, or if this function is being called
			// by our middleware, the middleware will call the OnSuperTokensAPIError callback
			return err
		}
		return nil
	}
	for _, recipe := range s.RecipeModules {
		LogDebugMessage("errorHandler: Checking recipe for match: " + recipe.recipeID)
		if recipe.HandleError != nil {
			LogDebugMessage("errorHandler: Matched with recipeId: " + recipe.recipeID)
			handled, err := recipe.HandleError(originalError, req, res)
			if err != nil {
				return err
			}
			if handled {
				return nil
			}
		}
	}
	return originalError
}

type UserPaginationResult struct {
	Users []struct {
		RecipeId string                 `json:"recipeId"`
		User     map[string]interface{} `json:"user"`
	}
	NextPaginationToken *string
}

// TODO: Add tests
func GetUsersWithSearchParams(tenantId string, timeJoinedOrder string, paginationToken *string, limit *int, includeRecipeIds *[]string, searchParams map[string]string) (UserPaginationResult, error) {

	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return UserPaginationResult{}, err
	}

	requestBody := map[string]string{}
	if searchParams != nil {
		requestBody = searchParams
	}
	requestBody["timeJoinedOrder"] = timeJoinedOrder

	if limit != nil {
		requestBody["limit"] = strconv.Itoa(*limit)
	}
	if paginationToken != nil {
		requestBody["paginationToken"] = *paginationToken
	}
	if includeRecipeIds != nil {
		requestBody["includeRecipeIds"] = strings.Join((*includeRecipeIds)[:], ",")
	}

	resp, err := querier.SendGetRequest(tenantId+"/users", requestBody)

	if err != nil {
		return UserPaginationResult{}, err
	}

	temporaryVariable, err := json.Marshal(resp)
	if err != nil {
		return UserPaginationResult{}, err
	}

	var result = UserPaginationResult{}

	err = json.Unmarshal(temporaryVariable, &result)

	if err != nil {
		return UserPaginationResult{}, err
	}

	return result, nil
}

// TODO: Add tests
func getUserCount(includeRecipeIds *[]string, tenantId string, includeAllTenants *bool) (float64, error) {

	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return -1, err
	}

	requestBody := map[string]string{}

	if includeRecipeIds != nil {
		requestBody["includeRecipeIds"] = strings.Join((*includeRecipeIds)[:], ",")
	}
	if includeAllTenants != nil {
		requestBody["includeAllTenants"] = strconv.FormatBool(*includeAllTenants)
	}

	resp, err := querier.SendGetRequest(tenantId+"/users/count", requestBody)

	if err != nil {
		return -1, err
	}

	return resp["count"].(float64), nil
}

func deleteUser(userId string) error {
	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return err
	}

	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		return err
	}

	if MaxVersion(cdiVersion, "2.10") == cdiVersion {
		_, err = querier.SendPostRequest("/user/remove", map[string]interface{}{
			"userId": userId,
		})

		if err != nil {
			return err
		}

		return nil
	} else {
		return errors.New("please upgrade the SuperTokens core to >= 3.7.0")
	}
}

func ResetForTest() {
	ResetQuerierForTest()
	superTokensInstance = nil
}

func IsRunningInTestMode() bool {
	return flag.Lookup("test.v") != nil || IsTestFlag
}

func getRequestFromUserContext(userContext UserContext) *http.Request {
	if userContext == nil {
		return nil
	}

	_userContext := *userContext
	defaultObj, ok := _userContext["_default"]

	if !ok {
		return nil
	}

	emptyMap := map[string]interface{}{}
	if reflect.TypeOf(defaultObj).Kind() != reflect.TypeOf(emptyMap).Kind() {
		return nil
	}

	return defaultObj.(map[string]interface{})["request"].(*http.Request)
}
