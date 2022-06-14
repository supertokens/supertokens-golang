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
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type APIOptions struct {
	RecipeImplementation RecipeInterface
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
	EmailDelivery        emaildelivery.Ingredient
	SmsDelivery          smsdelivery.Ingredient
}

type APIInterface struct {
	CreateCodePOST       *func(email *string, phoneNumber *string, options APIOptions, userContext supertokens.UserContext) (CreateCodePOSTResponse, error)
	ResendCodePOST       *func(deviceID string, preAuthSessionID string, options APIOptions, userContext supertokens.UserContext) (ResendCodePOSTResponse, error)
	ConsumeCodePOST      *func(userInput *UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, options APIOptions, userContext supertokens.UserContext) (ConsumeCodePOSTResponse, error)
	EmailExistsGET       *func(email string, options APIOptions, userContext supertokens.UserContext) (EmailExistsGETResponse, error)
	PhoneNumberExistsGET *func(phoneNumber string, options APIOptions, userContext supertokens.UserContext) (PhoneNumberExistsGETResponse, error)
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

type ResendCodePOSTResponse struct {
	OK             *struct{}
	ResetFlowError *struct{}
	GeneralError   *supertokens.GeneralErrorResponse
}

type CreateCodePOSTResponse struct {
	OK *struct {
		DeviceID         string
		PreAuthSessionID string
		FlowType         string
	}
	GeneralError *supertokens.GeneralErrorResponse
}

type EmailExistsGETResponse struct {
	OK           *struct{ Exists bool }
	GeneralError *supertokens.GeneralErrorResponse
}

type PhoneNumberExistsGETResponse struct {
	OK           *struct{ Exists bool }
	GeneralError *supertokens.GeneralErrorResponse
}
