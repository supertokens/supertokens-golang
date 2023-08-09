package multitenancy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestCreate(t *testing.T) {
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

	resp, err := CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	resp, err = CreateOrUpdateTenant("t2", multitenancymodels.TenantConfig{
		PasswordlessEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	resp, err = CreateOrUpdateTenant("t3", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	tenants, err := ListAllTenants()
	assert.Nil(t, err)
	assert.Equal(t, 4, len(tenants.OK.Tenants)) // public + 3 tenants
}

func TestGet(t *testing.T) {
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

	resp, err := CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	resp, err = CreateOrUpdateTenant("t2", multitenancymodels.TenantConfig{
		PasswordlessEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	resp, err = CreateOrUpdateTenant("t3", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	tenant1, err := GetTenant("t1")
	assert.Nil(t, err)
	assert.Equal(t, "t1", tenant1.TenantId)
	assert.True(t, tenant1.EmailPassword.Enabled)
	assert.False(t, tenant1.Passwordless.Enabled)
	assert.False(t, tenant1.ThirdParty.Enabled)
	assert.Empty(t, tenant1.CoreConfig)

	tenant2, err := GetTenant("t2")
	assert.Nil(t, err)
	assert.Equal(t, "t2", tenant2.TenantId)
	assert.True(t, tenant2.Passwordless.Enabled)
	assert.False(t, tenant2.EmailPassword.Enabled)
	assert.False(t, tenant2.ThirdParty.Enabled)
	assert.Empty(t, tenant2.CoreConfig)

	tenant3, err := GetTenant("t3")
	assert.Nil(t, err)
	assert.Equal(t, "t3", tenant3.TenantId)
	assert.True(t, tenant3.ThirdParty.Enabled)
	assert.False(t, tenant3.EmailPassword.Enabled)
	assert.False(t, tenant3.Passwordless.Enabled)
	assert.Empty(t, tenant3.CoreConfig)
}

func TestUpdate(t *testing.T) {
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
	False := false

	resp, err := CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	tenant, err := GetTenant("t1")
	assert.Nil(t, err)
	assert.Equal(t, "t1", tenant.TenantId)
	assert.True(t, tenant.EmailPassword.Enabled)
	assert.False(t, tenant.Passwordless.Enabled)

	resp, err = CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		PasswordlessEnabled: &True,
	})
	assert.Nil(t, err)
	assert.False(t, resp.OK.CreatedNew)

	tenant, err = GetTenant("t1")
	assert.Nil(t, err)
	assert.Equal(t, "t1", tenant.TenantId)
	assert.True(t, tenant.EmailPassword.Enabled)
	assert.True(t, tenant.Passwordless.Enabled)

	resp, err = CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &False,
	})
	assert.Nil(t, err)
	assert.False(t, resp.OK.CreatedNew)

	tenant, err = GetTenant("t1")
	assert.Nil(t, err)
	assert.Equal(t, "t1", tenant.TenantId)
	assert.False(t, tenant.EmailPassword.Enabled)
	assert.True(t, tenant.Passwordless.Enabled)
}

func TestDelete(t *testing.T) {
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

	resp, err := CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	resp, err = CreateOrUpdateTenant("t2", multitenancymodels.TenantConfig{
		PasswordlessEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	resp, err = CreateOrUpdateTenant("t3", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	tenants, err := ListAllTenants()
	assert.Nil(t, err)
	assert.Equal(t, 4, len(tenants.OK.Tenants)) // public + 3 tenants

	deleteRes, err := DeleteTenant("t3")
	assert.Nil(t, err)
	assert.True(t, deleteRes.OK.DidExist)

	tenants, err = ListAllTenants()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(tenants.OK.Tenants)) // public + 2 tenants
}

func TestCreateThirdPartyConfig(t *testing.T) {
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

	resp, err := CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	res, err := CreateOrUpdateThirdPartyConfig("t1", tpmodels.ProviderConfig{
		ThirdPartyId: "google",
		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "abcd",
			},
		},
	}, nil)
	assert.Nil(t, err)
	assert.True(t, res.OK.CreatedNew)

	tenant, err := GetTenant("t1")
	assert.Nil(t, err)
	assert.Equal(t, "t1", tenant.TenantId)
	assert.True(t, tenant.ThirdParty.Enabled)
	assert.Equal(t, 1, len(tenant.ThirdParty.Providers))
	assert.Equal(t, "google", tenant.ThirdParty.Providers[0].ThirdPartyId)
	assert.Equal(t, 1, len(tenant.ThirdParty.Providers[0].Clients))
	assert.Equal(t, "abcd", tenant.ThirdParty.Providers[0].Clients[0].ClientID)
}

func TestDeleteThirdPartyConfig(t *testing.T) {
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

	resp, err := CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	res, err := CreateOrUpdateThirdPartyConfig("t1", tpmodels.ProviderConfig{
		ThirdPartyId: "google",
		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "abcd",
			},
		},
	}, nil)
	assert.Nil(t, err)
	assert.True(t, res.OK.CreatedNew)

	tenant, err := GetTenant("t1")
	assert.Nil(t, err)
	assert.Equal(t, "t1", tenant.TenantId)
	assert.True(t, tenant.ThirdParty.Enabled)
	assert.Equal(t, 1, len(tenant.ThirdParty.Providers))

	delRes, err := DeleteThirdPartyConfig("t1", "google")
	assert.Nil(t, err)
	assert.True(t, delRes.OK.DidConfigExist)

	tenant, err = GetTenant("t1")
	assert.Nil(t, err)
	assert.Equal(t, "t1", tenant.TenantId)
	assert.True(t, tenant.ThirdParty.Enabled)
	assert.Equal(t, 0, len(tenant.ThirdParty.Providers))
}

func TestUpdateThirdPartyConfig(t *testing.T) {
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

	resp, err := CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	res, err := CreateOrUpdateThirdPartyConfig("t1", tpmodels.ProviderConfig{
		ThirdPartyId: "google",
		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "abcd",
			},
		},
	}, nil)
	assert.Nil(t, err)
	assert.True(t, res.OK.CreatedNew)

	tenant, err := GetTenant("t1")
	assert.Nil(t, err)
	assert.Equal(t, "t1", tenant.TenantId)
	assert.True(t, tenant.ThirdParty.Enabled)
	assert.Equal(t, 1, len(tenant.ThirdParty.Providers))

	res, err = CreateOrUpdateThirdPartyConfig("t1", tpmodels.ProviderConfig{
		ThirdPartyId: "google",
		Name:         "Custom name",
		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "efgh",
			},
		},
	}, nil)
	assert.Nil(t, err)
	assert.False(t, res.OK.CreatedNew)

	tenant, err = GetTenant("t1")
	assert.Nil(t, err)
	assert.Equal(t, "t1", tenant.TenantId)
	assert.True(t, tenant.ThirdParty.Enabled)
	assert.Equal(t, 1, len(tenant.ThirdParty.Providers))
	assert.Equal(t, "Custom name", tenant.ThirdParty.Providers[0].Name)
	assert.Equal(t, "efgh", tenant.ThirdParty.Providers[0].Clients[0].ClientID)
}
