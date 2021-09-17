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

// import (
// 	"encoding/json"
// 	"strings"

// 	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
// )

// const appleID = "apple"

// func Apple(config tpmodels.AppleConfig) tpmodels.TypeProvider {
// 	return tpmodels.TypeProvider{
// 		ID: appleID,
// 		Get: func(redirectURI, authCodeFromRequest *string) tpmodels.TypeProviderGetResponse {
// 			accessTokenAPIURL := "https://appleid.apple.com/auth/token"
// 			clientSecret, _ := getClientSecret(config.ClientID, config.ClientSecret.KeyId, config.ClientSecret.TeamId, config.ClientSecret.PrivateKey)
// 			accessTokenAPIParams := map[string]string{
// 				"client_id":     config.ClientID,
// 				"client_secret": clientSecret,
// 				"grant_type":    "authorization_code",
// 			}
// 			if authCodeFromRequest != nil {
// 				accessTokenAPIParams["code"] = *authCodeFromRequest
// 			}
// 			if redirectURI != nil {
// 				accessTokenAPIParams["redirect_uri"] = *redirectURI
// 			}

// 			authorisationRedirectURL := "https://appleid.apple.com/auth/authorize"
// 			scopes := []string{"name", "email"}
// 			if config.Scope != nil {
// 				scopes = append(scopes, config.Scope...)
// 			}

// 			var additionalParams map[string]interface{} = nil
// 			if config.AuthorisationRedirect != nil && config.AuthorisationRedirect.Params != nil {
// 				additionalParams = config.AuthorisationRedirect.Params
// 			}

// 			authorizationRedirectParams := map[string]interface{}{
// 				"scope":         strings.Join(scopes, " "),
// 				"response_mode": "form_post",
// 				"response_type": "code",
// 				"client_id":     config.ClientID,
// 			}
// 			for key, value := range additionalParams {
// 				authorizationRedirectParams[key] = value
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
// 				GetProfileInfo: func(authCodeResponse interface{}) (tpmodels.UserInfo, error) {
// 					authCodeResponseJson, err := json.Marshal(authCodeResponse)
// 					if err != nil {
// 						return tpmodels.UserInfo{}, err
// 					}
// 					var accessTokenAPIResponse appleGetProfileInfoInput
// 					err = json.Unmarshal(authCodeResponseJson, &accessTokenAPIResponse)
// 					if err != nil {
// 						return tpmodels.UserInfo{}, err
// 					}
// 					return tpmodels.UserInfo{}, nil
// 				},
// 			}
// 		},
// 	}
// }

// func getClientSecret(clientId, keyId, teamId, privateKey string) (string, error) {
// 	return "", nil
// }

// type appleGetProfileInfoInput struct {
// 	AccessToken  string `json:"access_token"`
// 	ExpiresIn    int    `json:"expires_in"`
// 	TokenType    string `json:"token_type"`
// 	RefreshToken string `json:"refresh_token"`
// 	IDToken      string `json:"id_token"`
// }
