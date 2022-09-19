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
	emailPasswordBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/emailpassword/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(recipeInterfaceImpl tpepmodels.RecipeInterface, emailPasswordRecipeInterfaceImpl epmodels.RecipeInterface, appInfo supertokens.NormalisedAppinfo, sendResetPasswordEmail func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext)) emaildelivery.EmailDeliveryInterface {
	// We are using evmodels.User as opposed to tpmodels.User because TypeInput of thirdparty accepts evmodels.TypeInput for EmailVerificationFeature
	// Similarly with epmodels.User as well
	emailPasswordService := emailPasswordBackwardsCompatibilityService.MakeBackwardCompatibilityService(emailPasswordRecipeInterfaceImpl, appInfo, sendResetPasswordEmail)

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.PasswordReset != nil {
			return (*emailPasswordService.SendEmail)(input, userContext)

		} else {
			return errors.New("should never come here")
		}
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
