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

func TestAddNewRoleToUser(t *testing.T) {
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

	// Create a new role
	createResult, err := CreateNewRoleOrAddPermissions("role", []string{}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.True(t, createResult.OK.CreatedNewRole)

	// Add role to the user
	addResult, err := AddRoleToUser("userId", "role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, addResult.OK)
	assert.False(t, addResult.OK.DidUserAlreadyHaveRole)

	// Check user has new role
	listResult, err := GetRolesForUser("userId", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, listResult.OK)
	assert.Contains(t, listResult.OK.Roles, "role")
}

func TestAddDuplicateRoleToUser(t *testing.T) {
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

	// Create a new role
	createResult, err := CreateNewRoleOrAddPermissions("role", []string{}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.True(t, createResult.OK.CreatedNewRole)

	// Add role to the user
	addResult, err := AddRoleToUser("userId", "role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, addResult.OK)
	assert.False(t, addResult.OK.DidUserAlreadyHaveRole)

	// Add role to the user
	addResult, err = AddRoleToUser("userId", "role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, addResult.OK)
	assert.True(t, addResult.OK.DidUserAlreadyHaveRole)

	// Check user has new role
	listResult, err := GetRolesForUser("userId", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, listResult.OK)
	assert.Contains(t, listResult.OK.Roles, "role")
}

func TestAddUnknownRoleToUser(t *testing.T) {
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

	// Add role to the user
	addResult, err := AddRoleToUser("userId", "role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Nil(t, addResult.OK)
	assert.NotNil(t, addResult.UnknownRoleError)

	// Check user has new role
	listResult, err := GetRolesForUser("userId", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, listResult.OK)
	assert.NotContains(t, listResult.OK.Roles, "role")
}

func TestGetUsersThatHaveARole(t *testing.T) {
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

	// Create a new role
	createResult, err := CreateNewRoleOrAddPermissions("role", []string{}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.True(t, createResult.OK.CreatedNewRole)

	// Add role to the users
	users := []string{"user1", "user2", "user3"}
	for _, user := range users {
		addResult, err := AddRoleToUser(user, "role", nil, &map[string]interface{}{})
		assert.NoError(t, err)
		assert.NotNil(t, addResult.OK)
		assert.False(t, addResult.OK.DidUserAlreadyHaveRole)
	}

	// Check user has new role
	listResult, err := GetUsersThatHaveRole("role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, listResult.OK)
	assert.Contains(t, listResult.OK.Users, "user1")
	assert.Contains(t, listResult.OK.Users, "user2")
	assert.Contains(t, listResult.OK.Users, "user3")
	assert.Equal(t, 3, len(listResult.OK.Users))
}

func TestGetUsersThatHaveUnknownRole(t *testing.T) {
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

	// Check user has new role
	listResult, err := GetUsersThatHaveRole("role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Nil(t, listResult.OK)
	assert.NotNil(t, listResult.UnknownRoleError)
}

func TestRemoveUserRole(t *testing.T) {
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

	// Create a new roles and assign them to the users
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

	rolesResult, err := GetRolesForUser("userId", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, rolesResult.OK)
	assert.Contains(t, rolesResult.OK.Roles, "role1")
	assert.Contains(t, rolesResult.OK.Roles, "role2")
	assert.Contains(t, rolesResult.OK.Roles, "role3")
	assert.Equal(t, 3, len(rolesResult.OK.Roles))

	// Remove role from the user
	removeResult, err := RemoveUserRole("userId", "role2", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, removeResult.OK)
	assert.True(t, removeResult.OK.DidUserHaveRole)

	rolesResult, err = GetRolesForUser("userId", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, rolesResult.OK)
	assert.Contains(t, rolesResult.OK.Roles, "role1")
	assert.Contains(t, rolesResult.OK.Roles, "role3")
	assert.Equal(t, 2, len(rolesResult.OK.Roles))
}

func TestRemoveUnassignedUserRole(t *testing.T) {
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

	// Create a new roles and assign them to the users
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
	createResult, err := CreateNewRoleOrAddPermissions("role4", []string{}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.True(t, createResult.OK.CreatedNewRole)

	rolesResult, err := GetRolesForUser("userId", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, rolesResult.OK)
	assert.Contains(t, rolesResult.OK.Roles, "role1")
	assert.Contains(t, rolesResult.OK.Roles, "role2")
	assert.Contains(t, rolesResult.OK.Roles, "role3")
	assert.Equal(t, 3, len(rolesResult.OK.Roles))

	// Remove role from the user
	removeResult, err := RemoveUserRole("userId", "role4", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, removeResult.OK)
	assert.False(t, removeResult.OK.DidUserHaveRole)

	rolesResult, err = GetRolesForUser("userId", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, rolesResult.OK)
	assert.Contains(t, rolesResult.OK.Roles, "role1")
	assert.Contains(t, rolesResult.OK.Roles, "role2")
	assert.Contains(t, rolesResult.OK.Roles, "role3")
	assert.Equal(t, 3, len(rolesResult.OK.Roles))
}

func TestRemoveUnknownUserRole(t *testing.T) {
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

	// Remove role from the user
	removeResult, err := RemoveUserRole("userId", "role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Nil(t, removeResult.OK)
	assert.NotNil(t, removeResult.UnknownRoleError)
}
