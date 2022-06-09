package thirdparty

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultBackwardCompatibilityEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	tpConfig := &tpmodels.TypeInput{
		SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
			Providers: []tpmodels.TypeProvider{
				customProviderForEmailVerification,
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpConfig))
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

func TestBackwardCompatibilityEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""
	var thirdparty struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
	}

	tpConfig := &tpmodels.TypeInput{
		EmailVerificationFeature: &tpmodels.TypeInputEmailVerificationFeature{
			CreateAndSendCustomEmail: func(user tpmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
				email = user.Email
				emailVerifyLink = emailVerificationURLWithToken
				thirdparty = user.ThirdParty
				customCalled = true
			},
		},
		SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
			Providers: []tpmodels.TypeProvider{customProviderForEmailVerification},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpConfig))
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

func TestCustomOverrideEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""

	tpConfig := &tpmodels.TypeInput{
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
		SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
			Providers: []tpmodels.TypeProvider{customProviderForEmailVerification},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpConfig))
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
	tpConfig := &tpmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Service: &smtpService,
		},
		SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
			Providers: []tpmodels.TypeProvider{customProviderForEmailVerification},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpConfig))
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
