package thirdpartypasswordless

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/smsdelivery/twilioService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestTwilioServiceOverrideForContactPhoneMethodThroughAPI(t *testing.T) {
	getContentCalled := false
	sendRawSmsCalled := false
	customCalled := false

	fromPhoneNumber := "someNumber"
	twilioService, err := twilioService.MakeTwilioService(
		smsdelivery.TwilioTypeInput{
			TwilioSettings: smsdelivery.TwilioServiceConfig{
				AccountSid:          "sid",
				AuthToken:           "token",
				From:                &fromPhoneNumber,
				MessagingServiceSid: nil,
			},
			Override: func(originalImplementation smsdelivery.TwilioServiceInterface) smsdelivery.TwilioServiceInterface {
				(*originalImplementation.GetContent) = func(input smsdelivery.SmsType, userContext supertokens.UserContext) (smsdelivery.TwilioGetContentResult, error) {
					getContentCalled = true
					return smsdelivery.TwilioGetContentResult{}, nil
				}

				(*originalImplementation.SendRawSms) = func(input smsdelivery.TwilioGetContentResult, userContext supertokens.UserContext) error {
					sendRawSmsCalled = true
					return nil
				}

				return originalImplementation
			},
		},
	)
	assert.Nil(t, err)

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

				SmsDelivery: &smsdelivery.TypeInput{
					Service: &twilioService,
				},
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						customCalled = true
						return nil
					},
					ValidatePhoneNumber: func(phoneNumber interface{}) *string {
						return nil
					},
				},
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err = supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	unittesting.PasswordlessPhoneLoginRequest("somePhone", testServer.URL)

	assert.Equal(t, customCalled, false)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawSmsCalled, true)
}
