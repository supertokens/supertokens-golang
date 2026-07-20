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

package supertokens

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUserJSONShapeMatchesCore locks in the JSON serialization of the User
// object so it matches the shape emitted by supertokens-node/core and consumed
// by web-js (thirdParty[], webauthn.credentialIds, and the per-loginMethod
// optional thirdParty / webauthn fields).
func TestUserJSONShapeMatchesCore(t *testing.T) {
	email := "test@example.com"
	user := User{
		ID:            "primary-user-id",
		TimeJoined:    1000,
		IsPrimaryUser: true,
		TenantIDs:     []string{"public"},
		Emails:        []string{email},
		PhoneNumbers:  []string{},
		ThirdParty:    []ThirdParty{{ID: "google", UserID: "google-user-id"}},
		Webauthn:      &Webauthn{CredentialIds: []string{"cred-1"}},
		LoginMethods: []LoginMethod{
			{
				RecipeID:     "webauthn",
				RecipeUserID: "recipe-user-id",
				Verified:     true,
				TenantIDs:    []string{"public"},
				TimeJoined:   1000,
				Email:        &email,
				Webauthn:     &Webauthn{CredentialIds: []string{"cred-1"}},
			},
		},
	}

	data, err := json.Marshal(user)
	assert.NoError(t, err)

	var got map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &got))

	// Top-level fields.
	assert.Equal(t, "primary-user-id", got["id"])
	assert.Equal(t, true, got["isPrimaryUser"])
	assert.Equal(t, []interface{}{"public"}, got["tenantIds"])
	assert.Equal(t, []interface{}{email}, got["emails"])
	assert.Equal(t, []interface{}{}, got["phoneNumbers"])
	assert.Equal(t, []interface{}{map[string]interface{}{"id": "google", "userId": "google-user-id"}}, got["thirdParty"])
	assert.Equal(t, map[string]interface{}{"credentialIds": []interface{}{"cred-1"}}, got["webauthn"])

	// Login method fields.
	loginMethods := got["loginMethods"].([]interface{})
	assert.Len(t, loginMethods, 1)
	lm := loginMethods[0].(map[string]interface{})
	assert.Equal(t, "webauthn", lm["recipeId"])
	assert.Equal(t, "recipe-user-id", lm["recipeUserId"])
	assert.Equal(t, email, lm["email"])
	assert.Equal(t, map[string]interface{}{"credentialIds": []interface{}{"cred-1"}}, lm["webauthn"])

	// Optional fields must be omitted when unset (matches web-js optional keys).
	_, hasThirdParty := lm["thirdParty"]
	assert.False(t, hasThirdParty, "loginMethod.thirdParty should be omitted when nil")
	_, hasPhoneNumber := lm["phoneNumber"]
	assert.False(t, hasPhoneNumber, "loginMethod.phoneNumber should be omitted when nil")
}

// TestUserJSONRoundTripFromCore ensures a core-style user payload (with a
// thirdparty login method) unmarshals into the User struct and marshals back
// with the same key names.
func TestUserJSONRoundTripFromCore(t *testing.T) {
	corePayload := `{
		"id": "primary-user-id",
		"timeJoined": 1000,
		"isPrimaryUser": true,
		"tenantIds": ["public"],
		"emails": ["test@example.com"],
		"phoneNumbers": [],
		"thirdParty": [{"id": "google", "userId": "google-user-id"}],
		"loginMethods": [{
			"recipeId": "thirdparty",
			"recipeUserId": "recipe-user-id",
			"verified": true,
			"tenantIds": ["public"],
			"timeJoined": 1000,
			"email": "test@example.com",
			"thirdParty": {"id": "google", "userId": "google-user-id"}
		}]
	}`

	var user User
	assert.NoError(t, json.Unmarshal([]byte(corePayload), &user))

	assert.NotNil(t, user.LoginMethods[0].ThirdParty)
	assert.Equal(t, "google", user.LoginMethods[0].ThirdParty.ID)
	assert.Equal(t, "google-user-id", user.LoginMethods[0].ThirdParty.UserID)

	data, err := json.Marshal(user)
	assert.NoError(t, err)

	var got map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &got))
	lm := got["loginMethods"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, map[string]interface{}{"id": "google", "userId": "google-user-id"}, lm["thirdParty"])
	// webauthn was absent in the payload, so it must stay omitted.
	_, hasWebauthn := lm["webauthn"]
	assert.False(t, hasWebauthn)
}
