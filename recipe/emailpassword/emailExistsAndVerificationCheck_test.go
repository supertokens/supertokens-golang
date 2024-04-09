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
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

// Email exists tests
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
			Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						*originalImplementation.EmailExistsGET = nil
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

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	q := req.URL.Query()
	q.Add("email", "random@gmail.com")
	req.URL.RawQuery = q.Encode()
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)

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

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)

	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(res.Body)
	assert.Equal(t, "Please provide the email as a GET param", result["message"])
}

func TestGenerateTokenAPIWithValidInputAndEmailNotVerified(t *testing.T) {
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customCSRFval,
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

	resp, err := unittesting.SignupRequest("random@gmail.com", "validPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cookieData := unittesting.ExtractInfoFromResponse(resp)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", result["status"])
	user := result["user"].(map[string]interface{})

	rep1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, user["id"].(string), cookieData["sAccessToken"], cookieData["antiCsrf"])
	if err != nil {
		t.Error(err.Error())
	}
	result1 := *unittesting.HttpResponseToConsumableInformation(rep1.Body)
	assert.Equal(t, "OK", result1["status"])
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customCSRFval,
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

	resp, err := unittesting.SignupRequest("random@gmail.com", "validPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cookieData := unittesting.ExtractInfoFromResponse(resp)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", result["status"])
	user := result["user"].(map[string]interface{})

	verifyToken, err := emailverification.CreateEmailVerificationToken("public", user["id"].(string), nil)
	if err != nil {
		t.Error(err.Error())
	}
	emailverification.VerifyEmailUsingToken("public", verifyToken.OK.Token)

	rep1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, user["id"].(string), cookieData["sAccessToken"], cookieData["antiCsrf"])
	if err != nil {
		t.Error(err.Error())
	}
	result1 := *unittesting.HttpResponseToConsumableInformation(rep1.Body)
	assert.Equal(t, "EMAIL_ALREADY_VERIFIED_ERROR", result1["status"])
}

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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
		},
	}

	BeforeEach()
	unittesting.SetKeyValueInConfig("access_token_validity", strconv.Itoa(2))
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}

	resp.Body.Close()

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", response["status"])

	userId := response["user"].(map[string]interface{})["id"]
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	time.Sleep(5 * time.Second)

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)
	assert.Equal(t, 401, resp1.StatusCode)
	data1, err := io.ReadAll(resp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp1.Body.Close()
	var response1 map[string]interface{}
	err = json.Unmarshal(data1, &response1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "try refresh token", response1["message"])

	res, err := unittesting.SessionRefresh(testServer.URL, cookieData["sRefreshToken"], cookieData["antiCsrf"])
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	cookieData2 := unittesting.ExtractInfoFromResponse(res)

	res1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData2["sAccessToken"], cookieData2["antiCsrf"])
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, res1.StatusCode)

	data2, err := io.ReadAll(res1.Body)
	if err != nil {
		t.Error(err.Error())
	}

	res1.Body.Close()

	var response2 map[string]interface{}
	err = json.Unmarshal(data2, &response2)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", response2["status"])

}

func TestProvidingYourOwnEmailCallBackAndMakeSureItsCalled(t *testing.T) {
	var userInfo evmodels.User
	var emailToken string
	customAntiCsrfVal := "VIA_TOKEN"

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		userInfo = evmodels.User{
			ID:    input.EmailVerification.User.ID,
			Email: input.EmailVerification.User.Email,
		}
		emailToken = input.EmailVerification.EmailVerifyLink
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
			}),
			Init(&epmodels.TypeInput{}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["antiCsrf"])

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

}

func TestEmailVerifyApiWithValidInput(t *testing.T) {
	var token string
	customAntiCsrfVal := "VIA_TOKEN"
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		u, err := url.Parse(input.EmailVerification.EmailVerifyLink)
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
			}),
			Init(&epmodels.TypeInput{}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["antiCsrf"])

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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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
}

func TestThatTheHandlePostEmailVerificationCallBackIsCalledOnSuccessFullVerificationIfGiven(t *testing.T) {
	var userInfoFromCallback evmodels.User
	var token string
	customAntiCsrfVal := "VIA_TOKEN"
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		u, err := url.Parse(input.EmailVerification.EmailVerifyLink)
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				Override: &evmodels.OverrideStruct{
					APIs: func(originalImplementation evmodels.APIInterface) evmodels.APIInterface {
						originalVerifyEmailPost := *originalImplementation.VerifyEmailPOST
						*originalImplementation.VerifyEmailPOST = func(token string, sessionContainer sessmodels.SessionContainer, tenantId string, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.VerifyEmailPOSTResponse, error) {
							res, err := originalVerifyEmailPost(token, sessionContainer, tenantId, options, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							userInfoFromCallback = res.OK.User
							return res, nil
						}
						return originalImplementation
					},
				},
			}),
			Init(&epmodels.TypeInput{}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["antiCsrf"])

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

}

func TestEmailVerifyWithValidInputUsingTheGetMehtod(t *testing.T) {
	var token string
	customAntiCsrfVal := "VIA_TOKEN"
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		u, err := url.Parse(input.EmailVerification.EmailVerifyLink)
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
			}),
			Init(&epmodels.TypeInput{}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["antiCsrf"])

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
	req1.Header.Set("Cookie", "sAccessToken="+cookieData["sAccessToken"])
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
	assert.Equal(t, true, response3["isVerified"])
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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
}

func TestTheEmailVerifyAPIwithValidInputOverridingAPIs(t *testing.T) {
	var token string
	var user evmodels.User
	customAntiCsrfVal := "VIA_TOKEN"
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		u, err := url.Parse(input.EmailVerification.EmailVerifyLink)
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				Override: &evmodels.OverrideStruct{
					APIs: func(originalImplementation evmodels.APIInterface) evmodels.APIInterface {
						originalVerifyEmailPost := *originalImplementation.VerifyEmailPOST
						*originalImplementation.VerifyEmailPOST = func(token string, sessionContainer sessmodels.SessionContainer, tenantId string, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.VerifyEmailPOSTResponse, error) {
							res, err := originalVerifyEmailPost(token, sessionContainer, tenantId, options, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							user = res.OK.User
							return res, nil
						}
						return originalImplementation
					},
				},
			}),
			Init(&epmodels.TypeInput{}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["antiCsrf"])

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
}

func TestTheEmailVerifyAPIwithValidInputAndOverridingFunctions(t *testing.T) {
	var token string
	var user evmodels.User
	customAntiCsrfVal := "VIA_TOKEN"
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		u, err := url.Parse(input.EmailVerification.EmailVerifyLink)
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				Override: &evmodels.OverrideStruct{
					Functions: func(originalImplementation evmodels.RecipeInterface) evmodels.RecipeInterface {
						originalVerifyUsingToken := *originalImplementation.VerifyEmailUsingToken
						*originalImplementation.VerifyEmailUsingToken = func(token string, tenantId string, userContext supertokens.UserContext) (evmodels.VerifyEmailUsingTokenResponse, error) {
							res, err := originalVerifyUsingToken(token, tenantId, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							user = res.OK.User
							return res, nil
						}
						return originalImplementation
					},
				},
			}),
			Init(&epmodels.TypeInput{}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["antiCsrf"])

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
}

func TestTheEmailVerifyAPIwithValidInputThrowsErrorOnSuchOverriding(t *testing.T) {
	var token string
	var user evmodels.User
	customAntiCsrfVal := "VIA_TOKEN"
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		u, err := url.Parse(input.EmailVerification.EmailVerifyLink)
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				Override: &evmodels.OverrideStruct{
					APIs: func(originalImplementation evmodels.APIInterface) evmodels.APIInterface {
						originalVerifyEmailPost := *originalImplementation.VerifyEmailPOST
						*originalImplementation.VerifyEmailPOST = func(token string, sessionContainer sessmodels.SessionContainer, tenantId string, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.VerifyEmailPOSTResponse, error) {
							res, err := originalVerifyEmailPost(token, sessionContainer, tenantId, options, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							user = res.OK.User
							return evmodels.VerifyEmailPOSTResponse{}, errors.New("Verify Email Error")
						}
						return originalImplementation
					},
				},
			}),
			Init(&epmodels.TypeInput{}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["antiCsrf"])

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
}

func TestTheEmailVerifyAPIWithValidInputOverridingFunctionsThrowsError(t *testing.T) {
	var token string
	var user evmodels.User
	customAntiCsrfVal := "VIA_TOKEN"
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		u, err := url.Parse(input.EmailVerification.EmailVerifyLink)
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				Override: &evmodels.OverrideStruct{
					Functions: func(originalImplementation evmodels.RecipeInterface) evmodels.RecipeInterface {
						originalVerifyUsingToken := *originalImplementation.VerifyEmailUsingToken
						*originalImplementation.VerifyEmailUsingToken = func(token string, tenantId string, userContext supertokens.UserContext) (evmodels.VerifyEmailUsingTokenResponse, error) {
							res, err := originalVerifyUsingToken(token, tenantId, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							user = res.OK.User
							return evmodels.VerifyEmailUsingTokenResponse{}, errors.New("Verify Email Error")
						}
						return originalImplementation
					},
				},
			}),

			Init(&epmodels.TypeInput{}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["antiCsrf"])

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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	res, err := emailverification.CreateEmailVerificationToken("public", userId.(string), nil)
	if err != nil {
		t.Error(err.Error())
	}
	verifyToken := res.OK.Token

	emailverification.RevokeEmailVerificationTokens("public", userId.(string), nil)

	res1, err := emailverification.VerifyEmailUsingToken("public", verifyToken)
	assert.NoError(t, err)
	assert.NotNil(t, res1.EmailVerificationInvalidTokenError)
	assert.Nil(t, res1.OK)
}

func TestEmailVerifyWithDeletedUser(t *testing.T) {
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
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}),
			Init(&epmodels.TypeInput{}),
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

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.10", cdiVersion) != cdiVersion {
		return
	}

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
	supertokens.DeleteUser(userId.(string))

	resp1, err := unittesting.EmailVerifyTokenRequest(testServer.URL, userId.(string), cookieData["sAccessToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	cookieData1 := unittesting.ExtractInfoFromResponse(resp1)

	assert.Equal(t, 401, resp1.StatusCode)
	assert.Equal(t, "", cookieData1["sAccessToken"])
	assert.Equal(t, "", cookieData1["sRefreshToken"])

	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["refreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["accessTokenExpiry"])

	assert.Equal(t, "", cookieData1["accessTokenDomain"])
	assert.Equal(t, "", cookieData1["refreshTokenDomain"])
}
