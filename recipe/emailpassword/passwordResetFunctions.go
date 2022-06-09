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

package emailpassword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var PasswordResetEmailSentForTest = false
var PasswordResetDataForTest = struct {
	User                      epmodels.User
	PasswordResetURLWithToken string
	UserContext               supertokens.UserContext
}{}

func defaultGetResetPasswordURL(appInfo supertokens.NormalisedAppinfo) func(_ epmodels.User, userContext supertokens.UserContext) (string, error) {
	return func(_ epmodels.User, userContext supertokens.UserContext) (string, error) {
		return appInfo.WebsiteDomain.GetAsStringDangerous() + appInfo.WebsiteBasePath.GetAsStringDangerous() + "/reset-password", nil
	}
}

// TODO: add test to see query
func DefaultCreateAndSendCustomPasswordResetEmail(appInfo supertokens.NormalisedAppinfo) func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext) {
	return func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext) {
		if supertokens.IsRunningInTestMode() {
			// if running in test mode, we do not want to send this.
			PasswordResetEmailSentForTest = true
			PasswordResetDataForTest.User = user
			PasswordResetDataForTest.PasswordResetURLWithToken = passwordResetURLWithToken
			PasswordResetDataForTest.UserContext = userContext
			return
		}
		url := "https://api.supertokens.io/0/st/auth/password/reset"
		data := map[string]string{
			"email":            user.Email,
			"appName":          appInfo.AppName,
			"passwordResetURL": passwordResetURLWithToken,
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
			supertokens.LogDebugMessage(fmt.Sprintf("Password reset email sent to %s", user.Email))
			return
		}

		supertokens.LogDebugMessage("Error sending password reset email")
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
