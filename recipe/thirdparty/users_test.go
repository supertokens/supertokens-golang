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

package thirdparty

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestGetUsersOldesFirst(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.ProviderInput{
							customProvider2,
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

	unittesting.SigninupCustomRequest(testServer.URL, "test@gmail.com", "testPass0")
	unittesting.SigninupCustomRequest(testServer.URL, "test1@gmail.com", "testPass1")
	unittesting.SigninupCustomRequest(testServer.URL, "john@gmail.com", "testPass2")
	unittesting.SigninupCustomRequest(testServer.URL, "test3@gmail.com", "testPass3")
	unittesting.SigninupCustomRequest(testServer.URL, "test4@gmail.com", "testPass4")

	userPaginationResult, err := supertokens.GetUsersOldestFirst("public", nil, nil, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 5, len(userPaginationResult.Users))
	assert.Nil(t, userPaginationResult.NextPaginationToken)

	customLimit := 1
	userPaginationResult, err = supertokens.GetUsersOldestFirst("public", nil, &customLimit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 1, len(userPaginationResult.Users))
	assert.Equal(t, "test@gmail.com", userPaginationResult.Users[0].User["email"])
	assert.Equal(t, "*string", reflect.TypeOf(userPaginationResult.NextPaginationToken).String())

	userPaginationResult, err = supertokens.GetUsersOldestFirst("public", userPaginationResult.NextPaginationToken, &customLimit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 1, len(userPaginationResult.Users))
	assert.Equal(t, "test1@gmail.com", userPaginationResult.Users[0].User["email"])
	assert.Equal(t, "*string", reflect.TypeOf(userPaginationResult.NextPaginationToken).String())

	customLimit = 5
	userPaginationResult, err = supertokens.GetUsersOldestFirst("public", userPaginationResult.NextPaginationToken, &customLimit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 3, len(userPaginationResult.Users))

	customInvalidPaginationToken := "invalid-pagination-token"
	userPaginationResult, err = supertokens.GetUsersOldestFirst("public", &customInvalidPaginationToken, &customLimit, nil, nil)
	if err != nil {
		assert.Contains(t, err.Error(), "invalid pagination token")
	} else {
		t.Fail()
	}

	customLimit = -1
	userPaginationResult, err = supertokens.GetUsersOldestFirst("public", nil, &customLimit, nil, nil)
	if err != nil {
		assert.Contains(t, err.Error(), "limit must a positive integer with min value 1")
	} else {
		t.Fail()
	}

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("thirdparty")
	if err != nil {
		t.Fail()
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Fail()
	}

	if supertokens.MaxVersion(cdiVersion, "2.20") != cdiVersion {
		t.Skip()
	}

	customLimit = 10
	query := make(map[string]string)
	query["email"] = "doe"
	userPaginationResult, err = supertokens.GetUsersOldestFirst("public", nil, &customLimit, nil, query)
	if err != nil {
		t.Fail()
	} else {
		assert.Equal(t, len(userPaginationResult.Users), 0)
	}

	query["email"] = "john"
	userPaginationResult, err = supertokens.GetUsersOldestFirst("public", nil, &customLimit, nil, query)
	if err != nil {
		t.Fail()
	} else {
		assert.Equal(t, len(userPaginationResult.Users), 1)
	}
}

func TestGetUsersNewestFirst(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.ProviderInput{
							customProvider2,
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

	unittesting.SigninupCustomRequest(testServer.URL, "test@gmail.com", "testPass0")
	unittesting.SigninupCustomRequest(testServer.URL, "test1@gmail.com", "testPass1")
	unittesting.SigninupCustomRequest(testServer.URL, "john@gmail.com", "testPass2")
	unittesting.SigninupCustomRequest(testServer.URL, "test3@gmail.com", "testPass3")
	unittesting.SigninupCustomRequest(testServer.URL, "test4@gmail.com", "testPass4")

	userPaginationResult, err := supertokens.GetUsersNewestFirst("public", nil, nil, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 5, len(userPaginationResult.Users))
	assert.Nil(t, userPaginationResult.NextPaginationToken)

	customLimit := 1
	userPaginationResult, err = supertokens.GetUsersNewestFirst("public", nil, &customLimit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 1, len(userPaginationResult.Users))
	assert.Equal(t, "test4@gmail.com", userPaginationResult.Users[0].User["email"])
	assert.Equal(t, "*string", reflect.TypeOf(userPaginationResult.NextPaginationToken).String())

	userPaginationResult, err = supertokens.GetUsersNewestFirst("public", userPaginationResult.NextPaginationToken, &customLimit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 1, len(userPaginationResult.Users))
	assert.Equal(t, "test3@gmail.com", userPaginationResult.Users[0].User["email"])
	assert.Equal(t, "*string", reflect.TypeOf(userPaginationResult.NextPaginationToken).String())

	customLimit = 5
	userPaginationResult, err = supertokens.GetUsersNewestFirst("public", userPaginationResult.NextPaginationToken, &customLimit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 3, len(userPaginationResult.Users))

	customInvalidPaginationToken := "invalid-pagination-token"
	customLimit = 10
	userPaginationResult, err = supertokens.GetUsersNewestFirst("public", &customInvalidPaginationToken, &customLimit, nil, nil)
	if err != nil {
		assert.Contains(t, err.Error(), "invalid pagination token")
	} else {
		t.Fail()
	}

	customLimit = -1
	userPaginationResult, err = supertokens.GetUsersNewestFirst("public", nil, &customLimit, nil, nil)
	if err != nil {
		assert.Contains(t, err.Error(), "limit must a positive integer with min value 1")
	} else {
		t.Fail()
	}

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("thirdparty")
	if err != nil {
		t.Fail()
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Fail()
	}

	if supertokens.MaxVersion(cdiVersion, "2.20") != cdiVersion {
		t.Skip()
	}

	customLimit = 10
	query := make(map[string]string)
	query["email"] = "doe"
	userPaginationResult, err = supertokens.GetUsersNewestFirst("public", nil, &customLimit, nil, query)
	if err != nil {
		t.Fail()
	} else {
		assert.Equal(t, len(userPaginationResult.Users), 0)
	}

	query["email"] = "john"
	userPaginationResult, err = supertokens.GetUsersNewestFirst("public", nil, &customLimit, nil, query)
	if err != nil {
		t.Fail()
	} else {
		assert.Equal(t, len(userPaginationResult.Users), 1)
	}
}

func TestGetUserCount(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(
				&tpmodels.TypeInput{
					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
						Providers: []tpmodels.ProviderInput{
							customProvider2,
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

	userCount, err := supertokens.GetUserCount(nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 0.0, userCount)

	unittesting.SigninupCustomRequest(testServer.URL, "test@gmail.com", "testPass0")

	userCount, err = supertokens.GetUserCount(nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 1.0, userCount)

	unittesting.SigninupCustomRequest(testServer.URL, "test1@gmail.com", "testPass1")
	unittesting.SigninupCustomRequest(testServer.URL, "test2@gmail.com", "testPass2")
	unittesting.SigninupCustomRequest(testServer.URL, "test3@gmail.com", "testPass3")
	unittesting.SigninupCustomRequest(testServer.URL, "test4@gmail.com", "testPass4")

	userCount, err = supertokens.GetUserCount(nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 5.0, userCount)
}
