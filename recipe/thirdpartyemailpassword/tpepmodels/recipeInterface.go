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
	"github.com/supertokens/supertokens-golang/supertokens"
)

type RecipeInterface struct {
	GetUserByID              *func(userID string, userContext supertokens.UserContext) (*User, error)
	GetUsersByEmail          *func(email string, userContext supertokens.UserContext) ([]User, error)
	GetUserByThirdPartyInfo  *func(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*User, error)
	ThirdPartySignInUp       *func(thirdPartyID string, thirdPartyUserID string, email string, userContext supertokens.UserContext) (SignInUpResponse, error)
	EmailPasswordSignUp      *func(email string, password string, userContext supertokens.UserContext) (SignUpResponse, error)
	EmailPasswordSignIn      *func(email string, password string, userContext supertokens.UserContext) (SignInResponse, error)
	CreateResetPasswordToken *func(userID string, userContext supertokens.UserContext) (epmodels.CreateResetPasswordTokenResponse, error)
	ResetPasswordUsingToken  *func(token string, newPassword string, userContext supertokens.UserContext) (epmodels.ResetPasswordUsingTokenResponse, error)
	UpdateEmailOrPassword    *func(userId string, email *string, password *string, applyPasswordPolicy *bool, userContext supertokens.UserContext) (epmodels.UpdateEmailOrPasswordResponse, error)
}

type SignInUpResponse struct {
	OK *struct {
		CreatedNewUser bool
		User           User
	}
}

type SignUpResponse struct {
	OK *struct {
		User User
	}
	EmailAlreadyExistsError *struct{}
}

type SignInResponse struct {
	OK *struct {
		User User
	}
	WrongCredentialsError *struct{}
}
