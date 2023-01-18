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

package thirdpartypasswordless

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestMinimumConfigForThirdPartyPasswordlessWithEmailOrPhoneContactMethod(t *testing.T) {
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
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

	thirdPartyPasswordlessRecipe, err := getRecipeInstanceOrThrowError()
	assert.NoError(t, err)
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", thirdPartyPasswordlessRecipe.Config.FlowType)
}

func TestForThirdPartyPasswordLessCreateAndSendCustomTextMessageWithFlowTypeMagicLinkAndPhoneContactMethod(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						if userInputCode == nil && urlWithLinkCode != nil {
							isUserInputCodeAndUrlWithLinkCodeValid = true
						}
						return nil
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

	phone := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneBody, err := json.Marshal(phone)
	if err != nil {
		t.Error(err.Error())
	}

	phoneResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneResp.StatusCode)

	phoneDataInBytes, err := io.ReadAll(phoneResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneResp.Body.Close()

	var phoneResult map[string]interface{}
	err = json.Unmarshal(phoneDataInBytes, &phoneResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneResult["status"])
	assert.True(t, isUserInputCodeAndUrlWithLinkCodeValid)
}

func TestForThirdPartyPasswordlessCreateAndSendCustomMessageWithFlowTypeUserInputCodeAndMagicLinkAndPhoneContactMethod(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						if userInputCode != nil && urlWithLinkCode != nil {
							isUserInputCodeAndUrlWithLinkCodeValid = true
						}
						return nil
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

	phone := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneBody, err := json.Marshal(phone)
	if err != nil {
		t.Error(err.Error())
	}

	phoneResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneResp.StatusCode)

	phoneDataInBytes, err := io.ReadAll(phoneResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneResp.Body.Close()

	var phoneResult map[string]interface{}
	err = json.Unmarshal(phoneDataInBytes, &phoneResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneResult["status"])
	assert.True(t, isUserInputCodeAndUrlWithLinkCodeValid)
}

func TestWithThirdPartyPasswordLessCreateAndSendCustomTextMessageIfErrorIsThrownItShouldReturnA500Error(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						isUserInputCodeAndUrlWithLinkCodeValid = true
						return errors.New("test message")
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

	phone := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneBody, err := json.Marshal(phone)
	if err != nil {
		t.Error(err.Error())
	}

	phoneResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneBody))

	assert.NoError(t, err)
	assert.Equal(t, 500, phoneResp.StatusCode)
	assert.True(t, isUserInputCodeAndUrlWithLinkCodeValid)
}

func TestWithThirdPartyPasswordLessMinimumConfigWithEmailContactMethod(t *testing.T) {
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
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

	thirdPartyPasswordlessRecipe, err := getRecipeInstanceOrThrowError()
	assert.NoError(t, err)
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", thirdPartyPasswordlessRecipe.Config.FlowType)
}

func TestWithThirdPartyPasswordlessIfValidateEmailAdressIsCalledWithContactMethod(t *testing.T) {
	isValidateEmailAddressCalled := false
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
					ValidateEmailAddress: func(email interface{}) *string {
						isValidateEmailAddressCalled = true
						return nil
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

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	emailDataInBytes, err := io.ReadAll(emailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailResp.Body.Close()

	var emailResult map[string]interface{}
	err = json.Unmarshal(emailDataInBytes, &emailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailResult["status"])
	assert.True(t, isValidateEmailAddressCalled)
}

func TestWithThirdPartyPasswordlessIfValidateEmailAdressThrowsGenericErrorInCaseOfReturningAString(t *testing.T) {
	isValidateEmailAddressCalled := false
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
					ValidateEmailAddress: func(email interface{}) *string {
						isValidateEmailAddressCalled = true
						message := "test error"
						return &message
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

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	emailDataInBytes, err := io.ReadAll(emailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailResp.Body.Close()

	var emailResult map[string]interface{}
	err = json.Unmarshal(emailDataInBytes, &emailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "GENERAL_ERROR", emailResult["status"])
	assert.Equal(t, "test error", emailResult["message"])
	assert.True(t, isValidateEmailAddressCalled)
}

func TestForThirdPartyPasswordlessCreateAndSendCustomEmailWithFlowTypeUserInputCodeAndEmailContactMethod(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						if userInputCode != nil && urlWithLinkCode == nil {
							isUserInputCodeAndUrlWithLinkCodeValid = true
						}
						return nil
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

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	emailDataInBytes, err := io.ReadAll(emailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailResp.Body.Close()

	var emailResult map[string]interface{}
	err = json.Unmarshal(emailDataInBytes, &emailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailResult["status"])
	assert.True(t, isUserInputCodeAndUrlWithLinkCodeValid)
}

func TestWithThirdPartyPasswordlessCreateAndSendCustomEmailWithFlowTypeMagicLinkAndEmailContactMethod(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						if userInputCode == nil && urlWithLinkCode != nil {
							isUserInputCodeAndUrlWithLinkCodeValid = true
						}
						return nil
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

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	emailDataInBytes, err := io.ReadAll(emailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailResp.Body.Close()

	var emailResult map[string]interface{}
	err = json.Unmarshal(emailDataInBytes, &emailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailResult["status"])
	assert.True(t, isUserInputCodeAndUrlWithLinkCodeValid)
}

func TestWithThirdPartyPasswordlessCreateAndSendCustomEmailWithFlowTypeUserInputCodeAndMagicLinkAndEmailContactMethod(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						if userInputCode != nil && urlWithLinkCode != nil {
							isUserInputCodeAndUrlWithLinkCodeValid = true
						}
						return nil
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

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	emailDataInBytes, err := io.ReadAll(emailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailResp.Body.Close()

	var emailResult map[string]interface{}
	err = json.Unmarshal(emailDataInBytes, &emailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailResult["status"])
	assert.True(t, isUserInputCodeAndUrlWithLinkCodeValid)
}

func TestForThirdPartyPasswordLessThatForCreateAndCustomEmailIfErrorIsThrownTheStatusInTheResponseShouldBeA500Error(t *testing.T) {
	isCreateAndSendCustomEmailCalled := false
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						isCreateAndSendCustomEmailCalled = true
						return errors.New("test message")
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

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, 500, emailResp.StatusCode)
	assert.True(t, isCreateAndSendCustomEmailCalled)
}
