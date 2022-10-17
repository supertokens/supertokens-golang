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

type TypeGoogleInput struct {
	Config   []GoogleConfig
	Override func(provider *GoogleProvider) *GoogleProvider
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type GoogleProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (GoogleConfig, error)
	*tpmodels.TypeProvider
}

func Google(input TypeGoogleInput) tpmodels.TypeProvider {
	googleProvider := &GoogleProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: googleID,
		},
	}

	getConfig := func(clientID *string, userContext supertokens.UserContext) (GoogleConfig, error) {
		if len(input.Config) == 0 {
			return GoogleConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if clientID == nil && len(input.Config) > 1 {
			return GoogleConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if clientID == nil {
			return input.Config[0], nil
		}

		for _, config := range input.Config {
			if config.ClientID == *clientID {
				return config, nil
			}
		}

		return GoogleConfig{}, errors.New("config for specified clientID not found")
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

		url := "https://accounts.google.com/o/oauth2/v2/auth"
		queryParams := map[string]interface{}{
			"client_id":              config.ClientID,
			"scope":                  strings.Join(scopes, " "),
			"redirect_uri":           redirectURIOnProviderDashboard,
			"response_type":          "code",
			"access_type":            "offline",
			"include_granted_scopes": "true",
		}

		var codeVerifier *string
		if config.ClientSecret == "" {
			challenge, verifier, err := generateCodeChallengeS256(32)
			if err != nil {
				return tpmodels.TypeAuthorisationRedirect{}, err
			}
			queryParams["access_type"] = "online"
			queryParams["code_challenge"] = challenge
			queryParams["code_challenge_method"] = "S256"
			codeVerifier = &verifier
		}

		/* Transformation needed for dev keys BEGIN */
		if isUsingDevelopmentClientId(config.ClientID) {
			queryParams["client_id"] = getActualClientIdFromDevelopmentClientId(config.ClientID)
			queryParams["actual_redirect_uri"] = url
			url = DevOauthAuthorisationUrl
		}
		/* Transformation needed for dev keys END */

		queryParamsStr, err := qs.Marshal(queryParams)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}

		return tpmodels.TypeAuthorisationRedirect{
			URLWithQueryParams: url + "?" + queryParamsStr,
			PKCECodeVerifier:   codeVerifier,
		}, nil
	}

	exchangeAuthCodeForOAuthTokens := func(clientID *string, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		config, err := googleProvider.GetConfig(clientID, userContext)
		if err != nil {
			return nil, err
		}

		accessTokenAPIURL := "https://accounts.google.com/o/oauth2/token"
		accessTokenAPIParams := map[string]string{
			"client_id":    config.ClientID,
			"code":         redirectURIInfo.RedirectURIQueryParams["code"].(string),
			"redirect_url": redirectURIInfo.RedirectURIOnProviderDashboard,
			"grant_type":   "authorization_code",
		}
		if config.ClientSecret == "" {
			if redirectURIInfo.PKCECodeVerifier == nil {
				return nil, errors.New("code verifier not found")
			}
			accessTokenAPIParams["code_verifier"] = *redirectURIInfo.PKCECodeVerifier
		} else {
			accessTokenAPIParams["client_secret"] = config.ClientSecret
		}

		/* Transformation needed for dev keys BEGIN */
		if isUsingDevelopmentClientId(config.ClientID) {
			accessTokenAPIParams["client_id"] = getActualClientIdFromDevelopmentClientId(config.ClientID)
			accessTokenAPIParams["redirect_uri"] = DevOauthRedirectUrl
		}
		/* Transformation needed for dev keys END */

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
		response, err := getGoogleAuthRequest(accessTokenAPIResponse.AccessToken)
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

		isVerified, isVerifiedOk := userInfo["verified_email"].(bool)
		userInfoResult := tpmodels.TypeUserInfo{
			ThirdPartyUserId: ID,
			EmailInfo: &tpmodels.EmailStruct{
				ID:         email,
				IsVerified: isVerified && isVerifiedOk,
			},
			RawUserInfoFromProvider: userInfo,
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

	return *googleProvider.TypeProvider
}

func getGoogleAuthRequest(accessToken string) (interface{}, error) {
	url := "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	return doGetRequest(req)
}

type googleGetProfileInfoInput struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}
