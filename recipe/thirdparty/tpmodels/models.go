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

package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeRedirectURIQueryParams = map[string]interface{}
type TypeOAuthTokens = map[string]interface{}

type TypeRawUserInfoFromProvider struct {
	FromIdTokenPayload map[string]interface{}
	FromUserInfoAPI    map[string]interface{}
}

type TypeUserInfo struct {
	ThirdPartyUserId        string
	Email                   *EmailStruct
	RawUserInfoFromProvider TypeRawUserInfoFromProvider
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}

type TypeAuthorisationRedirect struct {
	URLWithQueryParams string
	PKCECodeVerifier   *string
}

type TypeRedirectURIInfo struct {
	RedirectURIOnProviderDashboard string                     `json:"redirectURIOnProviderDashboard"`
	RedirectURIQueryParams         TypeRedirectURIQueryParams `json:"redirectURIQueryParams"`
	PKCECodeVerifier               *string                    `json:"pkceCodeVerifier"`
}

type TypeFrom string

const (
	FromIdTokenPayload TypeFrom = "idTokenPayload"
	FromUserInfoAPI    TypeFrom = "userInfoAPI"
)

type TypeUserInfoMapFields struct {
	UserId        string `json:"userId,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified string `json:"emailVerified,omitempty"`
}

type TypeUserInfoMap struct {
	FromIdTokenPayload TypeUserInfoMapFields `json:"fromIdTokenPayload,omitempty"`
	FromUserInfoAPI    TypeUserInfoMapFields `json:"fromUserInfoAPI,omitempty"`
}

type User struct {
	ID         string `json:"id"`
	TimeJoined uint64 `json:"timeJoined"`
	Email      string `json:"email"`
	ThirdParty struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
	} `json:"thirdParty"`
	TenantIds []string `json:"tenantIds"`
}

type TypeInputSignInAndUp struct {
	Providers []ProviderInput
}

type TypeNormalisedInputSignInAndUp struct {
	Providers []ProviderInput
}

type TypeInput struct {
	SignInAndUpFeature TypeInputSignInAndUp
	Override           *OverrideStruct
}

type TypeNormalisedInput struct {
	SignInAndUpFeature TypeNormalisedInputSignInAndUp
	Override           OverrideStruct
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}

type ProviderInput struct {
	Config   ProviderConfig
	Override func(originalImplementation *TypeProvider) *TypeProvider
}

type ProviderConfig struct {
	ThirdPartyId string `json:"thirdPartyId"`
	Name         string `json:"name,omitempty"`

	Clients []ProviderClientConfig `json:"clients,omitempty"`

	// Fields below are optional for built-in providers
	AuthorizationEndpoint            string                 `json:"authorizationEndpoint,omitempty"`
	AuthorizationEndpointQueryParams map[string]interface{} `json:"authorizationEndpointQueryParams,omitempty"`
	TokenEndpoint                    string                 `json:"tokenEndpoint,omitempty"`
	TokenEndpointBodyParams          map[string]interface{} `json:"tokenEndpointBodyParams,omitempty"`
	UserInfoEndpoint                 string                 `json:"userInfoEndpoint,omitempty"`
	UserInfoEndpointQueryParams      map[string]interface{} `json:"userInfoEndpointQueryParams,omitempty"`
	UserInfoEndpointHeaders          map[string]interface{} `json:"userInfoEndpointHeaders,omitempty"`
	JwksURI                          string                 `json:"jwksURI,omitempty"`
	OIDCDiscoveryEndpoint            string                 `json:"oidcDiscoveryEndpoint,omitempty"`
	UserInfoMap                      TypeUserInfoMap        `json:"userInfoMap,omitempty"`
	RequireEmail                     *bool                  `json:"requireEmail,omitempty"`

	ValidateIdTokenPayload func(idTokenPayload map[string]interface{}, clientConfig ProviderConfigForClientType, userContext supertokens.UserContext) error `json:"-"`

	/**
	 * This function is responsible for validating the access token received from the third party provider.
	 * This check can include checking the expiry of the access token, checking the audience of the access token, etc.
	 *
	 * This function should return an error if the access token should be considered invalid, or return nothing if it is valid
	 */
	ValidateAccessToken func(accessToken string, clientConfig ProviderConfigForClientType, userContext supertokens.UserContext) error `json:"-"`
	GenerateFakeEmail   func(thirdPartyUserId string, tenantId string, userContext supertokens.UserContext) string                    `json:"-"`
}

type ProviderClientConfig struct {
	ClientType       string                 `json:"clientType,omitempty"` // optional
	ClientID         string                 `json:"clientId"`
	ClientSecret     string                 `json:"clientSecret,omitempty"`
	Scope            []string               `json:"scope,omitempty"`
	ForcePKCE        *bool                  `json:"forcePKCE,omitempty"`
	AdditionalConfig map[string]interface{} `json:"additionalConfig,omitempty"`
}

type ProviderConfigForClientType struct {
	Name string

	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenEndpointBodyParams          map[string]interface{}
	ForcePKCE                        *bool // Providers like twitter expects PKCE to be used along with secret
	UserInfoEndpoint                 string
	UserInfoEndpointQueryParams      map[string]interface{}
	UserInfoEndpointHeaders          map[string]interface{}
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig ProviderConfigForClientType, userContext supertokens.UserContext) error
	ValidateAccessToken              func(accessToken string, clientConfig ProviderConfigForClientType, userContext supertokens.UserContext) error `json:"-"`

	RequireEmail      *bool
	GenerateFakeEmail func(thirdPartyUserId string, tenantId string, userContext supertokens.UserContext) string
}

type TypeProvider struct {
	ID     string
	Config ProviderConfigForClientType

	GetConfigForClientType         func(clientType *string, userContext supertokens.UserContext) (ProviderConfigForClientType, error)
	GetAuthorisationRedirectURL    func(redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(redirectURIInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}
