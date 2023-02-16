/* Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestCreateAssignAndDeleteRole(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "Supertokens Demo",
			APIDomain:     "https://api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	})

	if !canRunTest(t) {
		return
	}

	roles := []string{"role1", "role2", "role3"}

	for _, role := range roles {
		createResult, err := CreateNewRoleOrAddPermissions(role, []string{}, nil, &map[string]interface{}{})
		assert.NoError(t, err)
		assert.NotNil(t, createResult.OK)
		assert.True(t, createResult.OK.CreatedNewRole)

		addResult, err := AddRoleToUser("userId", role, nil, &map[string]interface{}{})
		assert.NoError(t, err)
		assert.NotNil(t, addResult.OK)
		assert.False(t, addResult.OK.DidUserAlreadyHaveRole)
	}

	listResult, err := GetAllRoles(nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, listResult.OK)
	assert.Contains(t, listResult.OK.Roles, "role1")
	assert.Contains(t, listResult.OK.Roles, "role2")
	assert.Contains(t, listResult.OK.Roles, "role3")
	assert.Equal(t, 3, len(listResult.OK.Roles))

	// Delete a role
	deleteResult, err := DeleteRole("role2", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, deleteResult.OK)
	assert.True(t, deleteResult.OK.DidRoleExist)

	listResult, err = GetAllRoles(nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, listResult.OK)
	assert.Contains(t, listResult.OK.Roles, "role1")
	assert.Contains(t, listResult.OK.Roles, "role3")
	assert.Equal(t, 2, len(listResult.OK.Roles))

	userRolesResult, err := GetRolesForUser("userId", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, userRolesResult.OK)
	assert.Contains(t, userRolesResult.OK.Roles, "role1")
	assert.Contains(t, userRolesResult.OK.Roles, "role3")
	assert.NotContains(t, userRolesResult.OK.Roles, "role2")
	assert.Equal(t, 2, len(userRolesResult.OK.Roles))
}

func TestDeleteUnknownRole(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "Supertokens Demo",
			APIDomain:     "https://api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	})

	if !canRunTest(t) {
		return
	}

	// Delete a role
	deleteResult, err := DeleteRole("role1", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, deleteResult.OK)
	assert.False(t, deleteResult.OK.DidRoleExist)
}
