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
		panic("Please call supertokens.Init function before using the Middleware")
	}
	return instance.middleware(theirHandler)
}

func ErrorHandler(err error, req *http.Request, res http.ResponseWriter) error {
	instance, instanceErr := GetInstanceOrThrowError()
	if instanceErr != nil {
		return instanceErr
	}
	return instance.errorHandler(err, req, res)
}

func GetAllCORSHeaders() []string {
	instance, err := GetInstanceOrThrowError()
	if err != nil {
		panic("Please call supertokens.Init before using the GetAllCORSHeaders function")
	}
	return instance.getAllCORSHeaders()
}

func GetUserCount(includeRecipeIds *[]string) (float64, error) {
	return getUserCount(includeRecipeIds)
}

func GetUsersOldestFirst(paginationToken *string, limit *int, includeRecipeIds *[]string) (UserPaginationResult, error) {
	return GetUsersWithSearchParams("ASC", paginationToken, limit, includeRecipeIds, nil)
}

func GetUsersNewestFirst(paginationToken *string, limit *int, includeRecipeIds *[]string) (UserPaginationResult, error) {
	return GetUsersWithSearchParams("DESC", paginationToken, limit, includeRecipeIds, nil)
}

func DeleteUser(userId string) error {
	return deleteUser(userId)
}
