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

package thirdpartypasswordless

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultBackwardCompatibilityPasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	tplConfig := tplmodels.TypeInput{
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

	assert.True(t, passwordless.PasswordlessLoginEmailSentForTest)
	assert.Equal(t, passwordless.PasswordlessLoginEmailDataForTest.Email, "test@example.com")
	assert.NotNil(t, passwordless.PasswordlessLoginEmailDataForTest.UrlWithLinkCode)
	assert.NotNil(t, passwordless.PasswordlessLoginEmailDataForTest.UserInputCode)

	// Test resend
	ResetForTest()
	resp, err = unittesting.PasswordlessLoginResendRequest(body["deviceId"], body["preAuthSessionId"], testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	assert.True(t, passwordless.PasswordlessLoginEmailSentForTest)
	assert.Equal(t, passwordless.PasswordlessLoginEmailDataForTest.Email, "test@example.com")
	assert.NotNil(t, passwordless.PasswordlessLoginEmailDataForTest.UrlWithLinkCode)
	assert.NotNil(t, passwordless.PasswordlessLoginEmailDataForTest.UserInputCode)
}

func TestBackwardCompatibilityPasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	plessEmail := ""
	var code, urlWithCode *string
	var codeLife uint64

	tplConfig := tplmodels.TypeInput{
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
	assert.False(t, passwordless.PasswordlessLoginEmailSentForTest)
	assert.Empty(t, passwordless.PasswordlessLoginEmailDataForTest.Email)
	assert.Nil(t, passwordless.PasswordlessLoginEmailDataForTest.UserInputCode)
	assert.Nil(t, passwordless.PasswordlessLoginEmailDataForTest.UrlWithLinkCode)

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

	tplConfig := tplmodels.TypeInput{
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
	assert.False(t, passwordless.PasswordlessLoginEmailSentForTest)
	assert.Empty(t, passwordless.PasswordlessLoginEmailDataForTest.Email)
	assert.Nil(t, passwordless.PasswordlessLoginEmailDataForTest.UserInputCode)
	assert.Nil(t, passwordless.PasswordlessLoginEmailDataForTest.UrlWithLinkCode)

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
	tplConfig := tplmodels.TypeInput{
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
	assert.False(t, passwordless.PasswordlessLoginEmailSentForTest)
	assert.Empty(t, passwordless.PasswordlessLoginEmailDataForTest.Email)
	assert.Nil(t, passwordless.PasswordlessLoginEmailDataForTest.UserInputCode)
	assert.Nil(t, passwordless.PasswordlessLoginEmailDataForTest.UrlWithLinkCode)

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

func TestDefaultBackwardCompatibilityEmailVerifyForPasswordlessUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	tplConfig := tplmodels.TypeInput{
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

	resp, err = unittesting.PasswordlessLoginWithCodeRequest(response["deviceId"].(string), response["preAuthSessionId"].(string), *passwordless.PasswordlessLoginEmailDataForTest.UserInputCode, testServer.URL)
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

func TestDefaultBackwardCompatibilityEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	tplConfig := tplmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
		Providers: []tpmodels.ProviderInput{
			customProviderForEmailVerification,
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

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		OAuthTokens: map[string]interface{}{
			"access_token": "saodiasjodai",
		},
	}

	postBody, err := json.Marshal(signinupPostData)
	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)

	cookies := resp.Cookies()

	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, emailverification.EmailVerificationEmailSentForTest)
	assert.Equal(t, emailverification.EmailVerificationDataForTest.User.Email, "test@example.com")
	assert.NotEmpty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)
}

// func TestBackwardCompatibilityEmailVerifyForPasswordlessUser(t *testing.T) {
// 	BeforeEach()
// 	unittesting.StartUpST("localhost", "8080")
// 	defer AfterEach()

// 	customCalled := false
// 	email := ""
// 	emailVerifyLink := ""

// 	tplConfig := tplmodels.TypeInput{
// 		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
// 		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
// 			Enabled: true,
// 		},
// 		EmailVerificationFeature: &tplmodels.TypeInputEmailVerificationFeature{
// 			CreateAndSendCustomEmail: func(user tplmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
// 				email = *user.Email
// 				emailVerifyLink = emailVerificationURLWithToken
// 				customCalled = true
// 			},
// 		},
// 	}
// 	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
// 	defer testServer.Close()

// 	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	cdiVersion, err := querier.GetQuerierAPIVersion()
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	if unittesting.MaxVersion("2.10", cdiVersion) == "2.10" {
// 		return
// 	}

// 	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)
// 	bodyBytes, err := ioutil.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var response map[string]interface{}
// 	json.Unmarshal(bodyBytes, &response)

// 	resp, err = unittesting.PasswordlessLoginWithCodeRequest(response["deviceId"].(string), response["preAuthSessionId"].(string), *passwordless.PasswordlessLoginEmailDataForTest.UserInputCode, testServer.URL)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)

// 	cookies := resp.Cookies()
// 	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)

// 	bodyBytes, err = ioutil.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	json.Unmarshal(bodyBytes, &response)
// 	assert.Equal(t, response["status"], "EMAIL_ALREADY_VERIFIED_ERROR")

// 	// Default handler not called
// 	assert.False(t, emailverification.EmailVerificationEmailSentForTest)
// 	assert.Empty(t, emailverification.EmailVerificationDataForTest.User.Email)
// 	assert.Empty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)

// 	// Custom handler called
// 	assert.Empty(t, email)
// 	assert.Empty(t, emailVerifyLink)
// 	assert.False(t, customCalled)
// }

// func TestBackwardCompatibilityEmailVerifyForThirdpartyUser(t *testing.T) {
// 	BeforeEach()
// 	unittesting.StartUpST("localhost", "8080")
// 	defer AfterEach()

// 	customCalled := false
// 	email := ""
// 	emailVerifyLink := ""
// 	var thirdparty *struct {
// 		ID     string `json:"id"`
// 		UserID string `json:"userId"`
// 	}

// 	tplConfig := tplmodels.TypeInput{
// 		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
// 		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
// 			Enabled: true,
// 		},
// 		EmailVerificationFeature: &tplmodels.TypeInputEmailVerificationFeature{
// 			CreateAndSendCustomEmail: func(user tplmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
// 				email = *user.Email
// 				emailVerifyLink = emailVerificationURLWithToken
// 				thirdparty = user.ThirdParty
// 				customCalled = true
// 			},
// 		},
// 		Providers: []tpmodels.TypeProvider{customProviderForEmailVerification},
// 	}
// 	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
// 	defer testServer.Close()

// 	signinupPostData := PostDataForCustomProvider{
// 		ThirdPartyId: "custom",
// 		AuthCodeResponse: map[string]string{
// 			"access_token": "saodiasjodai",
// 		},
// 		RedirectUri: "http://127.0.0.1/callback",
// 	}

// 	postBody, err := json.Marshal(signinupPostData)
// 	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
// 	assert.NoError(t, err)

// 	cookies := resp.Cookies()
// 	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)

// 	// Default handler not called
// 	assert.False(t, emailverification.EmailVerificationEmailSentForTest)
// 	assert.Empty(t, emailverification.EmailVerificationDataForTest.User.Email)
// 	assert.Empty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)

// 	// Custom handler called
// 	assert.Equal(t, email, "test@example.com")
// 	assert.NotEmpty(t, emailVerifyLink)
// 	assert.NotNil(t, thirdparty)
// 	assert.True(t, customCalled)
// }

func TestCustomOverrideEmailVerifyForPasswordlessUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""

	tplConfig := tplmodels.TypeInput{
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

	resp, err = unittesting.PasswordlessLoginWithCodeRequest(response["deviceId"].(string), response["preAuthSessionId"].(string), *passwordless.PasswordlessLoginEmailDataForTest.UserInputCode, testServer.URL)
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

func TestCustomOverrideEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""

	tplConfig := tplmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},

		Providers: []tpmodels.ProviderInput{customProviderForEmailVerification},
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

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		OAuthTokens: map[string]interface{}{
			"access_token": "saodiasjodai",
		},
	}

	postBody, err := json.Marshal(signinupPostData)
	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailverification.EmailVerificationEmailSentForTest)
	assert.Empty(t, emailverification.EmailVerificationDataForTest.User.Email)
	assert.Empty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.True(t, customCalled)
}

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
	tplConfig := tplmodels.TypeInput{
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

func TestSMTPOverrideEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	emailVerifyLink := ""

	smtpService := smtpService.MakeSMTPService(emaildelivery.SMTPServiceConfig{
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
	tplConfig := tplmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
		Providers: []tpmodels.ProviderInput{customProviderForEmailVerification},
	}
	testServer := supertokensInitForTest(
		t,
		emailverification.Init(evmodels.TypeInput{
			Mode: evmodels.ModeOptional,
			EmailDelivery: &emaildelivery.TypeInput{
				Service: smtpService,
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

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		OAuthTokens: map[string]interface{}{
			"access_token": "saodiasjodai",
		},
	}

	postBody, err := json.Marshal(signinupPostData)
	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailverification.EmailVerificationEmailSentForTest)
	assert.Empty(t, emailverification.EmailVerificationDataForTest.User.Email)
	assert.Empty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)

	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}
