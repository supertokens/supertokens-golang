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
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"gopkg.in/h2non/gock.v1"
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
								},
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "google", nil)
	assert.NoError(t, err)

	providerInfo := providerRes

	assert.Equal(t, "google", providerInfo.ID)

	assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth", providerInfo.Config.AuthorizationEndpoint)
	assert.Equal(t, "https://oauth2.googleapis.com/token", providerInfo.Config.TokenEndpoint)
	assert.Equal(t, "https://openidconnect.googleapis.com/v1/userinfo", providerInfo.Config.UserInfoEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":              {"test"},
		"access_type":            {"offline"},
		"include_granted_scopes": {"true"},
		"response_type":          {"code"},
		"redirect_uri":           {"redirect"},
		"scope":                  {"openid email"},
	}, authParams)

	tokenParams := url.Values{}

	defer gock.OffAll()
	gock.New("https://oauth2.googleapis.com").
		Post("/token").
		Persist().
		Map(func(r *http.Request) *http.Request {
			data, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			tokenParams, err = url.ParseQuery(string(data))
			assert.NoError(t, err)
			return r
		}).
		Reply(200).
		JSON(map[string]string{
			"access_token": "abcd",
		})

	_, err = providerInfo.ExchangeAuthCodeForOAuthTokens(tpmodels.TypeRedirectURIInfo{
		RedirectURIOnProviderDashboard: "redirect",
		RedirectURIQueryParams: map[string]interface{}{
			"code": "abcd",
		},
	}, &map[string]interface{}{})
	assert.NoError(t, err)

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"client_secret": {"test-secret"},
		"grant_type":    {"authorization_code"},
		"code":          {"abcd"},
		"redirect_uri":  {"redirect"},
	}, tokenParams)
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
								},
							},
							AuthorizationEndpointQueryParams: map[string]interface{}{
								"key1": "value1",
								"key2": "value2",
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "google", nil)
	assert.NoError(t, err)

	providerInfo := providerRes
	assert.Equal(t, "google", providerInfo.ID)

	assert.Equal(t, "https://oauth2.googleapis.com/token", providerInfo.Config.TokenEndpoint)
	assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth", providerInfo.Config.AuthorizationEndpoint)
	assert.Equal(t, "https://openidconnect.googleapis.com/v1/userinfo", providerInfo.Config.UserInfoEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":              {"test"},
		"access_type":            {"offline"},
		"include_granted_scopes": {"true"},
		"response_type":          {"code"},
		"redirect_uri":           {"redirect"},
		"scope":                  {"openid email"},
		"key1":                   {"value1"},
		"key2":                   {"value2"},
	}, authParams)
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
									Scope:        []string{"test-scope-1", "test-scope-2"},
								},
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "google", nil)
	assert.NoError(t, err)

	providerInfo := providerRes

	assert.Equal(t, "google", providerInfo.ID)

	assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth", providerInfo.Config.AuthorizationEndpoint)
	assert.Equal(t, "https://oauth2.googleapis.com/token", providerInfo.Config.TokenEndpoint)
	assert.Equal(t, "https://openidconnect.googleapis.com/v1/userinfo", providerInfo.Config.UserInfoEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":              {"test"},
		"access_type":            {"offline"},
		"include_granted_scopes": {"true"},
		"response_type":          {"code"},
		"redirect_uri":           {"redirect"},
		"scope":                  {"test-scope-1 test-scope-2"},
	}, authParams)
}

func TestForThirdPartyPasswordlessMinimumConfigForThirdPartyProviderFacebook(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "facebook",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
								},
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "facebook", nil)
	assert.NoError(t, err)

	providerInfo := providerRes

	assert.Equal(t, "facebook", providerInfo.ID)

	assert.Equal(t, "https://graph.facebook.com/v12.0/oauth/access_token", providerInfo.Config.TokenEndpoint)
	assert.Equal(t, "https://www.facebook.com/v12.0/dialog/oauth", providerInfo.Config.AuthorizationEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"response_type": {"code"},
		"redirect_uri":  {"redirect"},
		"scope":         {"email"},
	}, authParams)

	tokenParams := url.Values{}

	defer gock.OffAll()
	gock.New("https://graph.facebook.com").
		Post("/v12.0/oauth/access_token").
		Persist().
		Map(func(r *http.Request) *http.Request {
			data, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			tokenParams, err = url.ParseQuery(string(data))
			assert.NoError(t, err)
			return r
		}).
		Reply(200).
		JSON(map[string]string{
			"access_token": "abcd",
		})

	_, err = providerInfo.ExchangeAuthCodeForOAuthTokens(tpmodels.TypeRedirectURIInfo{
		RedirectURIOnProviderDashboard: "redirect",
		RedirectURIQueryParams: map[string]interface{}{
			"code": "abcd",
		},
	}, &map[string]interface{}{})
	assert.NoError(t, err)

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"client_secret": {"test-secret"},
		"grant_type":    {"authorization_code"},
		"code":          {"abcd"},
		"redirect_uri":  {"redirect"},
	}, tokenParams)
}

func TestWithThirdPartyPasswordlessPassingScopesInConfigForThirdPartyProviderFacebook(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "facebook",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
									Scope:        []string{"test-scope-1", "test-scope-2"},
								},
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "facebook", nil)
	assert.NoError(t, err)

	providerInfo := providerRes

	assert.Equal(t, "facebook", providerInfo.ID)

	assert.Equal(t, "https://graph.facebook.com/v12.0/oauth/access_token", providerInfo.Config.TokenEndpoint)
	assert.Equal(t, "https://www.facebook.com/v12.0/dialog/oauth", providerInfo.Config.AuthorizationEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"response_type": {"code"},
		"redirect_uri":  {"redirect"},
		"scope":         {"test-scope-1 test-scope-2"},
	}, authParams)
}

func TestWithThirdPartyPasswordlessMinimumConfigForThirdPartyProviderGithub(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "github",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
								},
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "github", nil)
	assert.NoError(t, err)

	providerInfo := providerRes

	assert.Equal(t, "github", providerInfo.ID)

	assert.Equal(t, "https://github.com/login/oauth/authorize", providerInfo.Config.AuthorizationEndpoint)
	assert.Equal(t, "https://github.com/login/oauth/access_token", providerInfo.Config.TokenEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"response_type": {"code"},
		"redirect_uri":  {"redirect"},
		"scope":         {"read:user user:email"},
	}, authParams)

	tokenParams := url.Values{}

	defer gock.OffAll()
	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Persist().
		Map(func(r *http.Request) *http.Request {
			data, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			tokenParams, err = url.ParseQuery(string(data))
			assert.NoError(t, err)
			return r
		}).
		Reply(200).
		JSON(map[string]string{
			"access_token": "abcd",
		})

	_, err = providerInfo.ExchangeAuthCodeForOAuthTokens(tpmodels.TypeRedirectURIInfo{
		RedirectURIOnProviderDashboard: "redirect",
		RedirectURIQueryParams: map[string]interface{}{
			"code": "abcd",
		},
	}, &map[string]interface{}{})
	assert.NoError(t, err)

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"client_secret": {"test-secret"},
		"grant_type":    {"authorization_code"},
		"code":          {"abcd"},
		"redirect_uri":  {"redirect"},
	}, tokenParams)
}

func TestWithThirdPartyPasswordlessParamCheckTheyArePresentInAuthorizationURLForThirdPartyProviderGithub(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "github",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
								},
							},
							AuthorizationEndpointQueryParams: map[string]interface{}{
								"key1": "value1",
								"key2": "value2",
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "github", nil)
	assert.NoError(t, err)

	providerInfo := providerRes

	assert.Equal(t, "github", providerInfo.ID)

	assert.Equal(t, "https://github.com/login/oauth/authorize", providerInfo.Config.AuthorizationEndpoint)
	assert.Equal(t, "https://github.com/login/oauth/access_token", providerInfo.Config.TokenEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"response_type": {"code"},
		"redirect_uri":  {"redirect"},
		"scope":         {"read:user user:email"},
		"key1":          {"value1"},
		"key2":          {"value2"},
	}, authParams)
}

func TestWithThirdPartyPasswordlessPassingScopesInConfigForThirdPartyProviderGithub(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "github",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
									Scope:        []string{"test-scope-1", "test-scope-2"},
								},
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "github", nil)
	assert.NoError(t, err)

	providerInfo := providerRes

	assert.Equal(t, "github", providerInfo.ID)

	assert.Equal(t, "https://github.com/login/oauth/authorize", providerInfo.Config.AuthorizationEndpoint)
	assert.Equal(t, "https://github.com/login/oauth/access_token", providerInfo.Config.TokenEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"response_type": {"code"},
		"redirect_uri":  {"redirect"},
		"scope":         {"test-scope-1 test-scope-2"},
	}, authParams)
}

func TestWithThirdPartyPasswordlessMinimumConfigForThirdPartyProviderApple(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "apple",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID: clientId,
									AdditionalConfig: map[string]interface{}{
										"keyId":      "test-key",
										"privateKey": privateKey,
										"teamId":     "test-team-id",
									},
								},
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "apple", nil)
	assert.NoError(t, err)

	providerInfo := providerRes

	assert.Equal(t, "apple", providerInfo.ID)

	assert.Equal(t, "https://appleid.apple.com/auth/authorize", providerInfo.Config.AuthorizationEndpoint)
	assert.Equal(t, "https://appleid.apple.com/auth/token", providerInfo.Config.TokenEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"response_mode": {"form_post"},
		"response_type": {"code"},
		"redirect_uri":  {"redirect"},
		"scope":         {"openid email"},
	}, authParams)

	tokenParams := url.Values{}

	defer gock.OffAll()
	gock.New("https://appleid.apple.com").
		Post("/auth/token").
		Persist().
		Map(func(r *http.Request) *http.Request {
			data, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			tokenParams, err = url.ParseQuery(string(data))
			assert.NoError(t, err)
			return r
		}).
		Reply(200).
		JSON(map[string]string{
			"id_token": "abcd",
		})

	_, err = providerInfo.ExchangeAuthCodeForOAuthTokens(tpmodels.TypeRedirectURIInfo{
		RedirectURIOnProviderDashboard: "redirect",
		RedirectURIQueryParams: map[string]interface{}{
			"code": "abcd",
		},
	}, &map[string]interface{}{})
	assert.NoError(t, err)

	assert.NotEmpty(t, tokenParams.Get("client_secret"))
	tokenParams.Del("client_secret")

	assert.Equal(t, url.Values{
		"client_id":    {"test"},
		"grant_type":   {"authorization_code"},
		"code":         {"abcd"},
		"redirect_uri": {"redirect"},
	}, tokenParams)
}

func TestWithThirdPartyPasswordlessPassingAdditionalParamsCheckTheyArePresentInAuthorizationURLForThirdPartyProviderApple(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "apple",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID: clientId,
									AdditionalConfig: map[string]interface{}{
										"keyId":      "test-key",
										"privateKey": privateKey,
										"teamId":     "test-team-id",
									},
								},
							},
							AuthorizationEndpointQueryParams: map[string]interface{}{
								"key1": "value1",
								"key2": "value2",
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "apple", nil)
	assert.NoError(t, err)

	providerInfo := providerRes

	assert.Equal(t, "apple", providerInfo.ID)

	assert.Equal(t, "https://appleid.apple.com/auth/authorize", providerInfo.Config.AuthorizationEndpoint)
	assert.Equal(t, "https://appleid.apple.com/auth/token", providerInfo.Config.TokenEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"response_mode": {"form_post"},
		"response_type": {"code"},
		"redirect_uri":  {"redirect"},
		"scope":         {"openid email"},
		"key1":          {"value1"},
		"key2":          {"value2"},
	}, authParams)
}

func TestWithThirdPartyProviderPasswordlessPassingScopesInConfigForThirdPartyProviderApple(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "apple",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID: clientId,
									Scope:    []string{"test-scope-1", "test-scope-2"},
									AdditionalConfig: map[string]interface{}{
										"keyId":      "test-key",
										"privateKey": privateKey,
										"teamId":     "test-team-id",
									},
								},
							},
						},
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

	providerRes, err := ThirdPartyGetProvider("public", "apple", nil)
	assert.NoError(t, err)

	providerInfo := providerRes

	assert.Equal(t, "apple", providerInfo.ID)

	assert.Equal(t, "https://appleid.apple.com/auth/authorize", providerInfo.Config.AuthorizationEndpoint)
	assert.Equal(t, "https://appleid.apple.com/auth/token", providerInfo.Config.TokenEndpoint)

	authUrlRes, err := providerInfo.GetAuthorisationRedirectURL("redirect", &map[string]interface{}{})
	assert.NoError(t, err)

	urlObj, err := url.Parse(authUrlRes.URLWithQueryParams)
	assert.NoError(t, err)

	authParams := urlObj.Query()

	assert.Equal(t, url.Values{
		"client_id":     {"test"},
		"response_mode": {"form_post"},
		"response_type": {"code"},
		"redirect_uri":  {"redirect"},
		"scope":         {"test-scope-1 test-scope-2"},
	}, authParams)
}

func TestWithThirdPartyPasswordlessDuplicateProvider(t *testing.T) {
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
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
								},
							},
						},
					},
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "test",
									ClientSecret: "test-secret",
								},
							},
						},
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
		assert.Equal(t, "The providers array has multiple entries for the same third party provider.", err.Error())
	}
}
