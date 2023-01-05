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

func Facebook(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
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

	input.Override = func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := originalImplementation.GetConfigForClientType
		originalImplementation.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"email"}
			}

			return config, nil
		}

		oGetUserInfo := originalImplementation.GetUserInfo
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			if originalImplementation.Config.UserInfoEndpointQueryParams == nil {
				originalImplementation.Config.UserInfoEndpointQueryParams = map[string]interface{}{}
			}

			originalImplementation.Config.UserInfoEndpointQueryParams["access_token"] = oAuthTokens["access_token"]
			originalImplementation.Config.UserInfoEndpointQueryParams["fields"] = "id,email"
			originalImplementation.Config.UserInfoEndpointQueryParams["format"] = "json"

			if originalImplementation.Config.UserInfoEndpointHeaders == nil {
				originalImplementation.Config.UserInfoEndpointHeaders = map[string]interface{}{}
			}

			originalImplementation.Config.UserInfoEndpointHeaders["Authorization"] = nil

			return oGetUserInfo(oAuthTokens, userContext)
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}
		return originalImplementation
	}

	return NewProvider(input)
}
