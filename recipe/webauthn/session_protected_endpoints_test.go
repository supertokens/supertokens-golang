/* Copyright (c) 2025, VRAI Labs and/or its affiliates. All rights reserved.
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

package webauthn

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supertokens/supertokens-golang/test/unittesting"
)

// The session-protected webauthn endpoints (registerCredential,
// listCredentials, removeCredential) all route through
// getSessionWithoutClaimValidation, which calls session.GetSession with a
// VerifySessionOptions whose SessionRequired field is left nil. nil means
// REQUIRED (matching supertokens-node / supertokens-python), so these
// endpoints must respond with a clean 401 Unauthorized — not panic / 500 —
// even when a present-but-malformed session token is sent.
//
// The "no session token" variants of these endpoints live in
// api_corner_cases_test.go. The cases below cover the *present-but-invalid*
// token path specifically: the session recipe used to dereference a nil
// *bool (options.SessionRequired) on that path and panic.

// sendWithCookie issues a request to url with an explicit Cookie header so a
// malformed access token can be attached. The existing postJSON/getReq
// helpers don't allow custom headers.
func sendWithCookie(t *testing.T, method, url string, body interface{}, cookie string) *http.Response {
	t.Helper()
	var buf *bytes.Buffer
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req, err := http.NewRequest(method, url, buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("rid", "webauthn")
	req.Header.Set("Cookie", cookie)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

const malformedAccessTokenCookie = "sAccessToken=this-is-not-a-valid-jwt"

func assertCleanUnauthorised(t *testing.T, resp *http.Response) {
	t.Helper()
	// A clean 401 — not a panic surfaced as 500.
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode, "expected 401 Unauthorized, got %d", resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "unauthorised", result["message"])
}

// A registerCredential body that passes body validation so execution reaches
// the getSessionWithoutClaimValidation call.
func validRegisterCredentialBody() map[string]interface{} {
	return map[string]interface{}{
		"webauthnGeneratedOptionsId": "some-id",
		"credential": map[string]interface{}{
			"id":                      "cred-id",
			"rawId":                   "cred-id",
			"response":                map[string]interface{}{},
			"authenticatorAttachment": "platform",
			"type":                    "public-key",
			"clientExtensionResults":  map[string]interface{}{},
		},
	}
}

func TestRegisterCredentialWithMalformedTokenReturnsUnauthorised(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := sendWithCookie(t, http.MethodPost, ts.URL+"/auth/webauthn/credential", validRegisterCredentialBody(), malformedAccessTokenCookie)
	assertCleanUnauthorised(t, resp)
}

func TestListCredentialsWithMalformedTokenReturnsUnauthorised(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := sendWithCookie(t, http.MethodGet, ts.URL+"/auth/webauthn/credential/list", nil, malformedAccessTokenCookie)
	assertCleanUnauthorised(t, resp)
}

func TestRemoveCredentialWithMalformedTokenReturnsUnauthorised(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := sendWithCookie(t, http.MethodPost, ts.URL+"/auth/webauthn/credential/remove", map[string]interface{}{
		"webauthnCredentialId": "some-credential-id",
	}, malformedAccessTokenCookie)
	assertCleanUnauthorised(t, resp)
}
