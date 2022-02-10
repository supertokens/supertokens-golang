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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

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

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
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

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
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

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}
