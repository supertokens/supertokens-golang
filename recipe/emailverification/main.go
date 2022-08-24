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

func CreateEmailVerificationTokenWithContext(userID string, email *string, userContext supertokens.UserContext) (evmodels.CreateEmailVerificationTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	if email == nil {
		emailInfo, err := instance.GetEmailForUserID(userID, userContext)
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
	return (*instance.RecipeImpl.CreateEmailVerificationToken)(userID, *email, userContext)
}

func VerifyEmailUsingTokenWithContext(token string, userContext supertokens.UserContext) (evmodels.VerifyEmailUsingTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.VerifyEmailUsingTokenResponse{}, err
	}
	return (*instance.RecipeImpl.VerifyEmailUsingToken)(token, userContext)
}

func IsEmailVerifiedWithContext(userID string, email *string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	if email == nil {
		emailInfo, err := instance.GetEmailForUserID(userID, userContext)
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
	return (*instance.RecipeImpl.IsEmailVerified)(userID, *email, userContext)
}

func RevokeEmailVerificationTokensWithContext(userID string, email *string, userContext supertokens.UserContext) (evmodels.RevokeEmailVerificationTokensResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.RevokeEmailVerificationTokensResponse{}, err
	}
	if email == nil {
		emailInfo, err := instance.GetEmailForUserID(userID, userContext)
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
	return (*instance.RecipeImpl.RevokeEmailVerificationTokens)(userID, *email, userContext)
}

func UnverifyEmailWithContext(userID string, email *string, userContext supertokens.UserContext) (evmodels.UnverifyEmailResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	if email == nil {
		emailInfo, err := instance.GetEmailForUserID(userID, userContext)
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
	return (*instance.RecipeImpl.UnverifyEmail)(userID, *email, userContext)
}

func SendEmailWithContext(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(input, userContext)
}

func CreateEmailVerificationToken(userID string, email *string) (evmodels.CreateEmailVerificationTokenResponse, error) {
	return CreateEmailVerificationTokenWithContext(userID, email, &map[string]interface{}{})
}

func VerifyEmailUsingToken(token string) (evmodels.VerifyEmailUsingTokenResponse, error) {
	return VerifyEmailUsingTokenWithContext(token, &map[string]interface{}{})
}

func IsEmailVerified(userID string, email *string) (bool, error) {
	return IsEmailVerifiedWithContext(userID, email, &map[string]interface{}{})
}

func RevokeEmailVerificationTokens(userID string, email *string) (evmodels.RevokeEmailVerificationTokensResponse, error) {
	return RevokeEmailVerificationTokensWithContext(userID, email, &map[string]interface{}{})
}

func UnverifyEmail(userID string, email *string) (evmodels.UnverifyEmailResponse, error) {
	return UnverifyEmailWithContext(userID, email, &map[string]interface{}{})
}

func SendEmail(input emaildelivery.EmailType) error {
	return SendEmailWithContext(input, &map[string]interface{}{})
}

func MakeSMTPService(config emaildelivery.SMTPServiceConfig) *emaildelivery.EmailDeliveryInterface {
	return smtpService.MakeSMTPService(config)
}
