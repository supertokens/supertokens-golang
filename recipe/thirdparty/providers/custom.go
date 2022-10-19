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
	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeCustomProviderInput struct {
	ThirdPartyID string
	Config       []CustomProviderConfig
	Override     func(provider *TypeCustomProvider) *TypeCustomProvider
}

type CustomProviderConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string

	AuthorizationURL            *string
	AuthorizationURLQueryParams map[string]interface{}
	AccessTokenURL              *string
	AccessTokenMethod           *string
	AccessTokenParams           map[string]interface{}
	UserInfoURL                 *string
	UserInfoMethod              *string
	DefaultScope                []string
	ScopeParameter              *string
	ScopeSeparator              *string
	JwksURL                     *string
	OIDCEndpoint                *string

	GetSupertokensUserInfoFromRawUserInfoResponse func(rawUserInfoResponse map[string]interface{}, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error)
}

type TypeCustomProvider struct {
	GetConfig func(ID *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error)
	*tpmodels.TypeProvider
}

func normalizeCustomProviderInput(config CustomProviderConfig) CustomProviderConfig {
	if config.AccessTokenMethod == nil {
		post := http.MethodPost
		config.AccessTokenMethod = &post
	}
	if config.UserInfoMethod == nil {
		get := http.MethodGet
		config.UserInfoMethod = &get
	}
	if config.AuthorizationURLQueryParams == nil {
		config.AuthorizationURLQueryParams = map[string]interface{}{
			"response_type": "code",
		}
	}
	if config.AccessTokenParams == nil {
		config.AccessTokenParams = map[string]interface{}{
			"grant_type": "authorization_code",
		}
	}
	if config.ScopeParameter == nil {
		scope := "scope"
		config.ScopeParameter = &scope
	}
	if config.ScopeSeparator == nil {
		separator := " "
		config.ScopeSeparator = &separator
	}
	return config
}

func CustomProvider(input TypeCustomProviderInput) tpmodels.TypeProvider {

	customProvider := &TypeCustomProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: input.ThirdPartyID,
		},
	}

	customProvider.GetConfig = func(ID *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
		if ID == nil && len(input.Config) == 0 {
			return CustomProviderConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if ID == nil && len(input.Config) > 1 {
			return CustomProviderConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if ID == nil {
			return input.Config[0], nil
		}

		if ID.Type == tpmodels.TypeClientID {
			for _, config := range input.Config {
				if config.ClientID == ID.ID {
					return config, nil
				}
			}
		} else {
			// TODO Multitenant
		}

		return CustomProviderConfig{}, errors.New("config for specified ID not found")
	}

	customProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		config, err := customProvider.GetConfig(id, userContext)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		config = normalizeCustomProviderInput(config)

		scopes := config.DefaultScope
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		if config.Scope != nil {
			scopes = config.Scope
		}

		queryParams := map[string]interface{}{
			*config.ScopeParameter: strings.Join(scopes, *config.ScopeSeparator),
			"client_id":            config.ClientID,
		}
		for k, v := range config.AuthorizationURLQueryParams {
			queryParams[k] = v
		}

		url := *config.AuthorizationURL
		queryParams["redirect_uri"] = redirectURIOnProviderDashboard

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
			PKCECodeVerifier:   nil,
		}, nil
	}

	customProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		config, err := customProvider.GetConfig(id, userContext)
		if err != nil {
			return nil, err
		}
		config = normalizeCustomProviderInput(config)

		accessTokenAPIURL := *config.AccessTokenURL
		accessTokenAPIParams := map[string]interface{}{
			"client_id":    config.ClientID,
			"code":         redirectURIInfo.RedirectURIQueryParams["code"].(string),
			"redirect_uri": redirectURIInfo.RedirectURIOnProviderDashboard,
		}

		for k, v := range config.AccessTokenParams {
			accessTokenAPIParams[k] = v
		}

		if config.ClientSecret != "" {
			accessTokenAPIParams["client_secret"] = config.ClientSecret
		}

		if redirectURIInfo.PKCECodeVerifier != nil {
			accessTokenAPIParams["code_verifier"] = *redirectURIInfo.PKCECodeVerifier
		}

		/* Transformation needed for dev keys BEGIN */
		if isUsingDevelopmentClientId(config.ClientID) {
			accessTokenAPIParams["client_id"] = getActualClientIdFromDevelopmentClientId(config.ClientID)
			accessTokenAPIParams["redirect_uri"] = DevOauthRedirectUrl
		}
		/* Transformation needed for dev keys END */

		var status int
		var oAuthTokens interface{}
		if *config.AccessTokenMethod == http.MethodPost {
			status, oAuthTokens, err = doPostRequest(accessTokenAPIURL, accessTokenAPIParams, nil)
		} else {
			status, oAuthTokens, err = doGetRequest(accessTokenAPIURL, accessTokenAPIParams, nil)
		}
		if err != nil {
			return nil, err
		}

		if status >= 300 {
			return nil, errors.New("AccessToken API returned a non 2xx response")
		}

		return oAuthTokens.(map[string]interface{}), nil
	}

	customProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		config, err := customProvider.GetConfig(id, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		config = normalizeCustomProviderInput(config)

		var userInfo interface{}
		accessToken, accessTokenOk := oAuthTokens["access_token"].(string)
		idToken, idTokenOk := oAuthTokens["id_token"].(string)

		if idTokenOk && config.JwksURL != nil {
			claims := jwt.MapClaims{}
			jwksURL := *config.JwksURL
			jwks, err := getJWKSFromURL(jwksURL)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
			token, err := jwt.ParseWithClaims(idToken, claims, jwks.Keyfunc)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}

			if !token.Valid {
				return tpmodels.TypeUserInfo{}, errors.New("invalid id_token supplied")
			}
			userInfo = map[string]interface{}(claims)
		} else if accessTokenOk && config.UserInfoURL != nil {
			var status int

			if *config.UserInfoMethod == http.MethodGet {
				headers := map[string]string{
					"Authorization": "Bearer " + accessToken,
				}
				status, userInfo, err = doGetRequest(*config.UserInfoURL, nil, headers)
			} else {
				params := map[string]interface{}{
					"access_token": accessToken,
				}
				status, userInfo, err = doPostRequest(*config.UserInfoURL, params, nil)
			}
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}

			if status >= 300 {
				return tpmodels.TypeUserInfo{}, errors.New("UserInfo API returned a non 2xx response")
			}
		} else {
			return tpmodels.TypeUserInfo{}, errors.New("misconfigured custom provider, unable to fetch user info using access_token or id_token")
		}

		userInfoResult, err := config.GetSupertokensUserInfoFromRawUserInfoResponse(userInfo.(map[string]interface{}), userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		return userInfoResult, nil
	}

	if input.Override != nil {
		customProvider = input.Override(customProvider)
	}

	return *customProvider.TypeProvider
}
