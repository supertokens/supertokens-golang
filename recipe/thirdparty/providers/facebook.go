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
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const facebookID = "facebook"

type TypeFacebookInput struct {
	Config   []FacebookConfig
	Override func(provider *FacebookProvider) *FacebookProvider
}

type FacebookConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type FacebookProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (FacebookConfig, error)
	*tpmodels.TypeProvider
}

// func Facebook(config tpmodels.FacebookConfig) tpmodels.TypeProvider {
// 	return tpmodels.TypeProvider{
// 		ID: facebookID,
// 		Get: func(redirectURI, authCodeFromRequest *string, userContext supertokens.UserContext) tpmodels.TypeProviderGetResponse {
// 			accessTokenAPIURL := "https://graph.facebook.com/v9.0/oauth/access_token"
// 			accessTokenAPIParams := map[string]string{
// 				"client_id":     config.ClientID,
// 				"client_secret": config.ClientSecret,
// 			}
// 			if authCodeFromRequest != nil {
// 				accessTokenAPIParams["code"] = *authCodeFromRequest
// 			}
// 			if redirectURI != nil {
// 				accessTokenAPIParams["redirect_uri"] = *redirectURI
// 			}

// 			authorisationRedirectURL := "https://www.facebook.com/v9.0/dialog/oauth"
// 			scopes := []string{"email"}
// 			if config.Scope != nil {
// 				scopes = config.Scope
// 			}

// 			authorizationRedirectParams := map[string]interface{}{
// 				"scope":         strings.Join(scopes, " "),
// 				"response_type": "code",
// 				"client_id":     config.ClientID,
// 			}

// 			return tpmodels.TypeProviderGetResponse{
// 				AccessTokenAPI: tpmodels.AccessTokenAPI{
// 					URL:    accessTokenAPIURL,
// 					Params: accessTokenAPIParams,
// 				},
// 				AuthorisationRedirect: tpmodels.AuthorisationRedirect{
// 					URL:    authorisationRedirectURL,
// 					Params: authorizationRedirectParams,
// 				},
// 				GetProfileInfo: func(authCodeResponse interface{}, userContext supertokens.UserContext) (tpmodels.UserInfo, error) {
// 					authCodeResponseJson, err := json.Marshal(authCodeResponse)
// 					if err != nil {
// 						return tpmodels.UserInfo{}, err
// 					}
// 					var accessTokenAPIResponse facebookGetProfileInfoInput
// 					err = json.Unmarshal(authCodeResponseJson, &accessTokenAPIResponse)
// 					if err != nil {
// 						return tpmodels.UserInfo{}, err
// 					}
// 					accessToken := accessTokenAPIResponse.AccessToken
// 					response, err := getFacebookAuthRequest(accessToken)
// 					if err != nil {
// 						return tpmodels.UserInfo{}, err
// 					}
// 					userInfo := response.(map[string]interface{})
// 					ID := userInfo["id"].(string)
// 					email, emailOk := userInfo["email"].(string)
// 					if !emailOk {
// 						return tpmodels.UserInfo{
// 							ID: ID,
// 						}, nil
// 					}
// 					isVerified, isVerifiedOk := userInfo["verified_email"].(bool)
// 					return tpmodels.UserInfo{
// 						ID: ID,
// 						Email: &tpmodels.EmailStruct{
// 							ID:         email,
// 							IsVerified: isVerified && isVerifiedOk,
// 						},
// 					}, nil
// 				},
// 				GetClientId: func(userContext supertokens.UserContext) string {
// 					return config.ClientID
// 				},
// 			}
// 		},
// 		IsDefault: config.IsDefault,
// 	}
// }

// func getFacebookAuthRequest(accessToken string) (interface{}, error) {
// 	url := "https://graph.facebook.com/me"
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	q := req.URL.Query()
// 	q.Add("access_token", accessToken)
// 	q.Add("fields", "id,email")
// 	q.Add("format", "json")
// 	req.URL.RawQuery = q.Encode()
// 	return doGetRequest(req)
// }

// type facebookGetProfileInfoInput struct {
// 	AccessToken string `json:"access_token"`
// 	ExpiresIn   int    `json:"expires_in"`
// 	TokenType   string `json:"token_type"`
// }
