/* Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultUserContext(t *testing.T) {
	signInContextWorks := false
	signInAPIContextWorks := false
	createNewSessionContextWorks := false

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
			Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					Functions: func(originalImplementation epmodels.RecipeInterface) epmodels.RecipeInterface {
						originalSignIn := *originalImplementation.SignIn
						newSignIn := func(email string, password string, userContext supertokens.UserContext) (epmodels.SignInResponse, error) {
							if _default, ok := (*userContext)["_default"].(map[string]interface{}); ok {
								if _, ok := _default["request"].(*http.Request); ok {
									signInContextWorks = true
								}
							}
							return originalSignIn(email, password, userContext)
						}
						*originalImplementation.SignIn = newSignIn
						return originalImplementation
					},

					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						originalSignInPOST := *originalImplementation.SignInPOST
						newSignInPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignInPOSTResponse, error) {
							if _default, ok := (*userContext)["_default"].(map[string]interface{}); ok {
								if _, ok := _default["request"].(*http.Request); ok {
									signInAPIContextWorks = true
								}
							}
							return originalSignInPOST(formFields, options, userContext)
						}
						*originalImplementation.SignInPOST = newSignInPOST
						return originalImplementation
					},
				},
			}),
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},

				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						originalCreateNewSession := *originalImplementation.CreateNewSession
						newCreateNewSession := func(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							if _default, ok := (*userContext)["_default"].(map[string]interface{}); ok {
								if _, ok := _default["request"].(*http.Request); ok {
									createNewSessionContextWorks = true
								}
							}
							return originalCreateNewSession(req, res, userID, accessTokenPayload, sessionData, userContext)
						}
						*originalImplementation.CreateNewSession = newCreateNewSession
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)

	res1, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 200, res1.StatusCode)
	assert.True(t, signInContextWorks)
	assert.True(t, signInAPIContextWorks)
	assert.True(t, createNewSessionContextWorks)
}
