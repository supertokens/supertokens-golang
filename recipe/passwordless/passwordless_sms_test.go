/*
 * Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/smsdelivery/twilioService"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestSmsDefaultBackwardCompatibilityPasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	plessConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
			Enabled: true,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(plessConfig))
	defer testServer.Close()

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.10", cdiVersion) == "2.10" {
		return
	}

	resp, err := unittesting.PasswordlessPhoneLoginRequest("+919876543210", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	assert.True(t, PasswordlessLoginSmsSentForTest)
	assert.Equal(t, PasswordlessLoginSmsDataForTest.Phone, "+919876543210")
	assert.NotNil(t, PasswordlessLoginSmsDataForTest.UrlWithLinkCode)
	assert.NotNil(t, PasswordlessLoginSmsDataForTest.UserInputCode)
}

func TestSmsBackwardCompatibilityPasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	plessPhone := ""
	var code, urlWithCode *string
	var codeLife uint64

	plessConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
			Enabled: true,
			CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
				plessPhone = phoneNumber
				code = userInputCode
				urlWithCode = urlWithLinkCode
				codeLife = codeLifetime
				customCalled = true
				return nil
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(plessConfig))
	defer testServer.Close()

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.10", cdiVersion) == "2.10" {
		return
	}

	resp, err := unittesting.PasswordlessPhoneLoginRequest("+919876543210", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordlessLoginSmsSentForTest)
	assert.Empty(t, PasswordlessLoginSmsDataForTest.Phone)
	assert.Nil(t, PasswordlessLoginSmsDataForTest.UserInputCode)
	assert.Nil(t, PasswordlessLoginSmsDataForTest.UrlWithLinkCode)

	// Custom handler called
	assert.Equal(t, plessPhone, "+919876543210")
	assert.NotNil(t, code)
	assert.NotNil(t, urlWithCode)
	assert.NotZero(t, codeLife)
	assert.True(t, customCalled)
}

func TestSmsCustomOverridePasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	plessPhone := ""
	var code, urlWithCode *string
	var codeLife uint64

	plessConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
			Enabled: true,
		},
		SmsDelivery: &smsdelivery.TypeInput{
			Override: func(originalImplementation smsdelivery.SmsDeliveryInterface) smsdelivery.SmsDeliveryInterface {
				*originalImplementation.SendSms = func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
					if input.PasswordlessLogin != nil {
						customCalled = true
						plessPhone = input.PasswordlessLogin.PhoneNumber
						code = input.PasswordlessLogin.UserInputCode
						urlWithCode = input.PasswordlessLogin.UrlWithLinkCode
						codeLife = input.PasswordlessLogin.CodeLifetime
					}
					return nil
				}
				return originalImplementation
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(plessConfig))
	defer testServer.Close()

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.10", cdiVersion) == "2.10" {
		return
	}

	resp, err := unittesting.PasswordlessPhoneLoginRequest("+919876543210", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordlessLoginSmsSentForTest)
	assert.Empty(t, PasswordlessLoginSmsDataForTest.Phone)
	assert.Nil(t, PasswordlessLoginSmsDataForTest.UserInputCode)
	assert.Nil(t, PasswordlessLoginSmsDataForTest.UrlWithLinkCode)

	// Custom handler called
	assert.Equal(t, plessPhone, "+919876543210")
	assert.NotNil(t, code)
	assert.NotNil(t, urlWithCode)
	assert.NotZero(t, codeLife)
	assert.True(t, customCalled)
}

func TestSmsTwilioOverridePasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawSmsCalled := false
	plessPhone := ""
	var code, urlWithCode *string
	var codeLife uint64

	serviceSid := "MS123"
	twilioService, err := twilioService.MakeTwilioService(smsdelivery.TwilioTypeInput{
		TwilioSettings: smsdelivery.TwilioServiceConfig{
			AccountSid:          "AC123",
			AuthToken:           "123",
			MessagingServiceSid: &serviceSid,
		},
		Override: func(originalImplementation smsdelivery.TwilioServiceInterface) smsdelivery.TwilioServiceInterface {
			*originalImplementation.GetContent = func(input smsdelivery.SmsType, userContext supertokens.UserContext) (smsdelivery.TwilioGetContentResult, error) {
				if input.PasswordlessLogin != nil {
					plessPhone = input.PasswordlessLogin.PhoneNumber
					code = input.PasswordlessLogin.UserInputCode
					urlWithCode = input.PasswordlessLogin.UrlWithLinkCode
					codeLife = input.PasswordlessLogin.CodeLifetime
					getContentCalled = true
				}
				return smsdelivery.TwilioGetContentResult{}, nil
			}

			*originalImplementation.SendRawSms = func(input smsdelivery.TwilioGetContentResult, userContext supertokens.UserContext) error {
				sendRawSmsCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	assert.NoError(t, err)

	plessConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
			Enabled: true,
		},
		SmsDelivery: &smsdelivery.TypeInput{
			Service: &twilioService,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(plessConfig))
	defer testServer.Close()

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.10", cdiVersion) == "2.10" {
		return
	}

	resp, err := unittesting.PasswordlessPhoneLoginRequest("+919876543210", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordlessLoginSmsSentForTest)
	assert.Empty(t, PasswordlessLoginSmsDataForTest.Phone)
	assert.Nil(t, PasswordlessLoginSmsDataForTest.UserInputCode)
	assert.Nil(t, PasswordlessLoginSmsDataForTest.UrlWithLinkCode)

	assert.Equal(t, plessPhone, "+919876543210")
	assert.NotNil(t, code)
	assert.NotNil(t, urlWithCode)
	assert.NotZero(t, codeLife)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawSmsCalled, true)
}

// func TestSupertokensServiceManually(t *testing.T) {
// 	serviceImpl := supertokensService.MakeSupertokensService(smsdelivery.SupertokensServiceConfig{
// 		ApiKey: "...",
// 	})

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
// 			Init(plessmodels.TypeInput{
// 				FlowType: "USER_INPUT_CODE",
// 				SmsDelivery: &smsdelivery.TypeInput{
// 					Service: &serviceImpl,
// 				},
// 				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
// 					Enabled: true,
// 				},
// 			}),
// 		},
// 	}

// 	BeforeEach()
// 	defer AfterEach()
// 	err := supertokens.Init(configValue)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	code := "123456"
// 	(*singletonInstance.SmsDelivery.IngredientInterfaceImpl.SendSms)(
// 		smsdelivery.SmsType{
// 			PasswordlessLogin: &smsdelivery.PasswordlessLoginType{
// 				PhoneNumber:      "...",
// 				UserInputCode:    &code,
// 				UrlWithLinkCode:  nil,
// 				CodeLifetime:     3600,
// 				PreAuthSessionId: "someSession",
// 			},
// 		},
// 		nil,
// 	)
// }

// func TestTwilioServiceManually(t *testing.T) {
// 	fromPhoneNumber := "..."
// 	// msgServiceSid := "someSid"
// 	twilioService, err := twilioService.MakeTwilioService(
// 		smsdelivery.TwilioTypeInput{
// 			TwilioSettings: smsdelivery.TwilioServiceConfig{
// 				AccountSid: "...",
// 				AuthToken:  "...",
// 				From:       &fromPhoneNumber,
// 				// MessagingServiceSid: &msgServiceSid,
// 			},
// 		},
// 	)
// 	assert.Nil(t, err)

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
// 			Init(plessmodels.TypeInput{
// 				FlowType: "USER_INPUT_CODE",

// 				SmsDelivery: &smsdelivery.TypeInput{
// 					Service: &twilioService,
// 				},
// 				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
// 					Enabled: true,
// 					ValidatePhoneNumber: func(phoneNumber interface{}) *string {
// 						return nil
// 					},
// 				},
// 			}),
// 		},
// 	}

// 	BeforeEach()
// 	unittesting.StartUpST("localhost", "8080")
// 	defer AfterEach()
// 	err = supertokens.Init(configValue)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	code := "123456"
// 	(*singletonInstance.SmsDelivery.IngredientInterfaceImpl.SendSms)(
// 		smsdelivery.SmsType{
// 			PasswordlessLogin: &smsdelivery.PasswordlessLoginType{
// 				PhoneNumber:      "...",
// 				UserInputCode:    &code,
// 				UrlWithLinkCode:  nil,
// 				CodeLifetime:     3600,
// 				PreAuthSessionId: "someSession",
// 			},
// 		},
// 		nil,
// 	)
// }
