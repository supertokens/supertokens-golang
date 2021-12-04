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

func Init(config plessmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func CreateCodeWithEmail(email string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateCode)(&email, nil, userInputCode, userContext)
}

func CreateCodeWithPhoneNumber(phoneNumber string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateCode)(nil, &phoneNumber, userInputCode, userContext)
}

func CreateNewCodeForDevice(deviceID string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.ResendCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.ResendCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateNewCodeForDevice)(deviceID, userInputCode, userContext)
}

func ConsumeCodeWithUserInputCode(deviceID string, userInputCode string, userContext supertokens.UserContext) (plessmodels.ConsumeCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.ConsumeCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ConsumeCode)(&plessmodels.UserInputCodeWithDeviceID{
		Code:     userInputCode,
		DeviceID: deviceID,
	}, nil, userContext)
}

func ConsumeCodeWithLinkCode(linkCode string, userContext supertokens.UserContext) (plessmodels.ConsumeCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.ConsumeCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ConsumeCode)(nil, &linkCode, userContext)
}

func GetUserByID(userID string, userContext supertokens.UserContext) (*plessmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUserByEmail(email string, userContext supertokens.UserContext) (*plessmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByEmail)(email, userContext)
}

func GetUserByPhoneNumber(phoneNumber string, userContext supertokens.UserContext) (*plessmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByPhoneNumber)(phoneNumber, userContext)
}

func UpdateUser(userID string, email *string, phoneNumber *string, userContext supertokens.UserContext) (plessmodels.UpdateUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.UpdateUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.UpdateUser)(userID, email, phoneNumber, userContext)
}

func RevokeAllCodesByEmail(email string, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(&email, nil, userContext)
}

func RevokeAllCodesByPhoneNumber(phoneNumber string, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(nil, &phoneNumber, userContext)
}

func RevokeCode(codeID string, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeCode)(codeID, userContext)
}

// TODO:
