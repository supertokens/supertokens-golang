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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

//SigninFeature Tests
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

	res, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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
			Init(&epmodels.TypeInput{
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
			Init(&epmodels.TypeInput{
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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

	user, err := GetUserByEmail("random@gmail.com")
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

	user1, err := GetUserByEmail("random@gmail.com")
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, user1.Email, data["user"].(map[string]interface{})["email"])
	assert.Equal(t, user1.ID, data["user"].(map[string]interface{})["id"])

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "OK", data["status"])

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
			Init(nil),
			session.Init(nil),
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

	user, err := GetUserByID("randomId")
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

	user1, err := GetUserByID(data["user"].(map[string]interface{})["id"].(string))

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, user1.Email, data["user"].(map[string]interface{})["email"])
	assert.Equal(t, user1.ID, data["user"].(map[string]interface{})["id"])

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "OK", data["status"])

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
			Init(&epmodels.TypeInput{
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

}

//Signout Feature tests
func TestDefaultSignoutRouteRevokesSession(t *testing.T) {
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
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	cookieData := unittesting.ExtractInfoFromResponse(res)

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

	res1, err := unittesting.SignoutRequest(testServer.URL, cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}

	cookieData1 := unittesting.ExtractInfoFromResponseWhenAntiCSRFisNone(res1)

	assert.Equal(t, "", cookieData1["sAccessToken"])
	assert.Equal(t, "", cookieData1["sRefreshToken"])
	assert.Equal(t, "", cookieData1["sIdRefreshToken"])

	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["refreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["accessTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["idRefreshTokenExpiry"])

	assert.Equal(t, "", cookieData1["accessTokenDomain"])
	assert.Equal(t, "", cookieData1["refreshTokenDomain"])
	assert.Equal(t, "", cookieData1["idRefreshTokenDomain"])
}

func TestCallingTheAPIwithoutSessionShouldReturnOk(t *testing.T) {
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
			session.Init(nil),
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
	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/signout", nil)

	if err != nil {
		t.Error(err.Error())
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)

	if err != nil {
		t.Error(err.Error())
	}

	dataInbytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInbytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "OK", data["status"])
	assert.Nil(t, req.Header["Cookie"])
}

func TestSignoutAPIreturnsTryRefreshTokenAndSignoutShouldReturnOK(t *testing.T) {
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
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}

	BeforeEach()

	unittesting.SetKeyAndNumberValueInConfig("access_token_validity", 2)

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

	cookieData := unittesting.ExtractInfoFromResponse(res)

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

	time.Sleep(5 * time.Second)

	res1, err := unittesting.SignoutRequest(testServer.URL, cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusUnauthorized, res1.StatusCode)

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
	assert.Equal(t, "try refresh token", data1["message"])

	res2, err := unittesting.SessionRefresh(testServer.URL, cookieData["sRefreshToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}

	cookieData1 := unittesting.ExtractInfoFromResponse(res2)

	res3, err := unittesting.SignoutRequest(testServer.URL, cookieData1["sAccessToken"], cookieData1["sIdRefreshToken"], cookieData1["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}

	cookieData2 := unittesting.ExtractInfoFromResponseWhenAntiCSRFisNone(res3)

	assert.Equal(t, "", cookieData2["sAccessToken"])
	assert.Equal(t, "", cookieData2["sRefreshToken"])
	assert.Equal(t, "", cookieData2["sIdRefreshToken"])

	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData2["refreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData2["accessTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData2["idRefreshTokenExpiry"])

	assert.Equal(t, "", cookieData2["accessTokenDomain"])
	assert.Equal(t, "", cookieData2["refreshTokenDomain"])
	assert.Equal(t, "", cookieData2["idRefreshTokenDomain"])
}

//Signup Feature tests
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
			Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						*originalImplementation.SignUpPOST = nil
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

	res, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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
			Init(&epmodels.TypeInput{
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
			Init(&epmodels.TypeInput{
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
			Init(&epmodels.TypeInput{
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

	dataInBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	assert.Equal(t, 500, res.StatusCode)
	assert.Equal(t, "Are you sending too many / too few formFields?\n", string(dataInBytes))

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
			Init(&epmodels.TypeInput{
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
			Init(&epmodels.TypeInput{
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(&epmodels.TypeInput{
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
			Init(&epmodels.TypeInput{
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
			Init(&epmodels.TypeInput{
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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
			Init(nil),
			session.Init(nil),
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

}
