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

func SendTwilioSms(config TwilioServiceConfig, content TwilioGetContentResult) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSid,
		Password: config.AuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(content.ToPhoneNumber)
	params.SetBody(content.Body)

	if config.From != nil {
		params.SetFrom(*config.From)
	} else if config.MessagingServiceSid != nil {
		params.SetMessagingServiceSid(*config.MessagingServiceSid)
	} else {
		return errors.New("should not come here")
	}

	_, err := client.Api.CreateMessage(params)

	return err
}
