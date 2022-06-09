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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	assert.True(t, PasswordlessLoginEmailSentForTest)
	assert.Equal(t, PasswordlessLoginDataForTest.Email, "test@example.com")
	assert.NotNil(t, PasswordlessLoginDataForTest.UrlWithLinkCode)
	assert.NotNil(t, PasswordlessLoginDataForTest.UserInputCode)
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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordlessLoginEmailSentForTest)
	assert.Empty(t, PasswordlessLoginDataForTest.Email)
	assert.Nil(t, PasswordlessLoginDataForTest.UserInputCode)
	assert.Nil(t, PasswordlessLoginDataForTest.UrlWithLinkCode)

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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordlessLoginEmailSentForTest)
	assert.Empty(t, PasswordlessLoginDataForTest.Email)
	assert.Nil(t, PasswordlessLoginDataForTest.UserInputCode)
	assert.Nil(t, PasswordlessLoginDataForTest.UrlWithLinkCode)

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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordlessLoginEmailSentForTest)
	assert.Empty(t, PasswordlessLoginDataForTest.Email)
	assert.Nil(t, PasswordlessLoginDataForTest.UserInputCode)
	assert.Nil(t, PasswordlessLoginDataForTest.UrlWithLinkCode)

	assert.Equal(t, plessEmail, "test@example.com")
	assert.NotNil(t, code)
	assert.NotNil(t, urlWithCode)
	assert.NotZero(t, codeLife)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}
