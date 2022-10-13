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

package recipeimplementation

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeThirdPartyRecipeImplementation(recipeImplementation tplmodels.RecipeInterface) tpmodels.RecipeInterface {

	getUserByThirdPartyInfo := func(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
		user, err := (*recipeImplementation.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, userContext)
		if err != nil {
			return nil, err
		}
		if user == nil || user.ThirdParty == nil {
			return nil, nil
		}
		return &tpmodels.User{
			ID:         user.ID,
			Email:      *user.Email,
			TimeJoined: user.TimeJoined,
			ThirdParty: *user.ThirdParty,
		}, nil
	}

	signInUp := func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens tpmodels.TypeOAuthTokens, rawUserInfoFromProvider map[string]interface{}, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) {
		result, err := (*recipeImplementation.ThirdPartySignInUp)(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, userContext)
		if err != nil {
			return tpmodels.SignInUpResponse{}, err
		}

		return tpmodels.SignInUpResponse{
			OK: &struct {
				CreatedNewUser          bool
				User                    tpmodels.User
				OAuthTokens             tpmodels.TypeOAuthTokens
				RawUserInfoFromProvider map[string]interface{}
			}{
				CreatedNewUser: result.OK.CreatedNewUser,
				User: tpmodels.User{
					ID:         result.OK.User.ID,
					Email:      *result.OK.User.Email,
					TimeJoined: result.OK.User.TimeJoined,
					ThirdParty: *result.OK.User.ThirdParty,
				},
				OAuthTokens:             result.OK.OAuthTokens,
				RawUserInfoFromProvider: result.OK.RawUserInfoFromProvider,
			},
		}, nil
	}

	createUser := func(thirdPartyID string, thirdPartyUserID string, email string, userContext supertokens.UserContext) (tpmodels.CreateUserResponse, error) {
		result, err := (*recipeImplementation.ThirdPartyCreateUser)(thirdPartyID, thirdPartyUserID, email, userContext)
		if err != nil {
			return tpmodels.CreateUserResponse{}, err
		}

		return tpmodels.CreateUserResponse{
			OK: &struct {
				CreatedNewUser bool
				User           tpmodels.User
			}{
				CreatedNewUser: result.OK.CreatedNewUser,
				User: tpmodels.User{
					ID:         result.OK.User.ID,
					Email:      *result.OK.User.Email,
					TimeJoined: result.OK.User.TimeJoined,
					ThirdParty: *result.OK.User.ThirdParty,
				},
			},
		}, nil
	}

	getUserByID := func(userID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
		user, err := (*recipeImplementation.GetUserByID)(userID, userContext)
		if err != nil {
			return nil, err
		}
		if user == nil || user.ThirdParty == nil {
			return nil, nil
		}
		return &tpmodels.User{
			ID:         user.ID,
			Email:      *user.Email,
			TimeJoined: user.TimeJoined,
			ThirdParty: *user.ThirdParty,
		}, nil
	}

	getUserByEmail := func(email string, userContext supertokens.UserContext) ([]tpmodels.User, error) {
		users, err := (*recipeImplementation.GetUsersByEmail)(email, userContext)
		if err != nil {
			return nil, err
		}

		finalResult := []tpmodels.User{}

		for _, tpepUser := range users {
			if tpepUser.ThirdParty != nil {
				finalResult = append(finalResult, tpmodels.User{
					ID:         tpepUser.ID,
					TimeJoined: tpepUser.TimeJoined,
					Email:      *tpepUser.Email,
					ThirdParty: *tpepUser.ThirdParty,
				})
			}
		}
		return finalResult, nil
	}

	return tpmodels.RecipeInterface{
		GetUserByID:             &getUserByID,
		GetUsersByEmail:         &getUserByEmail,
		GetUserByThirdPartyInfo: &getUserByThirdPartyInfo,
		SignInUp:                &signInUp,
		CreateUser:              &createUser,
	}
}
