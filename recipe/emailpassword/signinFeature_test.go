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

package emailpassword

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestAfterDisablingTheDefaultSigninAPIdoesNotWork(t *testing.T) {
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
			Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						*originalImplementation.SignInPOST = nil
						return originalImplementation
					},
				},
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

	resp, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestSignInAPIWorksWhenInputIsFine(t *testing.T) {
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
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

	resp, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]interface{}
	json.Unmarshal(dataInBytes, &result)
	resp.Body.Close()

	assert.Equal(t, "OK", result["status"])

	signupUserInfo := result["user"].(map[string]interface{})

	resp1, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	dataInBytes1, err := io.ReadAll(resp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result1 map[string]interface{}
	json.Unmarshal(dataInBytes1, &result1)
	resp1.Body.Close()

	assert.Equal(t, "OK", result1["status"])

	signInUserInfo := result1["user"].(map[string]interface{})

	assert.Equal(t, signInUserInfo["id"], signupUserInfo["id"])
	assert.Equal(t, signInUserInfo["email"], signupUserInfo["email"])
}

func TestSigninAPIthrowsAnErrorWhenEmailDoesNotExist(t *testing.T) {
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
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

	resp, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]interface{}
	json.Unmarshal(dataInBytes, &result)
	resp.Body.Close()

	assert.Equal(t, "OK", result["status"])

	resp1, err := unittesting.SignInRequest("rand@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	dataInBytes1, err := io.ReadAll(resp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result1 map[string]interface{}
	json.Unmarshal(dataInBytes1, &result1)
	resp1.Body.Close()

	assert.Equal(t, "WRONG_CREDENTIALS_ERROR", result1["status"])
}
