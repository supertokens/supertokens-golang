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

import "github.com/supertokens/supertokens-golang/supertokens"

type RecipeInterface struct {
	SignUp                   *func(email string, password string, tenantId string, userContext supertokens.UserContext) (SignUpResponse, error)
	SignIn                   *func(email string, password string, tenantId string, userContext supertokens.UserContext) (SignInResponse, error)
	GetUserByID              *func(userID string, userContext supertokens.UserContext) (*User, error)
	GetUserByEmail           *func(email string, tenantId string, userContext supertokens.UserContext) (*User, error)
	CreateResetPasswordToken *func(userID string, tenantId string, userContext supertokens.UserContext) (CreateResetPasswordTokenResponse, error)
	ResetPasswordUsingToken  *func(token string, newPassword string, tenantId string, userContext supertokens.UserContext) (ResetPasswordUsingTokenResponse, error)
	UpdateEmailOrPassword    *func(userId string, email *string, password *string, applyPasswordPolicy *bool, tenantIdForPasswordPolicy string, userContext supertokens.UserContext) (UpdateEmailOrPasswordResponse, error)
}

type SignUpResponse struct {
	OK *struct {
		User User
	}
	EmailAlreadyExistsError *struct{}
}

func (resp SignUpResponse) ToJsonableMap() map[string]interface{} {
	if resp.OK != nil {
		return map[string]interface{}{
			"status": "OK",
			"user":   resp.OK.User,
		}
	} else {
		return map[string]interface{}{
			"status": "EMAIL_ALREADY_EXISTS_ERROR",
		}
	}
}

type SignInResponse struct {
	OK *struct {
		User User
	}
	WrongCredentialsError *struct{}
}

func (resp SignInResponse) ToJsonableMap() map[string]interface{} {
	if resp.OK != nil {
		return map[string]interface{}{
			"status": "OK",
			"user":   resp.OK.User,
		}
	} else {
		return map[string]interface{}{
			"status": "WRONG_CREDENTIALS_ERROR",
		}
	}
}

type CreateResetPasswordTokenResponse struct {
	OK *struct {
		Token string
	}
	UnknownUserIdError *struct{}
}

func (resp CreateResetPasswordTokenResponse) ToJsonableMap() map[string]interface{} {
	if resp.OK != nil {
		return map[string]interface{}{
			"status": "OK",
			"token":  resp.OK.Token,
		}
	} else {
		return map[string]interface{}{
			"status": "UNKNOWN_USER_ID_ERROR",
		}
	}
}

type ResetPasswordUsingTokenResponse struct {
	OK *struct {
		UserId *string
	}
	ResetPasswordInvalidTokenError *struct{}
}

func (resp ResetPasswordUsingTokenResponse) ToJsonableMap() map[string]interface{} {
	if resp.OK != nil {
		return map[string]interface{}{
			"status": "OK",
			"userId": resp.OK.UserId,
		}
	} else {
		return map[string]interface{}{
			"status": "RESET_PASSWORD_INVALID_TOKEN_ERROR",
		}
	}
}

type UpdateEmailOrPasswordResponse struct {
	OK                          *struct{}
	UnknownUserIdError          *struct{}
	EmailAlreadyExistsError     *struct{}
	PasswordPolicyViolatedError *PasswordPolicyViolatedError
}

func (resp UpdateEmailOrPasswordResponse) ToJsonableMap() map[string]interface{} {
	if resp.OK != nil {
		return map[string]interface{}{
			"status": "OK",
		}
	} else if resp.UnknownUserIdError != nil {
		return map[string]interface{}{
			"status": "UNKNOWN_USER_ID_ERROR",
		}
	} else if resp.EmailAlreadyExistsError != nil {
		return map[string]interface{}{
			"status": "EMAIL_ALREADY_EXISTS_ERROR",
		}
	} else if resp.PasswordPolicyViolatedError != nil {
		return map[string]interface{}{
			"status":        "PASSWORD_POLICY_VIOLATED_ERROR",
			"failureReason": resp.PasswordPolicyViolatedError.FailureReason,
		}
	}
	return map[string]interface{}{
		"status": "UNKNOWN_ERROR",
	}
}

type PasswordPolicyViolatedError struct {
	FailureReason string
}

type CreateResetPasswordLinkResponse struct {
	OK *struct {
		Link string
	}
	UnknownUserIdError *struct{}
}

func (resp CreateResetPasswordLinkResponse) ToJsonableMap() map[string]interface{} {
	if resp.OK != nil {
		return map[string]interface{}{
			"status": "OK",
			"link":   resp.OK.Link,
		}
	} else if resp.UnknownUserIdError != nil {
		return map[string]interface{}{
			"status": "UNKNOWN_USER_ID_ERROR",
		}
	}
	return map[string]interface{}{
		"status": "UNKNOWN_ERROR",
	}
}

type SendResetPasswordEmailResponse struct {
	OK                 *struct{}
	UnknownUserIdError *struct{}
}

func (resp SendResetPasswordEmailResponse) ToJsonableMap() map[string]interface{} {
	if resp.OK != nil {
		return map[string]interface{}{
			"status": "OK",
		}
	} else if resp.UnknownUserIdError != nil {
		return map[string]interface{}{
			"status": "UNKNOWN_USER_ID_ERROR",
		}
	}
	return map[string]interface{}{
		"status": "UNKNOWN_ERROR",
	}
}
