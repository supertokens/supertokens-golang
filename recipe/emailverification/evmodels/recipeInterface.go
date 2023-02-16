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

package evmodels

import "github.com/supertokens/supertokens-golang/supertokens"

type RecipeInterface struct {
	CreateEmailVerificationToken  *func(userID, email string, tenantId *string, userContext supertokens.UserContext) (CreateEmailVerificationTokenResponse, error)
	VerifyEmailUsingToken         *func(token string, tenantId *string, userContext supertokens.UserContext) (VerifyEmailUsingTokenResponse, error)
	IsEmailVerified               *func(userID, email string, tenantId *string, userContext supertokens.UserContext) (bool, error)
	RevokeEmailVerificationTokens *func(userId, email string, tenantId *string, userContext supertokens.UserContext) (RevokeEmailVerificationTokensResponse, error)
	UnverifyEmail                 *func(userId, email string, tenantId *string, userContext supertokens.UserContext) (UnverifyEmailResponse, error)
}

type CreateEmailVerificationTokenResponse struct {
	OK *struct {
		Token string
	}
	EmailAlreadyVerifiedError *struct{}
}

type VerifyEmailUsingTokenResponse struct {
	OK *struct {
		User User
	}
	EmailVerificationInvalidTokenError *struct{}
}

type RevokeEmailVerificationTokensResponse struct {
	OK *struct{}
}

type UnverifyEmailResponse struct {
	OK *struct{}
}
