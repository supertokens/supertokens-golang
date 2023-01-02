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
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeThirdPartyRecipeImplementation(recipeImplementation tpepmodels.RecipeInterface) tpmodels.RecipeInterface {

	getUserByThirdPartyInfo := func(thirdPartyID string, thirdPartyUserID string, tenantId *string, userContext supertokens.UserContext) (*tpmodels.User, error) {
		user, err := (*recipeImplementation.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, tenantId, userContext)
		if err != nil {
			return nil, err
		}
		if user == nil || user.ThirdParty == nil {
			return nil, nil
		}
		return &tpmodels.User{
			ID:         user.ID,
			Email:      user.Email,
			TimeJoined: user.TimeJoined,
			ThirdParty: *user.ThirdParty,
		}, nil
	}

	signInUp := func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens tpmodels.TypeOAuthTokens, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, tenantId *string, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) {
		result, err := (*recipeImplementation.ThirdPartySignInUp)(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, tenantId, userContext)
		if err != nil {
			return tpmodels.SignInUpResponse{}, err
		}

		return tpmodels.SignInUpResponse{
			OK: &struct {
				CreatedNewUser          bool
				User                    tpmodels.User
				OAuthTokens             map[string]interface{}
				RawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider
			}{
				CreatedNewUser: result.OK.CreatedNewUser,
				User: tpmodels.User{
					ID:         result.OK.User.ID,
					Email:      result.OK.User.Email,
					TimeJoined: result.OK.User.TimeJoined,
					ThirdParty: *result.OK.User.ThirdParty,
				},
				OAuthTokens:             result.OK.OAuthTokens,
				RawUserInfoFromProvider: result.OK.RawUserInfoFromProvider,
			},
		}, nil
	}

	manuallyCreateOrUpdateUser := func(thirdPartyID string, thirdPartyUserID string, email string, tenantId *string, userContext supertokens.UserContext) (tpmodels.ManuallyCreateOrUpdateUserResponse, error) {
		result, err := (*recipeImplementation.ThirdPartyManuallyCreateOrUpdateUser)(thirdPartyID, thirdPartyUserID, email, tenantId, userContext)
		if err != nil {
			return tpmodels.ManuallyCreateOrUpdateUserResponse{}, err
		}
		return tpmodels.ManuallyCreateOrUpdateUserResponse{
			OK: &struct {
				CreatedNewUser bool
				User           tpmodels.User
			}{
				CreatedNewUser: result.OK.CreatedNewUser,
				User: tpmodels.User{
					ID:         result.OK.User.ID,
					Email:      result.OK.User.Email,
					TimeJoined: result.OK.User.TimeJoined,
					ThirdParty: struct {
						ID     string "json:\"id\""
						UserID string "json:\"userId\""
					}{
						ID:     result.OK.User.ThirdParty.ID,
						UserID: result.OK.User.ThirdParty.UserID,
					},
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
			Email:      user.Email,
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
					Email:      tpepUser.Email,
					ThirdParty: *tpepUser.ThirdParty,
				})
			}
		}
		return finalResult, nil
	}

	getProvider := func(thirdPartyID string, tenantId *string, userContext supertokens.UserContext) (tpmodels.GetProviderResponse, error) {
		return (*recipeImplementation.ThirdPartyGetProvider)(thirdPartyID, tenantId, userContext)
	}

	return tpmodels.RecipeInterface{
		GetUserByID:                &getUserByID,
		GetUsersByEmail:            &getUserByEmail,
		GetUserByThirdPartyInfo:    &getUserByThirdPartyInfo,
		SignInUp:                   &signInUp,
		ManuallyCreateOrUpdateUser: &manuallyCreateOrUpdateUser,
		GetProvider:                &getProvider,
	}
}
