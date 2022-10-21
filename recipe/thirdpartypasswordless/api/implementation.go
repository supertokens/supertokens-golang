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
	plessapi "github.com/supertokens/supertokens-golang/recipe/passwordless/api"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	tpapi "github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() tplmodels.APIInterface {
	passwordlessImplementation := plessapi.MakeAPIImplementation()
	thirdPartyImplementation := tpapi.MakeAPIImplementation()

	ogSignInUpPOST := *thirdPartyImplementation.SignInUpPOST
	thirdPartySignInUpPOST := func(provider tpmodels.TypeProvider, clientID *string, input tpmodels.TypeSignInUpInput, options tpmodels.APIOptions, userContext supertokens.UserContext) (tplmodels.ThirdPartySignInUpOutput, error) {
		response, err := ogSignInUpPOST(provider, clientID, input, options, userContext)
		if err != nil {
			return tplmodels.ThirdPartySignInUpOutput{}, err
		}
		if response.GeneralError != nil {
			return tplmodels.ThirdPartySignInUpOutput{
				GeneralError: response.GeneralError,
			}, nil
		} else if response.NoEmailGivenByProviderError != nil {
			return tplmodels.ThirdPartySignInUpOutput{
				NoEmailGivenByProviderError: &struct{}{},
			}, nil
		} else {
			return tplmodels.ThirdPartySignInUpOutput{
				OK: &struct {
					CreatedNewUser          bool
					User                    tplmodels.User
					Session                 *sessmodels.TypeSessionContainer
					OAuthTokens             tpmodels.TypeOAuthTokens
					RawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider
				}{
					CreatedNewUser: response.OK.CreatedNewUser,
					User: tplmodels.User{
						ID:          response.OK.User.ID,
						TimeJoined:  response.OK.User.TimeJoined,
						Email:       &response.OK.User.Email,
						PhoneNumber: nil,
						ThirdParty:  &response.OK.User.ThirdParty,
					},
					Session:                 response.OK.Session,
					OAuthTokens:             response.OK.OAuthTokens,
					RawUserInfoFromProvider: response.OK.RawUserInfoFromProvider,
				},
			}, nil
		}
	}

	ogAuthorisationUrlGET := *thirdPartyImplementation.AuthorisationUrlGET
	authorisationUrlGET := func(provider tpmodels.TypeProvider, clientID *string, redirectURIOnProviderDashboard string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error) {
		return ogAuthorisationUrlGET(provider, clientID, redirectURIOnProviderDashboard, options, userContext)
	}

	ogAppleRedirectHandlerPOST := *thirdPartyImplementation.AppleRedirectHandlerPOST
	appleRedirectHandlerPOST := func(formPostInfoFromProvider map[string]interface{}, options tpmodels.APIOptions, userContext supertokens.UserContext) error {
		return ogAppleRedirectHandlerPOST(formPostInfoFromProvider, options, userContext)
	}

	ogConsumeCodePOST := *passwordlessImplementation.ConsumeCodePOST
	consumeCodePOST := func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, options plessmodels.APIOptions, userContext supertokens.UserContext) (tplmodels.ConsumeCodePOSTResponse, error) {
		resp, err := ogConsumeCodePOST(userInput, linkCode, preAuthSessionID, options, userContext)
		if err != nil {
			return tplmodels.ConsumeCodePOSTResponse{}, err
		}
		if resp.OK != nil {
			return tplmodels.ConsumeCodePOSTResponse{
				OK: &struct {
					CreatedNewUser bool
					User           tplmodels.User
					Session        sessmodels.SessionContainer
				}{
					CreatedNewUser: resp.OK.CreatedNewUser,
					Session:        resp.OK.Session,
					User: tplmodels.User{
						ID:          resp.OK.User.ID,
						TimeJoined:  resp.OK.User.TimeJoined,
						Email:       resp.OK.User.Email,
						PhoneNumber: resp.OK.User.PhoneNumber,
						ThirdParty:  nil,
					},
				},
			}, nil
		} else if resp.ExpiredUserInputCodeError != nil {
			return tplmodels.ConsumeCodePOSTResponse{
				ExpiredUserInputCodeError: resp.ExpiredUserInputCodeError,
			}, nil
		} else if resp.IncorrectUserInputCodeError != nil {
			return tplmodels.ConsumeCodePOSTResponse{
				IncorrectUserInputCodeError: resp.IncorrectUserInputCodeError,
			}, nil
		} else if resp.RestartFlowError != nil {
			return tplmodels.ConsumeCodePOSTResponse{
				RestartFlowError: &struct{}{},
			}, nil
		} else {
			return tplmodels.ConsumeCodePOSTResponse{
				GeneralError: resp.GeneralError,
			}, nil
		}
	}

	ogCreateCodePOST := *passwordlessImplementation.CreateCodePOST
	createCodePOST := func(email *string, phoneNumber *string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.CreateCodePOSTResponse, error) {
		return ogCreateCodePOST(email, phoneNumber, options, userContext)
	}

	ogEmailExistGET := *passwordlessImplementation.EmailExistsGET
	passwordlessEmailExistsGET := func(email string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.EmailExistsGETResponse, error) {
		return ogEmailExistGET(email, options, userContext)
	}

	ogPhoneNumberExistsGET := *passwordlessImplementation.PhoneNumberExistsGET
	passwordlessPhoneNumberExistsGET := func(phoneNumber string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.PhoneNumberExistsGETResponse, error) {
		return ogPhoneNumberExistsGET(phoneNumber, options, userContext)
	}

	ogResendCodePOST := *passwordlessImplementation.ResendCodePOST
	resendCodePOST := func(deviceID string, preAuthSessionID string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.ResendCodePOSTResponse, error) {
		return ogResendCodePOST(deviceID, preAuthSessionID, options, userContext)
	}

	result := tplmodels.APIInterface{
		AuthorisationUrlGET:              &authorisationUrlGET,
		ThirdPartySignInUpPOST:           &thirdPartySignInUpPOST,
		AppleRedirectHandlerPOST:         &appleRedirectHandlerPOST,
		CreateCodePOST:                   &createCodePOST,
		ResendCodePOST:                   &resendCodePOST,
		ConsumeCodePOST:                  &consumeCodePOST,
		PasswordlessEmailExistsGET:       &passwordlessEmailExistsGET,
		PasswordlessPhoneNumberExistsGET: &passwordlessPhoneNumberExistsGET,
	}

	modifiedPwdless := GetPasswordlessIterfaceImpl(result)
	(*passwordlessImplementation.ConsumeCodePOST) = *modifiedPwdless.ConsumeCodePOST
	(*passwordlessImplementation.CreateCodePOST) = *modifiedPwdless.CreateCodePOST
	(*passwordlessImplementation.EmailExistsGET) = *modifiedPwdless.EmailExistsGET
	(*passwordlessImplementation.PhoneNumberExistsGET) = *modifiedPwdless.PhoneNumberExistsGET
	(*passwordlessImplementation.ResendCodePOST) = *modifiedPwdless.ResendCodePOST

	modifiedTP := GetThirdPartyIterfaceImpl(result)
	(*thirdPartyImplementation.AuthorisationUrlGET) = *modifiedTP.AuthorisationUrlGET
	(*thirdPartyImplementation.SignInUpPOST) = *modifiedTP.SignInUpPOST
	(*thirdPartyImplementation.AppleRedirectHandlerPOST) = *modifiedTP.AppleRedirectHandlerPOST

	return result
}
