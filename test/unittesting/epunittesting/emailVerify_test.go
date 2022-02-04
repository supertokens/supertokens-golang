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
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestGenerateTokenAPIWithValidInputAndEmailNotVerified(t *testing.T) {
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

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	verifyToken, err := emailpassword.CreateEmailVerificationToken(userId.(string))
	if err != nil {
		t.Error(err.Error())
	}
	emailpassword.VerifyEmailUsingToken(verifyToken.OK.Token)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)
	data1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	var response1 map[string]interface{}
	_ = json.Unmarshal(data1, &response1)

	assert.Equal(t, "EMAIL_ALREADY_VERIFIED_ERROR", response1["status"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestGenerateTokenAPIWithValidInputNoSessionAndCheckOutput(t *testing.T) {
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

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/user/email/verify/token", nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "unauthorised", response["message"])
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

//!THIS NEEDS TO FIGURED OUT
func TestGenerateTokenAPIWithExpiredAccessToken(t *testing.T) {
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
	unittesting.SetKeyAndNumberValueInConfig("access_token_validity", 2)
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(resp1)
	// assert.NoError(t, err)
	// assert.Equal(t, 200, resp1.StatusCode)
	// data1, _ := io.ReadAll(resp1.Body)
	// resp1.Body.Close()
	// var response1 map[string]interface{}
	// _ = json.Unmarshal(data1, &response1)

	// assert.Equal(t, "EMAIL_ALREADY_VERIFIED_ERROR", response1["status"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestProvidingYourOwnEmailCallBackAndMakeSureItsCalled(t *testing.T) {
	var userInfo epmodels.User
	var emailToken string
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
				EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, emailVerificationURLWithToken string) {
						userInfo = user
						emailToken = emailVerificationURLWithToken
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

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)
	data1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	var response1 map[string]interface{}
	_ = json.Unmarshal(data1, &response1)

	assert.Equal(t, "OK", response1["status"])
	assert.Equal(t, "test@gmail.com", userInfo.Email)
	assert.Equal(t, userId, userInfo.ID)
	assert.NotNil(t, emailToken)

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestEmailVerifyApiWithValidInput(t *testing.T) {
	var token string
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
				EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, emailVerificationURLWithToken string) {
						token = strings.Split(strings.Split(emailVerificationURLWithToken, "?token=")[1], "&rid=")[0]
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

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)
	data1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	var response1 map[string]interface{}
	_ = json.Unmarshal(data1, &response1)

	assert.Equal(t, "OK", response1["status"])

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/user/email/verify", strings.NewReader(`{"method":"token","token":"`+token+`"}`))
	if err != nil {
		t.Error(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(datainBytes, &response2)
	assert.Equal(t, "OK", response2["status"])
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestTheEmailVerifyApiWithInvalidTokenAndCheckError(t *testing.T) {
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

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/user/email/verify", strings.NewReader(`{"method":"token","token":"randomToken"}`))
	if err != nil {
		t.Error(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(datainBytes, &response2)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "EMAIL_VERIFICATION_INVALID_TOKEN_ERROR", response2["status"])
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestEmailVerifyAPIWithTokenOfNotTypeString(t *testing.T) {
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

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/user/email/verify", strings.NewReader(`{"method":"token","token":`+strconv.Itoa(2000)+`}`))
	if err != nil {
		t.Error(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(datainBytes, &response2)
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "The email verification token must be a string", response2["message"])
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestThatTheHandlePostEmailVerificationCallBackIsCalledOnSuccessFullVerificationIfGiven(t *testing.T) {
	var userInfoFromCallback evmodels.User
	var token string
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
				EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, emailVerificationURLWithToken string) {
						token = strings.Split(strings.Split(emailVerificationURLWithToken, "?token=")[1], "&rid=")[0]
					},
				},
				Override: &epmodels.OverrideStruct{
					EmailVerificationFeature: &evmodels.OverrideStruct{
						APIs: func(originalImplementation evmodels.APIInterface) evmodels.APIInterface {
							originalVerifyEmailPost := *originalImplementation.VerifyEmailPOST
							*originalImplementation.VerifyEmailPOST = func(token string, options evmodels.APIOptions) (evmodels.VerifyEmailUsingTokenResponse, error) {
								res, err := originalVerifyEmailPost(token, options)
								if err != nil {
									log.Fatal(err.Error())
								}
								userInfoFromCallback = res.OK.User
								return res, nil
							}
							return originalImplementation
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

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)
	data1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	var response1 map[string]interface{}
	_ = json.Unmarshal(data1, &response1)

	assert.Equal(t, "OK", response1["status"])
	assert.NotNil(t, token)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/user/email/verify", strings.NewReader(`{"method":"token","token":"`+token+`"}`))
	if err != nil {
		t.Error(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(datainBytes, &response2)
	assert.Equal(t, "OK", response2["status"])

	assert.Equal(t, "test@gmail.com", userInfoFromCallback.Email)
	assert.Equal(t, userId, userInfoFromCallback.ID)

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestEmailVerifyWithValidInputUsingTheGetMehtod(t *testing.T) {
	var token string
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
				EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, emailVerificationURLWithToken string) {
						token = strings.Split(strings.Split(emailVerificationURLWithToken, "?token=")[1], "&rid=")[0]
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

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)
	data1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	var response1 map[string]interface{}
	_ = json.Unmarshal(data1, &response1)

	assert.Equal(t, "OK", response1["status"])

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/user/email/verify", strings.NewReader(`{"method":"token","token":"`+token+`"}`))
	if err != nil {
		t.Error(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(datainBytes, &response2)
	assert.Equal(t, "OK", response2["status"])

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/user/email/verify", nil)
	req1.Header.Set("Cookie", "sAccessToken="+cookieData["sAccessToken"]+"; sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req1.Header.Add("anti-csrf", cookieData["antiCsrf"])
	if err != nil {
		t.Error(err.Error())
	}
	res1, err := http.DefaultClient.Do(req1)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes1, err := io.ReadAll(res1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res1.Body.Close()
	var response3 map[string]interface{}
	_ = json.Unmarshal(datainBytes1, &response3)
	assert.Equal(t, "OK", response3["status"])
	fmt.Println(response3)
	assert.Equal(t, true, response3["isVerified"])
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestVerifySessionWithNoSessionUsingTheGetMethod(t *testing.T) {
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

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/user/email/verify", nil)
	if err != nil {
		t.Error(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(datainBytes, &response2)
	assert.Equal(t, 401, res.StatusCode)
	assert.Equal(t, "unauthorised", response2["message"])
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

//!needs to figured out
func TestTheEmailVerifyWithAnExpiredAccessTokenUsingTheGetMethod(t *testing.T) {

}

func TestTheEmailVerifyAPIwithValidInputOverridingAPIs(t *testing.T) {
	var token string
	var user evmodels.User
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
				EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, emailVerificationURLWithToken string) {
						token = strings.Split(strings.Split(emailVerificationURLWithToken, "?token=")[1], "&rid=")[0]
					},
				},
				Override: &epmodels.OverrideStruct{
					EmailVerificationFeature: &evmodels.OverrideStruct{
						APIs: func(originalImplementation evmodels.APIInterface) evmodels.APIInterface {
							originalVerifyEmailPost := *originalImplementation.VerifyEmailPOST
							*originalImplementation.VerifyEmailPOST = func(token string, options evmodels.APIOptions) (evmodels.VerifyEmailUsingTokenResponse, error) {
								res, err := originalVerifyEmailPost(token, options)
								if err != nil {
									log.Fatal(err.Error())
								}
								user = res.OK.User
								return res, nil
							}
							return originalImplementation
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

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)
	data1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	var response1 map[string]interface{}
	_ = json.Unmarshal(data1, &response1)

	assert.Equal(t, "OK", response1["status"])
	assert.NotNil(t, token)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/user/email/verify", strings.NewReader(`{"method":"token","token":"`+token+`"}`))
	if err != nil {
		t.Error(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(datainBytes, &response2)
	assert.Equal(t, "OK", response2["status"])

	assert.Equal(t, "test@gmail.com", user.Email)
	assert.Equal(t, userId, user.ID)
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestTheEmailVerifyAPIwithValidInputAndOverridingFunctions(t *testing.T) {
	var token string
	var user evmodels.User
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
				EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, emailVerificationURLWithToken string) {
						token = strings.Split(strings.Split(emailVerificationURLWithToken, "?token=")[1], "&rid=")[0]
					},
				},
				Override: &epmodels.OverrideStruct{
					EmailVerificationFeature: &evmodels.OverrideStruct{
						Functions: func(originalImplementation evmodels.RecipeInterface) evmodels.RecipeInterface {
							originalVerifyUsingToken := *originalImplementation.VerifyEmailUsingToken
							*originalImplementation.VerifyEmailUsingToken = func(token string) (evmodels.VerifyEmailUsingTokenResponse, error) {
								res, err := originalVerifyUsingToken(token)
								if err != nil {
									log.Fatal(err.Error())
								}
								user = res.OK.User
								return res, nil
							}
							return originalImplementation
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

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)
	data1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	var response1 map[string]interface{}
	_ = json.Unmarshal(data1, &response1)

	assert.Equal(t, "OK", response1["status"])
	assert.NotNil(t, token)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/user/email/verify", strings.NewReader(`{"method":"token","token":"`+token+`"}`))
	if err != nil {
		t.Error(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	var response2 map[string]interface{}
	_ = json.Unmarshal(datainBytes, &response2)
	assert.Equal(t, "OK", response2["status"])

	assert.Equal(t, "test@gmail.com", user.Email)
	assert.Equal(t, userId, user.ID)
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestTheEmailVerifyAPIwithValidInputThrowsErrorOnSuchOverriding(t *testing.T) {
	var token string
	var user evmodels.User
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
				EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, emailVerificationURLWithToken string) {
						token = strings.Split(strings.Split(emailVerificationURLWithToken, "?token=")[1], "&rid=")[0]
					},
				},
				Override: &epmodels.OverrideStruct{
					EmailVerificationFeature: &evmodels.OverrideStruct{
						APIs: func(originalImplementation evmodels.APIInterface) evmodels.APIInterface {
							originalVerifyEmailPost := *originalImplementation.VerifyEmailPOST
							*originalImplementation.VerifyEmailPOST = func(token string, options evmodels.APIOptions) (evmodels.VerifyEmailUsingTokenResponse, error) {
								res, err := originalVerifyEmailPost(token, options)
								if err != nil {
									log.Fatal(err.Error())
								}
								user = res.OK.User
								return evmodels.VerifyEmailUsingTokenResponse{}, errors.New("Verify Email Error")
							}
							return originalImplementation
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

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)
	data1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	var response1 map[string]interface{}
	_ = json.Unmarshal(data1, &response1)

	assert.Equal(t, "OK", response1["status"])
	assert.NotNil(t, token)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/user/email/verify", strings.NewReader(`{"method":"token","token":"`+token+`"}`))
	if err != nil {
		t.Error(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()

	assert.Equal(t, "Verify Email Error\n", string(datainBytes))
	assert.Equal(t, "test@gmail.com", user.Email)
	assert.Equal(t, userId, user.ID)
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestTheEmailVerifyAPIWithValidInputOverridingFunctionsThrowsError(t *testing.T) {
	var token string
	var user evmodels.User
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
				EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, emailVerificationURLWithToken string) {
						token = strings.Split(strings.Split(emailVerificationURLWithToken, "?token=")[1], "&rid=")[0]
					},
				},
				Override: &epmodels.OverrideStruct{
					EmailVerificationFeature: &evmodels.OverrideStruct{
						Functions: func(originalImplementation evmodels.RecipeInterface) evmodels.RecipeInterface {
							originalVerifyUsingToken := *originalImplementation.VerifyEmailUsingToken
							*originalImplementation.VerifyEmailUsingToken = func(token string) (evmodels.VerifyEmailUsingTokenResponse, error) {
								res, err := originalVerifyUsingToken(token)
								if err != nil {
									log.Fatal(err.Error())
								}
								user = res.OK.User
								return evmodels.VerifyEmailUsingTokenResponse{}, errors.New("Verify Email Error")
							}
							return originalImplementation
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

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["sIdRefreshToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)
	data1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	var response1 map[string]interface{}
	_ = json.Unmarshal(data1, &response1)

	assert.Equal(t, "OK", response1["status"])
	assert.NotNil(t, token)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/user/email/verify", strings.NewReader(`{"method":"token","token":"`+token+`"}`))
	if err != nil {
		t.Error(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	datainBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()

	assert.Equal(t, "Verify Email Error\n", string(datainBytes))
	assert.Equal(t, "test@gmail.com", user.Email)
	assert.Equal(t, userId, user.ID)
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestTheGenerateTokenAPIWithValidInputAndThenRemoveToken(t *testing.T) {
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

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]

	res, err := emailpassword.CreateEmailVerificationToken(userId.(string))
	if err != nil {
		t.Error(err.Error())
	}
	verifyToken := res.OK.Token

	emailpassword.RevokeEmailVerificationTokens(userId.(string))

	res1, err := emailpassword.VerifyEmailUsingToken(verifyToken)
	assert.Nil(t, res1)
	assert.Equal(t, "email verification token is invalid", err.Error())
	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}
