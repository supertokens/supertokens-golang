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

func CreateCodeWithEmailWithContext(email string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateCode)(&email, nil, userInputCode, userContext)
}

func CreateCodeWithPhoneNumberWithContext(phoneNumber string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateCode)(nil, &phoneNumber, userInputCode, userContext)
}

func CreateNewCodeForDeviceWithContext(deviceID string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.ResendCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.ResendCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateNewCodeForDevice)(deviceID, userInputCode, userContext)
}

func ConsumeCodeWithUserInputCodeWithContext(deviceID string, userInputCode string, preAuthSessionID string, userContext supertokens.UserContext) (plessmodels.ConsumeCodeResponse, error) {
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
	}, nil, preAuthSessionID, userContext)
}

func ConsumeCodeWithLinkCodeWithContext(linkCode string, preAuthSessionID string, userContext supertokens.UserContext) (plessmodels.ConsumeCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.ConsumeCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ConsumeCode)(nil, &linkCode, preAuthSessionID, userContext)
}

func GetUserByIDWithContext(userID string, userContext supertokens.UserContext) (*plessmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUserByEmailWithContext(email string, userContext supertokens.UserContext) (*plessmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByEmail)(email, userContext)
}

func GetUserByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) (*plessmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByPhoneNumber)(phoneNumber, userContext)
}

func UpdateUserWithContext(userID string, email *string, phoneNumber *string, userContext supertokens.UserContext) (plessmodels.UpdateUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.UpdateUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.UpdateUser)(userID, email, phoneNumber, userContext)
}

func RevokeAllCodesByEmailWithContext(email string, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(&email, nil, userContext)
}

func RevokeAllCodesByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(nil, &phoneNumber, userContext)
}

func RevokeCodeWithContext(codeID string, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeCode)(codeID, userContext)
}

func ListCodesByEmailWithContext(email string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByEmail)(email, userContext)
}

func ListCodesByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByPhoneNumber)(phoneNumber, userContext)
}

func ListCodesByDeviceIDWithContext(deviceID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByDeviceID)(deviceID, userContext)
}

func ListCodesByPreAuthSessionIDWithContext(preAuthSessionID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByPreAuthSessionID)(preAuthSessionID, userContext)
}

func CreateMagicLinkByEmailWithContext(email string, userContext supertokens.UserContext) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.CreateMagicLink(&email, nil, userContext)
}

func CreateMagicLinkByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.CreateMagicLink(nil, &phoneNumber, userContext)
}

func SignInUpByEmailWithContext(email string, userContext supertokens.UserContext) (struct {
	PreAuthSessionID string
	CreatedNewUser   bool
	User             plessmodels.User
}, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return struct {
			PreAuthSessionID string
			CreatedNewUser   bool
			User             plessmodels.User
		}{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.SignInUp(&email, nil, userContext)
}

func SignInUpByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) (struct {
	PreAuthSessionID string
	CreatedNewUser   bool
	User             plessmodels.User
}, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return struct {
			PreAuthSessionID string
			CreatedNewUser   bool
			User             plessmodels.User
		}{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.SignInUp(nil, &phoneNumber, userContext)
}

func DeleteEmailForUserWithContext(userID string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.DeleteUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.DeleteEmailForUser)(userID, userContext)
}

func DeletePhoneNumberForUserWithContext(userID string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.DeleteUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.DeletePhoneNumberForUser)(userID, userContext)
}

func CreateCodeWithEmail(email string, userInputCode *string) (plessmodels.CreateCodeResponse, error) {
	return CreateCodeWithEmailWithContext(email, userInputCode, &map[string]interface{}{})
}

func CreateCodeWithPhoneNumber(phoneNumber string, userInputCode *string) (plessmodels.CreateCodeResponse, error) {
	return CreateCodeWithPhoneNumberWithContext(phoneNumber, userInputCode, &map[string]interface{}{})
}

func CreateNewCodeForDevice(deviceID string, userInputCode *string) (plessmodels.ResendCodeResponse, error) {
	return CreateNewCodeForDeviceWithContext(deviceID, userInputCode, &map[string]interface{}{})
}

func ConsumeCodeWithUserInputCode(deviceID string, userInputCode string, preAuthSessionID string) (plessmodels.ConsumeCodeResponse, error) {
	return ConsumeCodeWithUserInputCodeWithContext(deviceID, userInputCode, preAuthSessionID, &map[string]interface{}{})
}

func ConsumeCodeWithLinkCode(linkCode string, preAuthSessionID string) (plessmodels.ConsumeCodeResponse, error) {
	return ConsumeCodeWithLinkCodeWithContext(linkCode, preAuthSessionID, &map[string]interface{}{})
}

func GetUserByID(userID string) (*plessmodels.User, error) {
	return GetUserByIDWithContext(userID, &map[string]interface{}{})
}

func GetUserByEmail(email string) (*plessmodels.User, error) {
	return GetUserByEmailWithContext(email, &map[string]interface{}{})
}

func GetUserByPhoneNumber(phoneNumber string) (*plessmodels.User, error) {
	return GetUserByPhoneNumberWithContext(phoneNumber, &map[string]interface{}{})
}

func UpdateUser(userID string, email *string, phoneNumber *string) (plessmodels.UpdateUserResponse, error) {
	return UpdateUserWithContext(userID, email, phoneNumber, &map[string]interface{}{})
}

func RevokeAllCodesByEmail(email string) error {
	return RevokeAllCodesByEmailWithContext(email, &map[string]interface{}{})
}

func RevokeAllCodesByPhoneNumber(phoneNumber string) error {
	return RevokeAllCodesByPhoneNumberWithContext(phoneNumber, &map[string]interface{}{})
}

func RevokeCode(codeID string) error {
	return RevokeCodeWithContext(codeID, &map[string]interface{}{})
}

func ListCodesByEmail(email string) ([]plessmodels.DeviceType, error) {
	return ListCodesByEmailWithContext(email, &map[string]interface{}{})
}

func ListCodesByPhoneNumber(phoneNumber string) ([]plessmodels.DeviceType, error) {
	return ListCodesByPhoneNumberWithContext(phoneNumber, &map[string]interface{}{})
}

func ListCodesByDeviceID(deviceID string) (*plessmodels.DeviceType, error) {
	return ListCodesByDeviceIDWithContext(deviceID, &map[string]interface{}{})
}

func ListCodesByPreAuthSessionID(preAuthSessionID string) (*plessmodels.DeviceType, error) {
	return ListCodesByPreAuthSessionIDWithContext(preAuthSessionID, &map[string]interface{}{})
}

func CreateMagicLinkByEmail(email string) (string, error) {
	return CreateMagicLinkByEmailWithContext(email, &map[string]interface{}{})
}

func CreateMagicLinkByPhoneNumber(phoneNumber string) (string, error) {
	return CreateMagicLinkByPhoneNumberWithContext(phoneNumber, &map[string]interface{}{})
}

func DeleteEmailForUser(userID string) (plessmodels.DeleteUserResponse, error) {
	return DeleteEmailForUserWithContext(userID, &map[string]interface{}{})
}

func DeletePhoneNumberForUser(userID string) (plessmodels.DeleteUserResponse, error) {
	return DeletePhoneNumberForUserWithContext(userID, &map[string]interface{}{})
}

func SignInUpByEmail(email string) (struct {
	PreAuthSessionID string
	CreatedNewUser   bool
	User             plessmodels.User
}, error) {
	return SignInUpByEmailWithContext(email, &map[string]interface{}{})
}

func SignInUpByPhoneNumber(phoneNumber string) (struct {
	PreAuthSessionID string
	CreatedNewUser   bool
	User             plessmodels.User
}, error) {
	return SignInUpByPhoneNumberWithContext(phoneNumber, &map[string]interface{}{})
}
