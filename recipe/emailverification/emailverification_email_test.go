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

package emailverification

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestBackwardCompatibilityServiceWithoutCustomFunction(t *testing.T) {
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
			Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
				GetEmailForUserID: func(userID string, userContext supertokens.UserContext) (evmodels.TypeEmailInfo, error) {
					return evmodels.TypeEmailInfo{
						OK: &struct{ Email string }{
							Email: "someEmail",
						},
					}, nil
				},
			}),
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
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

	SendEmail(emaildelivery.EmailType{
		EmailVerification: &emaildelivery.EmailVerificationType{
			User: emaildelivery.User{
				ID:    "someId",
				Email: "someEmail",
			},
			TenantId: "public",
		},
	})

	assert.Equal(t, EmailVerificationEmailSentForTest, true)
}

func TestBackwardCompatibilityServiceWithOverride(t *testing.T) {
	funcCalled := false
	overrideCalled := false
	ridInfo := ""
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
			Init(evmodels.TypeInput{
				Mode: "OPTIONAL",
				EmailDelivery: &emaildelivery.TypeInput{
					Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
						(*originalImplementation.SendEmail) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
							u, err := url.Parse(input.EmailVerification.EmailVerifyLink)
							if err != nil {
								return err
							}
							ridInfo = u.Query().Get("rid")
							overrideCalled = true
							return nil
						}
						return originalImplementation
					},
				},
				GetEmailForUserID: func(userID string, userContext supertokens.UserContext) (evmodels.TypeEmailInfo, error) {
					return evmodels.TypeEmailInfo{}, nil
				},
			}),
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
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

	SendEmail(emaildelivery.EmailType{
		EmailVerification: &emaildelivery.EmailVerificationType{
			User: emaildelivery.User{
				ID:    "someId",
				Email: "someEmail",
			},
			TenantId: "public",
		},
	})

	assert.Equal(t, EmailVerificationEmailSentForTest, false)
	assert.Equal(t, funcCalled, false)
	assert.Equal(t, overrideCalled, true)
	assert.Equal(t, ridInfo, "")
}

func TestSMTPServiceOverride(t *testing.T) {
	getContentCalled := false
	sendRawEmailCalled := false
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
				getContentCalled = true
				return emaildelivery.EmailContent{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.EmailContent, userContext supertokens.UserContext) error {
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
			Init(evmodels.TypeInput{
				Mode: "OPTIONAL",
				EmailDelivery: &emaildelivery.TypeInput{
					Service: smtpService,
				},
				GetEmailForUserID: func(userID string, userContext supertokens.UserContext) (evmodels.TypeEmailInfo, error) {
					return evmodels.TypeEmailInfo{}, nil
				},
			}),
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
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
			TenantId: "public",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPServiceOverrideDefaultEmailTemplate(t *testing.T) {
	sendRawEmailCalled := false
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
				assert.Contains(t, emailBody, "Please verify your email address")
				assert.Contains(t, emailBody, "SuperTokens")
				assert.Contains(t, emailBody, "some@email.com")

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
			Init(evmodels.TypeInput{
				Mode: "OPTIONAL",
				EmailDelivery: &emaildelivery.TypeInput{
					Service: smtpService,
				},
				GetEmailForUserID: func(userID string, userContext supertokens.UserContext) (evmodels.TypeEmailInfo, error) {
					return evmodels.TypeEmailInfo{}, nil
				},
			}),
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
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
				Email: "some@email.com",
			},
			TenantId: "public",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestThatLinkUsesResultFromOriginFunction(t *testing.T) {
	link := ""
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain: "api.supertokens.io",
			AppName:   "SuperTokens",
			GetOrigin: func(request *http.Request, userContext supertokens.UserContext) (string, error) {
				return (*userContext)["link"].(string), nil
			},
		},
		RecipeList: []supertokens.Recipe{
			Init(evmodels.TypeInput{
				Mode: "OPTIONAL",
				EmailDelivery: &emaildelivery.TypeInput{
					Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
						(*originalImplementation.SendEmail) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
							link = input.EmailVerification.EmailVerifyLink
							return nil
						}
						return originalImplementation
					},
				},
			}),
			session.Init(nil),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	email := "test@exmaple.com"
	resp, err := SendEmailVerificationEmail("public", "userId", &email, &map[string]interface{}{
		"link": "localhost:8080",
	})
	if err != nil {
		t.Error(err.Error())
	}
	assert.True(t, resp.OK != nil)

	assert.Equal(t, EmailVerificationEmailSentForTest, false)
	// assert that link starts with http://localhost:8080. We use starts with because the link
	// can continue a path and random query params too
	assert.Equal(t, link[:21], "http://localhost:8080")

}

// func TestSMTPServiceManually(t *testing.T) {
// 	targetEmail := "..."
// 	fromEmail := "no-reply@supertokens.com"
// 	host := "smtp.gmail.com"
// 	password := "..."
// 	// secure := false
// 	// port := 587
// 	secure := true
// 	port := 465

// 	smtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
// 		Settings: emaildelivery.SMTPSettings{
// 			Host: host,
// 			From: emaildelivery.SMTPFrom{
// 				Name:  "Test User",
// 				Email: fromEmail,
// 			},
// 			Secure:   &secure,
// 			Port:     port,
// 			Password: password,
// 		},
// 	})
// 	configValue := supertokens.TypeInput{
// 		Supertokens: &supertokens.ConnectionInfo{
// 			ConnectionURI: "http://localhost:8080",
// 		},
// 		AppInfo: supertokens.AppInfo{
// 			APIDomain:     "api.supertokens.io",
// 			AppName:       "SuperTokens",
// 			WebsiteDomain: "supertokens.io",
// 		},
// 		RecipeList: []supertokens.Recipe{
// 			Init(evmodels.TypeInput{
// 				EmailDelivery: &emaildelivery.TypeInput{
// 					Service: smtpService,
// 				},
// 				GetEmailForUserID: func(userID string, userContext supertokens.UserContext) (string, error) {
// 					return targetEmail, nil
// 				},
// 			}),
// 		},
// 	}

// 	BeforeEach()
// 	defer AfterEach()
// 	err := supertokens.Init(configValue)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	err = SendEmail(emaildelivery.EmailType{
// 		EmailVerification: &emaildelivery.EmailVerificationType{
// 			User: emaildelivery.User{
// 				ID:    "someId",
// 				Email: targetEmail,
// 			},
// 		},
// 	})

// 	assert.Nil(t, err)
// }
