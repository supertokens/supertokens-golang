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

package providers

import (
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tperrors"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func NewProvider(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	impl := &tpmodels.TypeProvider{
		ID: input.Config.ThirdPartyId,
	}

	// These are safe defaults common to most providers. Each provider implementations override these
	// as necessary
	if input.Config.UserInfoMap.FromIdTokenPayload.UserId == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.UserId = "sub"
	}
	if input.Config.UserInfoMap.FromIdTokenPayload.Email == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.Email = "email"
	}
	if input.Config.UserInfoMap.FromIdTokenPayload.EmailVerified == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.EmailVerified = "email_verified"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "sub"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified = "email_verified"
	}

	if input.Config.GenerateFakeEmail == nil {
		input.Config.GenerateFakeEmail = func(thirdPartyUserId string, tenantId string, userContext supertokens.UserContext) string {
			return fmt.Sprintf("%s@%s.fakeemail.com", thirdPartyUserId, input.Config.ThirdPartyId)
		}
	}

	impl.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
		inputConfig := input.Config

		if clientType == nil {
			if len(inputConfig.Clients) != 1 {
				return tpmodels.ProviderConfigForClientType{}, tperrors.ClientTypeNotFoundError{Msg: "please provide exactly one client config or pass clientType or tenantId"}
			}

			config := getProviderConfigForClient(inputConfig, inputConfig.Clients[0])
			return config, nil
		}

		for _, client := range inputConfig.Clients {
			if client.ClientType == *clientType {
				config := getProviderConfigForClient(input.Config, client)
				return config, nil
			}
		}

		return tpmodels.ProviderConfigForClientType{}, tperrors.ClientTypeNotFoundError{Msg: "Could not find client config for clientType: " + *clientType}
	}

	impl.GetAuthorisationRedirectURL = func(redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		return oauth2_GetAuthorisationRedirectURL(impl.Config, redirectURIOnProviderDashboard, userContext)
	}

	impl.ExchangeAuthCodeForOAuthTokens = func(redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		return oauth2_ExchangeAuthCodeForOAuthTokens(impl.Config, redirectURIInfo, userContext)
	}

	impl.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		return oauth2_GetUserInfo(impl.Config, oAuthTokens, userContext)
	}

	if input.Override != nil {
		impl = input.Override(impl)
	}

	return impl
}
