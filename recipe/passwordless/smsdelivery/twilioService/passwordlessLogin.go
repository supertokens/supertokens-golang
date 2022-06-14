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

package twilioService

import (
	"strings"

	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const magicLinkLoginTemplate = `Click ${magicLink} to login to ${appname}

This is valid for ${time}.`
const otpLoginTemplate = `OTP to login is ${otp} for ${appname}

This is valid for ${time}.`
const magicLinkAndOtpLoginTemplate = `OTP to login is ${otp} for ${appname}

Or click ${magicLink} to login.

This is valid for ${time}.`

func getPasswordlessLoginSmsContent(input smsdelivery.PasswordlessLoginType) smsdelivery.TwilioGetContentResult {
	stInstance, err := supertokens.GetInstanceOrThrowError()
	if err != nil {
		panic("Please call supertokens.Init function before using the Middleware")
	}
	return smsdelivery.TwilioGetContentResult{
		Body:          getPasswordlessLoginSmsBody(stInstance.AppInfo.AppName, input.CodeLifetime, input.UrlWithLinkCode, input.UserInputCode),
		ToPhoneNumber: input.PhoneNumber,
	}
}

func getPasswordlessLoginSmsBody(appName string, codeLifetime uint64, urlWithLinkCode *string, userInputCode *string) string {
	var smsBody string

	if urlWithLinkCode != nil && userInputCode != nil {
		smsBody = magicLinkAndOtpLoginTemplate
	} else if urlWithLinkCode != nil {
		smsBody = magicLinkLoginTemplate
	} else if userInputCode != nil {
		smsBody = otpLoginTemplate
	} else {
		// Should never come here
	}

	humanisedCodeLifetime := supertokens.HumaniseMilliseconds(codeLifetime)

	smsBody = strings.Replace(smsBody, "*|MC:SUBJECT|*", "Login to your account", -1)
	smsBody = strings.Replace(smsBody, "${appname}", appName, -1)
	smsBody = strings.Replace(smsBody, "${time}", humanisedCodeLifetime, -1)
	if urlWithLinkCode != nil {
		smsBody = strings.Replace(smsBody, "${magicLink}", *urlWithLinkCode, -1)
	}
	if userInputCode != nil {
		smsBody = strings.Replace(smsBody, "${otp}", *userInputCode, -1)
	}

	return smsBody
}
