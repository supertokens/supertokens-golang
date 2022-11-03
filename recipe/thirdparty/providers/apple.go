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
	Config   AppleConfig
	Override func(provider *AppleProvider) *AppleProvider
}

type AppleConfig struct {
	Clients []AppleClientConfig

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	ForcePKCE                        *bool
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      tpmodels.TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, config tpmodels.TypeNormalisedProviderConfig) (bool, error)
}

type AppleClientConfig struct {
	ClientID         string
	ClientSecret     AppleClientSecret
	Scope            []string
	AdditionalConfig map[string]interface{}
}

type AppleClientSecret struct {
	KeyId      string
	PrivateKey string
	TeamId     string
}

type AppleProvider struct {
	*tpmodels.TypeProvider
}

func Apple(input TypeAppleInput) tpmodels.TypeProvider {
	appleProvider := &AppleProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: appleID,
		},
	}

	var oAuth2ProviderClientConfig []OAuth2ProviderClientConfig
	if input.Config.Clients != nil {
		oAuth2ProviderClientConfig = make([]OAuth2ProviderClientConfig, len(input.Config.Clients))
		for idx, client := range input.Config.Clients {
			oAuth2ProviderClientConfig[idx] = appleConfigToOAuth2ProviderConfig(client) // TODO: recompute this after 180 days
		}
	}

	oAuth2Provider := oAuth2Provider(TypeOAuth2ProviderInput{
		ThirdPartyID: googleID,
		Config: OAuth2ProviderConfig{
			Clients:                          oAuth2ProviderClientConfig,
			AuthorizationEndpoint:            input.Config.AuthorizationEndpoint,
			AuthorizationEndpointQueryParams: input.Config.AuthorizationEndpointQueryParams,
			TokenEndpoint:                    input.Config.TokenEndpoint,
			TokenParams:                      input.Config.TokenParams,
			ForcePKCE:                        input.Config.ForcePKCE,
			UserInfoEndpoint:                 input.Config.UserInfoEndpoint,
			JwksURI:                          input.Config.JwksURI,
			OIDCDiscoveryEndpoint:            input.Config.OIDCDiscoveryEndpoint,
			UserInfoMap:                      input.Config.UserInfoMap,
			ValidateIdTokenPayload:           input.Config.ValidateIdTokenPayload,
		},
	})

	{
		// Apple provider APIs call into oAuth2 provider APIs

		appleProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
			return oAuth2Provider.GetConfig(clientType, tenantId, userContext)
		}

		appleProvider.GetAuthorisationRedirectURL = func(config tpmodels.TypeNormalisedProviderConfig, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return oAuth2Provider.GetAuthorisationRedirectURL(config, redirectURIOnProviderDashboard, userContext)
		}

		appleProvider.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.TypeNormalisedProviderConfig, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return oAuth2Provider.ExchangeAuthCodeForOAuthTokens(config, redirectInfo, userContext)
		}

		appleProvider.GetUserInfo = func(config tpmodels.TypeNormalisedProviderConfig, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return oAuth2Provider.GetUserInfo(config, oAuthTokens, userContext)
		}
	}

	if input.Override != nil {
		appleProvider = input.Override(appleProvider)
	}

	{
		// We want to always normalize (for apple) the config before returning it
		oGetConfig := appleProvider.GetConfig
		appleProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
			config, err := oGetConfig(clientType, tenantId, userContext)
			if err != nil {
				return config, err
			}
			return normalizeAppleConfig(config), nil
		}
	}

	return *appleProvider.TypeProvider
}

func normalizeAppleConfig(config tpmodels.TypeNormalisedProviderConfig) tpmodels.TypeNormalisedProviderConfig {
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

	if config.UserInfoMap.From == "" {
		config.UserInfoMap.From = tpmodels.FromIdTokenPayload
	}

	if config.UserInfoMap.IdField == "" {
		config.UserInfoMap.IdField = "sub"
	}

	if config.UserInfoMap.EmailField == "" {
		config.UserInfoMap.EmailField = "email"
	}

	if config.UserInfoMap.EmailVerifiedField == "" {
		config.UserInfoMap.EmailVerifiedField = "email_verified"
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

func appleConfigToOAuth2ProviderConfig(appleConfig AppleClientConfig) OAuth2ProviderClientConfig {
	clientSecret, _ := getClientSecret(appleConfig.ClientID, appleConfig.ClientSecret)

	additionalConfig := map[string]interface{}{}

	for k, v := range appleConfig.AdditionalConfig {
		additionalConfig[k] = v
	}
	additionalConfig["_clientSecret"] = appleConfig.ClientSecret

	return OAuth2ProviderClientConfig{
		ClientID:     appleConfig.ClientID,
		ClientSecret: clientSecret,
		Scope:        appleConfig.Scope,

		AdditionalConfig: additionalConfig,
	}
}
