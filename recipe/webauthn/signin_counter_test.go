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

package webauthn

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Guards against the double-verification bug that broke passkey sign-in in the
// Node and Python SDKs (https://github.com/supertokens/supertokens-core/issues/1195):
// their sign-in API handler verified the same assertion against the core twice,
// and since the core persists the signature counter on every verification, the
// second call presented a non-increasing signCount and tripped the core's
// clone detection ("Malicious counter value is detected...").
//
// The bug is only observable with authenticators that increment the counter
// (Windows Hello, security keys, Chrome virtual authenticators) — Apple/Google
// passkeys keep it at 0, which skips the check. This test simulates the
// incrementing kind across two consecutive logins, so it fails if signInPOST
// ever starts verifying an assertion against the core more than once.
func TestSignInAPIWithCounterIncrementingAuthenticator(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initWebauthnWithSession(t)
	ts := newTestServer(t)
	defer ts.Close()

	u := signUpNewUser(t, ts.URL, "counter-canary@example.com")

	// Real hardware increments the signature counter on every assertion.
	for _, count := range []uint32{1, 2} {
		u.auth.signCount = count

		optionsID, challenge, rpID := signInOptions(t, ts.URL)
		resp := postJSON(t, ts.URL+"/auth/webauthn/signin", map[string]interface{}{
			"webauthnGeneratedOptionsId": optionsID,
			"credential":                 u.auth.assertionResponse(rpID, challenge, testOrigin, u.userID),
		})
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body := toJSONMap(t, resp)
		require.Equal(t, "OK", body["status"],
			"sign in with signCount=%d must succeed: a single API sign-in presents the assertion "+
				"to the core exactly once; INVALID_CREDENTIALS_ERROR here means the handler verified "+
				"it twice and tripped the core's clone detection (supertokens-core#1195): %v",
			count, body)
	}
}
