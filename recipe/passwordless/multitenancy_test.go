package passwordless

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
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
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
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
		PasswordlessEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)
	resp, err = multitenancy.CreateOrUpdateTenant("t2", multitenancymodels.TenantConfig{
		PasswordlessEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	resp, err = multitenancy.CreateOrUpdateTenant("t3", multitenancymodels.TenantConfig{
		PasswordlessEnabled: &True,
	})
	assert.Nil(t, err)
	assert.True(t, resp.OK.CreatedNew)

	code := "123456"
	code1, err := CreateCodeWithEmail("t1", "test@example.com", &code)
	assert.Nil(t, err)
	assert.NotNil(t, code1.OK)

	code = "45678"
	code2, err := CreateCodeWithEmail("t2", "test@example.com", &code)
	assert.Nil(t, err)
	assert.NotNil(t, code2.OK)

	code = "789123"
	code3, err := CreateCodeWithEmail("t3", "test@example.com", &code)
	assert.Nil(t, err)
	assert.NotNil(t, code3.OK)

	user1, err := ConsumeCodeWithUserInputCode("t1", code1.OK.DeviceID, "123456", code1.OK.PreAuthSessionID)
	assert.Nil(t, err)

	user2, err := ConsumeCodeWithUserInputCode("t2", code2.OK.DeviceID, "45678", code2.OK.PreAuthSessionID)
	assert.Nil(t, err)

	user3, err := ConsumeCodeWithUserInputCode("t3", code3.OK.DeviceID, "789123", code3.OK.PreAuthSessionID)
	assert.Nil(t, err)

	assert.NotEqual(t, user1.OK.User.ID, user2.OK.User.ID)
	assert.NotEqual(t, user1.OK.User.ID, user3.OK.User.ID)
	assert.NotEqual(t, user2.OK.User.ID, user3.OK.User.ID)

	assert.Equal(t, []string{"t1"}, user1.OK.User.TenantIds)
	assert.Equal(t, []string{"t2"}, user2.OK.User.TenantIds)
	assert.Equal(t, []string{"t3"}, user3.OK.User.TenantIds)

	// get user by id

	gUser1, err := GetUserByID(user1.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, user1.OK.User, *gUser1)

	gUser2, err := GetUserByID(user2.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, user2.OK.User, *gUser2)

	gUser3, err := GetUserByID(user3.OK.User.ID)
	assert.Nil(t, err)
	assert.Equal(t, user3.OK.User, *gUser3)

	// get user by email
	gUserByEmail1, err := GetUserByEmail("t1", "test@example.com")
	assert.Nil(t, err)
	assert.Equal(t, user1.OK.User, *gUserByEmail1)

	gUserByEmail2, err := GetUserByEmail("t2", "test@example.com")
	assert.Nil(t, err)
	assert.Equal(t, user2.OK.User, *gUserByEmail2)

	gUserByEmail3, err := GetUserByEmail("t3", "test@example.com")
	assert.Nil(t, err)
	assert.Equal(t, user3.OK.User, *gUserByEmail3)
}
