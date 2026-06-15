/*
 * Copyright (c) 2025, VRAI Labs and/or its affiliates. All rights reserved.
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

// These tests exercise the getSession closure in recipeImplementation.go
// directly through RecipeImpl.GetSession. They guard against a regression
// where the "session is optional" branches dereferenced a nil
// options.SessionRequired (*bool) and panicked on a present-but-invalid
// access token. The canonical semantics (matching supertokens-node and
// supertokens-python) are: a nil SessionRequired means REQUIRED, so the
// recipe must return an UnauthorizedError rather than panic, and only an
// explicit SessionRequired=false returns a nil session with no error.

func initSessionRecipeForGetSessionTests(t *testing.T) {
	connectionURI := unittesting.StartUpST("localhost", "8080")
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
			Init(nil),
		},
	}
	require.NoError(t, supertokens.Init(configValue))
}

func getSessionRecipeImpl(t *testing.T) sessmodels.RecipeInterface {
	instance, err := getRecipeInstanceOrThrowError()
	require.NoError(t, err)
	return instance.RecipeImpl
}

// nil SessionRequired + nil access token => UnauthorizedError (not panic).
func TestGetSessionNilSessionRequiredNoTokenIsUnauthorized(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initSessionRecipeForGetSessionTests(t)

	impl := getSessionRecipeImpl(t)

	options := &sessmodels.VerifySessionOptions{} // SessionRequired left nil
	sess, err := (*impl.GetSession)(nil, nil, options, &map[string]interface{}{})

	assert.Nil(t, sess)
	require.Error(t, err)
	assert.IsType(t, errors.UnauthorizedError{}, err)
}

// nil SessionRequired + present-but-unparseable token => UnauthorizedError
// (not panic). This is the exact path the webauthn recipe hits via
// getSessionWithoutClaimValidation.
func TestGetSessionNilSessionRequiredUnparseableTokenIsUnauthorized(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initSessionRecipeForGetSessionTests(t)

	impl := getSessionRecipeImpl(t)

	badToken := "this-is-not-a-jwt"
	options := &sessmodels.VerifySessionOptions{} // SessionRequired left nil
	sess, err := (*impl.GetSession)(&badToken, nil, options, &map[string]interface{}{})

	assert.Nil(t, sess)
	require.Error(t, err)
	assert.IsType(t, errors.UnauthorizedError{}, err)
	assert.Equal(t, "Token parsing failed", err.(errors.UnauthorizedError).Msg)
}

// nil SessionRequired + token that parses as a JWT but fails access-token
// structure validation => UnauthorizedError (exercises the second buggy
// branch, after ValidateAccessTokenStructure).
func TestGetSessionNilSessionRequiredInvalidStructureTokenIsUnauthorized(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initSessionRecipeForGetSessionTests(t)

	impl := getSessionRecipeImpl(t)

	// A syntactically valid JWS (header.payload.signature, all base64url) whose
	// payload is an empty JSON object, so it parses but has none of the fields
	// ValidateAccessTokenStructure requires.
	// header {"alg":"RS256","typ":"JWT","version":"3","kid":"x"} / payload {}
	structurallyInvalidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsInZlcnNpb24iOiIzIiwia2lkIjoieCJ9.e30.c2ln"
	options := &sessmodels.VerifySessionOptions{} // SessionRequired left nil
	sess, err := (*impl.GetSession)(&structurallyInvalidToken, nil, options, &map[string]interface{}{})

	assert.Nil(t, sess)
	require.Error(t, err)
	assert.IsType(t, errors.UnauthorizedError{}, err)
}

// explicit SessionRequired=false + nil token => nil session, no error.
func TestGetSessionSessionRequiredFalseNoTokenReturnsNilSession(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initSessionRecipeForGetSessionTests(t)

	impl := getSessionRecipeImpl(t)

	sessionRequired := false
	options := &sessmodels.VerifySessionOptions{SessionRequired: &sessionRequired}
	sess, err := (*impl.GetSession)(nil, nil, options, &map[string]interface{}{})

	require.NoError(t, err)
	assert.Nil(t, sess)
}

// explicit SessionRequired=false + present-but-invalid token => nil session,
// no error.
func TestGetSessionSessionRequiredFalseInvalidTokenReturnsNilSession(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initSessionRecipeForGetSessionTests(t)

	impl := getSessionRecipeImpl(t)

	badToken := "this-is-not-a-jwt"
	sessionRequired := false
	options := &sessmodels.VerifySessionOptions{SessionRequired: &sessionRequired}
	sess, err := (*impl.GetSession)(&badToken, nil, options, &map[string]interface{}{})

	require.NoError(t, err)
	assert.Nil(t, sess)
}

// valid token => a real session is returned.
func TestGetSessionValidTokenReturnsSession(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initSessionRecipeForGetSessionTests(t)

	testServer := GetTestServer(t)
	defer testServer.Close()

	created, err := CreateNewSessionWithoutRequestResponse("public", "testuser", map[string]interface{}{}, map[string]interface{}{}, nil)
	require.NoError(t, err)
	accessToken := created.GetAccessToken()

	impl := getSessionRecipeImpl(t)

	// A non-nil pointer where the code expects one: SessionRequired=true.
	sessionRequired := true
	options := &sessmodels.VerifySessionOptions{SessionRequired: &sessionRequired}
	sess, err := (*impl.GetSession)(&accessToken, nil, options, &map[string]interface{}{})

	require.NoError(t, err)
	require.NotNil(t, sess)
	assert.Equal(t, "testuser", sess.GetUserID())
}
