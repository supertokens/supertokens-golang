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

const githubID = "github"

type TypeGithubInput struct {
	Config   []GithubConfig
	Override func(provider *GithubProvider) *GithubProvider
}

type GithubConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type GithubProvider struct {
	GetConfig func(id *tpmodels.TypeID, userContext supertokens.UserContext) (GithubConfig, error)
	*tpmodels.TypeProvider
}

func Github(input TypeGithubInput) tpmodels.TypeProvider {
	githubProvider := &GithubProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: githubID,
		},
	}

	githubProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (GithubConfig, error) {
		if id == nil && len(input.Config) == 0 {
			return GithubConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if id == nil && len(input.Config) > 1 {
			return GithubConfig{}, errors.New("please specify a clientID as there are multiple configs")
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

		return GithubConfig{}, errors.New("config for specified clientID not found")
	}

	customProvider := CustomProvider(TypeCustomProviderInput{
		ThirdPartyID: githubID,
		Override: func(provider *TypeCustomProvider) *TypeCustomProvider {
			provider.GetConfig = func(ID *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
				githubConfig, err := githubProvider.GetConfig(ID, userContext)
				if err != nil {
					return CustomProviderConfig{}, err
				}

				authURL := "https://github.com/login/oauth/authorize"
				tokenURL := "https://github.com/login/oauth/access_token"
				userInfoURL := "https://github.com/api/users/@me"

				return CustomProviderConfig{
					ClientID:     githubConfig.ClientID,
					ClientSecret: githubConfig.ClientSecret,
					Scope:        githubConfig.Scope,

					AuthorizationURL: &authURL,
					AccessTokenURL:   &tokenURL,
					UserInfoURL:      &userInfoURL,
					DefaultScope:     []string{"read:user", "user:email"},

					GetSupertokensUserInfoFromRawUserInfoResponse: func(rawUserInfoResponse map[string]interface{}, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
						result := tpmodels.TypeUserInfo{
							ThirdPartyUserId: fmt.Sprint(rawUserInfoResponse["user"].(map[string]interface{})["id"]),
						}

						emailsInfo := rawUserInfoResponse["emails"].([]interface{})
						for _, info := range emailsInfo {
							emailInfoMap := info.(map[string]interface{})
							if emailInfoMap["primary"].(bool) {
								verified, verifiedOk := emailInfoMap["verified"].(bool)
								result.EmailInfo = &tpmodels.EmailStruct{
									ID:         emailInfoMap["email"].(string),
									IsVerified: verified && verifiedOk,
								}
								break
							}
						}

						result.RawUserInfoFromProvider = rawUserInfoResponse

						return result, nil
					},
				}, nil
			}

			provider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
				config, err := provider.GetConfig(id, userContext)
				if err != nil {
					return tpmodels.TypeUserInfo{}, err
				}

				headers := map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", oAuthTokens["access_token"]),
					"Accept":        "application/vnd.github.v3+json",
				}
				rawResponse := map[string]interface{}{}
				status, emailInfo, err := doGetRequest("https://api.github.com/user/emails", nil, headers)
				if err != nil {
					return tpmodels.TypeUserInfo{}, err
				}
				if status >= 300 {
					return tpmodels.TypeUserInfo{}, errors.New("get email info returned a non 2xx status")
				}
				rawResponse["emails"] = emailInfo

				status, userInfo, err := doGetRequest("https://api.github.com/user", nil, headers)
				if err != nil {
					return tpmodels.TypeUserInfo{}, err
				}
				if status >= 300 {
					return tpmodels.TypeUserInfo{}, errors.New("get user info returned a non 2xx status")
				}
				rawResponse["user"] = userInfo

				userInfoResult, err := config.GetSupertokensUserInfoFromRawUserInfoResponse(rawResponse, userContext)
				if err != nil {
					return tpmodels.TypeUserInfo{}, err
				}
				return userInfoResult, nil
			}

			return provider
		},
	})

	githubProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		return customProvider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
	}

	githubProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		return customProvider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
	}

	githubProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		return customProvider.GetUserInfo(id, oAuthTokens, userContext)
	}

	if input.Override != nil {
		githubProvider = input.Override(githubProvider)
	}

	return *githubProvider.TypeProvider
}
