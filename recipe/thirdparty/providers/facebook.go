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

const facebookID = "facebook"

type TypeFacebookInput struct {
	Config   []FacebookConfig
	Override func(provider *FacebookProvider) *FacebookProvider
}

type FacebookConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type FacebookProvider struct {
	GetConfig func(id *tpmodels.TypeID, userContext supertokens.UserContext) (FacebookConfig, error)
	*tpmodels.TypeProvider
}

func Facebook(input TypeFacebookInput) tpmodels.TypeProvider {
	facebookProvider := &FacebookProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: facebookID,
		},
	}

	facebookProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (FacebookConfig, error) {
		if id == nil && len(input.Config) == 0 {
			return FacebookConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if id == nil && len(input.Config) > 1 {
			return FacebookConfig{}, errors.New("please specify a clientID as there are multiple configs")
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

		return FacebookConfig{}, errors.New("config for specified clientID not found")
	}

	customProvider := CustomProvider(TypeCustomProviderInput{
		ThirdPartyID: facebookID,
		Override: func(provider *TypeCustomProvider) *TypeCustomProvider {
			provider.GetConfig = func(ID *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
				facebookConfig, err := facebookProvider.GetConfig(ID, userContext)
				if err != nil {
					return CustomProviderConfig{}, err
				}

				authURL := "https://www.facebook.com/v12.0/dialog/oauth"
				tokenURL := "https://graph.facebook.com/v12.0/oauth/access_token"
				userInfoURL := "https://graph.facebook.com/me"

				return CustomProviderConfig{
					ClientID:     facebookConfig.ClientID,
					ClientSecret: facebookConfig.ClientSecret,
					Scope:        facebookConfig.Scope,

					AuthorizationURL: &authURL,
					AccessTokenURL:   &tokenURL,
					UserInfoURL:      &userInfoURL,
					DefaultScope:     []string{"email"},

					GetSupertokensUserInfoFromRawUserInfoResponse: func(rawUserInfoResponse map[string]interface{}, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
						result := tpmodels.TypeUserInfo{}
						result.ThirdPartyUserId = fmt.Sprint(rawUserInfoResponse["id"])
						result.EmailInfo = &tpmodels.EmailStruct{
							ID: fmt.Sprint(rawUserInfoResponse["email"]),
						}
						emailVerified, emailVerifiedOk := rawUserInfoResponse["email_verified"].(bool)
						result.EmailInfo.IsVerified = emailVerified && emailVerifiedOk

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

				accessToken := oAuthTokens["access_token"].(string)

				queryParams := map[string]interface{}{
					"access_token": accessToken,
					"fields":       "id,email",
					"format":       "json",
				}

				status, userInfo, err := doGetRequest(*config.UserInfoURL, queryParams, nil)
				if err != nil {
					return tpmodels.TypeUserInfo{}, err
				}

				if status >= 300 {
					return tpmodels.TypeUserInfo{}, errors.New("get user info returned a non 2xx response")
				}

				userInfoResult, err := config.GetSupertokensUserInfoFromRawUserInfoResponse(userInfo.(map[string]interface{}), userContext)
				if err != nil {
					return tpmodels.TypeUserInfo{}, err
				}
				return userInfoResult, nil
			}

			return provider
		},
	})

	facebookProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		return customProvider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
	}

	facebookProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		return customProvider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
	}

	facebookProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		return customProvider.GetUserInfo(id, oAuthTokens, userContext)
	}

	if input.Override != nil {
		facebookProvider = input.Override(facebookProvider)
	}

	return customProvider
}
