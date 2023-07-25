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

package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestSuperTokensInitWithJustTheCompulsoryFields(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}
	BeforeEach()

	unittesting.StartUpST("localhost", "8080")

	defer AfterEach()

	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	supertokensInstance, err := supertokens.GetInstanceOrThrowError()

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "/auth", supertokensInstance.AppInfo.APIBasePath.GetAsStringDangerous())
	assert.Equal(t, "/auth", supertokensInstance.AppInfo.WebsiteBasePath.GetAsStringDangerous())
}

func TestSuperTokensInitWithOptionalFields(t *testing.T) {
	apiBasePath := "test/"
	websiteBasePath := "test1/"
	configValue := supertokens.TypeInput{
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
			Init(&sessmodels.TypeInput{}),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")

	defer AfterEach()

	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	supertokensInstance, err := supertokens.GetInstanceOrThrowError()

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "/test", supertokensInstance.AppInfo.APIBasePath.GetAsStringDangerous())
	assert.Equal(t, "/test1", supertokensInstance.AppInfo.WebsiteBasePath.GetAsStringDangerous())
}

func TestSuperTokensInitWithoutApiDomain(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, err.Error(), "Please provide your apiDomain inside the appInfo object when calling supertokens.init")
	} else {
		t.Fail()
	}
}

func TestSuperTokensInitWithoutWebsiteDomain(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:   "SuperTokens",
			APIDomain: "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, err.Error(), "Please provide your websiteDomain inside the appInfo object when calling supertokens.init")
	} else {
		t.Fail()
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
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, err.Error(), "please provide at least one recipe to the supertokens.init function call")
	} else {
		t.Fail()
	}
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
			Init(
				&sessmodels.TypeInput{
					CookieDomain:             &cookieDomain,
					SessionExpiredStatusCode: &sessionExpiredStatusCode,
					CookieSecure:             &cookieSecure,
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

	sessionSingletonInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, *sessionSingletonInstance.Config.CookieDomain, "testdomain")
	assert.Equal(t, sessionSingletonInstance.Config.SessionExpiredStatusCode, 111)
	assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, true)
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
			Init(&sessmodels.TypeInput{
				CookieSameSite: &cookieSameSite0,
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
			Init(&sessmodels.TypeInput{
				CookieSameSite: &cookieSameSite0,
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
			Init(&sessmodels.TypeInput{
				CookieSameSite: &cookieSameSite0,
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
			Init(&sessmodels.TypeInput{
				CookieSameSite: &cookieSameSite0,
			}),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, err.Error(), `cookie same site must be one of "strict", "lax", or "none"`)
	} else {
		t.Fail()
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, *supertokens.QuerierAPIKey, configValue.Supertokens.APIKey)
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
			Init(
				&sessmodels.TypeInput{
					SessionExpiredStatusCode: &customSessionExpiredCode,
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
	sessionSingletonInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sessionSingletonInstance.Config.SessionExpiredStatusCode, 402)
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	hosts := supertokens.QuerierHosts

	assert.Equal(t, hosts[0].Domain.GetAsStringDangerous(), "http://localhost:8080")
	assert.Equal(t, hosts[1].Domain.GetAsStringDangerous(), "https://try.supertokens.io")
	assert.Equal(t, hosts[2].Domain.GetAsStringDangerous(), "https://try.supertokens.io:8080")
	assert.Equal(t, hosts[3].Domain.GetAsStringDangerous(), "http://localhost:90")
}

func TestSuperTokensInitWithNoneLaxFalseSessionConfigResults(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	configValue := supertokens.TypeInput{
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	sessionSingletonInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sessionSingletonInstance.Config.AntiCsrf, "NONE")
	assert.Equal(t, sessionSingletonInstance.Config.CookieSameSite, "lax")
	assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, false)
}

func TestSuperTokensInitWithCustomHeaderLaxTrueSessionConfigResults(t *testing.T) {
	apiBasePath0 := "/"
	customAntiCsrf := "VIA_CUSTOM_HEADER"
	configValue := supertokens.TypeInput{

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
			Init(
				&sessmodels.TypeInput{
					AntiCsrf: &customAntiCsrf,
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
	sessionSingletonInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sessionSingletonInstance.Config.AntiCsrf, "VIA_CUSTOM_HEADER")
	assert.Equal(t, sessionSingletonInstance.Config.CookieSameSite, "lax")
	assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, true)
}

func TestSuperTokensInitWithCustomHeaderLaxFalseSessionConfigResults(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	customAntiCsrf := "VIA_CUSTOM_HEADER"
	configValue := supertokens.TypeInput{
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
			Init(
				&sessmodels.TypeInput{
					AntiCsrf: &customAntiCsrf,
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
	sessionSingletonInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sessionSingletonInstance.Config.AntiCsrf, "VIA_CUSTOM_HEADER")
	assert.Equal(t, sessionSingletonInstance.Config.CookieSameSite, "lax")
	assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, false)
}

func TestSuperTokensInitWithCustomHeaderNoneTrueSessionConfigResultsWithNormalWebsiteDomain(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	configValue := supertokens.TypeInput{
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	sessionSingletonInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sessionSingletonInstance.Config.AntiCsrf, "VIA_CUSTOM_HEADER")
	assert.Equal(t, sessionSingletonInstance.Config.CookieSameSite, "none")
	assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, true)
}

func TestSuperTokensInitWithCustomHeaderNoneTrueSessionConfigResultsWithLocalWebsiteDomain(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	configValue := supertokens.TypeInput{
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	sessionSingletonInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sessionSingletonInstance.Config.AntiCsrf, "VIA_CUSTOM_HEADER")
	assert.Equal(t, sessionSingletonInstance.Config.CookieSameSite, "none")
	assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, true)
}

func TestSuperTokensWithAntiCSRFNone(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	customAntiCsrfVal := "NONE"
	True := true
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
			Init(
				&sessmodels.TypeInput{
					AntiCsrf:     &customAntiCsrfVal,
					CookieSecure: &True,
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
	singletoneSessionRecipeInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, singletoneSessionRecipeInstance.Config.AntiCsrf, "NONE")
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
			Init(
				&sessmodels.TypeInput{
					AntiCsrf: &customAntiCsrfVal,
				},
			),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, err.Error(), "antiCsrf config must be one of 'NONE' or 'VIA_CUSTOM_HEADER' or 'VIA_TOKEN'")
	} else {
		t.Fail()
	}
}

func TestSuperTokensInitWithDifferentWebAndApiDomainWithDefaultCookieSecure(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	configValue := supertokens.TypeInput{
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	// The following error is not expected to happen on supertokens init anymore
	// Since your API and website domain are different, for sessions to work, please use https on your apiDomain and dont set cookieSecure to false.
	assert.NoError(t, err)
}
func TestSuperTokensInitWithDifferentWebAndApiDomainWithCookieSecureValueSetToFalse(t *testing.T) {
	apiBasePath0 := "test/"
	websiteBasePath0 := "test1/"
	customCookieSecureVal := false
	configValue := supertokens.TypeInput{
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
			Init(&sessmodels.TypeInput{
				CookieSecure: &customCookieSecureVal,
			}),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	// The following error is not expected to happen on supertokens init anymore
	// Since your API and website domain are different, for sessions to work, please use https on your apiDomain and dont set cookieSecure to false.
	assert.NoError(t, err)
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletoneSessionRecipeInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, singletoneSessionRecipeInstance.Config.CookieSecure, true)
	assert.Equal(t, singletoneSessionRecipeInstance.Config.CookieSameSite, "none")
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, err.Error(), "please provide 'ConnectionURI' value. If you do not want to provide a connection URI, then set config.Supertokens to nil")
	} else {
		t.Fail()
	}
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, err.Error(), "Please provide your apiDomain inside the appInfo object when calling supertokens.init")
	} else {
		t.Fail()
	}
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, err.Error(), "Please provide your appName inside the appInfo object when calling supertokens.init")
	} else {
		t.Fail()
	}
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
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, err.Error(), "please provide at least one recipe to the supertokens.init function call")
	} else {
		t.Fail()
	}
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletoneSessionRecipeInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Nil(t, singletoneSessionRecipeInstance.Config.CookieDomain)
	assert.Equal(t, singletoneSessionRecipeInstance.Config.CookieSameSite, "lax")
	assert.Equal(t, singletoneSessionRecipeInstance.Config.CookieSecure, true)
	assert.Equal(t, singletoneSessionRecipeInstance.Config.RefreshTokenPath.GetAsStringDangerous(), "/auth/session/refresh")
	assert.Equal(t, singletoneSessionRecipeInstance.Config.SessionExpiredStatusCode, 401)
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
			Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "ronit", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)

	req2.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])

	req2.Header.Add("anti-csrf", cookieData["antiCsrf"])

	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, res2.StatusCode)
	sp, err := supertokens.GetInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sp.AppInfo.APIBasePath.GetAsStringDangerous(), "/gateway/auth")
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
			Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "uniqueId", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/hello/session/refresh", nil)
	assert.NoError(t, err)

	req2.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])

	req2.Header.Add("anti-csrf", cookieData["antiCsrf"])

	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, res2.StatusCode)
	sp, err := supertokens.GetInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sp.AppInfo.APIBasePath.GetAsStringDangerous(), "/gateway/hello")
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
			Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "uniqueId", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	req2, err := http.NewRequest(http.MethodPost, testServer.URL+"/hello/session/refresh", nil)
	assert.NoError(t, err)

	req2.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"])

	req2.Header.Add("anti-csrf", cookieData["antiCsrf"])

	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, res2.StatusCode)
	sp, err := supertokens.GetInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, sp.AppInfo.APIBasePath.GetAsStringDangerous(), "/hello")
}

func TestInvalidSameSiteNoneConfig(t *testing.T) {
	domainCombinations := []struct {
		WebsiteDomain string
		APIDomain     string
	}{
		{"http://localhost:3000", "http://supertokensapi.io"},
		{"http://127.0.0.1:3000", "http://supertokensapi.io"},
		{"http://supertokens.io", "http://localhost:8000"},
		{"http://supertokens.io", "http://127.0.0.1:8000"},
		{"http://supertokens.io", "http://supertokensapi.io"},
	}

	None := "none"

	for _, domainCombination := range domainCombinations {
		BeforeEach()
		configValue := supertokens.TypeInput{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				WebsiteDomain: domainCombination.WebsiteDomain,
				APIDomain:     domainCombination.APIDomain,
			},
			RecipeList: []supertokens.Recipe{
				Init(&sessmodels.TypeInput{
					CookieSameSite: &None,
					GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
						return sessmodels.CookieTransferMethod
					},
				}),
			},
		}

		// The following error is not expected to come from supertokens.Init any more
		// Since your API and website domain are different, for sessions to work, please use https on your apiDomain and dont set cookieSecure to false.
		err := supertokens.Init(configValue)
		assert.NoError(t, err)
		AfterEach()
	}
}

func TestValidSameSiteNoneConfig(t *testing.T) {
	domainCombinations := []struct {
		WebsiteDomain string
		APIDomain     string
	}{
		{"http://localhost:3000", "http://localhost:8000"},
		{"http://127.0.0.1:3000", "http://localhost:8000"},
		{"http://localhost:3000", "http://127.0.0.1:8000"},
		{"http://127.0.0.1:3000", "http://127.0.0.1:8000"},

		{"https://localhost:3000", "https://localhost:8000"},
		{"https://127.0.0.1:3000", "https://localhost:8000"},
		{"https://localhost:3000", "https://127.0.0.1:8000"},
		{"https://127.0.0.1:3000", "https://127.0.0.1:8000"},

		{"https://supertokens.io", "https://api.supertokens.io"},
		{"https://supertokens.io", "https://supertokensapi.io"},

		{"http://localhost:3000", "https://supertokensapi.io"},
		{"http://127.0.0.1:3000", "https://supertokensapi.io"},
	}

	None := "none"

	for _, domainCombination := range domainCombinations {
		BeforeEach()
		configValue := supertokens.TypeInput{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				WebsiteDomain: domainCombination.WebsiteDomain,
				APIDomain:     domainCombination.APIDomain,
			},
			RecipeList: []supertokens.Recipe{
				Init(&sessmodels.TypeInput{
					CookieSameSite: &None,
				}),
			},
		}
		err := supertokens.Init(configValue)
		assert.NoError(t, err)
		AfterEach()
	}
}

func TestThatJWKSAndOpenIdEndpointsAreExposed(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
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

	recipe, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	apisHandled, err := recipe.RecipeModule.GetAPIsHandled()
	if err != nil {
		t.Error(err.Error())
	}

	var jwksAPI supertokens.APIHandled
	var openIdAPI supertokens.APIHandled

	for _, apiHandled := range apisHandled {
		if apiHandled.ID == "/jwt/jwks.json" {
			jwksAPI = apiHandled
		}

		if apiHandled.ID == "/.well-known/openid-configuration" {
			openIdAPI = apiHandled
		}
	}

	assert.NotNil(t, jwksAPI)
	assert.Equal(t, jwksAPI.PathWithoutAPIBasePath.GetAsStringDangerous(), "/jwt/jwks.json")

	assert.NotNil(t, openIdAPI)
	assert.Equal(t, openIdAPI.PathWithoutAPIBasePath.GetAsStringDangerous(), "/.well-known/openid-configuration")
}
