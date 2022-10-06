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

const facebookID = "facebook"

func Facebook(input tpmodels.TypeFacebookInput) tpmodels.TypeProvider {
	facebookProvider := &tpmodels.FacebookProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: facebookID,
		},
	}

	getConfig := func(clientID *string, userContext supertokens.UserContext) (tpmodels.FacebookConfig, error) {
		if input.Config == nil || len(input.Config) == 0 {
			return tpmodels.FacebookConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if clientID == nil && len(input.Config) > 1 {
			return tpmodels.FacebookConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if clientID == nil && len(input.Config) == 1 {
			return input.Config[0], nil
		}

		for _, config := range input.Config {
			if config.ClientID == *clientID {
				return config, nil
			}
		}

		return tpmodels.FacebookConfig{}, errors.New("config for specified clientID not found")
	}

	getAuthorisationRedirectURL := func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		scopes := []string{"email"}
		config, err := (facebookProvider.GetConfig)(clientID, userContext)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		if config.Scope != nil {
			scopes = config.Scope
		}

		queryParams := map[string]interface{}{
			"scope":         strings.Join(scopes, " "),
			"response_type": "code",
			"client_id":     getActualClientIdFromDevelopmentClientId(config.ClientID),
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

		url := "https://www.facebook.com/v9.0/dialog/oauth"
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
			PKCECodeVerifier:   codeVerifier,
		}, nil
	}

	exchangeAuthCodeForOAuthTokens := func(clientID *string, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		config, err := facebookProvider.GetConfig(clientID, userContext)
		if err != nil {
			return nil, err
		}

		accessTokenAPIURL := "https://graph.facebook.com/v9.0/oauth/access_token"
		accessTokenAPIParams := map[string]string{
			"client_id":     getActualClientIdFromDevelopmentClientId(config.ClientID),
			"client_secret": config.ClientSecret,
			"code":          redirectURIInfo.RedirectURIQueryParams["code"].(string),
		}
		if config.ClientSecret == "" {
			if redirectURIInfo.PKCECodeVerifier == nil {
				return nil, errors.New("code verifier not found")
			}
			accessTokenAPIParams["code_verifier"] = *redirectURIInfo.PKCECodeVerifier
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
		var accessTokenAPIResponse facebookGetProfileInfoInput
		err = json.Unmarshal(authResponseJson, &accessTokenAPIResponse)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		response, err := getFacebookAuthRequest(accessTokenAPIResponse.AccessToken)
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

		isVerified, isVerifiedOk := userInfo["verified_email"].(bool)
		userInfoResult := tpmodels.TypeUserInfo{
			ThirdPartyUserId: ID,
			EmailInfo: &tpmodels.EmailStruct{
				ID:         email,
				IsVerified: isVerifiedOk && isVerified,
			},
			ResponseFromProvider: userInfo,
		}
		return userInfoResult, nil
	}

	facebookProvider.GetConfig = getConfig
	facebookProvider.GetAuthorisationRedirectURL = getAuthorisationRedirectURL
	facebookProvider.ExchangeAuthCodeForOAuthTokens = exchangeAuthCodeForOAuthTokens
	facebookProvider.GetUserInfo = getUserInfo

	if input.Override != nil {
		facebookProvider = input.Override(facebookProvider)
	}

	return *facebookProvider.TypeProvider
}

func getFacebookAuthRequest(accessToken string) (interface{}, error) {
	url := "https://graph.facebook.com/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("access_token", accessToken)
	q.Add("fields", "id,email")
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()
	return doGetRequest(req)
}

type facebookGetProfileInfoInput struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
