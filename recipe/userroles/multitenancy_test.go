package userroles

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDifferentRolesCanBeAssignedToSameUserAcrossTenants(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			session.Init(nil),
			Init(nil),
		},
	}

	BeforeEach()
	unittesting.StartUpSTWithMultitenancy("localhost", "8080")
	defer AfterEach()

	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	True := true
	resp, err := multitenancy.CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)
	resp, err = multitenancy.CreateOrUpdateTenant("t2", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)
	resp, err = multitenancy.CreateOrUpdateTenant("t3", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	user, err := emailpassword.SignUp("public", "test@example.com", "password")
	assert.Nil(t, err)

	multitenancy.AssociateUserToTenant("t1", user.OK.User.ID)
	multitenancy.AssociateUserToTenant("t2", user.OK.User.ID)
	multitenancy.AssociateUserToTenant("t3", user.OK.User.ID)

	CreateNewRoleOrAddPermissions("role1", []string{})
	CreateNewRoleOrAddPermissions("role2", []string{})
	CreateNewRoleOrAddPermissions("role3", []string{})

	AddRoleToUser("t1", user.OK.User.ID, "role1")
	AddRoleToUser("t1", user.OK.User.ID, "role2")

	AddRoleToUser("t2", user.OK.User.ID, "role2")
	AddRoleToUser("t2", user.OK.User.ID, "role3")

	AddRoleToUser("t3", user.OK.User.ID, "role1")
	AddRoleToUser("t3", user.OK.User.ID, "role3")

	roles, err := GetRolesForUser("t1", user.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(roles.OK.Roles))
	assert.Contains(t, roles.OK.Roles, "role1")
	assert.Contains(t, roles.OK.Roles, "role2")

	roles, err = GetRolesForUser("t2", user.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(roles.OK.Roles))
	assert.Contains(t, roles.OK.Roles, "role2")
	assert.Contains(t, roles.OK.Roles, "role3")

	roles, err = GetRolesForUser("t3", user.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(roles.OK.Roles))
	assert.Contains(t, roles.OK.Roles, "role1")
	assert.Contains(t, roles.OK.Roles, "role3")
}
