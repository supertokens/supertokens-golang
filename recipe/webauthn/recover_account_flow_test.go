package webauthn

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// markEmailAsVerifiedInCore verifies the email of the recipe user directly in the core.
func markEmailAsVerifiedInCore(t *testing.T, recipeUserId string, email string) {
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	require.NoError(t, err)
	tokenResp, err := querier.SendPostRequest("public/recipe/user/email/verify/token", map[string]interface{}{
		"userId": recipeUserId,
		"email":  email,
	}, nil)
	require.NoError(t, err)
	require.Equal(t, "OK", tokenResp["status"])
	verifyResp, err := querier.SendPostRequest("public/recipe/user/email/verify", map[string]interface{}{
		"method": "token",
		"token":  tokenResp["token"],
	}, nil)
	require.NoError(t, err)
	require.Equal(t, "OK", verifyResp["status"])
}

func generateRecoverAccountTokenForTest(t *testing.T, u signedUpUser) string {
	instance, err := GetRecipeInstanceOrThrowError()
	require.NoError(t, err)
	tokenResp, err := (*instance.RecipeImpl.GenerateRecoverAccountToken)(u.recipeUserID, u.email, "public", &map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, tokenResp.OK)
	return tokenResp.OK.Token
}

// registerOptionsForRecovery calls /options/register with a recover account token and returns
// (optionsId, challenge, rpId). This does not consume the token.
func registerOptionsForRecovery(t *testing.T, baseURL string, token string) (string, string, string) {
	resp := postJSON(t, baseURL+"/auth/webauthn/options/register", map[string]interface{}{"recoverAccountToken": token})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body := toJSONMap(t, resp)
	require.Equal(t, "OK", body["status"], "registerOptions failed: %v", body)
	rp := body["rp"].(map[string]interface{})
	return body["webauthnGeneratedOptionsId"].(string), body["challenge"].(string), rp["id"].(string)
}

func recoverAccount(t *testing.T, baseURL string, token string, optionsId string, credential map[string]interface{}) map[string]interface{} {
	resp := postJSON(t, baseURL+"/auth/user/webauthn/reset", map[string]interface{}{
		"token":                      token,
		"webauthnGeneratedOptionsId": optionsId,
		"credential":                 credential,
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return toJSONMap(t, resp)
}

func listCredentialIds(t *testing.T, recipeUserId string) []string {
	listed, err := ListCredentials(recipeUserId)
	require.NoError(t, err)
	require.NotNil(t, listed.OK)
	ids := []string{}
	for _, c := range listed.OK.Credentials {
		ids = append(ids, c.WebauthnCredentialId)
	}
	return ids
}

func TestRecoverAccountRejectsOptionsGeneratedForAnotherEmail(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	u := signUpNewUser(t, ts.URL, "owner@example.com")
	token := generateRecoverAccountTokenForTest(t, u)

	optionsId, challenge, rpId := registerOptions(t, ts.URL, "someone-else@example.com")
	newAuth := newVirtualAuthenticator()
	body := recoverAccount(t, ts.URL, token, optionsId, newAuth.registrationResponse(rpId, challenge, testOrigin))
	assert.Equal(t, "INVALID_OPTIONS_ERROR", body["status"])

	assert.Equal(t, []string{u.auth.credentialIDB64()}, listCredentialIds(t, u.recipeUserID))
}

func TestRecoverAccountOfUnverifiedUserRemovesOlderCredentialsAndSessions(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	// The account (and its credential and session) was created without the email being proven
	existing := signUpNewUser(t, ts.URL, "unverified@example.com")
	handles, err := session.GetAllSessionHandlesForUser(existing.userID, nil)
	require.NoError(t, err)
	require.Len(t, handles, 1)

	token := generateRecoverAccountTokenForTest(t, existing)
	optionsId, challenge, rpId := registerOptionsForRecovery(t, ts.URL, token)
	newAuth := newVirtualAuthenticator()
	body := recoverAccount(t, ts.URL, token, optionsId, newAuth.registrationResponse(rpId, challenge, testOrigin))
	require.Equal(t, "OK", body["status"], "recoverAccount failed: %v", body)

	assert.Equal(t, []string{newAuth.credentialIDB64()}, listCredentialIds(t, existing.recipeUserID))

	handles, err = session.GetAllSessionHandlesForUser(existing.userID, nil)
	require.NoError(t, err)
	assert.Empty(t, handles)
}

func TestSessionRevokedByRecoveryCannotRegisterCredentials(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	// Someone signed up with the owner's email (without proving it) and kept the session
	earlierSignUp := signUpNewUser(t, ts.URL, "owner@example.com")

	// The owner recovers the account
	token := generateRecoverAccountTokenForTest(t, earlierSignUp)
	optionsId, challenge, rpId := registerOptionsForRecovery(t, ts.URL, token)
	ownerAuth := newVirtualAuthenticator()
	body := recoverAccount(t, ts.URL, token, optionsId, ownerAuth.registrationResponse(rpId, challenge, testOrigin))
	require.Equal(t, "OK", body["status"], "recoverAccount failed: %v", body)

	// The access token of the revoked session must not be usable to add a credential afterwards
	rcOptionsId, rcChallenge, rcRpId := registerOptions(t, ts.URL, earlierSignUp.email)
	resp := postJSONClient(t, earlierSignUp.client, ts.URL+"/auth/webauthn/credential", map[string]interface{}{
		"webauthnGeneratedOptionsId": rcOptionsId,
		"credential":                 newVirtualAuthenticator().registrationResponse(rcRpId, rcChallenge, testOrigin),
	})
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	assert.Equal(t, []string{ownerAuth.credentialIDB64()}, listCredentialIds(t, earlierSignUp.recipeUserID))
}

func TestRecoverAccountOfVerifiedUserAddsCredential(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	existing := signUpNewUser(t, ts.URL, "verified@example.com")
	markEmailAsVerifiedInCore(t, existing.recipeUserID, existing.email)

	token := generateRecoverAccountTokenForTest(t, existing)
	optionsId, challenge, rpId := registerOptionsForRecovery(t, ts.URL, token)
	newAuth := newVirtualAuthenticator()
	body := recoverAccount(t, ts.URL, token, optionsId, newAuth.registrationResponse(rpId, challenge, testOrigin))
	require.Equal(t, "OK", body["status"], "recoverAccount failed: %v", body)

	assert.ElementsMatch(t, []string{existing.auth.credentialIDB64(), newAuth.credentialIDB64()}, listCredentialIds(t, existing.recipeUserID))

	handles, err := session.GetAllSessionHandlesForUser(existing.userID, nil)
	require.NoError(t, err)
	assert.Len(t, handles, 1)

	// The token is spent
	body = recoverAccount(t, ts.URL, token, optionsId, newVirtualAuthenticator().registrationResponse(rpId, challenge, testOrigin))
	assert.Equal(t, "RECOVER_ACCOUNT_TOKEN_INVALID_ERROR", body["status"])
}

func TestRecoverAccountOfUnverifiedUserRevokesSessionsEvenIfRegistrationFails(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	existing := signUpNewUser(t, ts.URL, "unverified@example.com")
	handles, err := session.GetAllSessionHandlesForUser(existing.userID, nil)
	require.NoError(t, err)
	require.Len(t, handles, 1)

	// The recovery token is valid, but the credential doesn't match the options
	token := generateRecoverAccountTokenForTest(t, existing)
	optionsId, _, rpId := registerOptionsForRecovery(t, ts.URL, token)
	body := recoverAccount(t, ts.URL, token, optionsId, newVirtualAuthenticator().registrationResponse(rpId, "not-the-challenge", testOrigin))
	require.NotEqual(t, "OK", body["status"])

	// The sessions created before the email was proven are revoked anyway
	handles, err = session.GetAllSessionHandlesForUser(existing.userID, nil)
	require.NoError(t, err)
	assert.Empty(t, handles)
	assert.Equal(t, []string{existing.auth.credentialIDB64()}, listCredentialIds(t, existing.recipeUserID))
}
