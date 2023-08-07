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

const bitbucketID = "bitbucket"

func Bitbucket(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.Name == "" {
		input.Config.Name = "Bitbucket"
	}

	if input.Config.AuthorizationEndpoint == "" {
		input.Config.AuthorizationEndpoint = "https://bitbucket.org/site/oauth2/authorize"
	}

	if input.Config.TokenEndpoint == "" {
		input.Config.TokenEndpoint = "https://bitbucket.org/site/oauth2/access_token"
	}

	if input.Config.AuthorizationEndpointQueryParams == nil {
		input.Config.AuthorizationEndpointQueryParams = map[string]interface{}{
			"audience": "api.atlassian.com",
		}
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
				config.Scope = []string{"account", "email"}
			}

			return config, nil
		}

		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			accessToken, ok := oAuthTokens["access_token"].(string)
			if !ok {
				return tpmodels.TypeUserInfo{}, errors.New("access token not found")
			}

			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			rawUserInfoFromProvider := tpmodels.TypeRawUserInfoFromProvider{}
			userInfoFromAccessToken, err := doGetRequest(
				"https://api.bitbucket.org/2.0/user",
				nil,
				headers,
			)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
			rawUserInfoFromProvider.FromUserInfoAPI = userInfoFromAccessToken.(map[string]interface{})

			userInfoFromEmail, err := doGetRequest(
				"https://api.bitbucket.org/2.0/user/emails",
				nil,
				headers,
			)
			rawUserInfoFromProvider.FromUserInfoAPI["email"] = userInfoFromEmail

			email := ""
			isVerified := false

			for _, emailInfo := range userInfoFromEmail.(map[string]interface{})["values"].([]interface{}) {
				emailInfoMap := emailInfo.(map[string]interface{})
				if emailInfoMap["is_primary"].(bool) {
					email = emailInfoMap["email"].(string)
					isVerified = emailInfoMap["is_confirmed"].(bool)
					break
				}
			}

			if email == "" {
				return tpmodels.TypeUserInfo{
					ThirdPartyUserId:        fmt.Sprint(rawUserInfoFromProvider.FromUserInfoAPI["uuid"]),
					RawUserInfoFromProvider: rawUserInfoFromProvider,
				}, nil
			} else {
				return tpmodels.TypeUserInfo{
					ThirdPartyUserId: fmt.Sprint(rawUserInfoFromProvider.FromUserInfoAPI["uuid"]),
					Email: &tpmodels.EmailStruct{
						ID:         email,
						IsVerified: isVerified,
					},
					RawUserInfoFromProvider: rawUserInfoFromProvider,
				}, nil
			}
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}
		return originalImplementation
	}

	return NewProvider(input)
}
