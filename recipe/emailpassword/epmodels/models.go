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

package epmodels

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeNormalisedInput struct {
	SignUpFeature                  TypeNormalisedInputSignUp
	SignInFeature                  TypeNormalisedInputSignIn
	ResetPasswordUsingTokenFeature TypeNormalisedInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       evmodels.TypeInput
	Override                       OverrideStruct
	GetEmailDeliveryConfig         func(recipeImpl RecipeInterface) emaildelivery.TypeInputWithService
}

type OverrideStruct struct {
	Functions                func(originalImplementation RecipeInterface) RecipeInterface
	APIs                     func(originalImplementation APIInterface) APIInterface
	EmailVerificationFeature *evmodels.OverrideStruct
}

type TypeInputEmailVerificationFeature struct {
	GetEmailVerificationURL  func(user User, userContext supertokens.UserContext) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string, userContext supertokens.UserContext) // Deprecated: Use EmailDelivery instead.
}

type TypeInputFormField struct {
	ID       string
	Validate func(value interface{}) *string
	Optional *bool
}

type TypeInputSignUp struct {
	FormFields []TypeInputFormField
}

type NormalisedFormField struct {
	ID       string
	Validate func(value interface{}) *string
	Optional bool
}

type TypeNormalisedInputSignUp struct {
	FormFields []NormalisedFormField
}

type TypeNormalisedInputSignIn struct {
	FormFields []NormalisedFormField
}

type TypeInputResetPasswordUsingTokenFeature struct {
	GetResetPasswordURL      func(user User, userContext supertokens.UserContext) (string, error)
	CreateAndSendCustomEmail func(user User, passwordResetURLWithToken string, userContext supertokens.UserContext) // Deprecated: Use EmailDelivery instead.
}

type TypeNormalisedInputResetPasswordUsingTokenFeature struct {
	GetResetPasswordURL            func(user User, userContext supertokens.UserContext) (string, error)
	FormFieldsForGenerateTokenForm []NormalisedFormField
	FormFieldsForPasswordResetForm []NormalisedFormField
}

type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	TimeJoined uint64 `json:"timejoined"`
}

type TypeInput struct {
	SignUpFeature                  *TypeInputSignUp
	ResetPasswordUsingTokenFeature *TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       *TypeInputEmailVerificationFeature
	Override                       *OverrideStruct
	EmailDelivery                  *emaildelivery.TypeInput
}

type TypeFormField struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
