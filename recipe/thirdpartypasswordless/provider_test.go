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

package thirdpartypasswordless

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

const privateKey = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----"

func TestForThirdPartyPasswordlessTheMinimumConfigForThirdPartyProviderGoogle(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     "test",
						ClientSecret: "test-secret",
					}),
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "google", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, "https://oauth2.googleapis.com/token", providerInfoGetResult.AccessTokenAPI.URL)
	assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth", providerInfoGetResult.AuthorisationRedirect.URL)

	assert.Equal(t, "test", providerInfoGetResult.AccessTokenAPI.Params["client_id"])
	assert.Equal(t, "test-secret", providerInfoGetResult.AccessTokenAPI.Params["client_secret"])
	assert.Equal(t, "authorization_code", providerInfoGetResult.AccessTokenAPI.Params["grant_type"])

	assert.Equal(t, "test", providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "offline", providerInfoGetResult.AuthorisationRedirect.Params["access_type"])
	assert.Equal(t, "true", providerInfoGetResult.AuthorisationRedirect.Params["include_granted_scopes"])
	assert.Equal(t, "code", providerInfoGetResult.AuthorisationRedirect.Params["response_type"])
	assert.Equal(t, "https://www.googleapis.com/auth/userinfo.email", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
}

func TestWithThirdPartyPasswordlessPassingAdditionalParamsCheckTheyArePresentInAuthorizationUrlForThirdPartyProviderGoogle(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "google", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, "test", providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "offline", providerInfoGetResult.AuthorisationRedirect.Params["access_type"])
	assert.Equal(t, "true", providerInfoGetResult.AuthorisationRedirect.Params["include_granted_scopes"])
	assert.Equal(t, "code", providerInfoGetResult.AuthorisationRedirect.Params["response_type"])
	assert.Equal(t, "https://www.googleapis.com/auth/userinfo.email", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
	assert.Equal(t, "value1", providerInfoGetResult.AuthorisationRedirect.Params["key1"])
	assert.Equal(t, "value2", providerInfoGetResult.AuthorisationRedirect.Params["key2"])
}

func TestForThirdpartyPasswordlessPassingScopesInConfigForThirdpartyProviderGoogle(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     "test",
						ClientSecret: "test-secret",
						Scope: []string{
							"test-scope-1", "test-scope-2",
						},
					}),
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "google", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, "test", providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "offline", providerInfoGetResult.AuthorisationRedirect.Params["access_type"])
	assert.Equal(t, "true", providerInfoGetResult.AuthorisationRedirect.Params["include_granted_scopes"])
	assert.Equal(t, "code", providerInfoGetResult.AuthorisationRedirect.Params["response_type"])
	assert.Equal(t, "test-scope-1 test-scope-2", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
}

func TestForThirdPartyPasswordlessMinimumConfigForThirdPartyProviderFacebook(t *testing.T) {
	clientId := "test"
	clientSecret := "test-secret"
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Facebook(tpmodels.FacebookConfig{
						ClientID:     clientId,
						ClientSecret: clientSecret,
					}),
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "facebook", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, "https://graph.facebook.com/v9.0/oauth/access_token", providerInfoGetResult.AccessTokenAPI.URL)
	assert.Equal(t, "https://www.facebook.com/v9.0/dialog/oauth", providerInfoGetResult.AuthorisationRedirect.URL)

	assert.Equal(t, clientId, providerInfoGetResult.AccessTokenAPI.Params["client_id"])
	assert.Equal(t, clientSecret, providerInfoGetResult.AccessTokenAPI.Params["client_secret"])

	assert.Equal(t, clientId, providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "code", providerInfoGetResult.AuthorisationRedirect.Params["response_type"])
	assert.Equal(t, "email", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
}

func TestWithThirdPartyPasswordlessPassingScopesInConfigForThirdPartyProviderFacebook(t *testing.T) {
	clientId := "test"
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Facebook(tpmodels.FacebookConfig{
						ClientID:     clientId,
						ClientSecret: "test-secret",
						Scope: []string{
							"test-scope-1", "test-scope-2",
						},
					}),
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "facebook", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, clientId, providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "code", providerInfoGetResult.AuthorisationRedirect.Params["response_type"])
	assert.Equal(t, "test-scope-1 test-scope-2", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
}

func TestWithThirdPartyPasswordlessMinimumConfigForThirdPartyProviderGithub(t *testing.T) {
	clientId := "test"
	clientSecret := "test-secret"
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Github(tpmodels.GithubConfig{
						ClientID:     clientId,
						ClientSecret: clientSecret,
					}),
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "github", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, "https://github.com/login/oauth/access_token", providerInfoGetResult.AccessTokenAPI.URL)
	assert.Equal(t, "https://github.com/login/oauth/authorize", providerInfoGetResult.AuthorisationRedirect.URL)

	assert.Equal(t, clientId, providerInfoGetResult.AccessTokenAPI.Params["client_id"])
	assert.Equal(t, clientSecret, providerInfoGetResult.AccessTokenAPI.Params["client_secret"])

	assert.Equal(t, clientId, providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "read:user user:email", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
}

func TestWithThirdPartyPasswordlessParamCheckTheyArePresentInAuthorizationURLForThirdPartyProviderGithub(t *testing.T) {
	clientId := "test"
	clientSecret := "test-secret"
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Github(tpmodels.GithubConfig{
						ClientID:     clientId,
						ClientSecret: clientSecret,
						AuthorisationRedirect: &struct{ Params map[string]interface{} }{
							Params: map[string]interface{}{
								"key1": "value1",
								"key2": "value2",
							},
						},
					}),
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "github", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, clientId, providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "read:user user:email", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
	assert.Equal(t, "value1", providerInfoGetResult.AuthorisationRedirect.Params["key1"])
	assert.Equal(t, "value2", providerInfoGetResult.AuthorisationRedirect.Params["key2"])
}

func TestWithThirdPartyPasswordlessPassingScopesInConfigForThirdPartyProviderGithub(t *testing.T) {
	clientId := "test"
	clientSecret := "test-secret"
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Github(tpmodels.GithubConfig{
						ClientID:     clientId,
						ClientSecret: clientSecret,
						Scope: []string{
							"test-scope-1", "test-scope-2",
						},
					}),
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "github", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, clientId, providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "test-scope-1 test-scope-2", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
}

func TestWithThirdPartyPasswordlessMinimumConfigForThirdPartyProviderApple(t *testing.T) {
	clientId := "test"
	clientSecret := tpmodels.AppleClientSecret{
		KeyId:      "test-key",
		PrivateKey: privateKey,
		TeamId:     "test-team-id",
	}

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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Apple(tpmodels.AppleConfig{
						ClientID:     clientId,
						ClientSecret: clientSecret,
					}),
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "apple", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, "https://appleid.apple.com/auth/token", providerInfoGetResult.AccessTokenAPI.URL)
	assert.Equal(t, "https://appleid.apple.com/auth/authorize", providerInfoGetResult.AuthorisationRedirect.URL)

	assert.Equal(t, clientId, providerInfoGetResult.AccessTokenAPI.Params["client_id"])
	assert.NotNil(t, providerInfoGetResult.AccessTokenAPI.Params["client_secret"])
	assert.Equal(t, "authorization_code", providerInfoGetResult.AccessTokenAPI.Params["grant_type"])

	assert.Equal(t, clientId, providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "email", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
	assert.Equal(t, "form_post", providerInfoGetResult.AuthorisationRedirect.Params["response_mode"])
	assert.Equal(t, "code", providerInfoGetResult.AuthorisationRedirect.Params["response_type"])
}

func TestWithThirdPartyPasswordlessPassingAdditionalParamsCheckTheyArePresentInAuthorizationURLForThirdPartyProviderApple(t *testing.T) {
	clientId := "test"
	clientSecret := tpmodels.AppleClientSecret{
		KeyId:      "test-key",
		PrivateKey: privateKey,
		TeamId:     "test-team-id",
	}

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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Apple(tpmodels.AppleConfig{
						ClientID:     clientId,
						ClientSecret: clientSecret,
						AuthorisationRedirect: &struct{ Params map[string]interface{} }{
							Params: map[string]interface{}{
								"key1": "value1",
								"key2": "value2",
							},
						},
					}),
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "apple", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, clientId, providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "email", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
	assert.Equal(t, "form_post", providerInfoGetResult.AuthorisationRedirect.Params["response_mode"])
	assert.Equal(t, "code", providerInfoGetResult.AuthorisationRedirect.Params["response_type"])
	assert.Equal(t, "value1", providerInfoGetResult.AuthorisationRedirect.Params["key1"])
	assert.Equal(t, "value2", providerInfoGetResult.AuthorisationRedirect.Params["key2"])
}

func TestWithThirdPartyProviderPasswordlessPassingScopesInConfigForThirdPartyProviderApple(t *testing.T) {
	clientId := "test"
	clientSecret := tpmodels.AppleClientSecret{
		KeyId:      "test-key",
		PrivateKey: privateKey,
		TeamId:     "test-team-id",
	}

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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Apple(tpmodels.AppleConfig{
						ClientID:     clientId,
						ClientSecret: clientSecret,
						Scope: []string{
							"test-scope-1", "test-scope-2",
						},
					}),
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	thirdpartypasswordlessrecipeinstance, err := GetRecipeInstanceOrThrowError()
	assert.NoError(t, err)

	providerInfo := thirdpartypasswordlessrecipeinstance.Config.Providers[0]
	assert.Equal(t, "apple", providerInfo.ID)

	providerInfoGetResult := providerInfo.Get(nil, nil, nil)

	assert.Equal(t, clientId, providerInfoGetResult.AuthorisationRedirect.Params["client_id"])
	assert.Equal(t, "test-scope-1 test-scope-2", providerInfoGetResult.AuthorisationRedirect.Params["scope"])
	assert.Equal(t, "form_post", providerInfoGetResult.AuthorisationRedirect.Params["response_mode"])
	assert.Equal(t, "code", providerInfoGetResult.AuthorisationRedirect.Params["response_type"])
}

func TestWithThirdPartyPasswordlessDuplicateProviderWithoutAnyDefault(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     "test",
						ClientSecret: "test-secret",
						Scope: []string{
							"test-scope-1", "test-scope-2",
						},
					}),
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     "test",
						ClientSecret: "test-secret",
						Scope: []string{
							"test-scope-1", "test-scope-2",
						},
					}),
				},
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, "The providers array has multiple entries for the same third party provider. Please mark one of them as the default one by using 'IsDefault: true'", err.Error())
	}
}

func TestWithThirdPartyPasswordlessDuplicateProviderWithBothDefault(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     "test",
						ClientSecret: "test-secret",
						Scope: []string{
							"test-scope-1", "test-scope-2",
						},
						IsDefault: true,
					}),
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     "test",
						ClientSecret: "test-secret",
						Scope: []string{
							"test-scope-1", "test-scope-2",
						},
						IsDefault: true,
					}),
				},
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		assert.Equal(t, "You have provided multiple third party providers that have the id: google and are marked as 'IsDefault: true'. Please only mark one of them as isDefault", err.Error())
	}
}

func TestWithThirdPartyPasswordlessDuplicateProviderWithOneMarkedAsDefault(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     "test",
						ClientSecret: "test-secret",
						Scope: []string{
							"test-scope-1", "test-scope-2",
						},
						IsDefault: true,
					}),
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     "test",
						ClientSecret: "test-secret",
						Scope: []string{
							"test-scope-1", "test-scope-2",
						},
					}),
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
}
