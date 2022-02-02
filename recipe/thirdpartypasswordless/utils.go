// /* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
//  *
//  * This software is licensed under the Apache License, Version 2.0 (the
//  * "License") as published by the Apache Software Foundation.
//  *
//  * You may not use this file except in compliance with the License. You may
//  * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//  *
//  * Unless required by applicable law or agreed to in writing, software
//  * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  * License for the specific language governing permissions and limitations
//  * under the License.
//  */

package thirdpartypasswordless

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance *Recipe, appInfo supertokens.NormalisedAppinfo, config tplmodels.TypeInput) (tplmodels.TypeNormalisedInput, error) {
	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance, config)

	typeNormalisedInput.EmailVerificationFeature = validateAndNormaliseEmailVerificationConfig(recipeInstance, config)

	if config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
		if config.Override.EmailVerificationFeature != nil {
			typeNormalisedInput.Override.EmailVerificationFeature = config.Override.EmailVerificationFeature
		}
	}

	return typeNormalisedInput, nil
}

func makeTypeNormalisedInput(recipeInstance *Recipe, inputConfig tplmodels.TypeInput) tplmodels.TypeNormalisedInput {
	return tplmodels.TypeNormalisedInput{
		Providers:                 inputConfig.Providers,
		ContactMethodPhone:        inputConfig.ContactMethodPhone,
		ContactMethodEmail:        inputConfig.ContactMethodEmail,
		ContactMethodEmailOrPhone: inputConfig.ContactMethodEmailOrPhone,
		FlowType:                  inputConfig.FlowType,
		GetLinkDomainAndPath:      inputConfig.GetLinkDomainAndPath,
		GetCustomUserInputCode:    inputConfig.GetCustomUserInputCode,
		EmailVerificationFeature:  validateAndNormaliseEmailVerificationConfig(recipeInstance, inputConfig),
		Override: tplmodels.OverrideStruct{
			Functions: func(originalImplementation tplmodels.RecipeInterface) tplmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation tplmodels.APIInterface) tplmodels.APIInterface {
				return originalImplementation
			},
			EmailVerificationFeature: nil,
		},
	}
}

func validateAndNormaliseEmailVerificationConfig(recipeInstance *Recipe, config tplmodels.TypeInput) evmodels.TypeInput {
	emailverificationTypeInput := evmodels.TypeInput{
		GetEmailForUserID: recipeInstance.getEmailForUserIdForEmailVerification,
		Override:          nil,
	}

	if config.Override != nil {
		emailverificationTypeInput.Override = config.Override.EmailVerificationFeature
	}

	if config.EmailVerificationFeature != nil {
		if config.EmailVerificationFeature.CreateAndSendCustomEmail != nil {
			emailverificationTypeInput.CreateAndSendCustomEmail = func(user evmodels.User, link string, userContext supertokens.UserContext) {
				userInfo, err := (*recipeInstance.RecipeImpl.GetUserByID)(user.ID, userContext)
				if err != nil {
					return
				}
				if userInfo == nil {
					return
				}
				config.EmailVerificationFeature.CreateAndSendCustomEmail(*userInfo, link, userContext)
			}
		}

		if config.EmailVerificationFeature.GetEmailVerificationURL != nil {
			emailverificationTypeInput.GetEmailVerificationURL = func(user evmodels.User, userContext supertokens.UserContext) (string, error) {
				userInfo, err := (*recipeInstance.RecipeImpl.GetUserByID)(user.ID, userContext)
				if err != nil {
					return "", err
				}
				if userInfo == nil {
					return "", errors.New("Unknown User ID provided")
				}
				return config.EmailVerificationFeature.GetEmailVerificationURL(*userInfo, userContext)
			}
		}
	}

	return emailverificationTypeInput
}
