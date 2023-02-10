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

package thirdpartypasswordless

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"gopkg.in/h2non/gock.v1"
)

func TestTheDefaultRouteAndItShouldRevokeTheSession(t *testing.T) {
	customAntiCsrfValue := "VIA_TOKEN"
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfValue,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					signinupCustomProvider1,
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	defer gock.OffAll()
	gock.New("https://test.com/").
		Post("oauth/token").
		Reply(200).
		JSON(map[string]string{})

	postData := map[string]interface{}{
		"thirdPartyId": "custom",
		"redirectURIInfo": map[string]interface{}{
			"redirectURIOnProviderDashboard": "http://127.0.0.1/callback",
			"redirectURIQueryParams": map[string]interface{}{
				"code": "32432432",
			},
		},
	}

	postBody, err := json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}

	gock.New(testServer.URL).EnableNetworking().Persist()
	gock.New("http://localhost:8080/").EnableNetworking().Persist()

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", result["status"])

	cookieData := unittesting.ExtractInfoFromResponse(resp)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/signout", nil)
	if err != nil {
		t.Error(err.Error())
	}

	req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
	req.Header.Add("anti-csrf", cookieData["antiCsrf"])

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)

	cookieData1 := unittesting.ExtractInfoFromResponseWhenAntiCSRFisNone(res)

	assert.Equal(t, "", cookieData1["sAccessToken"])
	assert.Equal(t, "", cookieData1["sRefreshToken"])

	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["refreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["accessTokenExpiry"])

	assert.Equal(t, "", cookieData1["accessTokenDomain"])
	assert.Equal(t, "", cookieData1["refreshTokenDomain"])
}

func TestDisablingDefaultRouteAndCallingTheAPIReturns404(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Override: &tplmodels.OverrideStruct{
					APIs: func(originalImplementation tplmodels.APIInterface) tplmodels.APIInterface {
						*originalImplementation.ThirdPartySignInUpPOST = nil
						return originalImplementation
					},
				},
				Providers: []tpmodels.ProviderInput{
					signinupCustomProvider1,
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/signout", nil)
	if err != nil {
		t.Error(err.Error())
	}

	req.Header.Add("rid", "thirdparty")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestCallingAPIWithoutSessionShouldReturnOk(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					signinupCustomProvider1,
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := http.Post(testServer.URL+"/auth/signout", "application/json", nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, len(resp.Cookies()))
	assert.Equal(t, "", resp.Header.Get("set-cookie"))

	response := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", response["status"])
}

func TestThatSignoutAPIreturnsTryRefreshTokenRefreshSessionAndSignoutShouldReturnOk(t *testing.T) {
	customAntiCsrfValue := "VIA_TOKEN"
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfValue,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					signinupCustomProvider1,
				},
			}),
		},
	}

	BeforeEach()
	unittesting.SetKeyValueInConfig("access_token_validity", "2")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	defer gock.OffAll()
	gock.New("https://test.com/").
		Post("oauth/token").
		Reply(200).
		JSON(map[string]string{})

	postData := map[string]interface{}{
		"thirdPartyId": "custom",
		"redirectURIInfo": map[string]interface{}{
			"redirectURIOnProviderDashboard": "http://127.0.0.1/callback",
			"redirectURIQueryParams": map[string]interface{}{
				"code": "32432432",
			},
		},
	}

	postBody, err := json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}

	gock.New(testServer.URL).EnableNetworking().Persist()
	gock.New("http://localhost:8080/").EnableNetworking().Persist()

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	cookieData := unittesting.ExtractInfoFromResponse(resp)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", result["status"])

	time.Sleep(5 * time.Second)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/signout", nil)
	if err != nil {
		t.Error(err.Error())
	}

	req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
	req.Header.Add("anti-csrf", cookieData["antiCsrf"])

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	result1 := *unittesting.HttpResponseToConsumableInformation(res.Body)
	assert.Equal(t, "try refresh token", result1["message"])

	req1, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	if err != nil {
		t.Error(err.Error())
	}
	req1.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])
	req1.Header.Add("anti-csrf", cookieData["antiCsrf"])

	res1, err := http.DefaultClient.Do(req1)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, res1.StatusCode)
	cookieData1 := unittesting.ExtractInfoFromResponse(res1)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/signout", nil)
	if err != nil {
		t.Error(err.Error())
	}
	req2.Header.Add("Cookie", "sAccessToken="+cookieData1["sAccessToken"])
	req2.Header.Add("anti-csrf", cookieData1["antiCsrf"])

	res2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, res2.StatusCode)

	cookieData2 := unittesting.ExtractInfoFromResponseWhenAntiCSRFisNone(res2)

	assert.Equal(t, "", cookieData2["sAccessToken"])
	assert.Equal(t, "", cookieData2["sRefreshToken"])

	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData2["refreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData2["accessTokenExpiry"])

	assert.Equal(t, "", cookieData2["accessTokenDomain"])
	assert.Equal(t, "", cookieData2["refreshTokenDomain"])
}
