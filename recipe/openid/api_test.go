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

package openid

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestWithDefaultConfigCallingDiscoveryConfigurationEndpointWorksAsExpected(t *testing.T) {
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

	resp, err := http.Get(testServer.URL + "/auth/.well-known/openid-configuration")
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]interface{}

	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "https://api.supertokens.io/auth", result["issuer"].(string))
	assert.Equal(t, "https://api.supertokens.io/auth/jwt/jwks.json", result["jwks_uri"].(string))
}

func TestWithAPIBasePathCallingDiscoveryConfigurationEndpointWorksAsExpected(t *testing.T) {
	customAPIBasePath := "/"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customAPIBasePath,
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
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

	resp, err := http.Get(testServer.URL + "/.well-known/openid-configuration")
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Error(err.Error())
	}
	var result map[string]interface{}

	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "https://api.supertokens.io", result["issuer"].(string))
	assert.Equal(t, "https://api.supertokens.io/jwt/jwks.json", result["jwks_uri"].(string))
}

func TestDiscoveryEndPointDoesNotWorkWhenDisabled(t *testing.T) {
	customAPIBasePath := "/"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customAPIBasePath,
		},
		RecipeList: []supertokens.Recipe{
			Init(&openidmodels.TypeInput{
				Override: &openidmodels.OverrideStruct{
					APIs: func(originalImplementation openidmodels.APIInterface) openidmodels.APIInterface {
						*originalImplementation.GetOpenIdDiscoveryConfigurationGET = nil
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

	resp, err := http.Get(testServer.URL + "/.well-known/openid-configuration")
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
