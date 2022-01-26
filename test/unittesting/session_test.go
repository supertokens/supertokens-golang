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
	"fmt"
	"log"
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
	StartingHelper()
	err := supertokens.Init(configValue)
	if err != nil {
		log.Fatal(err.Error())
	}
	testServer := httptest.NewServer(supertokens.Middleware(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			session.CreateNewSession(rw, "", nil, nil)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
		),
	))
	req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
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
	defer EndingHelper()
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
	StartingHelper()
	err := supertokens.Init(configValue)
	if err != nil {
		log.Fatal(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "", map[string]interface{}{}, map[string]interface{}{})
	})

	mux.HandleFunc("/verifySession", func(rw http.ResponseWriter, r *http.Request) {
		session.GetSession(r, rw, nil)
	})

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
		log.Fatal(err.Error())
	}
	res3.Body.Close()
	assert.Equal(t, "token theft detected", jsonResponse["message"])
	assert.Equal(t, 401, res3.StatusCode)
	assert.NoError(t, err)
	defer EndingHelper()
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
	StartingHelper()
	err := supertokens.Init(configValue)
	if err != nil {
		log.Fatal(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "", map[string]interface{}{}, map[string]interface{}{})
	})

	mux.HandleFunc("/verifySession", func(rw http.ResponseWriter, r *http.Request) {
		session.GetSession(r, rw, nil)
	})

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
		log.Fatal(err.Error())
	}
	res3.Body.Close()
	assert.Equal(t, "token theft detected", jsonResponse["message"])
	assert.Equal(t, 401, res3.StatusCode)
	assert.NoError(t, err)
	defer EndingHelper()
	defer func() {
		testServer.Close()
	}()
}

//!NEED A BIT OF HELP
func TestQuerringToTheCoreWithoutAPIKey(t *testing.T) {
	SetKeyValueInConfig("api_keys", "shfo3h98308hOIHoei309saiho")
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
	StartingHelper()
	err := supertokens.Init(configValue)
	if err != nil {
		log.Fatal(err.Error())
	}
	querrier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		log.Fatal(err.Error())
	}
	apiVersion, err := querrier.GetQuerierAPIVersion()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(apiVersion)
	EndingHelper()
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
	StartingHelper()
	err := supertokens.Init(configValue)
	if err != nil {
		log.Fatal(err.Error())
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "", map[string]interface{}{}, map[string]interface{}{})
	})

	mux.HandleFunc("/getSession", func(rw http.ResponseWriter, r *http.Request) {
		customValForAntiCsrfCheck := true
		session.GetSession(r, rw, &sessmodels.VerifySessionOptions{
			AntiCsrfCheck: &customValForAntiCsrfCheck,
		})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := ExtractInfoFromResponse(res)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/getSession", nil)
	assert.NoError(t, err)
	req1.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	// req1.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, 200, res1.StatusCode)

	fmt.Println(res1)

	defer EndingHelper()
	defer func() {
		testServer.Close()
	}()
}
