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

package thirdparty

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultBackwardCompatibilityPasswordResetForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	testServer := supertokensInitForTest(
		t,
		session.Init(&sessmodels.TypeInput{
			GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
				return sessmodels.CookieTransferMethod
			},
		}),
		Init(nil),
		emailpassword.Init(nil),
	)
	defer testServer.Close()

	ManuallyCreateOrUpdateUser("public", "custom", "user-id", "test@example.com")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)
}

func TestCustomOverrideResetPasswordForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	epConfig := &epmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
				*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
					if input.PasswordReset != nil {
						customCalled = true
						email = input.PasswordReset.User.Email
						passwordResetLink = input.PasswordReset.PasswordResetLink
					}
					return nil
				}
				return originalImplementation
			},
		},
	}
	testServer := supertokensInitForTest(
		t,
		session.Init(&sessmodels.TypeInput{
			GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
				return sessmodels.CookieTransferMethod
			},
		}),
		Init(nil),
		emailpassword.Init(epConfig),
	)
	defer testServer.Close()

	ManuallyCreateOrUpdateUser("public", "custom", "user-id", "test@example.com")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, customCalled)
}

func TestSMTPOverridePasswordResetForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	passwordResetLink := ""

	smtpService := emailpassword.MakeSMTPService(emaildelivery.SMTPServiceConfig{
		Settings: emaildelivery.SMTPSettings{
			Host: "",
			From: emaildelivery.SMTPFrom{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPInterface) emaildelivery.SMTPInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.EmailContent, error) {
				if input.PasswordReset != nil {
					email = input.PasswordReset.User.Email
					passwordResetLink = input.PasswordReset.PasswordResetLink
					getContentCalled = true
				}
				return emaildelivery.EmailContent{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.EmailContent, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	epConfig := &epmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Service: smtpService,
		},
	}
	testServer := supertokensInitForTest(
		t,
		session.Init(&sessmodels.TypeInput{
			GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
				return sessmodels.CookieTransferMethod
			},
		}),
		Init(nil),
		emailpassword.Init(epConfig),
	)
	defer testServer.Close()

	ManuallyCreateOrUpdateUser("public", "custom", "user-id", "test@example.com")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, getContentCalled)
	assert.False(t, sendRawEmailCalled)
}
