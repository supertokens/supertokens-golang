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
		response, err := (*options.RecipeImplementation.ConsumeCode)(userInput, linkCode, userContext)
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

		session, err := session.CreateNewSession(options.Res, user.ID, map[string]interface{}{}, map[string]interface{}{})
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

	return plessmodels.APIInterface{
		ConsumeCodePOST: &consumeCodePOST,
		CreateCodePOST:  &createCodePOST,
		// TODO:
	}
}
