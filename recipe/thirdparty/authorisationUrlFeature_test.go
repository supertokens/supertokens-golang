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
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestReqWithThirdPartyEmailPasswordRecipe(t *testing.T) {
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
			emailpassword.Init(nil),
			Init(
				&tpmodels.TypeInput{
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

	client := &http.Client{}
	req, _ := http.NewRequest("GET", testServer.URL+"/auth/authorisationurl?thirdPartyId=google", nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "thirdpartyemailpassword")

	resp, err := client.Do(req)

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

	assert.Equal(t, "OK", data["status"])

	fetchedUrl, err := url.Parse(data["urlWithQueryParams"].(string))
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "supertokens.io", fetchedUrl.Host)
	assert.Equal(t, "/dev/oauth/redirect-to-provider", fetchedUrl.Path)
}

func TestReqWithThirdPartyEmailPasswordRecipe2(t *testing.T) {
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
											ClientID:     "4398792-test-id",
											ClientSecret: "test-secret",
										},
									},
								},
							},
						},
					},
				},
			),
			emailpassword.Init(nil),
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

	client := &http.Client{}
	req, _ := http.NewRequest("GET", testServer.URL+"/auth/authorisationurl?thirdPartyId=google", nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "thirdpartyemailpassword")

	resp, err := client.Do(req)

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

	assert.Equal(t, "OK", data["status"])

	fetchedUrl, err := url.Parse(data["urlWithQueryParams"].(string))
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "supertokens.io", fetchedUrl.Host)
	assert.Equal(t, "/dev/oauth/redirect-to-provider", fetchedUrl.Path)
}

func TestUsingDevOAuthKeysWillUseDevAuthUrl(t *testing.T) {
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
											ClientID:     "4398792-test-id",
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

	assert.Equal(t, "OK", data["status"])

	fetchedUrl, err := url.Parse(data["urlWithQueryParams"].(string))
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "supertokens.io", fetchedUrl.Host)
	assert.Equal(t, "/dev/oauth/redirect-to-provider", fetchedUrl.Path)
}

func TestMinimumConfigForThirdPartyModule(t *testing.T) {
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
							unittesting.ReturnCustomProviderWithAuthRedirectParams(),
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl?thirdPartyId=custom&dynamic=example.com")
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

	assert.Equal(t, "OK", data["status"])

	fetchedUrl, err := url.Parse(data["urlWithQueryParams"].(string))
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "supertokens", fetchedUrl.Query()["client_id"][0])
	assert.Equal(t, "test", fetchedUrl.Query()["scope"][0])
	assert.Equal(t, "example.com", fetchedUrl.Query()["dynamic"][0])
}

func TestThirdPartyProviderDoesnotExist(t *testing.T) {
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
							unittesting.ReturnCustomProviderWithAuthRedirectParams(),
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

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	assert.Equal(t, "the provider google could not be found in the configuration", data["message"])
}

func TestInvalidGetParamsForThirdPartyModule(t *testing.T) {
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
							unittesting.ReturnCustomProviderWithAuthRedirectParams(),
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

	resp, err := http.Get(testServer.URL + "/auth/authorisationurl")
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

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Please provide the thirdPartyId as a GET param", data["message"])
}
