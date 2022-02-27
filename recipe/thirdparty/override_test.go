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

package thirdparty

import (
	"bytes"
	"encoding/json"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"gopkg.in/h2non/gock.v1"
)

func TestOverrideFunctions(t *testing.T) {
	var createdNewUser bool
	var user tpmodels.User
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
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							CUSTOM_PROVIDER_1,
						},
					},
					Override: &tpmodels.OverrideStruct{
						Functions: func(originalImplementation tpmodels.RecipeInterface) tpmodels.RecipeInterface {
							originalSignInUp := *originalImplementation.SignInUp
							originalGetUserById := *originalImplementation.GetUserByID
							*originalImplementation.SignInUp = func(thirdPartyID, thirdPartyUserID string, email tpmodels.EmailStruct, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) {
								resp, err := originalSignInUp(thirdPartyID, thirdPartyUserID, email, userContext)
								if err != nil {
									return tpmodels.SignInUpResponse{}, err
								}
								user = resp.OK.User
								createdNewUser = resp.OK.CreatedNewUser
								return resp, nil
							}
							*originalImplementation.GetUserByID = func(userID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
								resp, err := originalGetUserById(userID, userContext)
								if err != nil {
									return nil, err
								}
								user = *resp
								return resp, nil
							}
							return originalImplementation
						},
					},
				},
			),
			session.Init(&sessmodels.TypeInput{}),
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

	defer gock.OffAll()
	gock.New("https://test.com").
		Post("/oauth/token").
		Persist().
		Reply(200).
		JSON(map[string]string{})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	formFields := map[string]string{
		"thirdPartyId": "custom",
		"code":         "abcdefghj",
		"redirectURI":  testServer.URL + "/callback",
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	gock.New(testServer.URL).EnableNetworking().Persist()
	gock.New("http://localhost:8080/").EnableNetworking().Persist()

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, user)
	assert.True(t, createdNewUser)
	assert.Equal(t, user.Email, data["user"].(map[string]interface{})["Email"])
	assert.Equal(t, user.ID, data["user"].(map[string]interface{})["ID"])

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/user", nil)
	if err != nil {
		t.Error(err.Error())
	}
	q := req.URL.Query()
	q.Add("userId", data["user"].(map[string]interface{})["ID"].(string))
	req.URL.RawQuery = q.Encode()
	res1, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	dataInBytes1, err := io.ReadAll(res1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res1.Body.Close()
	var data1 map[string]interface{}
	err = json.Unmarshal(dataInBytes1, &data1)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, user.Email, data1["Email"])
}

func TestOverrideAPIs(t *testing.T) {
	var createdNewUser bool
	var user tpmodels.User
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
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							CUSTOM_PROVIDER_1,
						},
					},
					Override: &tpmodels.OverrideStruct{
						APIs: func(originalImplementation tpmodels.APIInterface) tpmodels.APIInterface {
							originalSigniupPost := *originalImplementation.SignInUpPOST
							*originalImplementation.SignInUpPOST = func(provider tpmodels.TypeProvider, code string, authCodeResponse interface{}, redirectURI string, options tpmodels.APIOptions, userContext *map[string]interface{}) (tpmodels.SignInUpPOSTResponse, error) {
								res, err := originalSigniupPost(provider, code, authCodeResponse, redirectURI, options, userContext)
								if err != nil {
									t.Error(err.Error())
								}
								user = res.OK.User
								createdNewUser = res.OK.CreatedNewUser
								return res, err
							}
							return originalImplementation
						},
					},
				},
			),
			session.Init(&sessmodels.TypeInput{}),
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

	defer gock.OffAll()
	gock.New("https://test.com").
		Post("/oauth/token").
		Persist().
		Reply(200).
		JSON(map[string]string{})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	formFields := map[string]string{
		"thirdPartyId": "custom",
		"code":         "abcdefghj",
		"redirectURI":  testServer.URL + "/callback",
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	gock.New(testServer.URL).EnableNetworking().Persist()
	gock.New("http://localhost:8080/").EnableNetworking().Persist()

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err.Error())
	}

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, user)
	assert.True(t, createdNewUser)
	assert.Equal(t, user.Email, data["user"].(map[string]interface{})["Email"])
	assert.Equal(t, user.ID, data["user"].(map[string]interface{})["ID"])

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/user", nil)
	if err != nil {
		t.Error(err.Error())
	}
	q := req.URL.Query()
	q.Add("userId", data["user"].(map[string]interface{})["ID"].(string))
	req.URL.RawQuery = q.Encode()
	res1, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	dataInBytes1, err := io.ReadAll(res1.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res1.Body.Close()
	var data1 map[string]interface{}
	err = json.Unmarshal(dataInBytes1, &data1)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, user.Email, data1["Email"])
}
