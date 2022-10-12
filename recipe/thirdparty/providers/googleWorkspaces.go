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
	"fmt"
	"strings"

	"github.com/derekstavis/go-qs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const googleWorkspacesID = "google-workspaces"

func GoogleWorkspaces(input tpmodels.TypeGoogleWorkspacesInput) (tpmodels.TypeProvider, error) {
	googleWorkspacesProvider := &tpmodels.GoogleWorkspacesProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: googleWorkspacesID,
		},
	}

	getConfig := func(clientID *string, userContext supertokens.UserContext) (tpmodels.GoogleWorkspacesConfig, error) {
		if clientID == nil && len(input.Config) > 1 {
			return tpmodels.GoogleWorkspacesConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if clientID == nil {
			return input.Config[0], nil
		}

		for _, config := range input.Config {
			if config.ClientID == *clientID {
				return config, nil
			}
		}

		return tpmodels.GoogleWorkspacesConfig{}, errors.New("config for specified clientID not found")
	}

	getTenantID := func(clientID *string, userContext supertokens.UserContext) (string, error) {
		return "none", nil
	}

	getAuthorisationRedirectURL := func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		scopes := []string{"https://www.googleapis.com/auth/userinfo.email"}
		config, err := googleWorkspacesProvider.GetConfig(clientID, userContext)
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
			"hd":                     "*",
		}

		if config.Domain != nil {
			queryParams["hd"] = *config.Domain
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
			PKCECodeVerifier:   codeVerifier,
		}, nil
	}

	exchangeAuthCodeForOAuthTokens := func(clientID *string, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		config, err := googleWorkspacesProvider.GetConfig(clientID, userContext)
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
		config, err := googleWorkspacesProvider.GetConfig(clientID, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		authResponseJson, err := json.Marshal(oAuthTokens)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		var accessTokenAPIResponse googleGetProfileInfoInput
		err = json.Unmarshal(authResponseJson, &accessTokenAPIResponse)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		claims, err := verifyAndGetClaims(oAuthTokens["id_token"].(string), getActualClientIdFromDevelopmentClientId(config.ClientID))
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		var email string
		var isVerified bool
		var id string
		var hd string

		for key, val := range claims {
			if key == "sub" {
				id = val.(string)
			} else if key == "email" {
				email = val.(string)
			} else if key == "email_verified" {
				isVerified = val.(bool)
			} else if key == "hd" {
				hd = val.(string)
			}
		}

		if email == "" {
			return tpmodels.TypeUserInfo{}, errors.New("Could not get email. Please use a different login method")
		}

		if hd == "" {
			return tpmodels.TypeUserInfo{}, errors.New("Please use a GoogleWorkspaces Workspace ID to login")
		}

		domain := "*"
		if config.Domain != nil {
			domain = *config.Domain
		}

		if !strings.Contains(domain, "*") && hd != domain {
			return tpmodels.TypeUserInfo{}, errors.New("Please use emails from " + domain + " to login")
		}

		tenantID, err := googleWorkspacesProvider.GetTenantID(clientID, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		id = fmt.Sprintf("%s-%s", tenantID, id)

		return tpmodels.TypeUserInfo{
			ThirdPartyUserId: id,
			EmailInfo: &tpmodels.EmailStruct{
				ID:         email,
				IsVerified: isVerified,
			},
			ResponseFromProvider: claims,
		}, nil

	}

	googleWorkspacesProvider.GetConfig = getConfig
	googleWorkspacesProvider.GetAuthorisationRedirectURL = getAuthorisationRedirectURL
	googleWorkspacesProvider.ExchangeAuthCodeForOAuthTokens = exchangeAuthCodeForOAuthTokens
	googleWorkspacesProvider.GetUserInfo = getUserInfo
	googleWorkspacesProvider.GetTenantID = getTenantID

	if input.Override != nil {
		googleWorkspacesProvider = input.Override(googleWorkspacesProvider)
	}

	if len(input.Config) == 0 && (&googleWorkspacesProvider.GetConfig == &getConfig) {
		// no config is provided and GetConfig is not overridden
		return tpmodels.TypeProvider{}, errors.New("please specify a config or override GetConfig")
	}

	return *googleWorkspacesProvider.TypeProvider, nil
}

func verifyAndGetClaims(idToken string, clientId string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	// Get the JWKS URL.
	jwksURL := "https://www.googleapis.com/oauth2/v3/certs"

	// Create the JWKS from the resource at the given URL.
	jwks, err := getJWKSFromURL(jwksURL)
	if err != nil {
		return claims, err
	}

	// Parse the JWT.
	token, err := jwt.ParseWithClaims(idToken, claims, jwks.Keyfunc)
	if err != nil {
		return claims, err
	}

	// Check if the token is valid.
	if !token.Valid {
		return claims, errors.New("invalid id_token supplied")
	}

	if claims["iss"].(string) != "https://accounts.google.com" && claims["iss"].(string) != "accounts.google.com" {
		return claims, errors.New("invalid iss field")
	}

	if claims["aud"].(string) != clientId {
		return claims, errors.New("the client for whom this key is for is different than the one provided")
	}

	return claims, nil
}
