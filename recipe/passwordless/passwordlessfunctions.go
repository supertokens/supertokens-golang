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

package passwordless

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/passwordless/smsdelivery/supertokensService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var PasswordlessLoginEmailSentForTest bool = false
var PasswordlessLoginEmailDataForTest struct {
	Email            string
	UserInputCode    *string
	UrlWithLinkCode  *string
	CodeLifetime     uint64
	PreAuthSessionId string
	UserContext      supertokens.UserContext
}
var PasswordlessLoginSmsSentForTest bool = false
var PasswordlessLoginSmsDataForTest struct {
	Phone            string
	UserInputCode    *string
	UrlWithLinkCode  *string
	CodeLifetime     uint64
	PreAuthSessionId string
	UserContext      supertokens.UserContext
}

func logAndReturnError(resp *http.Response, err error) error {
	if err != nil {
		supertokens.LogDebugMessage(fmt.Sprintf("Error: %s", err.Error()))
		return err
	}

	supertokens.LogDebugMessage(fmt.Sprintf("Error status: %d", resp.StatusCode))
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		supertokens.LogDebugMessage(fmt.Sprintf("Error: %s", err.Error()))
		return err
	}

	supertokens.LogDebugMessage(fmt.Sprintf("Error response: %s", string(body)))

	var bodyObj map[string]interface{}
	if err = json.Unmarshal(body, &bodyObj); err == nil {
		if errMsg, ok := bodyObj["err"].(string); ok {
			return errors.New(errMsg)
		}
	}

	return errors.New(string(body))
}

func DefaultCreateAndSendCustomEmail(appInfo supertokens.NormalisedAppinfo) func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	return func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
		if supertokens.IsRunningInTestMode() {
			// if running in test mode, we do not want to send this.
			PasswordlessLoginEmailSentForTest = true
			PasswordlessLoginEmailDataForTest.Email = email
			PasswordlessLoginEmailDataForTest.UserInputCode = userInputCode
			PasswordlessLoginEmailDataForTest.UrlWithLinkCode = urlWithLinkCode
			PasswordlessLoginEmailDataForTest.CodeLifetime = codeLifetime
			PasswordlessLoginEmailDataForTest.PreAuthSessionId = preAuthSessionId
			PasswordlessLoginEmailDataForTest.UserContext = userContext
			return nil
		}
		url := "https://api.supertokens.io/0/st/auth/passwordless/login"
		data := map[string]interface{}{
			"email":        email,
			"appName":      appInfo.AppName,
			"codeLifetime": codeLifetime,
		}
		if urlWithLinkCode != nil {
			data["urlWithLinkCode"] = *urlWithLinkCode
		}
		if userInputCode != nil {
			data["userInputCode"] = *userInputCode
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("api-version", "0")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil && resp.StatusCode < 300 {
			supertokens.LogDebugMessage(fmt.Sprintf("Passwordless login email sent to %s", email))
			return nil
		}

		err = logAndReturnError(resp, err)
		supertokens.LogDebugMessage("Logging the input below:")
		supertokens.LogDebugMessage(string(jsonData))
		return err
	}
}

func DefaultCreateAndSendCustomTextMessage(appInfo supertokens.NormalisedAppinfo) func(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	return func(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
		if supertokens.IsRunningInTestMode() {
			// if running in test mode, we do not want to send this.
			PasswordlessLoginSmsSentForTest = true
			PasswordlessLoginSmsDataForTest.Phone = phoneNumber
			PasswordlessLoginSmsDataForTest.UserInputCode = userInputCode
			PasswordlessLoginSmsDataForTest.UrlWithLinkCode = urlWithLinkCode
			PasswordlessLoginSmsDataForTest.CodeLifetime = codeLifetime
			PasswordlessLoginSmsDataForTest.PreAuthSessionId = preAuthSessionId
			PasswordlessLoginSmsDataForTest.UserContext = userContext

			return nil
		}

		data := map[string]map[string]interface{}{
			"smsInput": {
				"type":         "PASSWORDLESS_LOGIN",
				"phoneNumber":  phoneNumber,
				"codeLifetime": codeLifetime,
				"appName":      appInfo.AppName,
			},
		}
		if urlWithLinkCode != nil {
			data["smsInput"]["urlWithLinkCode"] = *urlWithLinkCode
		}
		if userInputCode != nil {
			data["smsInput"]["userInputCode"] = *userInputCode
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", supertokensService.SUPERTOKENS_SMS_SERVICE_URL, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		req.Header.Set("content-type", "application/json")
		req.Header.Set("api-version", "0")
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil && resp.StatusCode < 300 {
			supertokens.LogDebugMessage(fmt.Sprintf("Passwordless login SMS sent to %s", phoneNumber))
			return nil
		}

		if err == nil && resp.StatusCode == 429 {
			smsData, err := json.Marshal(data["smsInput"])
			if err != nil {
				return err
			}
			fmt.Println("Free daily SMS quota reached. If you want to use SuperTokens to send SMS, please sign up on supertokens.com to get your SMS API key, else you can also define your own method by overriding the service. For now, we are logging it below:")
			fmt.Println()
			fmt.Printf("SMS content: %s\n", string(smsData))

			return nil
		}

		err = logAndReturnError(resp, err)
		supertokens.LogDebugMessage("Logging the input below:")
		supertokens.LogDebugMessage(string(jsonData))
		return err
	}
}
