/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package epunittesting

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestEmailExistGetStopsWorkingWhenDisabled(t *testing.T) {
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
			emailpassword.Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						*originalImplementation.EmailExistsGET = nil
						return originalImplementation
					},
				},
			}),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {

		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	q := req.URL.Query()
	q.Add("email", "random@gmail.com")
	req.URL.RawQuery = q.Encode()
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestGoodInputsEmailExists(t *testing.T) {
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
			emailpassword.Init(nil),
			session.Init(nil),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

	passwordVal := "validPass123"

	emailVal := "random@email.com"

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": emailVal,
			},
			{
				"id":    "password",
				"value": passwordVal,
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	assert.Equal(t, 200, resp.StatusCode)

	assert.NoError(t, err)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)

	assert.Equal(t, "OK", response["status"])

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	q := req.URL.Query()
	q.Add("email", "random@email.com")
	req.URL.RawQuery = q.Encode()
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	data2, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(data2, &response2)

	assert.Equal(t, "OK", response2["status"])
	assert.Equal(t, true, response2["exists"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestGoodInputsEmailDoesNotExists(t *testing.T) {
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
			emailpassword.Init(nil),
			session.Init(nil),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	q := req.URL.Query()
	q.Add("email", "random@gmail.com")
	req.URL.RawQuery = q.Encode()
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)

	assert.Equal(t, "OK", response["status"])
	assert.Equal(t, false, response["exists"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestEmailExistsWithSyntacticallyInvalidEmail(t *testing.T) {
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
			emailpassword.Init(nil),
			session.Init(nil),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

	passwordVal := "validPass123"

	emailVal := "random@email.com"

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": emailVal,
			},
			{
				"id":    "password",
				"value": passwordVal,
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	assert.Equal(t, 200, resp.StatusCode)

	assert.NoError(t, err)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)

	assert.Equal(t, "OK", response["status"])

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	q := req.URL.Query()
	q.Add("email", "randomemail.com")
	req.URL.RawQuery = q.Encode()
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	data2, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(data2, &response2)

	assert.Equal(t, "OK", response2["status"])
	assert.Equal(t, false, response2["exists"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestEmailExistsWithUnNormalizedEmail(t *testing.T) {
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
			emailpassword.Init(nil),
			session.Init(nil),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

	passwordVal := "validPass123"

	emailVal := "random@email.com"

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": emailVal,
			},
			{
				"id":    "password",
				"value": passwordVal,
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	assert.Equal(t, 200, resp.StatusCode)

	assert.NoError(t, err)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)

	assert.Equal(t, "OK", response["status"])

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	q := req.URL.Query()
	q.Add("email", "RaNDom@email.com")
	req.URL.RawQuery = q.Encode()
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	data2, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(data2, &response2)

	assert.Equal(t, "OK", response2["status"])
	assert.Equal(t, true, response2["exists"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestEmailDoesExistsWithBadInput(t *testing.T) {
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
			emailpassword.Init(nil),
			session.Init(nil),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)

	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}
