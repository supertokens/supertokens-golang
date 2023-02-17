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

func SignUpWithContext(email string, password string, tenantId *string, userContext supertokens.UserContext) (epmodels.SignUpResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.SignUpResponse{}, err
	}
	return (*instance.RecipeImpl.SignUp)(email, password, tenantId, userContext)
}

func SignInWithContext(email string, password string, tenantId *string, userContext supertokens.UserContext) (epmodels.SignInResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.SignInResponse{}, err
	}
	return (*instance.RecipeImpl.SignIn)(email, password, tenantId, userContext)
}

func GetUserByIDWithContext(userID string, tenantId *string, userContext supertokens.UserContext) (*epmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, tenantId, userContext)
}

func GetUserByEmailWithContext(email string, tenantId *string, userContext supertokens.UserContext) (*epmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByEmail)(email, tenantId, userContext)
}

func CreateResetPasswordTokenWithContext(userID string, tenantId *string, userContext supertokens.UserContext) (epmodels.CreateResetPasswordTokenResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.CreateResetPasswordTokenResponse{}, err
	}
	return (*instance.RecipeImpl.CreateResetPasswordToken)(userID, tenantId, userContext)
}

func ResetPasswordUsingTokenWithContext(token string, newPassword string, tenantId *string, userContext supertokens.UserContext) (epmodels.ResetPasswordUsingTokenResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.ResetPasswordUsingTokenResponse{}, nil
	}
	return (*instance.RecipeImpl.ResetPasswordUsingToken)(token, newPassword, tenantId, userContext)
}

func UpdateEmailOrPasswordWithContext(userId string, email *string, password *string, tenantId *string, userContext supertokens.UserContext) (epmodels.UpdateEmailOrPasswordResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.UpdateEmailOrPasswordResponse{}, nil
	}
	return (*instance.RecipeImpl.UpdateEmailOrPassword)(userId, email, password, tenantId, userContext)
}

func SendEmailWithContext(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(input, userContext)
}

func SignUp(email string, password string, tenantId *string) (epmodels.SignUpResponse, error) {
	return SignUpWithContext(email, password, tenantId, &map[string]interface{}{})
}

func SignIn(email string, password string, tenantId *string) (epmodels.SignInResponse, error) {
	return SignInWithContext(email, password, tenantId, &map[string]interface{}{})
}

func GetUserByID(userID string, tenantId *string) (*epmodels.User, error) {
	return GetUserByIDWithContext(userID, tenantId, &map[string]interface{}{})
}

func GetUserByEmail(email string, tenantId *string) (*epmodels.User, error) {
	return GetUserByEmailWithContext(email, tenantId, &map[string]interface{}{})
}

func CreateResetPasswordToken(userID string, tenantId *string) (epmodels.CreateResetPasswordTokenResponse, error) {
	return CreateResetPasswordTokenWithContext(userID, tenantId, &map[string]interface{}{})
}

func ResetPasswordUsingToken(token string, newPassword string, tenantId *string) (epmodels.ResetPasswordUsingTokenResponse, error) {
	return ResetPasswordUsingTokenWithContext(token, newPassword, tenantId, &map[string]interface{}{})
}

func UpdateEmailOrPassword(userId string, email *string, password *string, tenantId *string) (epmodels.UpdateEmailOrPasswordResponse, error) {
	return UpdateEmailOrPasswordWithContext(userId, email, password, tenantId, &map[string]interface{}{})
}

func SendEmail(input emaildelivery.EmailType) error {
	return SendEmailWithContext(input, &map[string]interface{}{})
}

func MakeSMTPService(config emaildelivery.SMTPServiceConfig) *emaildelivery.EmailDeliveryInterface {
	return smtpService.MakeSMTPService(config)
}
