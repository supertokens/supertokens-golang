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

package thirdpartypasswordless

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config tplmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func ThirdPartySignInUp(thirdPartyID string, thirdPartyUserID string, email tplmodels.EmailStruct, userContext supertokens.UserContext) (tplmodels.ThirdPartySignInUp, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ThirdPartySignInUp{}, err
	}
	return (*instance.RecipeImpl.ThirdPartySignInUp)(thirdPartyID, thirdPartyUserID, email, userContext)
}

func GetUserByThirdPartyInfo(thirdPartyID string, thirdPartyUserID string, email tpmodels.EmailStruct, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, userContext)
}

func GetUserById(userID string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUsersByEmail(email string, userContext supertokens.UserContext) ([]tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, userContext)
}

func CreateEmailVerificationToken(userID string, userContext supertokens.UserContext) (evmodels.CreateEmailVerificationTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	email, err := instance.getEmailForUserIdForEmailVerification(userID, userContext)
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	return (*instance.EmailVerificationRecipe.RecipeImpl.CreateEmailVerificationToken)(userID, email, userContext)
}

func VerifyEmailUsingToken(token string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	response, err := (*instance.EmailVerificationRecipe.RecipeImpl.VerifyEmailUsingToken)(token, userContext)
	if err != nil {
		return nil, err
	}
	if response.EmailVerificationInvalidTokenError != nil {
		return nil, errors.New("email verification token is invalid")
	}
	return (*instance.RecipeImpl.GetUserByID)(response.OK.User.ID, userContext)
}

func IsEmailVerified(userID string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	email, err := instance.getEmailForUserIdForEmailVerification(userID, userContext)
	if err != nil {
		return false, err
	}
	return (*instance.EmailVerificationRecipe.RecipeImpl.IsEmailVerified)(userID, email, userContext)
}

func RevokeEmailVerificationTokens(userID string, userContext supertokens.UserContext) (evmodels.RevokeEmailVerificationTokensResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.RevokeEmailVerificationTokensResponse{}, err
	}
	email, err := instance.getEmailForUserIdForEmailVerification(userID, userContext)
	if err != nil {
		return evmodels.RevokeEmailVerificationTokensResponse{}, err
	}
	return (*instance.EmailVerificationRecipe.RecipeImpl.RevokeEmailVerificationTokens)(userID, email, userContext)
}

func UnverifyEmail(userID string, userContext supertokens.UserContext) (evmodels.UnverifyEmailResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	email, err := instance.getEmailForUserIdForEmailVerification(userID, userContext)
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	return (*instance.EmailVerificationRecipe.RecipeImpl.UnverifyEmail)(userID, email, userContext)
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

func ConsumeCodeWithUserInputCode(deviceID string, userInputCode string, preAuthSessionID string, userContext supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ConsumeCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ConsumeCode)(&plessmodels.UserInputCodeWithDeviceID{
		Code:     userInputCode,
		DeviceID: deviceID,
	}, nil, preAuthSessionID, userContext)
}

func ConsumeCodeWithLinkCode(linkCode string, preAuthSessionID string, userContext supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ConsumeCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ConsumeCode)(nil, &linkCode, preAuthSessionID, userContext)
}

func GetUserByID(userID string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUserByPhoneNumber(phoneNumber string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByPhoneNumber)(phoneNumber, userContext)
}

func UpdatePasswordlessUser(userID string, email *string, phoneNumber *string, userContext supertokens.UserContext) (plessmodels.UpdateUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.UpdateUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.UpdatePasswordlessUser)(userID, email, phoneNumber, userContext)
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

func ListCodesByEmail(email string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByEmail)(email, userContext)
}

func ListCodesByPhoneNumber(phoneNumber string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByPhoneNumber)(phoneNumber, userContext)
}

func ListCodesByDeviceID(deviceID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByDeviceID)(deviceID, userContext)
}

func ListCodesByPreAuthSessionID(preAuthSessionID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByPreAuthSessionID)(preAuthSessionID, userContext)
}

func CreateMagicLinkByEmail(email string, userContext supertokens.UserContext) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.passwordlessRecipe.CreateMagicLink(&email, nil, userContext)
}

func CreateMagicLinkByPhoneNumber(phoneNumber string, userContext supertokens.UserContext) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.passwordlessRecipe.CreateMagicLink(nil, &phoneNumber, userContext)
}

func SignInUpByEmail(email string, userContext supertokens.UserContext) (struct {
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
	return instance.passwordlessRecipe.SignInUp(&email, nil, userContext)
}

func SignInUpByPhoneNumber(phoneNumber string, userContext supertokens.UserContext) (struct {
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
	return instance.passwordlessRecipe.SignInUp(nil, &phoneNumber, userContext)
}
