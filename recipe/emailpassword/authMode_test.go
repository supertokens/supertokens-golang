package emailpassword

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

const exampleJWTForTest string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestWithDefaultGetTokenTransferMethodCreateNewSessionShouldDefaultToHeaderBasedSession(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	resp := createNewSession(t, testServer.URL, nil, nil, nil, nil)

	assert.Equal(t, resp["sAccessToken"], "-not-present-")
	assert.Equal(t, resp["sRefreshToken"], "-not-present-")
	assert.Equal(t, resp["antiCsrf"], "-not-present-")
	assert.NotEmpty(t, resp["accessTokenFromHeader"])
	assert.NotEqual(t, resp["accessTokenFromHeader"], "-not-present-")
	assert.NotEmpty(t, resp["refreshTokenFromHeader"])
	assert.NotEqual(t, resp["refreshTokenFromHeader"], "-not-present-")
}

func TestWithDefaultGetTokenTransferMethodCreateNewSessionWithBadAuthModeHeaderShouldDefaultToHeaderBasedSession(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	authMode := "badauthmode"
	resp := createNewSession(t, testServer.URL, &authMode, nil, nil, nil)

	assert.Equal(t, resp["sAccessToken"], "-not-present-")
	assert.Equal(t, resp["sRefreshToken"], "-not-present-")
	assert.Equal(t, resp["antiCsrf"], "-not-present-")
	assert.NotEmpty(t, resp["accessTokenFromHeader"])
	assert.NotEqual(t, resp["accessTokenFromHeader"], "-not-present-")
	assert.NotEmpty(t, resp["refreshTokenFromHeader"])
	assert.NotEqual(t, resp["refreshTokenFromHeader"], "-not-present-")
}

func TestWithDefaultGetTokenTransferMethodCreateNewSessionWithAuthModeSpecifiedAsHeader(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	authMode := string(sessmodels.HeaderTransferMethod)
	resp := createNewSession(t, testServer.URL, &authMode, nil, nil, nil)

	assert.Equal(t, resp["sAccessToken"], "-not-present-")
	assert.Equal(t, resp["sRefreshToken"], "-not-present-")
	assert.Equal(t, resp["antiCsrf"], "-not-present-")
	assert.NotEmpty(t, resp["accessTokenFromHeader"])
	assert.NotEqual(t, resp["accessTokenFromHeader"], "-not-present-")
	assert.NotEmpty(t, resp["refreshTokenFromHeader"])
	assert.NotEqual(t, resp["refreshTokenFromHeader"], "-not-present-")
}

func TestWithDefaultGetTokenTransferMethodCreateNewSessionWithAuthModeSpecifiedAsCookie(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	authMode := string(sessmodels.CookieTransferMethod)
	resp := createNewSession(t, testServer.URL, &authMode, nil, nil, nil)

	assert.NotEmpty(t, resp["sAccessToken"])
	assert.NotEqual(t, resp["sAccessToken"], "-not-present-")
	assert.NotEmpty(t, resp["sRefreshToken"])
	assert.NotEqual(t, resp["sRefreshToken"], "-not-present-")
	assert.NotEmpty(t, resp["antiCsrf"])
	assert.NotEqual(t, resp["antiCsrf"], "-not-present-")
	assert.Equal(t, resp["accessTokenFromHeader"], "-not-present-")
	assert.Equal(t, resp["refreshTokenFromHeader"], "-not-present-")
}

func TestWithGetTokenTransferMethodProvidedCreateNewSessionWithShouldUseHeaderIfMethodReturnsAny(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.AnyTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	resp := createNewSession(t, testServer.URL, nil, nil, nil, nil)

	assert.Equal(t, resp["sAccessToken"], "-not-present-")
	assert.Equal(t, resp["sRefreshToken"], "-not-present-")
	assert.Equal(t, resp["antiCsrf"], "-not-present-")
	assert.NotEmpty(t, resp["accessTokenFromHeader"])
	assert.NotEqual(t, resp["accessTokenFromHeader"], "-not-present-")
	assert.NotEmpty(t, resp["refreshTokenFromHeader"])
	assert.NotEqual(t, resp["refreshTokenFromHeader"], "-not-present-")
}

func TestWithGetTokenTransferMethodProvidedCreateNewSessionWithShouldUseHeaderIfMethodReturnsHeader(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.HeaderTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	resp := createNewSession(t, testServer.URL, nil, nil, nil, nil)

	assert.Equal(t, resp["sAccessToken"], "-not-present-")
	assert.Equal(t, resp["sRefreshToken"], "-not-present-")
	assert.Equal(t, resp["antiCsrf"], "-not-present-")
	assert.NotEmpty(t, resp["accessTokenFromHeader"])
	assert.NotEqual(t, resp["accessTokenFromHeader"], "-not-present-")
	assert.NotEmpty(t, resp["refreshTokenFromHeader"])
	assert.NotEqual(t, resp["refreshTokenFromHeader"], "-not-present-")
}

func TestWithGetTokenTransferMethodProvidedCreateNewSessionWithShouldClearCookiesIfMethodReturnsHeader(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.HeaderTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	cookies := []http.Cookie{
		{
			Name:  "sAccessToken",
			Value: exampleJWTForTest,
		},
	}

	resp := createNewSession(t, testServer.URL, nil, nil, cookies, nil)

	assert.Empty(t, resp["sAccessToken"])
	assert.Equal(t, resp["accessTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Empty(t, resp["sRefreshToken"])
	assert.Equal(t, resp["refreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, resp["antiCsrf"], "-not-present-")
}

func TestWithGetTokenTransferMethodProvidedCreateNewSessionWithShouldUseCookieIfMethodReturnsCookie(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	resp := createNewSession(t, testServer.URL, nil, nil, nil, nil)

	assert.NotEmpty(t, resp["sAccessToken"])
	assert.NotEqual(t, resp["sAccessToken"], "-not-present-")
	assert.NotEmpty(t, resp["sRefreshToken"])
	assert.NotEqual(t, resp["sRefreshToken"], "-not-present-")
	assert.NotEmpty(t, resp["antiCsrf"])
	assert.NotEqual(t, resp["antiCsrf"], "-not-present-")
	assert.Equal(t, resp["accessTokenFromHeader"], "-not-present-")
	assert.Equal(t, resp["refreshTokenFromHeader"], "-not-present-")
}

func TestWithGetTokenTransferMethodProvidedCreateNewSessionWithShouldClearHeaderIfMethodReturnsCookie(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	headers := map[string]string{
		"Authorization": "Bearer " + exampleJWTForTest,
	}

	resp := createNewSession(t, testServer.URL, nil, nil, nil, headers)

	assert.Empty(t, resp["accessTokenFromHeader"])
	assert.Empty(t, resp["refreshTokenFromHeader"])
}

func TestVerifySessionBehaviour(t *testing.T) {
	viaToken := "VIA_TOKEN"

	behaviourTable := []struct {
		getTokenTransferMethodRes string
		sessionRequired           bool
		authHeader                bool
		authCookie                bool
		output                    string
	}{
		{getTokenTransferMethodRes: "any", sessionRequired: false, authHeader: false, authCookie: false, output: "undefined"},
		{getTokenTransferMethodRes: "header", sessionRequired: false, authHeader: false, authCookie: false, output: "undefined"},
		{getTokenTransferMethodRes: "cookie", sessionRequired: false, authHeader: false, authCookie: false, output: "undefined"},
		{getTokenTransferMethodRes: "cookie", sessionRequired: false, authHeader: true, authCookie: false, output: "undefined"},
		{getTokenTransferMethodRes: "header", sessionRequired: false, authHeader: false, authCookie: true, output: "undefined"},
		{getTokenTransferMethodRes: "any", sessionRequired: true, authHeader: false, authCookie: false, output: "UNAUTHORISED"},
		{getTokenTransferMethodRes: "header", sessionRequired: true, authHeader: false, authCookie: false, output: "UNAUTHORISED"},
		{getTokenTransferMethodRes: "cookie", sessionRequired: true, authHeader: false, authCookie: false, output: "UNAUTHORISED"},
		{getTokenTransferMethodRes: "cookie", sessionRequired: true, authHeader: true, authCookie: false, output: "UNAUTHORISED"},
		{getTokenTransferMethodRes: "header", sessionRequired: true, authHeader: false, authCookie: true, output: "UNAUTHORISED"},
		{getTokenTransferMethodRes: "any", sessionRequired: true, authHeader: true, authCookie: true, output: "validateheader"},
		{getTokenTransferMethodRes: "any", sessionRequired: false, authHeader: true, authCookie: true, output: "validateheader"},
		{getTokenTransferMethodRes: "header", sessionRequired: true, authHeader: true, authCookie: true, output: "validateheader"},
		{getTokenTransferMethodRes: "header", sessionRequired: false, authHeader: true, authCookie: true, output: "validateheader"},
		{getTokenTransferMethodRes: "cookie", sessionRequired: true, authHeader: true, authCookie: true, output: "validatecookie"},
		{getTokenTransferMethodRes: "cookie", sessionRequired: false, authHeader: true, authCookie: true, output: "validatecookie"},
		{getTokenTransferMethodRes: "any", sessionRequired: true, authHeader: true, authCookie: false, output: "validateheader"},
		{getTokenTransferMethodRes: "any", sessionRequired: false, authHeader: true, authCookie: false, output: "validateheader"},
		{getTokenTransferMethodRes: "header", sessionRequired: true, authHeader: true, authCookie: false, output: "validateheader"},
		{getTokenTransferMethodRes: "header", sessionRequired: false, authHeader: true, authCookie: false, output: "validateheader"},
		{getTokenTransferMethodRes: "any", sessionRequired: true, authHeader: false, authCookie: true, output: "validatecookie"},
		{getTokenTransferMethodRes: "any", sessionRequired: false, authHeader: false, authCookie: true, output: "validatecookie"},
		{getTokenTransferMethodRes: "cookie", sessionRequired: true, authHeader: false, authCookie: true, output: "validatecookie"},
		{getTokenTransferMethodRes: "cookie", sessionRequired: false, authHeader: false, authCookie: true, output: "validatecookie"},
	}

	for _, behaviour := range behaviourTable {
		t.Run(fmt.Sprintf("behaviour: %v with valid token", behaviour), func(t *testing.T) {
			configValue := supertokens.TypeInput{
				Supertokens: &supertokens.ConnectionInfo{
					ConnectionURI: "http://localhost:8080",
				},
				AppInfo: supertokens.AppInfo{
					APIDomain:     "api.supertokens.io",
					AppName:       "SuperTokens",
					WebsiteDomain: "supertokens.io",
				},
				RecipeList: []supertokens.Recipe{
					Init(nil),
					session.Init(&sessmodels.TypeInput{
						AntiCsrf: &viaToken,
						GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
							return sessmodels.TokenTransferMethod(behaviour.getTokenTransferMethodRes)
						},
					}),
				},
			}

			BeforeEach()
			unittesting.StartUpST("localhost", "8080")
			defer AfterEach()

			err := supertokens.Init(configValue)
			if err != nil {
				assert.NoError(t, err)
			}

			mux := http.NewServeMux()
			testServer := httptest.NewServer(supertokens.Middleware(mux))
			defer testServer.Close()
			setupRoutesForTest(t, mux)

			cookie := "cookie"
			createInfo := createNewSession(t, testServer.URL, &cookie, nil, nil, nil)
			fmt.Println(createInfo)

			authMode := ""
			if behaviour.authCookie && behaviour.authHeader {
				authMode = "both"
			} else if behaviour.authCookie {
				authMode = "cookie"
			} else if behaviour.authHeader {
				authMode = "header"
			} else {
				authMode = "none"
			}

			expectedStatus := 200
			if behaviour.output == "UNAUTHORISED" {
				expectedStatus = 401
			}

			testRes := testGetVerify(t, testServer.URL, createInfo, behaviour.sessionRequired, expectedStatus, authMode)
			switch behaviour.output {
			case "undefined":
				assert.Equal(t, testRes["sessionExists"], false)
			case "UNAUTHORISED":
				assert.Equal(t, testRes["message"], "unauthorised")
			case "validateheader":
				assert.Equal(t, testRes["sessionExists"], true)
			case "validatecookie":
				assert.Equal(t, testRes["sessionExists"], true)
			}
		})

		t.Run(fmt.Sprintf("behaviour: %v with expired token", behaviour), func(t *testing.T) {
			configValue := supertokens.TypeInput{
				Supertokens: &supertokens.ConnectionInfo{
					ConnectionURI: "http://localhost:8080",
				},
				AppInfo: supertokens.AppInfo{
					APIDomain:     "api.supertokens.io",
					AppName:       "SuperTokens",
					WebsiteDomain: "supertokens.io",
				},
				RecipeList: []supertokens.Recipe{
					Init(nil),
					session.Init(&sessmodels.TypeInput{
						AntiCsrf: &viaToken,
						GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
							return sessmodels.TokenTransferMethod(behaviour.getTokenTransferMethodRes)
						},
					}),
				},
			}

			BeforeEach()
			unittesting.SetKeyValueInConfig("access_token_validity", "2")
			unittesting.StartUpST("localhost", "8080")
			defer AfterEach()

			err := supertokens.Init(configValue)
			if err != nil {
				assert.NoError(t, err)
			}

			mux := http.NewServeMux()
			testServer := httptest.NewServer(supertokens.Middleware(mux))
			defer testServer.Close()
			setupRoutesForTest(t, mux)

			cookie := "cookie"
			createInfo := createNewSession(t, testServer.URL, &cookie, nil, nil, nil)
			time.Sleep(3 * time.Second)

			authMode := ""
			if behaviour.authCookie && behaviour.authHeader {
				authMode = "both"
			} else if behaviour.authCookie {
				authMode = "cookie"
			} else if behaviour.authHeader {
				authMode = "header"
			} else {
				authMode = "none"
			}

			expectedStatus := 401
			if behaviour.output == "undefined" {
				expectedStatus = 200
			}
			testRes := testGetVerify(t, testServer.URL, createInfo, behaviour.sessionRequired, expectedStatus, authMode)
			switch behaviour.output {
			case "undefined":
				assert.Equal(t, testRes["sessionExists"], false)
			case "UNAUTHORISED":
				assert.Equal(t, testRes["message"], "unauthorised")
			case "validateheader":
				assert.Equal(t, testRes["message"], "try refresh token")
			case "validatecookie":
				assert.Equal(t, testRes["message"], "try refresh token")
			}
		})
	}
}

func TestWithAccessTokenInBothHeaderAndCookieShouldUseHeadersIfMethodReturnsAny(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.AnyTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	header := "header"
	createInfoCookie := createNewSession(t, testServer.URL, &header, nil, nil, nil)
	createInfoHeader := createNewSession(t, testServer.URL, &header, nil, nil, nil)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+createInfoHeader["accessTokenFromHeader"])
	req.Header.Set("Cookie", "sAccessToken="+url.QueryEscape(createInfoCookie["accessTokenFromHeader"]))
	if createInfoCookie["antiCsrf"] != "" {
		req.Header.Set("anti-csrf", createInfoCookie["antiCsrf"])
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	result := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, result["sessionHandle"], createInfoHeader["sessionHandle"])
}

func TestWithAccessTokenInBothHeaderAndCookieShouldUseHeadersIfMethodReturnsHeader(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.HeaderTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	header := "header"
	createInfoCookie := createNewSession(t, testServer.URL, &header, nil, nil, nil)
	createInfoHeader := createNewSession(t, testServer.URL, &header, nil, nil, nil)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+createInfoHeader["accessTokenFromHeader"])
	req.Header.Set("Cookie", "sAccessToken="+url.QueryEscape(createInfoCookie["accessTokenFromHeader"]))
	if createInfoCookie["antiCsrf"] != "" {
		req.Header.Set("anti-csrf", createInfoCookie["antiCsrf"])
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	result := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, result["sessionHandle"], createInfoHeader["sessionHandle"])
}

func TestWithAccessTokenInBothHeaderAndCookieShouldUseCookieIfMethodReturnsCookie(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					if forCreateNewSession {
						return sessmodels.AnyTransferMethod
					}
					return sessmodels.CookieTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	header := "header"
	createInfoCookie := createNewSession(t, testServer.URL, &header, nil, nil, nil)
	createInfoHeader := createNewSession(t, testServer.URL, &header, nil, nil, nil)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+createInfoHeader["accessTokenFromHeader"])
	req.Header.Set("Cookie", "sAccessToken="+url.QueryEscape(createInfoCookie["accessTokenFromHeader"]))
	if createInfoCookie["antiCsrf"] != "" {
		req.Header.Set("anti-csrf", createInfoCookie["antiCsrf"])
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	result := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, result["sessionHandle"], createInfoCookie["sessionHandle"])
}

func TestWithAccessTokenInBothHeaderAndCookieShouldRejectRequestWithsIdRefreshToken(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	cookie := "cookie"
	createInfo := createNewSession(t, testServer.URL, &cookie, nil, nil, nil)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	assert.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "sAccessToken",
		Value: createInfo["sAccessToken"],
	})
	req.AddCookie((&http.Cookie{
		Name:  "sIdRefreshToken",
		Value: createInfo["sRefreshToken"], // The value doesn't actually matter
	}))
	if createInfo["antiCsrf"] != "" {
		req.Header.Set("anti-csrf", createInfo["antiCsrf"])
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	result := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, resp.StatusCode, 401)
	assert.Equal(t, result, map[string]interface{}{"message": "try refresh token"})
}

func TestWithNonSTAuthorizeHeaderShouldUseCookiesIfPresentAndMethodReturnsAny(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.AnyTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	header := "header"
	createInfo := createNewSession(t, testServer.URL, &header, nil, nil, nil)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	assert.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "sAccessToken",
		Value: url.QueryEscape(createInfo["accessTokenFromHeader"]),
	})
	if createInfo["antiCsrf"] != "" {
		req.Header.Set("anti-csrf", createInfo["antiCsrf"])
	}
	req.Header.Set("Authorization", "Bearer "+exampleJWTForTest)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	result := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, createInfo["sessionHandle"], result["sessionHandle"])
}

func TestWithNonSTAuthorizeHeaderShouldRejectWithUnauthorisedIfMethodReturnsHeader(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					if forCreateNewSession {
						return sessmodels.AnyTransferMethod
					}
					return sessmodels.HeaderTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	header := "header"
	createInfo := createNewSession(t, testServer.URL, &header, nil, nil, nil)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	assert.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "sAccessToken",
		Value: url.QueryEscape(createInfo["accessTokenFromHeader"]),
	})
	if createInfo["antiCsrf"] != "" {
		req.Header.Set("anti-csrf", createInfo["antiCsrf"])
	}
	req.Header.Set("Authorization", "Bearer "+exampleJWTForTest)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	result := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, resp.StatusCode, 401)
	assert.Equal(t, result, map[string]interface{}{"message": "unauthorised"})
}

func TestWithNonSTAuthorizeHeaderShouldRejectWithUnauthorisedIfCookiesAreNotPresent(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.AnyTransferMethod
				},
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+exampleJWTForTest)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	result := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, resp.StatusCode, 401)
	assert.Equal(t, result, map[string]interface{}{"message": "unauthorised"})
}

func TestMergeIntoAccessTokenPayloadShouldUpdateCookiesIfSessionWasCookieBased(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	header := "header"
	createInfo := createNewSession(t, testServer.URL, &header, nil, nil, nil)
	updateInfo := testGetUrl(t, testServer.URL, createInfo, "/update-payload", 200, "cookie", nil)

	// Didn't update
	assert.Equal(t, updateInfo["sRefreshToken"], "-not-present-")
	assert.Equal(t, updateInfo["antiCsrf"], "-not-present-")
	assert.Equal(t, updateInfo["accessTokenFromHeader"], "-not-present-")
	assert.Equal(t, updateInfo["refreshTokenFromHeader"], "-not-present-")

	// Updated access token
	assert.NotEmpty(t, updateInfo["sAccessToken"])
	assert.NotEqual(t, updateInfo["sAccessToken"], "-not-present-")
	assert.NotEqual(t, updateInfo["sAccessToken"], url.QueryEscape(createInfo["accessTokenFromHeader"]))

	// Update front token
	assert.NotEmpty(t, updateInfo["frontToken"])
	assert.NotEqual(t, updateInfo["frontToken"], "-not-present-")
	assert.NotEqual(t, updateInfo["frontToken"], createInfo["frontToken"])
}

func TestMergeIntoAccessTokenPayloadShouldAllowHeadersIfSessionWasHeaderBased(t *testing.T) {
	viaToken := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &viaToken,
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
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()
	setupRoutesForTest(t, mux)

	cookie := "cookie"
	createInfo := createNewSession(t, testServer.URL, &cookie, nil, nil, nil)

	updateInfo := testGetUrl(t, testServer.URL, createInfo, "/update-payload", 200, "header", nil)

	// Didn't update
	assert.Equal(t, updateInfo["sAccessToken"], "-not-present-")
	assert.Empty(t, updateInfo["sRefresthToken"])
	assert.Equal(t, updateInfo["antiCsrf"], "-not-present-")
	assert.Equal(t, updateInfo["refreshTokenFromHeader"], "-not-present-")

	// Updated access token
	assert.NotEmpty(t, updateInfo["accessTokenFromHeader"])
	assert.NotEqual(t, url.QueryEscape(updateInfo["accessTokenFromHeader"].(string)), createInfo["sAccessToken"])

	// Update front token
	assert.NotEmpty(t, updateInfo["frontToken"])
	assert.NotEqual(t, updateInfo["frontToken"], "-not-present-")
	assert.NotEqual(t, updateInfo["frontToken"], createInfo["frontToken"])
}

func TestRefreshTokenBehaviour(t *testing.T) {
	viaToken := "VIA_TOKEN"

	behaviourTable := []struct {
		getTokenTransferMethodRes string
		authHeader                bool
		authCookie                bool
		output                    string
		setTokens                 string
		clearedTokens             string
	}{
		{getTokenTransferMethodRes: "any", authHeader: false, authCookie: false, output: "unauthorised", setTokens: "none", clearedTokens: "none"},
		{getTokenTransferMethodRes: "header", authHeader: false, authCookie: false, output: "unauthorised", setTokens: "none", clearedTokens: "none"},
		{getTokenTransferMethodRes: "cookie", authHeader: false, authCookie: false, output: "unauthorised", setTokens: "none", clearedTokens: "none"},
		{getTokenTransferMethodRes: "any", authHeader: false, authCookie: true, output: "validatecookie", setTokens: "cookies", clearedTokens: "none"},
		{getTokenTransferMethodRes: "header", authHeader: false, authCookie: true, output: "unauthorised", setTokens: "none", clearedTokens: "none"},
		{getTokenTransferMethodRes: "cookie", authHeader: false, authCookie: true, output: "validatecookie", setTokens: "cookies", clearedTokens: "none"},
		{getTokenTransferMethodRes: "any", authHeader: true, authCookie: false, output: "validateheader", setTokens: "headers", clearedTokens: "none"},
		{getTokenTransferMethodRes: "header", authHeader: true, authCookie: false, output: "validateheader", setTokens: "headers", clearedTokens: "none"},
		{getTokenTransferMethodRes: "cookie", authHeader: true, authCookie: false, output: "unauthorised", setTokens: "none", clearedTokens: "none"},
		{getTokenTransferMethodRes: "any", authHeader: true, authCookie: true, output: "validateheader", setTokens: "headers", clearedTokens: "cookies"},
		{getTokenTransferMethodRes: "header", authHeader: true, authCookie: true, output: "validateheader", setTokens: "headers", clearedTokens: "cookies"},
		{getTokenTransferMethodRes: "cookie", authHeader: true, authCookie: true, output: "validatecookie", setTokens: "cookies", clearedTokens: "headers"},
	}

	for _, behaviour := range behaviourTable {
		t.Run(fmt.Sprintf("behaviour %v with valid token", behaviour), func(t *testing.T) {
			configValue := supertokens.TypeInput{
				Supertokens: &supertokens.ConnectionInfo{
					ConnectionURI: "http://localhost:8080",
				},
				AppInfo: supertokens.AppInfo{
					APIDomain:     "api.supertokens.io",
					AppName:       "SuperTokens",
					WebsiteDomain: "supertokens.io",
				},
				RecipeList: []supertokens.Recipe{
					Init(nil),
					session.Init(&sessmodels.TypeInput{
						AntiCsrf: &viaToken,
						GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
							return sessmodels.TokenTransferMethod(behaviour.getTokenTransferMethodRes)
						},
					}),
				},
			}

			BeforeEach()
			unittesting.StartUpST("localhost", "8080")
			defer AfterEach()

			err := supertokens.Init(configValue)
			if err != nil {
				assert.NoError(t, err)
			}

			mux := http.NewServeMux()
			testServer := httptest.NewServer(supertokens.Middleware(mux))
			defer testServer.Close()
			setupRoutesForTest(t, mux)

			header := "header" // Which we create doesn't really matter, since the token is the same
			createInfo := createNewSession(t, testServer.URL, &header, nil, nil, nil)

			authMode := ""
			if behaviour.authCookie && behaviour.authHeader {
				authMode = "both"
			} else if behaviour.authCookie {
				authMode = "cookie"
			} else if behaviour.authHeader {
				authMode = "header"
			} else {
				authMode = "none"
			}

			expectedStatus := 200
			if behaviour.output == "unauthorised" {
				expectedStatus = 401
			}

			refreshRes := refreshSession(t, testServer.URL, &behaviour.getTokenTransferMethodRes, authMode, createInfo, expectedStatus)

			if behaviour.clearedTokens == "headers" {
				assert.Empty(t, refreshRes["accessTokenFromHeader"])
				assert.Empty(t, refreshRes["refreshTokenFromHeader"])
			} else if behaviour.clearedTokens == "cookies" {
				assert.Empty(t, refreshRes["sAccessToken"])
				assert.Equal(t, refreshRes["accessTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
				assert.Empty(t, refreshRes["sRefreshToken"])
				assert.Equal(t, refreshRes["refreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
			}

			switch behaviour.setTokens {
			case "headers":
				assert.NotEmpty(t, refreshRes["accessTokenFromHeader"])
				assert.NotEqual(t, refreshRes["accessTokenFromHeader"], "-not-present-")
				assert.NotEmpty(t, refreshRes["refreshTokenFromHeader"])
				assert.NotEqual(t, refreshRes["refreshTokenFromHeader"], "-not-present-")

			case "cookies":
				assert.NotEmpty(t, refreshRes["sAccessToken"])
				assert.NotEqual(t, refreshRes["sAccessToken"], "-not-present-")
				assert.NotEqual(t, refreshRes["accessTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
				assert.NotEmpty(t, refreshRes["sRefreshToken"])
				assert.NotEqual(t, refreshRes["sRefreshToken"], "-not-present-")
				assert.NotEqual(t, refreshRes["refreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")

			case "none":
				if behaviour.clearedTokens == "none" {
					assert.Empty(t, refreshRes["frontToken"])
				}
			}

			if behaviour.setTokens != "cookies" && behaviour.clearedTokens != "cookies" {
				assert.Equal(t, refreshRes["sAccessToken"], "-not-present-")
				assert.Equal(t, refreshRes["accessTokenExpiry"], "-not-present-")
				assert.Equal(t, refreshRes["sRefreshToken"], "-not-present-")
				assert.Equal(t, refreshRes["refreshTokenExpiry"], "-not-present-")
			}
			if behaviour.setTokens != "headers" && behaviour.clearedTokens != "headers" {
				assert.Equal(t, refreshRes["accessTokenFromHeader"], "-not-present-")
				assert.Equal(t, refreshRes["refreshTokenFromHeader"], "-not-present-")
			}

		})

		t.Run(fmt.Sprintf("behaviour %v with invalid token", behaviour), func(t *testing.T) {
			configValue := supertokens.TypeInput{
				Supertokens: &supertokens.ConnectionInfo{
					ConnectionURI: "http://localhost:8080",
				},
				AppInfo: supertokens.AppInfo{
					APIDomain:     "api.supertokens.io",
					AppName:       "SuperTokens",
					WebsiteDomain: "supertokens.io",
				},
				RecipeList: []supertokens.Recipe{
					Init(nil),
					session.Init(&sessmodels.TypeInput{
						AntiCsrf: &viaToken,
						GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
							return sessmodels.TokenTransferMethod(behaviour.getTokenTransferMethodRes)
						},
					}),
				},
			}

			BeforeEach()
			unittesting.StartUpST("localhost", "8080")
			defer AfterEach()

			err := supertokens.Init(configValue)
			if err != nil {
				assert.NoError(t, err)
			}

			mux := http.NewServeMux()
			testServer := httptest.NewServer(supertokens.Middleware(mux))
			defer testServer.Close()
			setupRoutesForTest(t, mux)

			createInfo := map[string]string{
				"sRefreshToken": "invalid",
			}

			authMode := ""
			if behaviour.authCookie && behaviour.authHeader {
				authMode = "both"
			} else if behaviour.authCookie {
				authMode = "cookie"
			} else if behaviour.authHeader {
				authMode = "header"
			} else {
				authMode = "none"
			}

			refreshRes := refreshSession(t, testServer.URL, &behaviour.getTokenTransferMethodRes, authMode, createInfo, 401)

			if behaviour.output == "validateheader" {
				assert.Empty(t, refreshRes["accessTokenFromHeader"])
				assert.Empty(t, refreshRes["refreshTokenFromHeader"])
			} else if behaviour.output == "validatecookie" {
				assert.Empty(t, refreshRes["sAccessToken"])
				assert.Equal(t, refreshRes["accessTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
				assert.Empty(t, refreshRes["sRefreshToken"])
				assert.Equal(t, refreshRes["refreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
			}
		})
	}
}

func createNewSession(t *testing.T, baseURL string, authModeHeader *string, body map[string]interface{}, cookies []http.Cookie, headers map[string]string) map[string]string {
	req, err := http.NewRequest(http.MethodPost, baseURL+"/create", nil)
	assert.NoError(t, err)

	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if authModeHeader != nil {
		req.Header.Add("st-auth-mode", *authModeHeader)
	}

	createResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	result := unittesting.ExtractInfoFromResponseForAuthModeTests(createResp)

	bodyMap := map[string]interface{}{}
	err = json.NewDecoder(createResp.Body).Decode(&bodyMap)
	assert.NoError(t, err)

	for k, v := range bodyMap {
		result[k] = fmt.Sprint(v)
	}

	return result
}

func refreshSession(t *testing.T, baseURL string, authModeHeader *string, authMode string, info map[string]string, expectedStatus int) map[string]string {
	req, err := http.NewRequest(http.MethodPost, baseURL+"/auth/session/refresh", nil)
	assert.NoError(t, err)

	if authModeHeader != nil {
		req.Header.Add("st-auth-mode", *authModeHeader)
	}

	accessToken := ""
	if info["sAccessToken"] != "" && info["sAccessToken"] != "-not-present-" {
		accessToken = info["sAccessToken"]
	}
	if info["accessTokenFromHeader"] != "" && info["accessTokenFromHeader"] != "-not-present-" {
		accessToken = url.QueryEscape(info["accessTokenFromHeader"])
	}

	refreshToken := ""
	if info["sRefreshToken"] != "" && info["sRefreshToken"] != "-not-present-" {
		refreshToken = info["sRefreshToken"]
	}
	if info["refreshTokenFromHeader"] != "" && info["refreshTokenFromHeader"] != "-not-present-" {
		refreshToken = url.QueryEscape(info["refreshTokenFromHeader"])
	}

	if authMode == "cookie" || authMode == "both" {
		req.AddCookie(&http.Cookie{Name: "sAccessToken", Value: accessToken})
		req.AddCookie(&http.Cookie{Name: "sRefreshToken", Value: refreshToken})

		if info["antiCsrf"] != "" {
			req.Header.Set("anti-csrf", info["antiCsrf"])
		}
	}
	if authMode == "header" || authMode == "both" {
		refToken, err := url.QueryUnescape(refreshToken)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+refToken)
	}

	refreshResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, expectedStatus, refreshResp.StatusCode)

	result := unittesting.ExtractInfoFromResponseForAuthModeTests(refreshResp)

	return result
}

func testGetVerify(t *testing.T, baseURL string, info map[string]string, sessionRequired bool, expectedStatus int, authMode string) map[string]interface{} {
	endpoint := "/verify-optional"
	if sessionRequired {
		endpoint = "/verify"
	}

	return testGetUrl(t, baseURL, info, endpoint, expectedStatus, authMode, nil)
}

func testGetUrl(t *testing.T, baseURL string, info map[string]string, endpoint string, expectedStatus int, authMode string, authModeHeader interface{}) map[string]interface{} {
	accessToken := ""
	if info["sAccessToken"] != "" && info["sAccessToken"] != "-not-present-" {
		accessToken = info["sAccessToken"]
	}
	if info["accessTokenFromHeader"] != "" && info["accessTokenFromHeader"] != "-not-present-" {
		accessToken = url.QueryEscape(info["accessTokenFromHeader"])
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+endpoint, nil)
	assert.NoError(t, err)

	if authMode == "cookie" || authMode == "both" {
		req.Header.Add("Cookie", "sAccessToken="+accessToken)
	}
	if authMode == "header" || authMode == "both" {
		accToken, err := url.QueryUnescape(accessToken)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+accToken)
	}

	if info["antiCsrf"] != "" {
		req.Header.Set("anti-csrf", info["antiCsrf"])
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, resp.StatusCode, expectedStatus)

	result := map[string]interface{}{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	tokenInfo := unittesting.ExtractInfoFromResponseForAuthModeTests(resp)
	for k, v := range tokenInfo {
		result[k] = v
	}

	return result
}

func setupRoutesForTest(t *testing.T, mux *http.ServeMux) {
	mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {

		sessionContainer, err := session.CreateNewSession(r, w, "testuser", nil, nil, nil)
		assert.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       true,
			"sessionHandle": sessionContainer.GetHandle(),
		})
	})

	mux.HandleFunc("/update-payload", verifySession(true, func(w http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		sessionContainer.MergeIntoAccessTokenPayload(map[string]interface{}{"newValue": "test"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": true,
		})

	}))

	mux.HandleFunc("/verify", verifySession(true, func(w http.ResponseWriter, r *http.Request) {
		session := session.GetSessionFromRequestContext(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       true,
			"sessionHandle": session.GetHandle(),
			"sessionExists": true,
		})
	}))

	mux.HandleFunc("/verify-optional", verifySession(false, func(w http.ResponseWriter, r *http.Request) {
		session := session.GetSessionFromRequestContext(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if session == nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message":       true,
				"sessionHandle": nil,
				"sessionExists": false,
			})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message":       true,
				"sessionHandle": session.GetHandle(),
				"sessionExists": true,
			})
		}
	}))

	mux.HandleFunc("/logout", verifySession(true, func(w http.ResponseWriter, r *http.Request) {
		session := session.GetSessionFromRequestContext(r.Context())
		err := session.RevokeSession()
		assert.NoError(t, err)
	}))
}

func verifySession(sessionRequired bool, otherHandler http.HandlerFunc) http.HandlerFunc {
	return session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &sessionRequired,
	}, otherHandler)
}
