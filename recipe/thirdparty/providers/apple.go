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

					GetSupertokensUserInfoFromRawUserInfoResponse: func(rawUserInfoResponse map[string]interface{}, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
						result := tpmodels.TypeUserInfo{}
						result.ThirdPartyUserId = fmt.Sprint(rawUserInfoResponse["sub"])
						result.EmailInfo = &tpmodels.EmailStruct{
							ID: fmt.Sprint(rawUserInfoResponse["email"]),
						}
						emailVerified, emailVerifiedOk := rawUserInfoResponse["email_verified"].(bool)
						result.EmailInfo.IsVerified = emailVerified && emailVerifiedOk

						result.RawUserInfoFromProvider = rawUserInfoResponse

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

	return *appleProvider.TypeProvider
}

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
