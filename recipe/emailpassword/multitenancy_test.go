package emailpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
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

	// sign up
	user1, err := SignUp("t1", "test@example.com", "password1")
	assert.Nil(t, err)

	user2, err := SignUp("t2", "test@example.com", "password2")
	assert.Nil(t, err)

	user3, err := SignUp("t3", "test@example.com", "password3")
	assert.Nil(t, err)

	assert.NotEqual(t, user1.OK.User.ID, user2.OK.User.ID)
	assert.NotEqual(t, user1.OK.User.ID, user3.OK.User.ID)
	assert.NotEqual(t, user2.OK.User.ID, user3.OK.User.ID)

	assert.Equal(t, []string{"t1"}, user1.OK.User.TenantIds)
	assert.Equal(t, []string{"t2"}, user2.OK.User.TenantIds)
	assert.Equal(t, []string{"t3"}, user3.OK.User.TenantIds)

	// sign in
	sUser1, err := SignIn("t1", "test@example.com", "password1")
	assert.Nil(t, err)
	assert.Equal(t, user1.OK.User, sUser1.OK.User)

	sUser2, err := SignIn("t2", "test@example.com", "password2")
	assert.Nil(t, err)
	assert.Equal(t, user2.OK.User, sUser2.OK.User)

	sUser3, err := SignIn("t3", "test@example.com", "password3")
	assert.Nil(t, err)
	assert.Equal(t, user3.OK.User, sUser3.OK.User)

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

	gUserByEmail1, err := GetUserByEmail("t1", "test@example.com")
	assert.Nil(t, err)
	assert.Equal(t, user1.OK.User, *gUserByEmail1)

	gUserByEmail2, err := GetUserByEmail("t2", "test@example.com")
	assert.Nil(t, err)
	assert.Equal(t, user2.OK.User, *gUserByEmail2)

	gUserByEmail3, err := GetUserByEmail("t3", "test@example.com")
	assert.Nil(t, err)
	assert.Equal(t, user3.OK.User, *gUserByEmail3)

	// password reset token
	passwordResetToken1, err := CreateResetPasswordToken("t1", user1.OK.User.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, passwordResetToken1.OK.Token)

	passwordResetToken2, err := CreateResetPasswordToken("t2", user2.OK.User.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, passwordResetToken2.OK.Token)

	passwordResetToken3, err := CreateResetPasswordToken("t3", user3.OK.User.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, passwordResetToken3.OK.Token)

	// reset password
	_, err = ResetPasswordUsingToken("t1", passwordResetToken1.OK.Token, "newpassword1")
	assert.Nil(t, err)

	_, err = ResetPasswordUsingToken("t2", passwordResetToken2.OK.Token, "newpassword2")
	assert.Nil(t, err)

	_, err = ResetPasswordUsingToken("t3", passwordResetToken3.OK.Token, "newpassword3")
	assert.Nil(t, err)

	// sign in with new password
	sUser1, err = SignIn("t1", "test@example.com", "newpassword1")
	assert.Nil(t, err)
	assert.Equal(t, user1.OK.User, sUser1.OK.User)

	sUser2, err = SignIn("t2", "test@example.com", "newpassword2")
	assert.Nil(t, err)
	assert.Equal(t, user2.OK.User, sUser2.OK.User)

	sUser3, err = SignIn("t3", "test@example.com", "newpassword3")
	assert.Nil(t, err)
	assert.Equal(t, user3.OK.User, sUser3.OK.User)
}
