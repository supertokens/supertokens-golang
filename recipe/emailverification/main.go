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

package emailverification

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config evmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func CreateEmailVerificationToken(tenantId string, userID string, email *string, userContext ...supertokens.UserContext) (evmodels.CreateEmailVerificationTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	if email == nil {
		emailInfo, err := instance.GetEmailForUserID(userID, userContext[0])
		if err != nil {
			return evmodels.CreateEmailVerificationTokenResponse{}, err
		}
		if emailInfo.OK != nil {
			email = &emailInfo.OK.Email
		} else if emailInfo.EmailDoesNotExistError != nil {
			return evmodels.CreateEmailVerificationTokenResponse{
				EmailAlreadyVerifiedError: &struct{}{},
			}, nil
		} else {
			return evmodels.CreateEmailVerificationTokenResponse{}, errors.New("unknown user id provided without email")
		}
	}
	return (*instance.RecipeImpl.CreateEmailVerificationToken)(userID, *email, tenantId, userContext[0])
}

func VerifyEmailUsingToken(tenantId string, token string, userContext ...supertokens.UserContext) (evmodels.VerifyEmailUsingTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.VerifyEmailUsingTokenResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.VerifyEmailUsingToken)(token, tenantId, userContext[0])
}

func IsEmailVerified(userID string, email *string, userContext ...supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	if email == nil {
		emailInfo, err := instance.GetEmailForUserID(userID, userContext[0])
		if err != nil {
			return false, err
		}
		if emailInfo.OK != nil {
			email = &emailInfo.OK.Email
		} else if emailInfo.EmailDoesNotExistError != nil {
			return true, nil
		} else {
			return false, errors.New("unknown user id provided without email")
		}
	}
	return (*instance.RecipeImpl.IsEmailVerified)(userID, *email, userContext[0])
}

func RevokeEmailVerificationTokens(tenantId string, userID string, email *string, userContext ...supertokens.UserContext) (evmodels.RevokeEmailVerificationTokensResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.RevokeEmailVerificationTokensResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	if email == nil {
		emailInfo, err := instance.GetEmailForUserID(userID, userContext[0])
		if err != nil {
			return evmodels.RevokeEmailVerificationTokensResponse{}, err
		}
		if emailInfo.OK != nil {
			email = &emailInfo.OK.Email
		} else if emailInfo.EmailDoesNotExistError != nil {
			return evmodels.RevokeEmailVerificationTokensResponse{
				OK: &struct{}{},
			}, nil
		} else {
			return evmodels.RevokeEmailVerificationTokensResponse{}, errors.New("unknown user id provided without email")
		}
	}
	return (*instance.RecipeImpl.RevokeEmailVerificationTokens)(userID, *email, tenantId, userContext[0])
}

func UnverifyEmail(userID string, email *string, userContext ...supertokens.UserContext) (evmodels.UnverifyEmailResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	if email == nil {
		emailInfo, err := instance.GetEmailForUserID(userID, userContext[0])
		if err != nil {
			return evmodels.UnverifyEmailResponse{}, err
		}
		if emailInfo.OK != nil {
			email = &emailInfo.OK.Email
		} else if emailInfo.EmailDoesNotExistError != nil {
			return evmodels.UnverifyEmailResponse{
				OK: &struct{}{},
			}, nil
		} else {
			return evmodels.UnverifyEmailResponse{}, errors.New("unknown user id provided without email")
		}
	}
	return (*instance.RecipeImpl.UnverifyEmail)(userID, *email, userContext[0])
}

func SendEmail(input emaildelivery.EmailType, userContext ...supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
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
