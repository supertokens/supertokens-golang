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
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const bitbucketID = "bitbucket"

func Bitbucket(config tpmodels.BitbucketConfig) tpmodels.TypeProvider {
	return tpmodels.TypeProvider{
		ID: bitbucketID,
		Get: func(redirectURI, authCodeFromRequest *string, userContext supertokens.UserContext) tpmodels.TypeProviderGetResponse {
			accessTokenAPIURL := "https://bitbucket.org/site/oauth2/access_token"
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

			authorisationRedirectURL := "https://bitbucket.org/site/oauth2/authorize"
			scopes := []string{"account", "email"}
			if config.Scope != nil {
				scopes = config.Scope
			}

			var additionalParams map[string]interface{} = nil
			if config.AuthorisationRedirect != nil && config.AuthorisationRedirect.Params != nil {
				additionalParams = config.AuthorisationRedirect.Params
			}

			authorizationRedirectParams := map[string]interface{}{
				"scope":         strings.Join(scopes, " "),
				"access_type":   "offline",
				"response_type": "code",
				"client_id":     config.ClientID,
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
					authCodeResponseJson, err := json.Marshal(authCodeResponse)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					var accessTokenAPIResponse bitbucketGetProfileInfoInput
					err = json.Unmarshal(authCodeResponseJson, &accessTokenAPIResponse)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					accessToken := accessTokenAPIResponse.AccessToken
					authHeader := "Bearer " + accessToken
					response, err := getBitbucketAuthRequest(authHeader)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					userInfo := response.(map[string]interface{})
					ID := userInfo["uuid"].(string)

					emailResponse, err := getBitbucketEmailRequest(authHeader)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					var email string
					var isVerified bool = false
					emailResponseInfo := emailResponse.(map[string]interface{})
					for _, emailInfo := range emailResponseInfo["values"].([]interface{}) {
						emailInfoMap := emailInfo.(map[string]interface{})
						if emailInfoMap["is_primary"].(bool) {
							email = emailInfoMap["email"].(string)
							isVerified = emailInfoMap["is_confirmed"].(bool)
						}
					}
					if email == "" {
						return tpmodels.UserInfo{
							ID: ID,
						}, nil
					}
					return tpmodels.UserInfo{
						ID: ID,
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

func getBitbucketAuthRequest(authHeader string) (interface{}, error) {
	url := "https://api.bitbucket.org/2.0/user"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	return doGetRequest(req)
}

func getBitbucketEmailRequest(authHeader string) (interface{}, error) {
	url := "https://api.bitbucket.org/2.0/user/emails"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	return doGetRequest(req)
}

type bitbucketGetProfileInfoInput struct {
	AccessToken string `json:"access_token"`
}
