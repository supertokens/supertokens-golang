package userroles

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/userroles/userrolesclaims"
	"github.com/supertokens/supertokens-golang/recipe/userroles/userrolesmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func setUpRolesInTwoTenants(t *testing.T) string {
	connectionURI := unittesting.StartUpSTWithMultitenancy("localhost", "8080")
	err := supertokens.Init(supertokens.TypeInput{
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			// The claims are not added on session creation, so the values below only come from FetchAndSetClaim
			Init(&userrolesmodels.TypeInput{
				SkipAddingRolesToAccessToken:       true,
				SkipAddingPermissionsToAccessToken: true,
			}),
		},
	})
	assert.NoError(t, err)

	True := true
	tenantResp, err := multitenancy.CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &True,
	})
	assert.NoError(t, err)
	assert.NotNil(t, tenantResp.OK)

	user, err := emailpassword.SignUp("public", "test@example.com", "password1234")
	assert.NoError(t, err)
	userId := user.OK.User.ID

	_, err = multitenancy.AssociateUserToTenant("t1", userId)
	assert.NoError(t, err)

	_, err = CreateNewRoleOrAddPermissions("admin", []string{})
	assert.NoError(t, err)
	_, err = CreateNewRoleOrAddPermissions("editor", []string{})
	assert.NoError(t, err)

	addResp, err := AddRoleToUser("public", userId, "admin")
	assert.NoError(t, err)
	assert.NotNil(t, addResp.OK)
	addResp, err = AddRoleToUser("t1", userId, "editor")
	assert.NoError(t, err)
	assert.NotNil(t, addResp.OK)

	return userId
}

func TestFetchAndSetClaimUsesTheTenantOfTheSession(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	userId := setUpRolesInTwoTenants(t)
	if !canRunTest(t) {
		return
	}

	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)

	t1Session, err := session.CreateNewSession(req, fakeRes{}, "t1", userId, map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)
	assert.Nil(t, t1Session.GetClaimValue(userrolesclaims.UserRoleClaim))

	err = t1Session.FetchAndSetClaim(userrolesclaims.UserRoleClaim)
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"editor"}, t1Session.GetClaimValue(userrolesclaims.UserRoleClaim))
	assert.Equal(t, "t1", t1Session.GetTenantId())

	publicSession, err := session.CreateNewSession(req, fakeRes{}, "public", userId, map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	err = publicSession.FetchAndSetClaim(userrolesclaims.UserRoleClaim)
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"admin"}, publicSession.GetClaimValue(userrolesclaims.UserRoleClaim))
}
