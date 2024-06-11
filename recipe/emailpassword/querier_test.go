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
	assert.Nil(t, tpuser)
	assert.False(t, calledCore)

	user, err = GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.False(t, calledCore)
}

func TestNoCachingIfDisabledByUser(t *testing.T) {
	calledCore := false

	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
			NetworkInterceptor: func(r *http.Request, uc supertokens.UserContext) *http.Request {
				calledCore = true
				return r
			},
			DisableCoreCallCache: true,
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "ST",
			WebsiteDomain: "http://supertokens.io",
			APIDomain:     "http://api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
			Init(nil),
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
	assert.True(t, calledCore)
}

func TestNoCachingIfHeadersAreDifferent(t *testing.T) {
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

	calledCore = false

	tpuser, err := thirdparty.GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, tpuser)
	assert.True(t, calledCore)
}

func TestCachingGetsClearWhenQueryWithoutUserContext(t *testing.T) {
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

	_, err = SignUp("public", "test@example.com", "abcd1234")
	assert.Nil(t, err)

	calledCore = false

	user, err = GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)
}

func TestCachingDoesNotGetClearWithNonGetIfKeepAlive(t *testing.T) {
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
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(config)
	if err != nil {
		t.Error(err.Error())
	}

	userContext := map[string]interface{}{
		"_default": map[string]interface{}{
			"keepCacheAlive": true,
		},
	}
	userContext2 := map[string]interface{}{}

	user, err := GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)

	calledCore = false

	user, err = GetUserByID("random", &userContext2)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)

	_, err = SignUp("public", "test@example.com", "abcd1234", &userContext)
	assert.Nil(t, err)

	calledCore = false

	user, err = GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)

	calledCore = false

	user, err = GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.False(t, calledCore)

	user, err = GetUserByID("random", &userContext2)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.False(t, calledCore)
}

func TestCachingGetsClearWithNonGetIfKeepAliveIsFalse(t *testing.T) {
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
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(config)
	if err != nil {
		t.Error(err.Error())
	}

	userContext := map[string]interface{}{
		"_default": map[string]interface{}{
			"keepCacheAlive": false,
		},
	}
	userContext2 := map[string]interface{}{}

	user, err := GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)

	calledCore = false

	user, err = GetUserByID("random", &userContext2)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)

	_, err = SignUp("public", "test@example.com", "abcd1234", &userContext)
	assert.Nil(t, err)

	calledCore = false

	user, err = GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)

	calledCore = false

	user, err = GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.False(t, calledCore)

	user, err = GetUserByID("random", &userContext2)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)
}

func TestCachingGetsClearWithNonGetIfKeepAliveIsNotSet(t *testing.T) {
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
	userContext2 := map[string]interface{}{}

	user, err := GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)

	calledCore = false

	user, err = GetUserByID("random", &userContext2)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)

	_, err = SignUp("public", "test@example.com", "abcd1234", &userContext)
	assert.Nil(t, err)

	calledCore = false

	user, err = GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)

	calledCore = false

	user, err = GetUserByID("random", &userContext)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.False(t, calledCore)

	user, err = GetUserByID("random", &userContext2)
	assert.Nil(t, err)
	assert.Nil(t, user)
	assert.True(t, calledCore)
}
