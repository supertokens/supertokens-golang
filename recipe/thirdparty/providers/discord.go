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
	"errors"
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const discordID = "discord"

type TypeDiscordInput struct {
	Config   []DiscordConfig
	Override func(provider *DiscordProvider) *DiscordProvider
}

type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
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

	discordProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (DiscordConfig, error) {
		if id == nil && len(input.Config) == 0 {
			return DiscordConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if id == nil && len(input.Config) > 1 {
			return DiscordConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if id == nil {
			return input.Config[0], nil
		}

		if id.Type == tpmodels.TypeClientID {
			for _, config := range input.Config {
				if config.ClientID == id.ID {
					return config, nil
				}
			}
		} else {
			// TODO Multitenant
		}

		return DiscordConfig{}, errors.New("config for specified clientID not found")
	}

	customProvider := CustomProvider(TypeCustomProviderInput{
		ThirdPartyID: discordID,
		Override: func(provider *TypeCustomProvider) *TypeCustomProvider {
			provider.GetConfig = func(ID *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
				discordConfig, err := discordProvider.GetConfig(ID, userContext)
				if err != nil {
					return CustomProviderConfig{}, err
				}

				authURL := "https://discord.com/api/oauth2/authorize"
				tokenURL := "https://discord.com/api/oauth2/token"
				userInfoURL := "https://discord.com/api/users/@me"

				return CustomProviderConfig{
					ClientID:     discordConfig.ClientID,
					ClientSecret: discordConfig.ClientSecret,
					Scope:        discordConfig.Scope,

					AuthorizationURL: &authURL,
					AccessTokenURL:   &tokenURL,
					UserInfoURL:      &userInfoURL,
					DefaultScope:     []string{"email", "identify"},

					GetSupertokensUserInfoFromRawUserInfoResponse: func(rawUserInfoResponse map[string]interface{}, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
						result := tpmodels.TypeUserInfo{}
						result.ThirdPartyUserId = fmt.Sprint(rawUserInfoResponse["id"])
						result.EmailInfo = &tpmodels.EmailStruct{
							ID: fmt.Sprint(rawUserInfoResponse["email"]),
						}
						emailVerified, emailVerifiedOk := rawUserInfoResponse["verified"].(bool)
						result.EmailInfo.IsVerified = emailVerified && emailVerifiedOk

						result.RawUserInfoFromProvider = rawUserInfoResponse

						return result, nil
					},
				}, nil
			}

			return provider
		},
	})

	discordProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		return customProvider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
	}

	discordProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		return customProvider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
	}

	discordProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		return customProvider.GetUserInfo(id, oAuthTokens, userContext)
	}

	if input.Override != nil {
		discordProvider = input.Override(discordProvider)
	}

	return *discordProvider.TypeProvider
}
