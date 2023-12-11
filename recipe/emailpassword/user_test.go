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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestGetUsersOldestFirst(t *testing.T) {
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
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	_, err = unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test1@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test2@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test3@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test4@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	users, err := supertokens.GetUsersOldestFirst("public", nil, nil, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, len(users.Users), 5)
	assert.Nil(t, users.NextPaginationToken)

	limit := 1
	users, err = supertokens.GetUsersOldestFirst("public", nil, &limit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, len(users.Users), 1)
	assert.NotNil(t, users.NextPaginationToken)
	assert.Equal(t, "test@gmail.com", users.Users[0].Emails[0])

	users, err = supertokens.GetUsersOldestFirst("public", users.NextPaginationToken, &limit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, len(users.Users), 1)
	assert.NotNil(t, users.NextPaginationToken)
	assert.Equal(t, "test1@gmail.com", users.Users[0].Emails[0])

	limit = 5
	users, err = supertokens.GetUsersOldestFirst("public", users.NextPaginationToken, &limit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, len(users.Users), 3)
	assert.Nil(t, users.NextPaginationToken)

	customPaginationToken := "invalid-pagination-token"
	limit = 10
	users, err = supertokens.GetUsersOldestFirst("public", &customPaginationToken, &limit, nil, nil)
	if err != nil {
		assert.Equal(t, "SuperTokens core threw an error for a request to path: '/public/users' with status code: 400 and message: invalid pagination token\n", err.Error())
	} else {
		t.Fail()
	}

	limit = -1
	users, err = supertokens.GetUsersOldestFirst("public", nil, &limit, nil, nil)
	if err != nil {
		assert.Equal(t, "SuperTokens core threw an error for a request to path: '/public/users' with status code: 400 and message: limit must a positive integer with min value 1\n", err.Error())
	} else {
		t.Fail()
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
			Init(nil),
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
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	_, err = unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test1@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test2@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test3@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test4@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	users, err := supertokens.GetUsersNewestFirst("public", nil, nil, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, len(users.Users), 5)
	assert.Nil(t, users.NextPaginationToken)

	limit := 1
	users, err = supertokens.GetUsersNewestFirst("public", nil, &limit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, len(users.Users), 1)
	assert.NotNil(t, users.NextPaginationToken)
	assert.Equal(t, "test4@gmail.com", users.Users[0].Emails[0])

	users, err = supertokens.GetUsersNewestFirst("public", users.NextPaginationToken, &limit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, len(users.Users), 1)
	assert.NotNil(t, users.NextPaginationToken)
	assert.Equal(t, "test3@gmail.com", users.Users[0].Emails[0])

	limit = 5
	users, err = supertokens.GetUsersNewestFirst("public", users.NextPaginationToken, &limit, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, len(users.Users), 3)
	assert.Nil(t, users.NextPaginationToken)

	customPaginationToken := "invalid-pagination-token"
	limit = 10
	users, err = supertokens.GetUsersNewestFirst("public", &customPaginationToken, &limit, nil, nil)
	if err != nil {
		assert.Equal(t, "SuperTokens core threw an error for a request to path: '/public/users' with status code: 400 and message: invalid pagination token\n", err.Error())
	} else {
		t.Fail()
	}

	limit = -1
	users, err = supertokens.GetUsersNewestFirst("public", nil, &limit, nil, nil)
	if err != nil {
		assert.Equal(t, "SuperTokens core threw an error for a request to path: '/public/users' with status code: 400 and message: limit must a positive integer with min value 1\n", err.Error())
	} else {
		t.Fail()
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
			Init(nil),
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
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	_, err = unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	userCount, err := supertokens.GetUserCount(nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, float64(1), userCount)

	_, err = unittesting.SignupRequest("test1@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test2@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test3@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = unittesting.SignupRequest("test4@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	userCount, err = supertokens.GetUserCount(nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, float64(5), userCount)
}
