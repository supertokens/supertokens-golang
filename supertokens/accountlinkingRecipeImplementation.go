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

		// TODO: add helper functions to the user object.

		if err != nil {
			return UserPaginationResult{}, err
		}

		return result, nil
	}

	// TODO:...
	return AccountLinkingRecipeInterface{
		GetUsersWithSearchParams: &getUsers,
	}
}
