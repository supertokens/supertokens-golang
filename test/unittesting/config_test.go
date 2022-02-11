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

package unittesting

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestSuperTokensInitWithAbsentOptionalFieldsInAppInfo(t *testing.T) {
	apiBasePath := "test/"
	websiteBasePath := "test1/"
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				APIDomain:     "api.supertokens.io",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:         "SuperTokens",
				APIDomain:       "api.supertokens.io",
				WebsiteDomain:   "supertokens.io",
				APIBasePath:     &apiBasePath,
				WebsiteBasePath: &websiteBasePath,
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
	}
	for _, configValue := range configValues {

		BeforeEach()
		StartUpST("localhost", "8080")
		err := supertokens.Init(configValue)
		if err != nil {
			t.Error("Failed to get a supertokens instance")
		}
		supertokensInstance, err := supertokens.GetInstanceOrThrowError()

		if err != nil {
			t.Error("could not find a supertokens instance")
		}

		if configValue.AppInfo.APIBasePath != nil {
			assert.Equal(t, "/test", supertokensInstance.AppInfo.APIBasePath.GetAsStringDangerous())
			assert.Equal(t, "/test1", supertokensInstance.AppInfo.WebsiteBasePath.GetAsStringDangerous())
		} else {
			assert.Equal(t, "/auth", supertokensInstance.AppInfo.APIBasePath.GetAsStringDangerous())
			assert.Equal(t, "/auth", supertokensInstance.AppInfo.WebsiteBasePath.GetAsStringDangerous())
		}
		AfterEach()
	}
}

func TestSuperTokensInitWithAbsenceOfCompulsoryInputInAppInfo(t *testing.T) {
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:   "SuperTokens",
				APIDomain: "api.supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
	}
	for _, configValue := range configValues {
		BeforeEach()
		StartUpST("localhost", "8080")
		err := supertokens.Init(configValue)
		if err != nil {
			errMessage := err.Error()
			if configValue.AppInfo.AppName != "SuperTokens" {
				assert.Equal(t, errMessage, "Please provide your appName inside the appInfo object when calling supertokens.init")
			} else if configValue.AppInfo.APIDomain != "api.supertokens.io" {
				assert.Equal(t, errMessage, "Please provide your apiDomain inside the appInfo object when calling supertokens.init")
			} else if configValue.AppInfo.WebsiteDomain != "supertokens.io" {
				assert.Equal(t, errMessage, "Please provide your websiteDomain inside the appInfo object when calling supertokens.init")
			}
		}
		AfterEach()
	}
}

func TestSuperTokensInitWith0RecipeModules(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage != "please provide at least one recipe to the supertokens.init function call" {
			t.Errorf(err.Error())
		} else {
			assert.Equal(t, errorMessage, "please provide at least one recipe to the supertokens.init function call")
		}
	}

	defer AfterEach()
}

func TestSuperTokensInitWithConfigForSessionModules(t *testing.T) {
	cookieDomain := "testDomain"
	sessionExpiredStatusCode := 111
	cookieSecure := true
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
			session.Init(
				&sessmodels.TypeInput{
					CookieDomain:             &cookieDomain,
					SessionExpiredStatusCode: &sessionExpiredStatusCode,
					CookieSecure:             &cookieSecure,
				},
			),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Errorf(err.Error())
	}

	sessionSingletonInstance, err := session.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, *sessionSingletonInstance.Config.CookieDomain, "testdomain")
	assert.Equal(t, sessionSingletonInstance.Config.SessionExpiredStatusCode, 111)
	assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, true)
	AfterEach()
}

func TestSuperTokensInitWithConfigForSessionModulesWithSameSiteValueAsLax(t *testing.T) {
	cookieSameSite0 := " Lax "
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
				CookieSameSite: &cookieSameSite0,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Errorf(err.Error())
	}

	defer AfterEach()
}

func TestSuperTokensInitWithConfigForSessionModulesWithSameSiteValueAsNone(t *testing.T) {
	cookieSameSite0 := "None "
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
				CookieSameSite: &cookieSameSite0,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Errorf(err.Error())
	}

	defer AfterEach()
}

func TestSuperTokensInitWithConfigForSessionModulesWithSameSiteValueAsStrict(t *testing.T) {
	cookieSameSite0 := " STRICT "
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
				CookieSameSite: &cookieSameSite0,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Errorf(err.Error())
	}

	defer AfterEach()
}

func TestSuperTokensInitWithConfigForSessionModulesWithSameSiteValueAsRandom(t *testing.T) {
	cookieSameSite0 := " random "
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
				CookieSameSite: &cookieSameSite0,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		if err != nil {
			errorMessage := err.Error()
			if errorMessage != `cookie same site must be one of "strict", "lax", or "none"` {
				t.Errorf(err.Error())
			} else {
				assert.Equal(t, errorMessage, `cookie same site must be one of "strict", "lax", or "none"`)
			}
		}
	}

	defer AfterEach()
}

func TestSuperTokensWithVariousApiBasePath(t *testing.T) {
	apiBasePath0 := "/custom/a"
	apiBasePath1 := "/"
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
				APIBasePath:   &apiBasePath0,
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
				APIBasePath:   &apiBasePath1,
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
	}

	for _, configValue := range configValues {
		BeforeEach()
		StartUpST("localhost", "8080")
		err := supertokens.Init(configValue)
		if err != nil {
			t.Errorf(err.Error())
		}
		sessionSingletonInstance, err := session.GetRecipeInstanceOrThrowError()
		if err != nil {
			t.Errorf(err.Error())
		}
		if configValue.AppInfo.APIBasePath != nil {
			checker := RemoveTrailingSlashFromTheEndOfString(*configValue.AppInfo.APIBasePath) + "/session/refresh"
			assert.Equal(t, sessionSingletonInstance.Config.RefreshTokenPath.GetAsStringDangerous(), checker)
		} else {
			assert.Equal(t, sessionSingletonInstance.Config.RefreshTokenPath.GetAsStringDangerous(), "/auth/session/refresh")
		}
		AfterEach()
	}
}

func TestSuperTokensWithCustomApiKey(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
			APIKey:        "haha",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, *supertokens.QuerierAPIKey, configValue.Supertokens.APIKey)
	AfterEach()
}

func TestSuperTokensInitWithCustomSessionExpiredCodeInSessionRecipe(t *testing.T) {
	customAPIBasePath := "/custom"
	customSessionExpiredCode := 402
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customAPIBasePath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(
				&sessmodels.TypeInput{
					SessionExpiredStatusCode: &customSessionExpiredCode,
				},
			),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Errorf(err.Error())
	}
	sessionSingletonInstance, err := session.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, sessionSingletonInstance.Config.SessionExpiredStatusCode, 402)
	AfterEach()
}

func TestSuperTokensInitWithMultipleHosts(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080;try.supertokens.io;try.supertokens.io:8080;localhost:90",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Errorf(err.Error())
	}
	hosts := supertokens.QuerierHosts

	assert.Equal(t, hosts[0].Domain.GetAsStringDangerous(), "http://localhost:8080")
	assert.Equal(t, hosts[1].Domain.GetAsStringDangerous(), "https://try.supertokens.io")
	assert.Equal(t, hosts[2].Domain.GetAsStringDangerous(), "https://try.supertokens.io:8080")
	assert.Equal(t, hosts[3].Domain.GetAsStringDangerous(), "http://localhost:90")

	AfterEach()
}

// func TestSuperTokensInitWithNoneLaxTrueSessionConfigResults(t *testing.T) {
// 	apiBasePath0 := "test/"
// 	websiteBasePath0 := "test1/"
// 	configValues := []supertokens.TypeInput{
// 		{
// 			Supertokens: &supertokens.ConnectionInfo{
// 				ConnectionURI: "http://localhost:8080",
// 			},
// 			AppInfo: supertokens.AppInfo{
// 				APIDomain:       "https://api.supertokens.io",
// 				AppName:         "SuperTokens",
// 				WebsiteDomain:   "supertokens.io",
// 				APIBasePath:     &apiBasePath0,
// 				WebsiteBasePath: &websiteBasePath0,
// 			},
// 			RecipeList: []supertokens.Recipe{
// 				session.Init(nil),
// 			},
// 		},
// 		{
// 			Supertokens: &supertokens.ConnectionInfo{
// 				ConnectionURI: "http://localhost:8080",
// 			},
// 			AppInfo: supertokens.AppInfo{
// 				APIDomain:       "api.supertokens.io",
// 				AppName:         "SuperTokens",
// 				WebsiteDomain:   "supertokens.io",
// 				APIBasePath:     &apiBasePath0,
// 				WebsiteBasePath: &websiteBasePath0,
// 			},
// 			RecipeList: []supertokens.Recipe{
// 				session.Init(nil),
// 			},
// 		},
// 		{
// 			Supertokens: &supertokens.ConnectionInfo{
// 				ConnectionURI: "http://localhost:8080",
// 			},
// 			AppInfo: supertokens.AppInfo{
// 				APIDomain:       "api.supertokens.co.uk",
// 				AppName:         "SuperTokens",
// 				WebsiteDomain:   "supertokens.co.uk",
// 				APIBasePath:     &apiBasePath0,
// 				WebsiteBasePath: &websiteBasePath0,
// 			},
// 			RecipeList: []supertokens.Recipe{
// 				session.Init(nil),
// 			},
// 		},
// 	}
// 	for _, configValue := range configValues {
// 		BeforeEach()
// 		StartUpST("localhost", "8080")
// 		err := supertokens.Init(configValue)
// 		if err != nil {
// 			t.Error(err.Error())
// 		}
// 		sessionSingletonInstance, err := session.GetRecipeInstanceOrThrowError()
// 		if err != nil {
// 			t.Errorf(err.Error())
// 		}
// 		assert.Equal(t, sessionSingletonInstance.Config.AntiCsrf, "NONE")
// 		assert.Equal(t, sessionSingletonInstance.Config.CookieSameSite, "lax")
// 		assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, true)
// AfterEach()
// 	}
// }

func TestSuperTokensInitWithNoneLaxFalseSessionConfigResults(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:       "127.0.0.1:3000",
				AppName:         "SuperTokens",
				WebsiteDomain:   "127.0.0.1:9000",
				APIBasePath:     &apiBasePath0,
				WebsiteBasePath: &websiteBasePath0,
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
	}
	for _, configValue := range configValues {
		BeforeEach()
		StartUpST("localhost", "8080")
		err := supertokens.Init(configValue)
		if err != nil {
			t.Error(err.Error())
		}
		sessionSingletonInstance, err := session.GetRecipeInstanceOrThrowError()
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, sessionSingletonInstance.Config.AntiCsrf, "NONE")
		assert.Equal(t, sessionSingletonInstance.Config.CookieSameSite, "lax")
		assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, false)
		AfterEach()
	}
}

func TestSuperTokensInitWithCustomHeaderLaxTrueSessionConfigResults(t *testing.T) {
	apiBasePath0 := "/"
	customAntiCsrf := "VIA_CUSTOM_HEADER"
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
				APIBasePath:   &apiBasePath0,
			},
			RecipeList: []supertokens.Recipe{
				session.Init(
					&sessmodels.TypeInput{
						AntiCsrf: &customAntiCsrf,
					},
				),
			},
		},
	}
	for _, configValue := range configValues {
		BeforeEach()
		StartUpST("localhost", "8080")
		err := supertokens.Init(configValue)
		if err != nil {
			t.Error(err.Error())
		}
		sessionSingletonInstance, err := session.GetRecipeInstanceOrThrowError()
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, sessionSingletonInstance.Config.AntiCsrf, "VIA_CUSTOM_HEADER")
		assert.Equal(t, sessionSingletonInstance.Config.CookieSameSite, "lax")
		assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, true)
		AfterEach()
	}
}

func TestSuperTokensInitWithCustomHeaderLaxFalseSessionConfigResults(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	customAntiCsrf := "VIA_CUSTOM_HEADER"
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:       "127.0.0.1:3000",
				AppName:         "SuperTokens",
				WebsiteDomain:   "127.0.0.1:9000",
				APIBasePath:     &apiBasePath0,
				WebsiteBasePath: &websiteBasePath0,
			},
			RecipeList: []supertokens.Recipe{
				session.Init(
					&sessmodels.TypeInput{
						AntiCsrf: &customAntiCsrf,
					},
				),
			},
		},
	}
	for _, configValue := range configValues {
		BeforeEach()
		StartUpST("localhost", "8080")
		err := supertokens.Init(configValue)
		if err != nil {
			t.Error(err.Error())
		}
		sessionSingletonInstance, err := session.GetRecipeInstanceOrThrowError()
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, sessionSingletonInstance.Config.AntiCsrf, "VIA_CUSTOM_HEADER")
		assert.Equal(t, sessionSingletonInstance.Config.CookieSameSite, "lax")
		assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, false)
		AfterEach()
	}
}

func TestSuperTokensInitWithCustomHeaderNoneTrueSessionConfigResults(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:       "api.supertokens.com",
				AppName:         "SuperTokens",
				WebsiteDomain:   "supertokens.io",
				APIBasePath:     &apiBasePath0,
				WebsiteBasePath: &websiteBasePath0,
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:       "api.supertokens.com",
				AppName:         "SuperTokens",
				WebsiteDomain:   "127.0.0.1:9000",
				APIBasePath:     &apiBasePath0,
				WebsiteBasePath: &websiteBasePath0,
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
	}
	for _, configValue := range configValues {
		BeforeEach()
		StartUpST("localhost", "8080")
		err := supertokens.Init(configValue)
		if err != nil {
			t.Error(err.Error())
		}
		sessionSingletonInstance, err := session.GetRecipeInstanceOrThrowError()
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, sessionSingletonInstance.Config.AntiCsrf, "VIA_CUSTOM_HEADER")
		assert.Equal(t, sessionSingletonInstance.Config.CookieSameSite, "none")
		assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, true)
		AfterEach()
	}
}

func TestSuperTokensWithAntiCSRFNone(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	customAntiCsrfVal := "NONE"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:       "127.0.0.1:3000",
			AppName:         "SuperTokens",
			WebsiteDomain:   "google.com",
			APIBasePath:     &apiBasePath0,
			WebsiteBasePath: &websiteBasePath0,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(
				&sessmodels.TypeInput{
					AntiCsrf: &customAntiCsrfVal,
				},
			),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletoneSessionRecipeInstance, err := session.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, singletoneSessionRecipeInstance.Config.AntiCsrf, "NONE")
	AfterEach()
}

func TestSuperTokensWithAntiCSRFRandom(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	customAntiCsrfVal := "RANDOM"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:       "127.0.0.1:3000",
			AppName:         "SuperTokens",
			WebsiteDomain:   "google.com",
			APIBasePath:     &apiBasePath0,
			WebsiteBasePath: &websiteBasePath0,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(
				&sessmodels.TypeInput{
					AntiCsrf: &customAntiCsrfVal,
				},
			),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage == "antiCsrf config must be one of 'NONE' or 'VIA_CUSTOM_HEADER' or 'VIA_TOKEN'" {
			assert.Equal(t, errorMessage, "antiCsrf config must be one of 'NONE' or 'VIA_CUSTOM_HEADER' or 'VIA_TOKEN'")
		} else {
			t.Error(errorMessage)
		}
	}
	AfterEach()
}

func TestSuperTokensInitWithDifferentWebAndApiDomain(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	customCookieSecureVal := false
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:       "http://api.test.com:3000",
				AppName:         "SuperTokens",
				WebsiteDomain:   "google.com",
				APIBasePath:     &apiBasePath0,
				WebsiteBasePath: &websiteBasePath0,
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:       "https://api.test.com:3000",
				AppName:         "SuperTokens",
				WebsiteDomain:   "google.com",
				APIBasePath:     &apiBasePath0,
				WebsiteBasePath: &websiteBasePath0,
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					CookieSecure: &customCookieSecureVal,
				}),
			},
		},
	}
	for _, configValue := range configValues {
		BeforeEach()
		StartUpST("localhost", "8080")
		err := supertokens.Init(configValue)
		if err != nil {
			errorMessage := err.Error()
			if errorMessage == "Since your API and website domain are different, for sessions to work, please use https on your apiDomain and dont set cookieSecure to false." {
				assert.Equal(t, errorMessage, "Since your API and website domain are different, for sessions to work, please use https on your apiDomain and dont set cookieSecure to false.")
			} else {
				t.Error(errorMessage)
			}
		}
		AfterEach()
	}
}

func TestSuperTokensForTheDefaultCookieValues(t *testing.T) {

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "https://localhost",
			AppName:       "SuperTokens",
			WebsiteDomain: "http://localhost:3000",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletoneSessionRecipeInstance, err := session.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, singletoneSessionRecipeInstance.Config.CookieSecure, true)
	assert.Equal(t, singletoneSessionRecipeInstance.Config.CookieSameSite, "none")
	AfterEach()
}

func TestSuperTokensInitWithWrongConfigSchema(t *testing.T) {
	customAPIBasePath := "/custom/a"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customAPIBasePath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage == "please provide 'ConnectionURI' value. If you do not want to provide a connection URI, then set config.Supertokens to nil" {
			assert.Equal(t, errorMessage, "please provide 'ConnectionURI' value. If you do not want to provide a connection URI, then set config.Supertokens to nil")
		} else {
			t.Error(err.Error())
		}
	}
	AfterEach()
}

func TestSuperTokensInitWithoutAPIDomain(t *testing.T) {
	customAPIBasePath := "/custom/a"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customAPIBasePath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage == "Please provide your apiDomain inside the appInfo object when calling supertokens.init" {
			assert.Equal(t, errorMessage, "Please provide your apiDomain inside the appInfo object when calling supertokens.init")
		} else {
			t.Error(errorMessage)
		}
	}
	AfterEach()
}

func TestSuperTokensInitWithoutAppName(t *testing.T) {
	customAPIBasePath := "/custom/a"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customAPIBasePath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage == "Please provide your appName inside the appInfo object when calling supertokens.init" {
			assert.Equal(t, errorMessage, "Please provide your appName inside the appInfo object when calling supertokens.init")
		} else {
			t.Error(errorMessage)
		}
	}
	AfterEach()
}

func TestSuperTokensInitWithoutRecipeList(t *testing.T) {
	customAPIBasePath := "/custom/a"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "Supertokens",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customAPIBasePath,
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage == "please provide at least one recipe to the supertokens.init function call" {
			assert.Equal(t, errorMessage, "please provide at least one recipe to the supertokens.init function call")
		} else {
			t.Error(errorMessage)
		}
	}
	AfterEach()
}

func TestSuperTokensDefaultCookieConfig(t *testing.T) {
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
			session.Init(nil),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletoneSessionRecipeInstance, err := session.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Nil(t, singletoneSessionRecipeInstance.Config.CookieDomain)
	assert.Equal(t, singletoneSessionRecipeInstance.Config.CookieSameSite, "lax")
	assert.Equal(t, singletoneSessionRecipeInstance.Config.CookieSecure, true)
	assert.Equal(t, singletoneSessionRecipeInstance.Config.RefreshTokenPath.GetAsStringDangerous(), "/auth/session/refresh")
	assert.Equal(t, singletoneSessionRecipeInstance.Config.SessionExpiredStatusCode, 401)
	AfterEach()
}

func TestJwtFeatureIsDisabledByDefault(t *testing.T) {
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
			session.Init(nil),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletoneSessionRecipeInstance, err := session.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Nil(t, singletoneSessionRecipeInstance.Config.Jwt.Issuer)
	assert.Equal(t, singletoneSessionRecipeInstance.Config.Jwt.Enable, false)
	AfterEach()
}

func TestJWTFeatureDisabledOrEnabledWhenExplicitlyStatedSo(t *testing.T) {
	configValues := []supertokens.TypeInput{
		{
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
					Jwt: &sessmodels.JWTInputConfig{
						Enable: true,
					},
				}),
			},
		},
		{
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
					Jwt: &sessmodels.JWTInputConfig{
						Enable: false,
					},
				}),
			},
		},
	}
	for index, configValue := range configValues {
		BeforeEach()
		StartUpST("localhost", "8080")
		err := supertokens.Init(configValue)
		if err != nil {
			t.Error(err.Error())
		}
		singletoneSessionRecipeInstance, err := session.GetRecipeInstanceOrThrowError()
		if err != nil {
			t.Error(err.Error())
		}
		assert.Nil(t, singletoneSessionRecipeInstance.Config.Jwt.Issuer)
		if index == 0 {
			assert.Equal(t, singletoneSessionRecipeInstance.Config.Jwt.Enable, true)
		} else {
			assert.Equal(t, singletoneSessionRecipeInstance.Config.Jwt.Enable, false)
		}
		AfterEach()
	}
}

func TestJWTPropertyNameIsAccesedWhenSet(t *testing.T) {
	customJWTKey := "customJWTKey"
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
				Jwt: &sessmodels.JWTInputConfig{
					Enable:                           true,
					PropertyNameInAccessTokenPayload: &customJWTKey,
				},
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletoneSessionRecipeInstance, err := session.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Nil(t, singletoneSessionRecipeInstance.Config.Jwt.Issuer)
	assert.Equal(t, singletoneSessionRecipeInstance.Config.Jwt.PropertyNameInAccessTokenPayload, "customJWTKey")
	AfterEach()
}

func TestJWTPropertyNameIsSetCorrectlyByDefault(t *testing.T) {
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
				Jwt: &sessmodels.JWTInputConfig{
					Enable: true,
				},
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletoneSessionRecipeInstance, err := session.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Nil(t, singletoneSessionRecipeInstance.Config.Jwt.Issuer)
	assert.Equal(t, singletoneSessionRecipeInstance.Config.Jwt.PropertyNameInAccessTokenPayload, "jwt")
	AfterEach()
}

func TestJWTPropertyThrowsErrorWhenGetsReservedName(t *testing.T) {
	customJWTKey := "_jwtPName"
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
				Jwt: &sessmodels.JWTInputConfig{
					Enable:                           true,
					PropertyNameInAccessTokenPayload: &customJWTKey,
				},
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage == "_jwtPName is a reserved property name, please use a different key name for the jwt" {
			assert.Equal(t, errorMessage, "_jwtPName is a reserved property name, please use a different key name for the jwt")
		} else {
			t.Error(err.Error())
		}
	}
	AfterEach()
}

func TestSuperTokensInitWithAPIGateWayPath(t *testing.T) {
	customAPIGatewayPath := "/gateway"
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:      "api.supertokens.io",
			AppName:        "SuperTokens",
			WebsiteDomain:  "supertokens.io",
			APIGatewayPath: &customAPIGatewayPath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "ronit", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := ExtractInfoFromResponse(res)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)

	req2.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])

	req2.Header.Add("anti-csrf", cookieData["antiCsrf"])

	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, res2.StatusCode)
	sp, err := supertokens.GetInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sp.AppInfo.APIBasePath.GetAsStringDangerous(), "/gateway/auth")
	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestSuperTokensInitWithAPIGateWayPathAndAPIBasePath(t *testing.T) {
	customAPIGatewayPath := "/gateway"
	customAntiCsrfVal := "VIA_TOKEN"
	customAPIBasePath := "hello"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:      "api.supertokens.io",
			AppName:        "SuperTokens",
			WebsiteDomain:  "supertokens.io",
			APIBasePath:    &customAPIBasePath,
			APIGatewayPath: &customAPIGatewayPath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "uniqueId", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := ExtractInfoFromResponse(res)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/hello/session/refresh", nil)
	assert.NoError(t, err)

	req2.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])

	req2.Header.Add("anti-csrf", cookieData["antiCsrf"])

	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, res2.StatusCode)
	sp, err := supertokens.GetInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sp.AppInfo.APIBasePath.GetAsStringDangerous(), "/gateway/hello")
	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
}

func TestSuperTokensInitWithDefaultAPIGateWayPathandCustomAPIBasePath(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	customAPIBasePath := "hello"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customAPIBasePath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}
	BeforeEach()
	StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		session.CreateNewSession(rw, "uniqueId", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := ExtractInfoFromResponse(res)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/hello/session/refresh", nil)
	assert.NoError(t, err)

	req2.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])

	req2.Header.Add("anti-csrf", cookieData["antiCsrf"])

	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, res2.StatusCode)
	sp, err := supertokens.GetInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sp.AppInfo.APIBasePath.GetAsStringDangerous(), "/hello")
	defer AfterEach()
	defer func() {
		testServer.Close()
	}()
}
