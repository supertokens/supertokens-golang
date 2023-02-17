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
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	tpproviders "github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier, providers []tpmodels.ProviderInput) tpmodels.RecipeInterface {

	getProvider := func(thirdPartyID string, tenantId *string, clientType *string, userContext supertokens.UserContext) (tpmodels.GetProviderResponse, error) {

		tenantConfig, err := multitenancy.GetTenantConfigWithContext(tenantId, userContext)
		if err != nil {
			return tpmodels.GetProviderResponse{}, err
		}

		mergedProviders := tpproviders.MergeProvidersFromCoreAndStatic(tenantId, tenantConfig.OK.ThirdParty.Providers, providers)
		provider, err := tpproviders.FindAndCreateProviderInstance(mergedProviders, thirdPartyID, tenantId, clientType, userContext)
		if err != nil {
			return tpmodels.GetProviderResponse{}, err
		}

		return tpmodels.GetProviderResponse{
			OK: &struct {
				Provider          *tpmodels.TypeProvider
				ThirdPartyEnabled bool
			}{
				Provider:          provider,
				ThirdPartyEnabled: tenantConfig.OK.ThirdParty.Enabled,
			},
		}, nil
	}

	signInUp := func(thirdPartyID, thirdPartyUserID string, email string, oAuthTokens tpmodels.TypeOAuthTokens, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, tenantId *string, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) {
		response, err := querier.SendPostRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/signinup", map[string]interface{}{
			"thirdPartyId":     thirdPartyID,
			"thirdPartyUserId": thirdPartyUserID,
			"email":            map[string]interface{}{"id": email},
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

	manuallyCreateOrUpdateUser := func(thirdPartyID, thirdPartyUserID string, email string, tenantId *string, userContext supertokens.UserContext) (tpmodels.ManuallyCreateOrUpdateUserResponse, error) {
		response, err := querier.SendPostRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/signinup", map[string]interface{}{
			"thirdPartyId":     thirdPartyID,
			"thirdPartyUserId": thirdPartyUserID,
			"email":            map[string]interface{}{"id": email},
		})
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

	getUserByID := func(userID string, tenantId *string, userContext supertokens.UserContext) (*tpmodels.User, error) {
		response, err := querier.SendGetRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/user", map[string]string{
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
	}

	getUserByThirdPartyInfo := func(thirdPartyID, thirdPartyUserID string, tenantId *string, userContext supertokens.UserContext) (*tpmodels.User, error) {
		response, err := querier.SendGetRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/user", map[string]string{
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
	}

	getUsersByEmail := func(email string, tenantId *string, userContext supertokens.UserContext) ([]tpmodels.User, error) {
		response, err := querier.SendGetRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/users/by-email", map[string]string{
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
