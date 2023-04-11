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

package thirdpartyemailpassword

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"gopkg.in/h2non/gock.v1"
)

func TestDefaultRouteShouldRevokeSession(t *testing.T) {
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
			Init(&tpepmodels.TypeInput{
				Providers: []tpmodels.TypeProvider{
					customProvider2,
				},
			}),
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

	defer gock.OffAll()
	gock.New("https://test.com/").
		Post("oauth/token").
		Reply(200).
		JSON(map[string]string{})

	postData := map[string]string{
		"thirdPartyId": "custom",
		"code":         "abcdefghj",
		"redirectURI":  "http://127.0.0.1/callback",
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
	cookieData := unittesting.ExtractInfoFromResponse(resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var result map[string]interface{}

	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", result["status"])
	resp1, err := unittesting.SignoutRequest(testServer.URL, cookieData["sAccessToken"], cookieData["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	cookieData1 := unittesting.ExtractInfoFromResponseWhenAntiCSRFisNone(resp1)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	dataInBytes1, err := io.ReadAll(resp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp1.Body.Close()

	var result1 map[string]interface{}

	err = json.Unmarshal(dataInBytes1, &result1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", result1["status"])
	assert.Equal(t, "", cookieData1["sAccessToken"])
	assert.Equal(t, "", cookieData1["sRefreshToken"])

	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["refreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData1["accessTokenExpiry"])

	assert.Equal(t, "", cookieData1["accessTokenDomain"])
	assert.Equal(t, "", cookieData1["refreshTokenDomain"])

	resp2, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	cookieData2 := unittesting.ExtractInfoFromResponse(resp2)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	dataInBytes2, err := io.ReadAll(resp2.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp2.Body.Close()

	var result2 map[string]interface{}

	err = json.Unmarshal(dataInBytes2, &result2)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", result2["status"])

	resp3, err := unittesting.SignoutRequest(testServer.URL, cookieData2["sAccessToken"], cookieData2["antiCsrf"])

	if err != nil {
		t.Error(err.Error())
	}
	cookieData3 := unittesting.ExtractInfoFromResponseWhenAntiCSRFisNone(resp3)
	assert.Equal(t, http.StatusOK, resp3.StatusCode)
	dataInBytes3, err := io.ReadAll(resp3.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp3.Body.Close()

	var result3 map[string]interface{}

	err = json.Unmarshal(dataInBytes3, &result3)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", result3["status"])
	assert.Equal(t, "", cookieData3["sAccessToken"])
	assert.Equal(t, "", cookieData3["sRefreshToken"])

	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["refreshTokenExpiry"])
	assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 GMT", cookieData3["accessTokenExpiry"])

	assert.Equal(t, "", cookieData3["accessTokenDomain"])
	assert.Equal(t, "", cookieData3["refreshTokenDomain"])
}
