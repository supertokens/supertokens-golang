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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

// initWebauthnRecipeWithEmailCapture initialises the recipe with a custom email
// delivery service that records every WebauthnRecoverAccount email it is asked
// to send, so recover-account tests can assert on the email side-effect without
// hitting an external service.
func initWebauthnRecipeWithEmailCapture(t *testing.T, captured *[]emaildelivery.WebauthnRecoverAccountType) {
	connectionURI := unittesting.StartUpST("localhost", "8080")

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.WebauthnRecoverAccount != nil {
			*captured = append(*captured, *input.WebauthnRecoverAccount)
		}
		return nil
	}
	emailService := emaildelivery.EmailDeliveryInterface{SendEmail: &sendEmail}

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
			Init(&webauthnmodels.TypeInput{
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &emailService,
				},
			}),
		},
	}

	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
}

// A recover-account request for an email with no matching user must still
// return OK (to avoid leaking which emails are registered) but must NOT
// dispatch a recover-account email through the email delivery service.
func TestGenerateRecoverAccountTokenDoesNotSendEmailForUnknownEmail(t *testing.T) {
	BeforeEach()
	defer AfterEach()

	captured := []emaildelivery.WebauthnRecoverAccountType{}
	initWebauthnRecipeWithEmailCapture(t, &captured)

	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/user/webauthn/reset/token", map[string]interface{}{
		"email": "definitely-not-registered@example.com",
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "OK", result["status"])

	// No user resolved -> the email delivery service must not be invoked.
	assert.Empty(t, captured)
}

// A malformed recover-account request (missing email) must be rejected before
// any email is dispatched.
func TestGenerateRecoverAccountTokenDoesNotSendEmailWhenEmailMissing(t *testing.T) {
	BeforeEach()
	defer AfterEach()

	captured := []emaildelivery.WebauthnRecoverAccountType{}
	initWebauthnRecipeWithEmailCapture(t, &captured)

	ts := newTestServer(t)
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/auth/user/webauthn/reset/token", map[string]interface{}{})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	result := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	assert.Equal(t, "email is required", result["message"])

	assert.Empty(t, captured)
}
