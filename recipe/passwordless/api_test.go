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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestCreateCodeAPIWithRidAsThirdpartypasswordless(t *testing.T) {
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
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     "4398792-test-id",
										ClientSecret: "test-secret",
									},
								},
							},
						},
					},
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	validEmail := map[string]interface{}{
		"email": "test@example.com",
	}

	validEmailBody, err := json.Marshal(validEmail)
	if err != nil {
		t.Error(err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", testServer.URL+"/auth/signinup/code", bytes.NewBuffer(validEmailBody))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "thirdpartypasswordless")

	validEmailResp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validEmailResp.StatusCode)

	validEmailDataInBytes, err := io.ReadAll(validEmailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validEmailResp.Body.Close()

	var validEmailResult map[string]interface{}
	err = json.Unmarshal(validEmailDataInBytes, &validEmailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validEmailResult["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", validEmailResult["flowType"])
	assert.Equal(t, 4, len(validEmailResult))
}

func TestCreateCodeAPIWithRidAsRandom(t *testing.T) {
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
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     "4398792-test-id",
										ClientSecret: "test-secret",
									},
								},
							},
						},
					},
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	validEmail := map[string]interface{}{
		"email": "test@example.com",
	}

	validEmailBody, err := json.Marshal(validEmail)
	if err != nil {
		t.Error(err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", testServer.URL+"/auth/signinup/code", bytes.NewBuffer(validEmailBody))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "random")

	validEmailResp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validEmailResp.StatusCode)

	validEmailDataInBytes, err := io.ReadAll(validEmailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validEmailResp.Body.Close()

	var validEmailResult map[string]interface{}
	err = json.Unmarshal(validEmailDataInBytes, &validEmailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validEmailResult["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", validEmailResult["flowType"])
	assert.Equal(t, 4, len(validEmailResult))
}

func TestCreateCodeAPIWithWrongRid(t *testing.T) {
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
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     "4398792-test-id",
										ClientSecret: "test-secret",
									},
								},
							},
						},
					},
				},
			}),
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	validEmail := map[string]interface{}{
		"email": "test@example.com",
	}

	validEmailBody, err := json.Marshal(validEmail)
	if err != nil {
		t.Error(err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", testServer.URL+"/auth/signinup/code", bytes.NewBuffer(validEmailBody))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "emailpassword")

	validEmailResp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validEmailResp.StatusCode)

	validEmailDataInBytes, err := io.ReadAll(validEmailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validEmailResp.Body.Close()

	var validEmailResult map[string]interface{}
	err = json.Unmarshal(validEmailDataInBytes, &validEmailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validEmailResult["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", validEmailResult["flowType"])
	assert.Equal(t, 4, len(validEmailResult))
}

func TestWithEmailExistAPI(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	query := req.URL.Query()
	query.Add("email", "test@example.com")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	emailDoesNotExistResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailDoesNotExistResp.StatusCode)

	emailDoesNotExistResponse := *unittesting.HttpResponseToConsumableInformation(emailDoesNotExistResp.Body)

	assert.Equal(t, "OK", emailDoesNotExistResponse["status"])
	assert.False(t, emailDoesNotExistResponse["exists"].(bool))

	codeInfo, err := CreateCodeWithEmail("public", "test@example.com", nil)
	assert.NoError(t, err)

	ConsumeCodeWithLinkCode("public", codeInfo.OK.LinkCode, codeInfo.OK.PreAuthSessionID)

	req, err = http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	query = req.URL.Query()
	query.Add("email", "test@example.com")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	emailExistsResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailExistsResp.StatusCode)

	emailExistsResponse := *unittesting.HttpResponseToConsumableInformation(emailExistsResp.Body)

	assert.Equal(t, "OK", emailExistsResponse["status"])
	assert.True(t, emailExistsResponse["exists"].(bool))
}

func TestWithEmailExistAPINewPath(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/passwordless/email/exists", nil)
	query := req.URL.Query()
	query.Add("email", "test@example.com")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	emailDoesNotExistResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailDoesNotExistResp.StatusCode)

	emailDoesNotExistResponse := *unittesting.HttpResponseToConsumableInformation(emailDoesNotExistResp.Body)

	assert.Equal(t, "OK", emailDoesNotExistResponse["status"])
	assert.False(t, emailDoesNotExistResponse["exists"].(bool))

	codeInfo, err := CreateCodeWithEmail("public", "test@example.com", nil)
	assert.NoError(t, err)

	ConsumeCodeWithLinkCode("public", codeInfo.OK.LinkCode, codeInfo.OK.PreAuthSessionID)

	req, err = http.NewRequest(http.MethodGet, testServer.URL+"/auth/passwordless/email/exists", nil)
	query = req.URL.Query()
	query.Add("email", "test@example.com")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	emailExistsResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailExistsResp.StatusCode)

	emailExistsResponse := *unittesting.HttpResponseToConsumableInformation(emailExistsResp.Body)

	assert.Equal(t, "OK", emailExistsResponse["status"])
	assert.True(t, emailExistsResponse["exists"].(bool))
}

func TestMagicLinkFormatInCreateCodeAPI(t *testing.T) {
	var magicLinkURL *url.URL
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		magicLinkURL, _ = url.Parse(*input.PasswordlessLogin.UrlWithLinkCode)
		return nil
	}
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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	validCreateCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validCreateCodeResp.StatusCode)

	validCreateCodeResponse := *unittesting.HttpResponseToConsumableInformation(validCreateCodeResp.Body)

	assert.Equal(t, "OK", validCreateCodeResponse["status"])
	assert.Equal(t, "supertokens.io", magicLinkURL.Hostname())
	assert.Equal(t, "/auth/verify", magicLinkURL.Path)
	assert.Equal(t, "", magicLinkURL.Query().Get("rid"))
	assert.Equal(t, validCreateCodeResponse["preAuthSessionId"], magicLinkURL.Query().Get("preAuthSessionId"))
}

func TestPhoneNumberToAUsersInfoAndSigningInWillSignInTheSameUserUsingTheEmailOrPhoneContactMethod(t *testing.T) {
	var userInputCodeRef string
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		userInputCodeRef = *input.PasswordlessLogin.UserInputCode
		return nil
	}
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		userInputCodeRef = *input.PasswordlessLogin.UserInputCode
		return nil
	}

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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &smsdelivery.SmsDeliveryInterface{
						SendSms: &sendSms,
					},
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	emailCreateCodeResult := *unittesting.HttpResponseToConsumableInformation(emailResp.Body)

	assert.Equal(t, "OK", emailCreateCodeResult["status"])

	consumeCodePostData := map[string]interface{}{
		"preAuthSessionId": emailCreateCodeResult["preAuthSessionId"],
		"userInputCode":    userInputCodeRef,
		"deviceId":         emailCreateCodeResult["deviceId"],
	}

	consumeCodePostBody, err := json.Marshal(consumeCodePostData)
	if err != nil {
		t.Error(err.Error())
	}

	consumeCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, consumeCodeResp.StatusCode)

	emailUserInputCodeResponse := *unittesting.HttpResponseToConsumableInformation(consumeCodeResp.Body)

	assert.Equal(t, "OK", emailUserInputCodeResponse["status"])
	user := emailUserInputCodeResponse["user"].(map[string]interface{})

	phoneNumber := "+12345678901"
	UpdateUser(user["id"].(string), nil, &phoneNumber)

	phoneNumberPostData := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneNumberPostBody, err := json.Marshal(phoneNumberPostData)
	if err != nil {
		t.Error(err.Error())
	}

	phoneNumberPostResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneNumberPostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneNumberPostResp.StatusCode)

	phoneCreateCodeResponse := *unittesting.HttpResponseToConsumableInformation(phoneNumberPostResp.Body)

	assert.Equal(t, "OK", phoneCreateCodeResponse["status"])

	consumeCodePostData1 := map[string]interface{}{
		"preAuthSessionId": phoneCreateCodeResponse["preAuthSessionId"],
		"userInputCode":    userInputCodeRef,
		"deviceId":         phoneCreateCodeResponse["deviceId"],
	}

	consumeCodePostBody1, err := json.Marshal(consumeCodePostData1)
	if err != nil {
		t.Error(err.Error())
	}

	consumeCodeResp1, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody1))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, consumeCodeResp1.StatusCode)

	phoneUserInputCodeResponse := *unittesting.HttpResponseToConsumableInformation(consumeCodeResp1.Body)

	assert.Equal(t, "OK", phoneUserInputCodeResponse["status"])
	user1 := phoneUserInputCodeResponse["user"].(map[string]interface{})

	assert.Equal(t, user["id"], user1["id"])
}

func TestWithInvalidInputToCreateCodeAPIWhileUsingTheEmailOrPhoneContactMethod(t *testing.T) {
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
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	postData := map[string]interface{}{
		"email":       "test@example.com",
		"phoneNumber": "+12345678901",
	}

	postBody, err := json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(postBody))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NoError(t, err)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)

	assert.Equal(t, "Please provide exactly one of email or phoneNumber", result["message"])

	postData = map[string]interface{}{}

	postBody, err = json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}

	resp1, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(postBody))
	assert.Equal(t, http.StatusBadRequest, resp1.StatusCode)
	assert.NoError(t, err)

	result1 := *unittesting.HttpResponseToConsumableInformation(resp1.Body)

	assert.Equal(t, "Please provide exactly one of email or phoneNumber", result1["message"])
}

func TestForCreatingACodeWithEmailAndThenResendingTheCodeAndCheckThatTheSendingCustomEmailFunctionIsCalledWhileUsingTheEmailOrPhoneContactMethod(t *testing.T) {
	isCreateAndSendCustomEmailCalled := false
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		isCreateAndSendCustomEmailCalled = true
		return nil
	}
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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(emailResp.Body)

	assert.Equal(t, "OK", result["status"])
	assert.True(t, isCreateAndSendCustomEmailCalled)

	isCreateAndSendCustomEmailCalled = false

	codeResendPostData := map[string]interface{}{
		"deviceId":         result["deviceId"],
		"preAuthSessionId": result["preAuthSessionId"],
	}

	codeResendPostBody, err := json.Marshal(codeResendPostData)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendPostResp, err := http.Post(testServer.URL+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(codeResendPostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, codeResendPostResp.StatusCode)

	codeResendResult := *unittesting.HttpResponseToConsumableInformation(codeResendPostResp.Body)
	assert.Equal(t, "OK", codeResendResult["status"])
	assert.True(t, isCreateAndSendCustomEmailCalled)
}

func TestSignUpSignInFlowWithPhoneNumberUsingEmailOrPhoneContactMethod(t *testing.T) {
	var userInputCodeRef string
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		userInputCodeRef = *input.PasswordlessLogin.UserInputCode
		return nil
	}
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
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &smsdelivery.SmsDeliveryInterface{
						SendSms: &sendSms,
					},
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	phone := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneBody, err := json.Marshal(phone)
	if err != nil {
		t.Error(err.Error())
	}

	phoneResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneResp.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(phoneResp.Body)

	assert.Equal(t, "OK", result["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", result["flowType"])

	consumeCodePostData := map[string]interface{}{
		"preAuthSessionId": result["preAuthSessionId"],
		"userInputCode":    userInputCodeRef,
		"deviceId":         result["deviceId"],
	}

	consumeCodePostBody, err := json.Marshal(consumeCodePostData)
	if err != nil {
		t.Error(err.Error())
	}

	consumeCodeResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(consumeCodePostBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, consumeCodeResp.StatusCode)

	codeConsumeResult := *unittesting.HttpResponseToConsumableInformation(consumeCodeResp.Body)

	user := codeConsumeResult["user"].(map[string]interface{})
	assert.Equal(t, "OK", codeConsumeResult["status"])
	assert.True(t, codeConsumeResult["createdNewUser"].(bool))
	assert.NotNil(t, user)
	assert.Nil(t, user["email"])
	assert.NotNil(t, user["id"])
	assert.NotNil(t, user["timeJoined"])
	assert.NotNil(t, user["phoneNumber"])
}

func TestSignInUpFlowWithEmailUsingTheEmailOrPhoneContactMethod(t *testing.T) {
	var userInputCodeRef *string
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		userInputCodeRef = input.PasswordlessLogin.UserInputCode
		return nil
	}
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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	emailDataInBytes, err := io.ReadAll(emailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailResp.Body.Close()

	var emailResult map[string]interface{}
	err = json.Unmarshal(emailDataInBytes, &emailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailResult["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", emailResult["flowType"])
	assert.Equal(t, 4, len(emailResult))

	//consume code API
	codeResendPostBody := map[string]interface{}{
		"deviceId":         emailResult["deviceId"],
		"userInputCode":    *userInputCodeRef,
		"preAuthSessionId": emailResult["preAuthSessionId"],
	}

	codeResendPostBodyJson, err := json.Marshal(codeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(codeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, codeResendResp.StatusCode)

	codeResendRespInBytes, err := io.ReadAll(codeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	codeResendResp.Body.Close()

	var codeResendResult map[string]interface{}
	err = json.Unmarshal(codeResendRespInBytes, &codeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", codeResendResult["status"])
	assert.True(t, codeResendResult["createdNewUser"].(bool))
	assert.Equal(t, 3, len(codeResendResult))
	assert.Equal(t, 5, len(codeResendResult["user"].(map[string]interface{})))
	assert.Equal(t, "test@example.com", codeResendResult["user"].(map[string]interface{})["email"])
}

func TestSignInUpFlowWithPhoneNumberUsingTheEmailOrPhoneContactMethod(t *testing.T) {
	var userInputCodeRef *string
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		userInputCodeRef = input.PasswordlessLogin.UserInputCode
		return nil
	}
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
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &smsdelivery.SmsDeliveryInterface{
						SendSms: &sendSms,
					},
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	phone := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneBody, err := json.Marshal(phone)
	if err != nil {
		t.Error(err.Error())
	}

	phoneResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneResp.StatusCode)

	phoneDataInBytes, err := io.ReadAll(phoneResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneResp.Body.Close()

	var phoneResult map[string]interface{}
	err = json.Unmarshal(phoneDataInBytes, &phoneResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneResult["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", phoneResult["flowType"])
	assert.Equal(t, 4, len(phoneResult))

	//consume code API
	codeResendPostBody := map[string]interface{}{
		"deviceId":         phoneResult["deviceId"],
		"userInputCode":    *userInputCodeRef,
		"preAuthSessionId": phoneResult["preAuthSessionId"],
	}

	codeResendPostBodyJson, err := json.Marshal(codeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(codeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, codeResendResp.StatusCode)

	codeResendRespInBytes, err := io.ReadAll(codeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	codeResendResp.Body.Close()

	var codeResendResult map[string]interface{}
	err = json.Unmarshal(codeResendRespInBytes, &codeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", codeResendResult["status"])
	assert.True(t, codeResendResult["createdNewUser"].(bool))
	assert.Equal(t, 3, len(codeResendResult))
	assert.Equal(t, 5, len(codeResendResult["user"].(map[string]interface{})))
	assert.Equal(t, "+12345678901", codeResendResult["user"].(map[string]interface{})["phoneNumber"])
}

func TestInvalidInputToCreateCodeApiUsingTheEmailOrPhoneContactMethod(t *testing.T) {
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
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	radomData1 := map[string]interface{}{
		"phoneNumber": "+12345678901",
		"email":       "test@example.com",
	}

	randomBody1, err := json.Marshal(radomData1)
	if err != nil {
		t.Error(err.Error())
	}

	randomResp1, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(randomBody1))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, randomResp1.StatusCode)

	randomDataInBytes1, err := io.ReadAll(randomResp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	randomResp1.Body.Close()

	var randomResult1 map[string]interface{}
	err = json.Unmarshal(randomDataInBytes1, &randomResult1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "Please provide exactly one of email or phoneNumber", randomResult1["message"])

	radomData2 := map[string]interface{}{}

	randomBody2, err := json.Marshal(radomData2)
	if err != nil {
		t.Error(err.Error())
	}

	randomResp2, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(randomBody2))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, randomResp2.StatusCode)

	randomDataInBytes2, err := io.ReadAll(randomResp2.Body)
	if err != nil {
		t.Error(err.Error())
	}
	randomResp2.Body.Close()

	var randomResult2 map[string]interface{}
	err = json.Unmarshal(randomDataInBytes2, &randomResult2)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "Please provide exactly one of email or phoneNumber", randomResult2["message"])
}

func TestAddingPhoneNumberToAUsersInfoAndSigningInWillSignInTheSameUserUsingTheEmailOrPhoneContractMethod(t *testing.T) {
	var userInputCodeRef *string
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		userInputCodeRef = input.PasswordlessLogin.UserInputCode
		return nil
	}
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		userInputCodeRef = input.PasswordlessLogin.UserInputCode
		return nil
	}
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
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emaildelivery.EmailDeliveryInterface{
						SendEmail: &sendEmail,
					},
				},
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &smsdelivery.SmsDeliveryInterface{
						SendSms: &sendSms,
					},
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	email := map[string]interface{}{
		"email": "test@example.com",
	}

	emailBody, err := json.Marshal(email)
	if err != nil {
		t.Error(err.Error())
	}

	emailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(emailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	emailDataInBytes, err := io.ReadAll(emailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailResp.Body.Close()

	var emailResult map[string]interface{}
	err = json.Unmarshal(emailDataInBytes, &emailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailResult["status"])

	emailCodeResendPostBody := map[string]interface{}{
		"deviceId":         emailResult["deviceId"],
		"userInputCode":    *userInputCodeRef,
		"preAuthSessionId": emailResult["preAuthSessionId"],
	}

	emailCodeResendPostBodyJson, err := json.Marshal(emailCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	emailCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(emailCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailCodeResendResp.StatusCode)

	emailCodeResendRespInBytes, err := io.ReadAll(emailCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailCodeResendResp.Body.Close()

	var emailCodeResendResult map[string]interface{}
	err = json.Unmarshal(emailCodeResendRespInBytes, &emailCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailCodeResendResult["status"])

	emailForUpdating := emailCodeResendResult["user"].(map[string]interface{})["email"].(string)
	phoneNumberForUpdating := "+12345678901"

	_, err = UpdateUser(emailCodeResendResult["user"].(map[string]interface{})["id"].(string), &emailForUpdating, &phoneNumberForUpdating)

	assert.NoError(t, err)

	phone := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	phoneBody, err := json.Marshal(phone)
	if err != nil {
		t.Error(err.Error())
	}

	phoneResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(phoneBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneResp.StatusCode)

	phoneDataInBytes, err := io.ReadAll(phoneResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneResp.Body.Close()

	var phoneResult map[string]interface{}
	err = json.Unmarshal(phoneDataInBytes, &phoneResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneResult["status"])

	phoneCodeResendBody := map[string]interface{}{
		"deviceId":         phoneResult["deviceId"],
		"userInputCode":    *userInputCodeRef,
		"preAuthSessionId": phoneResult["preAuthSessionId"],
	}

	phoneCodeResendPostBodyJson, err := json.Marshal(phoneCodeResendBody)
	if err != nil {
		t.Error(err.Error())
	}

	phoneCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(phoneCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneCodeResendResp.StatusCode)

	phoneCodeResendRespInBytes, err := io.ReadAll(phoneCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneCodeResendResp.Body.Close()

	var phoneCodeResendResult map[string]interface{}
	err = json.Unmarshal(phoneCodeResendRespInBytes, &phoneCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneCodeResendResult["status"])

	assert.Equal(t, emailCodeResendResult["user"].(map[string]interface{})["id"], phoneCodeResendResult["user"].(map[string]interface{})["id"])
}

func TestNotPassingAnyFieldsToConsumeCodeAPI(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	codeResendPostBody := map[string]interface{}{
		"preAuthSessionId": "preAuthSessionId",
	}

	codeResendPostBodyJson, err := json.Marshal(codeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	codeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(codeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, codeResendResp.StatusCode)

	codeResendRespInBytes, err := io.ReadAll(codeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	codeResendResp.Body.Close()

	var codeResendResult map[string]interface{}
	err = json.Unmarshal(codeResendRespInBytes, &codeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "Please provide one of (linkCode) or (deviceId+userInputCode) and not both", codeResendResult["message"])
}

func TestConsumeCodeAPIWithMagicLink(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	codeInfo, err := CreateCodeWithEmail("public", "test@example.com", nil)
	assert.NoError(t, err)

	invalidCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"linkCode":         "invalidLinkCode",
	}

	invalidCodeResendPostBodyJson, err := json.Marshal(invalidCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	invalidCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(invalidCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, invalidCodeResendResp.StatusCode)

	invalidCodeResendRespInBytes, err := io.ReadAll(invalidCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	invalidCodeResendResp.Body.Close()

	var invalidCodeResendResult map[string]interface{}
	err = json.Unmarshal(invalidCodeResendRespInBytes, &invalidCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "RESTART_FLOW_ERROR", invalidCodeResendResult["status"])

	validCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"linkCode":         codeInfo.OK.LinkCode,
	}

	validCodeResendPostBodyJson, err := json.Marshal(validCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	validCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(validCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validCodeResendResp.StatusCode)

	validCodeResendRespInBytes, err := io.ReadAll(validCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validCodeResendResp.Body.Close()

	var validCodeResendResult map[string]interface{}
	err = json.Unmarshal(validCodeResendRespInBytes, &validCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validCodeResendResult["status"])
	assert.True(t, validCodeResendResult["createdNewUser"].(bool))
	assert.Equal(t, 3, len(validCodeResendResult))
	assert.Equal(t, 5, len(validCodeResendResult["user"].(map[string]interface{})))
	assert.Equal(t, "test@example.com", validCodeResendResult["user"].(map[string]interface{})["email"])
}

func TestConsumeCodeAPIWithCode(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	codeInfo, err := CreateCodeWithEmail("public", "test@example.com", nil)
	assert.NoError(t, err)

	invalidCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"userInputCode":    "invalidLinkCode",
		"deviceId":         codeInfo.OK.DeviceID,
	}

	invalidCodeResendPostBodyJson, err := json.Marshal(invalidCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	invalidCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(invalidCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, invalidCodeResendResp.StatusCode)

	invalidCodeResendRespInBytes, err := io.ReadAll(invalidCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	invalidCodeResendResp.Body.Close()

	var invalidCodeResendResult map[string]interface{}
	err = json.Unmarshal(invalidCodeResendRespInBytes, &invalidCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "INCORRECT_USER_INPUT_CODE_ERROR", invalidCodeResendResult["status"])
	assert.Equal(t, float64(1), invalidCodeResendResult["failedCodeInputAttemptCount"])
	assert.Equal(t, float64(5), invalidCodeResendResult["maximumCodeInputAttempts"])
	assert.Equal(t, 3, len(invalidCodeResendResult))

	validCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"userInputCode":    codeInfo.OK.UserInputCode,
		"deviceId":         codeInfo.OK.DeviceID,
	}

	validCodeResendPostBodyJson, err := json.Marshal(validCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	validCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(validCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validCodeResendResp.StatusCode)

	validCodeResendRespInBytes, err := io.ReadAll(validCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validCodeResendResp.Body.Close()

	var validCodeResendResult map[string]interface{}
	err = json.Unmarshal(validCodeResendRespInBytes, &validCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validCodeResendResult["status"])
	assert.True(t, validCodeResendResult["createdNewUser"].(bool))
	assert.Equal(t, 3, len(validCodeResendResult))
	assert.Equal(t, 5, len(validCodeResendResult["user"].(map[string]interface{})))
	assert.Equal(t, "test@example.com", validCodeResendResult["user"].(map[string]interface{})["email"])

	usedCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"userInputCode":    codeInfo.OK.UserInputCode,
		"deviceId":         codeInfo.OK.DeviceID,
	}

	usedCodeResendPostBodyJson, err := json.Marshal(usedCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	usedCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(usedCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, usedCodeResendResp.StatusCode)

	usedCodeResendRespInBytes, err := io.ReadAll(usedCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	usedCodeResendResp.Body.Close()

	var usedCodeResendResult map[string]interface{}
	err = json.Unmarshal(usedCodeResendRespInBytes, &usedCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "RESTART_FLOW_ERROR", usedCodeResendResult["status"])
}

func TestConsumeCodeAPIWithExpiredCode(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	codeInfo, err := CreateCodeWithEmail("public", "test@example.com", nil)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	expiredCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"userInputCode":    codeInfo.OK.UserInputCode,
		"deviceId":         codeInfo.OK.DeviceID,
	}

	expiredCodeResendPostBodyJson, err := json.Marshal(expiredCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	expiredCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(expiredCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, expiredCodeResendResp.StatusCode)

	expiredCodeResendRespInBytes, err := io.ReadAll(expiredCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	expiredCodeResendResp.Body.Close()

	var expiredCodeResendResult map[string]interface{}
	err = json.Unmarshal(expiredCodeResendRespInBytes, &expiredCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "EXPIRED_USER_INPUT_CODE_ERROR", expiredCodeResendResult["status"])
	assert.Equal(t, float64(1), expiredCodeResendResult["failedCodeInputAttemptCount"])
	assert.Equal(t, float64(5), expiredCodeResendResult["maximumCodeInputAttempts"])
	assert.Equal(t, 3, len(expiredCodeResendResult))
}

func TestCreateCodeAPIWithEmail(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	validEmail := map[string]interface{}{
		"email": "test@example.com",
	}

	validEmailBody, err := json.Marshal(validEmail)
	if err != nil {
		t.Error(err.Error())
	}

	validEmailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(validEmailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validEmailResp.StatusCode)

	validEmailDataInBytes, err := io.ReadAll(validEmailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validEmailResp.Body.Close()

	var validEmailResult map[string]interface{}
	err = json.Unmarshal(validEmailDataInBytes, &validEmailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validEmailResult["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", validEmailResult["flowType"])
	assert.Equal(t, 4, len(validEmailResult))

	inValidEmail := map[string]interface{}{
		"email": "testple",
	}

	inValidEmailBody, err := json.Marshal(inValidEmail)
	if err != nil {
		t.Error(err.Error())
	}

	inValidEmailResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(inValidEmailBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, inValidEmailResp.StatusCode)

	inValidEmailDataInBytes, err := io.ReadAll(inValidEmailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	inValidEmailResp.Body.Close()

	var inValidEmailResult map[string]interface{}
	err = json.Unmarshal(inValidEmailDataInBytes, &inValidEmailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "GENERAL_ERROR", inValidEmailResult["status"])
	assert.Equal(t, "Email is invalid", inValidEmailResult["message"])
}

func TestCreateCodeAPIWithPhoneNumber(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	validPhoneNumber := map[string]interface{}{
		"phoneNumber": "+12345678901",
	}

	validphoneNumberBody, err := json.Marshal(validPhoneNumber)
	if err != nil {
		t.Error(err.Error())
	}

	validphoneNumberResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(validphoneNumberBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validphoneNumberResp.StatusCode)

	validphoneNumberDataInBytes, err := io.ReadAll(validphoneNumberResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validphoneNumberResp.Body.Close()

	var validphoneNumberResult map[string]interface{}
	err = json.Unmarshal(validphoneNumberDataInBytes, &validphoneNumberResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validphoneNumberResult["status"])
	assert.Equal(t, "USER_INPUT_CODE_AND_MAGIC_LINK", validphoneNumberResult["flowType"])
	assert.Equal(t, 4, len(validphoneNumberResult))

	inValidphoneNumber := map[string]interface{}{
		"phoneNumber": "+123",
	}

	inValidphoneNumberBody, err := json.Marshal(inValidphoneNumber)
	if err != nil {
		t.Error(err.Error())
	}

	inValidphoneNumberResp, err := http.Post(testServer.URL+"/auth/signinup/code", "application/json", bytes.NewBuffer(inValidphoneNumberBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, inValidphoneNumberResp.StatusCode)

	inValidphoneNumberDataInBytes, err := io.ReadAll(inValidphoneNumberResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	inValidphoneNumberResp.Body.Close()

	var inValidphoneNumberResult map[string]interface{}
	err = json.Unmarshal(inValidphoneNumberDataInBytes, &inValidphoneNumberResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "GENERAL_ERROR", inValidphoneNumberResult["status"])
	assert.Equal(t, "Phone number is invalid", inValidphoneNumberResult["message"])
}

func TestEmailExistsAPI(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	query := req.URL.Query()
	query.Add("email", "test@example.com")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	emailResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	emailDataInBytes, err := io.ReadAll(emailResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailResp.Body.Close()

	var emailResult map[string]interface{}
	err = json.Unmarshal(emailDataInBytes, &emailResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailResult["status"])
	assert.Equal(t, false, emailResult["exists"])

	codeInfo, err := CreateCodeWithEmail("public", "test@example.com", nil)
	assert.NoError(t, err)

	_, err = ConsumeCodeWithLinkCode("public", codeInfo.OK.LinkCode, codeInfo.OK.PreAuthSessionID)
	assert.NoError(t, err)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	query1 := req.URL.Query()
	query1.Add("email", "test@example.com")
	req1.URL.RawQuery = query1.Encode()
	assert.NoError(t, err)
	emailResp1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, emailResp1.StatusCode)

	emailDataInBytes1, err := io.ReadAll(emailResp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	emailResp1.Body.Close()

	var emailResult1 map[string]interface{}
	err = json.Unmarshal(emailDataInBytes1, &emailResult1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", emailResult1["status"])
	assert.Equal(t, true, emailResult1["exists"])
}

func TestPhoneNumberExistsAPI(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/phonenumber/exists", nil)
	query := req.URL.Query()
	query.Add("phoneNumber", "+1234567890")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	phoneResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneResp.StatusCode)

	phoneDataInBytes, err := io.ReadAll(phoneResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneResp.Body.Close()

	var phoneResult map[string]interface{}
	err = json.Unmarshal(phoneDataInBytes, &phoneResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneResult["status"])
	assert.Equal(t, false, phoneResult["exists"])

	codeInfo, err := CreateCodeWithPhoneNumber("public", "+1234567890", nil)
	assert.NoError(t, err)

	_, err = ConsumeCodeWithLinkCode("public", codeInfo.OK.LinkCode, codeInfo.OK.PreAuthSessionID)
	assert.NoError(t, err)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/phonenumber/exists", nil)
	query1 := req.URL.Query()
	query1.Add("phoneNumber", "+1234567890")
	req1.URL.RawQuery = query1.Encode()
	assert.NoError(t, err)
	phoneResp1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneResp1.StatusCode)

	phoneDataInBytes1, err := io.ReadAll(phoneResp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneResp1.Body.Close()

	var phoneResult1 map[string]interface{}
	err = json.Unmarshal(phoneDataInBytes1, &phoneResult1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneResult1["status"])
	assert.Equal(t, true, phoneResult1["exists"])
}

func TestPhoneNumberExistsAPINewPath(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/passwordless/phonenumber/exists", nil)
	query := req.URL.Query()
	query.Add("phoneNumber", "+1234567890")
	req.URL.RawQuery = query.Encode()
	assert.NoError(t, err)
	phoneResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneResp.StatusCode)

	phoneDataInBytes, err := io.ReadAll(phoneResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneResp.Body.Close()

	var phoneResult map[string]interface{}
	err = json.Unmarshal(phoneDataInBytes, &phoneResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneResult["status"])
	assert.Equal(t, false, phoneResult["exists"])

	codeInfo, err := CreateCodeWithPhoneNumber("public", "+1234567890", nil)
	assert.NoError(t, err)

	_, err = ConsumeCodeWithLinkCode("public", codeInfo.OK.LinkCode, codeInfo.OK.PreAuthSessionID)
	assert.NoError(t, err)

	req1, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/passwordless/phonenumber/exists", nil)
	query1 := req.URL.Query()
	query1.Add("phoneNumber", "+1234567890")
	req1.URL.RawQuery = query1.Encode()
	assert.NoError(t, err)
	phoneResp1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, phoneResp1.StatusCode)

	phoneDataInBytes1, err := io.ReadAll(phoneResp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	phoneResp1.Body.Close()

	var phoneResult1 map[string]interface{}
	err = json.Unmarshal(phoneDataInBytes1, &phoneResult1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", phoneResult1["status"])
	assert.Equal(t, true, phoneResult1["exists"])
}

func TestResendCodeAPI(t *testing.T) {
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	codeInfo, err := CreateCodeWithPhoneNumber("public", "+1234567890", nil)
	assert.NoError(t, err)

	validCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": codeInfo.OK.PreAuthSessionID,
		"deviceId":         codeInfo.OK.DeviceID,
	}

	validCodeResendPostBodyJson, err := json.Marshal(validCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	validCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(validCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, validCodeResendResp.StatusCode)

	validCodeResendRespInBytes, err := io.ReadAll(validCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	validCodeResendResp.Body.Close()

	var validCodeResendResult map[string]interface{}
	err = json.Unmarshal(validCodeResendRespInBytes, &validCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", validCodeResendResult["status"])

	invalidCodeResendPostBody := map[string]interface{}{
		"preAuthSessionId": "asdasdasdasdsa",
		"deviceId":         "asdeflasdkjqee",
	}

	invalidCodeResendPostBodyJson, err := json.Marshal(invalidCodeResendPostBody)
	if err != nil {
		t.Error(err.Error())
	}

	invalidCodeResendResp, err := http.Post(testServer.URL+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(invalidCodeResendPostBodyJson))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, invalidCodeResendResp.StatusCode)

	invalidCodeResendRespInBytes, err := io.ReadAll(invalidCodeResendResp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	invalidCodeResendResp.Body.Close()

	var invalidCodeResendResult map[string]interface{}
	err = json.Unmarshal(invalidCodeResendRespInBytes, &invalidCodeResendResult)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "RESTART_FLOW_ERROR", invalidCodeResendResult["status"])
}
