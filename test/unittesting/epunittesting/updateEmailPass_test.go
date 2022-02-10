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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestUpdateEmailPass(t *testing.T) {
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

	_, err = unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	res, err := unittesting.SignInRequest("testrandom@gmail.com", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}
	dataInBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	email := "test2@gmail.com"
	password := "testPass"

	emailpassword.UpdateEmailOrPassword(data["user"].(map[string]interface{})["id"].(string), &email, &password)

	res1, err := unittesting.SignInRequest("testrandom@gmail.com", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}
	dataInBytes1, err := io.ReadAll(res1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res1.Body.Close()

	var data1 map[string]interface{}
	err = json.Unmarshal(dataInBytes1, &data1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "WRONG_CREDENTIALS_ERROR", data1["status"])

	res2, err := unittesting.SignInRequest(email, password, testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}
	dataInBytes2, err := io.ReadAll(res2.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res2.Body.Close()

	var data2 map[string]interface{}
	err = json.Unmarshal(dataInBytes2, &data2)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", data2["status"])
	assert.Equal(t, email, data2["user"].(map[string]interface{})["email"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}
