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

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evclaims"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/jwt"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
	"github.com/supertokens/supertokens-golang/recipe/userroles/userrolesclaims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type CustomDevice struct {
	PreAuthSessionID string
	Codes            []CustomCode
}

type CustomCode struct {
	UrlWithLinkCode *string
	UserInputCode   *string
}

func saveCode(_ string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	device, ok := deviceStore[preAuthSessionId]
	if !ok {
		device = CustomDevice{
			PreAuthSessionID: preAuthSessionId,
			Codes:            []CustomCode{},
		}
	}

	codes := device.Codes
	device.Codes = append(codes, CustomCode{
		UrlWithLinkCode: urlWithLinkCode,
		UserInputCode:   userInputCode,
	})
	deviceStore[preAuthSessionId] = device
	return nil
}

var latestURLWithToken string = ""
var apiPort string = "8083"
var webPort string = "3031"
var deviceStore map[string]CustomDevice

func callSTInit(passwordlessConfig *plessmodels.TypeInput) {
	supertokens.ResetForTest()
	emailpassword.ResetForTest()
	emailverification.ResetForTest()
	jwt.ResetForTest()
	passwordless.ResetForTest()
	session.ResetForTest()
	thirdparty.ResetForTest()
	thirdpartyemailpassword.ResetForTest()
	thirdpartypasswordless.ResetForTest()
	userroles.ResetForTest()

	sendPasswordlessLoginSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		return saveCode(input.PasswordlessLogin.PhoneNumber, input.PasswordlessLogin.UserInputCode, input.PasswordlessLogin.UrlWithLinkCode, input.PasswordlessLogin.CodeLifetime, input.PasswordlessLogin.PreAuthSessionId, userContext)
	}

	if passwordlessConfig == nil {
		passwordlessConfig = &plessmodels.TypeInput{
			SmsDelivery: &smsdelivery.TypeInput{
				Service: &smsdelivery.SmsDeliveryInterface{
					SendSms: &sendPasswordlessLoginSms,
				},
			},
			ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
				Enabled: true,
			},
			FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
		}
	}

	countryOptional := true
	formFields := []epmodels.TypeInputFormField{
		{
			ID: "name",
		},
		{
			ID: "age",
			Validate: func(value interface{}, tenantId string) *string {
				age, _ := strconv.Atoi(value.(string))
				if age >= 18 {
					// return nil to indicate success
					return nil
				}
				err := "You must be over 18 to register"
				return &err
			},
		},
		{
			ID:       "country",
			Optional: &countryOptional,
		},
	}
	sendEvEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		latestURLWithToken = input.EmailVerification.EmailVerifyLink
		return nil
	}

	sendPasswordResetEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		latestURLWithToken = input.PasswordReset.PasswordResetLink
		return nil
	}

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:9000",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "localhost:" + apiPort,
			WebsiteDomain: "http://localhost:" + webPort,
		},
		RecipeList: []supertokens.Recipe{
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEvEmail,
					},
				},

				Override: &evmodels.OverrideStruct{
					APIs: func(originalImplementation evmodels.APIInterface) evmodels.APIInterface {
						ogGenerateEmailVerifyTokenPOST := *originalImplementation.GenerateEmailVerifyTokenPOST
						ogVerifyEmailPOST := *originalImplementation.VerifyEmailPOST

						(*originalImplementation.GenerateEmailVerifyTokenPOST) = func(sessionContainer sessmodels.SessionContainer, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.GenerateEmailVerifyTokenPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API email verification code", false)
							if gr != nil {
								return evmodels.GenerateEmailVerifyTokenPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogGenerateEmailVerifyTokenPOST(sessionContainer, options, userContext)
						}

						(*originalImplementation.VerifyEmailPOST) = func(token string, sessionContainer sessmodels.SessionContainer, tenantId string, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.VerifyEmailPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API email verify", false)
							if gr != nil {
								return evmodels.VerifyEmailPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogVerifyEmailPOST(token, sessionContainer, tenantId, options, userContext)
						}
						return originalImplementation
					},
				},
			}),
			emailpassword.Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						ogPasswordResetPOST := *originalImplementation.PasswordResetPOST
						ogGeneratePasswordResetTokenPOST := *originalImplementation.GeneratePasswordResetTokenPOST
						ogEmailExistsGET := *originalImplementation.EmailExistsGET
						ogSignUpPOST := *originalImplementation.SignUpPOST
						ogSignInPOST := *originalImplementation.SignInPOST

						(*originalImplementation.PasswordResetPOST) = func(formFields []epmodels.TypeFormField, token string, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.ResetPasswordPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API reset password consume", false)
							if gr != nil {
								return epmodels.ResetPasswordPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogPasswordResetPOST(formFields, token, tenantId, options, userContext)
						}

						(*originalImplementation.GeneratePasswordResetTokenPOST) = func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.GeneratePasswordResetTokenPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API reset password", false)
							if gr != nil {
								return epmodels.GeneratePasswordResetTokenPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogGeneratePasswordResetTokenPOST(formFields, tenantId, options, userContext)
						}

						(*originalImplementation.EmailExistsGET) = func(email string, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.EmailExistsGETResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API email exists", true)
							if gr != nil {
								return epmodels.EmailExistsGETResponse{
									GeneralError: gr,
								}, nil
							}
							return ogEmailExistsGET(email, tenantId, options, userContext)
						}

						(*originalImplementation.SignUpPOST) = func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API sign up", false)
							if gr != nil {
								return epmodels.SignUpPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogSignUpPOST(formFields, tenantId, options, userContext)
						}

						(*originalImplementation.SignInPOST) = func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignInPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API sign in", false)
							if gr != nil {
								return epmodels.SignInPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogSignInPOST(formFields, tenantId, options, userContext)
						}
						return originalImplementation
					},
				},
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: formFields,
				},
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendPasswordResetEmail,
					},
				},
			}),
			thirdparty.Init(&tpmodels.TypeInput{
				Override: &tpmodels.OverrideStruct{
					APIs: func(originalImplementation tpmodels.APIInterface) tpmodels.APIInterface {
						ogAuthorisationUrlGET := *originalImplementation.AuthorisationUrlGET
						ogSignInUpPOST := *originalImplementation.SignInUpPOST

						(*originalImplementation.AuthorisationUrlGET) = func(provider *tpmodels.TypeProvider, redirectURIOnProviderDashboard string, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API authorisation url get", true)
							if gr != nil {
								return tpmodels.AuthorisationUrlGETResponse{
									GeneralError: gr,
								}, nil
							}
							return ogAuthorisationUrlGET(provider, redirectURIOnProviderDashboard, tenantId, options, userContext)
						}

						(*originalImplementation.SignInUpPOST) = func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.SignInUpPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API sign in up", false)
							if gr != nil {
								return tpmodels.SignInUpPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogSignInUpPOST(provider, input, tenantId, options, userContext)
						}
						return originalImplementation
					},
				},
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
										ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
									},
								},
							},
						},
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "github",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
										ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
									},
								},
							},
						},
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "facebook",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
										ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
									},
								},
							},
						},
						customAuth0Provider(),
					},
				},
			}),
			thirdpartyemailpassword.Init(&tpepmodels.TypeInput{
				Override: &tpepmodels.OverrideStruct{
					APIs: func(originalImplementation tpepmodels.APIInterface) tpepmodels.APIInterface {
						ogPasswordResetPOST := *originalImplementation.PasswordResetPOST
						ogGeneratePasswordResetTokenPOST := *originalImplementation.GeneratePasswordResetTokenPOST
						ogEmailExistsGET := *originalImplementation.EmailPasswordEmailExistsGET
						ogSignUpPOST := *originalImplementation.EmailPasswordSignUpPOST
						ogSignInPOST := *originalImplementation.EmailPasswordSignInPOST
						ogAuthorisationUrlGET := *originalImplementation.AuthorisationUrlGET
						ogSignInUpPOST := *originalImplementation.ThirdPartySignInUpPOST

						(*originalImplementation.AuthorisationUrlGET) = func(provider *tpmodels.TypeProvider, redirectURIOnProviderDashboard string, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API authorisation url get", true)
							if gr != nil {
								return tpmodels.AuthorisationUrlGETResponse{
									GeneralError: gr,
								}, nil
							}
							return ogAuthorisationUrlGET(provider, redirectURIOnProviderDashboard, tenantId, options, userContext)
						}

						(*originalImplementation.ThirdPartySignInUpPOST) = func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpepmodels.ThirdPartySignInUpPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API sign in up", false)
							if gr != nil {
								return tpepmodels.ThirdPartySignInUpPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogSignInUpPOST(provider, input, tenantId, options, userContext)
						}

						(*originalImplementation.PasswordResetPOST) = func(formFields []epmodels.TypeFormField, token string, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.ResetPasswordPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API reset password consume", false)
							if gr != nil {
								return epmodels.ResetPasswordPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogPasswordResetPOST(formFields, token, tenantId, options, userContext)
						}

						(*originalImplementation.GeneratePasswordResetTokenPOST) = func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.GeneratePasswordResetTokenPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API reset password", false)
							if gr != nil {
								return epmodels.GeneratePasswordResetTokenPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogGeneratePasswordResetTokenPOST(formFields, tenantId, options, userContext)
						}

						(*originalImplementation.EmailPasswordEmailExistsGET) = func(email string, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.EmailExistsGETResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API email exists", true)
							if gr != nil {
								return epmodels.EmailExistsGETResponse{
									GeneralError: gr,
								}, nil
							}
							return ogEmailExistsGET(email, tenantId, options, userContext)
						}

						(*originalImplementation.EmailPasswordSignUpPOST) = func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (tpepmodels.SignUpPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API sign up", false)
							if gr != nil {
								return tpepmodels.SignUpPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogSignUpPOST(formFields, tenantId, options, userContext)
						}

						(*originalImplementation.EmailPasswordSignInPOST) = func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (tpepmodels.SignInPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API sign in", false)
							if gr != nil {
								return tpepmodels.SignInPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogSignInPOST(formFields, tenantId, options, userContext)
						}
						return originalImplementation
					},
				},
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: formFields,
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
									ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
								},
							},
						},
					},
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "github",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
									ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
								},
							},
						},
					},
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "facebook",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
									ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
								},
							},
						},
					},
					customAuth0Provider(),
				},
			}),
			session.Init(&sessmodels.TypeInput{
				Override: &sessmodels.OverrideStruct{
					APIs: func(originalImplementation sessmodels.APIInterface) sessmodels.APIInterface {
						ogSignOutPOST := *originalImplementation.SignOutPOST
						(*originalImplementation.SignOutPOST) = func(sessionContainer sessmodels.SessionContainer, options sessmodels.APIOptions, userContext supertokens.UserContext) (sessmodels.SignOutPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from signout API", false)
							if gr != nil {
								return sessmodels.SignOutPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogSignOutPOST(sessionContainer, options, userContext)
						}
						return originalImplementation
					},
				},
			}),
			passwordless.Init(plessmodels.TypeInput{
				ContactMethodPhone:        passwordlessConfig.ContactMethodPhone,
				ContactMethodEmail:        passwordlessConfig.ContactMethodEmail,
				ContactMethodEmailOrPhone: passwordlessConfig.ContactMethodEmailOrPhone,
				FlowType:                  passwordlessConfig.FlowType,
				GetCustomUserInputCode:    passwordlessConfig.GetCustomUserInputCode,
				Override: &plessmodels.OverrideStruct{
					APIs: func(originalImplementation plessmodels.APIInterface) plessmodels.APIInterface {
						ogConsumeCodePOST := *originalImplementation.ConsumeCodePOST
						ogCreateCodePOST := *originalImplementation.CreateCodePOST
						ogResendCodePOST := *originalImplementation.ResendCodePOST

						(*originalImplementation.ConsumeCodePOST) = func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, tenantId string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.ConsumeCodePOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API consume code", false)
							if gr != nil {
								return plessmodels.ConsumeCodePOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogConsumeCodePOST(userInput, linkCode, preAuthSessionID, tenantId, options, userContext)
						}

						(*originalImplementation.CreateCodePOST) = func(email, phoneNumber *string, tenantId string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.CreateCodePOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API create code", false)
							if gr != nil {
								return plessmodels.CreateCodePOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogCreateCodePOST(email, phoneNumber, tenantId, options, userContext)
						}

						(*originalImplementation.ResendCodePOST) = func(deviceID, preAuthSessionID string, tenantId string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.ResendCodePOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API resend code", false)
							if gr != nil {
								return plessmodels.ResendCodePOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogResendCodePOST(deviceID, preAuthSessionID, tenantId, options, userContext)
						}
						return originalImplementation
					},
				},
			}),
			thirdpartypasswordless.Init(tplmodels.TypeInput{
				ContactMethodPhone:        passwordlessConfig.ContactMethodPhone,
				ContactMethodEmail:        passwordlessConfig.ContactMethodEmail,
				ContactMethodEmailOrPhone: passwordlessConfig.ContactMethodEmailOrPhone,
				FlowType:                  passwordlessConfig.FlowType,
				GetCustomUserInputCode:    passwordlessConfig.GetCustomUserInputCode,
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
									ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
								},
							},
						},
					},
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "github",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
									ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
								},
							},
						},
					},
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "facebook",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
									ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
								},
							},
						},
					},
					customAuth0Provider(),
				},
				Override: &tplmodels.OverrideStruct{
					APIs: func(originalImplementation tplmodels.APIInterface) tplmodels.APIInterface {
						ogConsumeCodePOST := *originalImplementation.ConsumeCodePOST
						ogCreateCodePOST := *originalImplementation.CreateCodePOST
						ogResendCodePOST := *originalImplementation.ResendCodePOST
						ogAuthorisationUrlGET := *originalImplementation.AuthorisationUrlGET
						ogSignInUpPOST := *originalImplementation.ThirdPartySignInUpPOST

						(*originalImplementation.AuthorisationUrlGET) = func(provider *tpmodels.TypeProvider, redirectURIOnProviderDashboard string, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API authorisation url get", true)
							if gr != nil {
								return tpmodels.AuthorisationUrlGETResponse{
									GeneralError: gr,
								}, nil
							}
							return ogAuthorisationUrlGET(provider, redirectURIOnProviderDashboard, tenantId, options, userContext)
						}

						(*originalImplementation.ThirdPartySignInUpPOST) = func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tplmodels.ThirdPartySignInUpPOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API sign in up", false)
							if gr != nil {
								return tplmodels.ThirdPartySignInUpPOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogSignInUpPOST(provider, input, tenantId, options, userContext)
						}

						(*originalImplementation.ConsumeCodePOST) = func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, tenantId string, options plessmodels.APIOptions, userContext supertokens.UserContext) (tplmodels.ConsumeCodePOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API consume code", false)
							if gr != nil {
								return tplmodels.ConsumeCodePOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogConsumeCodePOST(userInput, linkCode, preAuthSessionID, tenantId, options, userContext)
						}

						(*originalImplementation.CreateCodePOST) = func(email, phoneNumber *string, tenantId string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.CreateCodePOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API create code", false)
							if gr != nil {
								return plessmodels.CreateCodePOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogCreateCodePOST(email, phoneNumber, tenantId, options, userContext)
						}

						(*originalImplementation.ResendCodePOST) = func(deviceID, preAuthSessionID string, tenantId string, options plessmodels.APIOptions, userContext supertokens.UserContext) (plessmodels.ResendCodePOSTResponse, error) {
							gr := returnGeneralErrorIfNeeded(*options.Req, "general error from API resend code", false)
							if gr != nil {
								return plessmodels.ResendCodePOSTResponse{
									GeneralError: gr,
								}, nil
							}
							return ogResendCodePOST(deviceID, preAuthSessionID, tenantId, options, userContext)
						}
						return originalImplementation
					},
				},
			}),
			userroles.Init(nil),
		},
	})

	if err != nil {
		panic(err.Error())
	}

	middleware := supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/sessionInfo" && r.Method == "GET" {
			session.VerifySession(nil, sessioninfo).ServeHTTP(rw, r)
		} else if r.URL.Path == "/token" && r.Method == "GET" {
			rw.WriteHeader(200)
			rw.Header().Add("content-type", "application/json")
			bytes, _ := json.Marshal(map[string]interface{}{
				"latestURLWithToken": latestURLWithToken,
			})
			rw.Write(bytes)
		} else if r.URL.Path == "/beforeeach" && r.Method == "POST" {
			deviceStore = map[string]CustomDevice{}
			rw.WriteHeader(200)
			rw.Header().Add("content-type", "application/json")
			bytes, _ := json.Marshal(map[string]interface{}{})
			rw.Write(bytes)
		} else if r.URL.Path == "/test/setFlow" && r.Method == "POST" {
			reInitST(rw, r)
		} else if r.URL.Path == "/test/getDevice" && r.Method == "GET" {
			getDevice(rw, r)
		} else if r.URL.Path == "/test/featureFlags" && r.Method == "GET" {
			rw.WriteHeader(200)
			rw.Header().Add("content-type", "application/json")
			bytes, _ := json.Marshal(map[string]interface{}{
				"available": []string{"passwordless", "thirdpartypasswordless", "generalerror", "userroles"},
			})
			rw.Write(bytes)

		} else if r.URL.Path == "/unverifyEmail" && r.Method == "GET" {
			session.VerifySession(nil, func(w http.ResponseWriter, r *http.Request) {
				sessionContainer := session.GetSessionFromRequestContext(r.Context())
				emailverification.UnverifyEmail(sessionContainer.GetUserID(), nil)
				sessionContainer.FetchAndSetClaim(evclaims.EmailVerificationClaim)
				rw.Header().Add("content-type", "application/json")
				rw.WriteHeader(200)
				rw.Write([]byte("{\"status\": \"OK\"}"))
			}).ServeHTTP(rw, r)

		} else if r.URL.Path == "/setRole" && r.Method == "POST" {
			session.VerifySession(nil, func(w http.ResponseWriter, r *http.Request) {
				sessionContainer := session.GetSessionFromRequestContext(r.Context())
				bodyBytes, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return
				}
				var body map[string]interface{}
				err = json.Unmarshal(bodyBytes, &body)
				if err != nil {
					return
				}
				role := body["role"].(string)
				permissions := body["permissions"].([]interface{})
				permissionsStr := make([]string, len(permissions))
				for i, p := range permissions {
					permissionsStr[i] = p.(string)
				}
				_, err = userroles.CreateNewRoleOrAddPermissions(role, permissionsStr, &map[string]interface{}{})
				if err != nil {
					return
				}
				_, err = userroles.AddRoleToUser("public", sessionContainer.GetUserID(), role, &map[string]interface{}{})
				if err != nil {
					return
				}
				err = sessionContainer.FetchAndSetClaim(userrolesclaims.UserRoleClaim)
				if err != nil {
					return
				}
				err = sessionContainer.FetchAndSetClaim(userrolesclaims.PermissionClaim)
				if err != nil {
					return
				}
				rw.Header().Add("content-type", "application/json")
				rw.WriteHeader(200)
				rw.Write([]byte("{\"status\": \"OK\"}"))
			}).ServeHTTP(rw, r)

		} else if r.URL.Path == "/checkRole" && r.Method == "POST" {
			session.VerifySession(&sessmodels.VerifySessionOptions{
				OverrideGlobalClaimValidators: func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
					req := (*userContext)["_default"].(map[string]interface{})["request"].(*http.Request)
					bodyBytes, err := ioutil.ReadAll(req.Body)
					if err != nil {
						return nil, err
					}
					var body map[string]interface{}
					err = json.Unmarshal(bodyBytes, &body)
					if err != nil {
						return nil, err
					}

					getValidator := func(validator claims.PrimitiveArrayClaimValidators, validatorStr string, args []interface{}) claims.SessionClaimValidator {
						var maxAge *int64 = nil
						var id *string = nil
						if len(args) > 1 {
							maxAgeFloat := args[1].(float64)
							maxAgeInt := int64(maxAgeFloat)
							maxAge = &maxAgeInt
						}
						if len(args) > 2 {
							idStr := args[2].(string)
							id = &idStr
						}

						switch validatorStr {
						case "includes":
							return validator.Includes(args[0].(string), maxAge, id)
						case "excludes":
							return validator.Excludes(args[0].(string), maxAge, id)
						case "includesAll":
							return validator.IncludesAll(args[0].([]interface{}), maxAge, id)
						case "excludesAll":
							return validator.ExcludesAll(args[0].([]interface{}), maxAge, id)
						}

						return claims.SessionClaimValidator{}
					}

					if role, ok := body["role"].(map[string]interface{}); ok {
						validatorStr := role["validator"].(string)
						args := role["args"].([]interface{})
						globalClaimValidators = append(globalClaimValidators, getValidator(userrolesclaims.UserRoleClaimValidators, validatorStr, args))
					}

					if permission, ok := body["permission"].(map[string]interface{}); ok {
						validatorStr := permission["validator"].(string)
						args := permission["args"].([]interface{})
						globalClaimValidators = append(globalClaimValidators, getValidator(userrolesclaims.PermissionClaimValidators, validatorStr, args))
					}

					return globalClaimValidators, nil
				},
			}, func(w http.ResponseWriter, r *http.Request) {
				rw.Header().Add("content-type", "application/json")
				rw.WriteHeader(200)
				rw.Write([]byte("{\"status\": \"OK\"}"))
			}).ServeHTTP(rw, r)
		}
	}))

	routes = &middleware
}

func returnGeneralErrorIfNeeded(req http.Request, message string, useQueryParams bool) *supertokens.GeneralErrorResponse {
	if useQueryParams {
		generalError := req.URL.Query().Get("generalError")
		if generalError == "" {
			return nil
		}
		return &supertokens.GeneralErrorResponse{
			Message: message,
		}
	} else {
		var body map[string]interface{}
		_ = json.NewDecoder(req.Body).Decode(&body)
		_, ok := body["generalError"]
		customMessage, okCustomMessage := body["generalErrorMessage"]
		if ok {
			if okCustomMessage {
				return &supertokens.GeneralErrorResponse{
					Message: customMessage.(string),
				}
			}
			return &supertokens.GeneralErrorResponse{
				Message: message,
			}
		}
	}
	return nil
}

func customAuth0Provider() tpmodels.ProviderInput {

	var providerInput tpmodels.ProviderInput

	providerInput.Config.ThirdPartyId = "auth0"
	providerInput.Config.Clients = []tpmodels.ProviderClientConfig{
		{
			ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		},
	}
	providerInput.Config.AuthorizationEndpoint = "https://" + os.Getenv("AUTH0_DOMAIN") + "/authorize"
	providerInput.Config.TokenEndpoint = "https://" + os.Getenv("AUTH0_DOMAIN") + "/oauth/token"

	providerInput.Override = func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {

			authCodeResponseJson, err := json.Marshal(oAuthTokens)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}

			var accessTokenAPIResponse auth0GetProfileInfoInput
			err = json.Unmarshal(authCodeResponseJson, &accessTokenAPIResponse)

			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}

			accessToken := accessTokenAPIResponse.AccessToken
			authHeader := "Bearer " + accessToken

			// TODO maybe userInfoEndpoint is enough. Verify while testing
			response, err := getAuth0AuthRequest(authHeader)

			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}

			userInfo := response.(map[string]interface{})

			ID := userInfo["sub"].(string)
			email := userInfo["name"].(string)

			return tpmodels.TypeUserInfo{
				ThirdPartyUserId: ID,
				Email: &tpmodels.EmailStruct{
					ID:         email,
					IsVerified: true, // true if email is verified already
				},
			}, nil
		}

		return originalImplementation
	}
	return providerInput
}

func getAuth0AuthRequest(authHeader string) (interface{}, error) {
	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/userinfo"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	return doGetRequest(req)
}

func doGetRequest(req *http.Request) (interface{}, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type auth0GetProfileInfoInput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:"+webPort)
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			response.Header().Set("Access-Control-Allow-Headers", strings.Join(append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...), ","))
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.WriteHeader(204)
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
	})
}

var routes *http.Handler

func main() {
	deviceStore = map[string]CustomDevice{}
	godotenv.Load()
	if len(os.Args) >= 2 {
		apiPort = os.Args[1]
	}
	if len(os.Args) >= 3 {
		webPort = os.Args[2]
	}
	supertokens.IsTestFlag = true
	callSTInit(nil)

	http.ListenAndServe("0.0.0.0:"+apiPort, corsMiddleware(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			(*routes).ServeHTTP(rw, r)
		})))
}

func reInitST(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var readBody map[string]interface{}
	json.Unmarshal(body, &readBody)
	sendPasswordlessLoginEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		return saveCode(input.PasswordlessLogin.Email, input.PasswordlessLogin.UserInputCode, input.PasswordlessLogin.UrlWithLinkCode, input.PasswordlessLogin.CodeLifetime, input.PasswordlessLogin.PreAuthSessionId, userContext)
	}

	sendPasswordlessLoginSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		return saveCode(input.PasswordlessLogin.PhoneNumber, input.PasswordlessLogin.UserInputCode, input.PasswordlessLogin.UrlWithLinkCode, input.PasswordlessLogin.CodeLifetime, input.PasswordlessLogin.PreAuthSessionId, userContext)
	}

	config := &plessmodels.TypeInput{
		FlowType: readBody["flowType"].(string),
		EmailDelivery: &emaildelivery.TypeInput{
			Service: &emaildelivery.EmailDeliveryInterface{
				SendEmail: &sendPasswordlessLoginEmail,
			},
		},
		SmsDelivery: &smsdelivery.TypeInput{
			Service: &smsdelivery.SmsDeliveryInterface{
				SendSms: &sendPasswordlessLoginSms,
			},
		},
	}

	if readBody["contactMethod"].(string) == "PHONE" {
		config.ContactMethodPhone = plessmodels.ContactMethodPhoneConfig{
			Enabled: true,
		}
	} else if readBody["contactMethod"].(string) == "EMAIL" {
		config.ContactMethodEmail = plessmodels.ContactMethodEmailConfig{
			Enabled: true,
		}
	} else {
		config.ContactMethodEmailOrPhone = plessmodels.ContactMethodEmailOrPhoneConfig{
			Enabled: true,
		}
	}
	callSTInit(config)
	w.WriteHeader(200)
	w.Write([]byte("success"))
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	preAuthSessionId := r.URL.Query().Get("preAuthSessionId")
	device, ok := deviceStore[preAuthSessionId]
	if ok {
		w.WriteHeader(200)
		w.Header().Add("content-type", "application/json")
		codes := []map[string]interface{}{}
		for _, code := range device.Codes {
			codes = append(codes, map[string]interface{}{
				"urlWithLinkCode": code.UrlWithLinkCode,
				"userInputCode":   code.UserInputCode,
			})
		}
		result := map[string]interface{}{
			"preAuthSessionId": device.PreAuthSessionID,
			"codes":            codes,
		}
		bytes, _ := json.Marshal(result)
		w.Write(bytes)
	} else {
		w.WriteHeader(200)
		w.Write([]byte(""))
	}
}

func sessioninfo(w http.ResponseWriter, r *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())

	if sessionContainer == nil {
		w.WriteHeader(500)
		w.Write([]byte("no session found"))
		return
	}
	sessionData, err := sessionContainer.GetSessionDataInDatabase()
	if err != nil {
		err = supertokens.ErrorHandler(err, r, w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		return
	}
	w.WriteHeader(200)
	w.Header().Add("content-type", "application/json")
	bytes, err := json.Marshal(map[string]interface{}{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error in converting to json"))
	} else {
		w.Write(bytes)
	}
}
