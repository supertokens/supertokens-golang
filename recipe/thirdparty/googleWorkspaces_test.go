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

package thirdparty

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

type googleWorkspacesTestProvider struct {
	server       *httptest.Server
	key          *rsa.PrivateKey
	userInfoBody map[string]interface{}
}

func newGoogleWorkspacesTestProvider(t *testing.T) *googleWorkspacesTestProvider {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	p := &googleWorkspacesTestProvider{key: key}

	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"issuer":                 p.server.URL,
			"authorization_endpoint": p.server.URL + "/auth",
			"token_endpoint":         p.server.URL + "/token",
			"userinfo_endpoint":      p.server.URL + "/userinfo",
			"jwks_uri":               p.server.URL + "/jwks",
		})
	})
	mux.HandleFunc("/jwks", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"keys": []interface{}{
				map[string]interface{}{
					"kty": "RSA",
					"alg": "RS256",
					"use": "sig",
					"kid": "test-key",
					"n":   base64.RawURLEncoding.EncodeToString(key.PublicKey.N.Bytes()),
					"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.PublicKey.E)).Bytes()),
				},
			},
		})
	})
	mux.HandleFunc("/userinfo", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(p.userInfoBody)
	})
	p.server = httptest.NewServer(mux)
	return p
}

func (p *googleWorkspacesTestProvider) signIdToken(t *testing.T, claims map[string]interface{}) string {
	mapClaims := jwt.MapClaims{
		"iss": p.server.URL,
		"aud": "test-client-id",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	for k, v := range claims {
		mapClaims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, mapClaims)
	token.Header["kid"] = "test-key"
	signed, err := token.SignedString(p.key)
	assert.NoError(t, err)
	return signed
}

func initGoogleWorkspacesForTest(t *testing.T, p *googleWorkspacesTestProvider, hd interface{}) *httptest.Server {
	connectionURI := unittesting.StartUpST("localhost", "8080")
	var additionalConfig map[string]interface{}
	if hd != nil {
		additionalConfig = map[string]interface{}{"hd": hd}
	}
	return supertokensInitForTest(t, connectionURI,
		session.Init(nil),
		Init(&tpmodels.TypeInput{
			SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
				Providers: []tpmodels.ProviderInput{
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId:          "google-workspaces",
							OIDCDiscoveryEndpoint: p.server.URL + "/.well-known/openid-configuration",
							// These tests are about the hd check on the user info, so they accept any access token
							ValidateAccessToken: func(accessToken string, clientConfig tpmodels.ProviderConfigForClientType, userContext supertokens.UserContext) error {
								return nil
							},
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:         "test-client-id",
									ClientSecret:     "test-secret",
									AdditionalConfig: additionalConfig,
								},
							},
						},
					},
				},
			},
		}),
	)
}

func signInUpWithOAuthTokensForTest(t *testing.T, testServer *httptest.Server, oAuthTokens map[string]interface{}) (int, map[string]interface{}) {
	body, err := json.Marshal(map[string]interface{}{
		"thirdPartyId": "google-workspaces",
		"oAuthTokens":  oAuthTokens,
	})
	assert.NoError(t, err)
	res, err := http.Post(testServer.URL+"/auth/signinup", "application/json", strings.NewReader(string(body)))
	assert.NoError(t, err)
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	result := map[string]interface{}{}
	_ = json.Unmarshal(data, &result)
	return res.StatusCode, result
}

func TestGoogleWorkspacesRejectsAccessTokenOnlySignInForUserOutsideHostedDomain(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	p := newGoogleWorkspacesTestProvider(t)
	defer p.server.Close()
	testServer := initGoogleWorkspacesForTest(t, p, "example.com")
	defer testServer.Close()

	// No hd in the user info response
	p.userInfoBody = map[string]interface{}{
		"sub":            "outside-user-1",
		"email":          "outsider@gmail.com",
		"email_verified": true,
	}
	status, result := signInUpWithOAuthTokensForTest(t, testServer, map[string]interface{}{"access_token": "access-token"})
	assert.NotEqual(t, 200, status)
	assert.NotEqual(t, "OK", result["status"])

	// A different hd in the user info response
	p.userInfoBody = map[string]interface{}{
		"sub":            "outside-user-2",
		"email":          "outsider@other.com",
		"email_verified": true,
		"hd":             "other.com",
	}
	status, result = signInUpWithOAuthTokensForTest(t, testServer, map[string]interface{}{"access_token": "access-token"})
	assert.NotEqual(t, 200, status)
	assert.NotEqual(t, "OK", result["status"])
}

func TestGoogleWorkspacesAcceptsAccessTokenOnlySignInForUserInHostedDomain(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	p := newGoogleWorkspacesTestProvider(t)
	defer p.server.Close()
	testServer := initGoogleWorkspacesForTest(t, p, "example.com")
	defer testServer.Close()

	p.userInfoBody = map[string]interface{}{
		"sub":            "inside-user",
		"email":          "insider@example.com",
		"email_verified": true,
		"hd":             "example.com",
	}
	status, result := signInUpWithOAuthTokensForTest(t, testServer, map[string]interface{}{"access_token": "access-token"})
	assert.Equal(t, 200, status)
	assert.Equal(t, "OK", result["status"])
}

func TestGoogleWorkspacesWithoutSpecificHostedDomainAcceptsAccessTokenOnlySignIn(t *testing.T) {
	for _, hd := range []interface{}{nil, "*"} {
		BeforeEach()
		p := newGoogleWorkspacesTestProvider(t)
		testServer := initGoogleWorkspacesForTest(t, p, hd)

		p.userInfoBody = map[string]interface{}{
			"sub":            "any-user",
			"email":          "user@gmail.com",
			"email_verified": true,
		}
		status, result := signInUpWithOAuthTokensForTest(t, testServer, map[string]interface{}{"access_token": "access-token"})
		assert.Equal(t, 200, status)
		assert.Equal(t, "OK", result["status"])

		testServer.Close()
		p.server.Close()
		AfterEach()
	}
}

func TestGoogleWorkspacesChecksHostedDomainOfVerifiedIdToken(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	p := newGoogleWorkspacesTestProvider(t)
	defer p.server.Close()
	testServer := initGoogleWorkspacesForTest(t, p, "example.com")
	defer testServer.Close()

	// The user info response has no hd, the verified id_token carries the matching hd
	p.userInfoBody = map[string]interface{}{
		"sub":            "inside-user",
		"email":          "insider@example.com",
		"email_verified": true,
	}
	idToken := p.signIdToken(t, map[string]interface{}{
		"sub":            "inside-user",
		"email":          "insider@example.com",
		"email_verified": true,
		"hd":             "example.com",
	})
	status, result := signInUpWithOAuthTokensForTest(t, testServer, map[string]interface{}{"access_token": "access-token", "id_token": idToken})
	assert.Equal(t, 200, status)
	assert.Equal(t, "OK", result["status"])

	idToken = p.signIdToken(t, map[string]interface{}{
		"sub":            "outside-user",
		"email":          "outsider@other.com",
		"email_verified": true,
		"hd":             "other.com",
	})
	status, result = signInUpWithOAuthTokensForTest(t, testServer, map[string]interface{}{"access_token": "access-token", "id_token": idToken})
	assert.NotEqual(t, 200, status)
	assert.NotEqual(t, "OK", result["status"])
}
