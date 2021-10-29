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

type RecipeInterface struct {
	SignUp                   *func(email string, password string) (SignUpResponse, error)
	SignIn                   *func(email string, password string) (SignInResponse, error)
	GetUserByID              *func(userID string) (*User, error)
	GetUserByEmail           *func(email string) (*User, error)
	CreateResetPasswordToken *func(userID string) (CreateResetPasswordTokenResponse, error)
	ResetPasswordUsingToken  *func(token string, newPassword string) (ResetPasswordUsingTokenResponse, error)
	UpdateEmailOrPassword    *func(userId string, email *string, password *string) (UpdateEmailOrPasswordResponse, error)
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

type CreateResetPasswordTokenResponse struct {
	OK *struct {
		Token string
	}
	UnknownUserIdError *struct{}
}

type ResetPasswordUsingTokenResponse struct {
	OK                             *struct{}
	ResetPasswordInvalidTokenError *struct{}
}

type UpdateEmailOrPasswordResponse struct {
	OK                      *struct{}
	UnknownUserIdError      *struct{}
	EmailAlreadyExistsError *struct{}
}
