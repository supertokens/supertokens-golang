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

	"github.com/supertokens/supertokens-golang/supertokens"
)

type TwilioSettings struct {
	AccountSid          string
	AuthToken           string
	From                *string
	MessagingServiceSid *string
}

type TwilioContent struct {
	Body          string
	ToPhoneNumber string
}

type TwilioServiceInterface struct {
	SendRawSms *func(input TwilioContent, userContext supertokens.UserContext) error
	GetContent *func(input SmsType, userContext supertokens.UserContext) (TwilioContent, error)
}

type TwilioServiceConfig struct {
	Settings TwilioSettings
	Override func(originalImplementation TwilioServiceInterface) TwilioServiceInterface
}

func NormaliseTwilioServiceConfig(input TwilioServiceConfig) (TwilioServiceConfig, error) {
	if input.Settings.From == nil && input.Settings.MessagingServiceSid == nil {
		return TwilioServiceConfig{}, errors.New("either 'From' or 'MessagingServiceSid' must be set")
	}
	if input.Settings.From != nil && input.Settings.MessagingServiceSid != nil {
		return TwilioServiceConfig{}, errors.New("only one of 'From' or 'MessagingServiceSid' must be set")
	}
	return input, nil
}
