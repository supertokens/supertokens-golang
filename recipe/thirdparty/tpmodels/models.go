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
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
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
	GetProfileInfo        func(authCodeResponse interface{}, userContext supertokens.UserContext) (UserInfo, error)
	GetClientId           func(userContext supertokens.UserContext) string
	GetRedirectURI        func(userContext supertokens.UserContext) (string, error)
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
	ID        string
	Get       func(redirectURI *string, authCodeFromRequest *string, userContext supertokens.UserContext) TypeProviderGetResponse
	IsDefault bool
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

type TypeInputEmailVerificationFeature struct {
	GetEmailVerificationURL  func(user User, userContext supertokens.UserContext) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string, userContext supertokens.UserContext) // Deprecated: Use EmailDelivery instead.
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
	EmailDelivery            *emaildelivery.TypeInput
}

type TypeNormalisedInput struct {
	SignInAndUpFeature       TypeNormalisedInputSignInAndUp
	EmailVerificationFeature evmodels.TypeInput
	Override                 OverrideStruct
	GetEmailDeliveryConfig   func(recipeImpl RecipeInterface) emaildelivery.TypeInputWithService
}

type OverrideStruct struct {
	Functions                func(originalImplementation RecipeInterface) RecipeInterface
	APIs                     func(originalImplementation APIInterface) APIInterface
	EmailVerificationFeature *evmodels.OverrideStruct
}
