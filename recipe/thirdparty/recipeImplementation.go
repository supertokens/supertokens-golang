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
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	tpproviders "github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier, providers []tpmodels.ProviderInput) tpmodels.RecipeInterface {

	getProvider := func(thirdPartyID string, clientType *string, tenantId string, userContext supertokens.UserContext) (*tpmodels.TypeProvider, error) {

		tenantConfig, err := multitenancy.GetTenant(tenantId, userContext)
		if err != nil {
			return nil, err
		}

		if tenantConfig == nil {
			return nil, errors.New("tenant not found")
		}

		mergedProviders := tpproviders.MergeProvidersFromCoreAndStatic(tenantConfig.ThirdParty.Providers, providers, tenantId == multitenancymodels.DefaultTenantId)
		provider, err := tpproviders.FindAndCreateProviderInstance(mergedProviders, thirdPartyID, clientType, userContext)
		if err != nil {
			return nil, err
		}

		return provider, nil
	}

	signInUp := func(thirdPartyID, thirdPartyUserID string, email string, oAuthTokens tpmodels.TypeOAuthTokens, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, tenantId string, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) {
		response, err := querier.SendPostRequest(tenantId+"/recipe/signinup", map[string]interface{}{
			"thirdPartyId":     thirdPartyID,
			"thirdPartyUserId": thirdPartyUserID,
			"email":            map[string]interface{}{"id": email},
		}, userContext)
		if err != nil {
			return tpmodels.SignInUpResponse{}, err
		}
		user, err := parseUser(response["user"])
		if err != nil {
			return tpmodels.SignInUpResponse{}, err
		}
		return tpmodels.SignInUpResponse{
			OK: &struct {
				CreatedNewUser          bool
				User                    tpmodels.User
				OAuthTokens             tpmodels.TypeOAuthTokens
				RawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider
			}{
				CreatedNewUser:          response["createdNewUser"].(bool),
				User:                    *user,
				OAuthTokens:             oAuthTokens,
				RawUserInfoFromProvider: rawUserInfoFromProvider,
			},
		}, nil
	}

	manuallyCreateOrUpdateUser := func(thirdPartyID, thirdPartyUserID string, email string, tenantId string, userContext supertokens.UserContext) (tpmodels.ManuallyCreateOrUpdateUserResponse, error) {
		response, err := querier.SendPostRequest(tenantId+"/recipe/signinup", map[string]interface{}{
			"thirdPartyId":     thirdPartyID,
			"thirdPartyUserId": thirdPartyUserID,
			"email":            map[string]interface{}{"id": email},
		}, userContext)
		if err != nil {
			return tpmodels.ManuallyCreateOrUpdateUserResponse{}, err
		}
		user, err := parseUser(response["user"])
		if err != nil {
			return tpmodels.ManuallyCreateOrUpdateUserResponse{}, err
		}
		return tpmodels.ManuallyCreateOrUpdateUserResponse{
			OK: &struct {
				CreatedNewUser bool
				User           tpmodels.User
			}{
				CreatedNewUser: response["createdNewUser"].(bool),
				User:           *user,
			},
		}, nil
	}

	getUserByID := func(userID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
		response, err := querier.SendGetRequest("/recipe/user", map[string]string{
			"userId": userID,
		}, userContext)
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
	}

	getUserByThirdPartyInfo := func(thirdPartyID, thirdPartyUserID string, tenantId string, userContext supertokens.UserContext) (*tpmodels.User, error) {
		response, err := querier.SendGetRequest(tenantId+"/recipe/user", map[string]string{
			"thirdPartyId":     thirdPartyID,
			"thirdPartyUserId": thirdPartyUserID,
		}, userContext)
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
	}

	getUsersByEmail := func(email string, tenantId string, userContext supertokens.UserContext) ([]tpmodels.User, error) {
		response, err := querier.SendGetRequest(tenantId+"/recipe/users/by-email", map[string]string{
			"email": email,
		}, userContext)
		if err != nil {
			return []tpmodels.User{}, err
		}
		users, err := parseUsers(response["users"])
		if err != nil {
			return []tpmodels.User{}, err
		}
		return users, nil
	}

	return tpmodels.RecipeInterface{
		GetUserByID:                &getUserByID,
		GetUsersByEmail:            &getUsersByEmail,
		GetUserByThirdPartyInfo:    &getUserByThirdPartyInfo,
		GetProvider:                &getProvider,
		SignInUp:                   &signInUp,
		ManuallyCreateOrUpdateUser: &manuallyCreateOrUpdateUser,
	}
}
