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

package unittesting

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := ExtractInfoFromResponse(res)
	assert.Equal(t, []string{"front-token, id-refresh-token, anti-csrf"}, res.Header["Access-Control-Expose-Headers"])
	assert.Equal(t, "", cookieData["refreshTokenDomain"])
	assert.Equal(t, "", cookieData["idRefreshTokenDomain"])
	assert.Equal(t, "", cookieData["accessTokenDomain"])
	assert.NotNil(t, cookieData["sAccessToken"])
	assert.NotNil(t, cookieData["sRefreshToken"])
	assert.NotNil(t, cookieData["sIdRefreshToken"])
	assert.NotNil(t, cookieData["antiCsrf"])
	assert.NotNil(t, cookieData["accessTokenExpiry"])
	assert.NotNil(t, cookieData["refreshTokenExpiry"])
	assert.NotNil(t, cookieData["idRefreshTokenExpiry"])
	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "", map[string]interface{}{}, map[string]interface{}{})
	})

	customValForAntiCsrfCheck := true
	customSessionRequiredValue := true
	mux.HandleFunc("/verifySession", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue,
		AntiCsrfCheck:   &customValForAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		session.GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue,
			AntiCsrfCheck:   &customValForAntiCsrfCheck,
		})
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := ExtractInfoFromResponse(res)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req2.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req2.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res2, err := http.DefaultClient.Do(req2)
	cookieData2 := ExtractInfoFromResponse(res2)
	assert.NoError(t, err)

	reqV, err := http.NewRequest(http.MethodGet, testServer.URL+"/verifySession", nil)
	assert.NoError(t, err)
	reqV.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"]+";"+"sIdRefreshToken="+cookieData2["sIdRefreshToken"])
	reqV.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	resv, err := http.DefaultClient.Do(reqV)
	assert.NoError(t, err)
	assert.Equal(t, resv.StatusCode, 200)

	req3, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req3.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
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
	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "", map[string]interface{}{}, map[string]interface{}{})
	})
	customValForAntiCsrfCheck := true
	customSessionRequiredValue := true
	mux.HandleFunc("/verifySession", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue,
		AntiCsrfCheck:   &customValForAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		session.GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue,
			AntiCsrfCheck:   &customValForAntiCsrfCheck,
		})
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := ExtractInfoFromResponse(res)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req2.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req2.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res2, err := http.DefaultClient.Do(req2)
	cookieData2 := ExtractInfoFromResponse(res2)
	assert.NoError(t, err)

	reqV, err := http.NewRequest(http.MethodGet, testServer.URL+"/verifySession", nil)
	assert.NoError(t, err)
	reqV.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"]+";"+"sIdRefreshToken="+cookieData2["sIdRefreshToken"])
	reqV.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	resv, err := http.DefaultClient.Do(reqV)
	assert.NoError(t, err)
	assert.Equal(t, resv.StatusCode, 200)

	req3, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req3.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
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
	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestQuerringToTheCoreWithoutAPIKey(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
			// APIKey:        "shfo3h98308hOIHoei309saiho",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	// SetKeyValueInConfig("api_keys", "shfo3h98308hOIHoei309saiho")
	SetKeyValueInConfig("api_keys", `"shfo3h98308hOIHoei309saiho"`)
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	querrier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	_, err = querrier.GetQuerierAPIVersion()
	if err != nil {
		assert.Equal(t, "SuperTokens core threw an error for a request to path: '/apiversion' with status code: 401 and message: Invalid API key\n", err.Error())
	}

	AfterEach()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "", map[string]interface{}{}, map[string]interface{}{})
	})
	customValForAntiCsrfCheck := true
	customSessionRequiredValue := true
	mux.HandleFunc("/getSession", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue,
		AntiCsrfCheck:   &customValForAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		session.GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue,
			AntiCsrfCheck:   &customValForAntiCsrfCheck,
		})
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := ExtractInfoFromResponse(res)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/getSession", nil)
	assert.NoError(t, err)
	req1.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	// req1.Header.Add("anti-csrf", cookieData["antiCsrf
	res1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, 401, res1.StatusCode)

	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "someUniqueID", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	_, err = session.RevokeAllSessionsForUser("someUniqueID")
	if err != nil {
		t.Error(err.Error())
	}

	sessionHandlesAfterRevoke, err := session.GetAllSessionHandlesForUser("someUniqueID")
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

	sessionHandlesBeforeRevoke1, err := session.GetAllSessionHandlesForUser("someUniqueID")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 2, len(sessionHandlesBeforeRevoke1))

	revokedSessions, err := session.RevokeAllSessionsForUser("someUniqueID")
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 2, len(revokedSessions))

	sessionHandlesAfterRevoke1, err := session.GetAllSessionHandlesForUser("someUniqueID")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 0, len(sessionHandlesAfterRevoke1))

	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	sessionHandles, err := session.GetAllSessionHandlesForUser("")

	if err != nil {
		t.Error(err.Error())
	}

	session.UpdateSessionData(sessionHandles[0], map[string]interface{}{
		"name": "John",
	})

	sessionInfo, err := session.GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, "John", sessionInfo.SessionData["name"])

	session.UpdateSessionData(sessionHandles[0], map[string]interface{}{
		"name": "Joel",
	})

	sessionInfo, err = session.GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, "Joel", sessionInfo.SessionData["name"])

	//update session data with wrong session handle

	err = session.UpdateSessionData("random", map[string]interface{}{
		"name": "Ronit",
	})

	assert.Error(t, err)

	assert.Equal(t, "Session does not exist.", err.Error())

	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "uniqueId", map[string]interface{}{}, nil)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	sessionHandles, err := session.GetAllSessionHandlesForUser("uniqueId")

	if err != nil {
		t.Error(err.Error())
	}

	sessionInfo, err := session.GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, map[string]interface{}{}, sessionInfo.SessionData)

	session.UpdateSessionData(sessionHandles[0], map[string]interface{}{
		"name": "John",
	})
	sessionInfo, err = session.GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, "John", sessionInfo.SessionData["name"])

	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	sessionHandles, err := session.GetAllSessionHandlesForUser("uniqueId")

	if err != nil {
		t.Error(err.Error())
	}

	err = session.UpdateAccessTokenPayload(sessionHandles[0], map[string]interface{}{
		"key": "value",
	})

	assert.NoError(t, err)

	sessionInfo, err := session.GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, "value", sessionInfo.AccessTokenPayload["key"])

	err = session.UpdateAccessTokenPayload(sessionHandles[0], map[string]interface{}{
		"key": "value2",
	})

	assert.NoError(t, err)

	sessionInfo1, err := session.GetSessionInformation(sessionHandles[0])

	assert.NoError(t, err)

	assert.Equal(t, "value2", sessionInfo1.AccessTokenPayload["key"])

	err = session.UpdateAccessTokenPayload("random", map[string]interface{}{
		"key": "value3",
	})

	assert.Error(t, err)

	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestWhenAntiCsrfIsDisabledFromSTcoreNotHavingThatInInputToVerifySessionIsFine(t *testing.T) {
	customAntiCsrfVal := "NONE"
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "supertokens", map[string]interface{}{}, map[string]interface{}{})
	})

	customValForAntiCsrfCheck := false
	customSessionRequiredValue := true
	mux.HandleFunc("/getSessionWithAntiCsrfFalse", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue,
		AntiCsrfCheck:   &customValForAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sess, err := session.GetSession(r, rw, &sessmodels.VerifySessionOptions{
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
	mux.HandleFunc("/getSessionWithAntiCsrfTrue", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue1,
		AntiCsrfCheck:   &customValForAntiCsrfCheck1,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sess, err := session.GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue1,
			AntiCsrfCheck:   &customValForAntiCsrfCheck1,
		})
		if err != nil {
			t.Error(err.Error())
		}
		assert.NotNil(t, sess)
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	cookieDataWithoutAntiCsrf := ExtractInfoFromResponseWhenAntiCSRFisNone(res)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/getSessionWithAntiCsrfFalse", nil)
	assert.NoError(t, err)
	req1.Header.Add("Cookie", "sAccessToken="+cookieDataWithoutAntiCsrf["sAccessToken"]+";"+"sIdRefreshToken="+cookieDataWithoutAntiCsrf["sIdRefreshToken"])
	res1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, 200, res1.StatusCode)

	req2, err := http.NewRequest(http.MethodGet, testServer.URL+"/getSessionWithAntiCsrfTrue", nil)
	assert.NoError(t, err)
	req2.Header.Add("Cookie", "sAccessToken="+cookieDataWithoutAntiCsrf["sAccessToken"]+";"+"sIdRefreshToken="+cookieDataWithoutAntiCsrf["sIdRefreshToken"])
	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, res2.StatusCode)

	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf:       &customAntiCsrfVal,
				CookieSameSite: &customCookieSameSiteVal,
			}),
		},
	}

	BeforeEach()
	StartUpST("localhost", "8080")

	err := supertokens.Init(configValue)

	assert.NoError(t, err)

	AfterEach()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf:       &customAntiCsrfVal,
				CookieSameSite: &customCookieSameSiteVal,
			}),
		},
	}

	BeforeEach()
	StartUpST("localhost", "8080")

	err := supertokens.Init(configValue)

	assert.NoError(t, err)

	AfterEach()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf:       &customAntiCsrfVal,
				CookieSameSite: &customCookieSameSiteVal,
			}),
		},
	}

	BeforeEach()
	StartUpST("localhost", "8080")

	err := supertokens.Init(configValue)

	assert.NoError(t, err)

	AfterEach()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "ronit", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	sessionHandlers, err := session.GetAllSessionHandlesForUser("ronit")

	if err != nil {
		t.Error(err.Error())
	}

	sessionInfo, err := session.GetSessionInformation(sessionHandlers[0])
	assert.NoError(t, err)

	assert.Equal(t, "ronit", sessionInfo.UserId)
	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "ronit", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	sessionHandlers, err := session.GetAllSessionHandlesForUser("ronit")

	if err != nil {
		t.Error(err.Error())
	}

	sessionInfo, err := session.GetSessionInformation(sessionHandlers[0])
	assert.NoError(t, err)
	assert.Equal(t, "ronit", sessionInfo.UserId)
	_, err = session.RevokeMultipleSessions(sessionHandlers)
	assert.NoError(t, err)
	_, err = session.RevokeAllSessionsForUser("ronit")
	assert.NoError(t, err)
	_, err = session.GetSessionInformation(sessionHandlers[0])
	assert.Error(t, err)
	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
}
