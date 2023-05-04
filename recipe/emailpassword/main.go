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

package emailpassword

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *epmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func SignUpWithContext(email string, password string, userContext supertokens.UserContext) (epmodels.SignUpResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.SignUpResponse{}, err
	}
	return (*instance.RecipeImpl.SignUp)(email, password, userContext)
}

func SignInWithContext(email string, password string, userContext supertokens.UserContext) (epmodels.SignInResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.SignInResponse{}, err
	}
	return (*instance.RecipeImpl.SignIn)(email, password, userContext)
}

func GetUserByIDWithContext(userID string, userContext supertokens.UserContext) (*epmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUserByEmailWithContext(email string, userContext supertokens.UserContext) (*epmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByEmail)(email, userContext)
}

func CreateResetPasswordTokenWithContext(userID string, userContext supertokens.UserContext) (epmodels.CreateResetPasswordTokenResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.CreateResetPasswordTokenResponse{}, err
	}
	return (*instance.RecipeImpl.CreateResetPasswordToken)(userID, userContext)
}

func ResetPasswordUsingTokenWithContext(token string, newPassword string, userContext supertokens.UserContext) (epmodels.ResetPasswordUsingTokenResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.ResetPasswordUsingTokenResponse{}, nil
	}
	return (*instance.RecipeImpl.ResetPasswordUsingToken)(token, newPassword, userContext)
}

func UpdateEmailOrPasswordWithContext(userId string, email *string, password *string, applyPasswordPolicy *bool, userContext supertokens.UserContext) (epmodels.UpdateEmailOrPasswordResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.UpdateEmailOrPasswordResponse{}, nil
	}
	return (*instance.RecipeImpl.UpdateEmailOrPassword)(userId, email, password, applyPasswordPolicy, userContext)
}

func SendEmailWithContext(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(input, userContext)
}

func SignUp(email string, password string) (epmodels.SignUpResponse, error) {
	return SignUpWithContext(email, password, &map[string]interface{}{})
}

func SignIn(email string, password string) (epmodels.SignInResponse, error) {
	return SignInWithContext(email, password, &map[string]interface{}{})
}

func GetUserByID(userID string) (*epmodels.User, error) {
	return GetUserByIDWithContext(userID, &map[string]interface{}{})
}

func GetUserByEmail(email string) (*epmodels.User, error) {
	return GetUserByEmailWithContext(email, &map[string]interface{}{})
}

func CreateResetPasswordToken(userID string) (epmodels.CreateResetPasswordTokenResponse, error) {
	return CreateResetPasswordTokenWithContext(userID, &map[string]interface{}{})
}

func ResetPasswordUsingToken(token string, newPassword string) (epmodels.ResetPasswordUsingTokenResponse, error) {
	return ResetPasswordUsingTokenWithContext(token, newPassword, &map[string]interface{}{})
}

func UpdateEmailOrPassword(userId string, email *string, password *string, applyPasswordPolicy *bool) (epmodels.UpdateEmailOrPasswordResponse, error) {
	return UpdateEmailOrPasswordWithContext(userId, email, password, applyPasswordPolicy, &map[string]interface{}{})
}

func SendEmail(input emaildelivery.EmailType) error {
	return SendEmailWithContext(input, &map[string]interface{}{})
}

func MakeSMTPService(config emaildelivery.SMTPServiceConfig) *emaildelivery.EmailDeliveryInterface {
	return smtpService.MakeSMTPService(config)
}
