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
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type APIInterface struct {
	AuthorisationUrlGET      *func(provider *tpmodels.TypeProvider, redirectURIOnProviderDashboard string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error)
	AppleRedirectHandlerPOST *func(formPostInfoFromProvider map[string]interface{}, options tpmodels.APIOptions, userContext supertokens.UserContext) error
	ThirdPartySignInUpPOST   *func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, options tpmodels.APIOptions, userContext supertokens.UserContext) (ThirdPartySignInUpPOSTResponse, error)

	EmailPasswordEmailExistsGET    *func(email string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.EmailExistsGETResponse, error)
	GeneratePasswordResetTokenPOST *func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.GeneratePasswordResetTokenPOSTResponse, error)
	PasswordResetPOST              *func(formFields []epmodels.TypeFormField, token string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.ResetPasswordPOSTResponse, error)
	EmailPasswordSignInPOST        *func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (SignInPOSTResponse, error)
	EmailPasswordSignUpPOST        *func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (SignUpPOSTResponse, error)
}

type SignUpPOSTResponse struct {
	OK *struct {
		User    User
		Session sessmodels.SessionContainer
	}
	EmailAlreadyExistsError *struct{}
	GeneralError            *supertokens.GeneralErrorResponse
}

type SignInPOSTResponse struct {
	OK *struct {
		User    User
		Session sessmodels.SessionContainer
	}
	WrongCredentialsError *struct{}
	GeneralError          *supertokens.GeneralErrorResponse
}

type EmailpasswordInput struct {
	IsSignIn   bool
	FormFields []epmodels.TypeFormField
	Options    epmodels.APIOptions
}

type EmailpasswordOutput struct {
	OK *struct {
		User           User
		CreatedNewUser bool
	}
	EmailAlreadyExistsError *struct{}
	WrongCredentialsError   *struct{}
}

type ThirdPartySignInUpPOSTResponse struct {
	OK *struct {
		CreatedNewUser          bool
		User                    User
		Session                 sessmodels.SessionContainer
		OAuthTokens             tpmodels.TypeOAuthTokens
		RawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider
	}
	NoEmailGivenByProviderError *struct{}
	GeneralError                *supertokens.GeneralErrorResponse
}
