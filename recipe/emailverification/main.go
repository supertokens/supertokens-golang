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
	"github.com/supertokens/supertokens-golang/recipe/emailverification/api"
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

func CreateEmailVerificationLink(tenantId string, userID string, email *string, userContext ...supertokens.UserContext) (evmodels.CreateEmailVerificationLinkResponse, error) {
	st, err := supertokens.GetInstanceOrThrowError()
	if err != nil {
		return evmodels.CreateEmailVerificationLinkResponse{}, err
	}

	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	emailVerificationTokenResponse, err := CreateEmailVerificationToken(tenantId, userID, email, userContext...)
	if err != nil {
		return evmodels.CreateEmailVerificationLinkResponse{}, err
	}

	if emailVerificationTokenResponse.EmailAlreadyVerifiedError != nil {
		return evmodels.CreateEmailVerificationLinkResponse{
			EmailAlreadyVerifiedError: &struct{}{},
		}, nil
	}

	link, err := api.GetEmailVerifyLink(st.AppInfo, emailVerificationTokenResponse.OK.Token, tenantId, supertokens.GetRequestFromUserContext(userContext[0]), userContext[0])

	if err != nil {
		return evmodels.CreateEmailVerificationLinkResponse{}, err
	}

	return evmodels.CreateEmailVerificationLinkResponse{
		OK: &struct{ Link string }{
			Link: link,
		},
	}, nil
}

func SendEmailVerificationEmail(tenantId string, userID string, email *string, userContext ...supertokens.UserContext) (evmodels.SendEmailVerificationLinkResponse, error) {
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	if email == nil {
		instance, err := getRecipeInstanceOrThrowError()
		if err != nil {
			return evmodels.SendEmailVerificationLinkResponse{}, err
		}
		emailInfo, err := instance.GetEmailForUserID(userID, userContext[0])
		if err != nil {
			return evmodels.SendEmailVerificationLinkResponse{}, err
		}
		if emailInfo.EmailDoesNotExistError != nil {
			return evmodels.SendEmailVerificationLinkResponse{
				EmailAlreadyVerifiedError: &struct{}{},
			}, nil
		}

		if emailInfo.UnknownUserIDError != nil {
			return evmodels.SendEmailVerificationLinkResponse{}, errors.New("unknown user id provided without email")
		}

		email = &emailInfo.OK.Email
	}

	emailVerificationLinkResponse, err := CreateEmailVerificationLink(tenantId, userID, email, userContext...)
	if err != nil {
		return evmodels.SendEmailVerificationLinkResponse{}, err
	}

	if emailVerificationLinkResponse.EmailAlreadyVerifiedError != nil {
		return evmodels.SendEmailVerificationLinkResponse{
			EmailAlreadyVerifiedError: &struct{}{},
		}, nil
	}

	err = SendEmail(emaildelivery.EmailType{
		EmailVerification: &emaildelivery.EmailVerificationType{
			User: emaildelivery.User{
				ID:    userID,
				Email: *email,
			},
			EmailVerifyLink: emailVerificationLinkResponse.OK.Link,
			TenantId:        tenantId,
		},
	}, userContext...)
	if err != nil {
		return evmodels.SendEmailVerificationLinkResponse{}, err
	}
	return evmodels.SendEmailVerificationLinkResponse{
		OK: &struct{}{},
	}, nil
}

func MakeSMTPService(config emaildelivery.SMTPServiceConfig) *emaildelivery.EmailDeliveryInterface {
	return smtpService.MakeSMTPService(config)
}
