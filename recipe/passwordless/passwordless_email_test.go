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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestSMTPOverrideEmailVerifyForPasswordlessUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	emailVerifyLink := ""
	var userInputCode *string

	evSmtpService := smtpService.MakeSMTPService(emaildelivery.SMTPServiceConfig{
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
				if input.EmailVerification != nil {
					email = input.EmailVerification.User.Email
					emailVerifyLink = input.EmailVerification.EmailVerifyLink
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
	tplSmtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
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
				if input.PasswordlessLogin != nil {
					userInputCode = input.PasswordlessLogin.UserInputCode
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
	tplConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
		EmailDelivery: &emaildelivery.TypeInput{
			Service: tplSmtpService,
		},
	}
	testServer := supertokensInitForTest(
		t,
		emailverification.Init(evmodels.TypeInput{
			Mode: evmodels.ModeOptional,
			EmailDelivery: &emaildelivery.TypeInput{
				Service: evSmtpService,
			},
		}),
		session.Init(&sessmodels.TypeInput{
			GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
				return sessmodels.CookieTransferMethod
			},
		}),
		Init(tplConfig),
	)
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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	json.Unmarshal(bodyBytes, &response)

	resp, err = unittesting.PasswordlessLoginWithCodeRequest(response["deviceId"].(string), response["preAuthSessionId"].(string), *userInputCode, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	sendRawEmailCalled = false // it would be true for the passwordless login, so reset it

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	json.Unmarshal(bodyBytes, &response)
	assert.Equal(t, response["status"], "EMAIL_ALREADY_VERIFIED_ERROR")

	// Default handler not called
	assert.False(t, emailverification.EmailVerificationEmailSentForTest)
	assert.Empty(t, emailverification.EmailVerificationDataForTest.User.Email)
	assert.Empty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, emailVerifyLink)
	assert.False(t, getContentCalled)
	assert.False(t, sendRawEmailCalled)
}

func TestCustomOverrideEmailVerifyForPasswordlessUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""

	tplConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
	}
	testServer := supertokensInitForTest(
		t,
		emailverification.Init(evmodels.TypeInput{
			Mode: evmodels.ModeOptional,
			EmailDelivery: &emaildelivery.TypeInput{
				Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
					sendEmail := *originalImplementation.SendEmail
					*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
						if input.EmailVerification != nil {
							customCalled = true
							email = input.EmailVerification.User.Email
							emailVerifyLink = input.EmailVerification.EmailVerifyLink
							return nil
						}
						return sendEmail(input, userContext)
					}
					return originalImplementation
				},
			},
		}),
		session.Init(&sessmodels.TypeInput{
			GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
				return sessmodels.CookieTransferMethod
			},
		}),
		Init(tplConfig),
	)
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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	json.Unmarshal(bodyBytes, &response)

	resp, err = unittesting.PasswordlessLoginWithCodeRequest(response["deviceId"].(string), response["preAuthSessionId"].(string), *PasswordlessLoginEmailDataForTest.UserInputCode, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	json.Unmarshal(bodyBytes, &response)
	assert.Equal(t, response["status"], "EMAIL_ALREADY_VERIFIED_ERROR")

	// Default handler not called
	assert.False(t, emailverification.EmailVerificationEmailSentForTest)
	assert.Empty(t, emailverification.EmailVerificationDataForTest.User.Email)
	assert.Empty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, emailVerifyLink)
	assert.False(t, customCalled)
}

func TestDefaultBackwardCompatibilityEmailVerifyForUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	tplConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
	}
	testServer := supertokensInitForTest(
		t,
		emailverification.Init(evmodels.TypeInput{Mode: evmodels.ModeOptional}),
		session.Init(&sessmodels.TypeInput{
			GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
				return sessmodels.CookieTransferMethod
			},
		}),
		Init(tplConfig),
	)
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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	json.Unmarshal(bodyBytes, &response)

	resp, err = unittesting.PasswordlessLoginWithCodeRequest(response["deviceId"].(string), response["preAuthSessionId"].(string), *PasswordlessLoginEmailDataForTest.UserInputCode, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	json.Unmarshal(bodyBytes, &response)
	assert.Equal(t, response["status"], "EMAIL_ALREADY_VERIFIED_ERROR")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, emailverification.EmailVerificationEmailSentForTest)
	assert.Empty(t, emailverification.EmailVerificationDataForTest.User.Email)
	assert.Empty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)
}

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
	testServer := supertokensInitForTest(
		t,
		session.Init(&sessmodels.TypeInput{
			GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
				return sessmodels.CookieTransferMethod
			},
		}),
		Init(tplConfig),
	)
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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := map[string]string{}

	err = json.Unmarshal(bodyBytes, &body)
	assert.NoError(t, err)

	assert.True(t, PasswordlessLoginEmailSentForTest)
	assert.Equal(t, PasswordlessLoginEmailDataForTest.Email, "test@example.com")
	assert.NotNil(t, PasswordlessLoginEmailDataForTest.UrlWithLinkCode)
	assert.NotNil(t, PasswordlessLoginEmailDataForTest.UserInputCode)

	// Test resend
	supertokens.ResetForTest()
	resp, err = unittesting.PasswordlessLoginResendRequest(body["deviceId"], body["preAuthSessionId"], testServer.URL)
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

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		plessEmail = input.PasswordlessLogin.Email
		code = input.PasswordlessLogin.UserInputCode
		urlWithCode = input.PasswordlessLogin.UrlWithLinkCode
		codeLife = input.PasswordlessLogin.CodeLifetime
		customCalled = true
		return nil
	}

	tplConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		EmailDelivery: &emaildelivery.TypeInput{
			Service: &emaildelivery.EmailDeliveryInterface{
				SendEmail: &sendEmail,
			},
		},
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
	}
	testServer := supertokensInitForTest(
		t,
		session.Init(&sessmodels.TypeInput{
			GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
				return sessmodels.CookieTransferMethod
			},
		}),
		Init(tplConfig),
	)
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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := map[string]string{}

	err = json.Unmarshal(bodyBytes, &body)
	assert.NoError(t, err)

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

	// Test resend
	customCalled = false
	plessEmail = ""
	code = nil
	urlWithCode = nil
	codeLife = 0

	resp, err = unittesting.PasswordlessLoginResendRequest(body["deviceId"], body["preAuthSessionId"], testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

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
	testServer := supertokensInitForTest(
		t,
		session.Init(&sessmodels.TypeInput{
			GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
				return sessmodels.CookieTransferMethod
			},
		}),
		Init(tplConfig),
	)
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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := map[string]string{}

	err = json.Unmarshal(bodyBytes, &body)
	assert.NoError(t, err)

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

	// Test resend
	customCalled = false
	plessEmail = ""
	code = nil
	urlWithCode = nil
	codeLife = 0

	resp, err = unittesting.PasswordlessLoginResendRequest(body["deviceId"], body["preAuthSessionId"], testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

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

	smtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
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
				if input.PasswordlessLogin != nil {
					plessEmail = input.PasswordlessLogin.Email
					code = input.PasswordlessLogin.UserInputCode
					urlWithCode = input.PasswordlessLogin.UrlWithLinkCode
					codeLife = input.PasswordlessLogin.CodeLifetime
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
	tplConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
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
		Init(tplConfig),
	)
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

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := map[string]string{}

	err = json.Unmarshal(bodyBytes, &body)
	assert.NoError(t, err)

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

	// Test resend
	getContentCalled = false
	sendRawEmailCalled = false
	plessEmail = ""
	code = nil
	urlWithCode = nil
	codeLife = 0

	resp, err = unittesting.PasswordlessLoginResendRequest(body["deviceId"], body["preAuthSessionId"], testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

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
	smtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
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
			(*originalImplementation.SendRawEmail) = func(input emaildelivery.EmailContent, userContext supertokens.UserContext) error {
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
					Service: smtpService,
				},
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
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
			TenantId:         "public",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, customCalled, false)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPServiceOverrideEmailTemplateForOtp(t *testing.T) {
	sendRawEmailCalled := false
	customCalled := false
	smtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
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
			(*originalImplementation.SendRawEmail) = func(input emaildelivery.EmailContent, userContext supertokens.UserContext) error {
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
					Service: smtpService,
				},
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
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
			TenantId:         "public",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, customCalled, false)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPServiceOverrideEmailTemplateForMagicLinkAndOtp(t *testing.T) {
	sendRawEmailCalled := false
	customCalled := false
	smtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
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
			(*originalImplementation.SendRawEmail) = func(input emaildelivery.EmailContent, userContext supertokens.UserContext) error {
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
					Service: smtpService,
				},
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
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
			TenantId:         "public",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, customCalled, false)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestThatMagicLinkUsesRightValueFromOriginFunction(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	plessEmail := ""
	var code, urlWithCode *string
	var codeLife uint64

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		plessEmail = input.PasswordlessLogin.Email
		code = input.PasswordlessLogin.UserInputCode
		urlWithCode = input.PasswordlessLogin.UrlWithLinkCode
		codeLife = input.PasswordlessLogin.CodeLifetime
		customCalled = true
		return nil
	}

	tplConfig := plessmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		EmailDelivery: &emaildelivery.TypeInput{
			Service: &emaildelivery.EmailDeliveryInterface{
				SendEmail: &sendEmail,
			},
		},
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
	}

	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain: "api.supertokens.io",
			AppName:   "SuperTokens",
			GetOrigin: func(request *http.Request, userContext supertokens.UserContext) (string, error) {
				if request == nil {
					return "https://supertokens.com", nil
				}
				// read request body
				decoder := json.NewDecoder(request.Body)
				var requestBody map[string]interface{}
				err := decoder.Decode(&requestBody)
				if err != nil {
					return "https://supertokens.com", nil
				}
				if requestBody["origin"] == nil {
					return "https://supertokens.com", nil
				}
				return requestBody["origin"].(string), nil
			},
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
			Init(tplConfig),
		},
	}

	err := supertokens.Init(config)
	assert.NoError(t, err)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
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

	body := map[string]string{
		"email":  "test@example.com",
		"origin": "localhost:2000",
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		t.Error(err.Error())
		return
	}

	resp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	body = map[string]string{}

	err = json.Unmarshal(bodyBytes, &body)
	assert.NoError(t, err)

	// Default handler not called
	assert.False(t, PasswordlessLoginEmailSentForTest)
	assert.Empty(t, PasswordlessLoginEmailDataForTest.Email)
	assert.Nil(t, PasswordlessLoginEmailDataForTest.UserInputCode)
	assert.Nil(t, PasswordlessLoginEmailDataForTest.UrlWithLinkCode)

	// Custom handler called
	assert.Equal(t, plessEmail, "test@example.com")
	assert.NotNil(t, code)
	assert.Equal(t, (*urlWithCode)[:21], "http://localhost:2000")
	assert.NotZero(t, codeLife)
	assert.True(t, customCalled)
}
