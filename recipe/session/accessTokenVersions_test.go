package session

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestShouldCreateAV4Token(t *testing.T) {
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

	cookies := unittesting.ExtractInfoFromResponse(res)
	assert.True(t, cookies["accessTokenFromAny"] != "")
	assert.True(t, cookies["refreshTokenFromAny"] != "")
	assert.True(t, cookies["frontToken"] != "")

	parsedToken, parseErr := ParseJWTWithoutSignatureVerification(cookies["accessTokenFromAny"])
	if parseErr != nil {
		t.Error(parseErr.Error())
	}

	assert.Equal(t, parsedToken.Version, 4)
	bytes, err := base64.RawURLEncoding.DecodeString(parsedToken.Header)
	if err != nil {
		t.Error(err.Error())
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, reflect.TypeOf(result["kid"]).Kind(), reflect.String)
	assert.True(t, strings.HasPrefix(result["kid"].(string), "d-"))
}

func TestShouldCreateV4TokenSignedByStaticKeyIfSetInConfig(t *testing.T) {
	useDynamicKey := false
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
				UseDynamicAccessTokenSigningKey: &useDynamicKey,
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

	cookies := unittesting.ExtractInfoFromResponse(res)
	assert.True(t, cookies["accessTokenFromAny"] != "")
	assert.True(t, cookies["refreshTokenFromAny"] != "")
	assert.True(t, cookies["frontToken"] != "")

	parsedToken, parseErr := ParseJWTWithoutSignatureVerification(cookies["accessTokenFromAny"])
	if parseErr != nil {
		t.Error(parseErr.Error())
	}

	assert.Equal(t, parsedToken.Version, 4)
	bytes, err := base64.RawURLEncoding.DecodeString(parsedToken.Header)
	if err != nil {
		t.Error(err.Error())
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, reflect.TypeOf(result["kid"]).Kind(), reflect.String)
	assert.True(t, strings.HasPrefix(result["kid"].(string), "s-"))
}

func TestShouldThrowErrorWhenUsingProtectedProps(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	appSub := "asdf"
	body := map[string]map[string]*string{
		"payload": {
			"sub": &appSub,
		},
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		t.Error(err.Error())
	}
	res2, err2 := http.Post(testServer.URL+"/create", "application/json", bytes.NewBuffer(postBody))
	if err2 != nil {
		t.Error(err2.Error())
	}

	assert.Equal(t, 200, res2.StatusCode)
	cookies := unittesting.ExtractInfoFromResponse(res2)
	assert.False(t, cookies["accessTokenFromAny"] == "")
	assert.False(t, cookies["refreshTokenFromAny"] == "")
	assert.False(t, cookies["frontToken"] == "")

	parsedToken, err := ParseJWTWithoutSignatureVerification(cookies["accessTokenFromAny"])
	if err != nil {
		t.Error(err.Error())
	}

	assert.True(t, parsedToken.Payload["sub"] != "asdf")
}

func TestMergeIntoATShouldHelpMigratingV2TokenUsingProtectedProps(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("session")
	if err != nil {
		t.Error(err.Error())
	}

	querier.SetApiVersionForTests("2.18")
	response, err := querier.SendPostRequest("/recipe/session", map[string]interface{}{
		"userId":         "test-user-id",
		"enableAntiCsrf": false,
		"userDataInJWT": map[string]interface{}{
			"sub": "asdf",
		},
		"userDataInDatabase": map[string]interface{}{},
	})
	if err != nil {
		t.Error(err.Error())
	}

	responseByte, err := json.Marshal(response)
	if err != nil {
		t.Error(err.Error())
	}
	var resp sessmodels.CreateOrRefreshAPIResponse
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		t.Error(err.Error())
	}

	legacyAccessToken := resp.AccessToken.Token
	legacyRefreshToken := resp.RefreshToken.Token

	querier.SetApiVersionForTests("")

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
	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/merge-into-payload", bytes.NewBuffer(postBody))
	req.Header.Set("Authorization", "Bearer "+legacyAccessToken)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	mergeCookies := unittesting.ExtractInfoFromResponse(res)
	assert.True(t, mergeCookies["accessTokenFromAny"] != "")
	assert.True(t, mergeCookies["refreshTokenFromAny"] == "")
	assert.True(t, mergeCookies["frontToken"] != "")

	parsedTokenAfterMerge, err := ParseJWTWithoutSignatureVerification(mergeCookies["accessTokenFromAny"])
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, parsedTokenAfterMerge.Version, 2)

	req2, err2 := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	if err2 != nil {
		t.Error(err.Error())
	}

	req2.Header.Set("Authorization", "Bearer "+legacyRefreshToken)
	res2, err3 := http.DefaultClient.Do(req2)
	assert.NoError(t, err3)
	assert.Equal(t, 200, res2.StatusCode)

	cookiesAfterRefresh := unittesting.ExtractInfoFromResponse(res2)
	assert.True(t, cookiesAfterRefresh["accessTokenFromAny"] != "")
	assert.True(t, cookiesAfterRefresh["refreshTokenFromAny"] != "")
	assert.True(t, cookiesAfterRefresh["frontToken"] != "")

	parsedAtAfterRefresh, err := ParseJWTWithoutSignatureVerification(cookiesAfterRefresh["accessTokenFromAny"])
	if err2 != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, parsedAtAfterRefresh.Version, 4)
	assert.Equal(t, parsedAtAfterRefresh.Payload["sub"], "test-user-id")
	assert.Equal(t, parsedAtAfterRefresh.Payload["appSub"], "asdf")

	bytes, err := base64.RawURLEncoding.DecodeString(parsedAtAfterRefresh.Header)
	if err != nil {
		t.Error(err.Error())
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, reflect.TypeOf(result["kid"]).Kind(), reflect.String)
	assert.True(t, strings.HasPrefix(result["kid"].(string), "d-"))
}

func TestShouldHelpMigratingV2TokenUsingProtectedPropsWhenCalledUsingSessionHandle(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("session")
	if err != nil {
		t.Error(err.Error())
	}

	querier.SetApiVersionForTests("2.18")
	response, err := querier.SendPostRequest("/recipe/session", map[string]interface{}{
		"userId":         "test-user-id",
		"enableAntiCsrf": false,
		"userDataInJWT": map[string]interface{}{
			"sub": "asdf",
		},
		"userDataInDatabase": map[string]interface{}{},
	})
	if err != nil {
		t.Error(err.Error())
	}

	responseByte, err := json.Marshal(response)
	if err != nil {
		t.Error(err.Error())
	}
	var resp sessmodels.CreateOrRefreshAPIResponse
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		t.Error(err.Error())
	}

	legacyRefreshToken := resp.RefreshToken.Token

	querier.SetApiVersionForTests("")

	_, err = MergeIntoAccessTokenPayload(resp.Session.Handle, map[string]interface{}{
		"sub":    nil,
		"appSub": "asdf",
	})
	if err != nil {
		t.Error(err.Error())
	}

	req2, err2 := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	if err2 != nil {
		t.Error(err.Error())
	}

	req2.Header.Set("Authorization", "Bearer "+legacyRefreshToken)
	res2, err3 := http.DefaultClient.Do(req2)
	assert.NoError(t, err3)
	assert.Equal(t, 200, res2.StatusCode)

	cookiesAfterRefresh := unittesting.ExtractInfoFromResponse(res2)
	assert.True(t, cookiesAfterRefresh["accessTokenFromAny"] != "")
	assert.True(t, cookiesAfterRefresh["refreshTokenFromAny"] != "")
	assert.True(t, cookiesAfterRefresh["frontToken"] != "")

	parsedAtAfterRefresh, err := ParseJWTWithoutSignatureVerification(cookiesAfterRefresh["accessTokenFromAny"])
	if err2 != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, parsedAtAfterRefresh.Version, 4)
	assert.Equal(t, parsedAtAfterRefresh.Payload["sub"], "test-user-id")
	assert.Equal(t, parsedAtAfterRefresh.Payload["appSub"], "asdf")

	bytes, err := base64.RawURLEncoding.DecodeString(parsedAtAfterRefresh.Header)
	if err != nil {
		t.Error(err.Error())
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, reflect.TypeOf(result["kid"]).Kind(), reflect.String)
	assert.True(t, strings.HasPrefix(result["kid"].(string), "d-"))
}

func TestVerifyShouldValidateV2Tokens(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("session")
	if err != nil {
		t.Error(err.Error())
	}

	querier.SetApiVersionForTests("2.18")
	response, err := querier.SendPostRequest("/recipe/session", map[string]interface{}{
		"userId":         "test-user-id",
		"enableAntiCsrf": false,
		"userDataInJWT": map[string]interface{}{
			"sub": "asdf",
		},
		"userDataInDatabase": map[string]interface{}{},
	})
	if err != nil {
		t.Error(err.Error())
	}

	responseByte, err := json.Marshal(response)
	if err != nil {
		t.Error(err.Error())
	}
	var resp sessmodels.CreateOrRefreshAPIResponse
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		t.Error(err.Error())
	}

	legacyToken := resp.AccessToken.Token

	querier.SetApiVersionForTests("")

	req3, err3 := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	if err3 != nil {
		t.Error(err.Error())
	}

	req3.Header.Set("Authorization", "Bearer "+legacyToken)
	res3, err4 := http.DefaultClient.Do(req3)
	assert.NoError(t, err4)
	assert.Equal(t, 200, res3.StatusCode)

	dataInBytes, err := io.ReadAll(res3.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res3.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, data, map[string]interface{}{
		"message":       true,
		"sessionExists": true,
		"sessionHandle": resp.Session.Handle,
	})
}

func TestVerifyShouldValidateV2TokensWithCheckDatabaseEnabled(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("session")
	if err != nil {
		t.Error(err.Error())
	}

	querier.SetApiVersionForTests("2.18")
	response, err := querier.SendPostRequest("/recipe/session", map[string]interface{}{
		"userId":         "test-user-id",
		"enableAntiCsrf": false,
		"userDataInJWT": map[string]interface{}{
			"sub": "asdf",
		},
		"userDataInDatabase": map[string]interface{}{},
	})
	if err != nil {
		t.Error(err.Error())
	}

	responseByte, err := json.Marshal(response)
	if err != nil {
		t.Error(err.Error())
	}
	var resp sessmodels.CreateOrRefreshAPIResponse
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		t.Error(err.Error())
	}

	legacyToken := resp.AccessToken.Token

	querier.SetApiVersionForTests("")

	req3, err3 := http.NewRequest(http.MethodGet, testServer.URL+"/revoke-session", nil)
	if err3 != nil {
		t.Error(err.Error())
	}

	req3.Header.Set("Authorization", "Bearer "+legacyToken)
	res3, err4 := http.DefaultClient.Do(req3)
	assert.NoError(t, err4)
	assert.Equal(t, 200, res3.StatusCode)

	req3, err3 = http.NewRequest(http.MethodGet, testServer.URL+"/verify-checkdb", nil)
	if err3 != nil {
		t.Error(err.Error())
	}

	req3.Header.Set("Authorization", "Bearer "+legacyToken)
	res3, err4 = http.DefaultClient.Do(req3)
	assert.NoError(t, err4)
	assert.Equal(t, 401, res3.StatusCode)

	req3, err3 = http.NewRequest(http.MethodGet, testServer.URL+"/verify-no-db", nil)
	if err3 != nil {
		t.Error(err.Error())
	}

	req3.Header.Set("Authorization", "Bearer "+legacyToken)
	res3, err4 = http.DefaultClient.Do(req3)
	assert.NoError(t, err4)
	assert.Equal(t, 200, res3.StatusCode)

	dataInBytes, err := io.ReadAll(res3.Body)
	if err != nil {
		t.Error(err.Error())
	}
	res3.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(dataInBytes, &data)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, data, map[string]interface{}{
		"message":       true,
		"sessionExists": true,
		"sessionHandle": resp.Session.Handle,
	})
}

func TestVerifyShouldValidateV3TokensWithCheckDatabaseEnabled(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	cookies := unittesting.ExtractInfoFromResponse(res)
	assert.True(t, cookies["accessTokenFromAny"] != "")
	assert.True(t, cookies["refreshTokenFromAny"] != "")
	assert.True(t, cookies["frontToken"] != "")

	legacyToken := cookies["accessTokenFromAny"]

	req3, err3 := http.NewRequest(http.MethodGet, testServer.URL+"/revoke-session", nil)
	if err3 != nil {
		t.Error(err.Error())
	}

	req3.Header.Set("Authorization", "Bearer "+legacyToken)
	res3, err4 := http.DefaultClient.Do(req3)
	assert.NoError(t, err4)
	assert.Equal(t, 200, res3.StatusCode)

	req3, err3 = http.NewRequest(http.MethodGet, testServer.URL+"/verify-checkdb", nil)
	if err3 != nil {
		t.Error(err.Error())
	}

	req3.Header.Set("Authorization", "Bearer "+legacyToken)
	res3, err4 = http.DefaultClient.Do(req3)
	assert.NoError(t, err4)
	assert.Equal(t, 401, res3.StatusCode)

	req3, err3 = http.NewRequest(http.MethodGet, testServer.URL+"/verify-no-db", nil)
	if err3 != nil {
		t.Error(err.Error())
	}

	req3.Header.Set("Authorization", "Bearer "+legacyToken)
	res3, err4 = http.DefaultClient.Do(req3)
	assert.NoError(t, err4)
	assert.Equal(t, 200, res3.StatusCode)
}

func TestVerifyShouldNotValidateTokenSignedByStaticKeyIfNotSetInConfig(t *testing.T) {
	useDynamicKey := false
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
				UseDynamicAccessTokenSigningKey: &useDynamicKey,
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
	checkDBTrue := true

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

		_, err2 := CreateNewSession(r, rw, "public", "uniqueId", payload, map[string]interface{}{})

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

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	cookies := unittesting.ExtractInfoFromResponse(res)
	assert.True(t, cookies["accessTokenFromAny"] != "")
	assert.True(t, cookies["refreshTokenFromAny"] != "")
	assert.True(t, cookies["frontToken"] != "")

	resetAll()

	useDynamicKey = true
	configValue = supertokens.TypeInput{
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
				UseDynamicAccessTokenSigningKey: &useDynamicKey,
			}),
		},
	}
	err = supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux = http.NewServeMux()
	checkDBTrue = true

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

		_, err2 := CreateNewSession(r, rw, "public", "uniqueId", payload, map[string]interface{}{})

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

	testServer = httptest.NewServer(supertokens.Middleware(mux))

	req3, err3 := http.NewRequest(http.MethodGet, testServer.URL+"/verify", nil)
	if err3 != nil {
		t.Error(err.Error())
	}

	req3.Header.Set("Authorization", "Bearer "+cookies["accessTokenFromAny"])
	res3, err4 := http.DefaultClient.Do(req3)
	assert.NoError(t, err4)
	assert.Equal(t, 401, res3.StatusCode)
}

func TestShouldRefreshLegacySessionsToNewVersion(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("session")
	if err != nil {
		t.Error(err.Error())
	}

	querier.SetApiVersionForTests("2.18")
	response, err := querier.SendPostRequest("/recipe/session", map[string]interface{}{
		"userId":             "test-user-id",
		"enableAntiCsrf":     false,
		"userDataInJWT":      map[string]interface{}{},
		"userDataInDatabase": map[string]interface{}{},
	})
	if err != nil {
		t.Error(err.Error())
	}

	responseByte, err := json.Marshal(response)
	if err != nil {
		t.Error(err.Error())
	}
	var resp sessmodels.CreateOrRefreshAPIResponse
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		t.Error(err.Error())
	}

	legacyRefreshToken := resp.RefreshToken.Token
	querier.SetApiVersionForTests("")

	req2, err2 := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	if err2 != nil {
		t.Error(err.Error())
	}

	req2.Header.Set("Authorization", "Bearer "+legacyRefreshToken)
	res2, err3 := http.DefaultClient.Do(req2)
	assert.NoError(t, err3)
	assert.Equal(t, 200, res2.StatusCode)

	cookiesAfterRefresh := unittesting.ExtractInfoFromResponse(res2)
	assert.True(t, cookiesAfterRefresh["accessTokenFromAny"] != "")
	assert.True(t, cookiesAfterRefresh["refreshTokenFromAny"] != "")
	assert.True(t, cookiesAfterRefresh["frontToken"] != "")

	parsedAtAfterRefresh, err := ParseJWTWithoutSignatureVerification(cookiesAfterRefresh["accessTokenFromAny"])
	if err2 != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, parsedAtAfterRefresh.Version, 4)
}

func TestShouldThrowWhenRefreshInLegacySessionsWithProtectedProp(t *testing.T) {
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("session")
	if err != nil {
		t.Error(err.Error())
	}

	querier.SetApiVersionForTests("2.18")
	response, err := querier.SendPostRequest("/recipe/session", map[string]interface{}{
		"userId":         "test-user-id",
		"enableAntiCsrf": false,
		"userDataInJWT": map[string]interface{}{
			"sub": "asdf",
		},
		"userDataInDatabase": map[string]interface{}{},
	})
	if err != nil {
		t.Error(err.Error())
	}

	responseByte, err := json.Marshal(response)
	if err != nil {
		t.Error(err.Error())
	}
	var resp sessmodels.CreateOrRefreshAPIResponse
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		t.Error(err.Error())
	}

	legacyRefreshToken := resp.RefreshToken.Token
	querier.SetApiVersionForTests("")

	req2, err2 := http.NewRequest(http.MethodPost, testServer.URL+"/auth/session/refresh", nil)
	if err2 != nil {
		t.Error(err.Error())
	}

	req2.Header.Set("Authorization", "Bearer "+legacyRefreshToken)
	res2, err3 := http.DefaultClient.Do(req2)
	assert.NoError(t, err3)
	assert.Equal(t, 401, res2.StatusCode)

	cookiesAfterRefresh := unittesting.ExtractInfoFromResponse(res2)
	assert.True(t, cookiesAfterRefresh["accessTokenFromAny"] == "")
	assert.True(t, cookiesAfterRefresh["refreshTokenFromAny"] == "")
	assert.True(t, cookiesAfterRefresh["frontToken"] == "remove")
}

/**
We want to make sure that for access token claims that can be null, the SDK does not fail access token validation if the
core does not send them as part of the payload.

For this we verify that validation passes when the keys are nil, empty or a different type

For now this test checks for:
- antiCsrfToken
- parentRefreshTokenHash1

But this test should be updated to include any keys that the core considers optional in the payload (i.e either it sends
JSON null or skips them entirely)
*/
func TestValidationLogicWithKeysThatCanUseJSONNullValuesInClaims(t *testing.T) {
	version3 := 3

	payloadv3 := map[string]interface{}{
		"sessionHandle":     "",
		"sub":               "",
		"refreshTokenHash1": "",
		"exp":               float64(0),
		"iat":               float64(0),
	}

	err := ValidateAccessTokenStructure(payloadv3, version3)
	assert.NoError(t, err)

	payloadv3 = map[string]interface{}{
		"sessionHandle":           "",
		"sub":                     "",
		"refreshTokenHash1":       "",
		"exp":                     float64(0),
		"iat":                     float64(0),
		"parentRefreshTokenHash1": nil,
		"antiCsrfToken":           nil,
	}

	err = ValidateAccessTokenStructure(payloadv3, version3)
	assert.NoError(t, err)

	payloadv3 = map[string]interface{}{
		"sessionHandle":           "",
		"sub":                     "",
		"refreshTokenHash1":       "",
		"exp":                     float64(0),
		"iat":                     float64(0),
		"parentRefreshTokenHash1": "",
		"antiCsrfToken":           "",
	}

	err = ValidateAccessTokenStructure(payloadv3, version3)
	assert.NoError(t, err)

	payloadv3 = map[string]interface{}{
		"sessionHandle":           "",
		"sub":                     "",
		"refreshTokenHash1":       "",
		"exp":                     float64(0),
		"iat":                     float64(0),
		"parentRefreshTokenHash1": 1,
		"antiCsrfToken":           1,
	}

	err = ValidateAccessTokenStructure(payloadv3, version3)
	assert.NoError(t, err)
}
