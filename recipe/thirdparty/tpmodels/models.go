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
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
)

type UserInfo struct {
	ID    string
	Email *EmailStruct
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}

type TypeProviderGetResponse struct {
	AccessTokenAPI        AccessTokenAPI
	AuthorisationRedirect AuthorisationRedirect
	GetProfileInfo        func(authCodeResponse interface{}) (UserInfo, error)
}

type AccessTokenAPI struct {
	URL    string
	Params map[string]string
}

type AuthorisationRedirect struct {
	URL    string
	Params map[string]interface{}
}

type TypeProvider struct {
	ID  string
	Get func(redirectURI *string, authCodeFromRequest *string) TypeProviderGetResponse
}

type User struct {
	ID         string
	TimeJoined uint64
	Email      string
	ThirdParty struct {
		ID     string
		UserID string
	}
}

type TypeInputEmailVerificationFeature struct {
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
}

type TypeInputSignInAndUp struct {
	Providers []TypeProvider
}

type TypeNormalisedInputSignInAndUp struct {
	Providers []TypeProvider
}

type TypeInput struct {
	SignInAndUpFeature       TypeInputSignInAndUp
	EmailVerificationFeature *TypeInputEmailVerificationFeature
	Override                 *OverrideStruct
}

type TypeNormalisedInput struct {
	SignInAndUpFeature       TypeNormalisedInputSignInAndUp
	EmailVerificationFeature evmodels.TypeInput
	Override                 OverrideStruct
}

type OverrideStruct struct {
	Functions                func(originalImplementation RecipeInterface) RecipeInterface
	APIs                     func(originalImplementation APIInterface) APIInterface
	EmailVerificationFeature *evmodels.OverrideStruct
}
