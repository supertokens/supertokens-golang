/* Copyright (c) 2026, VRAI Labs and/or its affiliates. All rights reserved.
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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestConsumeCodeVerifiesEmailOnlyWhenCodeWasSentToThatEmail(t *testing.T) {
	BeforeEach()
	connectionURI := unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			session.Init(nil),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
				},
			}),
		},
	}
	err := supertokens.Init(configValue)
	assert.NoError(t, err)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	consumeCode := func(code plessmodels.CreateCodeResponse) map[string]interface{} {
		assert.NotNil(t, code.OK)
		postBody, err := json.Marshal(map[string]interface{}{
			"preAuthSessionId": code.OK.PreAuthSessionID,
			"deviceId":         code.OK.DeviceID,
			"userInputCode":    code.OK.UserInputCode,
		})
		assert.NoError(t, err)
		resp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(postBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
		assert.Equal(t, "OK", result["status"])
		return result
	}

	phoneNumber := "+919876543210"
	email := "test@example.com"

	// Sign up with a phone number and store an (unverified) email for the user
	code, err := CreateCodeWithPhoneNumber("public", phoneNumber, nil)
	assert.NoError(t, err)
	result := consumeCode(code)
	userId := result["user"].(map[string]interface{})["id"].(string)

	updateResp, err := UpdateUser(userId, &email, nil)
	assert.NoError(t, err)
	assert.NotNil(t, updateResp.OK)

	isVerified, err := emailverification.IsEmailVerified(userId, &email)
	assert.NoError(t, err)
	assert.False(t, isVerified)

	// Signing in with a code sent by SMS does not verify the stored email
	code, err = CreateCodeWithPhoneNumber("public", phoneNumber, nil)
	assert.NoError(t, err)
	result = consumeCode(code)
	assert.Equal(t, userId, result["user"].(map[string]interface{})["id"])

	isVerified, err = emailverification.IsEmailVerified(userId, &email)
	assert.NoError(t, err)
	assert.False(t, isVerified)

	// Signing in with a code sent to the email verifies it
	code, err = CreateCodeWithEmail("public", email, nil)
	assert.NoError(t, err)
	result = consumeCode(code)
	assert.Equal(t, userId, result["user"].(map[string]interface{})["id"])

	isVerified, err = emailverification.IsEmailVerified(userId, &email)
	assert.NoError(t, err)
	assert.True(t, isVerified)
}
