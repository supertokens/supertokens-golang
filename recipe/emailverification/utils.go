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
	"github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config evmodels.TypeInput) (evmodels.TypeNormalisedInput, error) {
	if config.Mode != "REQUIRED" && config.Mode != "OPTIONAL" {
		return evmodels.TypeNormalisedInput{}, errors.New("mode must be either REQUIRED or OPTIONAL")
	}
	typeNormalisedInput := makeTypeNormalisedInput(appInfo)

	typeNormalisedInput.Mode = config.Mode
	typeNormalisedInput.GetEmailForUserID = config.GetEmailForUserID

	typeNormalisedInput.GetEmailDeliveryConfig = func() emaildelivery.TypeInputWithService {
		createAndSendCustomEmail := DefaultCreateAndSendCustomEmail(appInfo)
		if config.CreateAndSendCustomEmail != nil {
			createAndSendCustomEmail = config.CreateAndSendCustomEmail
		}
		emailService := backwardCompatibilityService.MakeBackwardCompatibilityService(appInfo, createAndSendCustomEmail)
		if config.EmailDelivery != nil && config.EmailDelivery.Service != nil {
			emailService = *config.EmailDelivery.Service
		}
		result := emaildelivery.TypeInputWithService{
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
	return typeNormalisedInput, nil
}

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) evmodels.TypeNormalisedInput {
	return evmodels.TypeNormalisedInput{
		GetEmailForUserID: func(userID string, userContext supertokens.UserContext) (evmodels.TypeEmailInfo, error) {
			return evmodels.TypeEmailInfo{}, errors.New("not defined by user")
		},
		GetEmailDeliveryConfig: nil,
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
