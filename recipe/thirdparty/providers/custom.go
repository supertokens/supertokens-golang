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

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}

	TokenEndpoint string
	TokenParams   map[string]interface{}

	UserInfoEndpoint string

	JwksURI      string
	OIDCEndpoint string

	GetSupertokensUserInfoFromRawUserInfoResponse func(rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider, userContext supertokens.UserContext) (tpmodels.TypeSupertokensUserInfo, error)

	AdditionalConfig map[string]interface{}
}

const ScopeParameter = "scope"
const ScopeSeparator = " "

type TypeCustomProvider struct {
	GetConfig func(id *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error)
	*tpmodels.TypeProvider
}

func normalizeCustomProviderInput(config CustomProviderConfig) (CustomProviderConfig, error) {
	if config.Scope == nil {
		config.Scope = []string{}
	}
	if config.OIDCEndpoint != "" {
		// TODO cache this value for 24 hours
		status, oidcInfo, err := doGetRequest(config.OIDCEndpoint, nil, nil)

		if err == nil && status < 300 {
			oidcInfoMap := oidcInfo.(map[string]interface{})
			if authURL, ok := oidcInfoMap["authorization_endpoint"].(string); ok {
				if config.AuthorizationEndpoint == "" {
					config.AuthorizationEndpoint = authURL
				}
			}

			if tokenURL, ok := oidcInfoMap["token_endpoint"].(string); ok {
				if config.TokenEndpoint == "" {
					config.TokenEndpoint = tokenURL
				}
			}

			if userInfoURL, ok := oidcInfoMap["userinfo_endpoint"].(string); ok {
				if config.UserInfoEndpoint == "" {
					config.UserInfoEndpoint = userInfoURL
				}
			}

			if jwksUri, ok := oidcInfoMap["jwks_uri"].(string); ok {
				config.JwksURI = jwksUri
			}
		}
	}

	return config, nil
}

func customProvider(input TypeCustomProviderInput) *TypeCustomProvider {
	customProvider := &TypeCustomProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: input.ThirdPartyID,
		},
	}

	customProvider.GetConfig = func(ID *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
		if ID == nil {
			if len(input.Config) == 0 || len(input.Config) > 1 {
				return CustomProviderConfig{}, errors.New("please provide exactly one config or pass ClientID or TenantID")
			}

			return input.Config[0], nil
		}

		if ID.Type == tpmodels.TypeClientID {
			for _, config := range input.Config {
				if config.ClientID == ID.ID {
					return config, nil
				}
			}

			return CustomProviderConfig{}, errors.New("config for specified ClientID not found")
		} else {
			// TODO Multitenant
			return CustomProviderConfig{}, errors.New("needs implementation")
		}
	}

	customProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		config, err := customProvider.GetConfig(id, userContext)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		config, err = normalizeCustomProviderInput(config)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}

		queryParams := map[string]interface{}{
			ScopeParameter:  strings.Join(config.Scope, ScopeSeparator),
			"client_id":     config.ClientID,
			"redirect_uri":  redirectURIOnProviderDashboard,
			"response_type": "code",
		}
		var pkceCodeVerifier *string
		if config.ClientSecret == "" {
			challenge, verifier, err := generateCodeChallengeS256(32)
			if err != nil {
				return tpmodels.TypeAuthorisationRedirect{}, err
			}
			queryParams["code_challenge"] = challenge
			queryParams["code_challenge_method"] = "S256"
			pkceCodeVerifier = &verifier
		}

		for k, v := range config.AuthorizationEndpointQueryParams {
			if v == nil {
				delete(queryParams, k)
			} else {
				queryParams[k] = v
			}
		}

		url := config.AuthorizationEndpoint

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
			PKCECodeVerifier:   pkceCodeVerifier,
		}, nil
	}

	customProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		config, err := customProvider.GetConfig(id, userContext)
		if err != nil {
			return nil, err
		}
		config, err = normalizeCustomProviderInput(config)
		if err != nil {
			return nil, err
		}

		tokenAPIURL := config.TokenEndpoint
		accessTokenAPIParams := map[string]interface{}{
			"client_id":    config.ClientID,
			"redirect_uri": redirectURIInfo.RedirectURIOnProviderDashboard,
			"code":         redirectURIInfo.RedirectURIQueryParams["code"].(string),
			"grant_type":   "authorization_code",
		}
		if config.ClientSecret != "" {
			accessTokenAPIParams["client_secret"] = config.ClientSecret
		}
		if redirectURIInfo.PKCECodeVerifier != nil {
			accessTokenAPIParams["code_verifier"] = *redirectURIInfo.PKCECodeVerifier
		}

		for k, v := range config.TokenParams {
			if v == nil {
				delete(accessTokenAPIParams, k)
			} else {
				accessTokenAPIParams[k] = v
			}
		}

		/* Transformation needed for dev keys BEGIN */
		if isUsingDevelopmentClientId(config.ClientID) {
			accessTokenAPIParams["client_id"] = getActualClientIdFromDevelopmentClientId(config.ClientID)
			accessTokenAPIParams["redirect_uri"] = DevOauthRedirectUrl
		}
		/* Transformation needed for dev keys END */

		status, oAuthTokens, err := doPostRequest(tokenAPIURL, accessTokenAPIParams, nil)
		if err != nil {
			return nil, err
		}

		if status >= 300 {
			// TODO add debug logs
			return nil, errors.New("AccessToken API returned a non 2xx response") // TODO add status code and response
		}

		return oAuthTokens, nil
	}

	customProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		config, err := customProvider.GetConfig(id, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		config, err = normalizeCustomProviderInput(config)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		accessToken, accessTokenOk := oAuthTokens["access_token"].(string)
		idToken, idTokenOk := oAuthTokens["id_token"].(string)

		rawUserInfoFromProvider := tpmodels.TypeRawUserInfoFromProvider{}

		if idTokenOk && config.JwksURI != "" {
			claims := jwt.MapClaims{}
			jwksURL := config.JwksURI
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
			rawUserInfoFromProvider.FromIdToken = map[string]interface{}(claims)
		}
		if accessTokenOk && config.UserInfoEndpoint != "" {
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			status, userInfoFromAccessToken, err := doGetRequest(config.UserInfoEndpoint, nil, headers)
			rawUserInfoFromProvider.FromAccessToken = userInfoFromAccessToken.(map[string]interface{})

			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}

			if status >= 300 {
				return tpmodels.TypeUserInfo{}, errors.New("UserInfo API returned a non 2xx response") // TODO Add status code and response
			}
		}

		userInfoResult, err := config.GetSupertokensUserInfoFromRawUserInfoResponse(rawUserInfoFromProvider, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		return tpmodels.TypeUserInfo{
			ThirdPartyUserId:        userInfoResult.ThirdPartyUserId,
			EmailInfo:               userInfoResult.EmailInfo,
			RawUserInfoFromProvider: rawUserInfoFromProvider,
		}, nil
	}

	if input.Override != nil {
		customProvider = input.Override(customProvider)
	}

	return customProvider
}

func CustomProvider(input TypeCustomProviderInput) tpmodels.TypeProvider {
	return *customProvider(input).TypeProvider
}
