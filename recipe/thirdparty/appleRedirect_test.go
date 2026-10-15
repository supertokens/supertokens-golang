package thirdparty

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func setUpAppleRedirectTest(t *testing.T, appInfo supertokens.AppInfo, tpConfig *tpmodels.TypeInput) *httptest.Server {
	connectionURI := unittesting.StartUpST("localhost", "8080")
	appInfo.AppName = "SuperTokens"
	appInfo.APIDomain = "api.example.com"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
		},
		AppInfo:    appInfo,
		RecipeList: []supertokens.Recipe{Init(tpConfig)},
	})
	assert.NoError(t, err)

	return httptest.NewServer(supertokens.Middleware(http.NewServeMux()))
}

func encodeAppleState(rawState string) string {
	return base64.StdEncoding.EncodeToString([]byte(rawState))
}

func encodeAppleStateWithRedirect(frontendRedirectURI string) string {
	stateJSON, _ := json.Marshal(map[string]interface{}{
		"frontendRedirectURI": frontendRedirectURI,
	})
	return encodeAppleState(string(stateJSON))
}

// postAppleCallback posts the given form to the apple callback API without following redirects.
func postAppleCallback(t *testing.T, testServer *httptest.Server, form url.Values) (*http.Response, string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := client.PostForm(testServer.URL+"/auth/callback/apple", form)
	// A panic in the handler would close the connection, resulting in an error here
	if !assert.NoError(t, err) {
		return nil, ""
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	return res, string(body)
}

func assertAppleCallbackRejected(t *testing.T, testServer *httptest.Server, form url.Values, description string) {
	res, _ := postAppleCallback(t, testServer, form)
	if res == nil {
		t.Errorf("no response for: %s", description)
		return
	}
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, description)
	assert.Empty(t, res.Header.Get("Location"), description)
}

func appleCallbackForm(state string) url.Values {
	return url.Values{
		"state":    []string{state},
		"code":     []string{"testcode"},
		"id_token": []string{"testidtoken"},
	}
}

func TestAppleCallbackRedirectsToSameOrigin(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, nil)
	defer testServer.Close()

	state := encodeAppleStateWithRedirect("https://app.example.com/auth/callback/apple?existing=1")
	res, _ := postAppleCallback(t, testServer, appleCallbackForm(state))
	if !assert.NotNil(t, res) {
		return
	}
	assert.Equal(t, http.StatusSeeOther, res.StatusCode)

	location, err := url.Parse(res.Header.Get("Location"))
	assert.NoError(t, err)
	assert.Equal(t, "https", location.Scheme)
	assert.Equal(t, "app.example.com", location.Host)
	assert.Equal(t, "/auth/callback/apple", location.Path)
	assert.Equal(t, "1", location.Query().Get("existing"))
	assert.Equal(t, "testcode", location.Query().Get("code"))
	assert.Equal(t, "testidtoken", location.Query().Get("id_token"))
	assert.Equal(t, state, location.Query().Get("state"))
}

func TestAppleCallbackAcceptsSameOriginWithNonCanonicalConfiguredOrigin(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{Origin: "https://App.Example.com:443"}, nil)
	defer testServer.Close()

	for _, redirect := range []string{
		"https://app.example.com/auth/callback/apple",
		"https://APP.example.com:443/auth/callback/apple",
	} {
		res, _ := postAppleCallback(t, testServer, appleCallbackForm(encodeAppleStateWithRedirect(redirect)))
		if !assert.NotNil(t, res, redirect) {
			continue
		}
		assert.Equal(t, http.StatusSeeOther, res.StatusCode, redirect)
		assert.True(t, strings.HasPrefix(res.Header.Get("Location"), "https://"), redirect)
	}
}

func TestAppleCallbackRejectsRedirectToAnotherOrigin(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, nil)
	defer testServer.Close()

	for _, redirect := range []string{
		"https://evil.test/cb",
		"https://app.example.com@evil.test/cb",
		"https://app.example.com.evil.test/cb",
		"http://app.example.com/cb",
		"https://app.example.com:8443/cb",
		"https:evil.test",
		"//evil.test/cb",
		"/auth/callback/apple",
		"javascript:alert(1)",
		"",
	} {
		assertAppleCallbackRejected(t, testServer, appleCallbackForm(encodeAppleStateWithRedirect(redirect)), redirect)
	}
}

func TestAppleCallbackRejectsInvalidState(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, nil)
	defer testServer.Close()

	assertAppleCallbackRejected(t, testServer, url.Values{"code": []string{"testcode"}}, "missing state")

	for description, state := range map[string]string{
		"empty state":                      "",
		"non base64 state":                 "%%%not-base64%%%",
		"non JSON state":                   encodeAppleState("not json"),
		"JSON array state":                 encodeAppleState(`["https://app.example.com/cb"]`),
		"JSON null state":                  encodeAppleState("null"),
		"JSON string state":                encodeAppleState(`"https://app.example.com/cb"`),
		"missing frontendRedirectURI":      encodeAppleState(`{}`),
		"non string frontendRedirectURI":   encodeAppleState(`{"frontendRedirectURI": 1}`),
		"null frontendRedirectURI":         encodeAppleState(`{"frontendRedirectURI": null}`),
		"object frontendRedirectURI":       encodeAppleState(`{"frontendRedirectURI": {"a": "b"}}`),
		"unparseable frontendRedirectURI":  encodeAppleState(`{"frontendRedirectURI": "https://app.example.com:port/cb"}`),
		"control char frontendRedirectURI": encodeAppleState(`{"frontendRedirectURI": "https://app.example.com/\u0000"}`),
	} {
		assertAppleCallbackRejected(t, testServer, appleCallbackForm(state), description)
	}
}

func TestAppleCallbackHandlerRejectsNonStringStateWithoutPanicking(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, nil)
	defer testServer.Close()

	stInstance, err := supertokens.GetInstanceOrThrowError()
	assert.NoError(t, err)

	apiImpl := api.MakeAPIImplementation()
	for _, state := range []interface{}{nil, 1, map[string]interface{}{}, []interface{}{}} {
		rec := httptest.NewRecorder()
		options := tpmodels.APIOptions{
			Req:     httptest.NewRequest(http.MethodPost, "/auth/callback/apple", nil),
			Res:     rec,
			AppInfo: stInstance.AppInfo,
		}
		assert.NotPanics(t, func() {
			err = (*apiImpl.AppleRedirectHandlerPOST)(map[string]interface{}{"state": state}, options, &map[string]interface{}{})
		})
		assert.ErrorAs(t, err, &supertokens.BadInputError{})
		assert.Empty(t, rec.Header().Get("Location"))
	}
}
