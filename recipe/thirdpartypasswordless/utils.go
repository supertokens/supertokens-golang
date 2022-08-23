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

package thirdpartypasswordless

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/emaildelivery/backwardCompatibilityService"
	smsBackwardCompatibilityService "github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/smsdelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance *Recipe, appInfo supertokens.NormalisedAppinfo, config tplmodels.TypeInput) (tplmodels.TypeNormalisedInput, error) {
	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance, config)

	typeNormalisedInput.GetEmailDeliveryConfig = func() emaildelivery.TypeInputWithService {
		sendPasswordlessLoginEmail := passwordless.DefaultCreateAndSendCustomEmail(appInfo)
		if config.ContactMethodEmail.Enabled && config.ContactMethodEmail.CreateAndSendCustomEmail != nil {
			sendPasswordlessLoginEmail = config.ContactMethodEmail.CreateAndSendCustomEmail
		} else if config.ContactMethodEmailOrPhone.Enabled && config.ContactMethodEmailOrPhone.CreateAndSendCustomEmail != nil {
			sendPasswordlessLoginEmail = config.ContactMethodEmailOrPhone.CreateAndSendCustomEmail
		}

		sendEmailVerificationEmail := emailverification.DefaultCreateAndSendCustomEmail(appInfo)

		emailService := backwardCompatibilityService.MakeBackwardCompatibilityService(appInfo, sendEmailVerificationEmail, sendPasswordlessLoginEmail)
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

	typeNormalisedInput.GetSmsDeliveryConfig = func() smsdelivery.TypeInputWithService {
		sendPasswordlessLoginSms := passwordless.DefaultCreateAndSendCustomTextMessage(appInfo)

		if config.ContactMethodPhone.Enabled && config.ContactMethodPhone.CreateAndSendCustomTextMessage != nil {
			sendPasswordlessLoginSms = config.ContactMethodPhone.CreateAndSendCustomTextMessage
		} else if config.ContactMethodEmailOrPhone.Enabled && config.ContactMethodEmailOrPhone.CreateAndSendCustomTextMessage != nil {
			sendPasswordlessLoginSms = config.ContactMethodEmailOrPhone.CreateAndSendCustomTextMessage
		}

		smsService := smsBackwardCompatibilityService.MakeBackwardCompatibilityService(sendPasswordlessLoginSms)
		if config.SmsDelivery != nil && config.SmsDelivery.Service != nil {
			smsService = *config.SmsDelivery.Service
		}
		result := smsdelivery.TypeInputWithService{
			Service: smsService,
		}
		if config.SmsDelivery != nil && config.SmsDelivery.Override != nil {
			result.Override = config.SmsDelivery.Override
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
		Override: tplmodels.OverrideStruct{
			Functions: func(originalImplementation tplmodels.RecipeInterface) tplmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation tplmodels.APIInterface) tplmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}
