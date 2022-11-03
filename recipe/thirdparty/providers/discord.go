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

type DiscordConfig = OAuth2ProviderConfig

type TypeDiscordInput struct {
	Config   DiscordConfig
	Override func(provider *DiscordProvider) *DiscordProvider
}

type DiscordProvider = tpmodels.TypeProvider

func Discord(input TypeDiscordInput) tpmodels.TypeProvider {
	discordProvider := &DiscordProvider{
		ID: discordID,
	}

	oAuth2Provider := oAuth2Provider(TypeOAuth2ProviderInput{
		ThirdPartyID: discordID,
		Config:       input.Config,
	})

	{
		// Discord provider APIs call into oAuth2 provider APIs

		discordProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
			return oAuth2Provider.GetConfig(clientType, tenantId, userContext)
		}

		discordProvider.GetAuthorisationRedirectURL = func(config tpmodels.TypeNormalisedProviderConfig, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return oAuth2Provider.GetAuthorisationRedirectURL(config, redirectURIOnProviderDashboard, userContext)
		}

		discordProvider.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.TypeNormalisedProviderConfig, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return oAuth2Provider.ExchangeAuthCodeForOAuthTokens(config, redirectInfo, userContext)
		}

		discordProvider.GetUserInfo = func(config tpmodels.TypeNormalisedProviderConfig, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return oAuth2Provider.GetUserInfo(config, oAuthTokens, userContext)
		}
	}

	if input.Override != nil {
		discordProvider = input.Override(discordProvider)
	}

	{
		// We want to always normalize (for discord) the config before returning it
		oGetConfig := discordProvider.GetConfig
		discordProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
			config, err := oGetConfig(clientType, tenantId, userContext)
			if err != nil {
				return tpmodels.TypeNormalisedProviderConfig{}, err
			}
			return normalizeDiscordConfig(config), nil
		}
	}

	return *discordProvider
}

func normalizeDiscordConfig(config tpmodels.TypeNormalisedProviderConfig) tpmodels.TypeNormalisedProviderConfig {
	if config.AuthorizationEndpoint == "" {
		config.AuthorizationEndpoint = "https://discord.com/api/oauth2/authorize"
	}
	if config.TokenEndpoint == "" {
		config.TokenEndpoint = "https://discord.com/api/oauth2/token"
	}
	if config.UserInfoEndpoint == "" {
		config.UserInfoEndpoint = "https://discord.com/api/users/@me"
	}
	if len(config.Scope) == 0 {
		config.Scope = []string{"identify", "email"}
	}

	if config.UserInfoMap.From == "" {
		config.UserInfoMap.From = tpmodels.FromAccessTokenPayload
	}

	if config.UserInfoMap.IdField == "" {
		config.UserInfoMap.IdField = "id"
	}

	if config.UserInfoMap.EmailField == "" {
		config.UserInfoMap.EmailField = "email"
	}

	if config.UserInfoMap.EmailVerifiedField == "" {
		config.UserInfoMap.EmailVerifiedField = "verified"
	}
	return config
}
