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

	"github.com/golang-jwt/jwt/v5"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Apple(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.Name == "" {
		input.Config.Name = "Apple"
	}

	if input.Config.OIDCDiscoveryEndpoint == "" {
		input.Config.OIDCDiscoveryEndpoint = "https://appleid.apple.com/.well-known/openid-configuration"
	}

	if input.Config.AuthorizationEndpointQueryParams == nil {
		input.Config.AuthorizationEndpointQueryParams = map[string]interface{}{}
	}

	if _, ok := input.Config.AuthorizationEndpointQueryParams["response_mode"]; !ok {
		input.Config.AuthorizationEndpointQueryParams["response_mode"] = "form_post"
	}

	oOverride := input.Override

	input.Override = func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := originalImplementation.GetConfigForClientType
		originalImplementation.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"openid", "email"}
			}

			if config.ClientSecret == "" {

				if config.AdditionalConfig == nil || config.AdditionalConfig["teamId"] == nil || config.AdditionalConfig["keyId"] == nil || config.AdditionalConfig["privateKey"] == nil {
					return tpmodels.ProviderConfigForClientType{}, errors.New("please ensure that keyId, teamId and privateKey are provided in the AdditionalConfig")
				}

				clientSecret, err := getClientSecret(config.ClientID, config.AdditionalConfig)
				if err != nil {
					return tpmodels.ProviderConfigForClientType{}, err
				}
				config.ClientSecret = clientSecret
			}

			// The config could be coming from core where we didn't add the well-known previously
			config.OIDCDiscoveryEndpoint = normaliseOIDCEndpointToIncludeWellKnown(config.OIDCDiscoveryEndpoint)

			return config, nil
		}

		oExchangeAuthCodeForOAuthTokens := originalImplementation.ExchangeAuthCodeForOAuthTokens
		originalImplementation.ExchangeAuthCodeForOAuthTokens = func(redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			res, err := oExchangeAuthCodeForOAuthTokens(redirectURIInfo, userContext)
			if err != nil {
				return tpmodels.TypeOAuthTokens{}, err
			}

			if user, ok := redirectURIInfo.RedirectURIQueryParams["user"].(string); ok {
				userInfo := map[string]interface{}{}
				err := json.Unmarshal([]byte(user), &userInfo)
				if err != nil {
					return nil, err
				}
				res["user"] = userInfo

			} else if userInfo, ok := redirectURIInfo.RedirectURIQueryParams["user"].(map[string]interface{}); ok {
				res["user"] = userInfo
			}

			return res, nil
		}

		oGetUserInfo := originalImplementation.GetUserInfo
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			res, err := oGetUserInfo(oAuthTokens, userContext)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}

			if user, ok := oAuthTokens["user"].(string); ok {
				userInfo := map[string]interface{}{}
				err := json.Unmarshal([]byte(user), &userInfo)
				if err != nil {
					return tpmodels.TypeUserInfo{}, err
				}
				res.RawUserInfoFromProvider.FromIdTokenPayload["user"] = userInfo
			} else if userInfo, ok := oAuthTokens["user"].(map[string]interface{}); ok {
				res.RawUserInfoFromProvider.FromIdTokenPayload["user"] = userInfo
			}

			return res, nil
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}
		return originalImplementation
	}

	return NewProvider(input)
}

func getClientSecret(clientId string, additionalConfig map[string]interface{}) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Audience:  jwt.ClaimStrings{"https://appleid.apple.com"},
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
