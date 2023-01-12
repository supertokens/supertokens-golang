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

func Github(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.Name == "" {
		input.Config.Name = "Github"
	}

	if input.Config.AuthorizationEndpoint == "" {
		input.Config.AuthorizationEndpoint = "https://github.com/login/oauth/authorize"
	}

	if input.Config.TokenEndpoint == "" {
		input.Config.TokenEndpoint = "https://github.com/login/oauth/access_token"
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
				config.Scope = []string{"read:user", "user:email"}
			}

			return config, nil
		}

		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			headers := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", oAuthTokens["access_token"]),
				"Accept":        "application/vnd.github.v3+json",
			}
			rawResponse := map[string]interface{}{}
			emailInfo, err := doGetRequest("https://api.github.com/user/emails", nil, headers)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
			rawResponse["emails"] = emailInfo

			userInfo, err := doGetRequest("https://api.github.com/user", nil, headers)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
			rawResponse["user"] = userInfo

			rawUserInfoResponseFromProvider := tpmodels.TypeRawUserInfoFromProvider{FromUserInfoAPI: rawResponse}
			userInfoResult, err := getSupertokensUserInfoFromRawUserInfoResponseForGithub(rawUserInfoResponseFromProvider)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
			return tpmodels.TypeUserInfo{
				ThirdPartyUserId:        userInfoResult.ThirdPartyUserId,
				Email:                   userInfoResult.Email,
				RawUserInfoFromProvider: rawUserInfoResponseFromProvider,
			}, nil
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}
		return originalImplementation
	}

	return NewProvider(input)
}

func getSupertokensUserInfoFromRawUserInfoResponseForGithub(rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider) (tpmodels.TypeUserInfo, error) {
	if rawUserInfoResponse.FromUserInfoAPI == nil {
		return tpmodels.TypeUserInfo{}, errors.New("rawUserInfoResponse.FromUserInfoAPI is not available")
	}

	result := tpmodels.TypeUserInfo{
		ThirdPartyUserId: fmt.Sprint(rawUserInfoResponse.FromUserInfoAPI["user"].(map[string]interface{})["id"]),
	}

	emailsInfo := rawUserInfoResponse.FromUserInfoAPI["emails"].([]interface{})
	for _, info := range emailsInfo {
		emailInfoMap := info.(map[string]interface{})
		if emailInfoMap["primary"].(bool) {
			verified, verifiedOk := emailInfoMap["verified"].(bool)
			result.Email = &tpmodels.EmailStruct{
				ID:         emailInfoMap["email"].(string),
				IsVerified: verified && verifiedOk,
			}
			break
		}
	}

	return result, nil
}
