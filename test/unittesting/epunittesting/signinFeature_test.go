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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDisablingAPIDefaultSigninDoesNotWork(t *testing.T) {
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
						*originalImplementation.SignInPOST = nil
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

	res, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

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

func TestSignInAPIworksWithValidInput(t *testing.T) {
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

	res1, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

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

	assert.Equal(t, 200, res1.StatusCode)
	assert.Equal(t, "OK", data1["status"])

	signInUserInfo := data1["user"].(map[string]interface{})

	assert.Equal(t, signupUserInfo["id"], signInUserInfo["id"])
	assert.Equal(t, signupUserInfo["email"], signInUserInfo["email"])
	assert.Equal(t, signupUserInfo["timejoined"], signInUserInfo["timejoined"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestSigninAPIthrowsAnErrorWhenEmailDoesNotMatch(t *testing.T) {
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

	res1, err := unittesting.SignInRequest("ran@gmail.com", "validpass123", testServer.URL)

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

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestSigninAPIThrowsErrorWhenPasswordIsIncorrect(t *testing.T) {
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

	res1, err := unittesting.SignInRequest("random@gmail.com", "validpass12345", testServer.URL)

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

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestBadInputNoPostBodyToSignInAPI(t *testing.T) {
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

	resp, err := http.Post(testServer.URL+"/auth/signin", "application/json", nil)

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

func TestSuccessfullSigInYieldsSession(t *testing.T) {
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

	res1, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

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

	cookieData := unittesting.ExtractInfoFromResponse(res1)

	assert.Equal(t, "OK", data1["status"])

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

func TestCustomEmailValidatorsToSignupAndMakeSureTheyAreAppliedToSignIn(t *testing.T) {
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
							ID: "email",
							Validate: func(value interface{}) *string {
								customErrMessage := "email does not start with test"
								if strings.HasPrefix(value.(string), "test") {
									return nil
								}
								return &customErrMessage
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

	res, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
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

	res1, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

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

	assert.Equal(t, "email does not start with test", data1["formFields"].([]interface{})[0].(map[string]interface{})["error"])
	assert.Equal(t, "email", data1["formFields"].([]interface{})[0].(map[string]interface{})["id"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestCustomPasswordValidatorsToSignupAndMakeSureTheyAreAppliedToSignIn(t *testing.T) {
	failsValidatorCtr := 0
	passesValidatorCtr := 0
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
							ID: "password",
							Validate: func(value interface{}) *string {
								customErrMessage := "password is greater than 5 characters"
								if len(value.(string)) <= 5 {
									passesValidatorCtr++
									return nil
								}
								failsValidatorCtr++
								return &customErrMessage
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

	res, err := unittesting.SignupRequest("random@gmail.com", "valid", testServer.URL)
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
	assert.Equal(t, 1, passesValidatorCtr)
	assert.Equal(t, 0, failsValidatorCtr)

	res1, err := unittesting.SignInRequest("random@gmail.com", "invalid", testServer.URL)

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
	assert.Equal(t, 1, passesValidatorCtr)
	assert.Equal(t, 0, failsValidatorCtr)

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestPasswordFieldValidationError(t *testing.T) {
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

	res1, err := unittesting.SignInRequest("random@gmail.com", "invalidpass", testServer.URL)

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

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestEmailFieldValidationError(t *testing.T) {
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

	res, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
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

	res1, err := unittesting.SignInRequest("randomgmail.com", "validpass123", testServer.URL)

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

	assert.Equal(t, "Email is invalid", data1["formFields"].([]interface{})[0].(map[string]interface{})["error"])
	assert.Equal(t, "email", data1["formFields"].([]interface{})[0].(map[string]interface{})["id"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestFormFieldsHasNoEmailField(t *testing.T) {
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

	res, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
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

	resp, err := http.Post(testServer.URL+"/auth/signin", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes1, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err.Error())
	}

	resp.Body.Close()

	assert.Equal(t, "Are you sending too many / too few formFields?\n", string(dataInBytes1))
	assert.Equal(t, 500, resp.StatusCode)

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestFormFieldsHasNoPasswordField(t *testing.T) {
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

	res, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
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

	resp, err := http.Post(testServer.URL+"/auth/signin", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes1, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err.Error())
	}

	resp.Body.Close()

	assert.Equal(t, "Are you sending too many / too few formFields?\n", string(dataInBytes1))
	assert.Equal(t, 500, resp.StatusCode)

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestInvalidEmailAndWrongPassword(t *testing.T) {
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

	res, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
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

	res1, err := unittesting.SignInRequest("randomgmail.com", "invalid", testServer.URL)

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
	assert.Equal(t, "Email is invalid", data1["formFields"].([]interface{})[0].(map[string]interface{})["error"])
	assert.Equal(t, "email", data1["formFields"].([]interface{})[0].(map[string]interface{})["id"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestGetUserByEmail(t *testing.T) {
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

	user, err := emailpassword.GetUserByEmail("random@gmail.com")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Nil(t, user)

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

	user1, err := emailpassword.GetUserByEmail("random@gmail.com")
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, user1.Email, data["user"].(map[string]interface{})["email"])
	assert.Equal(t, user1.ID, data["user"].(map[string]interface{})["id"])

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "OK", data["status"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestGetUserById(t *testing.T) {
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

	user, err := emailpassword.GetUserByID("randomId")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Nil(t, user)

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

	user1, err := emailpassword.GetUserByID(data["user"].(map[string]interface{})["id"].(string))

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, user1.Email, data["user"].(map[string]interface{})["email"])
	assert.Equal(t, user1.ID, data["user"].(map[string]interface{})["id"])

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "OK", data["status"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestHandlePostSignInFunction(t *testing.T) {
	var customUser epmodels.User
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
						originalSignInPost := *originalImplementation.SignInPOST
						*originalImplementation.SignInPOST = func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.SignInResponse, error) {
							res, _ := originalSignInPost(formFields, options)
							customUser = res.OK.User
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

	_, err = unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	res, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)
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

	assert.NotNil(t, customUser)

	assert.Equal(t, customUser.Email, data["user"].(map[string]interface{})["email"])
	assert.Equal(t, customUser.ID, data["user"].(map[string]interface{})["id"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}
