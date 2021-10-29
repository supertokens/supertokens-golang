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
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

type APIInterface struct {
	AuthorisationUrlGET            *func(provider tpmodels.TypeProvider, options tpmodels.APIOptions) (tpmodels.AuthorisationUrlGETResponse, error)
	EmailExistsGET                 *func(email string, options epmodels.APIOptions) (epmodels.EmailExistsGETResponse, error)
	GeneratePasswordResetTokenPOST *func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.GeneratePasswordResetTokenPOSTResponse, error)
	PasswordResetPOST              *func(formFields []epmodels.TypeFormField, token string, options epmodels.APIOptions) (epmodels.ResetPasswordUsingTokenResponse, error)
	ThirdPartySignInUpPOST         *func(provider tpmodels.TypeProvider, code string, redirectURI string, options tpmodels.APIOptions) (ThirdPartyOutput, error)
	EmailPasswordSignInPOST        *func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (SignInResponse, error)
	EmailPasswordSignUpPOST        *func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (SignUpResponse, error)
}

type EmailpasswordInput struct {
	IsSignIn   bool
	FormFields []epmodels.TypeFormField
	Options    epmodels.APIOptions
}

type SignInUpAPIOutput struct {
	EmailpasswordOutput *EmailpasswordOutput
	ThirdPartyOutput    *ThirdPartyOutput
}

type EmailpasswordOutput struct {
	OK *struct {
		User           User
		CreatedNewUser bool
	}
	EmailAlreadyExistsError *struct{}
	WrongCredentialsError   *struct{}
}

type ThirdPartyOutput struct {
	OK *struct {
		CreatedNewUser   bool
		User             User
		AuthCodeResponse interface{}
	}
	NoEmailGivenByProviderError *struct{}
	FieldError                  *struct{ Error string }
}
