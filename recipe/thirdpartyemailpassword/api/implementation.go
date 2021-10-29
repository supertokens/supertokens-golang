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

package api

import (
	epapi "github.com/supertokens/supertokens-golang/recipe/emailpassword/api"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	tpapi "github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
)

func MakeAPIImplementation() tpepmodels.APIInterface {
	emailPasswordImplementation := epapi.MakeAPIImplementation()
	thirdPartyImplementation := tpapi.MakeAPIImplementation()

	emailExistsGET := func(email string, options epmodels.APIOptions) (epmodels.EmailExistsGETResponse, error) {
		return (*emailPasswordImplementation.EmailExistsGET)(email, options)

	}

	generatePasswordResetTokenPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.GeneratePasswordResetTokenPOSTResponse, error) {
		return (*emailPasswordImplementation.GeneratePasswordResetTokenPOST)(formFields, options)
	}

	passwordResetPOST := func(formFields []epmodels.TypeFormField, token string, options epmodels.APIOptions) (epmodels.ResetPasswordUsingTokenResponse, error) {
		return (*emailPasswordImplementation.PasswordResetPOST)(formFields, token, options)
	}

	emailPasswordSignInPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (tpepmodels.SignInResponse, error) {
		response, err := (*emailPasswordImplementation.SignInPOST)(formFields, options)
		if err != nil {
			return tpepmodels.SignInResponse{}, err
		}
		if response.OK != nil {
			return tpepmodels.SignInResponse{
				OK: &struct {
					User tpepmodels.User
				}{
					User: tpepmodels.User{
						ID:         response.OK.User.ID,
						Email:      response.OK.User.Email,
						TimeJoined: response.OK.User.TimeJoined,
						ThirdParty: nil,
					},
				},
			}, nil
		} else {
			return tpepmodels.SignInResponse{
				WrongCredentialsError: &struct{}{},
			}, nil
		}
	}

	emailPasswordSignUpPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (tpepmodels.SignUpResponse, error) {
		response, err := (*emailPasswordImplementation.SignUpPOST)(formFields, options)
		if err != nil {
			return tpepmodels.SignUpResponse{}, err
		}
		if response.OK != nil {
			return tpepmodels.SignUpResponse{
				OK: &struct {
					User tpepmodels.User
				}{
					User: tpepmodels.User{
						ID:         response.OK.User.ID,
						Email:      response.OK.User.Email,
						TimeJoined: response.OK.User.TimeJoined,
						ThirdParty: nil,
					},
				},
			}, nil
		} else {
			return tpepmodels.SignUpResponse{
				EmailAlreadyExistsError: &struct{}{},
			}, nil
		}
	}

	thirdPartySignInUpPOST := func(provider tpmodels.TypeProvider, code, redirectURI string, options tpmodels.APIOptions) (tpepmodels.ThirdPartyOutput, error) {
		response, err := (*thirdPartyImplementation.SignInUpPOST)(provider, code, redirectURI, options)
		if err != nil {
			return tpepmodels.ThirdPartyOutput{}, err
		}
		if response.FieldError != nil {
			return tpepmodels.ThirdPartyOutput{
				FieldError: &struct{ Error string }{},
			}, nil
		} else if response.NoEmailGivenByProviderError != nil {
			return tpepmodels.ThirdPartyOutput{
				NoEmailGivenByProviderError: &struct{}{},
			}, nil
		} else {
			return tpepmodels.ThirdPartyOutput{
				OK: &struct {
					CreatedNewUser   bool
					User             tpepmodels.User
					AuthCodeResponse interface{}
				}{
					CreatedNewUser:   response.OK.CreatedNewUser,
					AuthCodeResponse: response.OK.AuthCodeResponse,
					User: tpepmodels.User{
						ID:         response.OK.User.ID,
						TimeJoined: response.OK.User.TimeJoined,
						Email:      response.OK.User.Email,
						ThirdParty: &response.OK.User.ThirdParty,
					},
				},
			}, nil
		}
	}

	authorisationUrlGET := func(provider tpmodels.TypeProvider, options tpmodels.APIOptions) (tpmodels.AuthorisationUrlGETResponse, error) {
		return (*thirdPartyImplementation.AuthorisationUrlGET)(provider, options)
	}

	return tpepmodels.APIInterface{
		AuthorisationUrlGET:            &authorisationUrlGET,
		EmailExistsGET:                 &emailExistsGET,
		GeneratePasswordResetTokenPOST: &generatePasswordResetTokenPOST,
		PasswordResetPOST:              &passwordResetPOST,
		ThirdPartySignInUpPOST:         &thirdPartySignInUpPOST,
		EmailPasswordSignInPOST:        &emailPasswordSignInPOST,
		EmailPasswordSignUpPOST:        &emailPasswordSignUpPOST,
	}
}
