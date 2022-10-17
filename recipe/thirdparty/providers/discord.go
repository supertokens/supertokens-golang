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
	"net/http"
	"strings"

	"github.com/derekstavis/go-qs"
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
	GetConfig func(clientID *string, userContext supertokens.UserContext) (DiscordConfig, error)
	*tpmodels.TypeProvider
}

//	func Discord(config tpmodels.DiscordConfig) tpmodels.TypeProvider {
//		return tpmodels.TypeProvider{
//			ID: discordID,
//			Get: func(redirectURI, authCodeFromRequest *string, userContext supertokens.UserContext) tpmodels.TypeProviderGetResponse {
//				accessTokenAPIURL := "https://discord.com/api/oauth2/token"
//				accessTokenAPIParams := map[string]string{
//					"client_id":     config.ClientID,
//					"client_secret": config.ClientSecret,
//					"grant_type":    "authorization_code",
//				}
//				if authCodeFromRequest != nil {
//					accessTokenAPIParams["code"] = *authCodeFromRequest
//				}
//				if redirectURI != nil {
//					accessTokenAPIParams["redirect_uri"] = *redirectURI
//				}
func Discord(input TypeDiscordInput) tpmodels.TypeProvider {
	discordProvider := &DiscordProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: discordID,
		},
	}

	getConfig := func(clientID *string, userContext supertokens.UserContext) (DiscordConfig, error) {
		if input.Config == nil || len(input.Config) == 0 {
			return DiscordConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if clientID == nil && len(input.Config) > 1 {
			return DiscordConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if clientID == nil && len(input.Config) == 1 {
			return input.Config[0], nil
		}

		for _, config := range input.Config {
			if config.ClientID == *clientID {
				return config, nil
			}
		}

		return DiscordConfig{}, errors.New("config for specified clientID not found")
	}

	getAuthorisationRedirectURL := func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		scopes := []string{"email", "identify"}
		config, err := (discordProvider.GetConfig)(clientID, userContext)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		if config.Scope != nil {
			scopes = config.Scope
		}

		queryParams := map[string]interface{}{
			"scope":         strings.Join(scopes, " "),
			"client_id":     getActualClientIdFromDevelopmentClientId(config.ClientID),
			"response_type": "code",
		}

		url := "https://discord.com/api/oauth2/authorize"
		queryParams["redirect_uri"] = redirectURIOnProviderDashboard

		// url, queryParams, err = getAuthRedirectForDev(config.ClientID, url, queryParams)
		// if err != nil {
		// 	return tpmodels.TypeAuthorisationRedirect{}, err
		// }

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
		config, err := discordProvider.GetConfig(clientID, userContext)
		if err != nil {
			return nil, err
		}

		accessTokenAPIURL := "https://discord.com/api/oauth2/token"
		accessTokenAPIParams := map[string]string{
			"client_id":     getActualClientIdFromDevelopmentClientId(config.ClientID),
			"client_secret": config.ClientSecret,
			"grant_type":    "authorization_code",
			"code":          redirectURIInfo.RedirectURIQueryParams["code"].(string),
			"redirect_uri":  redirectURIInfo.RedirectURIOnProviderDashboard,
		}

		// redirectURI := checkDevAndGetRedirectURI(
		// 	config.ClientID,
		// 	redirectURIInfo.RedirectURIOnProviderDashboard,
		// 	userContext,
		// )

		// accessTokenAPIParams["redirect_uri"] = redirectURI

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
		response, err := getDiscordAuthRequest("Bearer " + oAuthTokens["access_token"].(string))
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		userInfo := response.(map[string]interface{})
		ID := userInfo["id"].(string)
		email := userInfo["email"].(string)
		if email == "" {
			userInfoResult := tpmodels.TypeUserInfo{
				ThirdPartyUserId:        ID,
				RawUserInfoFromProvider: userInfo,
			}
			return userInfoResult, nil
		}

		isVerified, isVerifiedOk := userInfo["verified"].(bool)
		userInfoResult := tpmodels.TypeUserInfo{
			ThirdPartyUserId: ID,
			EmailInfo: &tpmodels.EmailStruct{
				ID:         email,
				IsVerified: isVerifiedOk && isVerified,
			},
			RawUserInfoFromProvider: userInfo,
		}
		return userInfoResult, nil
	}

	discordProvider.GetConfig = getConfig
	discordProvider.GetAuthorisationRedirectURL = getAuthorisationRedirectURL
	discordProvider.ExchangeAuthCodeForOAuthTokens = exchangeAuthCodeForOAuthTokens
	discordProvider.GetUserInfo = getUserInfo

	if input.Override != nil {
		discordProvider = input.Override(discordProvider)
	}

	return *discordProvider.TypeProvider
}

func getDiscordAuthRequest(authHeader string) (interface{}, error) {
	url := "https://discord.com/api/users/@me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	return doGetRequest(req)
}
