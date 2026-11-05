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

package userroles

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/userroles/userrolesclaims"
	"github.com/supertokens/supertokens-golang/recipe/userroles/userrolesmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestValidateClaimsForSessionHandleFetchesRolesForTheSessionTenant(t *testing.T) {
	BeforeEach()
	connectionURI := unittesting.StartUpSTWithMultitenancy("localhost", "8080")
	defer AfterEach()
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
			emailpassword.Init(nil),
			session.Init(nil),
			Init(&userrolesmodels.TypeInput{
				SkipAddingRolesToAccessToken: true,
			}),
		},
	}
	err := supertokens.Init(configValue)
	assert.NoError(t, err)

	True := true
	tenantResp, err := multitenancy.CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &True,
	})
	assert.NoError(t, err)
	assert.True(t, tenantResp.OK.CreatedNew)

	_, err = CreateNewRoleOrAddPermissions("admin", []string{})
	assert.NoError(t, err)

	validateAdmin := func(sessionHandle string) sessmodels.ValidateClaimsResponse {
		res, err := session.ValidateClaimsForSessionHandle(sessionHandle, func(globalClaimValidators []claims.SessionClaimValidator, sessionInfo sessmodels.SessionInformation, userContext supertokens.UserContext) []claims.SessionClaimValidator {
			return []claims.SessionClaimValidator{
				userrolesclaims.UserRoleClaimValidators.Includes("admin", nil, nil),
			}
		})
		assert.NoError(t, err)
		assert.NotNil(t, res.OK)
		return res
	}

	// Role in t1 only, session in t1: valid
	user1, err := emailpassword.SignUp("public", "user1@example.com", "password123")
	assert.NoError(t, err)
	userId1 := user1.OK.User.ID
	_, err = multitenancy.AssociateUserToTenant("t1", userId1)
	assert.NoError(t, err)
	_, err = AddRoleToUser("t1", userId1, "admin")
	assert.NoError(t, err)

	session1, err := session.CreateNewSessionWithoutRequestResponse("t1", userId1, map[string]interface{}{}, map[string]interface{}{}, nil)
	assert.NoError(t, err)

	res := validateAdmin(session1.GetHandle())
	assert.Equal(t, 0, len(res.OK.InvalidClaims))

	sessionInfo, err := session.GetSessionInformation(session1.GetHandle())
	assert.NoError(t, err)
	assert.Equal(t, "t1", sessionInfo.TenantId)
	_, hasTenantId := sessionInfo.CustomClaimsInAccessTokenPayload["tId"]
	assert.False(t, hasTenantId)
	assert.NotNil(t, sessionInfo.CustomClaimsInAccessTokenPayload["st-role"])

	// Role in public only, session in t1: invalid
	user2, err := emailpassword.SignUp("public", "user2@example.com", "password123")
	assert.NoError(t, err)
	userId2 := user2.OK.User.ID
	_, err = multitenancy.AssociateUserToTenant("t1", userId2)
	assert.NoError(t, err)
	_, err = AddRoleToUser("public", userId2, "admin")
	assert.NoError(t, err)

	session2, err := session.CreateNewSessionWithoutRequestResponse("t1", userId2, map[string]interface{}{}, map[string]interface{}{}, nil)
	assert.NoError(t, err)

	res = validateAdmin(session2.GetHandle())
	assert.Equal(t, 1, len(res.OK.InvalidClaims))
	if len(res.OK.InvalidClaims) == 1 {
		assert.Equal(t, "st-role", res.OK.InvalidClaims[0].ID)
	}
}

func TestValidateClaimsForSessionHandleDoesNotMergeProtectedProps(t *testing.T) {
	BeforeEach()
	connectionURI := unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
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
			session.Init(&sessmodels.TypeInput{
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						// A claim (or an override) returns an update that includes protected props
						*originalImplementation.ValidateClaims = func(userId string, accessTokenPayload map[string]interface{}, claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) (sessmodels.ValidateClaimsResult, error) {
							return sessmodels.ValidateClaimsResult{
								InvalidClaims: []claims.ClaimValidationError{},
								AccessTokenPayloadUpdate: map[string]interface{}{
									"custom":        "value",
									"sub":           "someone-else",
									"sessionHandle": "another-handle",
									"exp":           1,
									"stt":           1,
									"tId":           "another-tenant",
								},
							}, nil
						}
						return originalImplementation
					},
				},
			}),
		},
	}
	err := supertokens.Init(configValue)
	assert.NoError(t, err)

	sess, err := session.CreateNewSessionWithoutRequestResponse("public", "user-id", map[string]interface{}{}, map[string]interface{}{}, nil)
	assert.NoError(t, err)

	res, err := session.ValidateClaimsForSessionHandle(sess.GetHandle(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, res.OK)

	sessionInfo, err := session.GetSessionInformation(sess.GetHandle())
	assert.NoError(t, err)
	assert.Equal(t, "value", sessionInfo.CustomClaimsInAccessTokenPayload["custom"])
	for _, protectedProp := range []string{"sub", "sessionHandle", "exp", "stt", "tId"} {
		_, found := sessionInfo.CustomClaimsInAccessTokenPayload[protectedProp]
		assert.False(t, found, protectedProp)
	}
}
