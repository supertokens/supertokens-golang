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

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"gopkg.in/h2non/gock.v1"
)

func TestWithDisabledAPIDefaultSigninupAPIdoesnNotWork(t *testing.T) {
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
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							Google(tpmodels.GoogleConfig{
								ClientID:     "test",
								ClientSecret: "test-secret",
							}),
						},
					},
					Override: &tpmodels.OverrideStruct{
						APIs: func(originalImplementation tpmodels.APIInterface) tpmodels.APIInterface {
							*originalImplementation.SignInUpPOST = nil
							return originalImplementation
						},
					},
				},
			),
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

	signinupPostData := map[string]string{
		"thirdPartyId": "google",
		"code":         "abcdefghj",
		"redirectURI":  "http://127.0.0.1/callback",
	}

	postBody, err := json.Marshal(signinupPostData)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestMinimumConfigWithoutCodeForThirdPartyModule(t *testing.T) {
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
			}),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							customProvider6,
						},
					},
				},
			),
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

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		AuthCodeResponse: map[string]string{
			"access_token": "saodiasjodai",
		},
		RedirectUri: "http://127.0.0.1/callback",
	}

	postBody, err := json.Marshal(signinupPostData)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	cookieData := unittesting.ExtractInfoFromResponse(resp)

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

	assert.Equal(t, true, result["createdNewUser"])
	assert.Equal(t, "OK", result["status"])

	user := result["user"].(map[string]interface{})

	assert.Equal(t, "email@test.com", user["email"])
	assert.Equal(t, "custom", user["thirdParty"].(map[string]interface{})["id"])
	assert.Equal(t, "user", user["thirdParty"].(map[string]interface{})["userId"])

	assert.NotNil(t, cookieData["antiCsrf"])
	assert.NotNil(t, cookieData["sAccessToken"])
	assert.NotNil(t, cookieData["sRefreshToken"])
	assert.NotNil(t, cookieData["sIdRefreshToken"])
	assert.NotNil(t, cookieData["refreshTokenExpiry"])
	assert.NotNil(t, cookieData["refreshTokenHttpOnly"])
	assert.NotNil(t, cookieData["idRefreshTokenExpiry"])
	assert.NotNil(t, cookieData["idRefreshTokenHttpOnly"])
	assert.NotNil(t, cookieData["accessTokenExpiry"])
	assert.NotNil(t, cookieData["accessTokenHttpOnly"])
}

func TestMissingCodeAndAuthCodeResponse(t *testing.T) {
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
			}),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							customProvider6,
						},
					},
				},
			),
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

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		RedirectUri:  "http://127.0.0.1/callback",
	}

	postBody, err := json.Marshal(signinupPostData)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestMinimumConfigForThirdPartyModuleWithCode(t *testing.T) {
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
			}),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							customProvider1,
						},
					},
				},
			),
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
		JSON(map[string]string{"access_token": "abcdefghj"})

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
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	cookieData := unittesting.ExtractInfoFromResponse(resp)

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

	assert.Equal(t, true, result["createdNewUser"])
	assert.Equal(t, "OK", result["status"])

	user := result["user"].(map[string]interface{})

	assert.Equal(t, "email@test.com", user["email"])
	assert.Equal(t, "custom", user["thirdParty"].(map[string]interface{})["id"])
	assert.Equal(t, "user", user["thirdParty"].(map[string]interface{})["userId"])

	assert.NotNil(t, cookieData["antiCsrf"])
	assert.NotNil(t, cookieData["sAccessToken"])
	assert.NotNil(t, cookieData["sRefreshToken"])
	assert.NotNil(t, cookieData["sIdRefreshToken"])
	assert.NotNil(t, cookieData["refreshTokenExpiry"])
	assert.NotNil(t, cookieData["refreshTokenHttpOnly"])
	assert.NotNil(t, cookieData["idRefreshTokenExpiry"])
	assert.NotNil(t, cookieData["idRefreshTokenHttpOnly"])
	assert.NotNil(t, cookieData["accessTokenExpiry"])
	assert.NotNil(t, cookieData["accessTokenHttpOnly"])
}

func TestMinimumConfigForThirdPartyModuleWithEmailUnverified(t *testing.T) {
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
			}),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							customProvider5,
						},
					},
				},
			),
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
		JSON(map[string]string{"access_token": "abcdefghj"})

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
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	cookieData := unittesting.ExtractInfoFromResponse(resp)

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

	assert.Equal(t, true, result["createdNewUser"])
	assert.Equal(t, "OK", result["status"])

	user := result["user"].(map[string]interface{})

	isVerified, err := emailverification.IsEmailVerified(user["id"].(string), "FIXME")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, false, isVerified)

	assert.Equal(t, "email@test.com", user["email"])
	assert.Equal(t, "custom", user["thirdParty"].(map[string]interface{})["id"])
	assert.Equal(t, "user", user["thirdParty"].(map[string]interface{})["userId"])

	assert.NotNil(t, cookieData["antiCsrf"])
	assert.NotNil(t, cookieData["sAccessToken"])
	assert.NotNil(t, cookieData["sRefreshToken"])
	assert.NotNil(t, cookieData["sIdRefreshToken"])
	assert.NotNil(t, cookieData["refreshTokenExpiry"])
	assert.NotNil(t, cookieData["refreshTokenHttpOnly"])
	assert.NotNil(t, cookieData["idRefreshTokenExpiry"])
	assert.NotNil(t, cookieData["idRefreshTokenHttpOnly"])
	assert.NotNil(t, cookieData["accessTokenExpiry"])
	assert.NotNil(t, cookieData["accessTokenHttpOnly"])
}

func TestThirdPartyProviderDoesNotExist(t *testing.T) {
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
			session.Init(nil),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							customProvider1,
						},
					},
				},
			),
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

	postData := map[string]string{
		"thirdPartyId": "google",
		"code":         "abcdefghj",
		"redirectURI":  "http://127.0.0.1/callback",
	}

	postBody, err := json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	dataInBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var response map[string]string

	err = json.Unmarshal(dataInBytes, &response)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "The third party provider google seems to be missing from the backend configs.", response["message"])
}

func TestInvalidPostParamsForThirdPartyModule(t *testing.T) {
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
			session.Init(nil),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							customProvider1,
						},
					},
				},
			),
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

	//request where the postData was empty
	postData := map[string]string{}
	postBody, err := json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}
	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	dataInBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()
	var response map[string]string
	err = json.Unmarshal(dataInBytes, &response)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "Please provide the thirdPartyId in request body", response["message"])

	//request where the post data just had the thirdpartyid
	postData1 := map[string]string{
		"thirdPartyId": "custom",
	}
	postBody1, err := json.Marshal(postData1)
	if err != nil {
		t.Error(err.Error())
	}
	resp1, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody1))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusBadRequest, resp1.StatusCode)
	dataInBytes1, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp1.Body.Close()
	var response1 map[string]string
	err = json.Unmarshal(dataInBytes1, &response1)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "Please provide one of code or authCodeResponse in the request body", response1["message"])

	//request where the post data without redirect-uri
	postData2 := map[string]interface{}{
		"thirdPartyId": "custom",
		"code":         "32432432",
	}
	postBody2, err := json.Marshal(postData2)
	if err != nil {
		t.Error(err.Error())
	}
	resp2, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody2))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)
	dataInBytes2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp2.Body.Close()
	var response2 map[string]string
	err = json.Unmarshal(dataInBytes2, &response2)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, "Please provide the redirectURI in request body", response2["message"])
}

func TestEmailNotReturnedInGetProfileInfoFunction(t *testing.T) {
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
			}),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							customProvider3,
						},
					},
				},
			),
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

	assert.Equal(t, "NO_EMAIL_GIVEN_BY_PROVIDER", result["status"])
}

func TestGetUserByIdWhenUserDoesNotExist(t *testing.T) {
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
			}),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							customProvider1,
						},
					},
				},
			),
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
		"code":         "32432432",
		"redirectURI":  "http://localhost.org",
	}

	postBody, err := json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}

	gock.New(testServer.URL).EnableNetworking().Persist()
	gock.New("http://localhost:8080/").EnableNetworking().Persist()

	userDataBeforeSignup, err := GetUserByID("as")

	if err != nil {
		t.Error(err.Error())
	}

	assert.Nil(t, userDataBeforeSignup)

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
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

	user := result["user"].(map[string]interface{})
	userInfoAfterSignup, err := GetUserByID(user["id"].(string))
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, userInfoAfterSignup.ID, user["id"].(string))
	assert.Equal(t, userInfoAfterSignup.Email, user["email"].(string))
}

func TestGetUserByThirdPartyInfoWhenUserDoesNotExist(t *testing.T) {
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
			}),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							customProvider1,
						},
					},
				},
			),
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
		"code":         "32432432",
		"redirectURI":  "http://localhost.org",
	}

	postBody, err := json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}

	gock.New(testServer.URL).EnableNetworking().Persist()
	gock.New("http://localhost:8080/").EnableNetworking().Persist()

	userBegoreSignup, err := GetUserByThirdPartyInfo("custom", "user")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Nil(t, userBegoreSignup)

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
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

	user := result["user"].(map[string]interface{})
	userInfoAfterSignup, err := GetUserByThirdPartyInfo("custom", "user")
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, userInfoAfterSignup.ID, user["id"].(string))
	assert.Equal(t, userInfoAfterSignup.Email, user["email"].(string))
}
