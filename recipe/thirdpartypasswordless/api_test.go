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
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestForThirdPartyPasswordlessSignInUpFlowWithEmailUsingTheEmailOrPhoneContactMethod(t *testing.T) {
	var userInputCodeRef string
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		userInputCodeRef = *input.PasswordlessLogin.UserInputCode
		return nil
	}
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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	var validCreateCodeResponse map[string]interface{}
	err = json.Unmarshal(emailDataInBytes, &validCreateCodeResponse)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validCreateCodeResponse["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", validCreateCodeResponse["flowType"])

	data := map[string]interface{}{
		"preAuthSessionId": validCreateCodeResponse["preAuthSessionId"],
		"userInputCode":    userInputCodeRef,
		"deviceId":         validCreateCodeResponse["deviceId"],
	}

	condeConsumePostBody, err := json.Marshal(data)
	if err != nil {
		t.Error(err.Error())
	}

	validUserInputCodeResponse, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(condeConsumePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validUserInputCodeResponse.StatusCode)

	validUserInputCodeDataInBytes, err := io.ReadAll(validUserInputCodeResponse.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validUserInputCodeResponse.Body.Close()

	var validUserInputCodeDataResponse map[string]interface{}
	err = json.Unmarshal(validUserInputCodeDataInBytes, &validUserInputCodeDataResponse)
	if err != nil {
		t.Error(err.Error())
	}

	user := validUserInputCodeDataResponse["user"].(map[string]interface{})
	assert.Equal(t, "OK", validUserInputCodeDataResponse["status"])
	assert.True(t, validUserInputCodeDataResponse["createdNewUser"].(bool))
	assert.NotNil(t, user)
	assert.NotNil(t, user["email"])
	assert.NotNil(t, user["id"])
	assert.NotNil(t, user["timejoined"])
	assert.Nil(t, user["phoneNumber"])
}

func TestForThirdPartyPasswordlessSignUpSignInFlowWithPhoneNumberUsingEmailOrPhoneContactMethod(t *testing.T) {
	var userInputCodeRef string
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		userInputCodeRef = *input.PasswordlessLogin.UserInputCode
		return nil
	}
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
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &smsdelivery.SmsDeliveryInterface{
						SendSms: &sendSms,
					},
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	result := *unittesting.HttpResponseToConsumableInformation(phoneResp.Body)

	assert.Equal(t, "OK", result["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", result["flowType"])

	consumeCodePostData := map[string]interface{}{
		"preAuthSessionId": result["preAuthSessionId"],
		"userInputCode":    userInputCodeRef,
		"deviceId":         result["deviceId"],
	}

	consumeCodePostBody, err := json.Marshal(consumeCodePostData)
	if err != nil {
		t.Error(err.Error())
	}

	consumeCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, consumeCodeResp.StatusCode)

	codeConsumeResult := *unittesting.HttpResponseToConsumableInformation(consumeCodeResp.Body)

	user := codeConsumeResult["user"].(map[string]interface{})
	assert.Equal(t, "OK", codeConsumeResult["status"])
	assert.True(t, codeConsumeResult["createdNewUser"].(bool))
	assert.NotNil(t, user)
	assert.Nil(t, user["email"])
	assert.NotNil(t, user["id"])
	assert.NotNil(t, user["timejoined"])
	assert.NotNil(t, user["phoneNumber"])
}

func TestForThirdPartyPasswordlessCreatingACodeWithEmailAndThenResendingTheCodeAndCheckThatTheSendingCustomEmailFunctionIsCalledWhileUsingTheEmailOrPhoneContactMethod(t *testing.T) {
	isCreateAndSendCustomEmailCalled := false
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		isCreateAndSendCustomEmailCalled = true
		return nil
	}
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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	result := *unittesting.HttpResponseToConsumableInformation(emailResp.Body)

	assert.Equal(t, "OK", result["status"])
	assert.True(t, isCreateAndSendCustomEmailCalled)

	isCreateAndSendCustomEmailCalled = false

	codeResendPostData := map[string]interface{}{
		"deviceId":         result["deviceId"],
		"preAuthSessionId": result["preAuthSessionId"],
	}

	codeResendPostBody, err := json.Marshal(codeResendPostData)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendPostResp, err := http.Post(testServer.URL+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(codeResendPostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, codeResendPostResp.StatusCode)

	codeResendResult := *unittesting.HttpResponseToConsumableInformation(codeResendPostResp.Body)
	assert.Equal(t, "OK", codeResendResult["status"])
	assert.True(t, isCreateAndSendCustomEmailCalled)
}

func TestWithThirdPartyPasswordlessInvalidInputToCreateCodeAPIWhileUsingTheEmailOrPhoneContactMethod(t *testing.T) {
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

	postData := map[string]interface{}{
		"email":       "test@example.com",
		"phoneNumber": "+12345678901",
	}

	postBody, err := json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(postBody))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NoError(t, err)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)

	assert.Equal(t, "Please provide exactly one of email or phoneNumber", result["message"])

	postData = map[string]interface{}{}

	postBody, err = json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}

	resp1, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(postBody))
	assert.Equal(t, http.StatusBadRequest, resp1.StatusCode)
	assert.NoError(t, err)

	result1 := *unittesting.HttpResponseToConsumableInformation(resp1.Body)

	assert.Equal(t, "Please provide exactly one of email or phoneNumber", result1["message"])
}

func TestWithThirdPartyPasswordLessAddingPhoneNumberToAUsersInfoAndSigningInWillSignInTheSameUserUsingTheEmailOrPhoneContactMethod(t *testing.T) {
	var userInputCodeRef string
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		userInputCodeRef = *input.PasswordlessLogin.UserInputCode
		return nil
	}
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		userInputCodeRef = *input.PasswordlessLogin.UserInputCode
		return nil
	}

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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &smsdelivery.SmsDeliveryInterface{
						SendSms: &sendSms,
					},
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	emailCreateCodeResult := *unittesting.HttpResponseToConsumableInformation(emailResp.Body)

	assert.Equal(t, "OK", emailCreateCodeResult["status"])

	consumeCodePostData := map[string]interface{}{
		"preAuthSessionId": emailCreateCodeResult["preAuthSessionId"],
		"userInputCode":    userInputCodeRef,
		"deviceId":         emailCreateCodeResult["deviceId"],
	}

	consumeCodePostBody, err := json.Marshal(consumeCodePostData)
	if err != nil {
		t.Error(err.Error())
	}

	consumeCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, consumeCodeResp.StatusCode)

	emailUserInputCodeResponse := *unittesting.HttpResponseToConsumableInformation(consumeCodeResp.Body)

	assert.Equal(t, "OK", emailUserInputCodeResponse["status"])
	user := emailUserInputCodeResponse["user"].(map[string]interface{})

	phoneNumber := "+12345678901"
	UpdatePasswordlessUser(user["id"].(string), nil, &phoneNumber)

	phoneNumberPostData := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneNumberPostBody, err := json.Marshal(phoneNumberPostData)
	if err != nil {
		t.Error(err.Error())
	}

	phoneNumberPostResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneNumberPostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneNumberPostResp.StatusCode)

	phoneCreateCodeResponse := *unittesting.HttpResponseToConsumableInformation(phoneNumberPostResp.Body)

	assert.Equal(t, "OK", phoneCreateCodeResponse["status"])

	consumeCodePostData1 := map[string]interface{}{
		"preAuthSessionId": phoneCreateCodeResponse["preAuthSessionId"],
		"userInputCode":    userInputCodeRef,
		"deviceId":         phoneCreateCodeResponse["deviceId"],
	}

	consumeCodePostBody1, err := json.Marshal(consumeCodePostData1)
	if err != nil {
		t.Error(err.Error())
	}

	consumeCodeResp1, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody1))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, consumeCodeResp1.StatusCode)

	phoneUserInputCodeResponse := *unittesting.HttpResponseToConsumableInformation(consumeCodeResp1.Body)

	assert.Equal(t, "OK", phoneUserInputCodeResponse["status"])
	user1 := phoneUserInputCodeResponse["user"].(map[string]interface{})

	assert.Equal(t, user["id"], user1["id"])
}

func TestNotPassingAnyFieldsToConsumeCodeAPI(t *testing.T) {
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

	consumeCodePostData := map[string]interface{}{
		"preAuthSessionId": "preAuthSessionId",
	}

	consumeCodePostBody, err := json.Marshal(consumeCodePostData)
	if err != nil {
		t.Error(err.Error())
	}

	consumeCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, consumeCodeResp.StatusCode)

	emailUserInputCodeResponse := *unittesting.HttpResponseToConsumableInformation(consumeCodeResp.Body)

	assert.Equal(t, "Please provide one of (linkCode) or (deviceId+userInputCode) and not both", emailUserInputCodeResponse["message"])
}

func TestWithThirdPartyPasswordlessConsumeCodeAPIWithMagicLink(t *testing.T) {
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

	codeInfo, err := CreateCodeWithEmail("public", "test@example.com", nil)
	assert.NoError(t, err)

	consumeCodePostData := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"linkCode":         "invalidLinkCode",
	}

	consumeCodePostBody, err := json.Marshal(consumeCodePostData)
	if err != nil {
		t.Error(t, err)
	}

	invalidConsumeCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, invalidConsumeCodeResp.StatusCode)

	invalidLinkCodeResponse := *unittesting.HttpResponseToConsumableInformation(invalidConsumeCodeResp.Body)
	assert.Equal(t, "RESTART_FLOW_ERROR", invalidLinkCodeResponse["status"])

	consumeCodePostData = map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"linkCode":         codeInfo.OK.LinkCode,
	}

	consumeCodePostBody, err = json.Marshal(consumeCodePostData)
	if err != nil {
		t.Error(t, err)
	}

	validConsumeCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validConsumeCodeResp.StatusCode)

	validLinkCodeResponse := *unittesting.HttpResponseToConsumableInformation(validConsumeCodeResp.Body)
	assert.Equal(t, "OK", validLinkCodeResponse["status"])
	assert.True(t, validLinkCodeResponse["createdNewUser"].(bool))
	user := validLinkCodeResponse["user"].(map[string]interface{})
	assert.NotNil(t, user)
	assert.NotNil(t, user["email"])
	assert.NotNil(t, user["id"])
	assert.NotNil(t, user["timejoined"])
	assert.Nil(t, user["phoneNumber"])
}

func TestWithThirdPartyPasswordlessConsumeCodeAPIWithCode(t *testing.T) {
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

	codeInfo, err := CreateCodeWithEmail("public", "test@example.com", nil)
	assert.NoError(t, err)

	consumeCodePostData := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"deviceId":         codeInfo.OK.DeviceID,
		"userInputCode":    "invalidCode",
	}

	consumeCodePostBody, err := json.Marshal(consumeCodePostData)
	if err != nil {
		t.Error(t, err)
	}

	incorrectUserInputCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, incorrectUserInputCodeResp.StatusCode)

	incorrectUserInputCodeResponse := *unittesting.HttpResponseToConsumableInformation(incorrectUserInputCodeResp.Body)
	assert.Equal(t, "INCORRECT_USER_INPUT_CODE_ERROR", incorrectUserInputCodeResponse["status"])
	assert.Equal(t, float64(1), incorrectUserInputCodeResponse["failedCodeInputAttemptCount"])
	assert.Equal(t, float64(5), incorrectUserInputCodeResponse["maximumCodeInputAttempts"])

	consumeCodePostData = map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"deviceId":         codeInfo.OK.DeviceID,
		"userInputCode":    codeInfo.OK.UserInputCode,
	}

	consumeCodePostBody, err = json.Marshal(consumeCodePostData)
	if err != nil {
		t.Error(t, err)
	}

	correctUserInputCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, correctUserInputCodeResp.StatusCode)

	correctUserInputCodeResponse := *unittesting.HttpResponseToConsumableInformation(correctUserInputCodeResp.Body)
	assert.Equal(t, "OK", correctUserInputCodeResponse["status"])
	assert.True(t, correctUserInputCodeResponse["createdNewUser"].(bool))

	user := correctUserInputCodeResponse["user"].(map[string]interface{})
	assert.NotNil(t, user)
	assert.NotNil(t, user["email"])
	assert.NotNil(t, user["id"])
	assert.NotNil(t, user["timejoined"])
	assert.Nil(t, user["phoneNumber"])

	consumeCodePostData = map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"deviceId":         codeInfo.OK.DeviceID,
		"userInputCode":    codeInfo.OK.UserInputCode,
	}

	consumeCodePostBody, err = json.Marshal(consumeCodePostData)
	if err != nil {
		t.Error(t, err)
	}

	usedUserInputCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, usedUserInputCodeResp.StatusCode)

	usedUserInputCodeResponse := *unittesting.HttpResponseToConsumableInformation(usedUserInputCodeResp.Body)
	assert.Equal(t, "RESTART_FLOW_ERROR", usedUserInputCodeResponse["status"])
}

func TestWithThirdPartyPasswordLessConsumeCodeAPIWithExpiredCode(t *testing.T) {
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
				},
			}),
		},
	}
	BeforeEach()
	unittesting.SetKeyValueInConfig("passwordless_code_lifetime", "1000")
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

	codeInfo, err := CreateCodeWithEmail("public", "test@example.com", nil)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	expiredCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"userInputCode":    codeInfo.OK.UserInputCode,
		"deviceId":         codeInfo.OK.DeviceID,
	}

	expiredCodeResendPostBodyJson, err := json.Marshal(expiredCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	expiredCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(expiredCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, expiredCodeResendResp.StatusCode)

	expiredCodeResendResponse := *unittesting.HttpResponseToConsumableInformation(expiredCodeResendResp.Body)
	assert.Equal(t, "EXPIRED_USER_INPUT_CODE_ERROR", expiredCodeResendResponse["status"])
	assert.Equal(t, float64(1), expiredCodeResendResponse["failedCodeInputAttemptCount"])
	assert.Equal(t, float64(5), expiredCodeResendResponse["maximumCodeInputAttempts"])
}

func TestWithThirdPartyPasswordlessCreateCodeAPIWithEmail(t *testing.T) {
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

	validCreateCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validCreateCodeResp.StatusCode)

	validCreateCodeResponse := *unittesting.HttpResponseToConsumableInformation(validCreateCodeResp.Body)

	assert.Equal(t, "OK", validCreateCodeResponse["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", validCreateCodeResponse["flowType"])

	email = map[string]interface{}{
		"email": "testmpeom",
	}

	emailBody, err = json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	inValidCreateCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, inValidCreateCodeResp.StatusCode)

	inValidCreateCodeResponse := *unittesting.HttpResponseToConsumableInformation(inValidCreateCodeResp.Body)

	assert.Equal(t, "GENERAL_ERROR", inValidCreateCodeResponse["status"])
	assert.Equal(t, "Email is invalid", inValidCreateCodeResponse["message"])
}

func TestWithThirdPartyPasswordlessCreateCodeAPIWithPhoneNumber(t *testing.T) {
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

	phoneNumber := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneNumberBody, err := json.Marshal(phoneNumber)
	if err != nil {
		t.Error(err.Error())
	}

	validCreateCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneNumberBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validCreateCodeResp.StatusCode)

	validCreateCodeResponse := *unittesting.HttpResponseToConsumableInformation(validCreateCodeResp.Body)

	assert.Equal(t, "OK", validCreateCodeResponse["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", validCreateCodeResponse["flowType"])

	phoneNumber = map[string]interface{}{
		"phoneNumber": "231",
	}

	phoneNumberBody, err = json.Marshal(phoneNumber)
	if err != nil {
		t.Error(err.Error())
	}

	inValidCreateCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneNumberBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, inValidCreateCodeResp.StatusCode)

	inValidCreateCodeResponse := *unittesting.HttpResponseToConsumableInformation(inValidCreateCodeResp.Body)

	assert.Equal(t, "GENERAL_ERROR", inValidCreateCodeResponse["status"])
	assert.Equal(t, "Phone number is invalid", inValidCreateCodeResponse["message"])
}

func TestWithThirdPartyPasswordlessMagicLinkFormatInCreateCodeAPI(t *testing.T) {
	var magicLinkURL *url.URL
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		magicLinkURL, _ = url.Parse(*input.PasswordlessLogin.UrlWithLinkCode)
		return nil
	}
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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
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

	validCreateCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validCreateCodeResp.StatusCode)

	validCreateCodeResponse := *unittesting.HttpResponseToConsumableInformation(validCreateCodeResp.Body)

	assert.Equal(t, "OK", validCreateCodeResponse["status"])
	assert.Equal(t, "supertokens.io", magicLinkURL.Hostname())
	assert.Equal(t, "/auth/verify", magicLinkURL.Path)
	assert.Equal(t, "thirdpartypasswordless", magicLinkURL.Query().Get("rid"))
	assert.Equal(t, validCreateCodeResponse["preAuthSessionId"], magicLinkURL.Query().Get("preAuthSessionId"))
}

func TestWithThirdPartyPasswordlessEmailExistAPI(t *testing.T) {
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

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	query := req.URL.Query()
	query.Add("email", "test@example.com")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	emailDoesNotExistResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailDoesNotExistResp.StatusCode)

	emailDoesNotExistResponse := *unittesting.HttpResponseToConsumableInformation(emailDoesNotExistResp.Body)

	assert.Equal(t, "OK", emailDoesNotExistResponse["status"])
	assert.False(t, emailDoesNotExistResponse["exists"].(bool))

	codeInfo, err := CreateCodeWithEmail("public", "test@example.com", nil)
	assert.NoError(t, err)

	ConsumeCodeWithLinkCode("public", codeInfo.OK.LinkCode, codeInfo.OK.PreAuthSessionID)

	req, err = http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	query = req.URL.Query()
	query.Add("email", "test@example.com")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	emailExistsResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailExistsResp.StatusCode)

	emailExistsResponse := *unittesting.HttpResponseToConsumableInformation(emailExistsResp.Body)

	assert.Equal(t, "OK", emailExistsResponse["status"])
	assert.True(t, emailExistsResponse["exists"].(bool))
}

func TestWithThirdPartyPasswordlessPhoneNumberExistsAPI(t *testing.T) {
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

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/phonenumber/exists", nil)
	query := req.URL.Query()
	query.Add("phoneNumber", "+1234567890")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	phoneNumberDoesNotExistResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneNumberDoesNotExistResp.StatusCode)

	phoneNumberDoesNotExistResponse := *unittesting.HttpResponseToConsumableInformation(phoneNumberDoesNotExistResp.Body)

	assert.Equal(t, "OK", phoneNumberDoesNotExistResponse["status"])
	assert.False(t, phoneNumberDoesNotExistResponse["exists"].(bool))

	codeInfo, err := CreateCodeWithPhoneNumber("public", "+1234567890", nil)
	assert.NoError(t, err)

	ConsumeCodeWithLinkCode("public", codeInfo.OK.LinkCode, codeInfo.OK.PreAuthSessionID)

	req, err = http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/phonenumber/exists", nil)
	query = req.URL.Query()
	query.Add("phoneNumber", "+1234567890")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	phoneNumberExistsResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneNumberExistsResp.StatusCode)

	phoneNumberExistsResponse := *unittesting.HttpResponseToConsumableInformation(phoneNumberExistsResp.Body)

	assert.Equal(t, "OK", phoneNumberExistsResponse["status"])
	assert.True(t, phoneNumberExistsResponse["exists"].(bool))
}

func TestWithThirdPartyPasswordlessResendCodeAPI(t *testing.T) {
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

	codeInfo, err := CreateCodeWithPhoneNumber("public", "+1234567890", nil)
	assert.NoError(t, err)

	codeResendPostData := map[string]interface{}{
		"deviceId":         codeInfo.OK.DeviceID,
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
	}

	codeResendPostBody, err := json.Marshal(codeResendPostData)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendPostResp, err := http.Post(testServer.URL+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(codeResendPostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, codeResendPostResp.StatusCode)

	codeResendResult := *unittesting.HttpResponseToConsumableInformation(codeResendPostResp.Body)
	assert.Equal(t, "OK", codeResendResult["status"])

	codeResendPostData = map[string]interface{}{
		"deviceId":         "codeInfo",
		"preAuthSessionId": "PreAuthSessionID",
	}

	codeResendPostBody, err = json.Marshal(codeResendPostData)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendPostResp, err = http.Post(testServer.URL+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(codeResendPostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, codeResendPostResp.StatusCode)

	codeResendResult = *unittesting.HttpResponseToConsumableInformation(codeResendPostResp.Body)
	assert.Equal(t, "RESTART_FLOW_ERROR", codeResendResult["status"])
}
