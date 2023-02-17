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
	GetUserByID             *func(userID string, tenantId *string, userContext supertokens.UserContext) (*User, error)
	GetUsersByEmail         *func(email string, tenantId *string, userContext supertokens.UserContext) ([]User, error)
	GetUserByThirdPartyInfo *func(thirdPartyID string, thirdPartyUserID string, tenantId *string, userContext supertokens.UserContext) (*User, error)

	GetProvider *func(thirdPartyID string, tenantId *string, clientType *string, userContext supertokens.UserContext) (GetProviderResponse, error)

	SignInUp                   *func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens TypeOAuthTokens, rawUserInfoFromProvider TypeRawUserInfoFromProvider, tenantId *string, userContext supertokens.UserContext) (SignInUpResponse, error)
	ManuallyCreateOrUpdateUser *func(thirdPartyID string, thirdPartyUserID string, email string, tenantId *string, userContext supertokens.UserContext) (ManuallyCreateOrUpdateUserResponse, error)
}

type GetProviderResponse struct {
	OK *struct {
		Provider          *TypeProvider
		ThirdPartyEnabled bool
	}
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
