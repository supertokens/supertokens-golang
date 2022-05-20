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
	"net/http"

	"github.com/supertokens/supertokens-golang/supertokens"
)

var PasswordlessLoginEmailSentForTest bool = false

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
		_, err = client.Do(req)
		if err != nil {
			return err
		}
		return nil
	}
}
