/* Copyright (c) 2025, VRAI Labs and/or its affiliates. All rights reserved.
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

package webauthn

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func recoverAccountInput() emaildelivery.WebauthnRecoverAccountType {
	return emaildelivery.WebauthnRecoverAccountType{
		User: emaildelivery.User{
			ID:    "user-id-1",
			Email: "test@example.com",
		},
		RecoverAccountLink: "https://supertokens.io/auth/webauthn/recover?token=token-123&tenantId=public",
		TenantId:           "public",
	}
}

// Mirrors the supertokens-node BackwardCompatibilityService: the request must hit
// the supertokens recover endpoint with api-version 0 and a JSON body of
// { email, appName, recoverAccountURL }.
func TestMakeRecoverAccountRequestMatchesNodeService(t *testing.T) {
	appInfo := supertokens.NormalisedAppinfo{AppName: "TestApp"}
	input := recoverAccountInput()

	req, err := makeRecoverAccountRequest(appInfo, input)
	assert.NoError(t, err)

	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, SUPERTOKENS_WEBAUTHN_RECOVER_ACCOUNT_URL, req.URL.String())
	assert.Equal(t, "https://api.supertokens.com/0/st/auth/webauthn/recover", req.URL.String())
	assert.Equal(t, "application/json; charset=utf-8", req.Header.Get("content-type"))
	assert.Equal(t, "0", req.Header.Get("api-version"))

	bodyBytes, err := io.ReadAll(req.Body)
	assert.NoError(t, err)

	var body map[string]any
	assert.NoError(t, json.Unmarshal(bodyBytes, &body))
	assert.Equal(t, "test@example.com", body["email"])
	assert.Equal(t, "TestApp", body["appName"])
	assert.Equal(t, input.RecoverAccountLink, body["recoverAccountURL"])
	assert.Len(t, body, 3)
}

// In test mode the service must short-circuit before issuing any network call,
// mirroring node's isTestEnv() guard.
func TestDefaultEmailServiceSkipsSendInTestMode(t *testing.T) {
	emailService := makeDefaultEmailService(supertokens.NormalisedAppinfo{AppName: "TestApp"})
	assert.NotNil(t, emailService.SendEmail)

	input := recoverAccountInput()
	err := (*emailService.SendEmail)(emaildelivery.EmailType{
		WebauthnRecoverAccount: &input,
	}, nil)
	assert.NoError(t, err)
}

// The default service only knows how to handle webauthn recover-account emails.
func TestDefaultEmailServiceErrorsForUnsupportedInput(t *testing.T) {
	emailService := makeDefaultEmailService(supertokens.NormalisedAppinfo{AppName: "TestApp"})

	err := (*emailService.SendEmail)(emaildelivery.EmailType{}, nil)
	assert.Error(t, err)
}
