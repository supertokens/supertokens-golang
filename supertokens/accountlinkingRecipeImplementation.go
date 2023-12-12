/* Copyright (c) 2023, VRAI Labs and/or its affiliates. All rights reserved.
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
	"strconv"
	"strings"
)

func makeRecipeImplementation(querier Querier) AccountLinkingRecipeInterface {

	getUsers := func(tenantID string, timeJoinedOrder string, paginationToken *string, limit *int, includeRecipeIds *[]string, searchParams map[string]string, userContext UserContext) (UserPaginationResult, error) {
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

		resp, err := querier.SendGetRequest(tenantID+"/users", requestBody, userContext)

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

	getUser := func(userId string, userContext UserContext) (*User, error) {
		requestBody := map[string]string{
			"userId": userId,
		}
		resp, err := querier.SendGetRequest("/user/id", requestBody, userContext)

		if err != nil {
			return nil, err
		}

		if resp["status"].(string) != "OK" {
			return nil, nil
		}

		var result = User{}

		temporaryVariable, err := json.Marshal(resp["user"])
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(temporaryVariable, &result)
		if err != nil {
			return nil, err
		}

		return &result, nil
	}

	canCreatePrimaryUser := func(recipeUserId RecipeUserID, userContext UserContext) (CanCreatePrimaryUserResponse, error) {
		requestBody := map[string]string{
			"recipeUserId": recipeUserId.GetAsString(),
		}
		resp, err := querier.SendGetRequest("/recipe/accountlinking/user/primary/check", requestBody, userContext)

		if err != nil {
			return CanCreatePrimaryUserResponse{}, err
		}

		if resp["status"].(string) == "OK" {
			return CanCreatePrimaryUserResponse{
				OK: &struct{ WasAlreadyAPrimaryUser bool }{
					WasAlreadyAPrimaryUser: resp["wasAlreadyAPrimaryUser"].(bool),
				},
			}, nil
		} else if resp["status"].(string) == "RECIPE_USER_ID_ALREADY_LINKED_WITH_PRIMARY_USER_ID_ERROR" {
			return CanCreatePrimaryUserResponse{
				RecipeUserIdAlreadyLinkedWithPrimaryUserIdError: &struct {
					PrimaryUserId string
					Description   string
				}{
					PrimaryUserId: resp["primaryUserId"].(string),
					Description:   resp["description"].(string),
				},
			}, nil
		} else {
			return CanCreatePrimaryUserResponse{
				AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError: &struct {
					PrimaryUserId string
					Description   string
				}{
					PrimaryUserId: resp["primaryUserId"].(string),
					Description:   resp["description"].(string),
				},
			}, nil
		}
	}

	createPrimaryUser := func(recipeUserId RecipeUserID, userContext UserContext) (CreatePrimaryUserResponse, error) {
		requestBody := map[string]interface{}{
			"recipeUserId": recipeUserId.GetAsString(),
		}
		resp, err := querier.SendPostRequest("/recipe/accountlinking/user/primary", requestBody, userContext)

		if err != nil {
			return CreatePrimaryUserResponse{}, err
		}

		if resp["status"].(string) == "OK" {
			var user = User{}

			temporaryVariable, err := json.Marshal(resp["user"])
			if err != nil {
				return CreatePrimaryUserResponse{}, err
			}

			err = json.Unmarshal(temporaryVariable, &user)
			if err != nil {
				return CreatePrimaryUserResponse{}, err
			}
			return CreatePrimaryUserResponse{
				OK: &struct {
					User                   User
					WasAlreadyAPrimaryUser bool
				}{
					WasAlreadyAPrimaryUser: resp["wasAlreadyAPrimaryUser"].(bool),
					User:                   user,
				},
			}, nil
		} else if resp["status"].(string) == "RECIPE_USER_ID_ALREADY_LINKED_WITH_PRIMARY_USER_ID_ERROR" {
			return CreatePrimaryUserResponse{
				RecipeUserIdAlreadyLinkedWithPrimaryUserIdError: &struct {
					PrimaryUserId string
				}{
					PrimaryUserId: resp["primaryUserId"].(string),
				},
			}, nil
		} else {
			return CreatePrimaryUserResponse{
				AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError: &struct {
					PrimaryUserId string
					Description   string
				}{
					PrimaryUserId: resp["primaryUserId"].(string),
					Description:   resp["description"].(string),
				},
			}, nil
		}
	}

	// TODO:...
	return AccountLinkingRecipeInterface{
		GetUsersWithSearchParams: &getUsers,
		GetUser:                  &getUser,
		CanCreatePrimaryUser:     &canCreatePrimaryUser,
		CreatePrimaryUser:        &createPrimaryUser,
	}
}
