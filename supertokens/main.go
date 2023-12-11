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
	"net/http"
)

func Init(config TypeInput) error {
	err := supertokensInit(config)
	if err != nil {
		return err
	}
	err = runPostInitCallbacks()
	if err != nil {
		return err
	}
	return nil
}

func Middleware(theirHandler http.Handler) http.Handler {
	instance, err := GetInstanceOrThrowError()
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		})
	}
	return instance.middleware(theirHandler)
}

func ErrorHandler(err error, req *http.Request, res http.ResponseWriter, userContext ...UserContext) error {
	instance, instanceErr := GetInstanceOrThrowError()
	if instanceErr != nil {
		return instanceErr
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return instance.errorHandler(err, req, res, userContext[0])
}

func GetAllCORSHeaders() []string {
	instance, err := GetInstanceOrThrowError()
	if err != nil {
		panic("Please call supertokens.Init before using the GetAllCORSHeaders function")
	}
	return instance.getAllCORSHeaders()
}

func GetUserCount(includeRecipeIds *[]string, tenantId *string) (float64, error) {
	var includeAllTenants *bool
	if tenantId == nil {
		defaultTenantId := DefaultTenantId
		tenantId = &defaultTenantId
		True := true
		includeAllTenants = &True
	}
	return getUserCount(includeRecipeIds, *tenantId, includeAllTenants)
}

func GetUsersOldestFirst(tenantId string, paginationToken *string, limit *int, includeRecipeIds *[]string, query map[string]string, userContext ...UserContext) (UserPaginationResult, error) {
	accountLinkingInstance, err := getAccountLinkingRecipeInstanceOrThrowError()
	if err != nil {
		return UserPaginationResult{}, err
	}

	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}

	return (*accountLinkingInstance.RecipeImpl.GetUsersWithSearchParams)(tenantId, "ASC", paginationToken, limit, includeRecipeIds, query, userContext[0])
}

func GetUsersNewestFirst(tenantId string, paginationToken *string, limit *int, includeRecipeIds *[]string, query map[string]string, userContext ...UserContext) (UserPaginationResult, error) {
	accountLinkingInstance, err := getAccountLinkingRecipeInstanceOrThrowError()
	if err != nil {
		return UserPaginationResult{}, err
	}

	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}

	return (*accountLinkingInstance.RecipeImpl.GetUsersWithSearchParams)(tenantId, "DESC", paginationToken, limit, includeRecipeIds, query, userContext[0])
}

func DeleteUser(userId string) error {
	return deleteUser(userId)
}

func GetRequestFromUserContext(userContext UserContext) *http.Request {
	return getRequestFromUserContext(userContext)
}
