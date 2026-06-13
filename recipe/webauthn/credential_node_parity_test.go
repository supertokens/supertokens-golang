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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
)

// bogusAuthenticationPayload builds a well-formed but meaningless credential,
// enough to exercise the verify path without real WebAuthn crypto.
func bogusAuthenticationPayload() webauthnmodels.AuthenticationPayload {
	return webauthnmodels.AuthenticationPayload{
		CredentialPayloadBase: webauthnmodels.CredentialPayloadBase{
			ID:                     "cred-id",
			RawID:                  "cred-id",
			Type:                   "public-key",
			ClientExtensionResults: map[string]interface{}{},
		},
		Response: webauthnmodels.AuthenticatorAssertionResponseJSON{
			ClientDataJSON:    "irrelevant",
			AuthenticatorData: "irrelevant",
			Signature:         "irrelevant",
		},
	}
}

// VerifyCredentials must reach a real core endpoint (not a 404) and must not
// report success for a bogus credential.
func TestVerifyCredentialsReachesSignInEndpoint(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)

	resp, err := VerifyCredentials(
		"00000000-0000-0000-0000-000000000000",
		bogusAuthenticationPayload(),
		"public",
	)
	require.NoError(t, err, "verifyCredentials must hit a real core endpoint, not 404")
	assert.Nil(t, resp.OK, "a bogus credential must not verify")
}

// An unknown credential must return CREDENTIAL_NOT_FOUND_ERROR with no Go error.
func TestGetCredentialUnknownReturnsNotFound(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnRecipeForAPITests(t)

	resp, err := GetCredential(
		"non-existent-credential-id",
		"00000000-0000-0000-0000-000000000000",
		"public",
	)
	require.NoError(t, err)
	assert.NotNil(t, resp.CredentialNotFoundError, "unknown credential must map to CREDENTIAL_NOT_FOUND_ERROR")
	assert.Nil(t, resp.OK)
}
