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
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/api"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *tpepmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func ThirdPartyManuallyCreateOrUpdateUser(tenantId string, thirdPartyID string, thirdPartyUserID string, email string, userContext ...supertokens.UserContext) (tpepmodels.ManuallyCreateOrUpdateUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tpepmodels.ManuallyCreateOrUpdateUserResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ThirdPartyManuallyCreateOrUpdateUser)(thirdPartyID, thirdPartyUserID, email, tenantId, userContext[0])
}

func ThirdPartyGetProvider(tenantId string, thirdPartyID string, clientType *string, userContext ...supertokens.UserContext) (*tpmodels.TypeProvider, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ThirdPartyGetProvider)(thirdPartyID, clientType, tenantId, userContext[0])
}

func GetUserByThirdPartyInfo(tenantId string, thirdPartyID string, thirdPartyUserID string, userContext ...supertokens.UserContext) (*tpepmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, tenantId, userContext[0])
}

func EmailPasswordSignUp(tenantId string, email, password string, userContext ...supertokens.UserContext) (tpepmodels.SignUpResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tpepmodels.SignUpResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.EmailPasswordSignUp)(email, password, tenantId, userContext[0])
}

func EmailPasswordSignIn(tenantId string, email, password string, userContext ...supertokens.UserContext) (tpepmodels.SignInResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tpepmodels.SignInResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.EmailPasswordSignIn)(email, password, tenantId, userContext[0])
}

func GetUserById(userID string, userContext ...supertokens.UserContext) (*tpepmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext[0])
}

func GetUsersByEmail(tenantId string, email string, userContext ...supertokens.UserContext) ([]tpepmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, tenantId, userContext[0])
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

func ResetPasswordUsingToken(tenantId string, token, newPassword string, userContext ...supertokens.UserContext) (epmodels.ResetPasswordUsingTokenResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.ResetPasswordUsingTokenResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ResetPasswordUsingToken)(token, newPassword, tenantId, userContext[0])
}

func UpdateEmailOrPassword(userId string, email *string, password *string, applyPasswordPolicy *bool, tenantIdForPasswordPolicy *string, userContext ...supertokens.UserContext) (epmodels.UpdateEmailOrPasswordResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.UpdateEmailOrPasswordResponse{}, err
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

func CreateResetPasswordLink(tenantId string, userID string, userContext ...supertokens.UserContext) (epmodels.CreateResetPasswordLinkResponse, error) {
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	tokenResponse, err := CreateResetPasswordToken(tenantId, userID, userContext...)
	if err != nil {
		return epmodels.CreateResetPasswordLinkResponse{}, err
	}
	if tokenResponse.UnknownUserIdError != nil {
		return epmodels.CreateResetPasswordLinkResponse{
			UnknownUserIdError: &struct{}{},
		}, nil
	}

	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.CreateResetPasswordLinkResponse{}, err
	}

	link, err := api.GetPasswordResetLink(
		instance.RecipeModule.GetAppInfo(),
		instance.RecipeModule.GetRecipeID(),
		tokenResponse.OK.Token,
		tenantId,
		supertokens.GetRequestFromUserContext(userContext[0]),
		userContext[0],
	)

	if err != nil {
		return epmodels.CreateResetPasswordLinkResponse{}, err
	}

	return epmodels.CreateResetPasswordLinkResponse{
		OK: &struct{ Link string }{
			Link: link,
		},
	}, nil
}

func SendResetPasswordEmail(tenantId string, userID string, userContext ...supertokens.UserContext) (epmodels.SendResetPasswordEmailResponse, error) {
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	linkResponse, err := CreateResetPasswordLink(tenantId, userID, userContext...)
	if err != nil {
		return epmodels.SendResetPasswordEmailResponse{}, err
	}
	if linkResponse.UnknownUserIdError != nil {
		return epmodels.SendResetPasswordEmailResponse{
			UnknownUserIdError: &struct{}{},
		}, nil
	}

	userInfo, err := GetUserById(userID, userContext...)
	if err != nil {
		return epmodels.SendResetPasswordEmailResponse{}, err
	}
	if userInfo == nil {
		return epmodels.SendResetPasswordEmailResponse{
			UnknownUserIdError: &struct{}{},
		}, nil
	}
	err = SendEmail(emaildelivery.EmailType{
		PasswordReset: &emaildelivery.PasswordResetType{
			User: emaildelivery.User{
				ID:    userInfo.ID,
				Email: userInfo.Email,
			},
			PasswordResetLink: linkResponse.OK.Link,
			TenantId:          tenantId,
		},
	}, userContext...)
	if err != nil {
		return epmodels.SendResetPasswordEmailResponse{}, err
	}

	return epmodels.SendResetPasswordEmailResponse{
		OK: &struct{}{},
	}, nil
}

func MakeSMTPService(config emaildelivery.SMTPServiceConfig) *emaildelivery.EmailDeliveryInterface {
	return smtpService.MakeSMTPService(config)
}
