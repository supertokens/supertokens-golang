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

package tplmodels

import (
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type APIInterface struct {
	AuthorisationUrlGET *func(provider tpmodels.TypeProvider, clientID *string, redirectURIOnProviderDashboard string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error)

	AppleRedirectHandlerPOST *func(formPostInfoFromProvider map[string]interface{}, options tpmodels.APIOptions, userContext supertokens.UserContext) error

	ThirdPartySignInUpPOST *func(provider tpmodels.TypeProvider, clientID *string, input tpmodels.TypeSignInUpInput, options tpmodels.APIOptions, userContext supertokens.UserContext) (ThirdPartySignInUpOutput, error)

	CreateCodePOST *func(email *string, phoneNumber *string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.CreateCodePOSTResponse, error)

	ResendCodePOST *func(deviceID string, preAuthSessionID string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.ResendCodePOSTResponse, error)

	ConsumeCodePOST *func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, options plessmodels.APIOptions, userContext supertokens.UserContext) (ConsumeCodePOSTResponse, error)

	PasswordlessEmailExistsGET *func(email string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.EmailExistsGETResponse, error)

	PasswordlessPhoneNumberExistsGET *func(email string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.PhoneNumberExistsGETResponse, error)
}

type ConsumeCodePOSTResponse struct {
	OK *struct {
		CreatedNewUser bool
		User           User
		Session        sessmodels.SessionContainer
	}
	IncorrectUserInputCodeError *struct {
		FailedCodeInputAttemptCount int
		MaximumCodeInputAttempts    int
	}
	ExpiredUserInputCodeError *struct {
		FailedCodeInputAttemptCount int
		MaximumCodeInputAttempts    int
	}
	RestartFlowError *struct{}
	GeneralError     *supertokens.GeneralErrorResponse
}

type ThirdPartySignInUpOutput struct {
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
