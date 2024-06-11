package emailpassword

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestCachingWorks(t *testing.T) {
	calledCore := false

	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
			NetworkInterceptor: func(r *http.Request, uc supertokens.UserContext) *http.Request {
				calledCore = true
				return r
			},
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "ST",
			WebsiteDomain: "http://supertokens.io",
			APIDomain:     "http://api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
			Init(nil),
			thirdparty.Init(nil),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(config)
	if err != nil {
		t.Error(err.Error())
	}

	userContext := map[string]interface{}{}
	user, err := GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)

	calledCore = false

	user, err = GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.False(t, calledCore)

	tpuser, err := thirdparty.GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, tpuser)
	assert.True(t, calledCore)

	calledCore = false

	tpuser, err = thirdparty.GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.False(t, calledCore)

	user, err = GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.False(t, calledCore)
}
