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
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const appleID = "apple"

type TypeAppleInput struct {
	Config   []AppleConfig
	Override func(provider *AppleProvider) *AppleProvider
}

type AppleConfig struct {
	ClientID     string
	ClientSecret AppleClientSecret
	Scope        []string
}

type AppleClientSecret struct {
	KeyId      string
	PrivateKey string
	TeamId     string
}

type AppleProvider struct {
	GetConfig func(id *tpmodels.TypeID, userContext supertokens.UserContext) (AppleConfig, error)
	*tpmodels.TypeProvider
}

func Apple(input TypeAppleInput) tpmodels.TypeProvider {
	appleProvider := &AppleProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: appleID,
		},
	}

	appleProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (AppleConfig, error) {
		if id == nil && len(input.Config) == 0 {
			return AppleConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if id == nil && len(input.Config) > 1 {
			return AppleConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if id == nil {
			return input.Config[0], nil
		}

		if id.Type == tpmodels.TypeClientID {
			for _, config := range input.Config {
				if config.ClientID == id.ID {
					return config, nil
				}
			}
		} else {
			// TODO Multitenant
		}

		return AppleConfig{}, errors.New("config for specified clientID not found")
	}

	customProvider := CustomProvider(TypeCustomProviderInput{
		ThirdPartyID: appleID,
		Override: func(provider *TypeCustomProvider) *TypeCustomProvider {
			provider.GetConfig = func(ID *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
				appleConfig, err := appleProvider.GetConfig(ID, userContext)
				if err != nil {
					return CustomProviderConfig{}, err
				}

				authURL := "https://appleid.apple.com/auth/authorize"
				tokenURL := "https://appleid.apple.com/auth/token"
				userInfoURL := "https://apple.com/api/users/@me"
				jwksURL := "https://appleid.apple.com/auth/keys"

				clientSecret, err := getClientSecret(appleConfig.ClientID, appleConfig.ClientSecret)
				if err != nil {
					return CustomProviderConfig{}, err
				}

				return CustomProviderConfig{
					ClientID:     appleConfig.ClientID,
					ClientSecret: clientSecret,
					Scope:        appleConfig.Scope,

					AuthorizationURL: &authURL,
					AccessTokenURL:   &tokenURL,
					UserInfoURL:      &userInfoURL,
					JwksURL:          &jwksURL,
					DefaultScope:     []string{"email"},

					AuthorizationURLQueryParams: map[string]interface{}{
						"response_mode": "form_post",
						"response_type": "code",
					},

					GetSupertokensUserFromRawResponse: func(rawResponse map[string]interface{}, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
						result := tpmodels.TypeUserInfo{}
						result.ThirdPartyUserId = fmt.Sprint(rawResponse["sub"])
						result.EmailInfo = &tpmodels.EmailStruct{
							ID: fmt.Sprint(rawResponse["email"]),
						}
						emailVerified, emailVerifiedOk := rawResponse["email_verified"].(bool)
						result.EmailInfo.IsVerified = emailVerified && emailVerifiedOk

						result.RawUserInfoFromProvider = rawResponse

						return result, nil
					},
				}, nil
			}

			return provider
		},
	})

	appleProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		return customProvider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
	}

	appleProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		return customProvider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
	}

	appleProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		return customProvider.GetUserInfo(id, oAuthTokens, userContext)
	}

	if input.Override != nil {
		appleProvider = input.Override(appleProvider)
	}

	return customProvider
}

// func Apple(config tpmodels.AppleConfig) tpmodels.TypeProvider {
// 	return tpmodels.TypeProvider{
// 		ID: appleID,
// 		Get: func(redirectURI, authCodeFromRequest *string, userContext supertokens.UserContext) tpmodels.TypeProviderGetResponse {
// 			accessTokenAPIURL := "https://appleid.apple.com/auth/token"
// 			clientSecret, err := getClientSecret(config.ClientID, config.ClientSecret.KeyId, config.ClientSecret.TeamId, config.ClientSecret.PrivateKey)
// 			if err != nil {
// 				panic(err)
// 			}
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
// 			scopes := []string{"email"}
// 			if config.Scope != nil {
// 				scopes = config.Scope
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
// 				GetProfileInfo: func(authCodeResponse interface{}, userContext supertokens.UserContext) (tpmodels.UserInfo, error) {
// 					claims, err := verifyAndGetClaimsAppleIdToken(authCodeResponse.(map[string]interface{})["id_token"].(string), api.GetActualClientIdFromDevelopmentClientId(config.ClientID))
// 					if err != nil {
// 						return tpmodels.UserInfo{}, err
// 					}

// 					var email string
// 					var isVerified bool
// 					var id string
// 					for key, val := range claims {
// 						if key == "sub" {
// 							id = val.(string)
// 						} else if key == "email" {
// 							email = val.(string)
// 						} else if key == "email_verified" {
// 							isVerified = val.(string) == "true"
// 						}
// 					}
// 					return tpmodels.UserInfo{
// 						ID: id,
// 						Email: &tpmodels.EmailStruct{
// 							ID:         email,
// 							IsVerified: isVerified,
// 						},
// 					}, nil
// 				},
// 				GetClientId: func(userContext supertokens.UserContext) string {
// 					return config.ClientID
// 				},
// 				GetRedirectURI: func(userContext supertokens.UserContext) (string, error) {
// 					supertokens, err := supertokens.GetInstanceOrThrowError()
// 					if err != nil {
// 						return "", err
// 					}
// 					return supertokens.AppInfo.APIDomain.GetAsStringDangerous() + supertokens.AppInfo.APIBasePath.GetAsStringDangerous() + "/callback/apple", nil
// 				},
// 			}
// 		},
// 		IsDefault: config.IsDefault,
// 	}
// }

func getClientSecret(clientId string, secret AppleClientSecret) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 86400*180,
		IssuedAt:  time.Now().Unix(),
		Audience:  "https://appleid.apple.com",
		Id:        secret.KeyId,
		Subject:   getActualClientIdFromDevelopmentClientId(clientId),
		Issuer:    secret.TeamId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	ecdsaPrivateKey, err := getECDSPrivateKey(secret.PrivateKey)
	if err != nil {
		return "", err
	}

	// Finally sign the token with the value of type *ecdsa.PrivateKey
	return token.SignedString(ecdsaPrivateKey)
}

func getECDSPrivateKey(privateKey string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	// Check if it's a private key
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	// Get the encoded bytes
	x509Encoded := block.Bytes

	// Now you need an instance of *ecdsa.PrivateKey
	parsedKey, err := x509.ParsePKCS8PrivateKey(x509Encoded) // EDIT to x509Encoded from p8bytes
	if err != nil {
		return nil, err
	}

	ecdsaPrivateKey, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("not ecdsa private key")
	}
	return ecdsaPrivateKey, nil
}

// func verifyAndGetClaimsAppleIdToken(idToken string, clientId string) (jwt.MapClaims, error) {
// 	/*
// 	   - Verify the JWS E256 signature using the server’s public key
// 	   - Verify that the iss field contains https://appleid.apple.com
// 	   - Verify that the aud field is the developer’s client_id
// 	   - Verify that the time is earlier than the exp value of the token */
// 	claims := jwt.MapClaims{}
// 	// Get the JWKS URL.
// 	jwksURL := "https://appleid.apple.com/auth/keys"

// 	// Create the JWKS from the resource at the given URL.
// 	jwks, err := getJWKSFromURL(jwksURL)
// 	if err != nil {
// 		return claims, err
// 	}

// 	// Parse the JWT.
// 	token, err := jwt.ParseWithClaims(idToken, claims, jwks.Keyfunc)
// 	if err != nil {
// 		return claims, err
// 	}

// 	// Check if the token is valid.
// 	if !token.Valid {
// 		return claims, errors.New("invalid id_token supplied")
// 	}

// 	if claims["iss"].(string) != "https://appleid.apple.com" {
// 		return claims, errors.New("invalid iss field in apple token")
// 	}

// 	if claims["aud"].(string) != clientId {
// 		return claims, errors.New("the client for whom this key is for is different than the one provided")
// 	}

// 	return claims, nil
// }
