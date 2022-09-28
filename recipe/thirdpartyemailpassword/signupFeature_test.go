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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"gopkg.in/h2non/gock.v1"
)

func TestDisablingDefaultAPIDoesNotWork(t *testing.T) {
	// TODO: fix this test
	// configValue := supertokens.TypeInput{
	// 	Supertokens: &supertokens.ConnectionInfo{
	// 		ConnectionURI: "http://localhost:8080",
	// 	},
	// 	AppInfo: supertokens.AppInfo{
	// 		APIDomain:     "api.supertokens.io",
	// 		AppName:       "SuperTokens",
	// 		WebsiteDomain: "supertokens.io",
	// 	},
	// 	RecipeList: []supertokens.Recipe{
	// 		Init(&tpepmodels.TypeInput{
	// 			Providers: []tpmodels.TypeProvider{
	// 				thirdparty.Google(tpmodels.GoogleConfig{
	// 					ClientID:     "test",
	// 					ClientSecret: "test-secret",
	// 				}),
	// 			},
	// 			Override: &tpepmodels.OverrideStruct{
	// 				APIs: func(originalImplementation tpepmodels.APIInterface) tpepmodels.APIInterface {
	// 					*originalImplementation.ThirdPartySignInUpPOST = nil
	// 					return originalImplementation
	// 				},
	// 			},
	// 		}),
	// 	},
	// }

	// BeforeEach()
	// unittesting.StartUpST("localhost", "8080")
	// defer AfterEach()
	// err := supertokens.Init(configValue)
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// mux := http.NewServeMux()
	// testServer := httptest.NewServer(supertokens.Middleware(mux))
	// defer testServer.Close()

	// signinupPostData := map[string]string{
	// 	"thirdPartyId": "google",
	// 	"code":         "abcdefghj",
	// 	"redirectURI":  "http://127.0.0.1/callback",
	// }

	// postBody, err := json.Marshal(signinupPostData)
	// if err != nil {
	// 	t.Error(err.Error())
	// }

	// resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestThatIfDisableAPIDefaultSignupAPIDoesNotWork(t *testing.T) {
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
						*originalImplementation.EmailPasswordSignUpPOST = nil
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestMinimumConfigForOneProvider(t *testing.T) {
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
				Providers: []tpmodels.TypeProvider{
					customProvider2,
				},
				Override: &tpepmodels.OverrideStruct{
					APIs: func(originalImplementation tpepmodels.APIInterface) tpepmodels.APIInterface {
						*originalImplementation.EmailPasswordSignUpPOST = nil
						return originalImplementation
					},
				},
			}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	postData := map[string]string{
		"thirdPartyId": "custom",
		"code":         "abcdefghj",
		"redirectURI":  "http://127.0.0.1/callback",
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

	assert.Equal(t, "OK", result["status"])
	assert.Equal(t, true, result["createdNewUser"])
	assert.Equal(t, "email@test.com", user["email"])
	assert.Equal(t, "custom", user["thirdParty"].(map[string]interface{})["id"])
	assert.Equal(t, "user", user["thirdParty"].(map[string]interface{})["userId"])
}

func TestSignUpAPIWorksWhenInputIsFine(t *testing.T) {
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	dataInBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var result map[string]interface{}

	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "OK", result["status"])
	assert.Equal(t, "random@gmail.com", result["user"].(map[string]interface{})["email"])
}
