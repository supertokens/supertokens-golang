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
	"testing"
)

func TestGetUsersOldesFirst(t *testing.T) {
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
	// 		session.Init(nil),
	// 		Init(
	// 			&tpmodels.TypeInput{
	// 				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
	// 					Providers: []tpmodels.TypeProvider{
	// 						customProvider2,
	// 					},
	// 				},
	// 			},
	// 		),
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

	// unittesting.SigninupCustomRequest(testServer.URL, "test@gmail.com", "testPass0")
	// unittesting.SigninupCustomRequest(testServer.URL, "test1@gmail.com", "testPass1")
	// unittesting.SigninupCustomRequest(testServer.URL, "test2@gmail.com", "testPass2")
	// unittesting.SigninupCustomRequest(testServer.URL, "test3@gmail.com", "testPass3")
	// unittesting.SigninupCustomRequest(testServer.URL, "test4@gmail.com", "testPass4")

	// userPaginationResult, err := supertokens.GetUsersOldestFirst(nil, nil, nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// assert.Equal(t, 5, len(userPaginationResult.Users))
	// assert.Nil(t, userPaginationResult.NextPaginationToken)

	// customLimit := 1
	// userPaginationResult, err = supertokens.GetUsersOldestFirst(nil, &customLimit, nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// assert.Equal(t, 1, len(userPaginationResult.Users))
	// assert.Equal(t, "test@gmail.com", userPaginationResult.Users[0].User["email"])
	// assert.Equal(t, "*string", reflect.TypeOf(userPaginationResult.NextPaginationToken).String())

	// userPaginationResult, err = supertokens.GetUsersOldestFirst(userPaginationResult.NextPaginationToken, &customLimit, nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// assert.Equal(t, 1, len(userPaginationResult.Users))
	// assert.Equal(t, "test1@gmail.com", userPaginationResult.Users[0].User["email"])
	// assert.Equal(t, "*string", reflect.TypeOf(userPaginationResult.NextPaginationToken).String())

	// customLimit = 5
	// userPaginationResult, err = supertokens.GetUsersOldestFirst(userPaginationResult.NextPaginationToken, &customLimit, nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// assert.Equal(t, 3, len(userPaginationResult.Users))

	// customInvalidPaginationToken := "invalid-pagination-token"
	// userPaginationResult, err = supertokens.GetUsersOldestFirst(&customInvalidPaginationToken, &customLimit, nil)
	// if err != nil {
	// 	assert.Contains(t, err.Error(), "invalid pagination token")
	// } else {
	// 	t.Fail()
	// }

	// customLimit = -1
	// userPaginationResult, err = supertokens.GetUsersOldestFirst(nil, &customLimit, nil)
	// if err != nil {
	// 	assert.Contains(t, err.Error(), "limit must a positive integer with min value 1")
	// } else {
	// 	t.Fail()
	// }
}

func TestGetUsersNewestFirst(t *testing.T) {
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
	// 		session.Init(nil),
	// 		Init(
	// 			&tpmodels.TypeInput{
	// 				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
	// 					Providers: []tpmodels.TypeProvider{
	// 						customProvider2,
	// 					},
	// 				},
	// 			},
	// 		),
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

	// unittesting.SigninupCustomRequest(testServer.URL, "test@gmail.com", "testPass0")
	// unittesting.SigninupCustomRequest(testServer.URL, "test1@gmail.com", "testPass1")
	// unittesting.SigninupCustomRequest(testServer.URL, "test2@gmail.com", "testPass2")
	// unittesting.SigninupCustomRequest(testServer.URL, "test3@gmail.com", "testPass3")
	// unittesting.SigninupCustomRequest(testServer.URL, "test4@gmail.com", "testPass4")

	// userPaginationResult, err := supertokens.GetUsersNewestFirst(nil, nil, nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// assert.Equal(t, 5, len(userPaginationResult.Users))
	// assert.Nil(t, userPaginationResult.NextPaginationToken)

	// customLimit := 1
	// userPaginationResult, err = supertokens.GetUsersNewestFirst(nil, &customLimit, nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// assert.Equal(t, 1, len(userPaginationResult.Users))
	// assert.Equal(t, "test4@gmail.com", userPaginationResult.Users[0].User["email"])
	// assert.Equal(t, "*string", reflect.TypeOf(userPaginationResult.NextPaginationToken).String())

	// userPaginationResult, err = supertokens.GetUsersNewestFirst(userPaginationResult.NextPaginationToken, &customLimit, nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// assert.Equal(t, 1, len(userPaginationResult.Users))
	// assert.Equal(t, "test3@gmail.com", userPaginationResult.Users[0].User["email"])
	// assert.Equal(t, "*string", reflect.TypeOf(userPaginationResult.NextPaginationToken).String())

	// customLimit = 5
	// userPaginationResult, err = supertokens.GetUsersNewestFirst(userPaginationResult.NextPaginationToken, &customLimit, nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// assert.Equal(t, 3, len(userPaginationResult.Users))

	// customInvalidPaginationToken := "invalid-pagination-token"
	// customLimit = 10
	// userPaginationResult, err = supertokens.GetUsersNewestFirst(&customInvalidPaginationToken, &customLimit, nil)
	// if err != nil {
	// 	assert.Contains(t, err.Error(), "invalid pagination token")
	// } else {
	// 	t.Fail()
	// }

	// customLimit = -1
	// userPaginationResult, err = supertokens.GetUsersNewestFirst(nil, &customLimit, nil)
	// if err != nil {
	// 	assert.Contains(t, err.Error(), "limit must a positive integer with min value 1")
	// } else {
	// 	t.Fail()
	// }
}

func TestGetUserCount(t *testing.T) {
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
	// 		session.Init(nil),
	// 		Init(
	// 			&tpmodels.TypeInput{
	// 				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
	// 					Providers: []tpmodels.TypeProvider{
	// 						customProvider2,
	// 					},
	// 				},
	// 			},
	// 		),
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

	// userCount, err := supertokens.GetUserCount(nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }

	// assert.Equal(t, 0.0, userCount)

	// unittesting.SigninupCustomRequest(testServer.URL, "test@gmail.com", "testPass0")

	// userCount, err = supertokens.GetUserCount(nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }

	// assert.Equal(t, 1.0, userCount)

	// unittesting.SigninupCustomRequest(testServer.URL, "test1@gmail.com", "testPass1")
	// unittesting.SigninupCustomRequest(testServer.URL, "test2@gmail.com", "testPass2")
	// unittesting.SigninupCustomRequest(testServer.URL, "test3@gmail.com", "testPass3")
	// unittesting.SigninupCustomRequest(testServer.URL, "test4@gmail.com", "testPass4")

	// userCount, err = supertokens.GetUserCount(nil)
	// if err != nil {
	// 	t.Error(err.Error())
	// }

	// assert.Equal(t, 5.0, userCount)
}
