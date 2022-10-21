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

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}

	TokenEndpoint string
	TokenParams   map[string]interface{}

	UserInfoEndpoint string

	JwksURI      string
	OIDCEndpoint string

	GetSupertokensUserInfoFromRawUserInfoResponse func(rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider, userContext supertokens.UserContext) (tpmodels.TypeSupertokensUserInfo, error)

	AdditionalConfig map[string]interface{}
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

	var customProviderConfig []CustomProviderConfig
	if input.Config != nil {
		customProviderConfig = make([]CustomProviderConfig, len(input.Config))
		for idx, config := range input.Config {
			customProviderConfig[idx] = appleConfigToCustomProviderConfig(config)
		}
	}

	customProvider := customProvider(TypeCustomProviderInput{
		ThirdPartyID: googleID,
		Config:       customProviderConfig,
	})

	{
		// Custom provider needs to use the config returned by apple provider GetConfig
		// Also, apple provider needs to use the default implementation of GetConfig provided by custom provider
		oGetConfigFromCustomProvider := customProvider.GetConfig
		customProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
			config, err := appleProvider.GetConfig(id, userContext)
			if err != nil {
				return CustomProviderConfig{}, err
			}
			return appleConfigToCustomProviderConfig(config), nil
		}
		appleProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (AppleConfig, error) {
			config, err := oGetConfigFromCustomProvider(id, userContext)
			if err != nil {
				return AppleConfig{}, err
			}
			return appleConfigFromCustomProviderConfig(config), nil
		}
	}

	{
		// Apple provider APIs call into custom provider APIs

		appleProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return customProvider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
		}

		appleProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return customProvider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
		}

		appleProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return customProvider.GetUserInfo(id, oAuthTokens, userContext)
		}
	}

	if input.Override != nil {
		appleProvider = input.Override(appleProvider)
	}

	{
		// We want to always normalize (for apple) the config before returning it
		oGetConfig := appleProvider.GetConfig
		appleProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (AppleConfig, error) {
			config, err := oGetConfig(id, userContext)
			if err != nil {
				return AppleConfig{}, err
			}
			return normalizeAppleConfig(config), nil
		}
	}

	return *appleProvider.TypeProvider
}

func normalizeAppleConfig(config AppleConfig) AppleConfig {
	if config.AuthorizationEndpoint == "" {
		config.AuthorizationEndpoint = "https://appleid.apple.com/auth/authorize"
	}

	if config.AuthorizationEndpointQueryParams == nil {
		config.AuthorizationEndpointQueryParams = map[string]interface{}{
			"response_mode": "form_post",
		}
	}

	if config.TokenEndpoint == "" {
		config.TokenEndpoint = "https://appleid.apple.com/auth/token"
	}

	if config.JwksURI == "" {
		config.JwksURI = "https://appleid.apple.com/auth/keys"
	}

	if len(config.Scope) == 0 {
		config.Scope = []string{"email"}
	}

	if config.GetSupertokensUserInfoFromRawUserInfoResponse == nil {
		config.GetSupertokensUserInfoFromRawUserInfoResponse = getSupertokensUserInfoFromRawUserInfo("sub", "email", "email_verified", "id_token")
	}

	return config
}

func getClientSecret(clientId string, secret AppleClientSecret) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 86400*180,
		IssuedAt:  time.Now().Unix(),
		Audience:  "https://appleid.apple.com",
		Subject:   getActualClientIdFromDevelopmentClientId(clientId),
		Issuer:    secret.TeamId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = secret.KeyId
	token.Header["alg"] = "ES256"

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

func appleConfigToCustomProviderConfig(appleConfig AppleConfig) CustomProviderConfig {
	clientSecret, _ := getClientSecret(appleConfig.ClientID, appleConfig.ClientSecret)

	additionalConfig := map[string]interface{}{}

	for k, v := range appleConfig.AdditionalConfig {
		additionalConfig[k] = v
	}
	additionalConfig["_clientSecret"] = appleConfig.ClientSecret

	return CustomProviderConfig{
		ClientID:     appleConfig.ClientID,
		ClientSecret: clientSecret,
		Scope:        appleConfig.Scope,

		AuthorizationEndpoint:            appleConfig.AuthorizationEndpoint,
		AuthorizationEndpointQueryParams: appleConfig.AuthorizationEndpointQueryParams,

		TokenEndpoint: appleConfig.TokenEndpoint,
		TokenParams:   appleConfig.TokenParams,

		UserInfoEndpoint: appleConfig.UserInfoEndpoint,

		JwksURI:      appleConfig.JwksURI,
		OIDCEndpoint: appleConfig.OIDCEndpoint,

		GetSupertokensUserInfoFromRawUserInfoResponse: appleConfig.GetSupertokensUserInfoFromRawUserInfoResponse,

		AdditionalConfig: additionalConfig,
	}
}

func appleConfigFromCustomProviderConfig(config CustomProviderConfig) AppleConfig {
	return AppleConfig{
		ClientID:     config.ClientID,
		ClientSecret: config.AdditionalConfig["_clientSecret"].(AppleClientSecret),

		Scope: config.Scope,

		AuthorizationEndpoint:            config.AuthorizationEndpoint,
		AuthorizationEndpointQueryParams: config.AuthorizationEndpointQueryParams,

		TokenEndpoint: config.TokenEndpoint,
		TokenParams:   config.TokenParams,

		UserInfoEndpoint: config.UserInfoEndpoint,

		JwksURI:      config.JwksURI,
		OIDCEndpoint: config.OIDCEndpoint,

		GetSupertokensUserInfoFromRawUserInfoResponse: config.GetSupertokensUserInfoFromRawUserInfoResponse,

		AdditionalConfig: config.AdditionalConfig,
	}
}
