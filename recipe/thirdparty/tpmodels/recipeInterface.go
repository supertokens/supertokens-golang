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

package tpmodels

import "github.com/supertokens/supertokens-golang/supertokens"

type RecipeInterface struct {
	GetUserByID                *func(userID string, userContext supertokens.UserContext) (*User, error)
	GetUsersByEmail            *func(email string, userContext supertokens.UserContext) ([]User, error)
	GetUserByThirdPartyInfo    *func(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*User, error)
	SignInUp                   *func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens TypeOAuthTokens, rawUserInfoFromProvider TypeRawUserInfoFromProvider, tenantId *string, userContext supertokens.UserContext) (SignInUpResponse, error)
	ManuallyCreateOrUpdateUser *func(thirdPartyID string, thirdPartyUserId string, email string, userContext supertokens.UserContext) (ManuallyCreateOrUpdateUserResponse, error)

	CreateOrUpdateThirdPartyConfig *func(thirdPartyId string, tenantId *string, config ProviderConfig, userContext supertokens.UserContext) (CreateOrUpdateTenantIdConfigResponse, error)
	FetchThirdPartyConfig          *func(thirdPartyId string, tenantId *string, userContext supertokens.UserContext) (FetchTenantIdConfigResponse, error)
	DeleteThirdPartyConfig         *func(thirdPartyId string, tenantId *string, userContext supertokens.UserContext) (DeleteTenantIdConfigResponse, error)
	ListThirdPartyConfigs          *func(tenantId *string, userContext supertokens.UserContext) (ListTenantConfigMappingsResponse, error)
}

type SignInUpResponse struct {
	OK *struct {
		CreatedNewUser          bool
		User                    User
		OAuthTokens             TypeOAuthTokens
		RawUserInfoFromProvider TypeRawUserInfoFromProvider
	}
}

type ManuallyCreateOrUpdateUserResponse struct {
	OK *struct {
		CreatedNewUser bool
		User           User
	}
}

type CreateOrUpdateTenantIdConfigResponse struct {
	OK *struct {
		CreatedNew bool
	}
}

type FetchTenantIdConfigResponse struct {
	OK *struct {
		Config ProviderConfig
	}
	UnknownMappingError *struct{}
}

type DeleteTenantIdConfigResponse struct {
	OK *struct {
		DidMappingExist bool
	}
}

type ListTenantConfigMappingsResponse struct {
	OK *struct {
		Configs []struct {
			ThirdPartyId string
			Config       ProviderConfig
		}
	}
}
