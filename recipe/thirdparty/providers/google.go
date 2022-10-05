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
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/derekstavis/go-qs"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const googleID = "google"

func Google(input tpmodels.TypeGoogleInput) (tpmodels.TypeProvider, error) {
	googleProvider := &tpmodels.GoogleProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: googleID,
		},
	}

	getConfig := func(clientID *string, userContext supertokens.UserContext) (tpmodels.GoogleConfig, error) {
		if clientID == nil && len(input.Config) > 1 {
			return tpmodels.GoogleConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if clientID == nil {
			return input.Config[0], nil
		}

		for _, config := range input.Config {
			if config.ClientID == *clientID {
				return config, nil
			}
		}

		return tpmodels.GoogleConfig{}, errors.New("config for specified clientID not found")
	}

	getAuthorisationRedirectURL := func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		scopes := []string{"https://www.googleapis.com/auth/userinfo.email"}
		config, err := googleProvider.GetConfig(clientID, userContext)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		if config.Scope != nil {
			scopes = config.Scope
		}

		queryParams := map[string]interface{}{
			"scope":                  strings.Join(scopes, " "),
			"access_type":            "offline",
			"include_granted_scopes": "true",
			"response_type":          "code",
			"client_id":              config.ClientID,
		}

		url := "https://accounts.google.com/o/oauth2/v2/auth"
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
		}, nil
	}

	exchangeAuthCodeForOAuthTokens := func(clientID *string, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		config, err := googleProvider.GetConfig(clientID, userContext)
		if err != nil {
			return nil, err
		}

		accessTokenAPIURL := "https://accounts.google.com/o/oauth2/token"
		accessTokenAPIParams := map[string]string{
			"client_id":     getActualClientIdFromDevelopmentClientId(config.ClientID),
			"client_secret": config.ClientSecret,
			"grant_type":    "authorization_code",
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
		authResponseJson, err := json.Marshal(oAuthTokens)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		var accessTokenAPIResponse googleGetProfileInfoInput
		err = json.Unmarshal(authResponseJson, &accessTokenAPIResponse)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		accessToken := accessTokenAPIResponse.AccessToken
		authHeader := "Bearer " + accessToken
		response, err := getGoogleAuthRequest(authHeader)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		userInfo := response.(map[string]interface{})
		ID := userInfo["id"].(string)
		email := userInfo["email"].(string)
		if email == "" {
			userInfoResult := tpmodels.TypeUserInfo{
				ThirdPartyUserId:     ID,
				ResponseFromProvider: userInfo,
			}
			return userInfoResult, nil
		}

		isVerified := userInfo["verified_email"].(bool)
		userInfoResult := tpmodels.TypeUserInfo{
			ThirdPartyUserId: ID,
			EmailInfo: &tpmodels.EmailStruct{
				ID:         email,
				IsVerified: isVerified,
			},
			ResponseFromProvider: userInfo,
		}
		return userInfoResult, nil
	}

	googleProvider.GetConfig = getConfig
	googleProvider.GetAuthorisationRedirectURL = getAuthorisationRedirectURL
	googleProvider.ExchangeAuthCodeForOAuthTokens = exchangeAuthCodeForOAuthTokens
	googleProvider.GetUserInfo = getUserInfo

	if input.Override != nil {
		googleProvider = input.Override(googleProvider)
	}

	if len(input.Config) == 0 && (&googleProvider.GetConfig == &getConfig) {
		// no config is provided and GetConfig is not overridden
		return tpmodels.TypeProvider{}, errors.New("please specify a config or override GetConfig")
	}

	return *googleProvider.TypeProvider, nil
}

func getGoogleAuthRequest(authHeader string) (interface{}, error) {
	url := "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	return doGetRequest(req)
}

type googleGetProfileInfoInput struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}
