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

package supertokensService

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const SUPERTOKENS_SMS_SERVICE_URL = "https://api.supertokens.com/0/services/sms"

func MakeSupertokensService(config smsdelivery.SupertokensServiceConfig) smsdelivery.SmsDeliveryInterface {
	sendPasswordlessLoginSms := func(input smsdelivery.PasswordlessLoginType, userContext supertokens.UserContext) error {
		instance, err := supertokens.GetInstanceOrThrowError()
		if err != nil {
			return err
		}

		data := map[string]interface{}{
			"apiKey": config.ApiKey,
			"smsInput": map[string]interface{}{
				"type":         "PASSWORDLESS_LOGIN",
				"phoneNumber":  input.PhoneNumber,
				"codeLifetime": input.CodeLifetime,
				"appName":      instance.AppInfo.AppName,
			},
		}
		if input.UrlWithLinkCode != nil {
			data["smsInput"].(map[string]interface{})["urlWithLinkCode"] = *input.UrlWithLinkCode
		}
		if input.UserInputCode != nil {
			data["smsInput"].(map[string]interface{})["userInputCode"] = *input.UserInputCode
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", SUPERTOKENS_SMS_SERVICE_URL, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		req.Header.Set("content-type", "application/json")
		req.Header.Set("api-version", "0")
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil && resp.StatusCode < 300 {
			supertokens.LogDebugMessage(fmt.Sprintf("Passwordless login SMS sent to %s", input.PhoneNumber))
			return nil
		}

		if err != nil {
			supertokens.LogDebugMessage(fmt.Sprintf("Error: %s", err.Error()))
		} else {
			supertokens.LogDebugMessage(fmt.Sprintf("Error status: %d", resp.StatusCode))
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				supertokens.LogDebugMessage(fmt.Sprintf("Error: %s", err.Error()))
			} else {
				supertokens.LogDebugMessage(fmt.Sprintf("Error response: %s", string(body)))
			}

			err = errors.New(fmt.Sprintf("Error sending SMS. API returned %d status.", resp.StatusCode))
		}

		supertokens.LogDebugMessage("Logging the input below:")
		supertokens.LogDebugMessage(string(jsonData))
		return err
	}

	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin != nil {
			return sendPasswordlessLoginSms(*input.PasswordlessLogin, userContext)
		} else {
			return errors.New("should never come here")
		}
	}

	return smsdelivery.SmsDeliveryInterface{
		SendSms: &sendSms,
	}
}
