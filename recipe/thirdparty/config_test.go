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

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
// 	"github.com/supertokens/supertokens-golang/supertokens"
// 	"github.com/supertokens/supertokens-golang/test/unittesting"
// )

// func TestConfigForValidInputForThirdPartyModule(t *testing.T) {
// 	configValue := supertokens.TypeInput{
// 		Supertokens: &supertokens.ConnectionInfo{
// 			ConnectionURI: "http://localhost:8080",
// 		},
// 		AppInfo: supertokens.AppInfo{
// 			APIDomain:     "api.supertokens.io",
// 			AppName:       "SuperTokens",
// 			WebsiteDomain: "supertokens.io",
// 		},
// 		RecipeList: []supertokens.Recipe{
// 			Init(&tpmodels.TypeInput{}),
// 		},
// 	}

// 	BeforeEach()
// 	unittesting.StartUpST("localhost", "8080")
// 	defer AfterEach()
// 	err := supertokens.Init(configValue)

// 	if err != nil {
// 		assert.Equal(t, "thirdparty recipe requires at least 1 provider to be passed in signInAndUpFeature.providers config", err.Error())
// 	} else {
// 		t.Fail()
// 	}
// }

// func TestConfigForInValidInputWithEmptyProviderSliceForThirdPartyModule(t *testing.T) {
// 	configValue := supertokens.TypeInput{
// 		Supertokens: &supertokens.ConnectionInfo{
// 			ConnectionURI: "http://localhost:8080",
// 		},
// 		AppInfo: supertokens.AppInfo{
// 			APIDomain:     "api.supertokens.io",
// 			AppName:       "SuperTokens",
// 			WebsiteDomain: "supertokens.io",
// 		},
// 		RecipeList: []supertokens.Recipe{
// 			Init(
// 				&tpmodels.TypeInput{
// 					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
// 						Providers: []tpmodels.TypeProvider{},
// 					},
// 				},
// 			),
// 		},
// 	}

// 	BeforeEach()
// 	unittesting.StartUpST("localhost", "8080")
// 	defer AfterEach()
// 	err := supertokens.Init(configValue)

// 	if err != nil {
// 		assert.Equal(t, "thirdparty recipe requires at least 1 provider to be passed in signInAndUpFeature.providers config", err.Error())
// 	} else {
// 		t.Fail()
// 	}
// }

// func TestMinimumConfigForThirdpartyModuleCustomProvider(t *testing.T) {
// 	configValue := supertokens.TypeInput{
// 		Supertokens: &supertokens.ConnectionInfo{
// 			ConnectionURI: "http://localhost:8080",
// 		},
// 		AppInfo: supertokens.AppInfo{
// 			APIDomain:     "api.supertokens.io",
// 			AppName:       "SuperTokens",
// 			WebsiteDomain: "supertokens.io",
// 		},
// 		RecipeList: []supertokens.Recipe{
// 			Init(
// 				&tpmodels.TypeInput{
// 					SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
// 						Providers: []tpmodels.TypeProvider{
// 							{
// 								ID: "custom",
// 								Get: func(redirectURI, authCodeFromRequest *string, userContext *map[string]interface{}) tpmodels.TypeProviderGetResponse {
// 									return tpmodels.TypeProviderGetResponse{
// 										AccessTokenAPI: tpmodels.AccessTokenAPI{
// 											URL: "test.com/oauth/token",
// 										},
// 										AuthorisationRedirect: tpmodels.AuthorisationRedirect{
// 											URL: "test.com/oauth/auth",
// 										},
// 										GetProfileInfo: func(authCodeResponse interface{}, userContext *map[string]interface{}) (tpmodels.UserInfo, error) {
// 											return tpmodels.UserInfo{
// 												ID: "user",
// 												Email: &tpmodels.EmailStruct{
// 													ID:         "email@test.com",
// 													IsVerified: true,
// 												},
// 											}, nil
// 										},
// 									}
// 								},
// 							},
// 						},
// 					},
// 				},
// 			),
// 		},
// 	}

// 	BeforeEach()
// 	unittesting.StartUpST("localhost", "8080")
// 	defer AfterEach()
// 	err := supertokens.Init(configValue)

// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// }
