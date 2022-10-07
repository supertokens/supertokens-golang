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
	"net/http"
	"strings"

	"github.com/derekstavis/go-qs"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const githubID = "github"

func Github(input tpmodels.TypeGithubInput) tpmodels.TypeProvider {
	githubProvider := &tpmodels.GithubProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: githubID,
		},
	}

	getConfig := func(clientID *string, userContext supertokens.UserContext) (tpmodels.GithubConfig, error) {
		if input.Config == nil || len(input.Config) == 0 {
			return tpmodels.GithubConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if clientID == nil && len(input.Config) > 1 {
			return tpmodels.GithubConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if clientID == nil && len(input.Config) == 1 {
			return input.Config[0], nil
		}

		for _, config := range input.Config {
			if config.ClientID == *clientID {
				return config, nil
			}
		}

		return tpmodels.GithubConfig{}, errors.New("config for specified clientID not found")
	}

	getAuthorisationRedirectURL := func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		scopes := []string{"read:user", "user:email"}
		config, err := (githubProvider.GetConfig)(clientID, userContext)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		if config.Scope != nil {
			scopes = config.Scope
		}

		queryParams := map[string]interface{}{
			"scope":     strings.Join(scopes, " "),
			"client_id": getActualClientIdFromDevelopmentClientId(config.ClientID),
		}

		url := "https://github.com/login/oauth/authorize"
		queryParams["redirect_uri"] = redirectURIOnProviderDashboard

		url, queryParams, err = getAuthRedirectForDev(config.ClientID, url, queryParams)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}

		queryParamsStr, err := qs.Marshal(queryParams)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}

		return tpmodels.TypeAuthorisationRedirect{
			URLWithQueryParams: url + "?" + queryParamsStr,
			PKCECodeVerifier:   nil,
		}, nil
	}

	exchangeAuthCodeForOAuthTokens := func(clientID *string, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		config, err := githubProvider.GetConfig(clientID, userContext)
		if err != nil {
			return nil, err
		}

		accessTokenAPIURL := "https://github.com/login/oauth/access_token"
		accessTokenAPIParams := map[string]string{
			"client_id":     getActualClientIdFromDevelopmentClientId(config.ClientID),
			"client_secret": config.ClientSecret,
			"code":          redirectURIInfo.RedirectURIQueryParams["code"].(string),
		}

		redirectURI := checkDevAndGetRedirectURI(
			config.ClientID,
			redirectURIInfo.RedirectURIOnProviderDashboard,
			userContext,
		)

		accessTokenAPIParams["redirect_uri"] = redirectURI

		authResponseFromRequest, err := postRequest(accessTokenAPIURL, accessTokenAPIParams)
		if err != nil {
			return nil, err
		}

		authResponse := tpmodels.TypeOAuthTokens{}

		for k, v := range authResponseFromRequest {
			authResponse[k] = v
		}

		return authResponse, nil
	}

	getUserInfo := func(clientID *string, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		authHeader := "Bearer " + oAuthTokens["access_token"].(string)
		response, err := getGithubAuthRequest(authHeader)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		userInfo := response.(map[string]interface{})
		ID := fmt.Sprint(userInfo["id"]) // github userId will be a number

		emailsInfoResponse, err := getGithubEmailsInfo(authHeader)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		emailsInfo := emailsInfoResponse.([]interface{})
		var emailInfo map[string]interface{}
		for _, info := range emailsInfo {
			emailInfoMap := info.(map[string]interface{})
			if primary, ok := emailInfoMap["primary"].(bool); primary && ok {
				emailInfo = emailInfoMap
				break
			}
		}
		if emailInfo == nil {
			return tpmodels.TypeUserInfo{
				ThirdPartyUserId: ID,
				ResponseFromProvider: map[string]interface{}{
					"userInfo":  userInfo,
					"emailInfo": emailsInfo,
				},
			}, nil
		}

		email := emailInfo["email"].(string)
		isVerified, isVerifiedOk := emailInfo["verified"].(bool)
		userInfoResult := tpmodels.TypeUserInfo{
			ThirdPartyUserId: ID,
			EmailInfo: &tpmodels.EmailStruct{
				ID:         email,
				IsVerified: isVerifiedOk && isVerified,
			},
			ResponseFromProvider: map[string]interface{}{
				"userInfo":  userInfo,
				"emailInfo": emailsInfo,
			},
		}
		return userInfoResult, nil
	}

	githubProvider.GetConfig = getConfig
	githubProvider.GetAuthorisationRedirectURL = getAuthorisationRedirectURL
	githubProvider.ExchangeAuthCodeForOAuthTokens = exchangeAuthCodeForOAuthTokens
	githubProvider.GetUserInfo = getUserInfo

	if input.Override != nil {
		githubProvider = input.Override(githubProvider)
	}

	return *githubProvider.TypeProvider
}

func getGithubAuthRequest(authHeader string) (interface{}, error) {
	url := "https://api.github.com/user"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	return doGetRequest(req)
}

func getGithubEmailsInfo(authHeader string) (interface{}, error) {
	url := "https://api.github.com/user/emails"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	return doGetRequest(req)
}
