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

package passwordless

import (
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// TODO: test all of these manually
func makeRecipeImplementation(querier supertokens.Querier) plessmodels.RecipeInterface {
	createCode := func(email *string, phoneNumber *string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
		body := map[string]interface{}{}
		if email != nil {
			body["email"] = *email
		} else if phoneNumber != nil {
			body["phoneNumber"] = *phoneNumber
		}
		if userInputCode != nil {
			body["userInputCode"] = *userInputCode
		}
		response, err := querier.SendPostRequest("/recipe/signinup/code", body)
		if err != nil {
			return plessmodels.CreateCodeResponse{}, err
		}
		return plessmodels.CreateCodeResponse{
			OK: &plessmodels.NewCode{
				PreAuthSessionID: response["preAuthSessionId"].(string),
				CodeID:           response["codeId"].(string),
				DeviceID:         response["deviceId"].(string),
				UserInputCode:    response["userInputCode"].(string),
				LinkCode:         response["linkCode"].(string),
				CodeLifetime:     uint64(response["codeLifetime"].(float64)),
				TimeCreated:      uint64(response["timeCreated"].(float64)),
			},
		}, nil
	}

	consumeCode := func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, userContext supertokens.UserContext) (plessmodels.ConsumeCodeResponse, error) {
		body := map[string]interface{}{
			"preAuthSessionId": preAuthSessionID,
		}
		if userInput != nil {
			body["userInputCode"] = userInput.Code
			body["deviceId"] = userInput.DeviceID
		} else if linkCode != nil {
			body["linkCode"] = *linkCode
		}
		response, err := querier.SendPostRequest("/recipe/signinup/code/consume", body)
		if err != nil {
			return plessmodels.ConsumeCodeResponse{}, err
		}
		status := response["status"].(string)
		if status == "OK" {
			return plessmodels.ConsumeCodeResponse{
				OK: &struct {
					CreatedNewUser bool
					User           plessmodels.User
				}{
					CreatedNewUser: response["createdNewUser"].(bool),
					User:           getUserFromJSONResponse(response["user"].(map[string]interface{})),
				},
			}, nil
		} else if status == "INCORRECT_USER_INPUT_CODE_ERROR" {
			return plessmodels.ConsumeCodeResponse{
				IncorrectUserInputCodeError: &struct {
					FailedCodeInputAttemptCount int
					MaximumCodeInputAttempts    int
				}{
					FailedCodeInputAttemptCount: int(response["failedCodeInputAttemptCount"].(float64)),
					MaximumCodeInputAttempts:    int(response["maximumCodeInputAttempts"].(float64)),
				},
			}, nil

		} else if status == "EXPIRED_USER_INPUT_CODE_ERROR" {
			return plessmodels.ConsumeCodeResponse{
				ExpiredUserInputCodeError: &struct {
					FailedCodeInputAttemptCount int
					MaximumCodeInputAttempts    int
				}{
					FailedCodeInputAttemptCount: int(response["failedCodeInputAttemptCount"].(float64)),
					MaximumCodeInputAttempts:    int(response["maximumCodeInputAttempts"].(float64)),
				},
			}, nil
		} else {
			return plessmodels.ConsumeCodeResponse{
				RestartFlowError: &struct{}{},
			}, nil
		}
	}

	createNewCodeForDevice := func(deviceID string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.ResendCodeResponse, error) {
		body := map[string]interface{}{
			"deviceId": deviceID,
		}

		if userInputCode != nil {
			body["userInputCode"] = *userInputCode
		}

		response, err := querier.SendPostRequest("/recipe/signinup/code", body)
		if err != nil {
			return plessmodels.ResendCodeResponse{}, err
		}

		status := response["status"].(string)

		if status == "OK" {
			return plessmodels.ResendCodeResponse{
				OK: &plessmodels.NewCode{
					PreAuthSessionID: response["preAuthSessionId"].(string),
					CodeID:           response["codeId"].(string),
					DeviceID:         response["deviceId"].(string),
					UserInputCode:    response["userInputCode"].(string),
					LinkCode:         response["linkCode"].(string),
					CodeLifetime:     uint64(response["codeLifetime"].(float64)),
					TimeCreated:      uint64(response["timeCreated"].(float64)),
				},
			}, nil
		} else if status == "USER_INPUT_CODE_ALREADY_USED_ERROR" {
			return plessmodels.ResendCodeResponse{
				UserInputCodeAlreadyUsedError: &struct{}{},
			}, nil
		} else {
			return plessmodels.ResendCodeResponse{
				RestartFlowError: &struct{}{},
			}, nil
		}
	}

	getUserByEmail := func(email string, userContext supertokens.UserContext) (*plessmodels.User, error) {
		response, err := querier.SendGetRequest("/recipe/user", map[string]string{
			"email": email,
		})
		if err != nil {
			return nil, err
		}
		status := response["status"].(string)

		if status == "OK" {
			user := getUserFromJSONResponse(response["user"].(map[string]interface{}))
			return &user, nil
		}
		return nil, nil
	}

	getUserByID := func(userID string, userContext supertokens.UserContext) (*plessmodels.User, error) {
		response, err := querier.SendGetRequest("/recipe/user", map[string]string{
			"userId": userID,
		})
		if err != nil {
			return nil, err
		}
		status := response["status"].(string)

		if status == "OK" {
			user := getUserFromJSONResponse(response["user"].(map[string]interface{}))
			return &user, nil
		}
		return nil, nil
	}

	getUserByPhoneNumber := func(phoneNumber string, userContext supertokens.UserContext) (*plessmodels.User, error) {
		response, err := querier.SendGetRequest("/recipe/user", map[string]string{
			"phoneNumber": phoneNumber,
		})
		if err != nil {
			return nil, err
		}
		status := response["status"].(string)

		if status == "OK" {
			user := getUserFromJSONResponse(response["user"].(map[string]interface{}))
			return &user, nil
		}
		return nil, nil
	}

	listCodesByDeviceID := func(deviceID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
		response, err := querier.SendGetRequest("/recipe/signinup/codes", map[string]string{
			"deviceId": deviceID,
		})

		if err != nil {
			return nil, err
		}

		devices := getDevicesFromResponse(response["devices"].([]map[string]interface{}))

		if len(devices) == 1 {
			return &devices[0], nil
		}

		return nil, nil
	}

	listCodesByEmail := func(email string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
		response, err := querier.SendGetRequest("/recipe/signinup/codes", map[string]string{
			"email": email,
		})

		if err != nil {
			return nil, err
		}

		return getDevicesFromResponse(response["devices"].([]map[string]interface{})), nil
	}

	listCodesByPhoneNumber := func(phoneNumber string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
		response, err := querier.SendGetRequest("/recipe/signinup/codes", map[string]string{
			"phoneNumber": phoneNumber,
		})

		if err != nil {
			return nil, err
		}

		return getDevicesFromResponse(response["devices"].([]map[string]interface{})), nil
	}

	listCodesByPreAuthSessionID := func(preAuthSessionID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
		response, err := querier.SendGetRequest("/recipe/signinup/codes", map[string]string{
			"preAuthSessionID": preAuthSessionID,
		})

		if err != nil {
			return nil, err
		}

		devices := getDevicesFromResponse(response["devices"].([]map[string]interface{}))

		if len(devices) == 1 {
			return &devices[0], nil
		}

		return nil, nil
	}

	revokeAllCodes := func(email *string, phoneNumber *string, userContext supertokens.UserContext) error {
		body := map[string]interface{}{}
		if email != nil {
			body["email"] = *email
		} else if phoneNumber != nil {
			body["phoneNumber"] = *phoneNumber
		}
		_, err := querier.SendPostRequest("/recipe/signinup/codes/remove", body)
		if err != nil {
			return err
		}
		return nil
	}

	revokeCode := func(codeID string, userContext supertokens.UserContext) error {
		body := map[string]interface{}{
			"codeId": codeID,
		}
		_, err := querier.SendPostRequest("/recipe/signinup/codes/remove", body)
		if err != nil {
			return err
		}
		return nil
	}

	updateUser := func(userID string, email *string, phoneNumber *string, userContext supertokens.UserContext) (plessmodels.UpdateUserResponse, error) {
		body := map[string]interface{}{
			"userId": userID,
		}
		if email != nil {
			body["email"] = *email
		}
		if phoneNumber != nil {
			body["phoneNumber"] = *phoneNumber
		}

		response, err := querier.SendPutRequest("/recipe/user", body)
		if err != nil {
			return plessmodels.UpdateUserResponse{}, err
		}

		status := response["status"].(string)

		if status == "OK" {
			return plessmodels.UpdateUserResponse{
				OK: &struct{}{},
			}, nil
		} else if status == "UNKNOWN_USER_ID_ERROR" {
			return plessmodels.UpdateUserResponse{
				UnknownUserIdError: &struct{}{},
			}, nil
		} else if status == "EMAIL_ALREADY_EXISTS_ERROR" {
			return plessmodels.UpdateUserResponse{
				EmailAlreadyExistsError: &struct{}{},
			}, nil
		} else {
			return plessmodels.UpdateUserResponse{
				PhoneNumberAlreadyExistsError: &struct{}{},
			}, nil
		}
	}

	return plessmodels.RecipeInterface{
		CreateCode:                  &createCode,
		ConsumeCode:                 &consumeCode,
		CreateNewCodeForDevice:      &createNewCodeForDevice,
		GetUserByEmail:              &getUserByEmail,
		GetUserByID:                 &getUserByID,
		GetUserByPhoneNumber:        &getUserByPhoneNumber,
		ListCodesByDeviceID:         &listCodesByDeviceID,
		ListCodesByEmail:            &listCodesByEmail,
		ListCodesByPhoneNumber:      &listCodesByPhoneNumber,
		ListCodesByPreAuthSessionID: &listCodesByPreAuthSessionID,
		RevokeAllCodes:              &revokeAllCodes,
		RevokeCode:                  &revokeCode,
		UpdateUser:                  &updateUser,
	}
}

func getDevicesFromResponse(devicesJSON []map[string]interface{}) []plessmodels.DeviceType {
	result := []plessmodels.DeviceType{}
	for _, deviceJSON := range devicesJSON {
		device := plessmodels.DeviceType{
			PreAuthSessionID:            deviceJSON["preAuthSessionId"].(string),
			FailedCodeInputAttemptCount: int(deviceJSON["failedCodeInputAttemptCount"].(float64)),
			Codes:                       getCodesFromDevicesResponse(deviceJSON["codes"].([]map[string]interface{})),
		}
		{
			email, ok := deviceJSON["email"]
			if ok {
				emailStr := email.(string)
				device.Email = &emailStr
			}
		}
		{
			phoneNumber, ok := deviceJSON["phoneNumber"]
			if ok {
				phoneNumberStr := phoneNumber.(string)
				device.PhoneNumber = &phoneNumberStr
			}
		}

		result = append(result, device)
	}
	return result
}

func getCodesFromDevicesResponse(codesJSON []map[string]interface{}) []plessmodels.Code {
	result := []plessmodels.Code{}
	for _, codeJSON := range codesJSON {
		code := plessmodels.Code{
			CodeID:       codeJSON["codeId"].(string),
			TimeCreated:  uint64(codeJSON["timeCreated"].(float64)),
			CodeLifetime: uint64(codeJSON["codeLifetime"].(float64)),
		}
		result = append(result, code)
	}
	return result
}

func getUserFromJSONResponse(userJSON map[string]interface{}) plessmodels.User {
	user := plessmodels.User{
		ID:         userJSON["id"].(string),
		TimeJoined: uint64(userJSON["timeJoined"].(float64)),
	}
	{
		email, ok := userJSON["email"]
		if ok {
			emailStr := email.(string)
			user.Email = &emailStr
		}
	}
	{
		phoneNumber, ok := userJSON["phoneNumber"]
		if ok {
			phoneNumberStr := phoneNumber.(string)
			user.PhoneNumber = &phoneNumberStr
		}
	}
	return user
}
