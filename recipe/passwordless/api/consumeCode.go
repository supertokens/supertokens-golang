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
	"encoding/json"
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ConsumeCode(apiImplementation plessmodels.APIInterface, tenantId string, options plessmodels.APIOptions, userContext supertokens.UserContext) error {

	if apiImplementation.ConsumeCodePOST == nil || (*apiImplementation.ConsumeCodePOST) == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	body, err := supertokens.ReadFromRequest(options.Req)
	if err != nil {
		return err
	}
	var readBody map[string]interface{}
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return err
	}

	preAuthSessionID, okPreAuthSessionID := readBody["preAuthSessionId"]
	linkCode, okLinkCode := readBody["linkCode"]
	deviceID, okDeviceID := readBody["deviceId"]
	userInputCode, okUserInputCode := readBody["userInputCode"]

	if !okPreAuthSessionID || reflect.ValueOf(preAuthSessionID).Kind() != reflect.String {
		return supertokens.BadInputError{Msg: "Please provide preAuthSessionId"}
	}

	if okUserInputCode || okDeviceID {
		// if either userInputCode or deviceId exists in the input
		if okLinkCode {
			return supertokens.BadInputError{Msg: "Please provide one of (linkCode) or (deviceId+userInputCode) and not both"}
		}

		if !okUserInputCode || !okDeviceID {
			return supertokens.BadInputError{Msg: "Please provide both deviceId and userInputCode"}
		}

		if reflect.ValueOf(userInputCode).Kind() != reflect.String {
			return supertokens.BadInputError{Msg: "Please make sure that userInputCode is a string"}
		}

		if reflect.ValueOf(deviceID).Kind() != reflect.String {
			return supertokens.BadInputError{Msg: "Please make sure that deviceId is a string"}
		}
	} else if !okLinkCode {
		return supertokens.BadInputError{Msg: "Please provide one of (linkCode) or (deviceId+userInputCode) and not both"}
	}

	if okLinkCode && reflect.ValueOf(linkCode).Kind() != reflect.String {
		return supertokens.BadInputError{Msg: "Please make sure that linkCode is a string"}
	}

	var userInput *plessmodels.UserInputCodeWithDeviceID

	if okUserInputCode {
		userInput = &plessmodels.UserInputCodeWithDeviceID{
			Code:     userInputCode.(string),
			DeviceID: deviceID.(string),
		}
	}

	var linkCodePointer *string
	if okLinkCode {
		t := linkCode.(string)
		linkCodePointer = &t
	}

	response, err := (*apiImplementation.ConsumeCodePOST)(userInput, linkCodePointer, preAuthSessionID.(string), tenantId, options, userContext)
	if err != nil {
		return err
	}

	var result map[string]interface{}

	if response.OK != nil {
		result = map[string]interface{}{
			"status":         "OK",
			"createdNewUser": response.OK.CreatedNewUser,
			"user":           response.OK.User,
		}
	} else if response.ExpiredUserInputCodeError != nil {
		result = map[string]interface{}{
			"status":                      "EXPIRED_USER_INPUT_CODE_ERROR",
			"failedCodeInputAttemptCount": response.ExpiredUserInputCodeError.FailedCodeInputAttemptCount,
			"maximumCodeInputAttempts":    response.ExpiredUserInputCodeError.MaximumCodeInputAttempts,
		}
	} else if response.IncorrectUserInputCodeError != nil {
		result = map[string]interface{}{
			"status":                      "INCORRECT_USER_INPUT_CODE_ERROR",
			"failedCodeInputAttemptCount": response.IncorrectUserInputCodeError.FailedCodeInputAttemptCount,
			"maximumCodeInputAttempts":    response.IncorrectUserInputCodeError.MaximumCodeInputAttempts,
		}
	} else if response.RestartFlowError != nil {
		result = map[string]interface{}{
			"status": "RESTART_FLOW_ERROR",
		}
	} else if response.GeneralError != nil {
		result = supertokens.ConvertGeneralErrorToJsonResponse(*response.GeneralError)
	} else {
		return supertokens.ErrorIfNoResponse(options.Res)
	}

	return supertokens.Send200Response(options.Res, result)
}
