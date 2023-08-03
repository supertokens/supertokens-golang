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
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestMinimumConfigWithEmailOrPhoneContactMethod(t *testing.T) {
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
			Init(plessmodels.TypeInput{
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

	passwordlessRecipe, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", passwordlessRecipe.Config.FlowType)
}

func TestAddingCustomValidatorsForPhoneAndEmailWithEmailOrPhoneContactMethod(t *testing.T) {
	isValidateEmailAddressCalled := false
	isValidatePhoneNumberCalled := false

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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					ValidateEmailAddress: func(email interface{}, tenantId string) *string {
						isValidateEmailAddressCalled = true
						return nil
					},
					ValidatePhoneNumber: func(phoneNumber interface{}, tenantId string) *string {
						isValidatePhoneNumberCalled = true
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
	assert.True(t, isValidateEmailAddressCalled)

	phoneNumber := map[string]interface{}{
		"phoneNumber": "+1234567890",
	}

	phoneBody, err := json.Marshal(phoneNumber)
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
	assert.True(t, isValidatePhoneNumberCalled)
}

func TestCustomFunctionToSendEmailWithEmailOrPhoneContactMethod(t *testing.T) {
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
			Init(plessmodels.TypeInput{
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

	var emailResult map[string]interface{}
	err = json.Unmarshal(emailDataInBytes, &emailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailResult["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", emailResult["flowType"])
	assert.True(t, isCreateAndSendCustomEmailCalled)
}

func TestCustomFunctionToSendTextSMSWithEmailOrPhoneContactMethod(t *testing.T) {
	isCreateAndSendCustomTextMessageCalled := false

	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		isCreateAndSendCustomTextMessageCalled = true
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
			Init(plessmodels.TypeInput{
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

	phoneNumber := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneBody, err := json.Marshal(phoneNumber)
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
	assert.True(t, isCreateAndSendCustomTextMessageCalled)
}

func TestMinimumConfigWithPhoneContactMethod(t *testing.T) {
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
			Init(plessmodels.TypeInput{
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

	passwordlessRecipe, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", passwordlessRecipe.Config.FlowType)
	assert.True(t, passwordlessRecipe.Config.ContactMethodPhone.Enabled)
}

func TestIfValidatePhoneNumberIsCalledWithPhoneContactMethod(t *testing.T) {
	isValidatePhoneNumberCalled := false
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					ValidatePhoneNumber: func(phoneNumber interface{}, tenantId string) *string {
						isValidatePhoneNumberCalled = true
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

	phoneNumber := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneBody, err := json.Marshal(phoneNumber)
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
	assert.True(t, isValidatePhoneNumberCalled)
}

func TestErrorMessageWithValidatePhoneNumberWithPhoneContactMethod(t *testing.T) {
	isValidatePhoneNumberCalled := false
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					ValidatePhoneNumber: func(phoneNumber interface{}, tenantId string) *string {
						message := "test error"
						isValidatePhoneNumberCalled = true
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

	phoneNumber := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneBody, err := json.Marshal(phoneNumber)
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

	assert.Equal(t, "GENERAL_ERROR", phoneResult["status"])
	assert.Equal(t, "test error", phoneResult["message"])
	assert.True(t, isValidatePhoneNumberCalled)
}

func TestCreateAndSendCustomMessageWithFlowTypeUserInputCodeAndPhoneContactNumber(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
	isOtherInputValid := false
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin.UserInputCode != nil && input.PasswordlessLogin.UrlWithLinkCode == nil {
			isUserInputCodeAndUrlWithLinkCodeValid = true
		}
		isOtherInputValid = true
		return nil

	}
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &smsdelivery.SmsDeliveryInterface{
						SendSms: &sendSms,
					},
				},
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

	phoneBody, err := json.Marshal(phoneNumber)
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
	assert.True(t, isOtherInputValid)
	assert.True(t, isUserInputCodeAndUrlWithLinkCodeValid)
}

func TestCreateAndSendCustomTextMessageWithFlowTypeMagicLinkAndPhoneContactMethod(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin.UserInputCode == nil && input.PasswordlessLogin.UrlWithLinkCode != nil {
			isUserInputCodeAndUrlWithLinkCodeValid = true
		}
		return nil
	}
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "MAGIC_LINK",
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &smsdelivery.SmsDeliveryInterface{
						SendSms: &sendSms,
					},
				},
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

	phoneBody, err := json.Marshal(phoneNumber)
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

func TestCreateAndSendCustomTextMessageWithFlowTypeUserInputCodeAndMagicLinkAndPhoneContactMethod(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin.UserInputCode != nil && input.PasswordlessLogin.UrlWithLinkCode != nil {
			isUserInputCodeAndUrlWithLinkCodeValid = true
		}
		return nil
	}
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &smsdelivery.SmsDeliveryInterface{
						SendSms: &sendSms,
					},
				},
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

	phoneBody, err := json.Marshal(phoneNumber)
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

func TestCreateAndSendCustomTextMessageIfErrorIsThrownItShouldContainA500Error(t *testing.T) {
	isCreateAndSendCustomTextMessageCalled := false
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		isCreateAndSendCustomTextMessageCalled = true
		return errors.New("test message")
	}
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &smsdelivery.SmsDeliveryInterface{
						SendSms: &sendSms,
					},
				},
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

	phoneBody, err := json.Marshal(phoneNumber)
	if err != nil {
		t.Error(err.Error())
	}

	phoneResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneBody))

	assert.NoError(t, err)
	assert.Equal(t, 500, phoneResp.StatusCode)
	assert.True(t, isCreateAndSendCustomTextMessageCalled)
}

func TestMinimumConfigWithEmailContactMethod(t *testing.T) {
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
			Init(plessmodels.TypeInput{
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

	passwordlessRecipe, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", passwordlessRecipe.Config.FlowType)
	assert.True(t, passwordlessRecipe.Config.ContactMethodEmail.Enabled)
}

func TestIfValidateEmailIsCalledWithEmailContactMethod(t *testing.T) {
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					ValidateEmailAddress: func(email interface{}, tenantId string) *string {
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

func TestValidateEmailWithGeneralErrorWithContactMethodSetToEmail(t *testing.T) {
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					ValidateEmailAddress: func(email interface{}, tenantId string) *string {
						message := "test error"
						isValidateEmailAddressCalled = true
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

func TestCreateAndSendCustomEmailWithFlowTypeMagicLinkAndCustomEmailContactMethod(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin.UserInputCode == nil && input.PasswordlessLogin.UrlWithLinkCode != nil {
			isUserInputCodeAndUrlWithLinkCodeValid = true
		}
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
			Init(plessmodels.TypeInput{
				FlowType: "MAGIC_LINK",
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

func TestCreateAndSendCustomTextMessageWithFlowTypeUserInputCodeAndMagicLinkAndEmailContactMethod(t *testing.T) {
	isUserInputCodeAndUrlWithLinkCodeValid := false
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin.UserInputCode != nil && input.PasswordlessLogin.UrlWithLinkCode != nil {
			isUserInputCodeAndUrlWithLinkCodeValid = true
		}
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
			Init(plessmodels.TypeInput{
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

func TestCreateAndSendCustomEmailIfErrorIsThrownTheStatusInTheResponseShouldBeA500Error(t *testing.T) {
	isCreateAndSendCustomEmailCalled := false

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		isCreateAndSendCustomEmailCalled = true
		return errors.New("test message")
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
			Init(plessmodels.TypeInput{
				FlowType: "MAGIC_LINK",
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

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, 500, emailResp.StatusCode)
	assert.True(t, isCreateAndSendCustomEmailCalled)
}

func TestPassingGetCustomUserInputCodeUsingDifferentCodes(t *testing.T) {
	var customCode string
	var userCodeSent *string

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		userCodeSent = input.PasswordlessLogin.UserInputCode
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
				},
				GetCustomUserInputCode: func(tenantId string, userContext supertokens.UserContext) (string, error) {
					customCode = unittesting.GenerateRandomCode(5)
					return customCode, nil
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
	assert.Equal(t, *userCodeSent, customCode)

	userCodeSent = nil
	customCode = ""

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
	assert.Equal(t, *userCodeSent, customCode)
}

func TestBasicOverrideUsageInPasswordLess(t *testing.T) {
	customDeviceId := "customDeviceId"
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
				},
				Override: &plessmodels.OverrideStruct{
					APIs: func(originalImplementation plessmodels.APIInterface) plessmodels.APIInterface {
						originalCodePost := *originalImplementation.CreateCodePOST
						*originalImplementation.CreateCodePOST = func(email, phoneNumber *string, tenantId string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.CreateCodePOSTResponse, error) {
							res, err := originalCodePost(email, phoneNumber, tenantId, options, userContext)
							res.OK.DeviceID = customDeviceId
							return res, err
						}
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

	assert.Equal(t, emailResult["deviceId"], customDeviceId)
}
