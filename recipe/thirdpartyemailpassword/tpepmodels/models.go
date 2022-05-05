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

package tpepmodels

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type User struct {
	ID         string `json:"id"`
	TimeJoined uint64 `json:"timeJoined"`
	Email      string `json:"email"`
	ThirdParty *struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
	} `json:"thirdParty"`
}

type TypeContext struct {
	FormFields                 []epmodels.TypeFormField
	ThirdPartyAuthCodeResponse interface{}
}

type TypeInputEmailVerificationFeature struct {
	GetEmailVerificationURL  func(user User, userContext supertokens.UserContext) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string, userContext supertokens.UserContext)
}

type TypeInput struct {
	SignUpFeature                  *epmodels.TypeInputSignUp
	Providers                      []tpmodels.TypeProvider
	ResetPasswordUsingTokenFeature *epmodels.TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       *TypeInputEmailVerificationFeature
	Override                       *OverrideStruct
}

type TypeNormalisedInput struct {
	SignUpFeature                  *epmodels.TypeInputSignUp
	Providers                      []tpmodels.TypeProvider
	ResetPasswordUsingTokenFeature *epmodels.TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       evmodels.TypeInput
	Override                       OverrideStruct
}

type OverrideStruct struct {
	Functions                func(originalImplementation RecipeInterface) RecipeInterface
	APIs                     func(originalImplementation APIInterface) APIInterface
	EmailVerificationFeature *evmodels.OverrideStruct
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}
