/*
 * Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package backwardCompatibilityService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	emailVerificationBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(recipeInterfaceImpl epmodels.RecipeInterface, appInfo supertokens.NormalisedAppinfo, sendResetPasswordEmail func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext), sendEmailVerificationEmail func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext)) emaildelivery.EmailDeliveryInterface {

	emailVerificationService := emailVerificationBackwardsCompatibilityService.MakeBackwardCompatibilityService(appInfo, sendEmailVerificationEmail)

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			return (*emailVerificationService.SendEmail)(input, userContext)
		} else if input.PasswordReset != nil {
			user, err := (*recipeInterfaceImpl.GetUserByID)(input.PasswordReset.User.ID, userContext)
			if err != nil {
				return err
			}
			if user == nil {
				return errors.New("should never come here")
			}
			sendResetPasswordEmail(*user, input.PasswordReset.PasswordResetLink, userContext)
		} else {
			return errors.New("should never come here")
		}
		return nil
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
