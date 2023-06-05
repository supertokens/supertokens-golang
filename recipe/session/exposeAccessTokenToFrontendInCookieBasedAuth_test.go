package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestThatShouldDefaultToFalse(t *testing.T) {
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
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	recipe, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	assert.False(t, recipe.Config.ExposeAccessTokenToFrontendInCookieBasedAuth)
}

func TestThatItAttachesTokensWithEnabled(t *testing.T) {
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
				ExposeAccessTokenToFrontendInCookieBasedAuth: true,
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	checkResponse(t, res, true)
}

func TestThatItAttachesTokensWithDisabled(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	checkResponse(t, res, false)
}

func TestShouldAttachTokensWhenEnabled(t *testing.T) {
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
				ExposeAccessTokenToFrontendInCookieBasedAuth: true,
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	appSub := "asdf"
	body := map[string]map[string]*string{
		"payload": {
			"sub":    nil,
			"appSub": &appSub,
		},
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		t.Error(err.Error())
	}
	res2, err2 := http.Post(testServer.URL+"/create", "application/json", bytes.NewBuffer(postBody))
	if err2 != nil {
		t.Error(err.Error())
	}

	checkResponse(t, res2, true)
}

func TestShouldAttachTokensWhenDisabled(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	appSub := "asdf"
	body := map[string]map[string]*string{
		"payload": {
			"appSub": &appSub,
		},
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		t.Error(err.Error())
	}
	res2, err2 := http.Post(testServer.URL+"/create", "application/json", bytes.NewBuffer(postBody))
	if err2 != nil {
		t.Error(err.Error())
	}

	checkResponse(t, res2, false)
}

func TestShouldAttachTokensAfterRefreshAndVerifyWhenEnabled(t *testing.T) {
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
				ExposeAccessTokenToFrontendInCookieBasedAuth: true,
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	info := unittesting.ExtractInfoFromResponse(res)

	req2, err2 := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	if err2 != nil {
		t.Error(err.Error())
	}

	req2.Header.Set("Cookie", "sRefreshToken="+info["sRefreshToken"])
	res2, err3 := http.DefaultClient.Do(req)
	assert.NoError(t, err3)
	assert.Equal(t, 200, res2.StatusCode)

	refreshInfo := unittesting.ExtractInfoFromResponse(res2)

	req3, err3 := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	if err3 != nil {
		t.Error(err.Error())
	}

	req3.Header.Set("Cookie", "sAccessToken="+refreshInfo["sAccessToken"])
	res3, err4 := http.DefaultClient.Do(req)
	assert.NoError(t, err4)
	assert.Equal(t, 200, res3.StatusCode)

	checkResponse(t, res3, true)
}

func TestShouldAttachTokensAfterRefreshAndVerifyWhenDisabled(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	info := unittesting.ExtractInfoFromResponse(res)

	req2, err2 := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	if err2 != nil {
		t.Error(err.Error())
	}

	req2.Header.Set("Cookie", "sRefreshToken="+info["sRefreshToken"])
	res2, err3 := http.DefaultClient.Do(req)
	assert.NoError(t, err3)
	assert.Equal(t, 200, res2.StatusCode)

	refreshInfo := unittesting.ExtractInfoFromResponse(res2)

	req3, err3 := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	if err3 != nil {
		t.Error(err.Error())
	}

	req3.Header.Set("Cookie", "sAccessToken="+refreshInfo["sAccessToken"])
	res3, err4 := http.DefaultClient.Do(req)
	assert.NoError(t, err4)
	assert.Equal(t, 200, res3.StatusCode)

	checkResponse(t, res3, false)
}

func TestShouldAttachTokensAfterRefreshWhenEnabled(t *testing.T) {
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
				ExposeAccessTokenToFrontendInCookieBasedAuth: true,
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	info := unittesting.ExtractInfoFromResponse(res)

	req2, err2 := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	if err2 != nil {
		t.Error(err.Error())
	}

	req2.Header.Set("Cookie", "sRefreshToken="+info["sRefreshToken"])
	res2, err3 := http.DefaultClient.Do(req)
	assert.NoError(t, err3)
	assert.Equal(t, 200, res2.StatusCode)

	checkResponse(t, res2, true)
}

func TestShouldAttachTokensAfterRefreshWhenDisabled(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	//create a newSession
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	info := unittesting.ExtractInfoFromResponse(res)

	req2, err2 := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	if err2 != nil {
		t.Error(err.Error())
	}

	req2.Header.Set("Cookie", "sRefreshToken="+info["sRefreshToken"])
	res2, err3 := http.DefaultClient.Do(req)
	assert.NoError(t, err3)
	assert.Equal(t, 200, res2.StatusCode)

	checkResponse(t, res2, false)
}

func TestThatRefreshTokenIsNotSentInHeadersWhenUsingCookies(t *testing.T) {
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
				ExposeAccessTokenToFrontendInCookieBasedAuth: true,
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	refreshTokenInHeader := res.Header.Get(refreshTokenHeaderKey)
	accessTokenInHeader := res.Header.Get(accessTokenHeaderKey)
	assert.NotEqual(t, accessTokenInHeader, "")
	assert.Equal(t, refreshTokenInHeader, "")
}

func checkResponse(t *testing.T, res *http.Response, exposed bool) {
	info := unittesting.ExtractInfoFromResponse(res)

	if exposed {
		assert.Equal(t, info["sAccessToken"], info["accessTokenFromHeader"])
	} else {
		assert.Equal(t, info["accessTokenFromHeader"], "")
		assert.NotEqual(t, info["sAccessToken"], "")
	}
}

func GetTestServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	checkDBTrue := true
	checkDBFalse := false

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		dataInBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		var result map[string]interface{}
		err = json.Unmarshal(dataInBytes, &result)

		var payload map[string]interface{}

		if result["payload"] != nil {
			payload = result["payload"].(map[string]interface{})
		}

		_, err2 := CreateNewSession(r, rw, "uniqueId", payload, map[string]interface{}{})

		if err2 != nil {
			http.Error(rw, fmt.Sprint(err2), 400)
		}
	})

	mux.HandleFunc("/verify", verifySession(true, &checkDBTrue, func(rw http.ResponseWriter, r *http.Request) {
		session := GetSessionFromRequestContext(r.Context())
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"message":       true,
			"sessionHandle": session.GetHandle(),
			"sessionExists": session != nil,
		})
	}))

	mux.HandleFunc("/verify-no-db", verifySession(true, &checkDBFalse, func(rw http.ResponseWriter, r *http.Request) {
		session := GetSessionFromRequestContext(r.Context())
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"message":       true,
			"sessionHandle": session.GetHandle(),
			"sessionExists": session != nil,
		})
	}))

	mux.HandleFunc("/verify-checkdb", verifySession(true, &checkDBTrue, func(rw http.ResponseWriter, r *http.Request) {
		session := GetSessionFromRequestContext(r.Context())
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"message":       true,
			"sessionHandle": session.GetHandle(),
			"sessionExists": session != nil,
		})
	}))

	mux.HandleFunc("/merge-into-payload", verifySession(true, nil, func(rw http.ResponseWriter, r *http.Request) {
		session := GetSessionFromRequestContext(r.Context())

		dataInBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		var result map[string]interface{}
		err = json.Unmarshal(dataInBytes, &result)

		err = session.MergeIntoAccessTokenPayload(result["payload"].(map[string]interface{}))
		assert.NoError(t, err)

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"message":       true,
			"sessionHandle": session.GetHandle(),
			"sessionExists": session != nil,
			"newPayload":    session.GetAccessTokenPayload(),
		})
	}))

	mux.HandleFunc("/revoke-session", verifySession(true, nil, func(rw http.ResponseWriter, r *http.Request) {
		session := GetSessionFromRequestContext(r.Context())
		err := session.RevokeSession()
		assert.NoError(t, err)

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"message":       true,
			"sessionHandle": session.GetHandle(),
		})
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	return testServer
}

func verifySession(sessionRequired bool, checkDatabase *bool, otherHandler http.HandlerFunc) http.HandlerFunc {
	return VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &sessionRequired,
		CheckDatabase:   checkDatabase,
	}, otherHandler)
}
