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
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() plessmodels.APIInterface {

	consumeCodePOST := func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.ConsumeCodePOSTResponse, error) {
		response, err := (*options.RecipeImplementation.ConsumeCode)(userInput, linkCode, preAuthSessionID, userContext)
		if err != nil {
			return plessmodels.ConsumeCodePOSTResponse{}, err
		}

		if response.OK == nil {
			return plessmodels.ConsumeCodePOSTResponse{
				IncorrectUserInputCodeError: response.IncorrectUserInputCodeError,
				ExpiredUserInputCodeError:   response.ExpiredUserInputCodeError,
				RestartFlowError:            response.RestartFlowError,
			}, nil
		}

		user := response.OK.User

		session, err := session.CreateNewSession(options.Res, user.ID, map[string]interface{}{}, map[string]interface{}{}, userContext)
		if err != nil {
			return plessmodels.ConsumeCodePOSTResponse{}, err
		}

		return plessmodels.ConsumeCodePOSTResponse{
			OK: &struct {
				CreatedNewUser bool
				User           plessmodels.User
				Session        sessmodels.SessionContainer
			}{
				CreatedNewUser: response.OK.CreatedNewUser,
				User:           response.OK.User,
				Session:        session,
			},
		}, nil
	}

	createCodePOST := func(email *string, phoneNumber *string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.CreateCodePOSTResponse, error) {

		var userInputCodeInput *string
		if options.Config.GetCustomUserInputCode != nil {
			c, err := options.Config.GetCustomUserInputCode(userContext)
			if err != nil {
				return plessmodels.CreateCodePOSTResponse{}, err
			}
			userInputCodeInput = &c
		}

		response, err := (*options.RecipeImplementation.CreateCode)(email, phoneNumber, userInputCodeInput, userContext)
		if err != nil {
			return plessmodels.CreateCodePOSTResponse{}, err
		}

		// now we will send an email / text message
		var magicLink *string
		var userInputCode *string
		flowType := options.Config.FlowType
		if flowType == "MAGIC_LINK" || flowType == "USER_INPUT_CODE_AND_MAGIC_LINK" {
			link, err := options.Config.GetLinkDomainAndPath(email, phoneNumber, userContext)
			if err != nil {
				return plessmodels.CreateCodePOSTResponse{}, err
			}
			link = link + "?rid=" + options.RecipeID + "&preAuthSessionId=" + response.OK.PreAuthSessionID + "#" + response.OK.LinkCode

			magicLink = &link
		}

		if flowType == "USER_INPUT_CODE" || flowType == "USER_INPUT_CODE_AND_MAGIC_LINK" {
			userInputCode = &response.OK.UserInputCode
		}

		if options.Config.ContactMethodPhone.Enabled {
			options.Config.ContactMethodPhone.CreateAndSendCustomTextMessage(*phoneNumber, userInputCode, magicLink, response.OK.CodeLifetime, response.OK.PreAuthSessionID, userContext)
		} else {
			options.Config.ContactMethodEmail.CreateAndSendCustomEmail(*email, userInputCode, magicLink, response.OK.CodeLifetime, response.OK.PreAuthSessionID, userContext)
		}

		return plessmodels.CreateCodePOSTResponse{
			OK: &struct {
				DeviceID         string
				PreAuthSessionID string
				FlowType         string
			}{
				DeviceID:         response.OK.DeviceID,
				PreAuthSessionID: response.OK.PreAuthSessionID,
				FlowType:         options.Config.FlowType,
			},
		}, nil
	}

	emailExistsGET := func(email string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.EmailExistsGETResponse, error) {
		response, err := (*options.RecipeImplementation.GetUserByEmail)(email, userContext)
		if err != nil {
			return plessmodels.EmailExistsGETResponse{}, err
		}

		return plessmodels.EmailExistsGETResponse{
			OK: &struct{ Exists bool }{
				Exists: response != nil,
			},
		}, nil
	}

	phoneNumberExistsGET := func(phoneNumber string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.PhoneNumberExistsGETResponse, error) {
		response, err := (*options.RecipeImplementation.GetUserByPhoneNumber)(phoneNumber, userContext)
		if err != nil {
			return plessmodels.PhoneNumberExistsGETResponse{}, err
		}

		return plessmodels.PhoneNumberExistsGETResponse{
			OK: &struct{ Exists bool }{
				Exists: response != nil,
			},
		}, nil
	}

	resendCodePOST := func(deviceID string, preAuthSessionID string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.ResendCodePOSTResponse, error) {
		deviceInfo, err := (*options.RecipeImplementation.ListCodesByDeviceID)(deviceID, userContext)
		if err != nil {
			return plessmodels.ResendCodePOSTResponse{}, err
		}

		if deviceInfo == nil {
			return plessmodels.ResendCodePOSTResponse{
				ResetFlowError: &struct{}{},
			}, nil
		}

		if (options.Config.ContactMethodEmail.Enabled && deviceInfo.Email == nil) || (options.Config.ContactMethodPhone.Enabled && deviceInfo.PhoneNumber == nil) {
			return plessmodels.ResendCodePOSTResponse{
				ResetFlowError: &struct{}{},
			}, nil
		}

		for numberOfTriesToCreateNewCode := 0; numberOfTriesToCreateNewCode < 3; numberOfTriesToCreateNewCode++ {
			var userInputCodeInput *string
			if options.Config.GetCustomUserInputCode != nil {
				c, err := options.Config.GetCustomUserInputCode(userContext)
				if err != nil {
					return plessmodels.ResendCodePOSTResponse{}, err
				}
				userInputCodeInput = &c
			}
			response, err := (*options.RecipeImplementation.CreateNewCodeForDevice)(deviceID, userInputCodeInput, userContext)
			if err != nil {
				return plessmodels.ResendCodePOSTResponse{}, err
			}

			if response.UserInputCodeAlreadyUsedError != nil {
				continue
			}

			if response.RestartFlowError != nil {
				return plessmodels.ResendCodePOSTResponse{
					ResetFlowError: response.RestartFlowError,
				}, nil
			}

			var magicLink *string
			var userInputCode *string
			flowType := options.Config.FlowType
			if flowType == "MAGIC_LINK" || flowType == "USER_INPUT_CODE_AND_MAGIC_LINK" {
				link, err := options.Config.GetLinkDomainAndPath(deviceInfo.Email, deviceInfo.PhoneNumber, userContext)
				if err != nil {
					return plessmodels.ResendCodePOSTResponse{}, err
				}
				link = link + "?rid=" + options.RecipeID + "&preAuthSessionId=" + response.OK.PreAuthSessionID + "#" + response.OK.LinkCode

				magicLink = &link
			}

			if flowType == "USER_INPUT_CODE" || flowType == "USER_INPUT_CODE_AND_MAGIC_LINK" {
				userInputCode = &response.OK.UserInputCode
			}

			if options.Config.ContactMethodPhone.Enabled {
				options.Config.ContactMethodPhone.CreateAndSendCustomTextMessage(*deviceInfo.PhoneNumber, userInputCode, magicLink, response.OK.CodeLifetime, response.OK.PreAuthSessionID, userContext)
			} else {
				options.Config.ContactMethodEmail.CreateAndSendCustomEmail(*deviceInfo.Email, userInputCode, magicLink, response.OK.CodeLifetime, response.OK.PreAuthSessionID, userContext)
			}

			return plessmodels.ResendCodePOSTResponse{
				OK: &struct{}{},
			}, nil

		}

		return plessmodels.ResendCodePOSTResponse{
			GeneralError: &struct{ Message string }{
				Message: "Failed to generate a one time code. Please try again",
			},
		}, nil
	}

	return plessmodels.APIInterface{
		ConsumeCodePOST:      &consumeCodePOST,
		CreateCodePOST:       &createCodePOST,
		EmailExistsGET:       &emailExistsGET,
		PhoneNumberExistsGET: &phoneNumberExistsGET,
		ResendCodePOST:       &resendCodePOST,
	}
}
