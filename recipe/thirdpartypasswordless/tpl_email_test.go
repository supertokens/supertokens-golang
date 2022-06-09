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
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/emaildelivery/smtpService"
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
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	assert.True(t, passwordless.PasswordlessLoginEmailSentForTest)
	assert.Equal(t, passwordless.PasswordlessLoginDataForTest.Email, "test@example.com")
	assert.NotNil(t, passwordless.PasswordlessLoginDataForTest.UrlWithLinkCode)
	assert.NotNil(t, passwordless.PasswordlessLoginDataForTest.UserInputCode)
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
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, passwordless.PasswordlessLoginEmailSentForTest)
	assert.Empty(t, passwordless.PasswordlessLoginDataForTest.Email)
	assert.Nil(t, passwordless.PasswordlessLoginDataForTest.UserInputCode)
	assert.Nil(t, passwordless.PasswordlessLoginDataForTest.UrlWithLinkCode)

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
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, passwordless.PasswordlessLoginEmailSentForTest)
	assert.Empty(t, passwordless.PasswordlessLoginDataForTest.Email)
	assert.Nil(t, passwordless.PasswordlessLoginDataForTest.UserInputCode)
	assert.Nil(t, passwordless.PasswordlessLoginDataForTest.UrlWithLinkCode)

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
	tplConfig := tplmodels.TypeInput{
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
	assert.False(t, passwordless.PasswordlessLoginEmailSentForTest)
	assert.Empty(t, passwordless.PasswordlessLoginDataForTest.Email)
	assert.Nil(t, passwordless.PasswordlessLoginDataForTest.UserInputCode)
	assert.Nil(t, passwordless.PasswordlessLoginDataForTest.UrlWithLinkCode)

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
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	json.Unmarshal(bodyBytes, &response)

	resp, err = unittesting.PasswordlessLoginWithCodeRequest(response["deviceId"].(string), response["preAuthSessionId"].(string), *passwordless.PasswordlessLoginDataForTest.UserInputCode, testServer.URL)
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
		Providers: []tpmodels.TypeProvider{
			customProviderForEmailVerification,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		AuthCodeResponse: map[string]string{
			"access_token": "saodiasjodai",
		},
		RedirectUri: "http://127.0.0.1/callback",
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

func TestBackwardCompatibilityEmailVerifyForPasswordlessUser(t *testing.T) {
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
		EmailVerificationFeature: &tplmodels.TypeInputEmailVerificationFeature{
			CreateAndSendCustomEmail: func(user tplmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
				email = *user.Email
				emailVerifyLink = emailVerificationURLWithToken
				customCalled = true
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	json.Unmarshal(bodyBytes, &response)

	resp, err = unittesting.PasswordlessLoginWithCodeRequest(response["deviceId"].(string), response["preAuthSessionId"].(string), *passwordless.PasswordlessLoginDataForTest.UserInputCode, testServer.URL)
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

	// Custom handler called
	assert.Empty(t, email)
	assert.Empty(t, emailVerifyLink)
	assert.False(t, customCalled)
}

func TestBackwardCompatibilityEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""
	var thirdparty *struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
	}

	tplConfig := tplmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
		EmailVerificationFeature: &tplmodels.TypeInputEmailVerificationFeature{
			CreateAndSendCustomEmail: func(user tplmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
				email = *user.Email
				emailVerifyLink = emailVerificationURLWithToken
				thirdparty = user.ThirdParty
				customCalled = true
			},
		},
		Providers: []tpmodels.TypeProvider{customProviderForEmailVerification},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		AuthCodeResponse: map[string]string{
			"access_token": "saodiasjodai",
		},
		RedirectUri: "http://127.0.0.1/callback",
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
	assert.NotNil(t, thirdparty)
	assert.True(t, customCalled)
}

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
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessEmailLoginRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	json.Unmarshal(bodyBytes, &response)

	resp, err = unittesting.PasswordlessLoginWithCodeRequest(response["deviceId"].(string), response["preAuthSessionId"].(string), *passwordless.PasswordlessLoginDataForTest.UserInputCode, testServer.URL)
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
		Providers: []tpmodels.TypeProvider{customProviderForEmailVerification},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		AuthCodeResponse: map[string]string{
			"access_token": "saodiasjodai",
		},
		RedirectUri: "http://127.0.0.1/callback",
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
				if input.EmailVerification != nil {
					email = input.EmailVerification.User.Email
					emailVerifyLink = input.EmailVerification.EmailVerifyLink
					getContentCalled = true
				} else if input.PasswordlessLogin != nil {
					userInputCode = input.PasswordlessLogin.UserInputCode
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
	tplConfig := tplmodels.TypeInput{
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
				if input.EmailVerification != nil {
					email = input.EmailVerification.User.Email
					emailVerifyLink = input.EmailVerification.EmailVerifyLink
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
	tplConfig := tplmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		},
		EmailDelivery: &emaildelivery.TypeInput{
			Service: &smtpService,
		},
		Providers: []tpmodels.TypeProvider{customProviderForEmailVerification},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		AuthCodeResponse: map[string]string{
			"access_token": "saodiasjodai",
		},
		RedirectUri: "http://127.0.0.1/callback",
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
