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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func initWebauthnRecipeForAPITests(t *testing.T) {
	connectionURI := unittesting.StartUpST("localhost", "8080")
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
			Init(nil),
		},
	}

	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestRegisterOptionsRequiresEmailOrRecoverToken(t *testing.T) {
	BeforeEach()
	defer AfterEach()

	initWebauthnRecipeForAPITests(t)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/webauthn/options/register", bytes.NewBuffer([]byte(`{}`)))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "webauthn")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(res.Body)
	assert.Equal(t, "Please provide either email or recoverAccountToken", result["message"])
}

func TestRegisterOptionsWithInvalidEmailReturnsInvalidEmailError(t *testing.T) {
	BeforeEach()
	defer AfterEach()

	initWebauthnRecipeForAPITests(t)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	requestBody, err := json.Marshal(map[string]interface{}{
		"email": "invalid-email",
	})
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/webauthn/options/register", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "webauthn")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(res.Body)
	assert.Equal(t, "INVALID_EMAIL_ERROR", result["status"])
	assert.Equal(t, "Email is not valid", result["err"])
}

func TestSignUpWithInvalidCredentialShapeReturnsInvalidCredentialsError(t *testing.T) {
	BeforeEach()
	defer AfterEach()

	initWebauthnRecipeForAPITests(t)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	requestBody, err := json.Marshal(map[string]interface{}{
		"webauthnGeneratedOptionsId": "options-id-1",
		"credential":                 "invalid-credential-shape",
	})
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/webauthn/signup", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "webauthn")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(res.Body)
	assert.Equal(t, "INVALID_CREDENTIALS_ERROR", result["status"])
}

func TestSignInRequiresWebauthnGeneratedOptionsId(t *testing.T) {
	BeforeEach()
	defer AfterEach()

	initWebauthnRecipeForAPITests(t)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	requestBody, err := json.Marshal(map[string]interface{}{
		"credential": map[string]interface{}{},
	})
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/webauthn/signin", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "webauthn")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(res.Body)
	assert.Equal(t, "webauthnGeneratedOptionsId is required", result["message"])
}

func TestEmailExistsRequiresEmailParam(t *testing.T) {
	BeforeEach()
	defer AfterEach()

	initWebauthnRecipeForAPITests(t)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/webauthn/email/exists", nil)
	assert.NoError(t, err)
	req.Header.Add("rid", "webauthn")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(res.Body)
	assert.Equal(t, "Please provide the email as a GET param", result["message"])
}
