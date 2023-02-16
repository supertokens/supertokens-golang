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

func TestCreateRole(t *testing.T) {
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

	createResult, err := CreateNewRoleOrAddPermissions("role", []string{}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.True(t, createResult.OK.CreatedNewRole)
}

func TestCreateRoleTwice(t *testing.T) {
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

	createResult, err := CreateNewRoleOrAddPermissions("role", []string{}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.True(t, createResult.OK.CreatedNewRole)

	createResult, err = CreateNewRoleOrAddPermissions("role", []string{}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.False(t, createResult.OK.CreatedNewRole)
}

func TestCreateRoleWithPermissions(t *testing.T) {
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

	createResult, err := CreateNewRoleOrAddPermissions("role", []string{"permission1"}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.True(t, createResult.OK.CreatedNewRole)

	permissionResult, err := GetPermissionsForRole("role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, permissionResult.OK)
	assert.Contains(t, permissionResult.OK.Permissions, "permission1")
}

func TestAddPermissionToRole(t *testing.T) {
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

	createResult, err := CreateNewRoleOrAddPermissions("role", []string{"permission1"}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.True(t, createResult.OK.CreatedNewRole)

	createResult, err = CreateNewRoleOrAddPermissions("role", []string{"permission2", "permission3"}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.False(t, createResult.OK.CreatedNewRole)

	permissionResult, err := GetPermissionsForRole("role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, permissionResult.OK)
	assert.Contains(t, permissionResult.OK.Permissions, "permission1")
	assert.Contains(t, permissionResult.OK.Permissions, "permission2")
	assert.Contains(t, permissionResult.OK.Permissions, "permission3")
}

func TestDuplicatePermission(t *testing.T) {
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

	createResult, err := CreateNewRoleOrAddPermissions("role", []string{"permission1", "permission2"}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.True(t, createResult.OK.CreatedNewRole)

	createResult, err = CreateNewRoleOrAddPermissions("role", []string{"permission1", "permission2"}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.False(t, createResult.OK.CreatedNewRole)

	permissionResult, err := GetPermissionsForRole("role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, permissionResult.OK)
	assert.Contains(t, permissionResult.OK.Permissions, "permission1")
	assert.Contains(t, permissionResult.OK.Permissions, "permission2")
	assert.Equal(t, 2, len(permissionResult.OK.Permissions))
}

func TestPermissionsOfUnknownRole(t *testing.T) {
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

	permissionResult, err := GetPermissionsForRole("role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Nil(t, permissionResult.OK)
	assert.NotNil(t, permissionResult.UnknownRoleError)
}

func TestGetRolesThatHavePermission(t *testing.T) {
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
		createResult, err := CreateNewRoleOrAddPermissions(role, []string{"permission"}, nil, &map[string]interface{}{})
		assert.NoError(t, err)
		assert.NotNil(t, createResult.OK)
		assert.True(t, createResult.OK.CreatedNewRole)
	}

	listResult, err := GetRolesThatHavePermission("permission", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, listResult.OK)
	assert.Contains(t, listResult.OK.Roles, "role1")
	assert.Contains(t, listResult.OK.Roles, "role2")
	assert.Contains(t, listResult.OK.Roles, "role3")
	assert.Equal(t, 3, len(listResult.OK.Roles))
}

func TestGetRolesThatHaveUnknownPermission(t *testing.T) {
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
	}

	listResult, err := GetRolesThatHavePermission("permission", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, listResult.OK)
	assert.Equal(t, 0, len(listResult.OK.Roles))
}

func TestDeletePermissionFromRole(t *testing.T) {
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

	createResult, err := CreateNewRoleOrAddPermissions("role", []string{"permission1", "permission2", "permission3"}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, createResult.OK)
	assert.True(t, createResult.OK.CreatedNewRole)

	permissionResult, err := GetPermissionsForRole("role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, permissionResult.OK)
	assert.Contains(t, permissionResult.OK.Permissions, "permission1")
	assert.Contains(t, permissionResult.OK.Permissions, "permission2")
	assert.Contains(t, permissionResult.OK.Permissions, "permission3")

	removeResult, err := RemovePermissionsFromRole("role", []string{"permission1", "permission3"}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, removeResult.OK)

	permissionResult, err = GetPermissionsForRole("role", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, permissionResult.OK)
	assert.NotContains(t, permissionResult.OK.Permissions, "permission1")
	assert.Contains(t, permissionResult.OK.Permissions, "permission2")
	assert.NotContains(t, permissionResult.OK.Permissions, "permission3")
	assert.Equal(t, 1, len(permissionResult.OK.Permissions))
}

func TestDeletePermissionFromUnknownRole(t *testing.T) {
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

	removeResult, err := RemovePermissionsFromRole("role", []string{"permission1", "permission2"}, nil, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Nil(t, removeResult.OK)
	assert.NotNil(t, removeResult.UnknownRoleError)
}
