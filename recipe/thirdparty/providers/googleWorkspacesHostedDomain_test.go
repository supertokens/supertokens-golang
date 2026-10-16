/* Copyright (c) 2026, VRAI Labs and/or its affiliates. All rights reserved.
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

package providers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

// These tests call the Google Workspaces provider directly with a fake OIDC provider, so they
// don't need a running core.

// newGoogleWorkspacesUserInfoServerForTest serves OIDC discovery info (without a jwks_uri) and a
// user info API that returns *userInfo.
func newGoogleWorkspacesUserInfoServerForTest(userInfo *map[string]interface{}) *httptest.Server {
	var server *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"issuer":                 server.URL,
			"authorization_endpoint": server.URL + "/auth",
			"token_endpoint":         server.URL + "/token",
			"userinfo_endpoint":      server.URL + "/userinfo",
		})
	})
	mux.HandleFunc("/userinfo", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(*userInfo)
	})
	server = httptest.NewServer(mux)
	return server
}

func createGoogleWorkspacesProviderForHdTest(t *testing.T, server *httptest.Server, additionalConfig map[string]interface{}) *tpmodels.TypeProvider {
	provider, err := FindAndCreateProviderInstance([]tpmodels.ProviderInput{
		{
			Config: tpmodels.ProviderConfig{
				ThirdPartyId:          "google-workspaces",
				OIDCDiscoveryEndpoint: server.URL + "/.well-known/openid-configuration",
				Clients: []tpmodels.ProviderClientConfig{
					{
						ClientID:         "test-client-id",
						ClientSecret:     "test-secret",
						AdditionalConfig: additionalConfig,
					},
				},
			},
		},
	}, "google-workspaces", nil, &map[string]interface{}{})
	if !assert.NoError(t, err) || !assert.NotNil(t, provider) {
		t.FailNow()
	}
	return provider
}

func getGoogleWorkspacesUserInfoForHdTest(provider *tpmodels.TypeProvider) (tpmodels.TypeUserInfo, error) {
	return provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "test-access-token"}, &map[string]interface{}{})
}

func TestGoogleWorkspacesWithEmptyHostedDomainRejectsAccessTokenOnlySignIn(t *testing.T) {
	userInfo := map[string]interface{}{}
	server := newGoogleWorkspacesUserInfoServerForTest(&userInfo)
	defer server.Close()
	provider := createGoogleWorkspacesProviderForHdTest(t, server, map[string]interface{}{"hd": ""})

	// The authorisation redirect still lets the user pick any account
	assert.Equal(t, "*", provider.Config.AuthorizationEndpointQueryParams["hd"])

	for _, testUserInfo := range []map[string]interface{}{
		// A personal account: no hd in the user info response
		{
			"sub":            "personal-user",
			"email":          "user@gmail.com",
			"email_verified": true,
		},
		// A Google Workspace account
		{
			"sub":            "workspace-user",
			"email":          "user@example.com",
			"email_verified": true,
			"hd":             "example.com",
		},
	} {
		userInfo = testUserInfo
		_, err := getGoogleWorkspacesUserInfoForHdTest(provider)
		assert.Error(t, err, testUserInfo["sub"])
	}

	// The id_token check treats the empty hd the same way
	for _, idTokenPayload := range []map[string]interface{}{
		{"sub": "personal-user"},
		{"sub": "workspace-user", "hd": "example.com"},
	} {
		err := provider.Config.ValidateIdTokenPayload(idTokenPayload, provider.Config, &map[string]interface{}{})
		assert.Error(t, err, idTokenPayload["sub"])
	}
}

func TestGoogleWorkspacesWithoutSpecificHostedDomainAllowsAccessTokenOnlySignIn(t *testing.T) {
	userInfo := map[string]interface{}{
		"sub":            "personal-user",
		"email":          "user@gmail.com",
		"email_verified": true,
	}
	server := newGoogleWorkspacesUserInfoServerForTest(&userInfo)
	defer server.Close()

	for _, additionalConfig := range []map[string]interface{}{
		nil,
		{},
		{"hd": nil},
		{"hd": "*"},
	} {
		provider := createGoogleWorkspacesProviderForHdTest(t, server, additionalConfig)
		assert.Equal(t, "*", provider.Config.AuthorizationEndpointQueryParams["hd"], additionalConfig)

		result, err := getGoogleWorkspacesUserInfoForHdTest(provider)
		if assert.NoError(t, err, additionalConfig) {
			assert.Equal(t, "personal-user", result.ThirdPartyUserId, additionalConfig)
		}
	}
}

func TestGoogleWorkspacesWithSpecificHostedDomainChecksTheUserInfoHd(t *testing.T) {
	userInfo := map[string]interface{}{}
	server := newGoogleWorkspacesUserInfoServerForTest(&userInfo)
	defer server.Close()
	provider := createGoogleWorkspacesProviderForHdTest(t, server, map[string]interface{}{"hd": "example.com"})
	assert.Equal(t, "example.com", provider.Config.AuthorizationEndpointQueryParams["hd"])

	userInfo = map[string]interface{}{
		"sub":            "workspace-user",
		"email":          "user@example.com",
		"email_verified": true,
		"hd":             "example.com",
	}
	result, err := getGoogleWorkspacesUserInfoForHdTest(provider)
	if assert.NoError(t, err) {
		assert.Equal(t, "workspace-user", result.ThirdPartyUserId)
	}

	for _, hd := range []interface{}{nil, "", "other.com", "EXAMPLE.COM"} {
		userInfo = map[string]interface{}{
			"sub":            "other-user",
			"email":          "user@other.com",
			"email_verified": true,
		}
		if hd != nil {
			userInfo["hd"] = hd
		}
		_, err := getGoogleWorkspacesUserInfoForHdTest(provider)
		assert.Error(t, err, hd)
	}
}
