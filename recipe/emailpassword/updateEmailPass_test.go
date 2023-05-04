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
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.7", cdiVersion) == "2.7" {
		return
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

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
	password := "testPass1"

	UpdateEmailOrPassword(data["user"].(map[string]interface{})["id"].(string), &email, &password, nil)
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

	password = "test"
	applyPasswordPolicy := true
	res3, err := UpdateEmailOrPassword(data["user"].(map[string]interface{})["id"].(string), &email, &password, &applyPasswordPolicy)
	assert.NotNil(t, res3.PasswordPolicyViolatedError)
	assert.Equal(t, "Password must contain at least 8 characters, including a number", res3.PasswordPolicyViolatedError.FailureReason)
}

func TestAPICustomResponse(t *testing.T) {
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
						oSignUpPost := originalImplementation.SignUpPOST
						nSignUpPost := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
							options.Res.Header().Set("Content-Type", "application/json; charset=utf-8")
							options.Res.WriteHeader(201)
							responseJson := map[string]interface{}{
								"message": "My custom response",
							}
							bytes, _ := json.Marshal(responseJson)
							options.Res.Write(bytes)
							return (*oSignUpPost)(formFields, options, userContext)
						}
						originalImplementation.SignUpPOST = &nSignUpPost
						return originalImplementation
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.7", cdiVersion) == "2.7" {
		return
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	res, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 201, res.StatusCode)
	dataInBytes, err := io.ReadAll(res.Body)
	data := map[string]interface{}{}
	json.Unmarshal(dataInBytes, &data)
	assert.Equal(t, "My custom response", data["message"])
}

func TestAPICustomResponseGeneralError(t *testing.T) {
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
						nSignUpPost := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
							options.Res.Header().Set("Content-Type", "application/json; charset=utf-8")
							options.Res.WriteHeader(201)
							responseJson := map[string]interface{}{
								"message": "My custom response",
							}
							bytes, _ := json.Marshal(responseJson)
							options.Res.Write(bytes)
							return epmodels.SignUpPOSTResponse{}, errors.New("My custom error")
						}
						originalImplementation.SignUpPOST = &nSignUpPost
						return originalImplementation
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.7", cdiVersion) == "2.7" {
		return
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	res, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 201, res.StatusCode)
	dataInBytes, err := io.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	data := map[string]interface{}{}
	json.Unmarshal(dataInBytes, &data)
	assert.Equal(t, "My custom response", data["message"])
}

func TestAPICustomResponseMalformedResult(t *testing.T) {
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
						nSignUpPost := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
							options.Res.Header().Set("Content-Type", "application/json; charset=utf-8")
							options.Res.WriteHeader(201)
							responseJson := map[string]interface{}{
								"message": "My custom response",
							}
							bytes, _ := json.Marshal(responseJson)
							options.Res.Write(bytes)
							return epmodels.SignUpPOSTResponse{}, nil
						}
						originalImplementation.SignUpPOST = &nSignUpPost
						return originalImplementation
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.7", cdiVersion) == "2.7" {
		return
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	res, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 201, res.StatusCode)
	dataInBytes, err := io.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	data := map[string]interface{}{}
	json.Unmarshal(dataInBytes, &data)
	assert.Equal(t, "My custom response", data["message"])
}

func TestAPICustomResponseMalformedResultWithoutCustomResponse(t *testing.T) {
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
						nSignUpPost := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
							return epmodels.SignUpPOSTResponse{}, nil
						}
						originalImplementation.SignUpPOST = &nSignUpPost
						return originalImplementation
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.7", cdiVersion) == "2.7" {
		return
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	res, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 500, res.StatusCode)
	dataInBytes, err := io.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "invalid return from API interface function\n", string(dataInBytes))
}

func TestAPIRequestBodyInAPIOverride(t *testing.T) {
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
						oSignUpPost := *originalImplementation.SignUpPOST
						nSignUpPost := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
							body := options.Req.Body
							bodyBytes, err := ioutil.ReadAll(body)
							assert.Nil(t, err)
							requestBody := map[string]interface{}{}
							json.Unmarshal(bodyBytes, &requestBody)

							requestFormFields := requestBody["formFields"].([]interface{})

							for _, formField := range requestFormFields {
								formFieldMap := formField.(map[string]interface{})
								if formFieldMap["id"] == "email" {
									assert.Equal(t, formFieldMap["value"], "testrandom@gmail.com")
								} else if formFieldMap["id"] == "password" {
									assert.Equal(t, formFieldMap["value"], "validpass123")
								}
							}

							return oSignUpPost(formFields, options, userContext)
						}
						originalImplementation.SignUpPOST = &nSignUpPost
						return originalImplementation
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.7", cdiVersion) == "2.7" {
		return
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	res, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, nil, err)
}
