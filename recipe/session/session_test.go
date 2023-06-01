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

package session

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestOutputHeadersAndSetCookieForCreateSessionIsFine(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "rope", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)
	assert.Equal(t, []string{"front-token, anti-csrf"}, res.Header["Access-Control-Expose-Headers"])
	assert.Equal(t, "", cookieData["refreshTokenDomain"])
	assert.Equal(t, "", cookieData["accessTokenDomain"])
	assert.NotNil(t, cookieData["sAccessToken"])
	assert.NotNil(t, cookieData["sRefreshToken"])
	assert.NotNil(t, cookieData["antiCsrf"])
	assert.NotNil(t, cookieData["accessTokenExpiry"])
	assert.NotNil(t, cookieData["refreshTokenExpiry"])
}

func TestTokenTheftDetection(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "user", map[string]interface{}{}, map[string]interface{}{})
	})

	customValForAntiCsrfCheck := true
	customSessionRequiredValue := true
	mux.HandleFunc("/verifySession", VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue,
		AntiCsrfCheck:   &customValForAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue,
			AntiCsrfCheck:   &customValForAntiCsrfCheck,
		})
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req2.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])
	req2.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res2, err := http.DefaultClient.Do(req2)
	cookieData2 := unittesting.ExtractInfoFromResponse(res2)
	assert.NoError(t, err)

	reqV, err := http.NewRequest(http.MethodGet, testServer.URL+"/verifySession", nil)
	assert.NoError(t, err)
	reqV.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"])
	reqV.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	resv, err := http.DefaultClient.Do(reqV)
	assert.NoError(t, err)
	assert.Equal(t, resv.StatusCode, 200)

	req3, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req3.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])
	req3.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res3, err := http.DefaultClient.Do(req3)
	assert.NoError(t, err)
	var jsonResponse map[string]interface{}
	err = json.NewDecoder(res3.Body).Decode(&jsonResponse)
	if err != nil {
		t.Error(err.Error())
	}
	res3.Body.Close()
	assert.Equal(t, "token theft detected", jsonResponse["message"])
	assert.Equal(t, 401, res3.StatusCode)
	assert.NoError(t, err)
}

func TestTokenTheftDetectionWithAPIKey(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
			APIKey:        "shfo3h98308hOIHoei309saiho",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
		},
	}
	BeforeEach()
	unittesting.SetKeyValueInConfig("api_keys", "shfo3h98308hOIHoei309saiho")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "userId", map[string]interface{}{}, map[string]interface{}{})
	})
	customValForAntiCsrfCheck := true
	customSessionRequiredValue := true
	mux.HandleFunc("/verifySession", VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue,
		AntiCsrfCheck:   &customValForAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue,
			AntiCsrfCheck:   &customValForAntiCsrfCheck,
		})
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req2.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])
	req2.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res2, err := http.DefaultClient.Do(req2)
	cookieData2 := unittesting.ExtractInfoFromResponse(res2)
	assert.NoError(t, err)

	reqV, err := http.NewRequest(http.MethodGet, testServer.URL+"/verifySession", nil)
	assert.NoError(t, err)
	reqV.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"])
	reqV.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	resv, err := http.DefaultClient.Do(reqV)
	assert.NoError(t, err)
	assert.Equal(t, resv.StatusCode, 200)

	req3, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req3.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])
	req3.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res3, err := http.DefaultClient.Do(req3)
	assert.NoError(t, err)
	var jsonResponse map[string]interface{}
	err = json.NewDecoder(res3.Body).Decode(&jsonResponse)
	if err != nil {
		t.Error(err.Error())
	}
	res3.Body.Close()
	assert.Equal(t, "token theft detected", jsonResponse["message"])
	assert.Equal(t, 401, res3.StatusCode)
	assert.NoError(t, err)
}

func TestSessionVerificationWithoutAntiCsrfPresent(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "someId", map[string]interface{}{}, map[string]interface{}{})
	})
	customValForAntiCsrfCheck := true
	customSessionRequiredValue := true
	mux.HandleFunc("/getSession", VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue,
		AntiCsrfCheck:   &customValForAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue,
			AntiCsrfCheck:   &customValForAntiCsrfCheck,
		})
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/getSession", nil)
	assert.NoError(t, err)
	req1.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
	res1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, 401, res1.StatusCode)
}

func TestRevokingOfSessions(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "someUniqueID", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	_, err = RevokeAllSessionsForUser("someUniqueID")
	if err != nil {
		t.Error(err.Error())
	}

	sessionHandlesAfterRevoke, err := GetAllSessionHandlesForUser("someUniqueID")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 0, len(sessionHandlesAfterRevoke))

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, 200, res1.StatusCode)

	req2, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, res2.StatusCode)

	sessionHandlesBeforeRevoke1, err := GetAllSessionHandlesForUser("someUniqueID")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 2, len(sessionHandlesBeforeRevoke1))

	revokedSessions, err := RevokeAllSessionsForUser("someUniqueID")
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 2, len(revokedSessions))

	sessionHandlesAfterRevoke1, err := GetAllSessionHandlesForUser("someUniqueID")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 0, len(sessionHandlesAfterRevoke1))
}

func TestManipulatingSessionData(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "rp", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	sessionHandles, err := GetAllSessionHandlesForUser("rp")

	if err != nil {
		t.Error(err.Error())
	}

	UpdateSessionDataInDatabase(sessionHandles[0], map[string]interface{}{
		"name": "John",
	})

	sessionInfo, err := GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, "John", sessionInfo.SessionDataInDatabase["name"])

	UpdateSessionDataInDatabase(sessionHandles[0], map[string]interface{}{
		"name": "Joel",
	})

	sessionInfo, err = GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, "Joel", sessionInfo.SessionDataInDatabase["name"])

	//update session data with wrong session handle

	sessionUpdated, err := UpdateSessionDataInDatabase("random", map[string]interface{}{
		"name": "Ronit",
	})

	assert.NoError(t, err)
	assert.False(t, sessionUpdated)
}

func TestNilValuesPassedForSessionData(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "uniqueId", map[string]interface{}{}, nil)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	sessionHandles, err := GetAllSessionHandlesForUser("uniqueId")

	if err != nil {
		t.Error(err.Error())
	}

	sessionInfo, err := GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, map[string]interface{}{}, sessionInfo.SessionDataInDatabase)

	UpdateSessionDataInDatabase(sessionHandles[0], map[string]interface{}{
		"name": "John",
	})
	sessionInfo, err = GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, "John", sessionInfo.SessionDataInDatabase["name"])
}

func TestManipulatingJWTpayload(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "uniqueId", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	sessionHandles, err := GetAllSessionHandlesForUser("uniqueId")

	if err != nil {
		t.Error(err.Error())
	}

	tokenUpdated, err := MergeIntoAccessTokenPayload(sessionHandles[0], map[string]interface{}{
		"key": "value",
	})

	assert.NoError(t, err)
	assert.True(t, tokenUpdated)

	sessionInfo, err := GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, "value", sessionInfo.CustomClaimsInAccessTokenPayload["key"])

	tokenUpdated, err = MergeIntoAccessTokenPayload(sessionHandles[0], map[string]interface{}{
		"key": "value2",
	})

	assert.NoError(t, err)
	assert.True(t, tokenUpdated)

	sessionInfo1, err := GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, "value2", sessionInfo1.CustomClaimsInAccessTokenPayload["key"])

	tokenUpdated, err = MergeIntoAccessTokenPayload("random", map[string]interface{}{
		"key": "value3",
	})

	assert.NoError(t, err)
	assert.False(t, tokenUpdated)
}

func TestWhenAntiCsrfIsDisabledFromSTcoreNotHavingThatInInputToVerifySessionIsFine(t *testing.T) {
	customAntiCsrfVal := "NONE"
	True := true
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				AntiCsrf:     &customAntiCsrfVal,
				CookieSecure: &True,
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "supertokens", map[string]interface{}{}, map[string]interface{}{})
	})

	customValForAntiCsrfCheck := false
	customSessionRequiredValue := true
	mux.HandleFunc("/getSessionWithAntiCsrfFalse", VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue,
		AntiCsrfCheck:   &customValForAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sess, err := GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue,
			AntiCsrfCheck:   &customValForAntiCsrfCheck,
		})
		if err != nil {
			t.Error(err.Error())
		}
		assert.NotNil(t, sess)
	}))

	customValForAntiCsrfCheck1 := true
	customSessionRequiredValue1 := true
	mux.HandleFunc("/getSessionWithAntiCsrfTrue", VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue1,
		AntiCsrfCheck:   &customValForAntiCsrfCheck1,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sess, err := GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue1,
			AntiCsrfCheck:   &customValForAntiCsrfCheck1,
		})
		if err != nil {
			t.Error(err.Error())
		}
		assert.NotNil(t, sess)
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	cookieDataWithoutAntiCsrf := unittesting.ExtractInfoFromResponseWhenAntiCSRFisNone(res)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/getSessionWithAntiCsrfFalse", nil)
	assert.NoError(t, err)
	req1.Header.Add("Cookie", "sAccessToken="+cookieDataWithoutAntiCsrf["sAccessToken"])
	res1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, 200, res1.StatusCode)

	req2, err := http.NewRequest(http.MethodGet, testServer.URL+"/getSessionWithAntiCsrfTrue", nil)
	assert.NoError(t, err)
	req2.Header.Add("Cookie", "sAccessToken="+cookieDataWithoutAntiCsrf["sAccessToken"])
	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, res2.StatusCode)
}

func TestAntiCsrfDisabledAndSameSiteNoneDoesNotThrowAnError(t *testing.T) {
	customAntiCsrfVal := "NONE"
	customCookieSameSiteVal := "none"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				AntiCsrf:       &customAntiCsrfVal,
				CookieSameSite: &customCookieSameSiteVal,
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

	assert.NoError(t, err)
}

func TestAntiCsrfDisabledAndSameSiteLaxDoesNotThrowAnError(t *testing.T) {
	customAntiCsrfVal := "NONE"
	customCookieSameSiteVal := "lax"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				AntiCsrf:       &customAntiCsrfVal,
				CookieSameSite: &customCookieSameSiteVal,
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
	assert.NoError(t, err)
}

func TestAntiCsrfDisabledAndSameSiteStrictDoesNotThrowAnError(t *testing.T) {
	customAntiCsrfVal := "NONE"
	customCookieSameSiteVal := "strict"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				AntiCsrf:       &customAntiCsrfVal,
				CookieSameSite: &customCookieSameSiteVal,
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
	assert.NoError(t, err)
}

func TestCustomUserIdIsReturnedCorrectly(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "ronit", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	sessionHandlers, err := GetAllSessionHandlesForUser("ronit")

	if err != nil {
		t.Error(err.Error())
	}

	sessionInfo, err := GetSessionInformation(sessionHandlers[0])
	assert.NoError(t, err)

	assert.Equal(t, "ronit", sessionInfo.UserId)
}

func TestRevokedSessionThrowsErrorWhenCallingGetSessionBySessionHandle(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "ronit", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	sessionHandlers, err := GetAllSessionHandlesForUser("ronit")

	if err != nil {
		t.Error(err.Error())
	}

	sessionInfo, err := GetSessionInformation(sessionHandlers[0])
	assert.NoError(t, err)
	assert.Equal(t, "ronit", sessionInfo.UserId)
	_, err = RevokeMultipleSessions(sessionHandlers)
	assert.NoError(t, err)
	_, err = RevokeAllSessionsForUser("ronit")
	assert.NoError(t, err)
	sessionInformation, err := GetSessionInformation(sessionHandlers[0])
	assert.Nil(t, sessionInformation)
	assert.NoError(t, err)
}

func TestSignoutWorksAfterSessionDeletedOnBackend(t *testing.T) {
	sessionHandle := ""
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		sess, _ := CreateNewSession(r, rw, "rope", map[string]interface{}{}, map[string]interface{}{})
		sessionHandle = sess.GetHandle()
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	RevokeSession(sessionHandle)

	resp1, err := unittesting.SignoutRequest(testServer.URL, cookieData["sAccessToken"], cookieData["antiCsrf"])
	cookieData = unittesting.ExtractInfoFromResponse(resp1)

	assert.Equal(t, cookieData["accessTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData["refreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData["accessToken"], "")
	assert.Equal(t, cookieData["refreshToken"], "")
}

func TestSessionContainerOverride(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						oGetSessionInformation := *originalImplementation.GetSessionInformation
						nGetSessionInformation := func(sessionHandle string, userContext supertokens.UserContext) (*sessmodels.SessionInformation, error) {
							info, err := oGetSessionInformation(sessionHandle, userContext)
							if err != nil {
								return nil, err
							}
							info.SessionDataInDatabase = map[string]interface{}{
								"test": 1,
							}
							return info, nil
						}
						*originalImplementation.GetSessionInformation = nGetSessionInformation
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
	res := MockResponseWriter{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	session, err := CreateNewSession(req, res, "testId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	data, err := session.GetSessionDataInDatabase()
	assert.NoError(t, err)

	assert.Equal(t, 1, data["test"])
}

func TestGetSessionReturnsNilForJWTWithoutSessionClaims(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
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
	False := false
	mux.HandleFunc("/getSession", func(rw http.ResponseWriter, r *http.Request) {
		response, err := GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &False,
		})

		assert.NoError(t, err)
		assert.Nil(t, response)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	validity := uint64(60)
	response, err := CreateJWT(map[string]interface{}{}, &validity, nil)

	if err != nil {
		t.Error(err.Error())
	}

	jwt := response.OK.Jwt

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/getSession", nil)
	req.Header.Add("Authorization", "Bearer "+jwt)
	assert.NoError(t, err)
	_, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetSessionReturnsNilForRequestWithNoSessionWithCheckDatabaseTrueAndSessionRequiredFalse(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
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
	False := false
	True := true
	mux.HandleFunc("/getSession", func(rw http.ResponseWriter, r *http.Request) {
		response, err := GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &False,
			CheckDatabase:   &True,
		})

		assert.NoError(t, err)
		assert.Nil(t, response)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/getSession", nil)
	assert.NoError(t, err)
	_, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestThatJWKSIsFetchedAsExpected(t *testing.T) {
	originalRefreshlimit := JWKRefreshRateLimit
	originalCacheAge := JWKCacheMaxAgeInMs

	JWKRefreshRateLimit = 100
	JWKCacheMaxAgeInMs = 2000

	lastLineBeforeTest := unittesting.GetInfoLogData(t, "").LastLine

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
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

	logInfoAfter := unittesting.GetInfoLogData(t, lastLineBeforeTest)
	var wellKnownCallLogs []string

	for _, line := range logInfoAfter.Output {
		if strings.Contains(line, "API called: /.well-known/jwks.json. Method: GET") {
			wellKnownCallLogs = append(wellKnownCallLogs, line)
		}
	}

	assert.Equal(t, len(wellKnownCallLogs), 1)

	session, err := CreateNewSessionWithoutRequestResponse("rope", map[string]interface{}{}, map[string]interface{}{}, nil)

	if err != nil {
		t.Error(err.Error())
	}

	tokens := session.GetAllSessionTokensDangerously()
	_, err = GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, &sessmodels.VerifySessionOptions{})

	if err != nil {
		t.Error(err.Error())
	}

	time.Sleep(time.Duration(JWKCacheMaxAgeInMs) * time.Millisecond)

	logInfoAfterWaiting := unittesting.GetInfoLogData(t, lastLineBeforeTest)
	wellKnownCallLogs = []string{}

	for _, line := range logInfoAfterWaiting.Output {
		if strings.Contains(line, "API called: /.well-known/jwks.json. Method: GET") {
			wellKnownCallLogs = append(wellKnownCallLogs, line)
		}
	}

	assert.Equal(t, len(wellKnownCallLogs), 2)

	JWKRefreshRateLimit = originalRefreshlimit
	JWKCacheMaxAgeInMs = originalCacheAge
}

func TestThatJWKSResultIsRefreshedProperly(t *testing.T) {
	originalRefreshlimit := JWKRefreshRateLimit
	originalCacheAge := JWKCacheMaxAgeInMs

	JWKRefreshRateLimit = 100
	JWKCacheMaxAgeInMs = 2000

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}
	BeforeEach()
	// 0.004 = 2 seconds roughly
	unittesting.SetKeyValueInConfig("access_token_dynamic_signing_key_update_interval", "0.0004")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	jwksAfterStart := jwksResults
	beforeKids := jwksAfterStart[0].JWKS.KIDs()

	time.Sleep(3 * time.Second)

	afterKids := jwksAfterStart[0].JWKS.KIDs()
	var newKeys []string

	for _, key := range afterKids {
		if !supertokens.DoesSliceContainString(key, beforeKids) {
			newKeys = append(newKeys, key)
		}
	}

	assert.True(t, len(newKeys) != 0)
	JWKRefreshRateLimit = originalRefreshlimit
	JWKCacheMaxAgeInMs = originalCacheAge
}

type MockResponseWriter struct{}

func (mw MockResponseWriter) Header() http.Header {
	return http.Header{}
}
func (mw MockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}
func (mw MockResponseWriter) WriteHeader(statusCode int) {
}
