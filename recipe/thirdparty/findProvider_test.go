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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestSingleConfigWithoutClientIDSpecified(t *testing.T) {
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
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-1",
									ClientSecret: "test-secret",
								},
							),
						},
					},
				},
			),
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=google")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	assert.NoError(t, err)

	authUrl := data["url"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-1", parsedURL.Query().Get("client_id"))
}

func TestSingleConfigWithDefaultWithoutClientIDSpecified(t *testing.T) {
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
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-1",
									ClientSecret: "test-secret",
									IsDefault:    true,
								},
							),
						},
					},
				},
			),
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=google")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	assert.NoError(t, err)

	authUrl := data["url"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-1", parsedURL.Query().Get("client_id"))
}

func TestMultipleProviderSingleConfigWithoutClientIDSpecified(t *testing.T) {
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
							Facebook(
								tpmodels.FacebookConfig{
									ClientID:     "client-id-2",
									ClientSecret: "test-secret",
								},
							),
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-1",
									ClientSecret: "test-secret",
								},
							),
						},
					},
				},
			),
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=google")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	assert.NoError(t, err)

	authUrl := data["url"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-1", parsedURL.Query().Get("client_id"))
}

func TestMultipleConfigWithoutClientIDSpecified1(t *testing.T) {
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
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-1",
									ClientSecret: "test-secret",
									IsDefault:    true,
								},
							),
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-2",
									ClientSecret: "test-secret",
								},
							),
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-3",
									ClientSecret: "test-secret",
								},
							),
						},
					},
				},
			),
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=google")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	assert.NoError(t, err)

	authUrl := data["url"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-1", parsedURL.Query().Get("client_id"))
}

func TestMultipleConfigWithoutClientIDSpecified2(t *testing.T) {
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
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-1",
									ClientSecret: "test-secret",
								},
							),
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-2",
									ClientSecret: "test-secret",
								},
							),
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-3",
									ClientSecret: "test-secret",
									IsDefault:    true,
								},
							),
						},
					},
				},
			),
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=google")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	assert.NoError(t, err)

	authUrl := data["url"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-3", parsedURL.Query().Get("client_id"))
}

func TestMultipleProviderMultipleConfigWithoutClientIDSpecified(t *testing.T) {
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
							Facebook(
								tpmodels.FacebookConfig{
									ClientID:     "client-id-1",
									ClientSecret: "test-secret",
									IsDefault:    true,
								},
							),
							Facebook(
								tpmodels.FacebookConfig{
									ClientID:     "client-id-2",
									ClientSecret: "test-secret",
								},
							),
							Facebook(
								tpmodels.FacebookConfig{
									ClientID:     "client-id-3",
									ClientSecret: "test-secret",
								},
							),
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-1",
									ClientSecret: "test-secret",
								},
							),
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-2",
									ClientSecret: "test-secret",
								},
							),
							Google(
								tpmodels.GoogleConfig{
									ClientID:     "client-id-3",
									ClientSecret: "test-secret",
									IsDefault:    true,
								},
							),
						},
					},
				},
			),
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=google")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	assert.NoError(t, err)

	authUrl := data["url"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-3", parsedURL.Query().Get("client_id"))
}
