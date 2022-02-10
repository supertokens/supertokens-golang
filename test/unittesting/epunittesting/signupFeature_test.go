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
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDisablingAPIDefaultSignUpDoesNotWork(t *testing.T) {
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
						*originalImplementation.SignUpPOST = nil
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

	res, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestSignUpAPIworksWithValidInput(t *testing.T) {
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

	res, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
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

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "OK", data["status"])

	signupUserInfo := data["user"].(map[string]interface{})

	assert.NotNil(t, signupUserInfo["id"])
	assert.Equal(t, "random@gmail.com", signupUserInfo["email"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestSignUpAPIThrowsErrorInCaseOfDuplicateEmail(t *testing.T) {
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

	res, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
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

	signupUserInfo := data["user"].(map[string]interface{})

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "OK", data["status"])

	assert.NotNil(t, signupUserInfo["id"])
	assert.Equal(t, "random@gmail.com", signupUserInfo["email"])

	res1, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
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

	assert.Equal(t, "FIELD_ERROR", data1["status"])
	assert.Equal(t, 1, len(data1["formFields"].([]interface{})))
	assert.Equal(t, "This email already exists. Please sign in instead.", data1["formFields"].([]interface{})[0].(map[string]interface{})["error"])
	assert.Equal(t, "email", data1["formFields"].([]interface{})[0].(map[string]interface{})["id"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestSignUpAPIThrowsErrorForInvalidEmailAndPassword(t *testing.T) {
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

	res, err := unittesting.SignupRequest("randomgmail.com", "invalidpass", testServer.URL)
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

	assert.Equal(t, "FIELD_ERROR", data["status"])
	assert.Equal(t, 2, len(data["formFields"].([]interface{})))

	formFields := data["formFields"].([]interface{})

	for _, formField := range formFields {
		if formField.(map[string]interface{})["id"] == "email" {
			assert.Equal(t, "Email is invalid", formField.(map[string]interface{})["error"])
		} else {
			assert.Equal(t, "Password must contain at least one number", formField.(map[string]interface{})["error"])
		}
	}

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestBadInputNoPostBodyToSignUpAPI(t *testing.T) {
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

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", nil)

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes1, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, "unexpected end of JSON input\n", string(dataInBytes1))

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestBadInputFormFieldsElementsHaveNoId(t *testing.T) {
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

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"randomKey": "randomValue",
			},
			{
				"randomKey2": "randomValue2",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes1, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes1, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "FIELD_ERROR", data["status"])
	assert.Equal(t, 2, len(data["formFields"].([]interface{})))

	formFields1 := data["formFields"].([]interface{})

	for _, formField := range formFields1 {
		if formField.(map[string]interface{})["id"] == "email" {
			assert.Equal(t, "Field is not optional", formField.(map[string]interface{})["error"])
		} else {
			assert.Equal(t, "Field is not optional", formField.(map[string]interface{})["error"])
		}
	}

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestSuccessfullSigUpYieldsSession(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	res, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
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

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "OK", data["status"])

	cookieData := unittesting.ExtractInfoFromResponse(res)

	assert.NotNil(t, cookieData["antiCsrf"])

	assert.NotNil(t, cookieData["sAccessToken"])
	assert.NotNil(t, cookieData["sRefreshToken"])
	assert.NotNil(t, cookieData["sIdRefreshToken"])

	assert.NotNil(t, cookieData["refreshTokenExpiry"])
	assert.NotNil(t, cookieData["refreshTokenDomain"])
	assert.NotNil(t, cookieData["refreshTokenHttpOnly"])

	assert.NotNil(t, cookieData["idRefreshTokenExpiry"])
	assert.NotNil(t, cookieData["idRefreshTokenDomain"])
	assert.NotNil(t, cookieData["idRefreshTokenHttpOnly"])

	assert.NotNil(t, cookieData["accessTokenExpiry"])
	assert.NotNil(t, cookieData["accessTokenDomain"])
	assert.NotNil(t, cookieData["accessTokenHttpOnly"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestExtraFieldAddingInSignupFormFieldWorks(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
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
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "testField",
						},
					},
				},
			}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "password",
				"value": "validpass123",
			},
			{
				"id":    "email",
				"value": "random@gmail.com",
			},
			{
				"id":    "testField",
				"value": "testValue",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", data["status"])
	assert.NotNil(t, data["user"].(map[string]interface{})["id"])
	assert.Equal(t, "random@gmail.com", data["user"].(map[string]interface{})["email"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestThatCustomFieldsAreSentUsingHandlePostSignup(t *testing.T) {
	var customFormFields []epmodels.TypeFormField
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
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "testField",
						},
					},
				},
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						originalSignUpPost := *originalImplementation.SignUpPOST
						*originalImplementation.SignUpPOST = func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.SignUpResponse, error) {
							res, _ := originalSignUpPost(formFields, options)
							customFormFields = formFields
							return res, nil
						}
						return originalImplementation
					},
				},
			}),
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

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "password",
				"value": "validpass123",
			},
			{
				"id":    "email",
				"value": "random@gmail.com",
			},
			{
				"id":    "testField",
				"value": "testValue",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", data["status"])

	assert.Equal(t, "password", customFormFields[0].ID)
	assert.Equal(t, "email", customFormFields[1].ID)
	assert.Equal(t, "testField", customFormFields[2].ID)

	assert.Equal(t, "validpass123", customFormFields[0].Value)
	assert.Equal(t, "random@gmail.com", customFormFields[1].Value)
	assert.Equal(t, "testValue", customFormFields[2].Value)

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestFormFieldsAddedInConfigButNotInInputToSignupCheckErrorAboutItBeingMissing(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
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
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "testField",
						},
					},
				},
			}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	res, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	assert.Equal(t, 500, res.StatusCode)
	assert.Equal(t, "Are you sending too many / too few formFields?\n", string(dataInBytes))

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestBadCaseInputWithoutOtional(t *testing.T) {
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
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "testField",
						},
					},
				},
			}),
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

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "password",
				"value": "validpass123",
			},
			{
				"id":    "email",
				"value": "random@gmail.com",
			},
			{
				"id":    "testField",
				"value": "",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "FIELD_ERROR", data["status"])
	assert.Equal(t, 1, len(data["formFields"].([]interface{})))
	assert.Equal(t, "testField", data["formFields"].([]interface{})[0].(map[string]interface{})["id"])
	assert.Equal(t, "Field is not optional", data["formFields"].([]interface{})[0].(map[string]interface{})["error"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestGoodCaseInputWithOtional(t *testing.T) {
	optionalVal := true
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
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID:       "testField",
							Optional: &optionalVal,
						},
					},
				},
			}),
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

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "password",
				"value": "validpass123",
			},
			{
				"id":    "email",
				"value": "random@gmail.com",
			},
			{
				"id":    "testField",
				"value": "",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", data["status"])
	assert.NotNil(t, data["user"].(map[string]interface{})["id"])
	assert.Equal(t, "random@gmail.com", data["user"].(map[string]interface{})["email"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestInputFormFieldWithoutEmailField(t *testing.T) {
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

	formFields := map[string][]map[string]string{
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

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 500, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, "Are you sending too many / too few formFields?\n", string(dataInBytes))

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestInputFormFieldWithoutPasswordField(t *testing.T) {
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

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 500, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, "Are you sending too many / too few formFields?\n", string(dataInBytes))

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestInputFormFieldHasADifferentNumberOfCustomFiledsThanInConfigFormFields(t *testing.T) {
	optionalVal := true
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
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID:       "testField",
							Optional: &optionalVal,
						},
						{
							ID: "testField2",
						},
					},
				},
			}),
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

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "password",
				"value": "validpass123",
			},
			{
				"id":    "email",
				"value": "random@gmail.com",
			},
			{
				"id":    "testField",
				"value": "",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 500, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, "Are you sending too many / too few formFields?\n", string(dataInBytes))

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestInputFormFieldHasSameNumberOfCustomFiledsThanInConfigFormFieldsButAMismatch(t *testing.T) {
	optionalVal := true
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
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID:       "testField",
							Optional: &optionalVal,
						},
						{
							ID: "testField2",
						},
					},
				},
			}),
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

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "password",
				"value": "validpass123",
			},
			{
				"id":    "email",
				"value": "random@gmail.com",
			},
			{
				"id":    "testField",
				"value": "",
			},
			{
				"id":    "testField3",
				"value": "",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()
	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "FIELD_ERROR", data["status"])
	assert.Equal(t, 1, len(data["formFields"].([]interface{})))
	assert.Equal(t, "testField2", data["formFields"].([]interface{})[0].(map[string]interface{})["id"])
	assert.Equal(t, "Field is not optional", data["formFields"].([]interface{})[0].(map[string]interface{})["error"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestCustomFieldValidationError(t *testing.T) {
	customErrorMessage := "testField validation error"
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
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "testField",
							Validate: func(value interface{}) *string {
								if len(value.(string)) <= 5 {
									return &customErrorMessage
								} else {
									return nil
								}
							},
						},
						{
							ID: "testField2",
							Validate: func(value interface{}) *string {
								if len(value.(string)) <= 5 {
									return &customErrorMessage
								} else {
									return nil
								}
							},
						},
					},
				},
			}),
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

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "password",
				"value": "validpass123",
			},
			{
				"id":    "email",
				"value": "random@gmail.com",
			},
			{
				"id":    "testField",
				"value": "test",
			},
			{
				"id":    "testField2",
				"value": "test",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()
	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "FIELD_ERROR", data["status"])
	assert.Equal(t, 2, len(data["formFields"].([]interface{})))

	formFields1 := data["formFields"].([]interface{})

	for _, formField := range formFields1 {
		if formField.(map[string]interface{})["id"] == "testField" {
			assert.Equal(t, "testField validation error", formField.(map[string]interface{})["error"])
		} else {
			assert.Equal(t, "testField validation error", formField.(map[string]interface{})["error"])
		}
	}

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestSignupPasswordFieldValidationError(t *testing.T) {
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

	resp, err := unittesting.SignupRequest("random@gmail.com", "invalid", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "FIELD_ERROR", data["status"])
	assert.Equal(t, 1, len(data["formFields"].([]interface{})))
	assert.Equal(t, "Password must contain at least 8 characters, including a number", data["formFields"].([]interface{})[0].(map[string]interface{})["error"])
	assert.Equal(t, "password", data["formFields"].([]interface{})[0].(map[string]interface{})["id"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestSignupEmailFieldValidationError(t *testing.T) {
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

	resp, err := unittesting.SignupRequest("randomgmail.com", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "FIELD_ERROR", data["status"])
	assert.Equal(t, 1, len(data["formFields"].([]interface{})))
	assert.Equal(t, "Email is invalid", data["formFields"].([]interface{})[0].(map[string]interface{})["error"])
	assert.Equal(t, "email", data["formFields"].([]interface{})[0].(map[string]interface{})["id"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestInputEmailIsTrimmed(t *testing.T) {
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

	resp, err := unittesting.SignupRequest("        random@gmail.com           ", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", data["status"])
	assert.NotNil(t, data["user"].(map[string]interface{})["id"])
	assert.Equal(t, "random@gmail.com", data["user"].(map[string]interface{})["email"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}
