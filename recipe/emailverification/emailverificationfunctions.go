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

package emailverification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// used for testing purposes.
var EmailVerificationEmailSentForTest = false
var EmailVerificationDataForTest = struct {
	User                    evmodels.User
	EmailVerifyURLWithToken string
	UserContext             supertokens.UserContext
}{}

func DefaultCreateAndSendCustomEmail(appInfo supertokens.NormalisedAppinfo) func(user evmodels.User, emailVerifyURLWithToken string, userContext supertokens.UserContext) {
	return func(user evmodels.User, emailVerifyURLWithToken string, userContext supertokens.UserContext) {
		if supertokens.IsRunningInTestMode() {
			EmailVerificationEmailSentForTest = true
			EmailVerificationDataForTest.User = user
			EmailVerificationDataForTest.EmailVerifyURLWithToken = emailVerifyURLWithToken
			EmailVerificationDataForTest.UserContext = userContext
			// if running in test mode, we do not want to send this.
			return
		}
		const url = "https://api.supertokens.io/0/st/auth/email/verify"

		data := map[string]string{
			"email":          user.Email,
			"appName":        appInfo.AppName,
			"emailVerifyURL": emailVerifyURLWithToken,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return
		}

		req.Header.Set("content-type", "application/json")
		req.Header.Set("api-version", "0")
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil && resp.StatusCode < 300 {
			supertokens.LogDebugMessage(fmt.Sprintf("Email verification email sent to %s", user.Email))
			return
		}

		supertokens.LogDebugMessage("Error sending verification email")
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
		}
		supertokens.LogDebugMessage("Logging the input below:")
		supertokens.LogDebugMessage(string(jsonData))
	}
}
