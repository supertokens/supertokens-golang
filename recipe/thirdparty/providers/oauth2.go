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
	"strings"

	"github.com/derekstavis/go-qs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeOAuth2ProviderInput struct {
	ThirdPartyID string
	Config       OAuth2ProviderConfig
	Override     func(provider *TypeOAuth2Provider) *TypeOAuth2Provider
}

type OAuth2ProviderConfig struct {
	Clients []OAuth2ProviderClientConfig

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	ForcePKCE                        *bool
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      tpmodels.TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, config tpmodels.TypeNormalisedProviderConfig) (bool, error)
}

type OAuth2ProviderClientConfig struct {
	ClientType       string
	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}
}

type TypeSupertokensUserInfoMap struct {
	IdField            string
	EmailField         string
	EmailVerifiedField string
}

const scopeParameter = "scope"
const scopeSeparator = " "

type TypeOAuth2Provider = tpmodels.TypeProvider

func normalizeOAuth2ProviderInput(config OAuth2ProviderConfig) (OAuth2ProviderConfig, error) {
	if config.OIDCDiscoveryEndpoint != "" {
		oidcInfo, err := getOIDCDiscoveryInfo(config.OIDCDiscoveryEndpoint)

		if err == nil {
			if authURL, ok := oidcInfo["authorization_endpoint"].(string); ok {
				if config.AuthorizationEndpoint == "" {
					config.AuthorizationEndpoint = authURL
				}
			}

			if tokenURL, ok := oidcInfo["token_endpoint"].(string); ok {
				if config.TokenEndpoint == "" {
					config.TokenEndpoint = tokenURL
				}
			}

			if userInfoURL, ok := oidcInfo["userinfo_endpoint"].(string); ok {
				if config.UserInfoEndpoint == "" {
					config.UserInfoEndpoint = userInfoURL
				}
			}

			if jwksUri, ok := oidcInfo["jwks_uri"].(string); ok {
				config.JwksURI = jwksUri
			}
		}
	}

	return config, nil
}

func oAuth2Provider(input TypeOAuth2ProviderInput) *TypeOAuth2Provider {
	oAuth2Provider := &tpmodels.TypeProvider{
		ID: input.ThirdPartyID,
	}

	oAuth2Provider.GetConfig = func(clientType *string, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
		if clientType == nil {
			if len(input.Config.Clients) == 0 || len(input.Config.Clients) > 1 {
				return tpmodels.TypeNormalisedProviderConfig{}, errors.New("please provide exactly one client config or pass clientType or tenantId")
			}

			return tpmodels.TypeNormalisedProviderConfig{
				ClientID:                         input.Config.Clients[0].ClientID,
				ClientSecret:                     input.Config.Clients[0].ClientSecret,
				Scope:                            input.Config.Clients[0].Scope,
				ClientType:                       input.Config.Clients[0].ClientType,
				AuthorizationEndpoint:            input.Config.AuthorizationEndpoint,
				AuthorizationEndpointQueryParams: input.Config.AuthorizationEndpointQueryParams,
				TokenEndpoint:                    input.Config.TokenEndpoint,
				TokenParams:                      input.Config.TokenParams,
				ForcePKCE:                        input.Config.ForcePKCE,
				UserInfoEndpoint:                 input.Config.UserInfoEndpoint,
				JwksURI:                          input.Config.JwksURI,
				OIDCDiscoveryEndpoint:            input.Config.OIDCDiscoveryEndpoint,
				UserInfoMap:                      input.Config.UserInfoMap,
				ValidateIdTokenPayload:           input.Config.ValidateIdTokenPayload,
			}, nil

		}

		// (else) clientType is not nil
		if tenantId == nil {
			for _, config := range input.Config.Clients {
				if config.ClientType == *clientType {
					return tpmodels.TypeNormalisedProviderConfig{
						ClientID:                         config.ClientID,
						ClientSecret:                     config.ClientSecret,
						Scope:                            config.Scope,
						ClientType:                       config.ClientType,
						AuthorizationEndpoint:            input.Config.AuthorizationEndpoint,
						AuthorizationEndpointQueryParams: input.Config.AuthorizationEndpointQueryParams,
						TokenEndpoint:                    input.Config.TokenEndpoint,
						TokenParams:                      input.Config.TokenParams,
						ForcePKCE:                        input.Config.ForcePKCE,
						UserInfoEndpoint:                 input.Config.UserInfoEndpoint,
						JwksURI:                          input.Config.JwksURI,
						OIDCDiscoveryEndpoint:            input.Config.OIDCDiscoveryEndpoint,
						UserInfoMap:                      input.Config.UserInfoMap,
						ValidateIdTokenPayload:           input.Config.ValidateIdTokenPayload,
					}, nil
				}
			}

			return tpmodels.TypeNormalisedProviderConfig{}, errors.New("config for specified clientType not found")
		} else {
			// TODO Multitenant
			return tpmodels.TypeNormalisedProviderConfig{}, errors.New("needs implementation")
		}
	}

	oAuth2Provider.GetAuthorisationRedirectURL = func(config tpmodels.TypeNormalisedProviderConfig, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {

		queryParams := map[string]interface{}{
			scopeParameter:  strings.Join(config.Scope, scopeSeparator),
			"client_id":     config.ClientID,
			"redirect_uri":  redirectURIOnProviderDashboard,
			"response_type": "code",
		}
		var pkceCodeVerifier *string
		if config.ClientSecret == "" || (config.ForcePKCE != nil && *config.ForcePKCE) {
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

	oAuth2Provider.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.TypeNormalisedProviderConfig, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
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

		oAuthTokens, err := doPostRequest(tokenAPIURL, accessTokenAPIParams, nil)
		if err != nil {
			return nil, err
		}

		return oAuthTokens, nil
	}

	oAuth2Provider.GetUserInfo = func(config tpmodels.TypeNormalisedProviderConfig, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
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
			userInfoFromAccessToken, err := doGetRequest(config.UserInfoEndpoint, nil, headers)
			rawUserInfoFromProvider.FromAccessToken = userInfoFromAccessToken.(map[string]interface{})

			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
		}

		userInfoResult, err := getSupertokensUserInfoResultFromRawUserInfo(rawUserInfoFromProvider, config)
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
		oAuth2Provider = input.Override(oAuth2Provider)
	}

	return oAuth2Provider
}

func OAuth2Provider(input TypeOAuth2ProviderInput) tpmodels.TypeProvider {
	return *oAuth2Provider(input)
}

func getSupertokensUserInfoResultFromRawUserInfo(rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider, config tpmodels.TypeNormalisedProviderConfig) (tpmodels.TypeSupertokensUserInfo, error) {
	var rawUserInfo map[string]interface{}

	if config.UserInfoMap.From == tpmodels.FromIdTokenPayload {
		if rawUserInfoResponse.FromIdToken == nil {
			return tpmodels.TypeSupertokensUserInfo{}, errors.New("rawUserInfoResponse.FromIdToken is not available")
		}
		rawUserInfo = rawUserInfoResponse.FromIdToken

		if config.ValidateIdTokenPayload != nil {
			valid, err := config.ValidateIdTokenPayload(rawUserInfo, config)
			if err != nil {
				return tpmodels.TypeSupertokensUserInfo{}, err
			}
			if !valid {
				return tpmodels.TypeSupertokensUserInfo{}, errors.New("id_token payload is invalid")
			}
		}
	} else {
		if rawUserInfoResponse.FromAccessToken == nil {
			return tpmodels.TypeSupertokensUserInfo{}, errors.New("rawUserInfoResponse.FromAccessToken is not available")
		}
		rawUserInfo = rawUserInfoResponse.FromAccessToken
	}

	result := tpmodels.TypeSupertokensUserInfo{}
	result.ThirdPartyUserId = fmt.Sprint(accessField(rawUserInfo, config.UserInfoMap.IdField))
	result.EmailInfo = &tpmodels.EmailStruct{
		ID: fmt.Sprint(accessField(rawUserInfo, config.UserInfoMap.EmailField)),
	}
	if emailVerified, ok := accessField(rawUserInfo, config.UserInfoMap.EmailVerifiedField).(bool); ok {
		result.EmailInfo.IsVerified = emailVerified
	} else if emailVerified, ok := accessField(rawUserInfo, config.UserInfoMap.EmailVerifiedField).(string); ok {
		result.EmailInfo.IsVerified = emailVerified == "true"
	}
	return result, nil
}
