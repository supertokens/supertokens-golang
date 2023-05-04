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

package jwt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestOverridingFunctions(t *testing.T) {
	var jwtCreated string
	var jwksKeys []jwtmodels.JsonWebKeys
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
			Init(&jwtmodels.TypeInput{
				Override: &jwtmodels.OverrideStruct{
					Functions: func(originalImplementation jwtmodels.RecipeInterface) jwtmodels.RecipeInterface {
						createJWToriginal := *originalImplementation.CreateJWT
						getJWKSOriginal := *originalImplementation.GetJWKS
						*originalImplementation.CreateJWT = func(payload map[string]interface{}, validitySeconds *uint64, useStaticSigningKey *bool, userContext supertokens.UserContext) (jwtmodels.CreateJWTResponse, error) {
							resp, err := createJWToriginal(payload, validitySeconds, useStaticSigningKey, userContext)
							if err != nil {
								t.Error(err.Error())
								return jwtmodels.CreateJWTResponse{}, err
							}
							jwtCreated = resp.OK.Jwt
							return resp, nil
						}
						*originalImplementation.GetJWKS = func(userContext supertokens.UserContext) (jwtmodels.GetJWKSResponse, error) {
							resp, err := getJWKSOriginal(userContext)
							if err != nil {
								t.Error(err.Error())
								return jwtmodels.GetJWKSResponse{}, err
							}
							jwksKeys = resp.OK.Keys
							return resp, nil
						}
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

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.8") == "2.8" {
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/jwtcreate", func(w http.ResponseWriter, r *http.Request) {
		dataInBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		var result map[string]interface{}
		err = json.Unmarshal(dataInBytes, &result)
		r.Body.Close()
		if err != nil {
			t.Error(err.Error())
		}
		payload := result["payload"]
		validity := uint64(1000)
		resp, err := CreateJWT(payload.(map[string]interface{}), &validity, nil)
		if err != nil {
			t.Error(err.Error())
		}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(200)
		w.Write(jsonResp)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	formFields := map[string]interface{}{
		"payload": map[string]interface{}{
			"someKey": "someValue",
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/jwtcreate", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Error(err.Error())
	}
	var result jwtmodels.CreateJWTResponse

	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, jwtCreated, result.OK.Jwt)

	resp1, err := http.Get(testServer.URL + "/auth/jwt/jwks.json")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	data1, err := ioutil.ReadAll(resp1.Body)
	resp1.Body.Close()
	if err != nil {
		t.Error(err.Error())
	}
	var result1 map[string]interface{}

	err = json.Unmarshal(data1, &result1)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, jwksKeys[0].Kty, result1["keys"].([]interface{})[0].(map[string]interface{})["kty"])
	assert.Equal(t, jwksKeys[0].Alg, result1["keys"].([]interface{})[0].(map[string]interface{})["alg"])
	assert.Equal(t, jwksKeys[0].E, result1["keys"].([]interface{})[0].(map[string]interface{})["e"])
	assert.Equal(t, jwksKeys[0].Kid, result1["keys"].([]interface{})[0].(map[string]interface{})["kid"])
	assert.Equal(t, jwksKeys[0].N, result1["keys"].([]interface{})[0].(map[string]interface{})["n"])
	assert.Equal(t, jwksKeys[0].Use, result1["keys"].([]interface{})[0].(map[string]interface{})["use"])
}

func TestOverridingAPI(t *testing.T) {
	var jwksKeys []jwtmodels.JsonWebKeys
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
			Init(&jwtmodels.TypeInput{
				Override: &jwtmodels.OverrideStruct{
					APIs: func(originalImplementation jwtmodels.APIInterface) jwtmodels.APIInterface {
						getJWKSOriginal := *originalImplementation.GetJWKSGET
						*originalImplementation.GetJWKSGET = func(options jwtmodels.APIOptions, userContext supertokens.UserContext) (jwtmodels.GetJWKSAPIResponse, error) {
							resp, err := getJWKSOriginal(options, userContext)
							if err != nil {
								t.Error(err.Error())
								return jwtmodels.GetJWKSAPIResponse{}, err
							}
							jwksKeys = resp.OK.Keys
							return resp, nil
						}
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

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.8") == "2.8" {
		return
	}
	mux := http.NewServeMux()

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp1, err := http.Get(testServer.URL + "/auth/jwt/jwks.json")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	data1, err := ioutil.ReadAll(resp1.Body)
	resp1.Body.Close()
	if err != nil {
		t.Error(err.Error())
	}
	var result1 map[string]interface{}

	err = json.Unmarshal(data1, &result1)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, jwksKeys[0].Kty, result1["keys"].([]interface{})[0].(map[string]interface{})["kty"])
	assert.Equal(t, jwksKeys[0].Alg, result1["keys"].([]interface{})[0].(map[string]interface{})["alg"])
	assert.Equal(t, jwksKeys[0].E, result1["keys"].([]interface{})[0].(map[string]interface{})["e"])
	assert.Equal(t, jwksKeys[0].Kid, result1["keys"].([]interface{})[0].(map[string]interface{})["kid"])
	assert.Equal(t, jwksKeys[0].N, result1["keys"].([]interface{})[0].(map[string]interface{})["n"])
	assert.Equal(t, jwksKeys[0].Use, result1["keys"].([]interface{})[0].(map[string]interface{})["use"])
}
