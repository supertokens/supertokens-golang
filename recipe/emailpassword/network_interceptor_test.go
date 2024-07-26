package emailpassword

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

var isNetworkIntercepted = false

func TestNetworkInterceptorDuringSignIn(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
			NetworkInterceptor: func(request *http.Request, context supertokens.UserContext) (*http.Request, error) {
				isNetworkIntercepted = true
				return request, nil
			},
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	res, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, true, isNetworkIntercepted)
}

func TestNetworkInterceptorNotSet(t *testing.T) {
	isNetworkIntercepted = false
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	res, err := unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, false, isNetworkIntercepted)
}

func TestNetworkInterceptorIncorrectCoreURL(t *testing.T) {
	isNetworkIntercepted = false
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
			NetworkInterceptor: func(request *http.Request, context supertokens.UserContext) (*http.Request, error) {
				isNetworkIntercepted = true
				newRequest := request
				newRequest.URL.Path = "/public/recipe/incorrect/path"
				return newRequest, nil
			},
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
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

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	_, err = SignIn("public", "random@gmail.com", "validpass123")

	assert.NotNil(t, err, "there should be an error")
	assert.Contains(t, err.Error(), "status code: 404")
	assert.Equal(t, true, isNetworkIntercepted)
}

func TestNetworkInterceptorIncorrectQueryParams(t *testing.T) {
	isNetworkIntercepted = false
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
			NetworkInterceptor: func(r *http.Request, context supertokens.UserContext) (*http.Request, error) {
				isNetworkIntercepted = true
				newRequest := r
				q := url.Values{}
				newRequest.URL.RawQuery = q.Encode()
				return newRequest, nil
			},
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			userroles.Init(nil),
		},
	}
	BeforeEach()

	unittesting.StartUpST("localhost", "8080")

	defer AfterEach()

	supertokens.Init(configValue)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, _ := SignUp("public", "random@gmail.com", "validpass123")
	_, err := userroles.GetRolesForUser("public", resp.OK.User.ID)
	assert.NotNil(t, err, "should err, because userId is not passed")
	assert.Contains(t, err.Error(), "status code: 400")
	assert.Equal(t, true, isNetworkIntercepted)
}

func TestNetworkInterceptorRequestBody(t *testing.T) {
	isNetworkIntercepted = false
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
			NetworkInterceptor: func(r *http.Request, context supertokens.UserContext) (*http.Request, error) {
				isNetworkIntercepted = true
				newBody := bytes.NewReader([]byte(`{"newKey": "newValue"}`))
				req, _ := http.NewRequest(r.Method, r.URL.String(), newBody)
				req.Header = r.Header
				return req, nil
			},
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}
	BeforeEach()

	unittesting.StartUpST("localhost", "8080")

	defer AfterEach()

	supertokens.Init(configValue)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	_, err := SignIn("public", "random@gmail.com", "validpass123")
	assert.NotNil(t, err, "should err, because request body is incorrect")
	assert.Contains(t, err.Error(), "status code: 400")
	assert.Equal(t, true, isNetworkIntercepted)
}
