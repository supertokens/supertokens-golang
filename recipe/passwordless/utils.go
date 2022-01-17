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

package passwordless

import (
	"reflect"
	"regexp"

	"github.com/nyaruka/phonenumbers"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config plessmodels.TypeInput) plessmodels.TypeNormalisedInput {

	if config.FlowType != "USER_INPUT_CODE" && config.FlowType != "MAGIC_LINK" && config.FlowType != "USER_INPUT_CODE_AND_MAGIC_LINK" {
		panic("FlowType config must be provided and must be one of \"USER_INPUT_CODE\", \"MAGIC_LINK\" or \"USER_INPUT_CODE_AND_MAGIC_LINK\"")
	}

	contactMethodEnabledCounter := 0

	if config.ContactMethodEmail.Enabled {
		contactMethodEnabledCounter++
	}

	if config.ContactMethodPhone.Enabled {
		contactMethodEnabledCounter++
	}

	if config.ContactMethodEmailOrPhone.Enabled {
		contactMethodEnabledCounter++
	}

	if contactMethodEnabledCounter != 1 {
		panic("Please enable only one of ContactMethodEmail, ContactMethodPhone or ContactMethodEmailOrPhone")
	}

	typeNormalisedInput := makeTypeNormalisedInput(appInfo, config)

	if config.ContactMethodPhone.Enabled {
		if config.ContactMethodPhone.CreateAndSendCustomTextMessage == nil {
			panic("Please pass a function (ContactMethodPhone.CreateAndSendCustomTextMessage) to send text messages.")
		}
		typeNormalisedInput.ContactMethodPhone.Enabled = true
		if config.ContactMethodPhone.CreateAndSendCustomTextMessage != nil {
			typeNormalisedInput.ContactMethodPhone.CreateAndSendCustomTextMessage = config.ContactMethodPhone.CreateAndSendCustomTextMessage
		}
		if config.ContactMethodPhone.ValidatePhoneNumber != nil {
			typeNormalisedInput.ContactMethodPhone.ValidatePhoneNumber = config.ContactMethodPhone.ValidatePhoneNumber
		}
	}

	if config.ContactMethodEmail.Enabled {
		if config.ContactMethodEmail.CreateAndSendCustomEmail == nil {
			panic("Please pass a function (ContactMethodEmail.CreateAndSendCustomEmail) to send emails.")
		}
		typeNormalisedInput.ContactMethodEmail.Enabled = true
		if config.ContactMethodEmail.CreateAndSendCustomEmail != nil {
			typeNormalisedInput.ContactMethodEmail.CreateAndSendCustomEmail = config.ContactMethodEmail.CreateAndSendCustomEmail
		}
		if config.ContactMethodEmail.ValidateEmailAddress != nil {
			typeNormalisedInput.ContactMethodEmail.ValidateEmailAddress = config.ContactMethodEmail.ValidateEmailAddress
		}
	}

	if config.ContactMethodEmailOrPhone.Enabled {
		if config.ContactMethodEmailOrPhone.CreateAndSendCustomTextMessage == nil {
			panic("Please pass a function (ContactMethodEmailOrPhone.CreateAndSendCustomTextMessage) to send text messages.")
		}
		if config.ContactMethodEmailOrPhone.CreateAndSendCustomEmail == nil {
			panic("Please pass a function (ContactMethodEmailOrPhone.CreateAndSendCustomEmail) to send emails.")
		}
		typeNormalisedInput.ContactMethodEmailOrPhone.Enabled = true
		if config.ContactMethodEmailOrPhone.CreateAndSendCustomEmail != nil {
			typeNormalisedInput.ContactMethodEmailOrPhone.CreateAndSendCustomEmail = config.ContactMethodEmailOrPhone.CreateAndSendCustomEmail
		}
		if config.ContactMethodEmailOrPhone.ValidateEmailAddress != nil {
			typeNormalisedInput.ContactMethodEmailOrPhone.ValidateEmailAddress = config.ContactMethodEmailOrPhone.ValidateEmailAddress
		}
		if config.ContactMethodEmailOrPhone.CreateAndSendCustomTextMessage != nil {
			typeNormalisedInput.ContactMethodEmailOrPhone.CreateAndSendCustomTextMessage = config.ContactMethodEmailOrPhone.CreateAndSendCustomTextMessage
		}
		if config.ContactMethodEmailOrPhone.ValidatePhoneNumber != nil {
			typeNormalisedInput.ContactMethodEmailOrPhone.ValidatePhoneNumber = config.ContactMethodEmailOrPhone.ValidatePhoneNumber
		}
	}

	// FlowType is initialized correctly in makeTypeNormalisedInput

	if config.GetLinkDomainAndPath != nil {
		typeNormalisedInput.GetLinkDomainAndPath = config.GetLinkDomainAndPath
	}

	// GetCustomUserInputCode is initialized correctly in makeTypeNormalisedInput

	if config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
	}
	return typeNormalisedInput
}

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo, inputConfig plessmodels.TypeInput) plessmodels.TypeNormalisedInput {
	return plessmodels.TypeNormalisedInput{
		FlowType: inputConfig.FlowType,
		ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
			Enabled:                        false,
			ValidateEmailAddress:           defaultValidateEmailAddress,
			CreateAndSendCustomEmail:       inputConfig.ContactMethodEmailOrPhone.CreateAndSendCustomEmail,
			ValidatePhoneNumber:            defaultValidatePhoneNumber,
			CreateAndSendCustomTextMessage: inputConfig.ContactMethodEmailOrPhone.CreateAndSendCustomTextMessage,
		},
		ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
			Enabled:                        false,
			ValidatePhoneNumber:            defaultValidatePhoneNumber,
			CreateAndSendCustomTextMessage: inputConfig.ContactMethodPhone.CreateAndSendCustomTextMessage,
		},
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled:                  false,
			ValidateEmailAddress:     defaultValidateEmailAddress,
			CreateAndSendCustomEmail: inputConfig.ContactMethodEmail.CreateAndSendCustomEmail,
		},
		GetLinkDomainAndPath:   getDefaultGetLinkDomainAndPath(appInfo),
		GetCustomUserInputCode: inputConfig.GetCustomUserInputCode,
		Override: plessmodels.OverrideStruct{
			Functions: func(originalImplementation plessmodels.RecipeInterface) plessmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation plessmodels.APIInterface) plessmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}

func getDefaultGetLinkDomainAndPath(appInfo supertokens.NormalisedAppinfo) func(email *string, phoneNumber *string, userContext supertokens.UserContext) (string, error) {
	return func(email *string, phoneNumber *string, userContext supertokens.UserContext) (string, error) {
		return appInfo.WebsiteDomain.GetAsStringDangerous() + appInfo.WebsiteBasePath.GetAsStringDangerous() + "/verify", nil
	}
}

func defaultValidateEmailAddress(value interface{}) *string {
	if reflect.TypeOf(value).Kind() != reflect.String {
		msg := "Development bug: Please make sure the email field yields a string"
		return &msg
	}
	check, err := regexp.Match(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, []byte(value.(string)))
	if err != nil || !check {
		msg := "Email is invalid"
		return &msg
	}
	return nil
}

func defaultValidatePhoneNumber(value interface{}) *string {
	if reflect.TypeOf(value).Kind() != reflect.String {
		msg := "Development bug: Please make sure the email field yields a string"
		return &msg
	}

	parsedPhoneNumber, err := phonenumbers.Parse(value.(string), "")
	if err != nil {
		msg := "Phone number is invalid"
		return &msg
	}
	if !phonenumbers.IsValidNumber(parsedPhoneNumber) {
		msg := "Phone number is invalid"
		return &msg
	}
	return nil
}

// func defaultCreateAndSendCustomEmail(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) {
// 	// TODO:
// }

// func defaultCreateAndSendCustomTextMessage(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) {
// 	// TODO:
// }
