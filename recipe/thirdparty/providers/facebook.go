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
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const facebookID = "facebook"

func Facebook(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.ThirdPartyId == "" {
		input.Config.ThirdPartyId = facebookID
	}
	if input.Config.Name == "" {
		input.Config.Name = "Facebook"
	}

	if input.Config.AuthorizationEndpoint == "" {
		input.Config.AuthorizationEndpoint = "https://www.facebook.com/v12.0/dialog/oauth"
	}

	if input.Config.TokenEndpoint == "" {
		input.Config.TokenEndpoint = "https://graph.facebook.com/v12.0/oauth/access_token"
	}

	if input.Config.UserInfoEndpoint == "" {
		input.Config.UserInfoEndpoint = "https://graph.facebook.com/me"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "id"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified = "email_verified"
	}

	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfigForClientType
		provider.GetConfigForClientType = func(clientType *string, input tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, input, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"email"}
			}

			return config, err
		}

		oGetUserInfo := provider.GetUserInfo
		provider.GetUserInfo = func(config tpmodels.ProviderConfigForClientType, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			if config.UserInfoEndpointQueryParams == nil {
				config.UserInfoEndpointQueryParams = map[string]interface{}{}
			}

			config.UserInfoEndpointQueryParams["access_token"] = oAuthTokens["access_token"]
			config.UserInfoEndpointQueryParams["fields"] = "id,email"
			config.UserInfoEndpointQueryParams["format"] = "json"

			if config.UserInfoEndpointHeaders == nil {
				config.UserInfoEndpointHeaders = map[string]interface{}{}
			}

			config.UserInfoEndpointHeaders["Authorization"] = nil

			return oGetUserInfo(config, oAuthTokens, userContext)
		}

		if oOverride != nil {
			provider = oOverride(provider)
		}
		return provider
	}

	return NewProvider(input)
}
