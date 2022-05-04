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
	"github.com/supertokens/supertokens-golang/supertokens"
)

type RecipeInterface struct {
	GetUserByID             *func(userID string, userContext supertokens.UserContext) (*User, error)
	GetUsersByEmail         *func(email string, userContext supertokens.UserContext) ([]User, error)
	GetUserByPhoneNumber    *func(phoneNumber string, userContext supertokens.UserContext) (*User, error)
	GetUserByThirdPartyInfo *func(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*User, error)
	ThirdPartySignInUp      *func(thirdPartyID string, thirdPartyUserID string, email EmailStruct, userContext supertokens.UserContext) (ThirdPartySignInUp, error)

	CreateCode *func(email *string, phoneNumber *string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error)

	CreateNewCodeForDevice *func(deviceID string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.ResendCodeResponse, error)

	ConsumeCode *func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, userContext supertokens.UserContext) (ConsumeCodeResponse, error)

	UpdatePasswordlessUser *func(userID string, email *string, phoneNumber *string, userContext supertokens.UserContext) (plessmodels.UpdateUserResponse, error)

	DeleteEmailForPasswordlessUser *func(userID string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error)

	DeletePhoneNumberForUser *func(userID string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error)

	RevokeAllCodes *func(email *string, phoneNumber *string, userContext supertokens.UserContext) error

	RevokeCode *func(codeID string, userContext supertokens.UserContext) error

	ListCodesByEmail *func(email string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error)

	ListCodesByPhoneNumber *func(phoneNumber string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error)

	ListCodesByDeviceID *func(deviceID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error)

	ListCodesByPreAuthSessionID *func(preAuthSessionID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error)
}

type ConsumeCodeResponse struct {
	OK *struct {
		CreatedNewUser bool
		User           User
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
}

type ThirdPartySignInUp struct {
	OK *struct {
		CreatedNewUser bool
		User           User
	}
	FieldError *struct{ ErrorMsg string }
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
