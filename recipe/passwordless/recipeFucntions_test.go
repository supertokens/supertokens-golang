/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestGetUser(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}

	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	user, err := GetUserByID("random", nil)
	assert.NoError(t, err)
	assert.Nil(t, user)

	result, err := SignInUpByEmail("test@example.com")
	assert.NoError(t, err)

	user = &result.User

	userData, err := GetUserByID(user.ID, nil)
	assert.NoError(t, err)

	assert.Equal(t, user.ID, userData.ID)
	assert.Equal(t, user.Email, userData.Email)
	assert.Nil(t, userData.PhoneNumber)

	user1, err := GetUserByID("random", nil)
	assert.NoError(t, err)
	assert.Nil(t, user1)

	result1, err := SignInUpByEmail("test@example.com")
	assert.NoError(t, err)

	user1 = &result1.User

	userData1, err := GetUserByEmail(*user1.Email)
	assert.NoError(t, err)

	assert.Equal(t, user1.ID, userData1.ID)
	assert.Equal(t, user1.Email, userData1.Email)
	assert.Nil(t, userData1.PhoneNumber)

	user2, err := GetUserByID("random", nil)
	assert.NoError(t, err)
	assert.Nil(t, user2)

	result2, err := SignInUpByPhoneNumber("+1234567890")
	assert.NoError(t, err)

	user2 = &result2.User

	userData2, err := GetUserByPhoneNumber(*user2.PhoneNumber)
	assert.NoError(t, err)

	assert.Equal(t, user2.ID, userData2.ID)
	assert.Equal(t, user2.PhoneNumber, userData2.PhoneNumber)
	assert.Nil(t, userData2.Email)
}

func TestCreateCode(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}

	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	resp, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	assert.NotNil(t, resp.OK.CodeID)
	assert.NotNil(t, resp.OK.CodeLifetime)
	assert.NotNil(t, resp.OK.DeviceID)
	assert.NotNil(t, resp.OK.LinkCode)
	assert.NotNil(t, resp.OK.PreAuthSessionID)
	assert.NotNil(t, resp.OK.TimeCreated)
	assert.NotNil(t, resp.OK.UserInputCode)

	userInputCode := "123"
	resp1, err := CreateCodeWithEmail("test@example.com", &userInputCode)
	assert.NoError(t, err)

	assert.NotNil(t, resp1.OK.CodeID)
	assert.NotNil(t, resp1.OK.CodeLifetime)
	assert.NotNil(t, resp1.OK.DeviceID)
	assert.NotNil(t, resp1.OK.LinkCode)
	assert.NotNil(t, resp1.OK.PreAuthSessionID)
	assert.NotNil(t, resp1.OK.TimeCreated)
	assert.NotNil(t, resp1.OK.UserInputCode)
}

func TestCreateNewCodeForDeviceTest(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}

	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	resp, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	newDeviceCodeResp, err := CreateNewCodeForDevice(resp.OK.DeviceID, nil)
	assert.NoError(t, err)

	assert.NotNil(t, newDeviceCodeResp.OK.CodeID)
	assert.NotNil(t, newDeviceCodeResp.OK.CodeLifetime)
	assert.NotNil(t, newDeviceCodeResp.OK.DeviceID)
	assert.NotNil(t, newDeviceCodeResp.OK.LinkCode)
	assert.NotNil(t, newDeviceCodeResp.OK.PreAuthSessionID)
	assert.NotNil(t, newDeviceCodeResp.OK.TimeCreated)
	assert.NotNil(t, newDeviceCodeResp.OK.UserInputCode)

	resp1, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	userInputCode := "123"
	newDeviceCodeResp1, err := CreateNewCodeForDevice(resp1.OK.DeviceID, &userInputCode)
	assert.NoError(t, err)

	assert.NotNil(t, newDeviceCodeResp1.OK.CodeID)
	assert.NotNil(t, newDeviceCodeResp1.OK.CodeLifetime)
	assert.NotNil(t, newDeviceCodeResp1.OK.DeviceID)
	assert.NotNil(t, newDeviceCodeResp1.OK.LinkCode)
	assert.NotNil(t, newDeviceCodeResp1.OK.PreAuthSessionID)
	assert.NotNil(t, newDeviceCodeResp1.OK.TimeCreated)
	assert.NotNil(t, newDeviceCodeResp1.OK.UserInputCode)

	_, err = CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	newDeviceCodeResp2, err := CreateNewCodeForDevice("asdasdasddas", nil)
	assert.NoError(t, err)

	assert.NotNil(t, newDeviceCodeResp2.RestartFlowError)
	assert.Nil(t, newDeviceCodeResp2.OK)

	resp2, err := CreateCodeWithEmail("test@example.com", &userInputCode)
	assert.NoError(t, err)

	newDeviceCodeResp3, err := CreateNewCodeForDevice(resp2.OK.DeviceID, &userInputCode)
	assert.NoError(t, err)

	assert.NotNil(t, newDeviceCodeResp3.UserInputCodeAlreadyUsedError)
	assert.Nil(t, newDeviceCodeResp3.OK)
}

func TestConsumeCode(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}

	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	codeInfo, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	resp, err := ConsumeCodeWithUserInputCode(codeInfo.OK.DeviceID, codeInfo.OK.UserInputCode, codeInfo.OK.PreAuthSessionID)
	assert.NoError(t, err)

	assert.True(t, resp.OK.CreatedNewUser)
	assert.NotNil(t, resp.OK.User)

	codeInfo1, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	resp1, err := ConsumeCodeWithUserInputCode(codeInfo1.OK.DeviceID, "qefefikjeii", codeInfo1.OK.PreAuthSessionID)
	assert.NoError(t, err)

	assert.NotNil(t, resp1.IncorrectUserInputCodeError)
	assert.Nil(t, resp1.OK)

	codeInfo2, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	_, err = ConsumeCodeWithUserInputCode(codeInfo2.OK.DeviceID, codeInfo2.OK.UserInputCode, "asdasdasdasds")
	assert.Contains(t, err.Error(), "preAuthSessionId and deviceId doesn't match")
}

func TestConsumeCodeWithExpiredUserInputCode(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
			}),
		},
	}
	BeforeEach()
	unittesting.SetKeyValueInConfig("passwordless_code_lifetime", "1000")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	codeInfo, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	resp, err := ConsumeCodeWithUserInputCode(codeInfo.OK.DeviceID, codeInfo.OK.UserInputCode, codeInfo.OK.PreAuthSessionID)
	assert.NoError(t, err)

	assert.NotNil(t, resp.ExpiredUserInputCodeError)
	assert.Equal(t, 1, resp.ExpiredUserInputCodeError.FailedCodeInputAttemptCount)
	assert.Equal(t, 5, resp.ExpiredUserInputCodeError.MaximumCodeInputAttempts)
}

func TestUpdateUserContactMethodEmail(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
			}),
		},
	}
	BeforeEach()
	unittesting.SetKeyValueInConfig("passwordless_code_lifetime", "1000")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	resp, err := SignInUpByEmail("test@example.com")
	assert.NoError(t, err)

	email := "test2@example.com"
	updatedResp, err := UpdateUser(resp.User.ID, &email, nil)
	assert.NoError(t, err)

	assert.NotNil(t, updatedResp.OK)

	updatedUser, err := GetUserByID(resp.User.ID, nil)
	assert.NoError(t, err)

	assert.Equal(t, *updatedUser.Email, email)

	updatedResp, err = UpdateUser("asdasdasdsads", &email, nil)
	assert.NoError(t, err)

	assert.Nil(t, updatedResp.OK)
	assert.NotNil(t, updatedResp.UnknownUserIdError)

	resp1, err := SignInUpByEmail("test3@example.com")
	assert.NoError(t, err)

	updatedResp, err = UpdateUser(resp1.User.ID, &email, nil)
	assert.NoError(t, err)

	assert.Nil(t, updatedResp.OK)
	assert.NotNil(t, updatedResp.EmailAlreadyExistsError)
}

func TestUpdateUserContactMethodPhone(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
			}),
		},
	}
	BeforeEach()
	unittesting.SetKeyValueInConfig("passwordless_code_lifetime", "1000")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	phoneNumber_1 := "+1234567891"
	phoneNumber_2 := "+1234567892"
	phoneNumber_3 := "+1234567893"

	userInfo, err := SignInUpByPhoneNumber(phoneNumber_1)
	assert.NoError(t, err)

	res1, err := UpdateUser(userInfo.User.ID, nil, &phoneNumber_2)
	assert.NoError(t, err)

	assert.NotNil(t, res1.OK)

	result, err := GetUserByID(userInfo.User.ID, nil)
	assert.NoError(t, err)

	assert.Equal(t, phoneNumber_2, *result.PhoneNumber)

	userInfo1, err := SignInUpByPhoneNumber(phoneNumber_3)
	assert.NoError(t, err)

	res1, err = UpdateUser(userInfo1.User.ID, nil, &phoneNumber_2)
	assert.NoError(t, err)

	assert.Nil(t, res1.OK)
	assert.NotNil(t, res1.PhoneNumberAlreadyExistsError)
}

func TestRevokeAllCodes(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
			}),
		},
	}
	BeforeEach()
	unittesting.SetKeyValueInConfig("passwordless_code_lifetime", "1000")
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	codeInfo1, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)
	codeInfo2, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	err = RevokeAllCodesByEmail("test@example.com")
	assert.NoError(t, err)

	result1, err := ConsumeCodeWithUserInputCode(codeInfo1.OK.DeviceID, codeInfo1.OK.UserInputCode, codeInfo1.OK.PreAuthSessionID)
	assert.NoError(t, err)

	assert.NotNil(t, result1.RestartFlowError)
	assert.Nil(t, result1.OK)

	result2, err := ConsumeCodeWithUserInputCode(codeInfo2.OK.DeviceID, codeInfo2.OK.UserInputCode, codeInfo2.OK.PreAuthSessionID)
	assert.NoError(t, err)

	assert.NotNil(t, result2.RestartFlowError)
	assert.Nil(t, result2.OK)
}

func TestRevokeCode(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	codeInfo1, err := CreateCodeWithEmail("random@example.com", nil)
	assert.NoError(t, err)
	codeInfo2, err := CreateCodeWithEmail("random@example.com", nil)
	assert.NoError(t, err)

	err = RevokeCode(codeInfo1.OK.CodeID)
	assert.NoError(t, err)

	result1, err := ConsumeCodeWithUserInputCode(codeInfo1.OK.DeviceID, codeInfo1.OK.UserInputCode, codeInfo1.OK.PreAuthSessionID)
	assert.NoError(t, err)
	assert.NotNil(t, result1.RestartFlowError)
	assert.Nil(t, result1.OK)

	result2, err := ConsumeCodeWithUserInputCode(codeInfo2.OK.DeviceID, codeInfo2.OK.UserInputCode, codeInfo2.OK.PreAuthSessionID)
	assert.NoError(t, err)
	assert.Nil(t, result2.RestartFlowError)
	assert.NotNil(t, result2.OK)
}

func TestListCodesByEmail(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	codeInfo1, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)
	codeInfo2, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	res, err := ListCodesByEmail("test@example.com")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(res))

	for _, dt := range res {
		for _, c := range dt.Codes {
			if !(c.CodeID == codeInfo1.OK.CodeID || c.CodeID == codeInfo2.OK.CodeID) {
				t.Fail()
			}
		}
	}
}

func TestListCodeByPhoneNumber(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	codeInfo1, err := CreateCodeWithPhoneNumber("+1234567890", nil)
	assert.NoError(t, err)
	codeInfo2, err := CreateCodeWithPhoneNumber("+1234567890", nil)
	assert.NoError(t, err)

	res, err := ListCodesByPhoneNumber("+1234567890")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(res))

	for _, dt := range res {
		for _, c := range dt.Codes {
			if !(c.CodeID == codeInfo1.OK.CodeID || c.CodeID == codeInfo2.OK.CodeID) {
				t.Fail()
			}
		}
	}
}

func TestCreatingMagicLink(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	link, err := CreateMagicLinkByPhoneNumber("+1234567890")
	assert.NoError(t, err)

	res, err := url.Parse(link)
	assert.NoError(t, err)

	assert.Equal(t, "supertokens.io", res.Host)
	assert.Equal(t, "/auth/verify", res.Path)
	assert.Equal(t, "passwordless", res.Query().Get("rid"))
}

func TestSignInUp(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	result, err := SignInUpByPhoneNumber("+1234567890")
	assert.NoError(t, err)

	assert.True(t, result.CreatedNewUser)
	assert.NotNil(t, result.User)
	assert.Equal(t, "+1234567890", *result.User.PhoneNumber)
	assert.NotNil(t, result.User.ID)
	assert.NotNil(t, result.User.TimeJoined)
}

func TestListCodesByPreAuthSessionID(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	codeInfo1, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	codeInfo2, err := CreateNewCodeForDevice(codeInfo1.OK.DeviceID, nil)
	assert.NoError(t, err)

	assert.Equal(t, codeInfo1.OK.PreAuthSessionID, codeInfo2.OK.PreAuthSessionID)

	res, err := ListCodesByPreAuthSessionID(codeInfo1.OK.PreAuthSessionID)
	assert.NoError(t, err)

	for _, c := range res.Codes {
		if !(c.CodeID == codeInfo1.OK.CodeID || c.CodeID == codeInfo2.OK.CodeID) {
			t.Fail()
		}
	}
}
