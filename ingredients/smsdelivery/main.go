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

package smsdelivery

import (
	"errors"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Ingredient struct {
	IngredientInterfaceImpl SmsDeliveryInterface
}

func MakeIngredient(config TypeInputWithService) Ingredient {
	result := Ingredient{
		IngredientInterfaceImpl: config.Service,
	}

	if config.Override != nil {
		result.IngredientInterfaceImpl = config.Override(result.IngredientInterfaceImpl)
	}

	return result
}

func SendTwilioSms(settings TwilioSettings, content SMSContent) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: settings.AccountSid,
		Password: settings.AuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(content.ToPhoneNumber)
	params.SetBody(content.Body)

	if settings.From != "" {
		params.SetFrom(settings.From)
	} else if settings.MessagingServiceSid != "" {
		params.SetMessagingServiceSid(settings.MessagingServiceSid)
	} else {
		return errors.New("should not come here")
	}

	_, err := client.Api.CreateMessage(params)

	return err
}
