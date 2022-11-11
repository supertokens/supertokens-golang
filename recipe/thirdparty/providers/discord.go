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

const discordID = "discord"

func Discord(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	if input.ThirdPartyID == "" {
		input.ThirdPartyID = discordID
	}

	if input.Config.AuthorizationEndpoint == "" {
		input.Config.AuthorizationEndpoint = "https://discord.com/oauth2/authorize"
	}

	if input.Config.TokenEndpoint == "" {
		input.Config.TokenEndpoint = "https://discord.com/api/oauth2/token"
	}

	if input.Config.UserInfoEndpoint == "" {
		input.Config.UserInfoEndpoint = "https://discord.com/api/users/@me"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "id"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified = "verified"
	}

	if input.Config.AuthorizationEndpointQueryParams == nil {
		input.Config.AuthorizationEndpointQueryParams = map[string]interface{}{}
	}

	if input.Config.AuthorizationEndpointQueryParams["response_type"] == nil {
		input.Config.AuthorizationEndpointQueryParams["response_type"] = "code"
	}

	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfig
		provider.GetConfig = func(clientType *string, input tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClient, error) {
			config, err := oGetConfig(clientType, input, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClient{}, err
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"identify", "email"}
			}

			return config, err
		}

		if oOverride != nil {
			provider = oOverride(provider)
		}
		return provider
	}

	return NewProvider(input)
}
