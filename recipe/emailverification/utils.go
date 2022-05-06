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

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/emaildeliverymodels"
	backwardcompatibility "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/services/backwardCompatibility"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config evmodels.TypeInput) evmodels.TypeNormalisedInput {
	typeNormalisedInput := makeTypeNormalisedInput(appInfo)

	if config.GetEmailVerificationURL != nil {
		typeNormalisedInput.GetEmailVerificationURL = config.GetEmailVerificationURL
	}

	typeNormalisedInput.GetEmailDeliveryConfig = func() emaildeliverymodels.TypeInputWithService {
		createAndSendCustomEmail := DefaultCreateAndSendCustomEmail(appInfo)
		if config.CreateAndSendCustomEmail != nil {
			createAndSendCustomEmail = config.CreateAndSendCustomEmail
		}
		emailService := backwardcompatibility.MakeBackwardCompatibilityService(appInfo, createAndSendCustomEmail)
		if config.EmailDelivery != nil && config.EmailDelivery.Service != nil {
			emailService = *config.EmailDelivery.Service
		}
		result := emaildeliverymodels.TypeInputWithService{
			Service: emailService,
		}
		if config.EmailDelivery != nil && config.EmailDelivery.Override != nil {
			result.Override = config.EmailDelivery.Override
		}
		return result
	}

	if config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
	}

	if config.GetEmailForUserID != nil {
		typeNormalisedInput.GetEmailForUserID = config.GetEmailForUserID
	}
	return typeNormalisedInput
}

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) evmodels.TypeNormalisedInput {
	return evmodels.TypeNormalisedInput{
		GetEmailForUserID: func(userID string, userContext supertokens.UserContext) (string, error) {
			return "", errors.New("not defined by user")
		},
		GetEmailVerificationURL: DefaultGetEmailVerificationURL(appInfo),
		GetEmailDeliveryConfig:  nil,
		Override: evmodels.OverrideStruct{
			Functions: func(originalImplementation evmodels.RecipeInterface) evmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation evmodels.APIInterface) evmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}
