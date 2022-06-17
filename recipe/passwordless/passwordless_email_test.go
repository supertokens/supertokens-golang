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
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultBackwardCompatibilityPasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	tplConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.10", cdiVersion) == cdiVersion {
		return
	}

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	assert.True(t, PasswordlessLoginEmailSentForTest)
	assert.Equal(t, PasswordlessLoginEmailDataForTest.Email, "test@example.com")
	assert.NotNil(t, PasswordlessLoginEmailDataForTest.UrlWithLinkCode)
	assert.NotNil(t, PasswordlessLoginEmailDataForTest.UserInputCode)
}

func TestBackwardCompatibilityPasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	plessEmail := ""
	var code, urlWithCode *string
	var codeLife uint64

	tplConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
			CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
				plessEmail = email
				code = userInputCode
				urlWithCode = urlWithLinkCode
				codeLife = codeLifetime
				customCalled = true
				return nil
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.10", cdiVersion) == cdiVersion {
		return
	}

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordlessLoginEmailSentForTest)
	assert.Empty(t, PasswordlessLoginEmailDataForTest.Email)
	assert.Nil(t, PasswordlessLoginEmailDataForTest.UserInputCode)
	assert.Nil(t, PasswordlessLoginEmailDataForTest.UrlWithLinkCode)

	// Custom handler called
	assert.Equal(t, plessEmail, "test@example.com")
	assert.NotNil(t, code)
	assert.NotNil(t, urlWithCode)
	assert.NotZero(t, codeLife)
	assert.True(t, customCalled)
}

func TestCustomOverridePasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	plessEmail := ""
	var code, urlWithCode *string
	var codeLife uint64

	tplConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
		EmailDelivery: &emaildelivery.TypeInput{
			Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
				*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
					if input.PasswordlessLogin != nil {
						customCalled = true
						plessEmail = input.PasswordlessLogin.Email
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
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.10", cdiVersion) == cdiVersion {
		return
	}

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordlessLoginEmailSentForTest)
	assert.Empty(t, PasswordlessLoginEmailDataForTest.Email)
	assert.Nil(t, PasswordlessLoginEmailDataForTest.UserInputCode)
	assert.Nil(t, PasswordlessLoginEmailDataForTest.UrlWithLinkCode)

	// Custom handler called
	assert.Equal(t, plessEmail, "test@example.com")
	assert.NotNil(t, code)
	assert.NotNil(t, urlWithCode)
	assert.NotZero(t, codeLife)
	assert.True(t, customCalled)
}

func TestSMTPOverridePasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	plessEmail := ""
	var code, urlWithCode *string
	var codeLife uint64

	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
				if input.PasswordlessLogin != nil {
					plessEmail = input.PasswordlessLogin.Email
					code = input.PasswordlessLogin.UserInputCode
					urlWithCode = input.PasswordlessLogin.UrlWithLinkCode
					codeLife = input.PasswordlessLogin.CodeLifetime
					getContentCalled = true
				}
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	tplConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
		EmailDelivery: &emaildelivery.TypeInput{
			Service: &smtpService,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.10", cdiVersion) == cdiVersion {
		return
	}

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordlessLoginEmailSentForTest)
	assert.Empty(t, PasswordlessLoginEmailDataForTest.Email)
	assert.Nil(t, PasswordlessLoginEmailDataForTest.UserInputCode)
	assert.Nil(t, PasswordlessLoginEmailDataForTest.UrlWithLinkCode)

	assert.Equal(t, plessEmail, "test@example.com")
	assert.NotNil(t, code)
	assert.NotNil(t, urlWithCode)
	assert.NotZero(t, codeLife)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPServiceOverrideEmailTemplateForMagicLink(t *testing.T) {
	sendRawEmailCalled := false
	customCalled := false
	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				emailBody := input.Body
				assert.Contains(t, emailBody, "Please click the button below to sign in / up")
				assert.Contains(t, emailBody, "SuperTokens")
				assert.Contains(t, emailBody, "some@email.com")
				assert.Contains(t, emailBody, "http://someUrl")
				assert.Contains(t, emailBody, "1 minute")

				assert.NotContains(t, emailBody, "${")
				return nil
			}

			return originalImplementation
		},
	})
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
			Init(plessmodels.TypeInput{
				FlowType: "MAGIC_LINK",
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
				},
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						customCalled = true
						return nil
					},
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someUrl := "http://someUrl"
	err = SendEmail(emaildelivery.EmailType{
		PasswordlessLogin: &emaildelivery.PasswordlessLoginType{
			Email:            "some@email.com",
			UrlWithLinkCode:  &someUrl,
			PreAuthSessionId: "someSession",
			CodeLifetime:     60000,
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, customCalled, false)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPServiceOverrideEmailTemplateForOtp(t *testing.T) {
	sendRawEmailCalled := false
	customCalled := false
	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				emailBody := input.Body
				assert.Contains(t, emailBody, "Enter the below OTP in your login screen.")
				assert.Contains(t, emailBody, "SuperTokens")
				assert.Contains(t, emailBody, "some@email.com")
				assert.Contains(t, emailBody, "123456")
				assert.Contains(t, emailBody, "1 minute")

				assert.NotContains(t, emailBody, "${")
				return nil
			}

			return originalImplementation
		},
	})
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
			Init(plessmodels.TypeInput{
				FlowType: "MAGIC_LINK",
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
				},
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						customCalled = true
						return nil
					},
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someCode := "123456"
	err = SendEmail(emaildelivery.EmailType{
		PasswordlessLogin: &emaildelivery.PasswordlessLoginType{
			Email:            "some@email.com",
			UserInputCode:    &someCode,
			PreAuthSessionId: "someSession",
			CodeLifetime:     60000,
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, customCalled, false)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPServiceOverrideEmailTemplateForMagicLinkAndOtp(t *testing.T) {
	sendRawEmailCalled := false
	customCalled := false
	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				emailBody := input.Body
				assert.Contains(t, emailBody, "Please click the button below to sign in / up")
				assert.Contains(t, emailBody, "Enter the below OTP in your login screen.")
				assert.Contains(t, emailBody, "SuperTokens")
				assert.Contains(t, emailBody, "some@email.com")
				assert.Contains(t, emailBody, "http://someUrl")
				assert.Contains(t, emailBody, "123456")
				assert.Contains(t, emailBody, "1 minute")

				assert.NotContains(t, emailBody, "${")
				return nil
			}

			return originalImplementation
		},
	})
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
			Init(plessmodels.TypeInput{
				FlowType: "MAGIC_LINK",
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
				},
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						customCalled = true
						return nil
					},
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someCode := "123456"
	someUrl := "http://someUrl"
	err = SendEmail(emaildelivery.EmailType{
		PasswordlessLogin: &emaildelivery.PasswordlessLoginType{
			Email:            "some@email.com",
			UserInputCode:    &someCode,
			UrlWithLinkCode:  &someUrl,
			PreAuthSessionId: "someSession",
			CodeLifetime:     60000,
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, customCalled, false)
	assert.Equal(t, sendRawEmailCalled, true)
}
