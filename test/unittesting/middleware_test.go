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
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestDisablingDefaultAPIActuallyDisablesIt(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				Override: &sessmodels.OverrideStruct{
					APIs: func(originalImplementation sessmodels.APIInterface) sessmodels.APIInterface {
						*originalImplementation.RefreshPOST = nil
						return originalImplementation
					},
				},
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
	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)

	defer AfterEach()
	defer testServer.Close()
}

func TestSessionVerifyMiddleware(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				ErrorHandlers: &sessmodels.ErrorHandlers{
					OnTokenTheftDetected: func(sessionHandle, userID string, req *http.Request, res http.ResponseWriter) error {
						res.WriteHeader(403)
						resp := make(map[string]string)
						resp["message"] = "Token theft detected"
						jsonResp, err := json.Marshal(resp)
						if err != nil {
							t.Errorf("Error happened in JSON marshal. Err: %s", err)
						}
						res.Write(jsonResp)
						return nil
					},
				},
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
		_, err := session.CreateNewSession(rw, "uniqueId", map[string]interface{}{}, map[string]interface{}{})
		if err != nil {
			rw.WriteHeader(500)
		}
		resp := make(map[string]bool)
		resp["message"] = true
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})

	mux.HandleFunc("/user/id", session.VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userId := sessionContainer.GetUserID()
		resp := make(map[string]string)
		resp["userId"] = userId
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

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

	customAntiCsrfCheck := true
	mux.HandleFunc("/user/handleV0", session.VerifySession(&sessmodels.VerifySessionOptions{
		AntiCsrfCheck: &customAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		handle := sessionContainer.GetHandle()
		resp := make(map[string]string)
		resp["handle"] = handle
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	customSessionRequiredVal := false
	mux.HandleFunc("/user/handleOptional", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredVal,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		resp := make(map[string]bool)
		if sessionContainer == nil {
			resp["message"] = false
		} else {
			resp["message"] = true
		}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	mux.HandleFunc("/logout", session.VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		err := sessionContainer.RevokeSession()
		if err != nil {
			rw.WriteHeader(500)
		}
		resp := make(map[string]bool)
		resp["message"] = true
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	//this is never used. Rather the default api is used
	mux.HandleFunc("/auth/session/refresh", session.VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		resp := make(map[string]bool)
		resp["message"] = true
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)

	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	cookieData := ExtractInfoFromResponse(res)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/user/id", nil)
	req1.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req1.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, 200, res1.StatusCode)

	dataInBytes, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]string
	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "uniqueId", result["userId"])
	res1.Body.Close()

	req2, err := http.NewRequest(http.MethodGet, testServer.URL+"/user/handleV0", nil)
	req2.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	assert.NoError(t, err)
	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 401, res2.StatusCode)
	dataInBytes2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result2 map[string]string
	err = json.Unmarshal(dataInBytes2, &result2)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "try refresh token", result2["message"])
	res2.Body.Close()

	req3, err := http.NewRequest(http.MethodGet, testServer.URL+"/user/handleOptional", nil)
	req3.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	assert.NoError(t, err)
	res3, err := http.DefaultClient.Do(req3)
	assert.NoError(t, err)
	assert.Equal(t, 200, res3.StatusCode)
	dataInBytes3, err := ioutil.ReadAll(res3.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result3 map[string]bool
	err = json.Unmarshal(dataInBytes3, &result3)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, true, result3["message"])
	res3.Body.Close()

	req4, err := http.NewRequest(http.MethodGet, testServer.URL+"/user/handleOptional", nil)
	assert.NoError(t, err)
	res4, err := http.DefaultClient.Do(req4)
	assert.NoError(t, err)
	assert.Equal(t, 200, res4.StatusCode)
	dataInBytes4, err := ioutil.ReadAll(res4.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result4 map[string]bool
	err = json.Unmarshal(dataInBytes4, &result4)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, false, result4["message"])
	res4.Body.Close()

	req5, err := http.NewRequest(http.MethodGet, testServer.URL+"/user/handleV0", nil)
	req5.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
	assert.NoError(t, err)
	req5.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res5, err := http.DefaultClient.Do(req5)
	assert.NoError(t, err)
	assert.Equal(t, 401, res5.StatusCode)
	dataInBytes5, err := ioutil.ReadAll(res5.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result5 map[string]string
	err = json.Unmarshal(dataInBytes5, &result5)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "unauthorised", result5["message"])
	res5.Body.Close()

	req6, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	req6.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req6.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res6, err := http.DefaultClient.Do(req6)
	cookieData2 := ExtractInfoFromResponse(res6)
	assert.NoError(t, err)
	assert.Equal(t, 200, res6.StatusCode)

	req7, err := http.NewRequest(http.MethodGet, testServer.URL+"/verifySession", nil)
	req7.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"]+";"+"sIdRefreshToken="+cookieData2["sIdRefreshToken"])
	req7.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	assert.NoError(t, err)
	res7, err := http.DefaultClient.Do(req7)
	assert.NoError(t, err)
	assert.Equal(t, 200, res7.StatusCode)

	req8, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	req8.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req8.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res8, err := http.DefaultClient.Do(req8)
	assert.NoError(t, err)
	assert.Equal(t, 403, res8.StatusCode)
	dataInBytes8, err := io.ReadAll(res8.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result8 map[string]string
	err = json.Unmarshal(dataInBytes8, &result8)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "Token theft detected", result8["message"])

	req9, err := http.NewRequest(http.MethodGet, testServer.URL+"/logout", nil)
	req9.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"]+";"+"sIdRefreshToken="+cookieData2["sIdRefreshToken"])
	assert.NoError(t, err)
	req9.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	res9, err := http.DefaultClient.Do(req9)
	assert.NoError(t, err)
	assert.Equal(t, 200, res9.StatusCode)
	cookieData3 := ExtractInfoFromResponseWhenAntiCSRFisNone(res9)
	assert.Equal(t, "", cookieData3["antiCsrf"])
	assert.Equal(t, "", cookieData3["sAccessToken"])
	assert.Equal(t, "", cookieData3["sIdRefreshToken"])
	assert.Equal(t, "", cookieData3["sRefreshToken"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["accessTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["idRefreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["refreshTokenExpiry"])
	dataInBytes9, err := ioutil.ReadAll(res9.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result9 map[string]bool
	err = json.Unmarshal(dataInBytes9, &result9)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, true, result9["message"])
	res9.Body.Close()

	defer AfterEach()
	defer testServer.Close()
}

func TestSessionVerifyMiddlewareWithAutoRefresh(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				ErrorHandlers: &sessmodels.ErrorHandlers{
					OnTokenTheftDetected: func(sessionHandle, userID string, req *http.Request, res http.ResponseWriter) error {
						res.WriteHeader(403)
						resp := make(map[string]string)
						resp["message"] = "Token theft detected"
						jsonResp, err := json.Marshal(resp)
						if err != nil {
							t.Errorf("Error happened in JSON marshal. Err: %s", err)
						}
						res.Write(jsonResp)
						return nil
					},
				},
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
		_, err := session.CreateNewSession(rw, "uniqueId", map[string]interface{}{}, map[string]interface{}{})
		if err != nil {
			rw.WriteHeader(500)
		}
		resp := make(map[string]bool)
		resp["message"] = true
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})

	mux.HandleFunc("/user/id", session.VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userId := sessionContainer.GetUserID()
		resp := make(map[string]string)
		resp["userId"] = userId
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

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

	customAntiCsrfCheck := true
	mux.HandleFunc("/user/handleV0", session.VerifySession(&sessmodels.VerifySessionOptions{
		AntiCsrfCheck: &customAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		handle := sessionContainer.GetHandle()
		resp := make(map[string]string)
		resp["handle"] = handle
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	customSessionRequiredVal := false
	mux.HandleFunc("/user/handleOptional", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredVal,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		resp := make(map[string]bool)
		if sessionContainer == nil {
			resp["message"] = false
		} else {
			resp["message"] = true
		}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	mux.HandleFunc("/logout", session.VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		err := sessionContainer.RevokeSession()
		if err != nil {
			rw.WriteHeader(500)
		}
		resp := make(map[string]bool)
		resp["message"] = true
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	cookieData := ExtractInfoFromResponse(res)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/user/id", nil)
	req1.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req1.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, 200, res1.StatusCode)

	dataInBytes, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]string
	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "uniqueId", result["userId"])
	res1.Body.Close()

	req2, err := http.NewRequest(http.MethodGet, testServer.URL+"/user/handleV0", nil)
	req2.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	assert.NoError(t, err)
	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 401, res2.StatusCode)
	dataInBytes2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result2 map[string]string
	err = json.Unmarshal(dataInBytes2, &result2)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "try refresh token", result2["message"])
	res2.Body.Close()

	req3, err := http.NewRequest(http.MethodGet, testServer.URL+"/user/handleOptional", nil)
	req3.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	assert.NoError(t, err)
	res3, err := http.DefaultClient.Do(req3)
	assert.NoError(t, err)
	assert.Equal(t, 200, res3.StatusCode)
	dataInBytes3, err := ioutil.ReadAll(res3.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result3 map[string]bool
	err = json.Unmarshal(dataInBytes3, &result3)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, true, result3["message"])
	res3.Body.Close()

	req4, err := http.NewRequest(http.MethodGet, testServer.URL+"/user/handleOptional", nil)
	assert.NoError(t, err)
	res4, err := http.DefaultClient.Do(req4)
	assert.NoError(t, err)
	assert.Equal(t, 200, res4.StatusCode)
	dataInBytes4, err := ioutil.ReadAll(res4.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result4 map[string]bool
	err = json.Unmarshal(dataInBytes4, &result4)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, false, result4["message"])
	res4.Body.Close()

	req5, err := http.NewRequest(http.MethodGet, testServer.URL+"/user/handleV0", nil)
	req5.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
	assert.NoError(t, err)
	req5.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res5, err := http.DefaultClient.Do(req5)
	assert.NoError(t, err)
	assert.Equal(t, 401, res5.StatusCode)
	dataInBytes5, err := ioutil.ReadAll(res5.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result5 map[string]string
	err = json.Unmarshal(dataInBytes5, &result5)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "unauthorised", result5["message"])
	res5.Body.Close()

	req6, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	req6.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req6.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res6, err := http.DefaultClient.Do(req6)
	cookieData2 := ExtractInfoFromResponse(res6)
	assert.NoError(t, err)
	assert.Equal(t, 200, res6.StatusCode)

	req7, err := http.NewRequest(http.MethodGet, testServer.URL+"/verifySession", nil)
	req7.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"]+";"+"sIdRefreshToken="+cookieData2["sIdRefreshToken"])
	req7.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	assert.NoError(t, err)
	res7, err := http.DefaultClient.Do(req7)
	assert.NoError(t, err)
	assert.Equal(t, 200, res7.StatusCode)

	req8, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	req8.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req8.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res8, err := http.DefaultClient.Do(req8)
	assert.NoError(t, err)
	assert.Equal(t, 403, res8.StatusCode)
	dataInBytes8, err := io.ReadAll(res8.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result8 map[string]string
	err = json.Unmarshal(dataInBytes8, &result8)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "Token theft detected", result8["message"])

	req9, err := http.NewRequest(http.MethodGet, testServer.URL+"/logout", nil)
	req9.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"]+";"+"sIdRefreshToken="+cookieData2["sIdRefreshToken"])
	assert.NoError(t, err)
	req9.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	res9, err := http.DefaultClient.Do(req9)
	assert.NoError(t, err)
	assert.Equal(t, 200, res9.StatusCode)
	cookieData3 := ExtractInfoFromResponseWhenAntiCSRFisNone(res9)
	assert.Equal(t, "", cookieData3["antiCsrf"])
	assert.Equal(t, "", cookieData3["sAccessToken"])
	assert.Equal(t, "", cookieData3["sIdRefreshToken"])
	assert.Equal(t, "", cookieData3["sRefreshToken"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["accessTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["idRefreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["refreshTokenExpiry"])
	dataInBytes9, err := ioutil.ReadAll(res9.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result9 map[string]bool
	err = json.Unmarshal(dataInBytes9, &result9)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, true, result9["message"])
	res9.Body.Close()

	defer AfterEach()
	defer testServer.Close()
}

func TestSessionVerifyMiddlewareWithDriverConfig(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	customapiBasePath := "/custom"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customapiBasePath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				ErrorHandlers: &sessmodels.ErrorHandlers{
					OnTokenTheftDetected: func(sessionHandle, userID string, req *http.Request, res http.ResponseWriter) error {
						res.WriteHeader(403)
						resp := make(map[string]string)
						resp["message"] = "Token theft detected"
						jsonResp, err := json.Marshal(resp)
						if err != nil {
							t.Errorf("Error happened in JSON marshal. Err: %s", err)
						}
						res.Write(jsonResp)
						return nil
					},
				},
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
		_, err := session.CreateNewSession(rw, "uniqueId", map[string]interface{}{}, map[string]interface{}{})
		if err != nil {
			rw.WriteHeader(500)
		}
		resp := make(map[string]bool)
		resp["message"] = true
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})

	mux.HandleFunc("/custom/user/id", session.VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userId := sessionContainer.GetUserID()
		resp := make(map[string]string)
		resp["userId"] = userId
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	customValForAntiCsrfCheck := true
	customSessionRequiredValue := true
	mux.HandleFunc("/custom/verifySession", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue,
		AntiCsrfCheck:   &customValForAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		session.GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue,
			AntiCsrfCheck:   &customValForAntiCsrfCheck,
		})
	}))

	customAntiCsrfCheck := true
	mux.HandleFunc("/custom/user/handleV0", session.VerifySession(&sessmodels.VerifySessionOptions{
		AntiCsrfCheck: &customAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		handle := sessionContainer.GetHandle()
		resp := make(map[string]string)
		resp["handle"] = handle
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	customSessionRequiredVal := false
	mux.HandleFunc("/custom/user/handleOptional", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredVal,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		resp := make(map[string]bool)
		if sessionContainer == nil {
			resp["message"] = false
		} else {
			resp["message"] = true
		}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	mux.HandleFunc("/custom/logout", session.VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		err := sessionContainer.RevokeSession()
		if err != nil {
			rw.WriteHeader(500)
		}
		resp := make(map[string]bool)
		resp["message"] = true
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	//this is never used. Rather the default api is used
	mux.HandleFunc("/custom/session/refresh", session.VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		resp := make(map[string]bool)
		resp["message"] = true
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)

	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	cookieData := ExtractInfoFromResponse(res)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/user/id", nil)
	req1.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req1.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, 200, res1.StatusCode)

	dataInBytes, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]string
	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "uniqueId", result["userId"])
	res1.Body.Close()

	req2, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/user/handleV0", nil)
	req2.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	assert.NoError(t, err)
	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 401, res2.StatusCode)
	dataInBytes2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result2 map[string]string
	err = json.Unmarshal(dataInBytes2, &result2)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "try refresh token", result2["message"])
	res2.Body.Close()

	req3, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/user/handleOptional", nil)
	req3.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	assert.NoError(t, err)
	res3, err := http.DefaultClient.Do(req3)
	assert.NoError(t, err)
	assert.Equal(t, 200, res3.StatusCode)
	dataInBytes3, err := ioutil.ReadAll(res3.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result3 map[string]bool
	err = json.Unmarshal(dataInBytes3, &result3)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, true, result3["message"])
	res3.Body.Close()

	req4, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/user/handleOptional", nil)
	assert.NoError(t, err)
	res4, err := http.DefaultClient.Do(req4)
	assert.NoError(t, err)
	assert.Equal(t, 200, res4.StatusCode)
	dataInBytes4, err := ioutil.ReadAll(res4.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result4 map[string]bool
	err = json.Unmarshal(dataInBytes4, &result4)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, false, result4["message"])
	res4.Body.Close()

	req5, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/user/handleV0", nil)
	req5.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
	assert.NoError(t, err)
	req5.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res5, err := http.DefaultClient.Do(req5)
	assert.NoError(t, err)
	assert.Equal(t, 401, res5.StatusCode)
	dataInBytes5, err := ioutil.ReadAll(res5.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result5 map[string]string
	err = json.Unmarshal(dataInBytes5, &result5)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "unauthorised", result5["message"])
	res5.Body.Close()

	req6, err := http.NewRequest(http.MethodPost, testServer.URL+"/custom/session/refresh", nil)
	req6.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req6.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res6, err := http.DefaultClient.Do(req6)
	cookieData2 := ExtractInfoFromResponse(res6)
	assert.NoError(t, err)
	assert.Equal(t, 200, res6.StatusCode)

	req7, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/verifySession", nil)
	req7.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"]+";"+"sIdRefreshToken="+cookieData2["sIdRefreshToken"])
	req7.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	assert.NoError(t, err)
	res7, err := http.DefaultClient.Do(req7)
	assert.NoError(t, err)
	assert.Equal(t, 200, res7.StatusCode)

	req8, err := http.NewRequest(http.MethodPost, testServer.URL+"/custom/session/refresh", nil)
	req8.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req8.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res8, err := http.DefaultClient.Do(req8)
	assert.NoError(t, err)
	assert.Equal(t, 403, res8.StatusCode)
	dataInBytes8, err := io.ReadAll(res8.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result8 map[string]string
	err = json.Unmarshal(dataInBytes8, &result8)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "Token theft detected", result8["message"])

	req9, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/logout", nil)
	req9.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"]+";"+"sIdRefreshToken="+cookieData2["sIdRefreshToken"])
	assert.NoError(t, err)
	req9.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	res9, err := http.DefaultClient.Do(req9)
	assert.NoError(t, err)
	assert.Equal(t, 200, res9.StatusCode)
	cookieData3 := ExtractInfoFromResponseWhenAntiCSRFisNone(res9)
	assert.Equal(t, "", cookieData3["antiCsrf"])
	assert.Equal(t, "", cookieData3["sAccessToken"])
	assert.Equal(t, "", cookieData3["sIdRefreshToken"])
	assert.Equal(t, "", cookieData3["sRefreshToken"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["accessTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["idRefreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["refreshTokenExpiry"])
	dataInBytes9, err := ioutil.ReadAll(res9.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result9 map[string]bool
	err = json.Unmarshal(dataInBytes9, &result9)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, true, result9["message"])
	res9.Body.Close()

	defer AfterEach()
	defer testServer.Close()
}

func TestSessionVerifyMiddlewareWithDriverConfigWithAutoRefresh(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	customapiBasePath := "/custom"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customapiBasePath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				ErrorHandlers: &sessmodels.ErrorHandlers{
					OnTokenTheftDetected: func(sessionHandle, userID string, req *http.Request, res http.ResponseWriter) error {
						res.WriteHeader(403)
						resp := make(map[string]string)
						resp["message"] = "Token theft detected"
						jsonResp, err := json.Marshal(resp)
						if err != nil {
							t.Errorf("Error happened in JSON marshal. Err: %s", err)
						}
						res.Write(jsonResp)
						return nil
					},
				},
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
		_, err := session.CreateNewSession(rw, "uniqueId", map[string]interface{}{}, map[string]interface{}{})
		if err != nil {
			rw.WriteHeader(500)
		}
		resp := make(map[string]bool)
		resp["message"] = true
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})

	mux.HandleFunc("/custom/user/id", session.VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userId := sessionContainer.GetUserID()
		resp := make(map[string]string)
		resp["userId"] = userId
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	customValForAntiCsrfCheck := true
	customSessionRequiredValue := true
	mux.HandleFunc("/custom/verifySession", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredValue,
		AntiCsrfCheck:   &customValForAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		session.GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &customSessionRequiredValue,
			AntiCsrfCheck:   &customValForAntiCsrfCheck,
		})
	}))

	customAntiCsrfCheck := true
	mux.HandleFunc("/custom/user/handleV0", session.VerifySession(&sessmodels.VerifySessionOptions{
		AntiCsrfCheck: &customAntiCsrfCheck,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		handle := sessionContainer.GetHandle()
		resp := make(map[string]string)
		resp["handle"] = handle
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	customSessionRequiredVal := false
	mux.HandleFunc("/custom/user/handleOptional", session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &customSessionRequiredVal,
	}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		resp := make(map[string]bool)
		if sessionContainer == nil {
			resp["message"] = false
		} else {
			resp["message"] = true
		}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	mux.HandleFunc("/custom/logout", session.VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		err := sessionContainer.RevokeSession()
		if err != nil {
			rw.WriteHeader(500)
		}
		resp := make(map[string]bool)
		resp["message"] = true
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	cookieData := ExtractInfoFromResponse(res)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/user/id", nil)
	req1.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req1.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, 200, res1.StatusCode)

	dataInBytes, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]string
	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "uniqueId", result["userId"])
	res1.Body.Close()

	req2, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/user/handleV0", nil)
	req2.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	assert.NoError(t, err)
	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 401, res2.StatusCode)
	dataInBytes2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result2 map[string]string
	err = json.Unmarshal(dataInBytes2, &result2)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "try refresh token", result2["message"])
	res2.Body.Close()

	req3, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/user/handleOptional", nil)
	req3.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	assert.NoError(t, err)
	res3, err := http.DefaultClient.Do(req3)
	assert.NoError(t, err)
	assert.Equal(t, 200, res3.StatusCode)
	dataInBytes3, err := ioutil.ReadAll(res3.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result3 map[string]bool
	err = json.Unmarshal(dataInBytes3, &result3)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, true, result3["message"])
	res3.Body.Close()

	req4, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/user/handleOptional", nil)
	assert.NoError(t, err)
	res4, err := http.DefaultClient.Do(req4)
	assert.NoError(t, err)
	assert.Equal(t, 200, res4.StatusCode)
	dataInBytes4, err := ioutil.ReadAll(res4.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result4 map[string]bool
	err = json.Unmarshal(dataInBytes4, &result4)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, false, result4["message"])
	res4.Body.Close()

	req5, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/user/handleV0", nil)
	req5.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
	assert.NoError(t, err)
	req5.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res5, err := http.DefaultClient.Do(req5)
	assert.NoError(t, err)
	assert.Equal(t, 401, res5.StatusCode)
	dataInBytes5, err := ioutil.ReadAll(res5.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result5 map[string]string
	err = json.Unmarshal(dataInBytes5, &result5)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "unauthorised", result5["message"])
	res5.Body.Close()

	req6, err := http.NewRequest(http.MethodPost, testServer.URL+"/custom/session/refresh", nil)
	req6.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req6.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res6, err := http.DefaultClient.Do(req6)
	cookieData2 := ExtractInfoFromResponse(res6)
	assert.NoError(t, err)
	assert.Equal(t, 200, res6.StatusCode)

	req7, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/verifySession", nil)
	req7.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"]+";"+"sIdRefreshToken="+cookieData2["sIdRefreshToken"])
	req7.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	assert.NoError(t, err)
	res7, err := http.DefaultClient.Do(req7)
	assert.NoError(t, err)
	assert.Equal(t, 200, res7.StatusCode)

	req8, err := http.NewRequest(http.MethodPost, testServer.URL+"/custom/session/refresh", nil)
	req8.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req8.Header.Add("anti-csrf", cookieData["antiCsrf"])
	assert.NoError(t, err)
	res8, err := http.DefaultClient.Do(req8)
	assert.NoError(t, err)
	assert.Equal(t, 403, res8.StatusCode)
	dataInBytes8, err := io.ReadAll(res8.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result8 map[string]string
	err = json.Unmarshal(dataInBytes8, &result8)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "Token theft detected", result8["message"])

	req9, err := http.NewRequest(http.MethodGet, testServer.URL+"/custom/logout", nil)
	req9.Header.Add("Cookie", "sAccessToken="+cookieData2["sAccessToken"]+";"+"sIdRefreshToken="+cookieData2["sIdRefreshToken"])
	assert.NoError(t, err)
	req9.Header.Add("anti-csrf", cookieData2["antiCsrf"])
	res9, err := http.DefaultClient.Do(req9)
	assert.NoError(t, err)
	assert.Equal(t, 200, res9.StatusCode)
	cookieData3 := ExtractInfoFromResponseWhenAntiCSRFisNone(res9)
	assert.Equal(t, "", cookieData3["antiCsrf"])
	assert.Equal(t, "", cookieData3["sAccessToken"])
	assert.Equal(t, "", cookieData3["sIdRefreshToken"])
	assert.Equal(t, "", cookieData3["sRefreshToken"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["accessTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["idRefreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["refreshTokenExpiry"])
	dataInBytes9, err := ioutil.ReadAll(res9.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result9 map[string]bool
	err = json.Unmarshal(dataInBytes9, &result9)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, true, result9["message"])
	res9.Body.Close()

	defer AfterEach()
	defer testServer.Close()
}
