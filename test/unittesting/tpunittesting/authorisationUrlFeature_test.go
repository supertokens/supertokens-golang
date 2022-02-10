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

package tpunittesting

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

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
			thirdparty.Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							thirdparty.Google(
								tpmodels.GoogleConfig{
									ClientID:     "4398792-test-id",
									ClientSecret: "test-secret",
								},
							),
						},
					},
				},
			),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

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

	fetchedUrl, err := url.Parse(data["url"].(string))
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "supertokens.io", fetchedUrl.Host)
	assert.Equal(t, "/dev/oauth/redirect-to-provider", fetchedUrl.Path)

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			thirdparty.Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							unittesting.ReturnCustomProviderWithAuthRedirectParams(),
						},
					},
				},
			),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

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

	fetchedUrl, err := url.Parse(data["url"].(string))
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "supertokens", fetchedUrl.Query()["client_id"][0])
	assert.Equal(t, "test", fetchedUrl.Query()["scope"][0])
	assert.Equal(t, "example.com", fetchedUrl.Query()["dynamic"][0])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			thirdparty.Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							unittesting.ReturnCustomProviderWithAuthRedirectParams(),
						},
					},
				},
			),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

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

	assert.Equal(t, "The third party provider google seems to not be missing from the backend configs", data["message"])

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
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
			thirdparty.Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.TypeProvider{
							unittesting.ReturnCustomProviderWithAuthRedirectParams(),
						},
					},
				},
			),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

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

	defer unittesting.AfterEach()
	defer func() {
		testServer.Close()
	}()
}
