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
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/derekstavis/go-qs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/supertokens"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

const appleID = "apple"

func Apple(input tpmodels.TypeAppleInput) tpmodels.TypeProvider {
	appleProvider := &tpmodels.AppleProvider{}

	getConfig := func(clientID *string, userContext supertokens.UserContext) (tpmodels.AppleConfig, error) {
		if input.Config == nil || len(input.Config) == 0 {
			return tpmodels.AppleConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if clientID == nil && len(input.Config) > 1 {
			return tpmodels.AppleConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if clientID == nil && len(input.Config) == 1 {
			return input.Config[0], nil
		}

		for _, config := range input.Config {
			if config.ClientID == *clientID {
				return config, nil
			}
		}

		return tpmodels.AppleConfig{}, errors.New("config for specified clientID not found")
	}

	getAuthorisationRedirectURL := func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		config, err := appleProvider.GetConfig(clientID, userContext)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		scopes := []string{"email"}
		if config.Scope != nil {
			scopes = config.Scope
		}

		url := "https://appleid.apple.com/auth/authorize"
		queryParams := map[string]interface{}{
			"scope":         strings.Join(scopes, " "),
			"response_mode": "form_post",
			"response_type": "code",
			"client_id":     getActualClientIdFromDevelopmentClientId(config.ClientID),
		}

		queryParams["redirect_uri"] = redirectURIOnProviderDashboard

		url, queryParams, err = getAuthRedirectForDev(config.ClientID, url, queryParams)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}

		queryParamsStr, err := qs.Marshal(queryParams)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}

		return tpmodels.TypeAuthorisationRedirect{
			URLWithQueryParams: url + "?" + queryParamsStr,
		}, nil
	}

	exchangeAuthCodeForOAuthTokens := func(clientID *string, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		config, err := appleProvider.GetConfig(clientID, userContext)
		if err != nil {
			return nil, err
		}

		clientSecret, err := getClientSecret(config.ClientID, config.ClientSecret.KeyId, config.ClientSecret.TeamId, config.ClientSecret.PrivateKey)
		if err != nil {
			return nil, err
		}
		accessTokenAPIURL := "https://appleid.apple.com/auth/token"
		accessTokenAPIParams := map[string]string{
			"client_id":     getActualClientIdFromDevelopmentClientId(config.ClientID),
			"client_secret": clientSecret,
			"grant_type":    "authorization_code",
			"code":          redirectURIInfo.RedirectURIQueryParams["code"].(string),
		}

		redirectURI := checkDevAndGetRedirectURI(
			config.ClientID,
			redirectURIInfo.RedirectURIOnProviderDashboard,
			userContext,
		)

		accessTokenAPIParams["redirect_uri"] = redirectURI

		authResponseFromRequest, err := postRequest(accessTokenAPIURL, accessTokenAPIParams)
		if err != nil {
			return nil, err
		}

		authResponse := tpmodels.TypeOAuthTokens{}

		for k, v := range authResponseFromRequest {
			authResponse[k] = v
		}

		return authResponse, nil
	}

	getUserInfo := func(clientID *string, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		config, err := appleProvider.GetConfig(clientID, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		claims, err := verifyAndGetClaimsAppleIdToken(
			oAuthTokens["id_token"].(string),
			getActualClientIdFromDevelopmentClientId(config.ClientID),
		)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		var email string
		var isVerified bool
		var id string
		for key, val := range claims {
			if key == "sub" {
				id = val.(string)
			} else if key == "email" {
				email = val.(string)
			} else if key == "email_verified" {
				isVerified = val.(string) == "true"
			}
		}
		userInfo := tpmodels.TypeUserInfo{
			ThirdPartyUserId: id,
			EmailInfo: &tpmodels.EmailStruct{
				ID:         email,
				IsVerified: isVerified,
			},
			ResponseFromProvider: claims,
		}
		return userInfo, nil
	}

	appleProvider.GetConfig = getConfig
	appleProvider.GetAuthorisationRedirectURL = getAuthorisationRedirectURL
	appleProvider.ExchangeAuthCodeForOAuthTokens = exchangeAuthCodeForOAuthTokens
	appleProvider.GetUserInfo = getUserInfo

	if input.Override != nil {
		appleProvider = input.Override(appleProvider)
	}

	return tpmodels.TypeProvider{
		ID: appleID,

		GetAuthorisationRedirectURL: func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return appleProvider.GetAuthorisationRedirectURL(clientID, redirectURIOnProviderDashboard, userContext)
		},

		ExchangeAuthCodeForOAuthTokens: func(clientID *string, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return appleProvider.ExchangeAuthCodeForOAuthTokens(clientID, redirectInfo, userContext)
		},

		GetUserInfo: func(clientID *string, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return appleProvider.GetUserInfo(clientID, oAuthTokens, userContext)
		},
	}
}

func getClientSecret(clientId, keyId, teamId, privateKey string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 86400*180,
		IssuedAt:  time.Now().Unix(),
		Audience:  "https://appleid.apple.com",
		Id:        keyId,
		Subject:   getActualClientIdFromDevelopmentClientId(clientId),
		Issuer:    teamId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	ecdsaPrivateKey, err := getECDSPrivateKey(privateKey)
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

func verifyAndGetClaimsAppleIdToken(idToken string, clientId string) (jwt.MapClaims, error) {
	/*
	   - Verify the JWS E256 signature using the server’s public key
	   - Verify that the iss field contains https://appleid.apple.com
	   - Verify that the aud field is the developer’s client_id
	   - Verify that the time is earlier than the exp value of the token */
	claims := jwt.MapClaims{}
	// Get the JWKS URL.
	jwksURL := "https://appleid.apple.com/auth/keys"

	// Create the keyfunc options. Refresh the JWKS every hour and log errors.
	refreshInterval := time.Hour
	options := keyfunc.Options{
		RefreshInterval: refreshInterval,
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		return claims, err
	}

	// Parse the JWT.
	token, err := jwt.ParseWithClaims(idToken, claims, jwks.Keyfunc)
	if err != nil {
		return claims, err
	}

	// Check if the token is valid.
	if !token.Valid {
		return claims, errors.New("invalid id_token supplied")
	}

	if claims["iss"].(string) != "https://appleid.apple.com" {
		return claims, errors.New("invalid iss field in apple token")
	}

	if claims["aud"].(string) != clientId {
		return claims, errors.New("the client for whom this key is for is different than the one provided")
	}

	return claims, nil
}
