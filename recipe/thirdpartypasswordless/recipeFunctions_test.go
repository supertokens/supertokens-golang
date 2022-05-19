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

package thirdpartypasswordless

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestWithThirdPartyPasswordlessForThirdPartyUserThatIsEmailVerifiedReturnsTheCorrectEmailVerificationStatus(t *testing.T) {
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
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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

	resp, err := ThirdPartySignInUp("customProvider", "verifiedUser", tplmodels.EmailStruct{
		ID:         "test@example.com",
		IsVerified: true,
	})
	assert.NoError(t, err)

	emailVerificationToken, err := CreateEmailVerificationToken(resp.OK.User.ID)
	assert.NoError(t, err)

	VerifyEmailUsingToken(emailVerificationToken.OK.Token)

	isVerfied, err := IsEmailVerified(resp.OK.User.ID)
	assert.NoError(t, err)
	assert.True(t, isVerfied)

	resp1, err := ThirdPartySignInUp("customProvider2", "NotVerifiedUser", tplmodels.EmailStruct{
		ID:         "test@example.com",
		IsVerified: false,
	})
	assert.NoError(t, err)

	isVerfied1, err := IsEmailVerified(resp1.OK.User.ID)
	assert.NoError(t, err)
	assert.False(t, isVerfied1)
}

func TestWithThirdPartyPasswordlessForPasswordlessUserThatIsEmailVerifiedReturnsTrueForBothEmailAndPhone(t *testing.T) {
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
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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

	response, err := PasswordlessSignInUpByEmail("test@example.com")
	assert.NoError(t, err)

	isVerified, err := IsEmailVerified(response.User.ID)
	assert.NoError(t, err)
	assert.True(t, isVerified)

	emailVerificationResp, err := CreateEmailVerificationToken(response.User.ID)
	assert.NoError(t, err)
	assert.NotNil(t, emailVerificationResp.EmailAlreadyVerifiedError)
	assert.Nil(t, emailVerificationResp.OK)

	response, err = PasswordlessSignInUpByPhoneNumber("+123456789012")
	assert.NoError(t, err)

	isVerified, err = IsEmailVerified(response.User.ID)
	assert.NoError(t, err)
	assert.True(t, isVerified)

	emailVerificationResp, err = CreateEmailVerificationToken(response.User.ID)
	assert.NoError(t, err)
	assert.NotNil(t, emailVerificationResp.EmailAlreadyVerifiedError)
	assert.Nil(t, emailVerificationResp.OK)
}

func TestWithThirdPartyPasswordlessGetUserFunctionality(t *testing.T) {
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
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
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

	user, err := GetUserByID("random")
	assert.NoError(t, err)
	assert.Nil(t, user)

	resp, err := PasswordlessSignInUpByEmail("test@example.com")
	assert.NoError(t, err)
	userId := resp.User.ID

	user, err = GetUserByID(userId)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, userId, user.ID)
	assert.Equal(t, resp.User.Email, user.Email)
	assert.Nil(t, user.PhoneNumber)

	users, err := GetUsersByEmail("random")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(users))

	users, err = GetUsersByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))

	userInfo := users[0]

	assert.Equal(t, user.Email, userInfo.Email)
	assert.Equal(t, user.ID, userInfo.ID)
	assert.Equal(t, user.PhoneNumber, userInfo.PhoneNumber)
	assert.Nil(t, userInfo.PhoneNumber)
	assert.Nil(t, userInfo.ThirdParty)
	assert.Equal(t, user.TimeJoined, userInfo.TimeJoined)

	user, err = GetUserByPhoneNumber("random")
	assert.NoError(t, err)
	assert.Nil(t, user)

	resp, err = PasswordlessSignInUpByPhoneNumber("+1234567890")
	assert.NoError(t, err)

	user, err = GetUserByPhoneNumber(*resp.User.PhoneNumber)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, user.Email, resp.User.Email)
	assert.Equal(t, user.ID, resp.User.ID)
	assert.Equal(t, user.PhoneNumber, resp.User.PhoneNumber)
	assert.Equal(t, user.ThirdParty, resp.User.ThirdParty)
	assert.Equal(t, user.ThirdParty, resp.User.ThirdParty)
	assert.Nil(t, user.Email)
	assert.Nil(t, user.ThirdParty)
}

func TestWithThirdPartyPasswordlessCreateCodeTest(t *testing.T) {
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
	assert.NotNil(t, resp.OK)
	assert.NotNil(t, resp.OK.CodeID)
	assert.NotNil(t, resp.OK.PreAuthSessionID)
	assert.NotNil(t, resp.OK.CodeLifetime)
	assert.NotNil(t, resp.OK.DeviceID)
	assert.NotNil(t, resp.OK.LinkCode)
	assert.NotNil(t, resp.OK.TimeCreated)
	assert.NotNil(t, resp.OK.UserInputCode)

	userInputCode := "23123"
	resp1, err := CreateCodeWithEmail("test@example.com", &userInputCode)
	assert.NoError(t, err)
	assert.NotNil(t, resp1.OK)
	assert.NotNil(t, resp1.OK.CodeID)
	assert.NotNil(t, resp1.OK.PreAuthSessionID)
	assert.NotNil(t, resp1.OK.CodeLifetime)
	assert.NotNil(t, resp1.OK.DeviceID)
	assert.NotNil(t, resp1.OK.LinkCode)
	assert.NotNil(t, resp1.OK.TimeCreated)
	assert.NotNil(t, resp1.OK.UserInputCode)
}

func TestThirdPartyPasswordlessCreateNewCodeFromDevice(t *testing.T) {
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

	resp1, err := CreateNewCodeForDevice(resp.OK.DeviceID, nil)
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.NotNil(t, resp1.OK)
	assert.NotNil(t, resp1.OK.CodeID)
	assert.NotNil(t, resp1.OK.PreAuthSessionID)
	assert.NotNil(t, resp1.OK.CodeLifetime)
	assert.NotNil(t, resp1.OK.DeviceID)
	assert.NotNil(t, resp1.OK.LinkCode)
	assert.NotNil(t, resp1.OK.TimeCreated)
	assert.NotNil(t, resp1.OK.UserInputCode)

	resp, err = CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	userInputCode := "2314"
	resp1, err = CreateNewCodeForDevice(resp.OK.DeviceID, &userInputCode)
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.NotNil(t, resp1.OK)
	assert.NotNil(t, resp1.OK.CodeID)
	assert.NotNil(t, resp1.OK.PreAuthSessionID)
	assert.NotNil(t, resp1.OK.CodeLifetime)
	assert.NotNil(t, resp1.OK.DeviceID)
	assert.NotNil(t, resp1.OK.LinkCode)
	assert.NotNil(t, resp1.OK.TimeCreated)
	assert.NotNil(t, resp1.OK.UserInputCode)

	resp, err = CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	resp1, err = CreateNewCodeForDevice("random", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp1.RestartFlowError)
	assert.Nil(t, resp1.OK)

	resp, err = CreateCodeWithEmail("test@example.com", &userInputCode)
	assert.NoError(t, err)

	resp1, err = CreateNewCodeForDevice(resp.OK.DeviceID, &userInputCode)
	assert.NoError(t, err)
	assert.NotNil(t, resp1.UserInputCodeAlreadyUsedError)
	assert.Nil(t, resp1.OK)
}

func TestThirdPartyPasswordlessConsumeCoed(t *testing.T) {
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

	result, err := ConsumeCodeWithUserInputCode(resp.OK.DeviceID, resp.OK.UserInputCode, resp.OK.PreAuthSessionID)
	assert.NoError(t, err)

	assert.NotNil(t, result.OK)
	assert.True(t, result.OK.CreatedNewUser)
	assert.Equal(t, "test@example.com", *result.OK.User.Email)
	assert.NotNil(t, result.OK)
	assert.NotNil(t, result.OK.User)
	assert.NotNil(t, result.OK.User.ID)
	assert.NotNil(t, result.OK.User.TimeJoined)
	assert.Nil(t, result.OK.User.PhoneNumber)
	assert.Nil(t, result.OK.User.ThirdParty)

	resp, err = CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	result, err = ConsumeCodeWithUserInputCode(resp.OK.DeviceID, "random", resp.OK.PreAuthSessionID)
	assert.NoError(t, err)
	assert.NotNil(t, result.IncorrectUserInputCodeError)
	assert.Equal(t, 1, result.IncorrectUserInputCodeError.FailedCodeInputAttemptCount)
	assert.Equal(t, 5, result.IncorrectUserInputCodeError.MaximumCodeInputAttempts)

	resp, err = CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	_, err = ConsumeCodeWithUserInputCode(resp.OK.DeviceID, resp.OK.UserInputCode, "random")
	assert.Contains(t, err.Error(), "preAuthSessionId and deviceId doesn't match")
}

func TestThirdPartyPasswordlessConsumeCodeTestWithExpiredUserInputCodeError(t *testing.T) {
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

	resp, err := CreateCodeWithEmail("test@example.com", nil)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	result, err := ConsumeCodeWithUserInputCode(resp.OK.DeviceID, resp.OK.UserInputCode, resp.OK.PreAuthSessionID)
	assert.NoError(t, err)
	assert.NotNil(t, result.ExpiredUserInputCodeError)
	assert.Equal(t, 1, result.ExpiredUserInputCodeError.FailedCodeInputAttemptCount)
	assert.Equal(t, 5, result.ExpiredUserInputCodeError.MaximumCodeInputAttempts)
	assert.Nil(t, result.OK)
}

func TestThirdPartyPasswordlessUpdateUserContactMethodEmailTest(t *testing.T) {
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

	userInfo, err := PasswordlessSignInUpByEmail("test@example.com")
	assert.NoError(t, err)

	updatedEmail := "test2@example.com"
	resp, err := UpdatePasswordlessUser(userInfo.User.ID, &updatedEmail, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp.OK)

	user, err := GetUserByID(userInfo.User.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedEmail, *user.Email)

	resp, err = UpdatePasswordlessUser("random", &updatedEmail, nil)
	assert.NoError(t, err)
	assert.Nil(t, resp.OK)
	assert.NotNil(t, resp.UnknownUserIdError)

	userInfo2, err := PasswordlessSignInUpByEmail("test3@example.com")
	assert.NoError(t, err)

	resp, err = UpdatePasswordlessUser(userInfo2.User.ID, &updatedEmail, nil)
	assert.NoError(t, err)
	assert.Nil(t, resp.OK)
	assert.NotNil(t, resp.EmailAlreadyExistsError)
}

func TestThirdPartyPasswordlessUpdateUserContactPhone(t *testing.T) {
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

	phoneNumber_1 := "+1234567891"
	phoneNumber_2 := "+1234567892"
	phoneNumber_3 := "+1234567893"

	userInfo, err := PasswordlessSignInUpByPhoneNumber(phoneNumber_1)
	assert.NoError(t, err)

	res1, err := UpdatePasswordlessUser(userInfo.User.ID, nil, &phoneNumber_2)
	assert.NoError(t, err)

	assert.NotNil(t, res1.OK)

	result, err := GetUserByID(userInfo.User.ID)
	assert.NoError(t, err)

	assert.Equal(t, phoneNumber_2, *result.PhoneNumber)

	userInfo1, err := PasswordlessSignInUpByPhoneNumber(phoneNumber_3)
	assert.NoError(t, err)

	res1, err = UpdatePasswordlessUser(userInfo1.User.ID, nil, &phoneNumber_2)
	assert.NoError(t, err)

	assert.Nil(t, res1.OK)
	assert.NotNil(t, res1.PhoneNumberAlreadyExistsError)
}

func TestThirdPartyPasswordlessRevokeAllCodesTest(t *testing.T) {
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

func TestThirdPartyPasswordlessRevokeCode(t *testing.T) {
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

func TestThirdPartyPasswordlessListCodesByEmail(t *testing.T) {
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

	result, err := ListCodesByEmail("test@example.com")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(result))
	for _, dt := range result {
		for _, c := range dt.Codes {
			if !(c.CodeID == codeInfo1.OK.CodeID || c.CodeID == codeInfo2.OK.CodeID) {
				t.Fail()
			}
		}
	}
}

func TestListCodesByPhoneNumber(t *testing.T) {
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

	codeInfo1, err := CreateCodeWithEmail("+1234567890", nil)
	assert.NoError(t, err)

	codeInfo2, err := CreateCodeWithEmail("+1234567890", nil)
	assert.NoError(t, err)

	result, err := ListCodesByEmail("+1234567890")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(result))
	for _, dt := range result {
		for _, c := range dt.Codes {
			if !(c.CodeID == codeInfo1.OK.CodeID || c.CodeID == codeInfo2.OK.CodeID) {
				t.Fail()
			}
		}
	}
}

func TestThirdPartyPasswordlessListCodesByDeviceIdAndListCodesByPreAuthSessionId(t *testing.T) {
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

	codeInfo1, err := CreateCodeWithEmail("+1234567890", nil)
	assert.NoError(t, err)

	result, err := ListCodesByDeviceID(codeInfo1.OK.DeviceID)
	assert.NoError(t, err)

	assert.Equal(t, codeInfo1.OK.CodeID, result.Codes[0].CodeID)

	result, err = ListCodesByPreAuthSessionID(codeInfo1.OK.PreAuthSessionID)
	assert.NoError(t, err)

	assert.Equal(t, codeInfo1.OK.CodeID, result.Codes[0].CodeID)
}

func TestCreateMagicLinkTest(t *testing.T) {
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

	result, err := CreateMagicLinkByPhoneNumber("+1234567890")
	assert.NoError(t, err)

	magicLinkURL, err := url.Parse(result)
	assert.NoError(t, err)

	assert.Equal(t, "supertokens.io", magicLinkURL.Host)
	assert.Equal(t, "/auth/verify", magicLinkURL.Path)
	assert.Equal(t, "thirdpartypasswordless", magicLinkURL.Query().Get("rid"))
	assert.NotNil(t, magicLinkURL.Query().Get("preAuthSessionId"))
}

func TestThirdPartyPasswordlessSignInUp(t *testing.T) {
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

	result, err := PasswordlessSignInUpByPhoneNumber("+12345678901")
	assert.NoError(t, err)
	assert.NotNil(t, result.User)
	assert.True(t, result.CreatedNewUser)
	assert.Equal(t, "+12345678901", *result.User.PhoneNumber)
	assert.NotNil(t, result.User.TimeJoined)
	assert.NotNil(t, result.User.ID)
	assert.Nil(t, result.User.ThirdParty)
}
