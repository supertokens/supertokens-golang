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

type GithubConfig = CustomProviderConfig

type TypeGithubInput struct {
	Config   []GithubConfig
	Override func(provider *GithubProvider) *GithubProvider
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

	var customProviderConfig []CustomProviderConfig
	if input.Config != nil {
		customProviderConfig = make([]CustomProviderConfig, len(input.Config))
		copy(customProviderConfig, input.Config)
	}

	customProvider := customProvider(TypeCustomProviderInput{
		ThirdPartyID: googleID,
		Config:       customProviderConfig,
	})

	{
		// Custom provider needs to use the config returned by google provider GetConfig
		// Also, google provider needs to use the default implementation of GetConfig provided by custom provider
		oGetConfig := customProvider.GetConfig
		customProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
			return githubProvider.GetConfig(id, userContext)
		}
		githubProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (GithubConfig, error) {
			return oGetConfig(id, userContext)
		}
	}

	{
		// Github provider APIs call into custom provider APIs

		githubProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return customProvider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
		}

		githubProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return customProvider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
		}
	}

	githubProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		config, err := githubProvider.GetConfig(id, userContext)
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

		rawUserInfoResponseFromProvider := tpmodels.TypeRawUserInfoFromProvider{FromAccessToken: rawResponse}
		userInfoResult, err := config.GetSupertokensUserInfoFromRawUserInfoResponse(rawUserInfoResponseFromProvider, userContext)
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
		githubProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (GithubConfig, error) {
			config, err := oGetConfig(id, userContext)
			if err != nil {
				return GithubConfig{}, err
			}
			return normalizeGithubConfig(config), nil
		}
	}

	return *githubProvider.TypeProvider
}

func normalizeGithubConfig(config GithubConfig) GithubConfig {
	if config.AuthorizationEndpoint == "" {
		config.AuthorizationEndpoint = "https://github.com/login/oauth/authorize"
	}
	if config.TokenEndpoint == "" {
		config.TokenEndpoint = "https://github.com/login/oauth/access_token"
	}
	if len(config.Scope) == 0 {
		config.Scope = []string{"read:user", "user:email"}
	}
	if config.GetSupertokensUserInfoFromRawUserInfoResponse == nil {
		config.GetSupertokensUserInfoFromRawUserInfoResponse = getSupertokensUserInfoFromRawUserInfoResponseForGithub
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
