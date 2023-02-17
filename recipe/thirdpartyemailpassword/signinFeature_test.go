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

package thirdpartyemailpassword

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"gopkg.in/h2non/gock.v1"
)

func TestAfterDisablingTheDefaultSigninupAPIdoesNotWork(t *testing.T) {
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
			Init(&tpepmodels.TypeInput{
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
								},
							},
						},
					},
				},
				Override: &tpepmodels.OverrideStruct{
					APIs: func(originalImplementation tpepmodels.APIInterface) tpepmodels.APIInterface {
						*originalImplementation.ThirdPartySignInUpPOST = nil
						return originalImplementation
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

	signinupPostData := map[string]interface{}{
		"thirdPartyId": "google",
		"redirectURIInfo": map[string]interface{}{
			"redirectURIOnProviderDashboard": "http://127.0.0.1/callback",
			"redirectURIQueryParams": map[string]interface{}{
				"code": "abcdefghj",
			},
		},
	}

	postBody, err := json.Marshal(signinupPostData)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAfterDisablingTheDefaultSigninAPIdoesNotWork(t *testing.T) {
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
			Init(&tpepmodels.TypeInput{
				Override: &tpepmodels.OverrideStruct{
					APIs: func(originalImplementation tpepmodels.APIInterface) tpepmodels.APIInterface {
						*originalImplementation.EmailPasswordSignInPOST = nil
						return originalImplementation
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

	resp, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestHandlePostSignUpInGetsSetCorrectly(t *testing.T) {
	userId := ""
	loginType := ""
	customAntiCsrfVal := "VIA_TOKEN"
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
			Init(&tpepmodels.TypeInput{
				Override: &tpepmodels.OverrideStruct{
					APIs: func(originalImplementation tpepmodels.APIInterface) tpepmodels.APIInterface {
						originalSignInUpPost := *originalImplementation.ThirdPartySignInUpPOST
						*originalImplementation.ThirdPartySignInUpPOST = func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, tenantId *string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpepmodels.ThirdPartySignInUpPOSTResponse, error) {
							resp, err := originalSignInUpPost(provider, input, tenantId, options, userContext)
							if err != nil {
								t.Error(err.Error())
							}
							userId = resp.OK.User.ID
							loginType = "thirdparty"
							return resp, err
						}
						return originalImplementation
					},
				},
				Providers: []tpmodels.ProviderInput{
					customProvider2,
				},
			}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
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

	defer gock.OffAll()
	gock.New("https://test.com/").
		Post("oauth/token").
		Reply(200).
		JSON(map[string]string{})

	postData := map[string]interface{}{
		"thirdPartyId": "custom",
		"redirectURIInfo": map[string]interface{}{
			"redirectURIOnProviderDashboard": "http://127.0.0.1/callback",
			"redirectURIQueryParams": map[string]interface{}{
				"code": "abcdefghj",
			},
		},
	}

	postBody, err := json.Marshal(postData)
	if err != nil {
		t.Error(err.Error())
	}

	gock.New(testServer.URL).EnableNetworking().Persist()
	gock.New("http://localhost:8080/").EnableNetworking().Persist()

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var result map[string]interface{}

	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		t.Error(err.Error())
	}

	user := result["user"].(map[string]interface{})

	assert.Equal(t, userId, user["id"])
	assert.Equal(t, "thirdparty", loginType)
}

func TestSignInAPIWorksWhenInputIsFine(t *testing.T) {
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
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
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

	resp, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]interface{}
	json.Unmarshal(dataInBytes, &result)
	resp.Body.Close()

	assert.Equal(t, "OK", result["status"])

	signupUserInfo := result["user"].(map[string]interface{})

	resp1, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	dataInBytes1, err := io.ReadAll(resp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result1 map[string]interface{}
	json.Unmarshal(dataInBytes1, &result1)
	resp1.Body.Close()

	assert.Equal(t, "OK", result1["status"])

	signInUserInfo := result1["user"].(map[string]interface{})

	assert.Equal(t, signInUserInfo["id"], signupUserInfo["id"])
	assert.Equal(t, signInUserInfo["email"], signupUserInfo["email"])
}

func TestSigninAPIthrowsAnErrorWhenEmailDoesNotExist(t *testing.T) {
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
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
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

	resp, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]interface{}
	json.Unmarshal(dataInBytes, &result)
	resp.Body.Close()

	assert.Equal(t, "OK", result["status"])

	resp1, err := unittesting.SignInRequest("rand@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	dataInBytes1, err := io.ReadAll(resp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result1 map[string]interface{}
	json.Unmarshal(dataInBytes1, &result1)
	resp1.Body.Close()

	assert.Equal(t, "WRONG_CREDENTIALS_ERROR", result1["status"])
}

func TestCustomEmailValidatorsToSignupAndMakeSureTheyAreAppliedToSignIn(t *testing.T) {
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
			Init(&tpepmodels.TypeInput{
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "email",
							Validate: func(value interface{}) *string {
								customErrorMessage := "email does not start with test"
								if strings.HasPrefix(value.(string), "test") {
									return nil
								}
								return &customErrorMessage
							},
						},
					},
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
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("testrandom@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]interface{}
	json.Unmarshal(dataInBytes, &result)
	resp.Body.Close()

	assert.Equal(t, "OK", result["status"])

	resp1, err := unittesting.SignInRequest("rand@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	dataInBytes1, err := io.ReadAll(resp1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var result1 map[string]interface{}
	json.Unmarshal(dataInBytes1, &result1)
	resp1.Body.Close()

	assert.Equal(t, "FIELD_ERROR", result1["status"])
	assert.Equal(t, "email does not start with test", result1["formFields"].([]interface{})[0].(map[string]interface{})["error"])
	assert.Equal(t, "email", result1["formFields"].([]interface{})[0].(map[string]interface{})["id"])
}
