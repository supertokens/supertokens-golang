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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

// initWebauthnWithSession initializes both webauthn and session recipes — the
// session recipe is needed for the credential management endpoints
// (registerCredential, listCredentials, removeCredential).
func initWebauthnWithSession(t *testing.T) {
	connectionURI := unittesting.StartUpST("localhost", "8080")
	cookieMethod := func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
		return sessmodels.CookieTransferMethod
	}
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "https://api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "https://api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: cookieMethod,
			}),
		},
	}
	if err := supertokens.Init(configValue); err != nil {
		t.Fatal(err.Error())
	}
}

func newTestServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	return httptest.NewServer(supertokens.Middleware(mux))
}

func postJSON(t *testing.T, url string, body interface{}) *http.Response {
	t.Helper()
	buf, err := json.Marshal(body)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(buf))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("rid", "webauthn")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func getReq(t *testing.T, url string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	req.Header.Set("rid", "webauthn")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

// -----------------------------------------------------------------------------
// registerOptions — POST /auth/webauthn/register/options
// -----------------------------------------------------------------------------

// Validates the FDI success response shape per
// https://supertokens.com/docs/references/fdi/webauthn (post-webauthn-register-options).
func TestRegisterOptionsWithValidEmailReturnsOKShape(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)

	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/register/options", map[string]interface{}{
		"email": "user@example.com",
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", result["status"])
	assert.NotEmpty(t, result["webauthnGeneratedOptionsId"])
	assert.NotEmpty(t, result["challenge"])
	assert.NotEmpty(t, result["createdAt"])
	assert.NotEmpty(t, result["expiresAt"])
	assert.NotZero(t, result["timeout"])
	assert.Contains(t, result, "attestation")

	rp, ok := result["rp"].(map[string]interface{})
	require.True(t, ok, "rp must be an object")
	assert.NotEmpty(t, rp["id"])
	assert.NotEmpty(t, rp["name"])

	user, ok := result["user"].(map[string]interface{})
	require.True(t, ok, "user must be an object")
	assert.NotEmpty(t, user["id"])
	assert.Equal(t, "user@example.com", user["name"])

	_, ok = result["excludeCredentials"].([]interface{})
	assert.True(t, ok, "excludeCredentials must be an array")

	pkc, ok := result["pubKeyCredParams"].([]interface{})
	require.True(t, ok, "pubKeyCredParams must be an array")
	assert.NotEmpty(t, pkc)

	authSel, ok := result["authenticatorSelection"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, authSel, "requireResidentKey")
	assert.Contains(t, authSel, "residentKey")
	assert.Contains(t, authSel, "userVerification")
}

// Per FDI, an invalid recoverAccountToken should return
// RECOVER_ACCOUNT_TOKEN_INVALID_ERROR (not surface a core error).
func TestRegisterOptionsWithInvalidRecoverAccountTokenReturnsTokenInvalidError(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)

	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/register/options", map[string]interface{}{
		"recoverAccountToken": "this-token-does-not-exist",
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "RECOVER_ACCOUNT_TOKEN_INVALID_ERROR", result["status"])
}

// -----------------------------------------------------------------------------
// signInOptions — POST /auth/webauthn/signin/options
// -----------------------------------------------------------------------------

// Validates FDI success shape for the signin options endpoint
// (post-webauthn-signin-options).
func TestSignInOptionsReturnsOKShape(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)

	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/signin/options", map[string]interface{}{})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)

	assert.Equal(t, "OK", result["status"])
	assert.NotEmpty(t, result["webauthnGeneratedOptionsId"])
	assert.NotEmpty(t, result["challenge"])
	assert.NotEmpty(t, result["createdAt"])
	assert.NotEmpty(t, result["expiresAt"])
	assert.NotEmpty(t, result["rpId"])
	assert.NotZero(t, result["timeout"])
	assert.Contains(t, result, "userVerification")
}

// -----------------------------------------------------------------------------
// signIn — POST /auth/webauthn/signin
// -----------------------------------------------------------------------------

func TestSignInRequiresCredential(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/signin", map[string]interface{}{
		"webauthnGeneratedOptionsId": "some-id",
	})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "credential is required", result["message"])
}

func TestSignInWithInvalidCredentialShapeReturnsInvalidCredentialsError(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/signin", map[string]interface{}{
		"webauthnGeneratedOptionsId": "some-id",
		"credential":                 "not-a-credential-object",
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "INVALID_CREDENTIALS_ERROR", result["status"])
}

// When the user submits a webauthnGeneratedOptionsId that the core has never
// seen, the FDI says we must respond with OPTIONS_NOT_FOUND_ERROR.
func TestSignInWithUnknownOptionsIdReturnsOptionsNotFoundError(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/signin", map[string]interface{}{
		"webauthnGeneratedOptionsId": "00000000-0000-0000-0000-000000000000",
		"credential": map[string]interface{}{
			"id":    "cred-id",
			"rawId": "cred-id",
			"response": map[string]interface{}{
				"clientDataJSON":    "irrelevant",
				"authenticatorData": "irrelevant",
				"signature":         "irrelevant",
				"userHandle":        "irrelevant",
			},
			"authenticatorAttachment": "platform",
			"type":                    "public-key",
			"clientExtensionResults":  map[string]interface{}{},
		},
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OPTIONS_NOT_FOUND_ERROR", result["status"])
}

// -----------------------------------------------------------------------------
// signUp — POST /auth/webauthn/signup
// -----------------------------------------------------------------------------

func TestSignUpRequiresWebauthnGeneratedOptionsId(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/signup", map[string]interface{}{
		"credential": map[string]interface{}{},
	})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "webauthnGeneratedOptionsId is required", result["message"])
}

func TestSignUpRequiresCredential(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/signup", map[string]interface{}{
		"webauthnGeneratedOptionsId": "some-id",
	})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "credential is required", result["message"])
}

// FDI: signup with an options id that the core does not recognise must return
// OPTIONS_NOT_FOUND_ERROR (no user is created, no panic).
func TestSignUpWithUnknownOptionsIdReturnsOptionsNotFoundError(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/signup", map[string]interface{}{
		"webauthnGeneratedOptionsId": "00000000-0000-0000-0000-000000000000",
		"credential": map[string]interface{}{
			"id":    "cred-id",
			"rawId": "cred-id",
			"response": map[string]interface{}{
				"clientDataJSON":    "irrelevant",
				"attestationObject": "irrelevant",
				"transports":        []string{"internal"},
				"userHandle":        "irrelevant",
			},
			"authenticatorAttachment": "platform",
			"type":                    "public-key",
			"clientExtensionResults":  map[string]interface{}{},
		},
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OPTIONS_NOT_FOUND_ERROR", result["status"])
}

// -----------------------------------------------------------------------------
// registerCredential — POST /auth/webauthn/credential
// -----------------------------------------------------------------------------

func TestRegisterCredentialRequiresWebauthnGeneratedOptionsId(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/credential", map[string]interface{}{
		"credential": map[string]interface{}{},
	})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "webauthnGeneratedOptionsId is required", result["message"])
}

func TestRegisterCredentialRequiresCredential(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/credential", map[string]interface{}{
		"webauthnGeneratedOptionsId": "some-id",
	})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "credential is required", result["message"])
}

// Invalid credential JSON is converted to INVALID_CREDENTIALS_ERROR by the
// API layer before the session check fires (the body is parsed first).
func TestRegisterCredentialWithInvalidCredentialShapeReturnsInvalidCredentialsError(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/credential", map[string]interface{}{
		"webauthnGeneratedOptionsId": "some-id",
		"credential":                 "not-a-credential-object",
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "INVALID_CREDENTIALS_ERROR", result["status"])
}

// Without a session cookie the request must fail authentication — the FDI
// describes credential management endpoints as session-protected.
func TestRegisterCredentialWithoutSessionReturnsUnauthorised(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/credential", map[string]interface{}{
		"webauthnGeneratedOptionsId": "some-id",
		"credential": map[string]interface{}{
			"id":                      "cred-id",
			"rawId":                   "cred-id",
			"response":                map[string]interface{}{},
			"authenticatorAttachment": "platform",
			"type":                    "public-key",
			"clientExtensionResults":  map[string]interface{}{},
		},
	})
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "unauthorised", result["message"])
}

// -----------------------------------------------------------------------------
// listCredentials — GET /auth/webauthn/credential/list
// -----------------------------------------------------------------------------

func TestListCredentialsWithoutSessionReturnsUnauthorised(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := getReq(t, ts.URL+"/auth/webauthn/credential/list")
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "unauthorised", result["message"])
}

// -----------------------------------------------------------------------------
// removeCredential — POST /auth/webauthn/credential/remove
// -----------------------------------------------------------------------------

// Body validation runs before the session check.
func TestRemoveCredentialRequiresWebauthnCredentialId(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/credential/remove", map[string]interface{}{})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "A valid webauthnCredentialId is required", result["message"])
}

func TestRemoveCredentialWithoutSessionReturnsUnauthorised(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/webauthn/credential/remove", map[string]interface{}{
		"webauthnCredentialId": "some-credential-id",
	})
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "unauthorised", result["message"])
}

// -----------------------------------------------------------------------------
// emailExists — GET /auth/webauthn/email/exists
// -----------------------------------------------------------------------------

// An unknown email must return {status: OK, exists: false} — the FDI never
// leaks UNKNOWN_USER_ID_ERROR back to the client for this endpoint.
func TestEmailExistsForUnknownEmailReturnsFalse(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := getReq(t, ts.URL+"/auth/webauthn/email/exists?email=unknown@example.com")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", result["status"])
	assert.Equal(t, false, result["exists"])
}

// -----------------------------------------------------------------------------
// generateRecoverAccountToken — POST /auth/user/webauthn/reset/token
// -----------------------------------------------------------------------------

func TestGenerateRecoverAccountTokenRequiresEmail(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/user/webauthn/reset/token", map[string]interface{}{})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "email is required", result["message"])
}

// For an unknown email, the FDI describes the response as 200 OK (no user
// enumeration leak). The API silently maps "user not found" to status OK.
func TestGenerateRecoverAccountTokenForUnknownEmailReturnsOK(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/user/webauthn/reset/token", map[string]interface{}{
		"email": "definitely-not-registered@example.com",
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", result["status"])
}

// -----------------------------------------------------------------------------
// recoverAccount — POST /auth/user/webauthn/reset
// -----------------------------------------------------------------------------

func TestRecoverAccountRequiresToken(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/user/webauthn/reset", map[string]interface{}{
		"webauthnGeneratedOptionsId": "some-id",
		"credential":                 map[string]interface{}{},
	})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "token is required", result["message"])
}

func TestRecoverAccountRequiresWebauthnGeneratedOptionsId(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/user/webauthn/reset", map[string]interface{}{
		"token":      "some-token",
		"credential": map[string]interface{}{},
	})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "webauthnGeneratedOptionsId is required", result["message"])
}

func TestRecoverAccountRequiresCredential(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/user/webauthn/reset", map[string]interface{}{
		"token":                      "some-token",
		"webauthnGeneratedOptionsId": "some-id",
	})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "credential is required", result["message"])
}

// An invalid recovery token must surface as RECOVER_ACCOUNT_TOKEN_INVALID_ERROR
// per the FDI — the API must not leak the underlying core error.
func TestRecoverAccountWithInvalidTokenReturnsTokenInvalidError(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)
	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/user/webauthn/reset", map[string]interface{}{
		"token":                      "definitely-not-a-real-token",
		"webauthnGeneratedOptionsId": "00000000-0000-0000-0000-000000000000",
		"credential": map[string]interface{}{
			"id":    "cred-id",
			"rawId": "cred-id",
			"response": map[string]interface{}{
				"clientDataJSON":    "irrelevant",
				"attestationObject": "irrelevant",
				"transports":        []string{"internal"},
			},
			"authenticatorAttachment": "platform",
			"type":                    "public-key",
			"clientExtensionResults":  map[string]interface{}{},
		},
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "RECOVER_ACCOUNT_TOKEN_INVALID_ERROR", result["status"])
}
