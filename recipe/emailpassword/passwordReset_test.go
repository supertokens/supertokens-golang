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

package emailpassword

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestEmailValidationCheckInGenerateTokenAPI(t *testing.T) {
	resetURL := ""
	tokenInfo := ""
	ridInfo := ""
	sendEmailFunc := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		u, err := url.Parse(input.PasswordReset.PasswordResetLink)
		if err != nil {
			return err
		}
		resetURL = u.Scheme + "://" + u.Host + u.Path
		tokenInfo = u.Query().Get("token")
		ridInfo = u.Query().Get("rid")
		return nil
	}
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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmailFunc,
					},
				},
			}),
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

	res, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	dataInBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	var result map[string]interface{}
	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "OK", result["status"])

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": "random@gmail.com",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/user/password/reset/token", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "https://supertokens.io/auth/reset-password", resetURL)
	assert.NotEmpty(t, tokenInfo)
	assert.True(t, strings.HasPrefix(ridInfo, "emailpassword"))
}

func TestPasswordValidation(t *testing.T) {
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

	formFields := map[string][]map[string]interface{}{
		"formFields": {
			{
				"id":    "password",
				"value": "invalid",
			},
		},
	}
	token := map[string]interface{}{
		"token": "RandomToken",
	}
	var data map[string]interface{}

	a, _ := json.Marshal(formFields)
	json.Unmarshal(a, &data)
	b, _ := json.Marshal(token)
	json.Unmarshal(b, &data)

	jData, _ := json.Marshal(data)

	resp, err := http.Post(testServer.URL+"/auth/user/password/reset", "application/json", bytes.NewBuffer(jData))

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)

	dataInBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var data1 map[string]interface{}
	json.Unmarshal(dataInBytes, &data1)

	assert.Equal(t, "FIELD_ERROR", data1["status"])
	assert.Equal(t, "Password must contain at least 8 characters, including a number", data1["formFields"].([]interface{})[0].(map[string]interface{})["error"])
	assert.Equal(t, "password", data1["formFields"].([]interface{})[0].(map[string]interface{})["id"])

	formFields1 := map[string][]map[string]interface{}{
		"formFields": {
			{
				"id":    "password",
				"value": "validpass123",
			},
		},
	}
	token1 := map[string]interface{}{
		"token": "RandomToken",
	}
	var data2 map[string]interface{}

	a1, _ := json.Marshal(formFields1)
	json.Unmarshal(a1, &data2)
	b1, _ := json.Marshal(token1)
	json.Unmarshal(b1, &data2)

	jData1, _ := json.Marshal(data2)

	resp1, err := http.Post(testServer.URL+"/auth/user/password/reset", "application/json", bytes.NewBuffer(jData1))

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)

	dataInBytes1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()

	var data3 map[string]interface{}
	json.Unmarshal(dataInBytes1, &data3)

	assert.Equal(t, "RESET_PASSWORD_INVALID_TOKEN_ERROR", data3["status"])
}

func TestTokenMissingFromInput(t *testing.T) {
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

	formFields := map[string][]map[string]interface{}{
		"formFields": {
			{
				"id":    "password",
				"value": "validpass123",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/user/password/reset", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	dataInBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var data1 map[string]interface{}
	json.Unmarshal(dataInBytes, &data1)

	assert.Equal(t, "Please provide the password reset token", data1["message"])

}

func TestValidTokenInputAndPasswordHasChanged(t *testing.T) {
	var token string
	sendEmailFunc := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		u, err := url.Parse(input.PasswordReset.PasswordResetLink)
		if err != nil {
			return err
		}
		token = u.Query().Get("token")
		return nil
	}
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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmailFunc,
					},
				},
			}),
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

	res, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, res.StatusCode)
	dataInBytes, _ := io.ReadAll(res.Body)
	res.Body.Close()

	var userData map[string]interface{}
	err = json.Unmarshal(dataInBytes, &userData)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "OK", userData["status"])
	userInfo := userData["user"].(map[string]interface{})

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": "random@gmail.com",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = http.Post(testServer.URL+"/auth/user/password/reset/token", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)

	formFields1 := map[string][]map[string]interface{}{
		"formFields": {
			{
				"id":    "password",
				"value": "validpass12345",
			},
		},
	}
	token1 := map[string]interface{}{
		"token": token,
	}
	var data2 map[string]interface{}

	a1, _ := json.Marshal(formFields1)
	json.Unmarshal(a1, &data2)
	b1, _ := json.Marshal(token1)
	json.Unmarshal(b1, &data2)

	jData1, _ := json.Marshal(data2)

	_, err = http.Post(testServer.URL+"/auth/user/password/reset", "application/json", bytes.NewBuffer(jData1))

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)

	res2, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)

	dataInBytes2, err := io.ReadAll(res2.Body)

	if err != nil {
		t.Error(err.Error())
	}

	res2.Body.Close()

	var result2 map[string]interface{}

	err = json.Unmarshal(dataInBytes2, &result2)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "WRONG_CREDENTIALS_ERROR", result2["status"])

	res3, err := unittesting.SignInRequest("random@gmail.com", "validpass12345", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)

	dataInBytes3, err := io.ReadAll(res3.Body)

	if err != nil {
		t.Error(err.Error())
	}

	res3.Body.Close()

	var result3 map[string]interface{}

	err = json.Unmarshal(dataInBytes3, &result3)
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, result3["user"])
	assert.Equal(t, userInfo["id"], result3["user"].(map[string]interface{})["id"].(string))
	assert.Equal(t, userInfo["email"], result3["user"].(map[string]interface{})["email"].(string))
}
