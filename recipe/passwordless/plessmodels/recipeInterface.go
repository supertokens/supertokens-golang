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

package plessmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type RecipeInterface struct {
	CreateCode *func(email *string, phoneNumber *string, userInputCode *string, userContext supertokens.UserContext) (CreateCodeResponse, error)

	ResendCode *func(deviceID string, userInputCode *string, userContext supertokens.UserContext) (ResendCodeResponse, error)

	ConsumeCode *func(userInput *UserInputCodeWithDeviceID, linkeCode *string, userContext supertokens.UserContext) (ConsumeCodeResponse, error)

	GetUserByID *func(userID string, userContext supertokens.UserContext) (*User, error)

	GetUserByEmail *func(email string, userContext supertokens.UserContext) (*User, error)

	GetUserByPhoneNumber *func(phoneNumber string, userContext supertokens.UserContext) (*User, error)

	UpdateUser *func(userID string, email *string, phoneNumber *string, userContext supertokens.UserContext) (UpdateUserResponse, error)

	RevokeAllCodes *func(email *string, phoneNumber *string, userContext supertokens.UserContext) error

	RevokeCode *func(codeID string, userContext supertokens.UserContext) error

	ListCodesByEmail *func(email string, userContext supertokens.UserContext) ([]DeviceType, error)

	ListCodesByPhoneNumber *func(phoneNumber string, userContext supertokens.UserContext) ([]DeviceType, error)

	ListCodesByDeviceID *func(deviceID string, userContext supertokens.UserContext) (*DeviceType, error)

	ListCodesByPreAuthSessionID *func(preAuthSessionID string, userContext supertokens.UserContext) (*DeviceType, error)
}

type DeviceType struct {
	PreAuthSessionID            string
	FailedCodeInputAttemptCount int
	Email                       *string
	PhoneNumber                 *string
	Codes                       []Code
}

type Code struct {
	CodeID       string
	TimeCreated  uint64
	CodeLifetime uint64
}

type ConsumeCodeResponse struct {
	OK *struct {
		PreAuthSessionID string
		CreatedNewUser   bool
		User             User
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

type UserInputCodeWithDeviceID struct {
	Code     string
	DeviceID string
}

type ResendCodeResponse struct {
	OK                            *NewCode
	RestartFlowError              *struct{}
	UserInputCodeAlreadyUsedError *struct{}
}

type CreateCodeResponse struct {
	OK *NewCode
}

type NewCode struct {
	PreAuthSessionID string
	CodeID           string
	DeviceID         string
	UserInputCode    string
	LinkCode         string
	CodeLifetime     uint64
	TimeCreated      uint64
}

type UpdateUserResponse struct {
	OK                            *struct{}
	UnknownUserIdError            *struct{}
	EmailAlreadyExistsError       *struct{}
	PhoneNumberAlreadyExistsError *struct{}
}
