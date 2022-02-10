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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestMinimumConfigForGoogleAsThirdPartyProvider(t *testing.T) {
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
							thirdparty.Google(tpmodels.GoogleConfig{
								ClientID:     "test",
								ClientSecret: "test-secret",
							}),
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

	singletonInstance, err := thirdparty.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	providerInfo := singletonInstance.Providers[0]

	assert.Equal(t, "google", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil)

	assert.Equal(t, "https://accounts.google.com/o/oauth2/token", providerInfoGetResult.AccessTokenAPI.URL)
	assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth", providerInfoGetResult.AuthorisationRedirect.URL)

	assert.Equal(t, map[string]string{
		"client_id":     "test",
		"client_secret": "test-secret",
		"grant_type":    "authorization_code",
	}, providerInfoGetResult.AccessTokenAPI.Params)

	assert.Equal(t, map[string]interface{}{
		"client_id":              "test",
		"access_type":            "offline",
		"include_granted_scopes": "true",
		"response_type":          "code",
		"scope":                  "https://www.googleapis.com/auth/userinfo.email",
	}, providerInfoGetResult.AuthorisationRedirect.Params)

	defer unittesting.AfterEach()
}

func TestPassingAdditionalParamsInAuthUrlForGoogleAndCheckItsPresense(t *testing.T) {
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
							thirdparty.Google(tpmodels.GoogleConfig{
								ClientID:     "test",
								ClientSecret: "test-secret",
								AuthorisationRedirect: &struct{ Params map[string]interface{} }{
									Params: map[string]interface{}{
										"key1": "value1",
										"key2": "value2",
									},
								},
							}),
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

	singletonInstance, err := thirdparty.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	providerInfo := singletonInstance.Providers[0]

	assert.Equal(t, "google", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil)

	assert.Equal(t, "https://accounts.google.com/o/oauth2/token", providerInfoGetResult.AccessTokenAPI.URL)
	assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth", providerInfoGetResult.AuthorisationRedirect.URL)

	assert.Equal(t, map[string]string{
		"client_id":     "test",
		"client_secret": "test-secret",
		"grant_type":    "authorization_code",
	}, providerInfoGetResult.AccessTokenAPI.Params)

	assert.Equal(t, map[string]interface{}{
		"client_id":              "test",
		"access_type":            "offline",
		"include_granted_scopes": "true",
		"response_type":          "code",
		"scope":                  "https://www.googleapis.com/auth/userinfo.email",
		"key1":                   "value1",
		"key2":                   "value2",
	}, providerInfoGetResult.AuthorisationRedirect.Params)

	defer unittesting.AfterEach()
}

func TestPassingScopesInConfigForGoogle(t *testing.T) {
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
							thirdparty.Google(tpmodels.GoogleConfig{
								ClientID:     "test",
								ClientSecret: "test-secret",
								Scope: []string{
									"test-scope-1", "test-scope-2",
								},
							}),
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

	singletonInstance, err := thirdparty.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	providerInfo := singletonInstance.Providers[0]

	assert.Equal(t, "google", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil)

	assert.Equal(t, "https://accounts.google.com/o/oauth2/token", providerInfoGetResult.AccessTokenAPI.URL)
	assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth", providerInfoGetResult.AuthorisationRedirect.URL)

	assert.Equal(t, map[string]string{
		"client_id":     "test",
		"client_secret": "test-secret",
		"grant_type":    "authorization_code",
	}, providerInfoGetResult.AccessTokenAPI.Params)

	assert.Equal(t, map[string]interface{}{
		"client_id":              "test",
		"access_type":            "offline",
		"include_granted_scopes": "true",
		"response_type":          "code",
		"scope":                  "test-scope-1 test-scope-2",
	}, providerInfoGetResult.AuthorisationRedirect.Params)

	defer unittesting.AfterEach()
}

func TestMinimumConfigForFacebookAsThirdPartyProvider(t *testing.T) {
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
							thirdparty.Facebook(tpmodels.FacebookConfig{
								ClientID:     "test",
								ClientSecret: "test-secret",
							}),
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

	singletonInstance, err := thirdparty.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	providerInfo := singletonInstance.Providers[0]

	assert.Equal(t, "facebook", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil)

	assert.Equal(t, "https://graph.facebook.com/v9.0/oauth/access_token", providerInfoGetResult.AccessTokenAPI.URL)
	assert.Equal(t, "https://www.facebook.com/v9.0/dialog/oauth", providerInfoGetResult.AuthorisationRedirect.URL)

	assert.Equal(t, map[string]string{
		"client_id":     "test",
		"client_secret": "test-secret",
	}, providerInfoGetResult.AccessTokenAPI.Params)

	assert.Equal(t, map[string]interface{}{
		"client_id":     "test",
		"response_type": "code",
		"scope":         "email",
	}, providerInfoGetResult.AuthorisationRedirect.Params)

	defer unittesting.AfterEach()
}

func TestPassingScopesInConfigForFacebook(t *testing.T) {
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
							thirdparty.Facebook(tpmodels.FacebookConfig{
								ClientID:     "test",
								ClientSecret: "test-secret",
								Scope: []string{
									"test-scope-1", "test-scope-2",
								},
							}),
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

	singletonInstance, err := thirdparty.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	providerInfo := singletonInstance.Providers[0]

	assert.Equal(t, "facebook", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil)

	assert.Equal(t, map[string]interface{}{
		"client_id":     "test",
		"response_type": "code",
		"scope":         "test-scope-1 test-scope-2",
	}, providerInfoGetResult.AuthorisationRedirect.Params)

	defer unittesting.AfterEach()
}

func TestMinimumConfigForGithubAsThirdPartyProvider(t *testing.T) {
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
							thirdparty.Github(tpmodels.GithubConfig{
								ClientID:     "test",
								ClientSecret: "test-secret",
							}),
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

	singletonInstance, err := thirdparty.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	providerInfo := singletonInstance.Providers[0]

	assert.Equal(t, "github", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil)

	assert.Equal(t, "https://github.com/login/oauth/access_token", providerInfoGetResult.AccessTokenAPI.URL)
	assert.Equal(t, "https://github.com/login/oauth/authorize", providerInfoGetResult.AuthorisationRedirect.URL)

	assert.Equal(t, map[string]string{
		"client_id":     "test",
		"client_secret": "test-secret",
	}, providerInfoGetResult.AccessTokenAPI.Params)

	assert.Equal(t, map[string]interface{}{
		"client_id": "test",
		"scope":     "read:user user:email",
	}, providerInfoGetResult.AuthorisationRedirect.Params)

	defer unittesting.AfterEach()
}

func TestPassingAdditionalParamsInAuthUrlForGithubAndCheckItsPresense(t *testing.T) {
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
							thirdparty.Github(tpmodels.GithubConfig{
								ClientID:     "test",
								ClientSecret: "test-secret",
								AuthorisationRedirect: &struct{ Params map[string]interface{} }{
									Params: map[string]interface{}{
										"key1": "value1",
										"key2": "value2",
									},
								},
							}),
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

	singletonInstance, err := thirdparty.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	providerInfo := singletonInstance.Providers[0]

	assert.Equal(t, "github", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil)

	assert.Equal(t, map[string]interface{}{
		"client_id": "test",
		"scope":     "read:user user:email",
		"key1":      "value1",
		"key2":      "value2",
	}, providerInfoGetResult.AuthorisationRedirect.Params)

	defer unittesting.AfterEach()
}

func TestPassingScopesInConfigForGithub(t *testing.T) {
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
							thirdparty.Github(tpmodels.GithubConfig{
								ClientID:     "test",
								ClientSecret: "test-secret",
								Scope: []string{
									"test-scope-1", "test-scope-2",
								},
							}),
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

	singletonInstance, err := thirdparty.GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	providerInfo := singletonInstance.Providers[0]

	assert.Equal(t, "github", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil)

	assert.Equal(t, map[string]interface{}{
		"client_id": "test",
		"scope":     "test-scope-1 test-scope-2",
	}, providerInfoGetResult.AuthorisationRedirect.Params)

	defer unittesting.AfterEach()
}
