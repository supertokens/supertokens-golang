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

package recipeimplementation

import (
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakePasswordlessRecipeImplementation(recipeImplementation tplmodels.RecipeInterface) plessmodels.RecipeInterface {
	createCode := func(email *string, phoneNumber *string, userInputCode *string, tenantId *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
		return (*recipeImplementation.CreateCode)(email, phoneNumber, userInputCode, tenantId, userContext)
	}

	consumeCode := func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, tenantId *string, userContext supertokens.UserContext) (plessmodels.ConsumeCodeResponse, error) {
		resp, err := (*recipeImplementation.ConsumeCode)(userInput, linkCode, preAuthSessionID, tenantId, userContext)

		if err != nil {
			return plessmodels.ConsumeCodeResponse{}, err
		}

		if resp.ExpiredUserInputCodeError != nil {
			return plessmodels.ConsumeCodeResponse{
				ExpiredUserInputCodeError: resp.ExpiredUserInputCodeError,
			}, nil
		} else if resp.IncorrectUserInputCodeError != nil {
			return plessmodels.ConsumeCodeResponse{
				IncorrectUserInputCodeError: resp.IncorrectUserInputCodeError,
			}, nil
		} else if resp.RestartFlowError != nil {
			return plessmodels.ConsumeCodeResponse{
				RestartFlowError: &struct{}{},
			}, nil
		} else {
			return plessmodels.ConsumeCodeResponse{
				OK: &struct {
					CreatedNewUser bool
					User           plessmodels.User
				}{
					CreatedNewUser: resp.OK.CreatedNewUser,
					User: plessmodels.User{
						ID:          resp.OK.User.ID,
						Email:       resp.OK.User.Email,
						PhoneNumber: resp.OK.User.PhoneNumber,
						TimeJoined:  resp.OK.User.TimeJoined,
					},
				},
			}, nil
		}
	}

	createNewCodeForDevice := func(deviceID string, userInputCode *string, tenantId *string, userContext supertokens.UserContext) (plessmodels.ResendCodeResponse, error) {
		return (*recipeImplementation.CreateNewCodeForDevice)(deviceID, userInputCode, tenantId, userContext)
	}

	getUserByEmail := func(email string, tenantId *string, userContext supertokens.UserContext) (*plessmodels.User, error) {
		resp, err := (*recipeImplementation.GetUsersByEmail)(email, tenantId, userContext)
		if err != nil {
			return nil, err
		}

		for _, user := range resp {
			if user.ThirdParty == nil {
				return &plessmodels.User{
					ID:          user.ID,
					Email:       user.Email,
					PhoneNumber: user.PhoneNumber,
					TimeJoined:  user.TimeJoined,
				}, nil
			}
		}

		return nil, nil
	}

	getUserByID := func(userID string, tenantId *string, userContext supertokens.UserContext) (*plessmodels.User, error) {
		resp, err := (*recipeImplementation.GetUserByID)(userID, tenantId, userContext)

		if err != nil {
			return nil, err
		}

		if resp == nil {
			return nil, nil
		}

		if resp.ThirdParty != nil {
			// this is a thirdparty user
			return nil, nil
		}

		return &plessmodels.User{
			ID:          resp.ID,
			Email:       resp.Email,
			PhoneNumber: resp.PhoneNumber,
			TimeJoined:  resp.TimeJoined,
		}, nil
	}

	getUserByPhoneNumber := func(phoneNumber string, tenantId *string, userContext supertokens.UserContext) (*plessmodels.User, error) {
		resp, err := (*recipeImplementation.GetUserByPhoneNumber)(phoneNumber, tenantId, userContext)

		if err != nil {
			return nil, err
		}

		if resp == nil {
			return nil, nil
		}

		if resp.ThirdParty != nil {
			// this is a thirdparty user
			return nil, nil
		}

		return &plessmodels.User{
			ID:          resp.ID,
			Email:       resp.Email,
			PhoneNumber: resp.PhoneNumber,
			TimeJoined:  resp.TimeJoined,
		}, nil
	}

	listCodesByDeviceID := func(deviceID string, tenantId *string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
		return (*recipeImplementation.ListCodesByDeviceID)(deviceID, tenantId, userContext)
	}

	listCodesByEmail := func(email string, tenantId *string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
		return (*recipeImplementation.ListCodesByEmail)(email, tenantId, userContext)
	}

	listCodesByPhoneNumber := func(phoneNumber string, tenantId *string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
		return (*recipeImplementation.ListCodesByPhoneNumber)(phoneNumber, tenantId, userContext)
	}

	listCodesByPreAuthSessionID := func(preAuthSessionID string, tenantId *string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
		return (*recipeImplementation.ListCodesByPreAuthSessionID)(preAuthSessionID, tenantId, userContext)
	}

	revokeAllCodes := func(email *string, phoneNumber *string, tenantId *string, userContext supertokens.UserContext) error {
		return (*recipeImplementation.RevokeAllCodes)(email, phoneNumber, tenantId, userContext)
	}

	revokeCode := func(codeID string, tenantId *string, userContext supertokens.UserContext) error {
		return (*recipeImplementation.RevokeCode)(codeID, tenantId, userContext)
	}

	updateUser := func(userID string, email *string, phoneNumber *string, tenantId *string, userContext supertokens.UserContext) (plessmodels.UpdateUserResponse, error) {
		return (*recipeImplementation.UpdatePasswordlessUser)(userID, email, phoneNumber, tenantId, userContext)
	}

	deleteEmailForUser := func(userID string, tenantId *string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
		return (*recipeImplementation.DeleteEmailForPasswordlessUser)(userID, tenantId, userContext)
	}

	deletePhoneNumberForUser := func(userID string, tenantId *string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
		return (*recipeImplementation.DeletePhoneNumberForUser)(userID, tenantId, userContext)
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
		DeleteEmailForUser:          &deleteEmailForUser,
		DeletePhoneNumberForUser:    &deletePhoneNumberForUser,
	}
}
