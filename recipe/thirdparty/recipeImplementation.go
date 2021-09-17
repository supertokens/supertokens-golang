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

package thirdparty

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier) tpmodels.RecipeInterface {
	return tpmodels.RecipeInterface{
		SignInUp: func(thirdPartyID, thirdPartyUserID string, email tpmodels.EmailStruct) (tpmodels.SignInUpResponse, error) {
			response, err := querier.SendPostRequest("/recipe/signinup", map[string]interface{}{
				"thirdPartyId":     thirdPartyID,
				"thirdPartyUserId": thirdPartyUserID,
				"email":            email,
			})
			if err != nil {
				return tpmodels.SignInUpResponse{}, err
			}
			user, err := parseUser(response["user"])
			if err != nil {
				return tpmodels.SignInUpResponse{}, err
			}
			return tpmodels.SignInUpResponse{
				OK: &struct {
					CreatedNewUser bool
					User           tpmodels.User
				}{
					CreatedNewUser: response["createdNewUser"].(bool),
					User:           *user,
				},
			}, nil
		},

		GetUserByID: func(userID string) (*tpmodels.User, error) {
			response, err := querier.SendGetRequest("/recipe/user", map[string]string{
				"userId": userID,
			})
			if err != nil {
				return nil, err
			}
			if response["status"] == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return nil, err
				}
				return user, nil
			}
			return nil, nil
		},

		GetUserByThirdPartyInfo: func(thirdPartyID, thirdPartyUserID string) (*tpmodels.User, error) {
			response, err := querier.SendGetRequest("/recipe/user", map[string]string{
				"thirdPartyId":     thirdPartyID,
				"thirdPartyUserId": thirdPartyUserID,
			})
			if err != nil {
				return nil, err
			}
			if response["status"] == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return nil, err
				}
				return user, nil
			}
			return nil, nil
		},

		GetUsersByEmail: func(email string) ([]tpmodels.User, error) {
			response, err := querier.SendGetRequest("/recipe/users/by-email", map[string]string{
				"email": email,
			})
			if err != nil {
				return []tpmodels.User{}, err
			}
			users, err := parseUsers(response["users"])
			if err != nil {
				return []tpmodels.User{}, err
			}
			return users, nil
		},
	}
}
