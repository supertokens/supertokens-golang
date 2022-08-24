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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"gopkg.in/h2non/gock.v1"
)

func TestThirdPartyPasswordlessThatIfYouDisableTheSignInUpAPIItDoesNotWork(t *testing.T) {
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
				Providers: []tpmodels.TypeProvider{
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     "test",
						ClientSecret: "test-secret",
					}),
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

func TestWithThirdPartyPasswordlessMinimumConfigWithoutCodeForThirdPartyModyule(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					signinupCustomProvider6,
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

	signinupPostData := thirdparty.PostDataForCustomProvider{
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

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)

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

func TestWithThirdPartyPasswordlessMissingCodeAndAuthCodeResponse(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					signinupCustomProvider6,
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

	signinupPostData := thirdparty.PostDataForCustomProvider{
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

func TestWithThirdPartyPasswordlessMinimumConfigForThirdpartyModule(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
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

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)

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

func TestWithThirdPartyPasswordlessWithMinimumConfigForThirdPartyModuleEmailUnverified(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					signinupCustomProvider5,
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

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)

	assert.Equal(t, true, result["createdNewUser"])
	assert.Equal(t, "OK", result["status"])

	user := result["user"].(map[string]interface{})

	isVerified, err := emailverification.IsEmailVerified(user["id"].(string), nil)
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

func TestWithThirdPartyPasswordlessThirdPartyProviderDoesNotExistInConfig(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
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

	response := *unittesting.HttpResponseToConsumableInformation(resp.Body)

	assert.Equal(t, "The third party provider google seems to be missing from the backend configs.", response["message"])
}

func TestWithThirdPartyPasswordlessEmailNotReturnedInGetProfileInfoFunction(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					signinupCustomProvider3,
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

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "NO_EMAIL_GIVEN_BY_PROVIDER", result["status"])
}

func TestWithThirdPartyPasswordlessErrorThrownFromGetProfileInfoFunction(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					signinupCustomProvider4,
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
	assert.Equal(t, 500, resp.StatusCode)
}

func TestWithThirdPartyPasswordlessInvalidPostParamsForThirdPartyModule(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
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
	response := *unittesting.HttpResponseToConsumableInformation(resp.Body)
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
	response1 := *unittesting.HttpResponseToConsumableInformation(resp1.Body)
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
	response2 := *unittesting.HttpResponseToConsumableInformation(resp2.Body)
	assert.Equal(t, "Please provide the redirectURI in request body", response2["message"])
}

func TestWithThirdPartyPasswordlessGetUserByIdWhenUserDoesNotExist(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
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

	userDataBeforeSignup, err := GetUserByID("randomId")

	if err != nil {
		t.Error(err.Error())
	}

	assert.Nil(t, userDataBeforeSignup)

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)

	assert.Equal(t, "OK", result["status"])

	user := result["user"].(map[string]interface{})
	userInfoAfterSignup, err := GetUserByID(user["id"].(string))
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, userInfoAfterSignup.ID, user["id"].(string))
	assert.Equal(t, *userInfoAfterSignup.Email, user["email"].(string))
}

func TestGetUserByThirdPartyInfoWhenUserDoesNotExist(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
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
	assert.Equal(t, *userInfoAfterSignup.Email, user["email"].(string))
}
