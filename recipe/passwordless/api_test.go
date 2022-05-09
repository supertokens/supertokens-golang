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

package passwordless

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestSignInUpFlowWithEmailUsingTheEmailOrPhoneContactMethod(t *testing.T) {
	var userInputCodeRef *string
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						userInputCodeRef = userInputCode
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
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", emailResult["flowType"])
	assert.Equal(t, 4, len(emailResult))

	//consume code API
	codeResendPostBody := map[string]interface{}{
		"deviceId":         emailResult["deviceId"],
		"userInputCode":    *userInputCodeRef,
		"preAuthSessionId": emailResult["preAuthSessionId"],
	}

	codeResendPostBodyJson, err := json.Marshal(codeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(codeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, codeResendResp.StatusCode)

	codeResendRespInBytes, err := io.ReadAll(codeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	codeResendResp.Body.Close()

	var codeResendResult map[string]interface{}
	err = json.Unmarshal(codeResendRespInBytes, &codeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", codeResendResult["status"])
	assert.True(t, codeResendResult["createdNewUser"].(bool))
	assert.Equal(t, 3, len(codeResendResult))
	assert.Equal(t, 4, len(codeResendResult["user"].(map[string]interface{})))
	assert.Equal(t, "test@example.com", codeResendResult["user"].(map[string]interface{})["email"])
}

func TestSignInUpFlowWithPhoneNumberUsingTheEmailOrPhoneContactMethod(t *testing.T) {
	var userInputCodeRef *string
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						userInputCodeRef = userInputCode
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
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", phoneResult["flowType"])
	assert.Equal(t, 4, len(phoneResult))

	//consume code API
	codeResendPostBody := map[string]interface{}{
		"deviceId":         phoneResult["deviceId"],
		"userInputCode":    *userInputCodeRef,
		"preAuthSessionId": phoneResult["preAuthSessionId"],
	}

	codeResendPostBodyJson, err := json.Marshal(codeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(codeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, codeResendResp.StatusCode)

	codeResendRespInBytes, err := io.ReadAll(codeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	codeResendResp.Body.Close()

	var codeResendResult map[string]interface{}
	err = json.Unmarshal(codeResendRespInBytes, &codeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", codeResendResult["status"])
	assert.True(t, codeResendResult["createdNewUser"].(bool))
	assert.Equal(t, 3, len(codeResendResult))
	assert.Equal(t, 4, len(codeResendResult["user"].(map[string]interface{})))
	assert.Equal(t, "+12345678901", codeResendResult["user"].(map[string]interface{})["phoneNumber"])
}

func TestCreatingACodeWithEmailAndThenResendingTheCodeAndCheckThatTheTheSendingCustomEmailFunctionIsWhileUsingTheEmailOrPhoneContactMethod(t *testing.T) {
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
			session.Init(nil),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						isCreateAndSendCustomEmailCalled = true
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
	assert.True(t, isCreateAndSendCustomEmailCalled)

	isCreateAndSendCustomEmailCalled = false

	codeResendPostBody := map[string]interface{}{
		"deviceId":         emailResult["deviceId"],
		"preAuthSessionId": emailResult["preAuthSessionId"],
	}

	codeResendPostBodyJson, err := json.Marshal(codeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(codeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, codeResendResp.StatusCode)

	codeResendRespInBytes, err := io.ReadAll(codeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	codeResendResp.Body.Close()

	var codeResendResult map[string]interface{}
	err = json.Unmarshal(codeResendRespInBytes, &codeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", codeResendResult["status"])
	assert.True(t, isCreateAndSendCustomEmailCalled)
}

func TestCreatingACodeWithPhoneAndThenResendingTheCodeAndCheckThatTheTheSendingCustomSmsFunctionIsWhileUsingTheEmailOrPhoneContactMethod(t *testing.T) {
	isCreateAndSendCustomTextMessageCalled := false
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						isCreateAndSendCustomTextMessageCalled = true
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
	assert.True(t, isCreateAndSendCustomTextMessageCalled)

	isCreateAndSendCustomTextMessageCalled = false

	codeResendPostBody := map[string]interface{}{
		"deviceId":         phoneResult["deviceId"],
		"preAuthSessionId": phoneResult["preAuthSessionId"],
	}

	codeResendPostBodyJson, err := json.Marshal(codeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(codeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, codeResendResp.StatusCode)

	codeResendRespInBytes, err := io.ReadAll(codeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	codeResendResp.Body.Close()

	var codeResendResult map[string]interface{}
	err = json.Unmarshal(codeResendRespInBytes, &codeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", codeResendResult["status"])
	assert.True(t, isCreateAndSendCustomTextMessageCalled)
}

func TestInvalidInputToCreateCodeApiUsingTheEmailOrPhoneContactMethod(t *testing.T) {
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
			Init(plessmodels.TypeInput{
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	radomData1 := map[string]interface{}{
		"phoneNumber": "+12345678901",
		"email":       "test@example.com",
	}

	randomBody1, err := json.Marshal(radomData1)
	if err != nil {
		t.Error(err.Error())
	}

	randomResp1, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(randomBody1))

	assert.Equal(t, http.StatusBadRequest, randomResp1.StatusCode)

	randomDataInBytes1, err := io.ReadAll(randomResp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	randomResp1.Body.Close()

	var randomResult1 map[string]interface{}
	err = json.Unmarshal(randomDataInBytes1, &randomResult1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "Please provide exactly one of email or phoneNumber", randomResult1["message"])

	radomData2 := map[string]interface{}{}

	randomBody2, err := json.Marshal(radomData2)
	if err != nil {
		t.Error(err.Error())
	}

	randomResp2, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(randomBody2))

	assert.Equal(t, http.StatusBadRequest, randomResp2.StatusCode)

	randomDataInBytes2, err := io.ReadAll(randomResp2.Body)
	if err != nil {
		t.Error(err.Error())
	}
	randomResp2.Body.Close()

	var randomResult2 map[string]interface{}
	err = json.Unmarshal(randomDataInBytes2, &randomResult2)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "Please provide exactly one of email or phoneNumber", randomResult2["message"])
}

func TestAddingPhoneNumberToAUsersInfoAndSigningInWillSignInTheSameUserUsingTheEmailOrPhoneContractMethod(t *testing.T) {
	var userInputCodeRef *string
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						userInputCodeRef = userInputCode
						return nil
					},
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						userInputCodeRef = userInputCode
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

	emailCodeResendPostBody := map[string]interface{}{
		"deviceId":         emailResult["deviceId"],
		"userInputCode":    *userInputCodeRef,
		"preAuthSessionId": emailResult["preAuthSessionId"],
	}

	emailCodeResendPostBodyJson, err := json.Marshal(emailCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	emailCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(emailCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailCodeResendResp.StatusCode)

	emailCodeResendRespInBytes, err := io.ReadAll(emailCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailCodeResendResp.Body.Close()

	var emailCodeResendResult map[string]interface{}
	err = json.Unmarshal(emailCodeResendRespInBytes, &emailCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailCodeResendResult["status"])

	emailForUpdating := emailCodeResendResult["user"].(map[string]interface{})["email"].(string)
	phoneNumberForUpdating := "+12345678901"

	_, err = UpdateUser(emailCodeResendResult["user"].(map[string]interface{})["id"].(string), &emailForUpdating, &phoneNumberForUpdating)

	assert.NoError(t, err)

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

	phoneCodeResendBody := map[string]interface{}{
		"deviceId":         phoneResult["deviceId"],
		"userInputCode":    *userInputCodeRef,
		"preAuthSessionId": phoneResult["preAuthSessionId"],
	}

	phoneCodeResendPostBodyJson, err := json.Marshal(phoneCodeResendBody)
	if err != nil {
		t.Error(err.Error())
	}

	phoneCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(phoneCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneCodeResendResp.StatusCode)

	phoneCodeResendRespInBytes, err := io.ReadAll(phoneCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneCodeResendResp.Body.Close()

	var phoneCodeResendResult map[string]interface{}
	err = json.Unmarshal(phoneCodeResendRespInBytes, &phoneCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneCodeResendResult["status"])

	assert.Equal(t, emailCodeResendResult["user"].(map[string]interface{})["id"], phoneCodeResendResult["user"].(map[string]interface{})["id"])
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
			session.Init(nil),
			Init(plessmodels.TypeInput{
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	codeResendPostBody := map[string]interface{}{
		"preAuthSessionId": "preAuthSessionId",
	}

	codeResendPostBodyJson, err := json.Marshal(codeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(codeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, codeResendResp.StatusCode)

	codeResendRespInBytes, err := io.ReadAll(codeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	codeResendResp.Body.Close()

	var codeResendResult map[string]interface{}
	err = json.Unmarshal(codeResendRespInBytes, &codeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "Please provide one of (linkCode) or (deviceId+userInputCode) and not both", codeResendResult["message"])
}

func TestConsumeCodeAPIWithMagicLink(t *testing.T) {
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
			Init(plessmodels.TypeInput{
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	codeInfo, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	invalidCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"linkCode":         "invalidLinkCode",
	}

	invalidCodeResendPostBodyJson, err := json.Marshal(invalidCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	invalidCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(invalidCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, invalidCodeResendResp.StatusCode)

	invalidCodeResendRespInBytes, err := io.ReadAll(invalidCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	invalidCodeResendResp.Body.Close()

	var invalidCodeResendResult map[string]interface{}
	err = json.Unmarshal(invalidCodeResendRespInBytes, &invalidCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "RESTART_FLOW_ERROR", invalidCodeResendResult["status"])

	validCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"linkCode":         codeInfo.OK.LinkCode,
	}

	validCodeResendPostBodyJson, err := json.Marshal(validCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	validCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(validCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validCodeResendResp.StatusCode)

	validCodeResendRespInBytes, err := io.ReadAll(validCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validCodeResendResp.Body.Close()

	var validCodeResendResult map[string]interface{}
	err = json.Unmarshal(validCodeResendRespInBytes, &validCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validCodeResendResult["status"])
	assert.True(t, validCodeResendResult["createdNewUser"].(bool))
	assert.Equal(t, 3, len(validCodeResendResult))
	assert.Equal(t, 4, len(validCodeResendResult["user"].(map[string]interface{})))
	assert.Equal(t, "test@example.com", validCodeResendResult["user"].(map[string]interface{})["email"])
}

func TestConsumeCodeAPIWithCode(t *testing.T) {
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
			Init(plessmodels.TypeInput{
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	codeInfo, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	invalidCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"userInputCode":    "invalidLinkCode",
		"deviceId":         codeInfo.OK.DeviceID,
	}

	invalidCodeResendPostBodyJson, err := json.Marshal(invalidCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	invalidCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(invalidCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, invalidCodeResendResp.StatusCode)

	invalidCodeResendRespInBytes, err := io.ReadAll(invalidCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	invalidCodeResendResp.Body.Close()

	var invalidCodeResendResult map[string]interface{}
	err = json.Unmarshal(invalidCodeResendRespInBytes, &invalidCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "INCORRECT_USER_INPUT_CODE_ERROR", invalidCodeResendResult["status"])
	assert.Equal(t, float64(1), invalidCodeResendResult["failedCodeInputAttemptCount"])
	assert.Equal(t, float64(5), invalidCodeResendResult["maximumCodeInputAttempts"])
	assert.Equal(t, 3, len(invalidCodeResendResult))

	validCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"userInputCode":    codeInfo.OK.UserInputCode,
		"deviceId":         codeInfo.OK.DeviceID,
	}

	validCodeResendPostBodyJson, err := json.Marshal(validCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	validCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(validCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validCodeResendResp.StatusCode)

	validCodeResendRespInBytes, err := io.ReadAll(validCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validCodeResendResp.Body.Close()

	var validCodeResendResult map[string]interface{}
	err = json.Unmarshal(validCodeResendRespInBytes, &validCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validCodeResendResult["status"])
	assert.True(t, validCodeResendResult["createdNewUser"].(bool))
	assert.Equal(t, 3, len(validCodeResendResult))
	assert.Equal(t, 4, len(validCodeResendResult["user"].(map[string]interface{})))
	assert.Equal(t, "test@example.com", validCodeResendResult["user"].(map[string]interface{})["email"])

	usedCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"userInputCode":    codeInfo.OK.UserInputCode,
		"deviceId":         codeInfo.OK.DeviceID,
	}

	usedCodeResendPostBodyJson, err := json.Marshal(usedCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	usedCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(usedCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, usedCodeResendResp.StatusCode)

	usedCodeResendRespInBytes, err := io.ReadAll(usedCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	usedCodeResendResp.Body.Close()

	var usedCodeResendResult map[string]interface{}
	err = json.Unmarshal(usedCodeResendRespInBytes, &usedCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "RESTART_FLOW_ERROR", usedCodeResendResult["status"])
}

func TestConsumeCodeAPIWithExpiredCode(t *testing.T) {
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
			Init(plessmodels.TypeInput{
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

	codeInfo, err := CreateCodeWithEmail("test@example.com", nil)
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

	expiredCodeResendRespInBytes, err := io.ReadAll(expiredCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	expiredCodeResendResp.Body.Close()

	var expiredCodeResendResult map[string]interface{}
	err = json.Unmarshal(expiredCodeResendRespInBytes, &expiredCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "EXPIRED_USER_INPUT_CODE_ERROR", expiredCodeResendResult["status"])
	assert.Equal(t, float64(1), expiredCodeResendResult["failedCodeInputAttemptCount"])
	assert.Equal(t, float64(5), expiredCodeResendResult["maximumCodeInputAttempts"])
	assert.Equal(t, 3, len(expiredCodeResendResult))
}

func TestCreateCodeAPIWithEmail(t *testing.T) {
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
			Init(plessmodels.TypeInput{
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

	validEmail := map[string]interface{}{
		"email": "test@example.com",
	}

	validEmailBody, err := json.Marshal(validEmail)
	if err != nil {
		t.Error(err.Error())
	}

	validEmailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(validEmailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validEmailResp.StatusCode)

	validEmailDataInBytes, err := io.ReadAll(validEmailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validEmailResp.Body.Close()

	var validEmailResult map[string]interface{}
	err = json.Unmarshal(validEmailDataInBytes, &validEmailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validEmailResult["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", validEmailResult["flowType"])
	assert.Equal(t, 4, len(validEmailResult))

	inValidEmail := map[string]interface{}{
		"email": "testple",
	}

	inValidEmailBody, err := json.Marshal(inValidEmail)
	if err != nil {
		t.Error(err.Error())
	}

	inValidEmailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(inValidEmailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, inValidEmailResp.StatusCode)

	inValidEmailDataInBytes, err := io.ReadAll(inValidEmailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	inValidEmailResp.Body.Close()

	var inValidEmailResult map[string]interface{}
	err = json.Unmarshal(inValidEmailDataInBytes, &inValidEmailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "GENERAL_ERROR", inValidEmailResult["status"])
	assert.Equal(t, "Email is invalid", inValidEmailResult["message"])
}

func TestCreateCodeAPIWithPhoneNumber(t *testing.T) {
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
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

	validPhoneNumber := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	validphoneNumberBody, err := json.Marshal(validPhoneNumber)
	if err != nil {
		t.Error(err.Error())
	}

	validphoneNumberResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(validphoneNumberBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validphoneNumberResp.StatusCode)

	validphoneNumberDataInBytes, err := io.ReadAll(validphoneNumberResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validphoneNumberResp.Body.Close()

	var validphoneNumberResult map[string]interface{}
	err = json.Unmarshal(validphoneNumberDataInBytes, &validphoneNumberResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validphoneNumberResult["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", validphoneNumberResult["flowType"])
	assert.Equal(t, 4, len(validphoneNumberResult))

	inValidphoneNumber := map[string]interface{}{
		"phoneNumber": "+123",
	}

	inValidphoneNumberBody, err := json.Marshal(inValidphoneNumber)
	if err != nil {
		t.Error(err.Error())
	}

	inValidphoneNumberResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(inValidphoneNumberBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, inValidphoneNumberResp.StatusCode)

	inValidphoneNumberDataInBytes, err := io.ReadAll(inValidphoneNumberResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	inValidphoneNumberResp.Body.Close()

	var inValidphoneNumberResult map[string]interface{}
	err = json.Unmarshal(inValidphoneNumberDataInBytes, &inValidphoneNumberResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "GENERAL_ERROR", inValidphoneNumberResult["status"])
	assert.Equal(t, "Phone number is invalid", inValidphoneNumberResult["message"])
}

func TestEmailExistsAPI(t *testing.T) {
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
			Init(plessmodels.TypeInput{
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

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	query := req.URL.Query()
	query.Add("email", "test@example.com")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	emailResp, err := http.DefaultClient.Do(req)
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
	assert.Equal(t, false, emailResult["exists"])

	codeInfo, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	_, err = ConsumeCodeWithLinkCode(codeInfo.OK.LinkCode, codeInfo.OK.PreAuthSessionID)
	assert.NoError(t, err)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	query1 := req.URL.Query()
	query1.Add("email", "test@example.com")
	req1.URL.RawQuery = query1.Encode()
	assert.NoError(t, err)
	emailResp1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp1.StatusCode)

	emailDataInBytes1, err := io.ReadAll(emailResp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailResp1.Body.Close()

	var emailResult1 map[string]interface{}
	err = json.Unmarshal(emailDataInBytes1, &emailResult1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailResult1["status"])
	assert.Equal(t, true, emailResult1["exists"])
}

func TestPhoneNumberExistsAPI(t *testing.T) {
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
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

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/phonenumber/exists", nil)
	query := req.URL.Query()
	query.Add("phoneNumber", "+1234567890")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	phoneResp, err := http.DefaultClient.Do(req)
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
	assert.Equal(t, false, phoneResult["exists"])

	codeInfo, err := CreateCodeWithPhoneNumber("+1234567890", nil)
	assert.NoError(t, err)

	_, err = ConsumeCodeWithLinkCode(codeInfo.OK.LinkCode, codeInfo.OK.PreAuthSessionID)
	assert.NoError(t, err)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/phonenumber/exists", nil)
	query1 := req.URL.Query()
	query1.Add("phoneNumber", "+1234567890")
	req1.URL.RawQuery = query1.Encode()
	assert.NoError(t, err)
	phoneResp1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneResp1.StatusCode)

	phoneDataInBytes1, err := io.ReadAll(phoneResp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneResp1.Body.Close()

	var phoneResult1 map[string]interface{}
	err = json.Unmarshal(phoneDataInBytes1, &phoneResult1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneResult1["status"])
	assert.Equal(t, true, phoneResult1["exists"])
}

func TestResendCodeAPI(t *testing.T) {
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
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

	codeInfo, err := CreateCodeWithPhoneNumber("+1234567890", nil)
	assert.NoError(t, err)

	validCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"deviceId":         codeInfo.OK.DeviceID,
	}

	validCodeResendPostBodyJson, err := json.Marshal(validCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	validCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(validCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validCodeResendResp.StatusCode)

	validCodeResendRespInBytes, err := io.ReadAll(validCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validCodeResendResp.Body.Close()

	var validCodeResendResult map[string]interface{}
	err = json.Unmarshal(validCodeResendRespInBytes, &validCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validCodeResendResult["status"])

	invalidCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": "asdasdasdasdsa",
		"deviceId":         "asdeflasdkjqee",
	}

	invalidCodeResendPostBodyJson, err := json.Marshal(invalidCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	invalidCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(invalidCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, invalidCodeResendResp.StatusCode)

	invalidCodeResendRespInBytes, err := io.ReadAll(invalidCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	invalidCodeResendResp.Body.Close()

	var invalidCodeResendResult map[string]interface{}
	err = json.Unmarshal(invalidCodeResendRespInBytes, &invalidCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "RESTART_FLOW_ERROR", invalidCodeResendResult["status"])
}
