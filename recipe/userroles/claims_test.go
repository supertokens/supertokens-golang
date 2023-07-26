package userroles

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	sessErrors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/userroles/userrolesclaims"
	"github.com/supertokens/supertokens-golang/recipe/userroles/userrolesmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestShouldAddClaimsToSessionWithoutConfig(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	if !canRunTest(t) {
		return
	}

	res := fakeRes{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	sessionContainer, err := session.CreateNewSession(req, res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	userroleClaimValue, err := session.GetClaimValue(sessionContainer.GetHandle(), userrolesclaims.UserRoleClaim)
	assert.NoError(t, err)

	assert.NotNil(t, userroleClaimValue.OK)
	assert.Equal(t, []interface{}{}, userroleClaimValue.OK.Value)

	permissionClaimValue, err := session.GetClaimValue(sessionContainer.GetHandle(), userrolesclaims.PermissionClaim)
	assert.NoError(t, err)

	assert.NotNil(t, permissionClaimValue.OK)
	assert.Equal(t, []interface{}{}, permissionClaimValue.OK.Value)
}

func TestShouldNotAddClaimsToSessionIfDisabledInConfig(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(&userrolesmodels.TypeInput{
				SkipAddingRolesToAccessToken:       true,
				SkipAddingPermissionsToAccessToken: true,
			}),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	if !canRunTest(t) {
		return
	}

	res := fakeRes{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	sessionContainer, err := session.CreateNewSession(req, res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	userroleClaimValue, err := session.GetClaimValue(sessionContainer.GetHandle(), userrolesclaims.UserRoleClaim)
	assert.NoError(t, err)

	assert.NotNil(t, userroleClaimValue.OK)
	assert.Equal(t, nil, userroleClaimValue.OK.Value)

	permissionClaimValue, err := session.GetClaimValue(sessionContainer.GetHandle(), userrolesclaims.PermissionClaim)
	assert.NoError(t, err)

	assert.NotNil(t, permissionClaimValue.OK)
	assert.Equal(t, nil, permissionClaimValue.OK.Value)
}

func TestShouldAddClaimsToSessionWithValues(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	if !canRunTest(t) {
		return
	}

	CreateNewRoleOrAddPermissions("test", []string{"a", "b"}, &map[string]interface{}{})
	AddRoleToUser("public", "userId", "test", &map[string]interface{}{})

	res := fakeRes{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	sessionContainer, err := session.CreateNewSession(req, res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	userroleClaimValue, err := session.GetClaimValue(sessionContainer.GetHandle(), userrolesclaims.UserRoleClaim)
	assert.NoError(t, err)

	assert.NotNil(t, userroleClaimValue.OK)
	assert.Equal(t, []interface{}{"test"}, userroleClaimValue.OK.Value)

	permissionClaimValue, err := session.GetClaimValue(sessionContainer.GetHandle(), userrolesclaims.PermissionClaim)
	assert.NoError(t, err)

	assert.NotNil(t, permissionClaimValue.OK)
	assert.Equal(t, 2, len(permissionClaimValue.OK.Value.([]interface{})))
	assert.Contains(t, permissionClaimValue.OK.Value, "a")
	assert.Contains(t, permissionClaimValue.OK.Value, "b")
}

func TestShouldValidateRoles(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	if !canRunTest(t) {
		return
	}

	CreateNewRoleOrAddPermissions("test", []string{"a", "b"}, &map[string]interface{}{})
	AddRoleToUser("public", "userId", "test", &map[string]interface{}{})

	res := fakeRes{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	sessionContainer, err := session.CreateNewSession(req, res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{
		userrolesclaims.UserRoleClaimValidators.Includes("test", nil, nil),
	})
	assert.NoError(t, err)

	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{
		userrolesclaims.UserRoleClaimValidators.Includes("nope", nil, nil),
	})
	assert.NotNil(t, err)
	assert.IsType(t, sessErrors.InvalidClaimError{}, err)

	invalidClaimErr := err.(sessErrors.InvalidClaimError)
	assert.Equal(t, 1, len(invalidClaimErr.InvalidClaims))
	assert.Equal(t, "st-role", invalidClaimErr.InvalidClaims[0].ID)
	assert.Equal(t, map[string]interface{}{
		"actualValue":       []interface{}{"test"},
		"expectedToInclude": "nope",
		"message":           "wrong value",
	}, invalidClaimErr.InvalidClaims[0].Reason)
}

func TestShouldValidateRolesAfterRefetching(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(&userrolesmodels.TypeInput{
				SkipAddingRolesToAccessToken: true,
			}),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	if !canRunTest(t) {
		return
	}

	CreateNewRoleOrAddPermissions("test", []string{"a", "b"}, &map[string]interface{}{})
	AddRoleToUser("public", "userId", "test", &map[string]interface{}{})

	res := fakeRes{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	sessionContainer, err := session.CreateNewSession(req, res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{
		userrolesclaims.UserRoleClaimValidators.Includes("test", nil, nil),
	})
	assert.NoError(t, err)

	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{
		userrolesclaims.UserRoleClaimValidators.Includes("nope", nil, nil),
	})
	assert.NotNil(t, err)
	assert.IsType(t, sessErrors.InvalidClaimError{}, err)

	invalidClaimErr := err.(sessErrors.InvalidClaimError)
	assert.Equal(t, 1, len(invalidClaimErr.InvalidClaims))
	assert.Equal(t, "st-role", invalidClaimErr.InvalidClaims[0].ID)
	assert.Equal(t, map[string]interface{}{
		"actualValue":       []interface{}{"test"},
		"expectedToInclude": "nope",
		"message":           "wrong value",
	}, invalidClaimErr.InvalidClaims[0].Reason)
}

func TestShouldValidatePermissions(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	if !canRunTest(t) {
		return
	}

	CreateNewRoleOrAddPermissions("test", []string{"a", "b"}, &map[string]interface{}{})
	AddRoleToUser("public", "userId", "test", &map[string]interface{}{})

	res := fakeRes{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	sessionContainer, err := session.CreateNewSession(req, res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{
		userrolesclaims.PermissionClaimValidators.Includes("a", nil, nil),
	})
	assert.NoError(t, err)

	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{
		userrolesclaims.PermissionClaimValidators.Includes("nope", nil, nil),
	})
	assert.NotNil(t, err)
	assert.IsType(t, sessErrors.InvalidClaimError{}, err)

	invalidClaimErr := err.(sessErrors.InvalidClaimError)
	assert.Equal(t, 1, len(invalidClaimErr.InvalidClaims))
	assert.Equal(t, "st-perm", invalidClaimErr.InvalidClaims[0].ID)
	reason := invalidClaimErr.InvalidClaims[0].Reason.(map[string]interface{})
	assert.Equal(t, "wrong value", reason["message"])
	assert.Equal(t, "nope", reason["expectedToInclude"])
	assert.Equal(t, 2, len(reason["actualValue"].([]interface{})))
	assert.Contains(t, reason["actualValue"], "a")
	assert.Contains(t, reason["actualValue"], "b")
}

func TestShouldValidatePermissionsAfterRefetching(t *testing.T) {
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(&userrolesmodels.TypeInput{
				SkipAddingPermissionsToAccessToken: true,
			}),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	if !canRunTest(t) {
		return
	}

	CreateNewRoleOrAddPermissions("test", []string{"a", "b"}, &map[string]interface{}{})
	AddRoleToUser("public", "userId", "test", &map[string]interface{}{})

	res := fakeRes{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	sessionContainer, err := session.CreateNewSession(req, res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{
		userrolesclaims.PermissionClaimValidators.Includes("a", nil, nil),
	})
	assert.NoError(t, err)

	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{
		userrolesclaims.PermissionClaimValidators.Includes("nope", nil, nil),
	})
	assert.NotNil(t, err)
	assert.IsType(t, sessErrors.InvalidClaimError{}, err)

	invalidClaimErr := err.(sessErrors.InvalidClaimError)
	assert.Equal(t, 1, len(invalidClaimErr.InvalidClaims))
	assert.Equal(t, "st-perm", invalidClaimErr.InvalidClaims[0].ID)
	reason := invalidClaimErr.InvalidClaims[0].Reason.(map[string]interface{})
	assert.Equal(t, "wrong value", reason["message"])
	assert.Equal(t, "nope", reason["expectedToInclude"])
	assert.Equal(t, 2, len(reason["actualValue"].([]interface{})))
	assert.Contains(t, reason["actualValue"], "a")
	assert.Contains(t, reason["actualValue"], "b")
}
