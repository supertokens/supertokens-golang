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

package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func validV3AccessTokenPayloadForTest() map[string]interface{} {
	now := time.Now().Unix()
	return map[string]interface{}{
		"sub":               "userId",
		"sessionHandle":     "sessionHandle",
		"refreshTokenHash1": "refreshTokenHash1",
		"exp":               float64(now + 3600),
		"iat":               float64(now),
	}
}

func TestValidateAccessTokenStructureChecksTokenType(t *testing.T) {
	payload := validV3AccessTokenPayloadForTest()
	assert.NoError(t, ValidateAccessTokenStructure(payload, 3))

	payload["stt"] = float64(0)
	assert.NoError(t, ValidateAccessTokenStructure(payload, 3))
	payload["tId"] = "public"
	assert.NoError(t, ValidateAccessTokenStructure(payload, 4))

	payload["stt"] = float64(1)
	assert.Error(t, ValidateAccessTokenStructure(payload, 3))

	payload["stt"] = "0"
	assert.Error(t, ValidateAccessTokenStructure(payload, 3))

	payload["stt"] = nil
	assert.Error(t, ValidateAccessTokenStructure(payload, 3))
}

func createSignedTokenForTest(t *testing.T, payload map[string]interface{}) string {
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	assert.NoError(t, err)
	response, err := querier.SendPostRequest("/recipe/jwt", map[string]interface{}{
		"payload":             payload,
		"validity":            3600,
		"algorithm":           "RS256",
		"jwksDomain":          "http://api.supertokens.io",
		"useStaticSigningKey": true,
	}, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, "OK", response["status"])
	return response["jwt"].(string)
}

func TestGetSessionChecksAccessTokenType(t *testing.T) {
	BeforeEach()
	connectionURI := unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	useDynamicKey := false
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				UseDynamicAccessTokenSigningKey: &useDynamicKey,
			}),
		},
	}
	err := supertokens.Init(configValue)
	assert.NoError(t, err)

	falseValue := false
	options := &sessmodels.VerifySessionOptions{
		AntiCsrfCheck: &falseValue,
	}

	// No token type: accepted
	token := createSignedTokenForTest(t, validV3AccessTokenPayloadForTest())
	sess, err := GetSessionWithoutRequestResponse(token, nil, options)
	assert.NoError(t, err)
	assert.NotNil(t, sess)
	if sess != nil {
		assert.Equal(t, "userId", sess.GetUserID())
	}

	// Access token type: accepted
	payload := validV3AccessTokenPayloadForTest()
	payload["stt"] = 0
	token = createSignedTokenForTest(t, payload)
	sess, err = GetSessionWithoutRequestResponse(token, nil, options)
	assert.NoError(t, err)
	assert.NotNil(t, sess)

	// Any other token type: rejected
	payload = validV3AccessTokenPayloadForTest()
	payload["stt"] = 1
	token = createSignedTokenForTest(t, payload)
	sess, err = GetSessionWithoutRequestResponse(token, nil, options)
	assert.Nil(t, sess)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &errors.UnauthorizedError{})
}

func TestTokenTypeIsAProtectedProp(t *testing.T) {
	assert.Contains(t, protectedProps, "stt")
}
