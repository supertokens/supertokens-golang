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

package thirdpartyemailpassword

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestTheGenerateTokenAPIwithValidInputEmailNotVerified(t *testing.T) {
	customCSRFval := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customCSRFval,
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("random@gmail.com", "validPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cookieData := unittesting.ExtractInfoFromResponse(resp)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", result["status"])
	user := result["user"].(map[string]interface{})

	rep1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, user["id"].(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])
	if err != nil {
		t.Error(err.Error())
	}
	result1 := *unittesting.HttpResponseToConsumableInformation(rep1.Body)
	assert.Equal(t, "OK", result1["status"])
}

func TestGenerateTokenAPIwithValidInputEmailVerifiedAndTestError(t *testing.T) {
	customCSRFval := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customCSRFval,
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("random@gmail.com", "validPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cookieData := unittesting.ExtractInfoFromResponse(resp)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", result["status"])
	user := result["user"].(map[string]interface{})

	verifyToken, err := CreateEmailVerificationToken(user["id"].(string))
	if err != nil {
		t.Error(err.Error())
	}
	VerifyEmailUsingToken(verifyToken.OK.Token)

	rep1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, user["id"].(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])
	if err != nil {
		t.Error(err.Error())
	}
	result1 := *unittesting.HttpResponseToConsumableInformation(rep1.Body)
	assert.Equal(t, "EMAIL_ALREADY_VERIFIED_ERROR", result1["status"])
}

func TestGenerateTokenAPIWithValidInputNoSessionAndCheckOutput(t *testing.T) {
	customCSRFval := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customCSRFval,
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := http.Post(testServer.URL+"/auth/user/email/verify/token", "application/json", nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "unauthorised", result["message"])
}

func TestEmailVerifyAPIwithInvalidTokenCheckError(t *testing.T) {
	customCSRFval := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customCSRFval,
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	formFields := map[string]string{
		"method": "token",
		"token":  "randomToken",
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/user/email/verify", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "EMAIL_VERIFICATION_INVALID_TOKEN_ERROR", result["status"])
}
