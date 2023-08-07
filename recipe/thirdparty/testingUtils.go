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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func resetAll() {
	supertokens.ResetForTest()
	ResetForTest()
	emailverification.ResetForTest()
	session.ResetForTest()
	multitenancy.ResetForTest()
}

func BeforeEach() {
	unittesting.KillAllST()
	resetAll()
	unittesting.SetUpST()
}

func AfterEach() {
	unittesting.KillAllST()
	resetAll()
	unittesting.CleanST()
}

type PostDataForCustomProvider struct {
	ThirdPartyId    string `json:"thirdPartyId"`
	RedirectURIInfo *struct {
		RedirectURIOnProviderDashboard string                 `json:"redirectURIOnProviderDashboard"`
		RedirectURIQueryParams         map[string]interface{} `json:"redirectURIQueryParams"`
	} `json:"redirectURIInfo,omitempty"`
	OAuthTokens map[string]interface{} `json:"oAuthTokens,omitempty"`
}

func supertokensInitForTest(t *testing.T, recipes ...supertokens.Recipe) *httptest.Server {
	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: recipes,
	}

	err := supertokens.Init(config)
	assert.NoError(t, err)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	return testServer
}

var customProviderForEmailVerification = tpmodels.ProviderInput{
	Config: tpmodels.ProviderConfig{
		ThirdPartyId:          "custom",
		AuthorizationEndpoint: "https://test.com/oauth/auth",
		TokenEndpoint:         "https://test.com/oauth/token",

		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "supertokens",
			},
		},
	},

	Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return tpmodels.TypeUserInfo{
				ThirdPartyUserId: "user",
				Email: &tpmodels.EmailStruct{
					ID:         "test@example.com",
					IsVerified: false,
				},
			}, nil
		}
		return originalImplementation
	},
}

var customProvider6 = tpmodels.ProviderInput{
	Config: tpmodels.ProviderConfig{
		ThirdPartyId:          "custom",
		AuthorizationEndpoint: "https://test.com/oauth/auth",
		TokenEndpoint:         "https://test.com/oauth/token",

		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "supertokens",
			},
		},
	},

	Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			if oAuthTokens["access_token"] == nil {
				return tpmodels.TypeUserInfo{}, nil
			}
			return tpmodels.TypeUserInfo{
				ThirdPartyUserId: "user",
				Email: &tpmodels.EmailStruct{
					ID:         "email@test.com",
					IsVerified: true,
				},
			}, nil
		}
		return originalImplementation
	},
}

var customProvider1 = tpmodels.ProviderInput{
	Config: tpmodels.ProviderConfig{
		ThirdPartyId:          "custom",
		AuthorizationEndpoint: "https://test.com/oauth/auth",
		TokenEndpoint:         "https://test.com/oauth/token",

		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "supertokens",
			},
		},
	},

	Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return tpmodels.TypeUserInfo{
				ThirdPartyUserId: "user",
				Email: &tpmodels.EmailStruct{
					ID:         "email@test.com",
					IsVerified: true,
				},
			}, nil
		}
		return originalImplementation
	},
}

var customProvider2 = tpmodels.ProviderInput{
	Config: tpmodels.ProviderConfig{
		ThirdPartyId:          "custom",
		AuthorizationEndpoint: "https://test.com/oauth/auth",
		TokenEndpoint:         "https://test.com/oauth/token",

		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "supertokens",
			},
		},
	},

	Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return tpmodels.TypeUserInfo{
				ThirdPartyUserId: oAuthTokens["id"].(string),
				Email: &tpmodels.EmailStruct{
					ID:         oAuthTokens["email"].(string),
					IsVerified: true,
				},
			}, nil
		}
		return originalImplementation
	},
}

var customProvider5 = tpmodels.ProviderInput{
	Config: tpmodels.ProviderConfig{
		ThirdPartyId:          "custom",
		AuthorizationEndpoint: "https://test.com/oauth/auth",
		TokenEndpoint:         "https://test.com/oauth/token",

		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "supertokens",
			},
		},
	},

	Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return tpmodels.TypeUserInfo{
				ThirdPartyUserId: "user",
				Email: &tpmodels.EmailStruct{
					ID:         "email@test.com",
					IsVerified: false,
				},
			}, nil
		}
		return originalImplementation
	},
}

var customProvider3 = tpmodels.ProviderInput{
	Config: tpmodels.ProviderConfig{
		ThirdPartyId:          "custom",
		AuthorizationEndpoint: "https://test.com/oauth/auth",
		TokenEndpoint:         "https://test.com/oauth/token",

		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "supertokens",
			},
		},
	},

	Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return tpmodels.TypeUserInfo{
				ThirdPartyUserId: "user",
			}, nil
		}
		return originalImplementation
	},
}
