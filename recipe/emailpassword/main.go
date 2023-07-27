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

func SignUp(tenantId string, email string, password string, userContext ...supertokens.UserContext) (epmodels.SignUpResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.SignUpResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.SignUp)(email, password, tenantId, userContext[0])
}

func SignIn(tenantId string, email string, password string, userContext ...supertokens.UserContext) (epmodels.SignInResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.SignInResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.SignIn)(email, password, tenantId, userContext[0])
}

func GetUserByID(userID string, userContext ...supertokens.UserContext) (*epmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext[0])
}

func GetUserByEmail(tenantId string, email string, userContext ...supertokens.UserContext) (*epmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserByEmail)(email, tenantId, userContext[0])
}

func CreateResetPasswordToken(tenantId string, userID string, userContext ...supertokens.UserContext) (epmodels.CreateResetPasswordTokenResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.CreateResetPasswordTokenResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.CreateResetPasswordToken)(userID, tenantId, userContext[0])
}

func ResetPasswordUsingToken(tenantId string, token string, newPassword string, userContext ...supertokens.UserContext) (epmodels.ResetPasswordUsingTokenResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.ResetPasswordUsingTokenResponse{}, nil
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ResetPasswordUsingToken)(token, newPassword, tenantId, userContext[0])
}

func UpdateEmailOrPassword(userId string, email *string, password *string, applyPasswordPolicy *bool, tenantIdForPasswordPolicy *string, userContext ...supertokens.UserContext) (epmodels.UpdateEmailOrPasswordResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.UpdateEmailOrPasswordResponse{}, nil
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	if tenantIdForPasswordPolicy == nil {
		tenantId := supertokens.DefaultTenantId
		tenantIdForPasswordPolicy = &tenantId
	}
	return (*instance.RecipeImpl.UpdateEmailOrPassword)(userId, email, password, applyPasswordPolicy, *tenantIdForPasswordPolicy, userContext[0])
}

func SendEmail(input emaildelivery.EmailType, userContext ...supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(input, userContext[0])
}

func MakeSMTPService(config emaildelivery.SMTPServiceConfig) *emaildelivery.EmailDeliveryInterface {
	return smtpService.MakeSMTPService(config)
}
