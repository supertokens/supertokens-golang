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

type GithubConfig = OAuth2ProviderConfig

type TypeGithubInput struct {
	Config   GithubConfig
	Override func(provider *GithubProvider) *GithubProvider
}

type GithubProvider struct {
	*tpmodels.TypeProvider
}

func Github(input TypeGithubInput) tpmodels.TypeProvider {
	githubProvider := &GithubProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: githubID,
		},
	}

	oAuth2Provider := OAuth2Provider(TypeOAuth2ProviderInput{
		ThirdPartyID: githubID,
		Config:       input.Config,
	})

	{
		// Github provider APIs call into OAuth2 provider APIs

		githubProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
			return oAuth2Provider.GetConfig(clientType, tenantId, userContext)
		}

		githubProvider.GetAuthorisationRedirectURL = func(config tpmodels.TypeNormalisedProviderConfig, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return oAuth2Provider.GetAuthorisationRedirectURL(config, redirectURIOnProviderDashboard, userContext)
		}

		githubProvider.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.TypeNormalisedProviderConfig, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return oAuth2Provider.ExchangeAuthCodeForOAuthTokens(config, redirectInfo, userContext)
		}
	}

	githubProvider.GetUserInfo = func(config tpmodels.TypeNormalisedProviderConfig, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
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

		rawUserInfoResponseFromProvider := tpmodels.TypeRawUserInfoFromProvider{FromAccessToken: rawResponse}
		userInfoResult, err := getSupertokensUserInfoFromRawUserInfoResponseForGithub(rawUserInfoResponseFromProvider, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		return tpmodels.TypeUserInfo{
			ThirdPartyUserId:        userInfoResult.ThirdPartyUserId,
			EmailInfo:               userInfoResult.EmailInfo,
			RawUserInfoFromProvider: rawUserInfoResponseFromProvider,
		}, nil
	}

	if input.Override != nil {
		githubProvider = input.Override(githubProvider)
	}

	{
		// We want to always normalize (for github) the config before returning it
		oGetConfig := githubProvider.GetConfig
		githubProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
			config, err := oGetConfig(clientType, tenantId, userContext)
			if err != nil {
				return tpmodels.TypeNormalisedProviderConfig{}, err
			}
			return normalizeGithubConfig(config), nil
		}
	}

	return *githubProvider.TypeProvider
}

func normalizeGithubConfig(config tpmodels.TypeNormalisedProviderConfig) tpmodels.TypeNormalisedProviderConfig {
	if config.AuthorizationEndpoint == "" {
		config.AuthorizationEndpoint = "https://github.com/login/oauth/authorize"
	}
	if config.TokenEndpoint == "" {
		config.TokenEndpoint = "https://github.com/login/oauth/access_token"
	}
	if len(config.Scope) == 0 {
		config.Scope = []string{"read:user", "user:email"}
	}
	return config
}

func getSupertokensUserInfoFromRawUserInfoResponseForGithub(rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider, userContext supertokens.UserContext) (tpmodels.TypeSupertokensUserInfo, error) {
	if rawUserInfoResponse.FromAccessToken == nil {
		return tpmodels.TypeSupertokensUserInfo{}, errors.New("rawUserInfoResponse.FromAccessToken is not available")
	}

	result := tpmodels.TypeSupertokensUserInfo{
		ThirdPartyUserId: fmt.Sprint(rawUserInfoResponse.FromAccessToken["user"].(map[string]interface{})["id"]),
	}

	emailsInfo := rawUserInfoResponse.FromAccessToken["emails"].([]interface{})
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

	return result, nil
}
