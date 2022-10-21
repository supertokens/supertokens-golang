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

type DiscordConfig = CustomProviderConfig

type TypeDiscordInput struct {
	Config   []DiscordConfig
	Override func(provider *DiscordProvider) *DiscordProvider
}

type DiscordProvider struct {
	GetConfig func(id *tpmodels.TypeID, userContext supertokens.UserContext) (DiscordConfig, error)
	*tpmodels.TypeProvider
}

func Discord(input TypeDiscordInput) tpmodels.TypeProvider {
	discordProvider := &DiscordProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: discordID,
		},
	}

	var customProviderConfig []CustomProviderConfig
	if input.Config != nil {
		customProviderConfig = make([]CustomProviderConfig, len(input.Config))
		for idx, config := range input.Config {
			customProviderConfig[idx] = config
		}
	}

	customProvider := customProvider(TypeCustomProviderInput{
		ThirdPartyID: discordID,
		Config:       customProviderConfig,
	})

	{
		// Custom provider needs to use the config returned by discord provider GetConfig
		// Also, discord provider needs to use the default implementation of GetConfig provided by custom provider
		oGetConfig := customProvider.GetConfig
		customProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
			return discordProvider.GetConfig(id, userContext)
		}
		discordProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (DiscordConfig, error) {
			return oGetConfig(id, userContext)
		}
	}

	{
		// Discord provider APIs call into custom provider APIs

		discordProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return customProvider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
		}

		discordProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return customProvider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
		}

		discordProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return customProvider.GetUserInfo(id, oAuthTokens, userContext)
		}
	}

	if input.Override != nil {
		discordProvider = input.Override(discordProvider)
	}

	{
		// We want to always normalize (for discord) the config before returning it
		oGetConfig := discordProvider.GetConfig
		discordProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (DiscordConfig, error) {
			config, err := oGetConfig(id, userContext)
			if err != nil {
				return DiscordConfig{}, err
			}
			return normalizeDiscordConfig(config), nil
		}
	}

	return *discordProvider.TypeProvider
}

func normalizeDiscordConfig(config DiscordConfig) DiscordConfig {
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
	if config.GetSupertokensUserInfoFromRawUserInfoResponse == nil {
		config.GetSupertokensUserInfoFromRawUserInfoResponse = getSupertokensUserInfoFromRawUserInfo("id", "email", "verified", "access_token")
	}
	return config
}
