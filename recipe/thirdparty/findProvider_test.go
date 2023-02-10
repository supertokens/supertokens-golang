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

func TestSingleConfigWithoutClientTypeSpecified(t *testing.T) {
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
						Providers: []tpmodels.ProviderInput{
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "google",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientID:     "client-id-1",
											ClientSecret: "test-secret",
										},
									},
								},
							},
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

	authUrl := data["urlWithQueryParams"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-1", parsedURL.Query().Get("client_id"))
}

func TestSingleConfigWithoutClientTypeSpecifiedOnlyInConfig(t *testing.T) {
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
						Providers: []tpmodels.ProviderInput{
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "google",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientType:   "web",
											ClientID:     "client-id-1",
											ClientSecret: "test-secret",
										},
									},
								},
							},
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

	authUrl := data["urlWithQueryParams"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-1", parsedURL.Query().Get("client_id"))
}

func TestSingleConfigWithClientTypeSpecified(t *testing.T) {
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
						Providers: []tpmodels.ProviderInput{
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "google",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientType:   "web",
											ClientID:     "client-id-1",
											ClientSecret: "test-secret",
										},
									},
								},
							},
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=google&clientType=web")
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

	authUrl := data["urlWithQueryParams"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-1", parsedURL.Query().Get("client_id"))
}

func TestSingleConfigWithDifferentClientTypeSpecified(t *testing.T) {
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
						Providers: []tpmodels.ProviderInput{
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "google",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientType:   "web",
											ClientID:     "client-id-1",
											ClientSecret: "test-secret",
										},
									},
								},
							},
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=google&clientType=ios")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 400, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	assert.NoError(t, err)

	assert.Equal(t, data["message"], "Could not find client config for clientType: ios")
}

func TestMultipleProviderSingleConfigWithoutClientTypeSpecified(t *testing.T) {
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
						Providers: []tpmodels.ProviderInput{
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "google",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientID:     "client-id-1",
											ClientSecret: "test-secret",
										},
									},
								},
							},
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "facebook",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientID:     "client-id-2",
											ClientSecret: "test-secret",
										},
									},
								},
							},
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

	authUrl := data["urlWithQueryParams"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-1", parsedURL.Query().Get("client_id"))
}

func TestMultipleConfigWithoutClientTypeSpecified(t *testing.T) {
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
						Providers: []tpmodels.ProviderInput{
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "google",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientType:   "web",
											ClientID:     "client-id-1",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "ios",
											ClientID:     "client-id-2",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "android",
											ClientID:     "client-id-3",
											ClientSecret: "test-secret",
										},
									},
								},
							},
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
	assert.Equal(t, 400, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	assert.NoError(t, err)

	assert.Equal(t, data["message"], "please provide exactly one client config or pass clientType or tenantId")
}

func TestMultipleConfigWithClientTypeSpecified(t *testing.T) {
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
						Providers: []tpmodels.ProviderInput{
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "google",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientType:   "web",
											ClientID:     "client-id-1",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "ios",
											ClientID:     "client-id-2",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "android",
											ClientID:     "client-id-3",
											ClientSecret: "test-secret",
										},
									},
								},
							},
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=google&clientType=ios")
	if err != nil {
		t.Error(err.Error())
	}
	// TODO this will result in an error
	assert.Equal(t, 200, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	assert.NoError(t, err)

	authUrl := data["urlWithQueryParams"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-2", parsedURL.Query().Get("client_id"))
}

func TestMultipleProviderMultipleConfigWithoutClientTypeSpecified(t *testing.T) {
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
						Providers: []tpmodels.ProviderInput{
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "google",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientType:   "web",
											ClientID:     "client-id-1",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "ios",
											ClientID:     "client-id-2",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "android",
											ClientID:     "client-id-3",
											ClientSecret: "test-secret",
										},
									},
								},
							},
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "facebook",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientType:   "web",
											ClientID:     "client-id-1",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "ios",
											ClientID:     "client-id-2",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "android",
											ClientID:     "client-id-3",
											ClientSecret: "test-secret",
										},
									},
								},
							},
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
	assert.Equal(t, 400, resp.StatusCode)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	assert.NoError(t, err)

	assert.Equal(t, data["message"], "please provide exactly one client config or pass clientType or tenantId")
}

func TestMultipleProviderMultipleConfigWithClientTypeSpecified(t *testing.T) {
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
						Providers: []tpmodels.ProviderInput{
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "google",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientType:   "web",
											ClientID:     "client-id-1",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "ios",
											ClientID:     "client-id-2",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "android",
											ClientID:     "client-id-3",
											ClientSecret: "test-secret",
										},
									},
								},
							},
							{
								Config: tpmodels.ProviderConfig{
									ThirdPartyId: "facebook",
									Clients: []tpmodels.ProviderClientConfig{
										{
											ClientType:   "web",
											ClientID:     "client-id-4",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "ios",
											ClientID:     "client-id-5",
											ClientSecret: "test-secret",
										},
										{
											ClientType:   "android",
											ClientID:     "client-id-6",
											ClientSecret: "test-secret",
										},
									},
								},
							},
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=facebook&clientType=android")
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

	authUrl := data["urlWithQueryParams"].(string)
	parsedURL, err := url.Parse(authUrl)
	assert.NoError(t, err)

	assert.Equal(t, "client-id-6", parsedURL.Query().Get("client_id"))
}
