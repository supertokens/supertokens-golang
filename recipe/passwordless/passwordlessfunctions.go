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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/passwordless/smsdelivery/supertokensService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var PasswordlessLoginEmailSentForTest bool = false
var PasswordlessLoginSmsSentForTest bool = false

func DefaultCreateAndSendCustomEmail(appInfo supertokens.NormalisedAppinfo) func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	return func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
		if supertokens.IsRunningInTestMode() {
			// if running in test mode, we do not want to send this.
			PasswordlessLoginEmailSentForTest = true
			return nil
		}
		url := "https://api.supertokens.io/0/st/auth/passwordless/login"
		data := map[string]interface{}{
			"email":           email,
			"appName":         appInfo.AppName,
			"codeLifetime":    codeLifetime,
			"urlWithLinkCode": urlWithLinkCode,
			"userInputCode":   userInputCode,
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
		if err != nil {
			return err
		}
		if resp.StatusCode >= 300 {
			return errors.New(fmt.Sprintf("Error sending email. API returned %d status.", resp.StatusCode))
		}
		return nil
	}
}

func DefaultCreateAndSendCustomTextMessage(appInfo supertokens.NormalisedAppinfo) func(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	return func(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
		if supertokens.IsRunningInTestMode() {
			// if running in test mode, we do not want to send this.
			PasswordlessLoginSmsSentForTest = true
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
		if err != nil {
			return err
		}
		if resp.StatusCode == 429 {
			smsData, err := json.Marshal(data["smsInput"])
			if err != nil {
				return err
			}
			fmt.Println("Free daily SMS quota reached. If using our managed service, please create a production environment to get dedicated API keys for SMS sending, or define your own method for sending SMS. For now, we are logging it below:")
			fmt.Printf("SMS content: %s\n", string(smsData))

			return nil
		}

		if resp.StatusCode >= 300 {
			return errors.New("error sending SMS")
		}
		return nil
	}
}
