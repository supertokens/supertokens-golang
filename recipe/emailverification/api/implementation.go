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
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() evmodels.APIInterface {
	verifyEmailPOST := func(token string, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.VerifyEmailUsingTokenResponse, error) {
		return (*options.RecipeImplementation.VerifyEmailUsingToken)(token, userContext)
	}

	isEmailVerifiedGET := func(options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.IsEmailVerifiedGETResponse, error) {
		session, err := session.GetSessionWithContext(options.Req, options.Res, nil, userContext)
		if err != nil {
			return evmodels.IsEmailVerifiedGETResponse{}, err
		}
		if session == nil {
			return evmodels.IsEmailVerifiedGETResponse{}, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
		}

		userID := session.GetUserIDWithContext(userContext)

		email, err := options.Config.GetEmailForUserID(userID, userContext)
		if err != nil {
			return evmodels.IsEmailVerifiedGETResponse{}, err
		}
		isVerified, err := (*options.RecipeImplementation.IsEmailVerified)(userID, email, userContext)
		if err != nil {
			return evmodels.IsEmailVerifiedGETResponse{}, err
		}
		return evmodels.IsEmailVerifiedGETResponse{
			OK: &struct{ IsVerified bool }{
				IsVerified: isVerified,
			},
		}, nil
	}

	generateEmailVerifyTokenPOST := func(options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.GenerateEmailVerifyTokenPOSTResponse, error) {
		session, err := session.GetSessionWithContext(options.Req, options.Res, nil, userContext)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}
		if session == nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
		}

		userID := session.GetUserIDWithContext(userContext)
		email, err := options.Config.GetEmailForUserID(userID, userContext)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}
		response, err := (*options.RecipeImplementation.CreateEmailVerificationToken)(userID, email, userContext)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}

		if response.EmailAlreadyVerifiedError != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{
				EmailAlreadyVerifiedError: &struct{}{},
			}, nil
		}

		user := evmodels.User{
			ID:    userID,
			Email: email,
		}
		emailVerificationURL, err := options.Config.GetEmailVerificationURL(user, userContext)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}
		emailVerifyLink := emailVerificationURL + "?token=" + response.OK.Token + "&rid=" + options.RecipeID

		(*options.EmailDelivery.IngredientInterfaceImpl.SendEmail)(emaildelivery.EmailType{
			EmailVerification: &emaildelivery.EmailVerificationType{
				User: emaildelivery.User{
					ID:    user.ID,
					Email: user.Email,
				},
				EmailVerifyLink: emailVerifyLink,
			},
		}, userContext)

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
