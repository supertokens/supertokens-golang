package thirdpartypasswordless

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/smsdelivery/twilioService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestSmsDefaultBackwardCompatibilityPasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	tplConfig := tplmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
			Enabled: true,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessPhoneLoginRequest("+919876543210", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	assert.True(t, passwordless.PasswordlessLoginSmsSentForTest)
	assert.Equal(t, passwordless.PasswordlessLoginSmsDataForTest.Phone, "+919876543210")
	assert.NotNil(t, passwordless.PasswordlessLoginSmsDataForTest.UrlWithLinkCode)
	assert.NotNil(t, passwordless.PasswordlessLoginSmsDataForTest.UserInputCode)
}

func TestSmsBackwardCompatibilityPasswordlessLogin(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	plessPhone := ""
	var code, urlWithCode *string
	var codeLife uint64

	tplConfig := tplmodels.TypeInput{
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
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessPhoneLoginRequest("+919876543210", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, passwordless.PasswordlessLoginSmsSentForTest)
	assert.Empty(t, passwordless.PasswordlessLoginSmsDataForTest.Phone)
	assert.Nil(t, passwordless.PasswordlessLoginSmsDataForTest.UserInputCode)
	assert.Nil(t, passwordless.PasswordlessLoginSmsDataForTest.UrlWithLinkCode)

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

	tplConfig := tplmodels.TypeInput{
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
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessPhoneLoginRequest("+919876543210", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, passwordless.PasswordlessLoginSmsSentForTest)
	assert.Empty(t, passwordless.PasswordlessLoginSmsDataForTest.Phone)
	assert.Nil(t, passwordless.PasswordlessLoginSmsDataForTest.UserInputCode)
	assert.Nil(t, passwordless.PasswordlessLoginSmsDataForTest.UrlWithLinkCode)

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

	tplConfig := tplmodels.TypeInput{
		FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
			Enabled: true,
		},
		SmsDelivery: &smsdelivery.TypeInput{
			Service: &twilioService,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tplConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessPhoneLoginRequest("+919876543210", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, passwordless.PasswordlessLoginSmsSentForTest)
	assert.Empty(t, passwordless.PasswordlessLoginSmsDataForTest.Phone)
	assert.Nil(t, passwordless.PasswordlessLoginSmsDataForTest.UserInputCode)
	assert.Nil(t, passwordless.PasswordlessLoginSmsDataForTest.UrlWithLinkCode)

	assert.Equal(t, plessPhone, "+919876543210")
	assert.NotNil(t, code)
	assert.NotNil(t, urlWithCode)
	assert.NotZero(t, codeLife)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawSmsCalled, true)
}
