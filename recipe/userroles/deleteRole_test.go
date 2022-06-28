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
		createResult, err := CreateNewRoleOrAddPermissions(role, []string{}, &map[string]interface{}{})
		assert.NoError(t, err)
		assert.NotNil(t, createResult.OK)
		assert.True(t, createResult.OK.CreatedNewRole)

		addResult, err := AddRoleToUser("userId", role, &map[string]interface{}{})
		assert.NoError(t, err)
		assert.NotNil(t, addResult.OK)
		assert.False(t, addResult.OK.DidUserAlreadyHaveRole)
	}

	listResult, err := GetAllRoles(&map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, listResult.OK)
	assert.Contains(t, listResult.OK.Roles, "role1")
	assert.Contains(t, listResult.OK.Roles, "role2")
	assert.Contains(t, listResult.OK.Roles, "role3")
	assert.Equal(t, 3, len(listResult.OK.Roles))

	// Delete a role
	deleteResult, err := DeleteRole("role2", &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, deleteResult.OK)
	assert.True(t, deleteResult.OK.DidRoleExist)

	listResult, err = GetAllRoles(&map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, listResult.OK)
	assert.Contains(t, listResult.OK.Roles, "role1")
	assert.Contains(t, listResult.OK.Roles, "role3")
	assert.Equal(t, 2, len(listResult.OK.Roles))

	userRolesResult, err := GetRolesForUser("userId", &map[string]interface{}{})
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
	deleteResult, err := DeleteRole("role1", &map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, deleteResult.OK)
	assert.False(t, deleteResult.OK.DidRoleExist)
}
