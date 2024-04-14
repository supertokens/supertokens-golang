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

package api

import (
	"errors"
	"fmt"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evclaims"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	sessErrors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() evmodels.APIInterface {
	verifyEmailPOST := func(token string, sessionContainer sessmodels.SessionContainer, tenantId string, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.VerifyEmailPOSTResponse, error) {
		resp, err := (*options.RecipeImplementation.VerifyEmailUsingToken)(token, tenantId, userContext)
		if err != nil {
			return evmodels.VerifyEmailPOSTResponse{}, err
		}
		if resp.OK != nil {
			if sessionContainer != nil {
				err := sessionContainer.FetchAndSetClaimWithContext(evclaims.EmailVerificationClaim, userContext)
				if err != nil {
					if err.Error() == "UNKNOWN_USER_ID" {
						supertokens.LogDebugMessage("verifyEmailPOST: Returning UnauthorizedError because the User Id provided is unknown")
						return evmodels.VerifyEmailPOSTResponse{}, sessErrors.UnauthorizedError{Msg: "Unknown User ID provided"}
					}
					return evmodels.VerifyEmailPOSTResponse{}, err
				}
			}
			return evmodels.VerifyEmailPOSTResponse{
				OK: resp.OK,
			}, err
		} else {
			return evmodels.VerifyEmailPOSTResponse{
				EmailVerificationInvalidTokenError: resp.EmailVerificationInvalidTokenError,
			}, nil
		}
	}

	isEmailVerifiedGET := func(sessionContainer sessmodels.SessionContainer, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.IsEmailVerifiedGETResponse, error) {
		if sessionContainer == nil {
			return evmodels.IsEmailVerifiedGETResponse{}, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
		}

		err := sessionContainer.FetchAndSetClaimWithContext(evclaims.EmailVerificationClaim, userContext)
		if err != nil {
			if err.Error() == "UNKNOWN_USER_ID" {
				supertokens.LogDebugMessage("isEmailVerifiedGET: Returning UnauthorizedError because the User Id provided is unknown")
				return evmodels.IsEmailVerifiedGETResponse{}, sessErrors.UnauthorizedError{Msg: "Unknown User ID provided"}
			}
			return evmodels.IsEmailVerifiedGETResponse{}, err
		}

		isVerified := sessionContainer.GetClaimValueWithContext(evclaims.EmailVerificationClaim, userContext)
		if isVerified == nil {
			return evmodels.IsEmailVerifiedGETResponse{}, errors.New("should never come here: EmailVerificationClaim failed to set value")
		}
		return evmodels.IsEmailVerifiedGETResponse{
			OK: &struct{ IsVerified bool }{
				IsVerified: isVerified.(bool),
			},
		}, nil
	}

	generateEmailVerifyTokenPOST := func(sessionContainer sessmodels.SessionContainer, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.GenerateEmailVerifyTokenPOSTResponse, error) {
		if sessionContainer == nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
		}

		userID := sessionContainer.GetUserIDWithContext(userContext)
		email, err := options.GetEmailForUserID(userID, userContext)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}
		if email.UnknownUserIDError != nil {
			supertokens.LogDebugMessage("generateEmailVerifyTokenPOST: Returning UnauthorizedError because the User Id provided is unknown")
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, sessErrors.UnauthorizedError{Msg: "Unknown User ID provided"}
		}
		if email.EmailDoesNotExistError != nil {
			supertokens.LogDebugMessage(fmt.Sprintf("Email verification email not sent to user %s because it doesn't have an email address.", userID))
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{
				EmailAlreadyVerifiedError: &struct{}{},
			}, nil
		}
		response, err := (*options.RecipeImplementation.CreateEmailVerificationToken)(userID, email.OK.Email, sessionContainer.GetTenantIdWithContext(userContext), userContext)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}

		if response.EmailAlreadyVerifiedError != nil {
			if sessionContainer.GetClaimValue(evclaims.EmailVerificationClaim) != true {
				sessionContainer.FetchAndSetClaimWithContext(evclaims.EmailVerificationClaim, userContext)
			}
			supertokens.LogDebugMessage(fmt.Sprintf("Email verification email not sent to %s because it is already verified", email.OK.Email))
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{
				EmailAlreadyVerifiedError: &struct{}{},
			}, nil
		}

		if sessionContainer.GetClaimValue(evclaims.EmailVerificationClaim) != false {
			sessionContainer.FetchAndSetClaimWithContext(evclaims.EmailVerificationClaim, userContext)
		}

		user := evmodels.User{
			ID:    userID,
			Email: email.OK.Email,
		}

		emailVerificationURL, err := GetEmailVerifyLink(
			options.AppInfo,
			response.OK.Token,
			sessionContainer.GetTenantIdWithContext(userContext),
			options.Req,
			userContext,
		)

		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}

		supertokens.LogDebugMessage(fmt.Sprintf("Sending email verification email to %s", email.OK.Email))
		err = (*options.EmailDelivery.IngredientInterfaceImpl.SendEmail)(emaildelivery.EmailType{
			EmailVerification: &emaildelivery.EmailVerificationType{
				User: emaildelivery.User{
					ID:    user.ID,
					Email: user.Email,
				},
				EmailVerifyLink: emailVerificationURL,
				TenantId:        sessionContainer.GetTenantIdWithContext(userContext),
			},
		}, userContext)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}

		return evmodels.GenerateEmailVerifyTokenPOSTResponse{
			OK: &struct{}{},
		}, nil
	}

	return evmodels.APIInterface{
		VerifyEmailPOST:              &verifyEmailPOST,
		IsEmailVerifiedGET:           &isEmailVerifiedGET,
		GenerateEmailVerifyTokenPOST: &generateEmailVerifyTokenPOST,
	}
}
