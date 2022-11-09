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
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	tpapi "github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() tpepmodels.APIInterface {
	emailPasswordImplementation := epapi.MakeAPIImplementation()
	thirdPartyImplementation := tpapi.MakeAPIImplementation()

	ogEmailExistsGET := *emailPasswordImplementation.EmailExistsGET
	emailExistsGET := func(email string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.EmailExistsGETResponse, error) {
		return ogEmailExistsGET(email, options, userContext)

	}

	ogGeneratePasswordResetTokenPOST := *emailPasswordImplementation.GeneratePasswordResetTokenPOST
	generatePasswordResetTokenPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.GeneratePasswordResetTokenPOSTResponse, error) {
		return ogGeneratePasswordResetTokenPOST(formFields, options, userContext)
	}

	ogPasswordResetPOST := *emailPasswordImplementation.PasswordResetPOST
	passwordResetPOST := func(formFields []epmodels.TypeFormField, token string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.ResetPasswordPOSTResponse, error) {
		return ogPasswordResetPOST(formFields, token, options, userContext)
	}

	ogSignInPOST := *emailPasswordImplementation.SignInPOST
	emailPasswordSignInPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (tpepmodels.SignInPOSTResponse, error) {
		response, err := ogSignInPOST(formFields, options, userContext)
		if err != nil {
			return tpepmodels.SignInPOSTResponse{}, err
		}
		if response.OK != nil {
			return tpepmodels.SignInPOSTResponse{
				OK: &struct {
					User    tpepmodels.User
					Session sessmodels.SessionContainer
				}{
					User: tpepmodels.User{
						ID:         response.OK.User.ID,
						Email:      response.OK.User.Email,
						TimeJoined: response.OK.User.TimeJoined,
						ThirdParty: nil,
					},
					Session: response.OK.Session,
				},
			}, nil
		} else if response.WrongCredentialsError != nil {
			return tpepmodels.SignInPOSTResponse{
				WrongCredentialsError: &struct{}{},
			}, nil
		} else {
			return tpepmodels.SignInPOSTResponse{
				GeneralError: response.GeneralError,
			}, nil
		}
	}

	ogSignUpPOST := *emailPasswordImplementation.SignUpPOST
	emailPasswordSignUpPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (tpepmodels.SignUpPOSTResponse, error) {
		response, err := ogSignUpPOST(formFields, options, userContext)
		if err != nil {
			return tpepmodels.SignUpPOSTResponse{}, err
		}
		if response.OK != nil {
			return tpepmodels.SignUpPOSTResponse{
				OK: &struct {
					User    tpepmodels.User
					Session sessmodels.SessionContainer
				}{
					User: tpepmodels.User{
						ID:         response.OK.User.ID,
						Email:      response.OK.User.Email,
						TimeJoined: response.OK.User.TimeJoined,
						ThirdParty: nil,
					},
					Session: response.OK.Session,
				},
			}, nil
		} else if response.EmailAlreadyExistsError != nil {
			return tpepmodels.SignUpPOSTResponse{
				EmailAlreadyExistsError: &struct{}{},
			}, nil
		} else {
			return tpepmodels.SignUpPOSTResponse{
				GeneralError: response.GeneralError,
			}, nil
		}
	}

	ogSignInUpPOST := *thirdPartyImplementation.SignInUpPOST
	thirdPartySignInUpPOST := func(provider tpmodels.TypeProvider, code string, authCodeResponse interface{}, redirectURI string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpepmodels.ThirdPartyOutput, error) {
		response, err := ogSignInUpPOST(provider, code, authCodeResponse, redirectURI, options, userContext)
		if err != nil {
			return tpepmodels.ThirdPartyOutput{}, err
		}
		if response.GeneralError != nil {
			return tpepmodels.ThirdPartyOutput{
				GeneralError: response.GeneralError,
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
					Session          sessmodels.SessionContainer
				}{
					CreatedNewUser:   response.OK.CreatedNewUser,
					AuthCodeResponse: response.OK.AuthCodeResponse,
					User: tpepmodels.User{
						ID:         response.OK.User.ID,
						TimeJoined: response.OK.User.TimeJoined,
						Email:      response.OK.User.Email,
						ThirdParty: &response.OK.User.ThirdParty,
					},
					Session: response.OK.Session,
				},
			}, nil
		}
	}

	ogAuthorisationUrlGET := *thirdPartyImplementation.AuthorisationUrlGET
	authorisationUrlGET := func(provider tpmodels.TypeProvider, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error) {
		return ogAuthorisationUrlGET(provider, options, userContext)
	}

	ogAppleRedirectHandlerPOST := *thirdPartyImplementation.AppleRedirectHandlerPOST
	appleRedirectHandlerPOST := func(code string, state string, options tpmodels.APIOptions, userContext supertokens.UserContext) error {
		return ogAppleRedirectHandlerPOST(code, state, options, userContext)
	}
	result := tpepmodels.APIInterface{
		AuthorisationUrlGET:            &authorisationUrlGET,
		EmailPasswordEmailExistsGET:    &emailExistsGET,
		GeneratePasswordResetTokenPOST: &generatePasswordResetTokenPOST,
		PasswordResetPOST:              &passwordResetPOST,
		ThirdPartySignInUpPOST:         &thirdPartySignInUpPOST,
		EmailPasswordSignInPOST:        &emailPasswordSignInPOST,
		EmailPasswordSignUpPOST:        &emailPasswordSignUpPOST,
		AppleRedirectHandlerPOST:       &appleRedirectHandlerPOST,
	}

	modifiedEP := GetEmailPasswordIterfaceImpl(result)
	(*emailPasswordImplementation.EmailExistsGET) = *modifiedEP.EmailExistsGET
	(*emailPasswordImplementation.GeneratePasswordResetTokenPOST) = *modifiedEP.GeneratePasswordResetTokenPOST
	(*emailPasswordImplementation.PasswordResetPOST) = *modifiedEP.PasswordResetPOST
	(*emailPasswordImplementation.SignInPOST) = *modifiedEP.SignInPOST
	(*emailPasswordImplementation.SignUpPOST) = *modifiedEP.SignUpPOST

	modifiedTP := GetThirdPartyIterfaceImpl(result)
	(*thirdPartyImplementation.AuthorisationUrlGET) = *modifiedTP.AuthorisationUrlGET
	(*thirdPartyImplementation.SignInUpPOST) = *modifiedTP.SignInUpPOST
	(*thirdPartyImplementation.AppleRedirectHandlerPOST) = *modifiedTP.AppleRedirectHandlerPOST

	return result
}
