package thirdparty

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestRecipeFunctionsWithMultitenancy(t *testing.T) {
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
	resp, err := multitenancy.CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)
	resp, err = multitenancy.CreateOrUpdateTenant("t2", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	resp, err = multitenancy.CreateOrUpdateTenant("t3", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	user1a, err := ManuallyCreateOrUpdateUser("t1", "google", "googleid", "test@example.com")
	assert.Nil(t, err)
	user1b, err := ManuallyCreateOrUpdateUser("t1", "facebook", "fbid", "test@example.com")
	assert.Nil(t, err)

	user2a, err := ManuallyCreateOrUpdateUser("t2", "google", "googleid", "test@example.com")
	assert.Nil(t, err)
	user2b, err := ManuallyCreateOrUpdateUser("t2", "facebook", "fbid", "test@example.com")
	assert.Nil(t, err)

	user3a, err := ManuallyCreateOrUpdateUser("t3", "google", "googleid", "test@example.com")
	assert.Nil(t, err)
	user3b, err := ManuallyCreateOrUpdateUser("t3", "facebook", "fbid", "test@example.com")
	assert.Nil(t, err)

	assert.Equal(t, []string{"t1"}, user1a.OK.User.TenantIds)
	assert.Equal(t, []string{"t1"}, user1b.OK.User.TenantIds)

	assert.Equal(t, []string{"t2"}, user2a.OK.User.TenantIds)
	assert.Equal(t, []string{"t2"}, user2b.OK.User.TenantIds)

	assert.Equal(t, []string{"t3"}, user3a.OK.User.TenantIds)
	assert.Equal(t, []string{"t3"}, user3b.OK.User.TenantIds)

	// get user by id
	gUser1a, err := GetUserByID(user1a.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, user1a.OK.User, *gUser1a)
	gUser1b, err := GetUserByID(user1b.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, user1b.OK.User, *gUser1b)

	gUser2a, err := GetUserByID(user2a.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, user2a.OK.User, *gUser2a)
	gUser2b, err := GetUserByID(user2b.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, user2b.OK.User, *gUser2b)

	gUser3a, err := GetUserByID(user3a.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, user3a.OK.User, *gUser3a)
	gUser3b, err := GetUserByID(user3b.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, user3b.OK.User, *gUser3b)

	// get users by email
	gUsers1, err := GetUsersByEmail("t1", "test@example.com")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(gUsers1))
	assert.Equal(t, user1a.OK.User, gUsers1[0])
	assert.Equal(t, user1b.OK.User, gUsers1[1])

	gUsers2, err := GetUsersByEmail("t2", "test@example.com")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(gUsers2))
	assert.Equal(t, user2a.OK.User, gUsers2[0])
	assert.Equal(t, user2b.OK.User, gUsers2[1])

	gUsers3, err := GetUsersByEmail("t3", "test@example.com")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(gUsers3))
	assert.Equal(t, user3a.OK.User, gUsers3[0])
	assert.Equal(t, user3b.OK.User, gUsers3[1])
}

func TestGetProvider(t *testing.T) {
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
	resp, err := multitenancy.CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)
	resp, err = multitenancy.CreateOrUpdateTenant("t2", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	resp, err = multitenancy.CreateOrUpdateTenant("t3", multitenancymodels.TenantConfig{
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	multitenancy.CreateOrUpdateThirdPartyConfig("t1", tpmodels.ProviderConfig{
		ThirdPartyId: "google",
		Clients:      []tpmodels.ProviderClientConfig{{ClientID: "a"}},
	}, nil)
	multitenancy.CreateOrUpdateThirdPartyConfig("t1", tpmodels.ProviderConfig{
		ThirdPartyId: "facebook",
		Clients:      []tpmodels.ProviderClientConfig{{ClientID: "a"}},
	}, nil)

	multitenancy.CreateOrUpdateThirdPartyConfig("t2", tpmodels.ProviderConfig{
		ThirdPartyId: "discord",
		Clients:      []tpmodels.ProviderClientConfig{{ClientID: "a"}},
	}, nil)
	multitenancy.CreateOrUpdateThirdPartyConfig("t2", tpmodels.ProviderConfig{
		ThirdPartyId: "facebook",
		Clients:      []tpmodels.ProviderClientConfig{{ClientID: "a"}},
	}, nil)

	multitenancy.CreateOrUpdateThirdPartyConfig("t3", tpmodels.ProviderConfig{
		ThirdPartyId: "discord",
		Clients:      []tpmodels.ProviderClientConfig{{ClientID: "a"}},
	}, nil)
	multitenancy.CreateOrUpdateThirdPartyConfig("t3", tpmodels.ProviderConfig{
		ThirdPartyId: "linkedin",
		Clients:      []tpmodels.ProviderClientConfig{{ClientID: "a"}},
	}, nil)

	provider, err := GetProvider("t1", "google", nil)
	assert.Nil(t, err)
	assert.Equal(t, "google", provider.ID)
	provider, err = GetProvider("t1", "facebook", nil)
	assert.Nil(t, err)
	assert.Equal(t, "facebook", provider.ID)
	provider, err = GetProvider("t1", "discord", nil)
	assert.Nil(t, err)
	assert.Nil(t, provider)

	provider, err = GetProvider("t2", "google", nil)
	assert.Nil(t, err)
	assert.Nil(t, provider)
	provider, err = GetProvider("t2", "facebook", nil)
	assert.Nil(t, err)
	assert.Equal(t, "facebook", provider.ID)
	provider, err = GetProvider("t2", "discord", nil)
	assert.Nil(t, err)
	assert.Equal(t, "discord", provider.ID)
}

func TestGetProviderMergesConfigFromStaticAndCore(t *testing.T) {
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
			Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{ClientID: "staticclientid", ClientSecret: "staticclientsecret"},
								},
							},
						},
					},
				},
			}),
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
		ThirdPartyEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	multitenancy.CreateOrUpdateThirdPartyConfig("t1", tpmodels.ProviderConfig{
		ThirdPartyId: "google",
		Clients: []tpmodels.ProviderClientConfig{
			{ClientID: "coreclientid", ClientSecret: "coreclientsecret"},
		},
	}, nil)

	provider, err := GetProvider("t1", "google", nil)
	assert.Nil(t, err)
	assert.Equal(t, "google", provider.ID)
	assert.Equal(t, "coreclientid", provider.Config.ClientID)
	assert.Equal(t, "coreclientsecret", provider.Config.ClientSecret)
}
