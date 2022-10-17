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

package thirdpartyemailpassword

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *tpepmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func ThirdPartyCreateOrUpdateUserWithContext(thirdPartyID string, thirdPartyUserID string, email string, userContext supertokens.UserContext) (tpepmodels.ThirdPartyManuallyCreateOrUpdateUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpepmodels.ThirdPartyManuallyCreateOrUpdateUserResponse{}, err
	}
	return (*instance.RecipeImpl.ThirdPartyManuallyCreateOrUpdateUser)(thirdPartyID, thirdPartyUserID, email, userContext)
}

func GetUserByThirdPartyInfoWithContext(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*tpepmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, userContext)
}

func EmailPasswordSignUpWithContext(email, password string, userContext supertokens.UserContext) (tpepmodels.SignUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpepmodels.SignUpResponse{}, err
	}
	return (*instance.RecipeImpl.EmailPasswordSignUp)(email, password, userContext)
}

func EmailPasswordSignInWithContext(email, password string, userContext supertokens.UserContext) (tpepmodels.SignInResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpepmodels.SignInResponse{}, err
	}
	return (*instance.RecipeImpl.EmailPasswordSignIn)(email, password, userContext)
}

func GetUserByIdWithContext(userID string, userContext supertokens.UserContext) (*tpepmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUsersByEmailWithContext(email string, userContext supertokens.UserContext) ([]tpepmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, userContext)
}

func CreateResetPasswordTokenWithContext(userID string, userContext supertokens.UserContext) (epmodels.CreateResetPasswordTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.CreateResetPasswordTokenResponse{}, err
	}
	return (*instance.RecipeImpl.CreateResetPasswordToken)(userID, userContext)
}

func ResetPasswordUsingTokenWithContext(token, newPassword string, userContext supertokens.UserContext) (epmodels.ResetPasswordUsingTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.ResetPasswordUsingTokenResponse{}, err
	}
	return (*instance.RecipeImpl.ResetPasswordUsingToken)(token, newPassword, userContext)
}

func UpdateEmailOrPasswordWithContext(userId string, email *string, password *string, userContext supertokens.UserContext) (epmodels.UpdateEmailOrPasswordResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.UpdateEmailOrPasswordResponse{}, err
	}
	return (*instance.RecipeImpl.UpdateEmailOrPassword)(userId, email, password, userContext)
}

func SendEmailWithContext(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(input, userContext)
}

func ThirdPartyManuallyCreateOrUpdateUser(thirdPartyID string, thirdPartyUserID string, email string) (tpepmodels.ThirdPartyManuallyCreateOrUpdateUserResponse, error) {
	return ThirdPartyCreateOrUpdateUserWithContext(thirdPartyID, thirdPartyUserID, email, &map[string]interface{}{})
}

func GetUserByThirdPartyInfo(thirdPartyID string, thirdPartyUserID string) (*tpepmodels.User, error) {
	return GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID, &map[string]interface{}{})
}

func EmailPasswordSignUp(email, password string) (tpepmodels.SignUpResponse, error) {
	return EmailPasswordSignUpWithContext(email, password, &map[string]interface{}{})
}

func EmailPasswordSignIn(email, password string) (tpepmodels.SignInResponse, error) {
	return EmailPasswordSignInWithContext(email, password, &map[string]interface{}{})
}

func GetUserById(userID string) (*tpepmodels.User, error) {
	return GetUserByIdWithContext(userID, &map[string]interface{}{})
}

func GetUsersByEmail(email string) ([]tpepmodels.User, error) {
	return GetUsersByEmailWithContext(email, &map[string]interface{}{})
}

func CreateResetPasswordToken(userID string) (epmodels.CreateResetPasswordTokenResponse, error) {
	return CreateResetPasswordTokenWithContext(userID, &map[string]interface{}{})
}

func ResetPasswordUsingToken(token, newPassword string) (epmodels.ResetPasswordUsingTokenResponse, error) {
	return ResetPasswordUsingTokenWithContext(token, newPassword, &map[string]interface{}{})
}

func UpdateEmailOrPassword(userId string, email *string, password *string) (epmodels.UpdateEmailOrPasswordResponse, error) {
	return UpdateEmailOrPasswordWithContext(userId, email, password, &map[string]interface{}{})
}

func SendEmail(input emaildelivery.EmailType) error {
	return SendEmailWithContext(input, &map[string]interface{}{})
}

func MakeSMTPService(config emaildelivery.SMTPServiceConfig) *emaildelivery.EmailDeliveryInterface {
	return smtpService.MakeSMTPService(config)
}
