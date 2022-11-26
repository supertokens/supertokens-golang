package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestRevokingSessionDuringRefreshWithRevokeSession(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
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
			Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
				Override: &sessmodels.OverrideStruct{
					APIs: func(originalImplementation sessmodels.APIInterface) sessmodels.APIInterface {
						oRefreshPOST := *originalImplementation.RefreshPOST
						refreshPost := func(options sessmodels.APIOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							sessionContainer, err := oRefreshPOST(options, userContext)
							if err != nil {
								return sessionContainer, err
							}
							err = sessionContainer.RevokeSession()
							if err != nil {
								return sessionContainer, err
							}
							return sessionContainer, nil
						}
						*originalImplementation.RefreshPOST = refreshPost
						return originalImplementation
					},
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(rw, "user", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	assert.NotEmpty(t, cookieData["sAccessToken"])
	assert.NotEmpty(t, cookieData["antiCsrf"])
	assert.NotEmpty(t, cookieData["idRefreshTokenFromHeader"])
	assert.NotEmpty(t, cookieData["sRefreshToken"])

	req, err = http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res, err = http.DefaultClient.Do(req)
	cookieData2 := unittesting.ExtractInfoFromResponse(res)
	assert.NoError(t, err)

	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, cookieData2["accessTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData2["refreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData2["idRefreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData2["accessToken"], "")
	assert.Equal(t, cookieData2["refreshToken"], "")
	assert.Equal(t, cookieData2["idRefreshTokenFromCookie"], "")
	assert.Equal(t, cookieData2["idRefreshTokenFromHeader"], "remove")
	assert.Greater(t, len(cookieData2["frontToken"]), 1)
}

func TestRevokingSessionDuringRefreshWithRevokeSessionAndSend401(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
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
			Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
				Override: &sessmodels.OverrideStruct{
					APIs: func(originalImplementation sessmodels.APIInterface) sessmodels.APIInterface {
						oRefreshPOST := *originalImplementation.RefreshPOST
						refreshPost := func(options sessmodels.APIOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							sessionContainer, err := oRefreshPOST(options, userContext)
							if err != nil {
								return sessionContainer, err
							}
							err = sessionContainer.RevokeSession()
							if err != nil {
								return sessionContainer, err
							}
							options.Res.Header().Add("Content-type", "application/json")
							options.Res.WriteHeader(401)
							options.Res.Write([]byte("{}"))
							return sessionContainer, nil
						}
						*originalImplementation.RefreshPOST = refreshPost
						return originalImplementation
					},
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(rw, "user", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	assert.NotEmpty(t, cookieData["sAccessToken"])
	assert.NotEmpty(t, cookieData["antiCsrf"])
	assert.NotEmpty(t, cookieData["idRefreshTokenFromHeader"])
	assert.NotEmpty(t, cookieData["sRefreshToken"])

	req, err = http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res, err = http.DefaultClient.Do(req)
	cookieData2 := unittesting.ExtractInfoFromResponse(res)
	assert.NoError(t, err)

	assert.Equal(t, res.StatusCode, 401)
	assert.Equal(t, cookieData2["accessTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData2["refreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData2["idRefreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData2["accessToken"], "")
	assert.Equal(t, cookieData2["refreshToken"], "")
	assert.Equal(t, cookieData2["idRefreshTokenFromCookie"], "")
	assert.Equal(t, cookieData2["idRefreshTokenFromHeader"], "remove")
	assert.Greater(t, len(cookieData2["frontToken"]), 1)
}

func TestRevokingSessionDuringRefreshWithThrowingUnauthorizedError(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
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
			Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
				Override: &sessmodels.OverrideStruct{
					APIs: func(originalImplementation sessmodels.APIInterface) sessmodels.APIInterface {
						oRefreshPOST := *originalImplementation.RefreshPOST
						refreshPost := func(options sessmodels.APIOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							sessionContainer, err := oRefreshPOST(options, userContext)
							if err != nil {
								return sessionContainer, err
							}
							return nil, errors.UnauthorizedError{
								Msg: "Unauthorized",
							}
						}
						*originalImplementation.RefreshPOST = refreshPost
						return originalImplementation
					},
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(rw, "user", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	assert.NotEmpty(t, cookieData["sAccessToken"])
	assert.NotEmpty(t, cookieData["antiCsrf"])
	assert.NotEmpty(t, cookieData["idRefreshTokenFromHeader"])
	assert.NotEmpty(t, cookieData["sRefreshToken"])

	req, err = http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res, err = http.DefaultClient.Do(req)
	cookieData2 := unittesting.ExtractInfoFromResponse(res)
	assert.NoError(t, err)

	assert.Equal(t, res.StatusCode, 401)
	assert.Equal(t, cookieData2["accessTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData2["refreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData2["idRefreshTokenExpiry"], "Thu, 01 Jan 1970 00:00:00 GMT")
	assert.Equal(t, cookieData2["accessToken"], "")
	assert.Equal(t, cookieData2["refreshToken"], "")
	assert.Equal(t, cookieData2["idRefreshTokenFromCookie"], "")
	assert.Equal(t, cookieData2["idRefreshTokenFromHeader"], "remove")
	assert.Greater(t, len(cookieData2["frontToken"]), 1)
}

func TestRevokingSessionDuringRefreshFailsIfJustSending401(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
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
			Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
				Override: &sessmodels.OverrideStruct{
					APIs: func(originalImplementation sessmodels.APIInterface) sessmodels.APIInterface {
						oRefreshPOST := *originalImplementation.RefreshPOST
						refreshPost := func(options sessmodels.APIOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							sessionContainer, err := oRefreshPOST(options, userContext)
							if err != nil {
								return sessionContainer, err
							}
							options.Res.Header().Add("Content-type", "application/json")
							options.Res.WriteHeader(401)
							options.Res.Write([]byte("{}"))
							return sessionContainer, nil
						}
						*originalImplementation.RefreshPOST = refreshPost
						return originalImplementation
					},
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

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(rw, "user", map[string]interface{}{}, map[string]interface{}{})
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	assert.NotEmpty(t, cookieData["sAccessToken"])
	assert.NotEmpty(t, cookieData["antiCsrf"])
	assert.NotEmpty(t, cookieData["idRefreshTokenFromHeader"])
	assert.NotEmpty(t, cookieData["sRefreshToken"])

	req, err = http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	assert.NoError(t, err)
	req.Header.Add("Cookie", "sRefreshToken="+cookieData["sRefreshToken"]+";"+"sIdRefreshToken="+cookieData["sIdRefreshToken"])
	req.Header.Add("anti-csrf", cookieData["antiCsrf"])
	res, err = http.DefaultClient.Do(req)
	cookieData2 := unittesting.ExtractInfoFromResponse(res)
	assert.NoError(t, err)

	assert.Equal(t, res.StatusCode, 401)
	assert.NotEmpty(t, cookieData2["sAccessToken"])
	assert.NotEmpty(t, cookieData2["antiCsrf"])
	assert.NotEmpty(t, cookieData2["idRefreshTokenFromHeader"])
	assert.NotEmpty(t, cookieData2["sRefreshToken"])
}
