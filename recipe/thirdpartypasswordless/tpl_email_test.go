package thirdpartypasswordless

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestEmailVerificationSMTPOverride(t *testing.T) {
	getContentCalled := false
	sendRawEmailCalled := false
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
				getContentCalled = true
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Google(tpmodels.GoogleConfig{ClientID: "id", ClientSecret: "secret"}),
				},
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
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

	err = SendEmail(emaildelivery.EmailType{
		EmailVerification: &emaildelivery.EmailVerificationType{
			User: emaildelivery.User{
				ID:    "someId",
				Email: "",
			},
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestPasswordlessLoginSMTPOverride(t *testing.T) {
	getContentCalled := false
	sendRawEmailCalled := false
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
				getContentCalled = true
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Google(tpmodels.GoogleConfig{ClientID: "id", ClientSecret: "secret"}),
				},
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
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

	err = SendEmail(emaildelivery.EmailType{
		PasswordlessLogin: &emaildelivery.PasswordlessLoginType{
			Email: "someEmail",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPServiceOverrideForContactEmailMethodThroughAPI(t *testing.T) {
	getContentCalled := false
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
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
				getContentCalled = true
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
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
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	unittesting.PasswordlessEmailLoginRequest("somone@email.com", testServer.URL)

	assert.Equal(t, customCalled, false)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestEmailVerificationSMTPOverrideThroughAPI(t *testing.T) {
	var customProviderForEmailVerification = tpmodels.TypeProvider{
		ID: "custom",
		Get: func(redirectURI, authCodeFromRequest *string, userContext *map[string]interface{}) tpmodels.TypeProviderGetResponse {
			return tpmodels.TypeProviderGetResponse{
				AccessTokenAPI: tpmodels.AccessTokenAPI{
					URL: "https://test.com/oauth/token",
				},
				AuthorisationRedirect: tpmodels.AuthorisationRedirect{
					URL: "https://test.com/oauth/auth",
				},
				GetProfileInfo: func(authCodeResponse interface{}, userContext *map[string]interface{}) (tpmodels.UserInfo, error) {
					if authCodeResponse.(map[string]interface{})["access_token"] == nil {
						return tpmodels.UserInfo{}, nil
					}
					return tpmodels.UserInfo{
						ID: "user",
						Email: &tpmodels.EmailStruct{
							ID:         "email@test.com",
							IsVerified: false,
						},
					}, nil
				},
				GetClientId: func(userContext *map[string]interface{}) string {
					return "supertokens"
				},
			}
		},
	}

	getContentCalled := false
	sendRawEmailCalled := false
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
				assert.NotNil(t, input.EmailVerification)
				assert.Equal(t, input.EmailVerification.User.Email, "email@test.com")
				getContentCalled = true
				return emaildelivery.SMTPGetContentResult{Body: "EmailVerification", ToEmail: input.EmailVerification.User.Email}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				assert.Equal(t, input.Body, "EmailVerification")
				assert.Equal(t, input.ToEmail, "email@test.com")
				sendRawEmailCalled = true
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
			session.Init(nil),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
				},
				Providers: []tpmodels.TypeProvider{
					customProviderForEmailVerification,
				},
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
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
	unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)

	assert.Nil(t, err)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPServiceOverrideEmailVerificationNotCalledForContactEmailMethodThroughAPI(t *testing.T) {
	getContentCalled := false
	sendRawEmailCalled := false
	customCalled := false
	var loginCode string
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
					loginCode = *input.PasswordlessLogin.UserInputCode
					return emaildelivery.SMTPGetContentResult{Body: "PasswordReset"}, nil
				} else if input.EmailVerification != nil {
					getContentCalled = true
					return emaildelivery.SMTPGetContentResult{Body: "EmailVerification"}, nil
				}
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				if input.Body == "PasswordReset" {
				} else if input.Body == "EmailVerification" {
					sendRawEmailCalled = true
				}
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
			session.Init(nil),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
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
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := unittesting.PasswordlessEmailLoginRequest("somone@email.com", testServer.URL)
	assert.NoError(t, err)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)

	resp, err = unittesting.PasswordlessLoginWithCodeRequest(
		response["deviceId"].(string),
		response["preAuthSessionId"].(string),
		loginCode,
		testServer.URL,
	)
	assert.NoError(t, err)
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	response = map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)
	assert.Equal(t, response["status"].(string), "OK")

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	response = map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)
	assert.Equal(t, response["status"].(string), "EMAIL_ALREADY_VERIFIED_ERROR")

	assert.Equal(t, customCalled, false)
	assert.Equal(t, getContentCalled, false)
	assert.Equal(t, sendRawEmailCalled, false)
}
