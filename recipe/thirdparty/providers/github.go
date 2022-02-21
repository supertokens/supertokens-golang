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
	"fmt"
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const githubID = "github"

func Github(config tpmodels.GithubConfig) tpmodels.TypeProvider {
	return tpmodels.TypeProvider{
		ID: githubID,
		Get: func(redirectURI, authCodeFromRequest *string, userContext supertokens.UserContext) tpmodels.TypeProviderGetResponse {
			accessTokenAPIURL := "https://github.com/login/oauth/access_token"
			accessTokenAPIParams := map[string]string{
				"client_id":     config.ClientID,
				"client_secret": config.ClientSecret,
			}
			if authCodeFromRequest != nil {
				accessTokenAPIParams["code"] = *authCodeFromRequest
			}
			if redirectURI != nil {
				accessTokenAPIParams["redirect_uri"] = *redirectURI
			}

			authorisationRedirectURL := "https://github.com/login/oauth/authorize"
			scopes := []string{"read:user", "user:email"}
			if config.Scope != nil {
				scopes = config.Scope
			}

			var additionalParams map[string]interface{} = nil
			if config.AuthorisationRedirect != nil && config.AuthorisationRedirect.Params != nil {
				additionalParams = config.AuthorisationRedirect.Params
			}

			authorizationRedirectParams := map[string]interface{}{
				"scope":     strings.Join(scopes, " "),
				"client_id": config.ClientID,
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
					var accessTokenAPIResponse githubGetProfileInfoInput
					err = json.Unmarshal(authCodeResponseJson, &accessTokenAPIResponse)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					accessToken := accessTokenAPIResponse.AccessToken
					authHeader := "Bearer " + accessToken
					response, err := getGithubAuthRequest(authHeader)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					userInfo := response.(map[string]interface{})
					emailsInfoResponse, err := getGithubEmailsInfo(authHeader)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					emailsInfo := emailsInfoResponse.([]interface{})
					ID := fmt.Sprintf("%f", userInfo["id"].(float64)) // github userId will be a number
					// if user has choosen not to show their email publicly, userInfo here will
					// have email as null. So we instead get the info from the emails api and
					// use the email which is maked as primary one.
					var emailInfo map[string]interface{}
					for _, info := range emailsInfo {
						emailInfoMap := info.(map[string]interface{})
						if emailInfoMap["primary"].(bool) {
							emailInfo = emailInfoMap
							break
						}
					}
					if emailInfo == nil {
						return tpmodels.UserInfo{
							ID: ID,
						}, nil
					}
					isVerified := false
					if emailInfo != nil {
						isVerified = emailInfo["verified"].(bool)
					}
					return tpmodels.UserInfo{
						ID: ID,
						Email: &tpmodels.EmailStruct{
							ID:         emailInfo["email"].(string),
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

func getGithubAuthRequest(authHeader string) (interface{}, error) {
	url := "https://api.github.com/user"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	return doGetRequest(req)
}

func getGithubEmailsInfo(authHeader string) (interface{}, error) {
	url := "https://api.github.com/user/emails"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	return doGetRequest(req)
}

type githubGetProfileInfoInput struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
