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
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GetPasswordlessIterfaceImpl(apiImplmentation tplmodels.APIInterface) plessmodels.APIInterface {

	result := plessmodels.APIInterface{
		CreateCodePOST:       apiImplmentation.CreateCodePOST,
		ResendCodePOST:       apiImplmentation.ResendCodePOST,
		EmailExistsGET:       apiImplmentation.PasswordlessEmailExistsGET,
		PhoneNumberExistsGET: apiImplmentation.PasswordlessPhoneNumberExistsGET,
		ConsumeCodePOST:      nil,
	}

	if apiImplmentation.ConsumeCodePOST != nil && (*apiImplmentation.ConsumeCodePOST) != nil {
		consumeCodePOST := func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, tenantId *string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.ConsumeCodePOSTResponse, error) {
			result, err := (*apiImplmentation.ConsumeCodePOST)(userInput, linkCode, preAuthSessionID, tenantId, options, userContext)
			if err != nil {
				return plessmodels.ConsumeCodePOSTResponse{}, err
			}
			if result.OK != nil {
				return plessmodels.ConsumeCodePOSTResponse{OK: &struct {
					CreatedNewUser bool
					User           plessmodels.User
					Session        sessmodels.SessionContainer
				}{
					CreatedNewUser: result.OK.CreatedNewUser,
					User: plessmodels.User{
						ID:          result.OK.User.ID,
						Email:       result.OK.User.Email,
						PhoneNumber: result.OK.User.PhoneNumber,
						TimeJoined:  result.OK.User.TimeJoined,
					},
					Session: result.OK.Session,
				}}, nil
			} else if result.ExpiredUserInputCodeError != nil {
				return plessmodels.ConsumeCodePOSTResponse{
					ExpiredUserInputCodeError: result.ExpiredUserInputCodeError,
				}, nil
			} else if result.IncorrectUserInputCodeError != nil {
				return plessmodels.ConsumeCodePOSTResponse{
					IncorrectUserInputCodeError: result.IncorrectUserInputCodeError,
				}, nil
			} else if result.RestartFlowError != nil {
				return plessmodels.ConsumeCodePOSTResponse{
					RestartFlowError: &struct{}{},
				}, nil
			} else {
				return plessmodels.ConsumeCodePOSTResponse{
					GeneralError: result.GeneralError,
				}, nil
			}
		}
		result.ConsumeCodePOST = &consumeCodePOST
	}

	return result
}
