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

	return plessmodels.APIInterface{
		ConsumeCodePOST: &consumeCodePOST,
		// TODO:
	}
}
