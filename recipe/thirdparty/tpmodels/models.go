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
	FromIdToken     map[string]interface{}
	FromAccessToken map[string]interface{}
}

type TypeSupertokensUserInfo struct {
	ThirdPartyUserId string
	EmailInfo        *EmailStruct
}

type TypeUserInfo struct {
	ThirdPartyUserId        string
	EmailInfo               *EmailStruct
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

type TypeProvider struct {
	ID string

	GetAuthorisationRedirectURL    func(id *TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(id *TypeID, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(id *TypeID, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}

type TypeRedirectURIInfo struct {
	RedirectURIOnProviderDashboard string                     `json:"redirectURIOnProviderDashboard"`
	RedirectURIQueryParams         TypeRedirectURIQueryParams `json:"redirectURIQueryParams"`
	PKCECodeVerifier               *string                    `json:"pkceCodeVerifier"`
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
	Providers []TypeProvider
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

type IDTypeEnum string

const (
	TypeClientID IDTypeEnum = "CLIENT_ID"
	TypeTenantID IDTypeEnum = "TENANT_ID"
)

type TypeID struct {
	ID   string
	Type IDTypeEnum
}
