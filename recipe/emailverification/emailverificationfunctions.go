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

package emailverification

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func DefaultGetEmailVerificationURL(appInfo supertokens.NormalisedAppinfo) func(evmodels.User) (string, error) {
	return func(user evmodels.User) (string, error) {
		return appInfo.WebsiteDomain.GetAsStringDangerous() + appInfo.WebsiteBasePath.GetAsStringDangerous() + "/verify-email", nil
	}
}

// TODO: add test to see query
func DefaultCreateAndSendCustomEmail(appInfo supertokens.NormalisedAppinfo) func(user evmodels.User, emailVerifyURLWithToken string) {
	return func(user evmodels.User, emailVerifyURLWithToken string) {
		if supertokens.IsRunningInTestMode() {
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
		_, err = client.Do(req)
		if err != nil {
			return
		}
		return
	}
}
