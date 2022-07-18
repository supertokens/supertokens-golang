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

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const googleWorkspacesID = "google-workspaces"

func GoogleWorkspaces(config tpmodels.GoogleWorkspacesConfig) tpmodels.TypeProvider {

	domain := "*"
	if config.Domain != nil {
		domain = *config.Domain
	}

	return tpmodels.TypeProvider{
		ID: googleWorkspacesID,
		Get: func(redirectURI, authCodeFromRequest *string, userContext supertokens.UserContext) tpmodels.TypeProviderGetResponse {
			accessTokenAPIURL := "https://accounts.google.com/o/oauth2/token"
			accessTokenAPIParams := map[string]string{
				"client_id":     config.ClientID,
				"client_secret": config.ClientSecret,
				"grant_type":    "authorization_code",
			}
			if authCodeFromRequest != nil {
				accessTokenAPIParams["code"] = *authCodeFromRequest
			}
			if redirectURI != nil {
				accessTokenAPIParams["redirect_uri"] = *redirectURI
			}

			authorisationRedirectURL := "https://accounts.google.com/o/oauth2/v2/auth"
			scopes := []string{"https://www.googleapis.com/auth/userinfo.email"}
			if config.Scope != nil {
				scopes = config.Scope
			}

			var additionalParams map[string]interface{} = nil
			if config.AuthorisationRedirect != nil && config.AuthorisationRedirect.Params != nil {
				additionalParams = config.AuthorisationRedirect.Params
			}

			authorizationRedirectParams := map[string]interface{}{
				"scope":                  strings.Join(scopes, " "),
				"access_type":            "offline",
				"include_granted_scopes": "true",
				"response_type":          "code",
				"client_id":              config.ClientID,
				"hd":                     domain,
			}
			for key, value := range additionalParams {
				authorizationRedirectParams[key] = value
			}

			return tpmodels.TypeProviderGetResponse{
				AccessTokenAPI: tpmodels.AccessTokenAPI{
					URL:    accessTokenAPIURL,
					Params: accessTokenAPIParams,
				},
				AuthorisationRedirect: tpmodels.AuthorisationRedirect{
					URL:    authorisationRedirectURL,
					Params: authorizationRedirectParams,
				},
				GetProfileInfo: func(authCodeResponse interface{}, userContext supertokens.UserContext) (tpmodels.UserInfo, error) {
					claims, err := verifyAndGetClaims(authCodeResponse.(map[string]interface{})["id_token"].(string), api.GetActualClientIdFromDevelopmentClientId(config.ClientID))
					if err != nil {
						return tpmodels.UserInfo{}, err
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
						return tpmodels.UserInfo{}, errors.New("Could not get email. Please use a different login method")
					}

					if hd == "" {
						return tpmodels.UserInfo{}, errors.New("Please use a Google Workspace ID to login")
					}

					if !strings.Contains(domain, "*") && hd != domain {
						return tpmodels.UserInfo{}, errors.New("Please use emails from " + domain + " to login")
					}

					return tpmodels.UserInfo{
						ID: id,
						Email: &tpmodels.EmailStruct{
							ID:         email,
							IsVerified: isVerified,
						},
					}, nil
				},
				GetClientId: func(userContext supertokens.UserContext) string {
					return config.ClientID
				},
			}
		},
		IsDefault: config.IsDefault,
	}
}

func verifyAndGetClaims(idToken string, clientId string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	// Get the JWKS URL.
	jwksURL := "https://www.googleapis.com/oauth2/v3/certs"

	// Create the keyfunc options. Refresh the JWKS every hour and log errors.
	options := keyfunc.Options{
		// https://github.com/supertokens/supertokens-golang/issues/155
		// This causes a leak as the pointer to JWKS would be held in the goroutine and
		// also results in compounding refresh requests
		// RefreshInterval: time.Hour,
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.Get(jwksURL, options)
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
