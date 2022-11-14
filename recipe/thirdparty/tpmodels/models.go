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

type TypeSupertokensUserInfo struct {
	ThirdPartyUserId string
	EmailInfo        *EmailStruct
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

type TypeUserInfoMap struct {
	FromIdTokenPayload struct {
		UserId        string `json:"userId"`
		Email         string `json:"email"`
		EmailVerified string `json:"emailVerified"`
	} `json:"fromIdTokenPayload"`
	FromUserInfoAPI struct {
		UserId        string `json:"userId"`
		Email         string `json:"email"`
		EmailVerified string `json:"emailVerified"`
	} `json:"fromUserInfoAPI"`
}

type User struct {
	ID         string `json:"id"`
	TimeJoined uint64 `json:"timeJoined"`
	Email      string `json:"email"`
	ThirdParty struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
	} `json:"thirdParty"`
}

type TypeInputSignInAndUp struct {
	Providers []TypeProvider
}

type TypeNormalisedInputSignInAndUp struct {
	Providers            []TypeProvider
	GetUserPoolForTenant func(tenantId string, userContext supertokens.UserContext) (string, error)
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

type TenantConfig struct {
	Clients []TenantClientConfig `json:"clients"`

	// Fields below are optional for built-in providers
	AuthorizationEndpoint            string                 `json:"authorizationEndpoint,omitempty"`
	AuthorizationEndpointQueryParams map[string]interface{} `json:"authorizationEndpointQueryParams,omitempty"`
	TokenEndpoint                    string                 `json:"tokenEndpoint,omitempty"`
	TokenParams                      map[string]interface{} `json:"tokenParams,omitempty"`
	ForcePKCE                        bool                   `json:"forcePKCE,omitempty"`
	UserInfoEndpoint                 string                 `json:"userInfoEndpoint,omitempty"`
	JwksURI                          string                 `json:"jwksURI,omitempty"`
	OIDCDiscoveryEndpoint            string                 `json:"oidcDiscoveryEndpoint,omitempty"`
	UserInfoMap                      TypeUserInfoMap        `json:"userInfoMap,omitempty"`

	FrontendInfo struct {
		Name string `json:"name"`
	} `json:"frontendInfo"`
}

type TenantClientConfig struct {
	ClientType       string                 // optional
	ClientID         string                 `json:"clientId"`
	ClientSecret     string                 `json:"clientSecret"`
	Scope            []string               `json:"scope"`
	AdditionalConfig map[string]interface{} `json:"additionalConfig"`
}
