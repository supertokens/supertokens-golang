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
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const discordID = "discord"

func Discord(config tpmodels.DiscordConfig) tpmodels.TypeProvider {
	return tpmodels.TypeProvider{
		ID: discordID,
		Get: func(redirectURI, authCodeFromRequest *string, userContext supertokens.UserContext) tpmodels.TypeProviderGetResponse {
			accessTokenAPIURL := "https://discord.com/api/oauth2/token"
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

			authorisationRedirectURL := "https://discord.com/api/oauth2/authorize"
			scopes := []string{"email", "identify"}
			if config.Scope != nil {
				scopes = config.Scope
			}

			var additionalParams map[string]interface{} = nil
			if config.AuthorisationRedirect != nil && config.AuthorisationRedirect.Params != nil {
				additionalParams = config.AuthorisationRedirect.Params
			}

			authorizationRedirectParams := map[string]interface{}{
				"scope":         strings.Join(scopes, " "),
				"client_id":     config.ClientID,
				"response_type": "code",
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

					accessToken := authCodeResponse.(map[string]interface{})["access_token"].(string)
					authHeader := "Bearer " + accessToken
					response, err := getAuthRequest(authHeader)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					userInfo := response.(map[string]interface{})
					_, emailOk := userInfo["email"]
					if !emailOk {
						return tpmodels.UserInfo{
							ID:    userInfo["id"].(string),
							Email: nil,
						}, nil
					}
					return tpmodels.UserInfo{
						ID: userInfo["id"].(string),
						Email: &tpmodels.EmailStruct{
							ID:         userInfo["email"].(string),
							IsVerified: userInfo["verified"].(bool),
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

func getAuthRequest(authHeader string) (interface{}, error) {
	url := "https://discord.com/api/users/@me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	return doGetRequest(req)
}
