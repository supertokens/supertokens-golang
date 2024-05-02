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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestCookieBasedAuth(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	cfgVal := func(tokenTransferMethod sessmodels.TokenTransferMethod, olderCookieDomain *string) supertokens.TypeInput {
		customAntiCsrfVal := "VIA_TOKEN"
		return supertokens.TypeInput{
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
					OlderCookieDomain: olderCookieDomain,
					AntiCsrf:          &customAntiCsrfVal,
					GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
						return tokenTransferMethod
					},
				}),
			},
		}
	}

	err := supertokens.Init(cfgVal(sessmodels.CookieTransferMethod, nil))
	assert.NoError(t, err)

	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "public", "rope", map[string]interface{}{}, map[string]interface{}{})
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

	assert.Equal(t, []string{"front-token, anti-csrf"}, res.Header["Access-Control-Expose-Headers"])
	assert.Equal(t, "", cookieData["refreshTokenDomain"])
	assert.Equal(t, "", cookieData["accessTokenDomain"])
	assert.NotNil(t, cookieData["sAccessToken"])
	assert.NotNil(t, cookieData["sRefreshToken"])
	assert.NotNil(t, cookieData["antiCsrf"])
	assert.NotNil(t, cookieData["accessTokenExpiry"])
	assert.NotNil(t, cookieData["refreshTokenExpiry"])

	t.Run("verifySession returns 401 if multiple tokens are passed in the request", func(t *testing.T) {
		req, err = http.NewRequest(http.MethodGet, testServer.URL+"/verifySession", nil)
		assert.NoError(t, err)
		req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
		req.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])
		req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
		req.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])
		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		content, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"message":"try refresh token"}`, string(content))
	})

	t.Run("refresh endpoint throws a 500 if multiple tokens are passed and olderCookieDomain is undefined", func(t *testing.T) {
		req, err = http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
		assert.NoError(t, err)
		req.Header.Add("Cookie", "sAccessToken=accessToken1")
		req.Header.Add("Cookie", "sAccessToken=accessToken2")
		req.Header.Add("Cookie", "sRefreshToken=refreshToken1")
		req.Header.Add("Cookie", "sRefreshToken=refreshToken2")
		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		cookieData = unittesting.ExtractInfoFromResponse(res)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		content, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, "The request contains multiple session cookies. This may happen if you've changed the 'cookieDomain' value in your configuration. To clear tokens from the previous domain, set 'olderCookieDomain' in your config.\n", string(content))
	})

	t.Run("all session tokens are cleared if refresh token api is called without the refresh token but with the access token", func(t *testing.T) {
		req, err = http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
		assert.NoError(t, err)
		req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
		req.Header.Add("anti-csrf", cookieData["antiCsrf"])
		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		cookieData = unittesting.ExtractInfoFromResponse(res)

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.Empty(t, cookieData["sAccessToken"])
		assert.Equal(t, cookieData["accessTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
		assert.Empty(t, cookieData["sRefreshToken"])
		assert.Equal(t, cookieData["refreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	})

	resetAll()
	olderCookieName := ".example.com"
	err = supertokens.Init(cfgVal(sessmodels.CookieTransferMethod, &olderCookieName))
	assert.NoError(t, err)

	mux = http.NewServeMux()
	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "public", "rope", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer = httptest.NewServer(supertokens.Middleware(mux))

	t.Run("access and refresh token for olderCookieDomain is cleared if multiple tokens are passed to the refresh endpoint", func(t *testing.T) {
		req, err = http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
		assert.NoError(t, err)
		req.Header.Add("Cookie", "sAccessToken=accessToken1")
		req.Header.Add("Cookie", "sAccessToken=accessToken2")
		req.Header.Add("Cookie", "sRefreshToken=refreshToken1")
		req.Header.Add("Cookie", "sRefreshToken=refreshToken2")
		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		cookieData = unittesting.ExtractInfoFromResponse(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Empty(t, cookieData["sAccessToken"])
		assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData["accessTokenExpiry"])
		assert.Equal(t, "example.com", cookieData["accessTokenDomain"]) // TODO: node sdk returns .example.com
		assert.Empty(t, cookieData["sRefreshToken"])
		assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData["refreshTokenExpiry"])
		assert.Equal(t, "example.com", cookieData["refreshTokenDomain"]) // TODO: node sdk returns .example.com
	})
}

func TestHeaderBasedAuthAndMultipleTokensInCookies(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	cfgVal := func(tokenTransferMethod sessmodels.TokenTransferMethod, olderCookieDomain *string) supertokens.TypeInput {
		customAntiCsrfVal := "VIA_TOKEN"
		return supertokens.TypeInput{
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
					OlderCookieDomain: olderCookieDomain,
					AntiCsrf:          &customAntiCsrfVal,
					GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
						return tokenTransferMethod
					},
				}),
			},
		}
	}

	err := supertokens.Init(cfgVal(sessmodels.HeaderTransferMethod, nil))
	assert.NoError(t, err)

	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "public", "testuserid", map[string]interface{}{}, map[string]interface{}{})
	})

	mux.HandleFunc("/verifySession", func(writer http.ResponseWriter, request *http.Request) {
		sessionResponse, _ := GetSession(request, writer, nil)
		userID := sessionResponse.GetUserID()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(userID))
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

	t.Run("verifySession returns 200 in header based auth even if multiple tokens are present in the cookie", func(t *testing.T) {
		req, err = http.NewRequest(http.MethodGet, testServer.URL+"/verifySession", nil)
		assert.NoError(t, err)
		req.Header.Add("Authorization", "Bearer "+cookieData["accessTokenFromHeader"])
		req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
		req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		content, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, `testuserid`, string(content))
	})

	t.Run("refresh endpoint refreshes the token in header based auth even if multiple tokens are present in the cookie", func(t *testing.T) {
		req, err = http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
		assert.NoError(t, err)
		req.Header.Add("Authorization", "Bearer "+cookieData["refreshTokenFromHeader"])
		req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
		req.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])
		req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
		req.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])

		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		cookieData := unittesting.ExtractInfoFromResponse(res)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.NotEmpty(t, cookieData["accessTokenFromHeader"])
		assert.NotEmpty(t, cookieData["refreshTokenFromHeader"])
	})
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
		CreateNewSession(r, rw, "public", "user", map[string]interface{}{}, map[string]interface{}{})
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
		CreateNewSession(r, rw, "public", "userId", map[string]interface{}{}, map[string]interface{}{})
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
		CreateNewSession(r, rw, "public", "someId", map[string]interface{}{}, map[string]interface{}{})
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
		CreateNewSession(r, rw, "public", "someUniqueID", map[string]interface{}{}, map[string]interface{}{})
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

	_, err = RevokeAllSessionsForUser("someUniqueID", nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	sessionHandlesAfterRevoke, err := GetAllSessionHandlesForUser("someUniqueID", nil, nil)
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

	sessionHandlesBeforeRevoke1, err := GetAllSessionHandlesForUser("someUniqueID", nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 2, len(sessionHandlesBeforeRevoke1))

	revokedSessions, err := RevokeAllSessionsForUser("someUniqueID", nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 2, len(revokedSessions))

	sessionHandlesAfterRevoke1, err := GetAllSessionHandlesForUser("someUniqueID", nil, nil)
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
		CreateNewSession(r, rw, "public", "rp", map[string]interface{}{}, map[string]interface{}{})
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

	sessionHandles, err := GetAllSessionHandlesForUser("rp", nil, nil)

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
		CreateNewSession(r, rw, "public", "uniqueId", map[string]interface{}{}, nil)
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

	sessionHandles, err := GetAllSessionHandlesForUser("uniqueId", nil, nil)

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
		CreateNewSession(r, rw, "public", "uniqueId", map[string]interface{}{}, map[string]interface{}{})
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

	sessionHandles, err := GetAllSessionHandlesForUser("uniqueId", nil, nil)

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
		CreateNewSession(r, rw, "public", "supertokens", map[string]interface{}{}, map[string]interface{}{})
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
		CreateNewSession(r, rw, "public", "ronit", map[string]interface{}{}, map[string]interface{}{})
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

	sessionHandlers, err := GetAllSessionHandlesForUser("ronit", nil, nil)

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
		CreateNewSession(r, rw, "public", "ronit", map[string]interface{}{}, map[string]interface{}{})
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

	sessionHandlers, err := GetAllSessionHandlesForUser("ronit", nil, nil)

	if err != nil {
		t.Error(err.Error())
	}

	sessionInfo, err := GetSessionInformation(sessionHandlers[0])
	assert.NoError(t, err)
	assert.Equal(t, "ronit", sessionInfo.UserId)
	_, err = RevokeMultipleSessions(sessionHandlers)
	assert.NoError(t, err)
	_, err = RevokeAllSessionsForUser("ronit", nil, nil)
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
		sess, _ := CreateNewSession(r, rw, "public", "rope", map[string]interface{}{}, map[string]interface{}{})
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
	session, err := CreateNewSession(req, res, "public", "testId", map[string]interface{}{}, map[string]interface{}{})
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

/*
*
This test verifies that the SDK calls the well known API properly in the normal flow

- Initialise the SDK and verify that the well known API was not called
- Create and verify a session
- Verify that the well known API was called to fetch the keys
*/
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

	assert.Equal(t, len(wellKnownCallLogs), 0)

	session, err := CreateNewSessionWithoutRequestResponse("public", "rope", map[string]interface{}{}, map[string]interface{}{}, nil)

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

	assert.Equal(t, len(wellKnownCallLogs), 1)

	JWKRefreshRateLimit = originalRefreshlimit
	JWKCacheMaxAgeInMs = originalCacheAge
}

/*
*
This test verifies that the cache used to store the pointer to the JWKS result is updated properly when the
cache expired and the keys need to be refetched.

- Init
- Call getJWKS to get the keys
- Wait for access token signing key to change
- Fetch the keys again
- Verify that the KIDs inside the pointer have changed
*/
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

	jwksBefore, err := getJWKS()

	if err != nil {
		t.Error(err.Error())
	}

	beforeKids := jwksBefore.KIDs()

	time.Sleep(3 * time.Second)

	jwksAfter, err := getJWKS()

	if err != nil {
		t.Error(err.Error())
	}

	afterKids := jwksAfter.KIDs()
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

/*
*
This test verifies that the SDK tried to re-fetch the keys from the core if the KID for the access token does not
exist in the keyfunc library's cache

- Init and verify no calls have been made
- Create and verify a session
- Verify that a call to the well known API was made
- Wait for access_token_dynamic_signing_key_update_interval so that the core uses a new key
- Create and verify another session
- Verify that the call to the well known API was made
- Create and verify another session
- Verify that no call is made
*/
func TestThatJWKSAreRefreshedIfKIDIsUnkown(t *testing.T) {
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
	unittesting.SetKeyValueInConfig("access_token_dynamic_signing_key_update_interval", "0.0014")
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

	assert.Equal(t, len(wellKnownCallLogs), 0)

	session, err := CreateNewSessionWithoutRequestResponse("public", "rope", map[string]interface{}{}, map[string]interface{}{}, nil)

	if err != nil {
		t.Error(err.Error())
	}

	tokens := session.GetAllSessionTokensDangerously()
	session, err = GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, &sessmodels.VerifySessionOptions{})

	if err != nil {
		t.Error(err.Error())
	}

	logInfoAfter = unittesting.GetInfoLogData(t, lastLineBeforeTest)
	wellKnownCallLogs = []string{}

	for _, line := range logInfoAfter.Output {
		if strings.Contains(line, "API called: /.well-known/jwks.json. Method: GET") {
			wellKnownCallLogs = append(wellKnownCallLogs, line)
		}
	}

	assert.Equal(t, len(wellKnownCallLogs), 1)

	time.Sleep(10 * time.Second)

	session, err = CreateNewSessionWithoutRequestResponse("public", "rope", map[string]interface{}{}, map[string]interface{}{}, nil)
	if err != nil {
		t.Error(err.Error())
	}

	tokens = session.GetAllSessionTokensDangerously()
	session, err = GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, &sessmodels.VerifySessionOptions{})

	if err != nil {
		t.Error(err.Error())
	}

	logInfoAfter = unittesting.GetInfoLogData(t, lastLineBeforeTest)
	wellKnownCallLogs = []string{}

	for _, line := range logInfoAfter.Output {
		if strings.Contains(line, "API called: /.well-known/jwks.json. Method: GET") {
			wellKnownCallLogs = append(wellKnownCallLogs, line)
		}
	}

	assert.Equal(t, len(wellKnownCallLogs), 2)

	tokens = session.GetAllSessionTokensDangerously()
	_, err = GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, &sessmodels.VerifySessionOptions{})

	logInfoAfter = unittesting.GetInfoLogData(t, lastLineBeforeTest)
	wellKnownCallLogs = []string{}

	for _, line := range logInfoAfter.Output {
		if strings.Contains(line, "API called: /.well-known/jwks.json. Method: GET") {
			wellKnownCallLogs = append(wellKnownCallLogs, line)
		}
	}

	assert.Equal(t, len(wellKnownCallLogs), 2)
}

/*
*
This test makes sure that initialising SuperTokens and Session with an invalid connection uri does not
result in an error during startup
*/
func TestThatInvalidConnectionUriDoesNotThrowDuringInitForJWKS(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "https://try.supertokens.io:8080",
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
	unittesting.SetKeyValueInConfig("access_token_dynamic_signing_key_update_interval", "0.0014")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
}

/*
*
This test verifies the behaviour of the JWKS cache maintained by the SDK

- Init
- Make sure the cache is empty
- Create and verify a session
- Make sure the cache has one entry now
- Wait for cache to expire
- Verify the session again
- Verify that an entry from the cache was deleted (because the SDK should try to re-fetch)
- Verify that the cache has an entry
*/
func TestJWKSCacheLogic(t *testing.T) {
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
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Nil(t, jwksCache)

	session, err := CreateNewSessionWithoutRequestResponse("public", "rope", map[string]interface{}{}, map[string]interface{}{}, nil)

	if err != nil {
		t.Error(err.Error())
	}

	assert.Nil(t, jwksCache)

	tokens := session.GetAllSessionTokensDangerously()
	session, err = GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, &sessmodels.VerifySessionOptions{})

	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, jwksCache)

	time.Sleep(3 * time.Second)

	session, err = GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, &sessmodels.VerifySessionOptions{})
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, jwksCache)

	JWKRefreshRateLimit = originalRefreshlimit
	JWKCacheMaxAgeInMs = originalCacheAge
}

/*
*
This test ensures that calling get combines JWKS results in an error if the connection uri is invalid. Note that
in this test we specifically expect a timeout but that does not mean that this is the only error the function can
throw
*/
func TestThatCombinedJWKSThrowsForInvalidConnectionUri(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "https://try.supertokens.io:8080",
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
	unittesting.SetKeyValueInConfig("access_token_dynamic_signing_key_update_interval", "0.0014")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	combinedJwks, err := GetCombinedJWKS()

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "timeout"))
	assert.Nil(t, combinedJwks)
}

/*
*
This test makes sure that when multiple core urls are provided, the get combined JWKS function does not throw an
error as long as one of the provided urls return a valid response

- Init with multiple core urls
- Call get combines jwks
- verify that there is a response and that there are no errors
*/
func TestThatGetCombinedJWKSDoesNotThrowIfAtleastOneCoreURLIsValid(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080;example.com:8080;localhost:90",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
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

	combinedJwks, err := GetCombinedJWKS()

	assert.Nil(t, err)
	assert.NotNil(t, combinedJwks)
}

/*
*
This test ensures that the SDK's caching logic for fetching JWKs works fine

- Init
- Create and verify a session
- Verify that the SDK did not use any cached values
- Verify the session again
- Verify that this time the SDK did return a cached value
- Wait for cache to expire
- Verify the session again
- This time the SDK should try to re-fetch and not return a cached value
*/
func TestThatJWKSReturnsFromCacheCorrectly(t *testing.T) {
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
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	session, err := CreateNewSessionWithoutRequestResponse("public", "rope", map[string]interface{}{}, map[string]interface{}{}, nil)

	if err != nil {
		t.Error(err.Error())
	}

	assert.Nil(t, jwksCache)

	tokens := session.GetAllSessionTokensDangerously()
	session, err = GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, &sessmodels.VerifySessionOptions{})

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, <-returnedFromCache, false)

	session, err = GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, &sessmodels.VerifySessionOptions{})

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, <-returnedFromCache, true)

	time.Sleep(3 * time.Second)

	session, err = GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, &sessmodels.VerifySessionOptions{})

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, <-returnedFromCache, false)

	JWKRefreshRateLimit = originalRefreshlimit
	JWKCacheMaxAgeInMs = originalCacheAge
}

/*
*
This test makes sure that the SDK tries to fetch for all core URLS if needed.
This test uses multiple hosts with only the last one being valid to make sure all URLs are used

- init with multiple core urls where only the last one is valid
- Call get combined jwks
- Make sure that the SDK tried fetching for all URLs (since only the last one would return a valid keyset)
*/
func TestThatTheSDKTriesFetchingJWKSForAllCoreHosts(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "example.com;localhost:90;http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
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

	assert.Equal(t, len(urlsAttemptedForJWKSFetch), 0)

	_, err = GetCombinedJWKS()

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, len(urlsAttemptedForJWKSFetch), 3)
}

/*
*
This test makes sure that the SDK stop fetching JWKS from multiple cores as soon as it gets a valid response

- init with multiple cores with the second one being valid (1st and 3rd invalid)
- call get combines jwks
- Verify that two urls were used to fetch and that the 3rd one was never used
*/
func TestThatTheSDKFetchesJWKSFromAllCoreHostsUntilAValidResponse(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "example.com;http://localhost:8080;localhost:90",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
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

	assert.Equal(t, len(urlsAttemptedForJWKSFetch), 0)

	_, err = GetCombinedJWKS()

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, len(urlsAttemptedForJWKSFetch), 2)
	assert.True(t, strings.Contains(urlsAttemptedForJWKSFetch[0], "example.com"))
	assert.True(t, strings.Contains(urlsAttemptedForJWKSFetch[1], "http://localhost:8080"))
}

func TestSessionVerificationOfJWTBasedOnSessionPayload(t *testing.T) {
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
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	session, err := CreateNewSessionWithoutRequestResponse("public", "testing", map[string]interface{}{}, map[string]interface{}{}, nil)
	if err != nil {
		t.Error(err.Error())
	}

	payload := session.GetAccessTokenPayload()
	delete(payload, "iat")
	delete(payload, "exp")

	currentTimeInSeconds := time.Now()
	jwtExpiry := uint64((currentTimeInSeconds.Add(10 * time.Second)).Unix())
	False := false
	jwt, err := CreateJWT(payload, &jwtExpiry, &False)
	if err != nil {
		t.Error(err.Error())
	}

	session, err = GetSessionWithoutRequestResponse(jwt.OK.Jwt, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, session.GetUserID(), "testing")
}

func TestSessionVerificationOfJWTBasedOnSessionPayloadWithCheckDatabase(t *testing.T) {
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
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("session")
	if err != nil {
		t.Fail()
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Fail()
	}

	// Only run test for cdi > 2.21 (not greater than equal to)
	if supertokens.MaxVersion(cdiVersion, "2.21") == "2.21" {
		t.Skip()
	}

	session, err := CreateNewSessionWithoutRequestResponse("public", "testing", map[string]interface{}{}, map[string]interface{}{}, nil)
	if err != nil {
		t.Error(err.Error())
	}

	payload := session.GetAccessTokenPayload()
	delete(payload, "iat")
	delete(payload, "exp")
	payload["tId"] = "public"
	payload["rsub"] = session.GetUserID()

	currentTimeInSeconds := time.Now()
	jwtExpiry := uint64((currentTimeInSeconds.Add(10 * time.Second)).Unix())
	False := false
	jwt, err := CreateJWT(payload, &jwtExpiry, &False)
	if err != nil {
		t.Error(err.Error())
	}

	True := true
	session, err = GetSessionWithoutRequestResponse(jwt.OK.Jwt, nil, &sessmodels.VerifySessionOptions{
		CheckDatabase: &True,
	})
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, session.GetUserID(), "testing")
}

func TestThatLockingForJWKSCacheWorksFine(t *testing.T) {
	originalRefreshlimit := JWKRefreshRateLimit
	originalCacheAge := JWKCacheMaxAgeInMs

	JWKRefreshRateLimit = 100
	JWKCacheMaxAgeInMs = 2000

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
		},
	}
	BeforeEach()
	unittesting.SetKeyValueInConfig("access_token_dynamic_signing_key_update_interval", "0.0014")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	differentKeyFoundCount := 0
	notReturnFromCacheCount := 0
	keys := []string{}
	shouldStop := false

	jwks, err := GetCombinedJWKS()
	if err != nil {
		t.Error(err.Error())
	}
	<-returnedFromCache // this value must be ignored

	for _, k := range jwks.KIDs() {
		keys = append(keys, k)
	}

	go func() {
		time.Sleep(11 * time.Second)
		shouldStop = true
	}()

	threadCount := 10
	var wg sync.WaitGroup
	wg.Add(threadCount)

	for i := 0; i < threadCount; i++ {
		go jwksLockTestRoutine(t, &shouldStop, i, &wg, func(_keys []string) {
			if <-returnedFromCache == false {
				notReturnFromCacheCount++
			}

			newKeys := []string{}

			for _, _k2 := range _keys {
				if !supertokens.DoesSliceContainString(_k2, keys) {
					newKeys = append(newKeys, _k2)
				}
			}

			if len(newKeys) != 0 {
				differentKeyFoundCount++
				keys = _keys
			}
		})
	}

	wg.Wait()

	// We test for both
	// - The keys changing
	// - The number of times the result is not returned from cache
	//
	// Because even if the keys change only twice it could still mean that the SDK's cache locking
	// does not work correctly and that it tried to query the core more times than it should have
	//
	// Checking for both the key change count and the cache miss count verifies the locking behaviour properly
	//
	// With the signing key interval as 5 seconds, and the test making requests for 11 seconds
	// You expect the keys to change twice
	assert.Equal(t, differentKeyFoundCount, 2)
	// With cache lifetime as 2 seconds, you expect the cache to miss 5 times
	assert.Equal(t, notReturnFromCacheCount, 5)

	JWKRefreshRateLimit = originalRefreshlimit
	JWKCacheMaxAgeInMs = originalCacheAge
}

func TestThatGetSessionThrowsWIthDynamicKeysIfSessionWasCreatedWithStaticKeys(t *testing.T) {
	False := false
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
			Init(&sessmodels.TypeInput{
				UseDynamicAccessTokenSigningKey: &False,
			}),
		},
	}

	BeforeEach()
	unittesting.SetKeyValueInConfig("access_token_dynamic_signing_key_update_interval", "0.0014")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	session, err := CreateNewSessionWithoutRequestResponse("public", "testing-user", map[string]interface{}{}, map[string]interface{}{}, nil)

	resetAll()
	True := true
	configValue = supertokens.TypeInput{
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
				UseDynamicAccessTokenSigningKey: &True,
			}),
		},
	}
	err = supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	sessionTokens := session.GetAllSessionTokensDangerously()
	session, err = GetSessionWithoutRequestResponse(sessionTokens.AccessToken, sessionTokens.AntiCsrfToken, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "The access token doesn't match the useDynamicAccessTokenSigningKey setting")
}

func TestThatRevokedAccessTokenThrowsUnauthorisedErrorWhenRegenerateTokenIsCalled(t *testing.T) {
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
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	sessionContainer, err := CreateNewSessionWithoutRequestResponse("public", "testing-user", map[string]interface{}{}, map[string]interface{}{}, nil)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = RevokeSession(sessionContainer.GetHandle())
	if err != nil {
		t.Error(err.Error())
	}

	err = sessionContainer.MergeIntoAccessTokenPayload(map[string]interface{}{"key": "value"})
	assert.NotNil(t, err)
	_, ok := err.(errors.UnauthorizedError)
	assert.True(t, ok)
}

func jwksLockTestRoutine(t *testing.T, shouldStop *bool, index int, group *sync.WaitGroup, doPost func([]string)) {
	jwks, err := GetCombinedJWKS()
	if err != nil {
		t.Error(err.Error())
	}

	doPost(jwks.KIDs())
	time.Sleep(100 * time.Millisecond)
	if *shouldStop == false {
		jwksLockTestRoutine(t, shouldStop, index, group, doPost)
	} else {
		group.Done()
	}
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
