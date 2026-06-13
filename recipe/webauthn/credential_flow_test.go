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
	"net/http/cookiejar"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

const testOrigin = "https://api.supertokens.io"

func toJSONMap(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()
	return *unittesting.HttpResponseToConsumableInformation(resp.Body)
}

func postJSONClient(t *testing.T, client *http.Client, url string, body interface{}) *http.Response {
	t.Helper()
	buf, err := json.Marshal(body)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(buf))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("rid", "webauthn")
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

func getReqClient(t *testing.T, client *http.Client, url string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	req.Header.Set("rid", "webauthn")
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

// registerOptions calls /options/register and returns (optionsId, challenge, rpId).
func registerOptions(t *testing.T, baseURL, email string) (string, string, string) {
	t.Helper()
	resp := postJSON(t, baseURL+"/auth/webauthn/options/register", map[string]interface{}{"email": email})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body := toJSONMap(t, resp)
	require.Equal(t, "OK", body["status"], "registerOptions failed: %v", body)
	rp := body["rp"].(map[string]interface{})
	return body["webauthnGeneratedOptionsId"].(string), body["challenge"].(string), rp["id"].(string)
}

// signInOptions calls /options/signin and returns (optionsId, challenge, rpId).
func signInOptions(t *testing.T, baseURL string) (string, string, string) {
	t.Helper()
	resp := postJSON(t, baseURL+"/auth/webauthn/options/signin", map[string]interface{}{})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body := toJSONMap(t, resp)
	require.Equal(t, "OK", body["status"], "signInOptions failed: %v", body)
	return body["webauthnGeneratedOptionsId"].(string), body["challenge"].(string), body["rpId"].(string)
}

type signedUpUser struct {
	auth         *virtualAuthenticator
	userID       string
	recipeUserID string
	email        string
	client       *http.Client // carries the session cookies set at sign up
}

// signUpNewUser runs the full register-options → sign-up ceremony with a fresh
// virtual authenticator and returns the resulting user (with an authenticated
// HTTP client).
func signUpNewUser(t *testing.T, baseURL, email string) signedUpUser {
	t.Helper()
	optionsID, challenge, rpID := registerOptions(t, baseURL, email)

	auth := newVirtualAuthenticator()
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := &http.Client{Jar: jar}

	resp := postJSONClient(t, client, baseURL+"/auth/webauthn/signup", map[string]interface{}{
		"webauthnGeneratedOptionsId": optionsID,
		"credential":                 auth.registrationResponse(rpID, challenge, testOrigin),
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body := toJSONMap(t, resp)
	require.Equal(t, "OK", body["status"], "signUp failed: %v", body)

	user := body["user"].(map[string]interface{})
	userID := user["id"].(string)
	loginMethods := user["loginMethods"].([]interface{})
	recipeUserID := loginMethods[0].(map[string]interface{})["recipeUserId"].(string)

	return signedUpUser{auth: auth, userID: userID, recipeUserID: recipeUserID, email: email, client: client}
}

func authenticationPayload(t *testing.T, m map[string]interface{}) webauthnmodels.AuthenticationPayload {
	t.Helper()
	buf, err := json.Marshal(m)
	require.NoError(t, err)
	var p webauthnmodels.AuthenticationPayload
	require.NoError(t, json.Unmarshal(buf, &p))
	return p
}

// Exercises the recipe-level credential functions against the real core after a
// genuine sign-up: GetCredential (incl. the restored UpdatedAt), ListCredentials,
// and RemoveCredential.
func TestRecipeCredentialFunctionsAfterSignUp(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	u := signUpNewUser(t, ts.URL, "creduser@example.com")
	credID := u.auth.credentialIDB64()

	// GetCredential — OK path with all fields (UpdatedAt must be populated).
	got, err := GetCredential(credID, u.recipeUserID, "public")
	require.NoError(t, err)
	require.NotNil(t, got.OK, "expected credential, got %+v", got)
	assert.Equal(t, credID, got.OK.Credential.WebauthnCredentialId)
	assert.Equal(t, u.recipeUserID, got.OK.Credential.RecipeUserId)
	assert.NotEmpty(t, got.OK.Credential.RelyingPartyId)
	assert.NotZero(t, got.OK.Credential.CreatedAt)
	assert.NotZero(t, got.OK.Credential.UpdatedAt, "UpdatedAt must be read from the core response")

	// ListCredentials — must contain the credential we just registered.
	listed, err := ListCredentials(u.recipeUserID)
	require.NoError(t, err)
	require.NotNil(t, listed.OK)
	found := false
	for _, c := range listed.OK.Credentials {
		if c.WebauthnCredentialId == credID {
			found = true
		}
	}
	assert.True(t, found, "ListCredentials must include the registered credential")

	// RemoveCredential — then GetCredential must report it gone.
	removed, err := RemoveCredential(credID, u.recipeUserID)
	require.NoError(t, err)
	assert.NotNil(t, removed.OK)

	gone, err := GetCredential(credID, u.recipeUserID, "public")
	require.NoError(t, err)
	assert.NotNil(t, gone.CredentialNotFoundError)
}

// VerifyCredentials must hit /recipe/webauthn/signin and succeed for a real
// assertion signed by the registered authenticator.
func TestVerifyCredentialsOKForRealAssertion(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	u := signUpNewUser(t, ts.URL, "verify@example.com")

	optionsID, challenge, rpID := signInOptions(t, ts.URL)
	assertion := authenticationPayload(t, u.auth.assertionResponse(rpID, challenge, testOrigin, u.userID))

	resp, err := VerifyCredentials(optionsID, assertion, "public")
	require.NoError(t, err)
	assert.NotNil(t, resp.OK, "real assertion must verify, got %+v", resp)
}

// Drives the session-protected API endpoints end to end: listCredentialsGET,
// signInPOST (OK path), registerCredentialPOST (incl. the email-match guard) and
// the reworked removeCredentialPOST.
func TestApiCredentialFlowWithSession(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	u := signUpNewUser(t, ts.URL, "apiflow@example.com")
	credID := u.auth.credentialIDB64()

	// listCredentialsGET (session-protected) lists the credential.
	listResp := getReqClient(t, u.client, ts.URL+"/auth/webauthn/credential/list")
	require.Equal(t, http.StatusOK, listResp.StatusCode)
	listBody := toJSONMap(t, listResp)
	require.Equal(t, "OK", listBody["status"], "listCredentials: %v", listBody)
	creds := listBody["credentials"].([]interface{})
	require.Len(t, creds, 1)
	// NOTE: the listCredentialsGET OK response currently emits PascalCase keys
	// (WebauthnCredentialId) because the handler marshals an anonymous struct
	// without json tags — a pre-existing FDI-conformance issue, unrelated to the
	// changes under test.
	assert.Equal(t, credID, creds[0].(map[string]interface{})["WebauthnCredentialId"])

	// signInPOST OK path with a real assertion.
	siOptionsID, siChallenge, siRpID := signInOptions(t, ts.URL)
	signinResp := postJSON(t, ts.URL+"/auth/webauthn/signin", map[string]interface{}{
		"webauthnGeneratedOptionsId": siOptionsID,
		"credential":                 u.auth.assertionResponse(siRpID, siChallenge, testOrigin, u.userID),
	})
	require.Equal(t, http.StatusOK, signinResp.StatusCode)
	assert.Equal(t, "OK", toJSONMap(t, signinResp)["status"])

	// registerCredentialPOST: a second credential for the SAME email passes the
	// email-match guard.
	rcOptionsID, rcChallenge, rcRpID := registerOptions(t, ts.URL, u.email)
	auth2 := newVirtualAuthenticator()
	regResp := postJSONClient(t, u.client, ts.URL+"/auth/webauthn/credential", map[string]interface{}{
		"webauthnGeneratedOptionsId": rcOptionsID,
		"credential":                 auth2.registrationResponse(rcRpID, rcChallenge, testOrigin),
	})
	require.Equal(t, http.StatusOK, regResp.StatusCode)
	assert.Equal(t, "OK", toJSONMap(t, regResp)["status"], "registerCredential should succeed for the user's own email")

	// removeCredentialPOST: the reworked scan finds the owning recipe user and
	// removes the original credential.
	rmResp := postJSONClient(t, u.client, ts.URL+"/auth/webauthn/credential/remove", map[string]interface{}{
		"webauthnCredentialId": credID,
	})
	require.Equal(t, http.StatusOK, rmResp.StatusCode)
	assert.Equal(t, "OK", toJSONMap(t, rmResp)["status"])
}

// registerOptions, when given a valid recoverAccountToken and no email, must
// resolve the token to the user's email in the recipe layer and return OK with
// that email — the account-recovery register-options path.
func TestRegisterOptionsResolvesRecoverAccountToken(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	u := signUpNewUser(t, ts.URL, "recover@example.com")

	instance, err := GetRecipeInstanceOrThrowError()
	require.NoError(t, err)
	tokenResp, err := (*instance.RecipeImpl.GenerateRecoverAccountToken)(u.recipeUserID, u.email, "public", &map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, tokenResp.OK)

	resp := postJSON(t, ts.URL+"/auth/webauthn/options/register", map[string]interface{}{
		"recoverAccountToken": tokenResp.OK.Token,
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body := toJSONMap(t, resp)
	require.Equal(t, "OK", body["status"], "registerOptions with valid token: %v", body)
	assert.Equal(t, u.email, body["user"].(map[string]interface{})["name"], "options must be generated for the token's email")
}

// registerCredentialPOST must reject options generated for a different email
// than the session user's, returning a GENERAL_ERROR (not registering it).
func TestRegisterCredentialRejectsEmailMismatch(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	u := signUpNewUser(t, ts.URL, "owner@example.com")

	// Options generated for a different email.
	rcOptionsID, rcChallenge, rcRpID := registerOptions(t, ts.URL, "someone-else@example.com")
	auth2 := newVirtualAuthenticator()
	resp := postJSONClient(t, u.client, ts.URL+"/auth/webauthn/credential", map[string]interface{}{
		"webauthnGeneratedOptionsId": rcOptionsID,
		"credential":                 auth2.registrationResponse(rcRpID, rcChallenge, testOrigin),
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body := toJSONMap(t, resp)
	assert.Equal(t, "GENERAL_ERROR", body["status"])
	assert.Equal(t, "Email mismatch", body["message"])
}
