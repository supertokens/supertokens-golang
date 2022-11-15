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
	"encoding/json"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const appleID = "apple"

func Apple(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	if input.ThirdPartyID == "" {
		input.ThirdPartyID = appleID
	}

	if input.Config.OIDCDiscoveryEndpoint == "" {
		input.Config.OIDCDiscoveryEndpoint = "https://appleid.apple.com/"
	}

	if input.Config.UserInfoMap.FromIdTokenPayload.UserId == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.UserId = "sub"
	}
	if input.Config.UserInfoMap.FromIdTokenPayload.Email == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.Email = "email"
	}

	if input.Config.UserInfoMap.FromIdTokenPayload.EmailVerified == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.EmailVerified = "email_verified"
	}

	if input.Config.AuthorizationEndpointQueryParams == nil {
		input.Config.AuthorizationEndpointQueryParams = map[string]interface{}{}
	}

	if input.Config.AuthorizationEndpointQueryParams["response_mode"] == nil {
		input.Config.AuthorizationEndpointQueryParams["response_mode"] = "form_post"
	}

	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfigForClientType
		provider.GetConfigForClientType = func(clientType *string, input tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, input, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"openid", "email"}
			}

			if config.ClientSecret == "" {
				clientSecret, err := getClientSecret(config.ClientID, config.AdditionalConfig)
				if err != nil {
					return tpmodels.ProviderConfigForClientType{}, err
				}
				config.ClientSecret = clientSecret
			}

			return config, err
		}

		oExchangeAuthCodeForOAuthTokens := provider.ExchangeAuthCodeForOAuthTokens
		provider.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.ProviderConfigForClientType, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			res, err := oExchangeAuthCodeForOAuthTokens(config, redirectURIInfo, userContext)
			if err != nil {
				return tpmodels.TypeOAuthTokens{}, err
			}

			if user, ok := redirectURIInfo.RedirectURIQueryParams["user"].(string); ok {
				userInfo := map[string]interface{}{}
				err := json.Unmarshal([]byte(user), &userInfo)
				if err != nil {
					res["user"] = user
				} else {
					res["user"] = userInfo
				}
			}

			return res, nil
		}

		oGetUserInfo := provider.GetUserInfo
		provider.GetUserInfo = func(config tpmodels.ProviderConfigForClientType, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			res, err := oGetUserInfo(config, oAuthTokens, userContext)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}

			if user, ok := oAuthTokens["user"]; ok {
				res.RawUserInfoFromProvider.FromIdTokenPayload["user"] = user
			}

			return res, nil
		}

		if oOverride != nil {
			provider = oOverride(provider)
		}
		return provider
	}

	return NewProvider(input)
}

func getClientSecret(clientId string, additionalConfig map[string]interface{}) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 86400*180,
		IssuedAt:  time.Now().Unix(),
		Audience:  "https://appleid.apple.com",
		Subject:   getActualClientIdFromDevelopmentClientId(clientId),
		Issuer:    additionalConfig["teamId"].(string),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = additionalConfig["keyId"].(string)
	token.Header["alg"] = "ES256"

	ecdsaPrivateKey, err := getECDSPrivateKey(additionalConfig["privateKey"].(string))
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
