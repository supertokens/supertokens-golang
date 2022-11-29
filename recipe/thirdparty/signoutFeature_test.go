/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package thirdparty

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/supertokens/supertokens-golang/recipe/session"
// 	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
// 	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
// 	"github.com/supertokens/supertokens-golang/supertokens"
// 	"github.com/supertokens/supertokens-golang/test/unittesting"
// 	"gopkg.in/h2non/gock.v1"
// )

// func TestThatCallingTheAPIwithoutASessionShouldReturnOk(t *testing.T) {
// 	configValue := supertokens.TypeInput{
// 		Supertokens: &supertokens.ConnectionInfo{
// 			ConnectionURI: "http://localhost:8080",
// 		},
// 		AppInfo: supertokens.AppInfo{
// 			APIDomain:     "api.supertokens.io",
// 			AppName:       "SuperTokens",
// 			WebsiteDomain: "supertokens.io",
// 		},
// 		RecipeList: []supertokens.Recipe{
// 			session.Init(nil),
// 			Init(
// 				&tpmodels.TypeInput{
// 					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
// 						Providers: []tpmodels.TypeProvider{
// 							customProvider1,
// 						},
// 					},
// 				},
// 			),
// 		},
// 	}

// 	BeforeEach()
// 	unittesting.StartUpST("localhost", "8080")
// 	defer AfterEach()
// 	err := supertokens.Init(configValue)

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	mux := http.NewServeMux()
// 	testServer := httptest.NewServer(supertokens.Middleware(mux))
// 	defer testServer.Close()

// 	resp, err := http.Post(testServer.URL+"/auth/signout", "application/json", nil)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)

// 	dataInBytes, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	resp.Body.Close()
// 	var response map[string]string

// 	err = json.Unmarshal(dataInBytes, &response)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	assert.Equal(t, "OK", response["status"])

// 	assert.Equal(t, 0, len(resp.Cookies()))
// }

// func TestTheDefaultRouteShouldRevokeSession(t *testing.T) {
// 	customAntiCsrfVal := "VIA_TOKEN"
// 	configValue := supertokens.TypeInput{
// 		Supertokens: &supertokens.ConnectionInfo{
// 			ConnectionURI: "http://localhost:8080",
// 		},
// 		AppInfo: supertokens.AppInfo{
// 			APIDomain:     "api.supertokens.io",
// 			AppName:       "SuperTokens",
// 			WebsiteDomain: "supertokens.io",
// 		},
// 		RecipeList: []supertokens.Recipe{
// 			session.Init(
// 				&sessmodels.TypeInput{
// 					AntiCsrf: &customAntiCsrfVal,
// 				},
// 			),
// 			Init(
// 				&tpmodels.TypeInput{
// 					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
// 						Providers: []tpmodels.TypeProvider{
// 							customProvider1,
// 						},
// 					},
// 				},
// 			),
// 		},
// 	}

// 	BeforeEach()
// 	unittesting.StartUpST("localhost", "8080")
// 	defer AfterEach()
// 	err := supertokens.Init(configValue)

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	mux := http.NewServeMux()
// 	testServer := httptest.NewServer(supertokens.Middleware(mux))
// 	defer testServer.Close()

// 	defer gock.OffAll()
// 	gock.New("https://test.com/").
// 		Post("oauth/token").
// 		Reply(200).
// 		JSON(map[string]string{})

// 	postData := map[string]string{
// 		"thirdPartyId": "custom",
// 		"code":         "32432432",
// 		"redirectURI":  "http://127.0.0.1/callback",
// 	}

// 	postBody, err := json.Marshal(postData)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	gock.New(testServer.URL).EnableNetworking().Persist()
// 	gock.New("http://localhost:8080/").EnableNetworking().Persist()

// 	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)

// 	dataInBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	resp.Body.Close()

// 	var result map[string]interface{}

// 	err = json.Unmarshal(dataInBytes, &result)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	assert.Equal(t, "OK", result["status"])

// 	cookieData := unittesting.ExtractInfoFromResponse(resp)

// 	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/signout", nil)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
// 	req.Header.Add("anti-csrf", cookieData["antiCsrf"])

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	assert.Equal(t, http.StatusOK, res.StatusCode)

// 	cookieData1 := unittesting.ExtractInfoFromResponseWhenAntiCSRFisNone(res)

// 	assert.Equal(t, "", cookieData1["sAccessToken"])
// 	assert.Equal(t, "", cookieData1["sRefreshToken"])
// 	assert.Equal(t, "", cookieData1["sIdRefreshToken"])

// 	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["refreshTokenExpiry"])
// 	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["accessTokenExpiry"])
// 	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["idRefreshTokenExpiry"])

// 	assert.Equal(t, "", cookieData1["accessTokenDomain"])
// 	assert.Equal(t, "", cookieData1["refreshTokenDomain"])
// 	assert.Equal(t, "", cookieData1["idRefreshTokenDomain"])
// }

// func TestThatSignoutAPIReturnsTryRefreshTokenRefreshSessionAndSignoutShouldReturnOk(t *testing.T) {
// 	customAntiCsrfVal := "VIA_TOKEN"
// 	configValue := supertokens.TypeInput{
// 		Supertokens: &supertokens.ConnectionInfo{
// 			ConnectionURI: "http://localhost:8080",
// 		},
// 		AppInfo: supertokens.AppInfo{
// 			APIDomain:     "api.supertokens.io",
// 			AppName:       "SuperTokens",
// 			WebsiteDomain: "supertokens.io",
// 		},
// 		RecipeList: []supertokens.Recipe{
// 			session.Init(
// 				&sessmodels.TypeInput{
// 					AntiCsrf: &customAntiCsrfVal,
// 				},
// 			),
// 			Init(
// 				&tpmodels.TypeInput{
// 					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
// 						Providers: []tpmodels.TypeProvider{
// 							customProvider1,
// 						},
// 					},
// 				},
// 			),
// 		},
// 	}

// 	BeforeEach()
// 	unittesting.SetKeyValueInConfig("access_token_validity", strconv.Itoa(2))
// 	unittesting.StartUpST("localhost", "8080")
// 	defer AfterEach()
// 	err := supertokens.Init(configValue)

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	mux := http.NewServeMux()
// 	testServer := httptest.NewServer(supertokens.Middleware(mux))
// 	defer testServer.Close()

// 	defer gock.OffAll()
// 	gock.New("https://test.com/").
// 		Post("oauth/token").
// 		Reply(200).
// 		JSON(map[string]string{})

// 	postData := map[string]string{
// 		"thirdPartyId": "custom",
// 		"code":         "32432432",
// 		"redirectURI":  "http://127.0.0.1/callback",
// 	}

// 	postBody, err := json.Marshal(postData)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	gock.New(testServer.URL).EnableNetworking().Persist()
// 	gock.New("http://localhost:8080/").EnableNetworking().Persist()

// 	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)

// 	dataInBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	resp.Body.Close()

// 	var result map[string]interface{}

// 	err = json.Unmarshal(dataInBytes, &result)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	assert.Equal(t, "OK", result["status"])

// 	cookieData := unittesting.ExtractInfoFromResponse(resp)

// 	time.Sleep(5 * time.Second)

// 	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/signout", nil)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	req.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
// 	req.Header.Add("anti-csrf", cookieData["antiCsrf"])

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
// 	dataInBytes1, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	res.Body.Close()

// 	var result1 map[string]string
// 	err = json.Unmarshal(dataInBytes1, &result1)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	assert.Equal(t, "try refresh token", result1["message"])

// 	req1, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	req1.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
// 	req1.Header.Add("anti-csrf", cookieData["antiCsrf"])

// 	res1, err := http.DefaultClient.Do(req1)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	assert.Equal(t, http.StatusOK, res1.StatusCode)

// 	cookieData1 := unittesting.ExtractInfoFromResponse(res1)

// 	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/signout", nil)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	req2.Header.Add("Cookie", "sAccessToken="+cookieData1["sAccessToken"]+";"+"sIdRefreshToken="+cookieData1["sIdRefreshToken"])
// 	req2.Header.Add("anti-csrf", cookieData1["antiCsrf"])

// 	res2, err := http.DefaultClient.Do(req2)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	assert.Equal(t, http.StatusOK, res2.StatusCode)

// 	cookieData2 := unittesting.ExtractInfoFromResponseWhenAntiCSRFisNone(res2)

// 	assert.Equal(t, "", cookieData2["sAccessToken"])
// 	assert.Equal(t, "", cookieData2["sRefreshToken"])
// 	assert.Equal(t, "", cookieData2["sIdRefreshToken"])

// 	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData2["refreshTokenExpiry"])
// 	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData2["accessTokenExpiry"])
// 	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData2["idRefreshTokenExpiry"])

// 	assert.Equal(t, "", cookieData2["accessTokenDomain"])
// 	assert.Equal(t, "", cookieData2["refreshTokenDomain"])
// 	assert.Equal(t, "", cookieData2["idRefreshTokenDomain"])
// }
