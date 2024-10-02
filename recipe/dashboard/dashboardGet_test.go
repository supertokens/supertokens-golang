package dashboard

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestThatDashboardGetNormalizesConnectionURIWithoutHTTP(t *testing.T) {
	connectionURI := "http://localhost:8080"
	connectionURIWithoutProtocol := strings.Replace(connectionURI, "http://", "", -1)
	config := supertokens.TypeInput{
		OnSuperTokensAPIError: func(err error, req *http.Request, res http.ResponseWriter) {
			print(err)
		},
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURIWithoutProtocol,
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			Init(&dashboardmodels.TypeInput{
				ApiKey: "testapikey",
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(config)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/dashboard", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer testapikey")
	res, err := http.DefaultClient.Do(req)
	assert.Equal(t, res.StatusCode, 200)

	if err != nil {
		t.Error(err.Error())
	}

	body, _ := io.ReadAll(res.Body)
	assert.True(t, strings.Contains(string(body), fmt.Sprintf("window.connectionURI = \"%s\"", connectionURI)))
}

func TestThatDashboardGetNormalizesConnectionURIWithoutHTTPS(t *testing.T) {
	connectionURI := "https://try.supertokens.com"
	connectionURIWithoutProtocol := strings.Replace(connectionURI, "https://", "", -1)
	config := supertokens.TypeInput{
		OnSuperTokensAPIError: func(err error, req *http.Request, res http.ResponseWriter) {
			print(err)
		},
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURIWithoutProtocol,
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			Init(&dashboardmodels.TypeInput{
				ApiKey: "testapikey",
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(config)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/dashboard", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer testapikey")
	res, err := http.DefaultClient.Do(req)
	assert.Equal(t, res.StatusCode, 200)

	if err != nil {
		t.Error(err.Error())
	}

	body, _ := io.ReadAll(res.Body)
	assert.True(t, strings.Contains(string(body), fmt.Sprintf("window.connectionURI = \"%s\"", connectionURI)))
}

func TestThatDashboardGetReturnsFirstURIWhenMultipleArePassed(t *testing.T) {
	firstConnectionURI := "http://localhost:8080"
	secondConnectionURI := "https://try.supertokens.com"
	multiplConnectionURIs := fmt.Sprintf("%s;%s", firstConnectionURI, secondConnectionURI)
	config := supertokens.TypeInput{
		OnSuperTokensAPIError: func(err error, req *http.Request, res http.ResponseWriter) {
			print(err)
		},
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: multiplConnectionURIs,
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			Init(&dashboardmodels.TypeInput{
				ApiKey: "testapikey",
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(config)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/dashboard", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer testapikey")
	res, err := http.DefaultClient.Do(req)
	assert.Equal(t, res.StatusCode, 200)

	if err != nil {
		t.Error(err.Error())
	}

	body, _ := io.ReadAll(res.Body)
	assert.True(t, strings.Contains(string(body), fmt.Sprintf("window.connectionURI = \"%s\"", firstConnectionURI)))
}
