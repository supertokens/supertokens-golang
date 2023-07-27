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

package emailpassword

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestOverridingFunctionCalls(t *testing.T) {
	var user *epmodels.User
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
			Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					Functions: func(originalImplementation epmodels.RecipeInterface) epmodels.RecipeInterface {
						originalSignup := *originalImplementation.SignUp
						originalSignIn := *originalImplementation.SignIn
						originalGetUserById := *originalImplementation.GetUserByID
						*originalImplementation.SignUp = func(email, password string, tenantId string, userContext supertokens.UserContext) (epmodels.SignUpResponse, error) {
							res, err := originalSignup(email, password, tenantId, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							user = &res.OK.User
							return res, nil
						}
						*originalImplementation.SignIn = func(email, password string, tenantId string, userContext supertokens.UserContext) (epmodels.SignInResponse, error) {
							res, err := originalSignIn(email, password, tenantId, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							user = &res.OK.User
							return res, nil
						}
						*originalImplementation.GetUserByID = func(userID string, userContext supertokens.UserContext) (*epmodels.User, error) {
							res, err := originalGetUserById(userID, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							user = res
							return res, nil
						}
						return originalImplementation
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
	mux.HandleFunc("/user", func(rw http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		fetchedUser, err := GetUserByID(userId)
		if err != nil {
			t.Error(err.Error())
		}
		jsonResp, err := json.Marshal(fetchedUser)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	res, err := unittesting.SignupRequest("user@test.com", "test123!", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)

	dataInBytes, err := io.ReadAll(res.Body)

	if err != nil {
		t.Error(err.Error())
	}

	res.Body.Close()

	var result map[string]interface{}

	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, result["user"])
	assert.Equal(t, user.ID, result["user"].(map[string]interface{})["id"].(string))
	assert.Equal(t, user.Email, result["user"].(map[string]interface{})["email"].(string))

	user = nil

	assert.Nil(t, user)

	res1, err := unittesting.SignInRequest("user@test.com", "test123!", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)

	dataInBytes1, err := io.ReadAll(res1.Body)

	if err != nil {
		t.Error(err.Error())
	}

	res1.Body.Close()

	var result1 map[string]interface{}

	err = json.Unmarshal(dataInBytes1, &result1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, result1["user"])
	assert.Equal(t, user.ID, result1["user"].(map[string]interface{})["id"].(string))
	assert.Equal(t, user.Email, result1["user"].(map[string]interface{})["email"].(string))

	user = nil

	assert.Nil(t, user)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/user", nil)

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)
	q := req.URL.Query()
	q.Add("userId", result1["user"].(map[string]interface{})["id"].(string))
	req.URL.RawQuery = q.Encode()
	assert.NoError(t, err)
	res2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	dataInBytes2, err := io.ReadAll(res2.Body)

	if err != nil {
		t.Error(err.Error())
	}

	res1.Body.Close()

	var result2 epmodels.User

	err = json.Unmarshal(dataInBytes2, &result2)
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, user)
	assert.Equal(t, user.ID, result2.ID)
	assert.Equal(t, user.Email, result2.Email)
	assert.Equal(t, user.TimeJoined, result2.TimeJoined)

}

func TestOverridingApiCalls(t *testing.T) {
	var user *epmodels.User
	var emailExists bool
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
			Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						originalSignupPost := *originalImplementation.SignUpPOST
						*originalImplementation.SignUpPOST = func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
							res, err := originalSignupPost(formFields, tenantId, options, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							user = &res.OK.User
							return res, err
						}
						originalSignInPOST := *originalImplementation.SignInPOST
						*originalImplementation.SignInPOST = func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignInPOSTResponse, error) {
							res, err := originalSignInPOST(formFields, tenantId, options, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							user = &res.OK.User
							return res, err
						}
						originalemailExistGet := *originalImplementation.EmailExistsGET
						*originalImplementation.EmailExistsGET = func(email string, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.EmailExistsGETResponse, error) {
							res, err := originalemailExistGet(email, tenantId, options, userContext)
							if err != nil {
								log.Fatal(err.Error())
							}
							emailExists = res.OK.Exists
							return res, err
						}
						return originalImplementation
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
	mux.HandleFunc("/user", func(rw http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		fetchedUser, err := GetUserByID(userId)
		if err != nil {
			t.Error(err.Error())
		}
		jsonResp, err := json.Marshal(fetchedUser)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	q := req.URL.Query()
	q.Add("email", "user@test.com")
	req.URL.RawQuery = q.Encode()
	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	dataInBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res.Body.Close()
	var result map[string]interface{}
	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		t.Error(err.Error())
	}
	assert.False(t, result["exists"].(bool))
	assert.False(t, emailExists)

	res1, err := unittesting.SignupRequest("user@test.com", "test123!", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	dataInBytes1, err := io.ReadAll(res1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res1.Body.Close()
	var result1 map[string]interface{}
	err = json.Unmarshal(dataInBytes1, &result1)
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, result1["user"])
	assert.Equal(t, user.ID, result1["user"].(map[string]interface{})["id"].(string))
	assert.Equal(t, user.Email, result1["user"].(map[string]interface{})["email"].(string))

	req3, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/signup/email/exists", nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	q3 := req3.URL.Query()
	q3.Add("email", "user@test.com")
	req3.URL.RawQuery = q3.Encode()
	httpClient1 := &http.Client{}
	res3, err := httpClient1.Do(req3)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, 200, res3.StatusCode)
	dataInBytes3, err := io.ReadAll(res3.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res3.Body.Close()
	var result3 map[string]interface{}
	err = json.Unmarshal(dataInBytes3, &result3)
	if err != nil {
		t.Error(err.Error())
	}

	assert.True(t, result3["exists"].(bool))
	assert.True(t, emailExists)

	res2, err := unittesting.SignInRequest("user@test.com", "test123!", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.NoError(t, err)

	dataInBytes2, err := io.ReadAll(res2.Body)

	if err != nil {
		t.Error(err.Error())
	}

	res2.Body.Close()

	var result2 map[string]interface{}

	err = json.Unmarshal(dataInBytes2, &result2)
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, result2["user"])
	assert.Equal(t, user.ID, result2["user"].(map[string]interface{})["id"].(string))
	assert.Equal(t, user.Email, result2["user"].(map[string]interface{})["email"].(string))

}
