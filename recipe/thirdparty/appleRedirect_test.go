package thirdparty

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
		"myapp://callback",
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

func TestAppleCallbackCustomRedirectValidatorReplacesOriginCheck(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	validatedURIs := []string{}
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, &tpmodels.TypeInput{
		IsValidFrontendRedirectURI: func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
			assert.NotNil(t, req)
			validatedURIs = append(validatedURIs, frontendRedirectURI)
			return frontendRedirectURI == "https://other.example.com/auth/callback/apple", nil
		},
	})
	defer testServer.Close()

	// A different origin is allowed when the custom validator allows it
	state := encodeAppleStateWithRedirect("https://other.example.com/auth/callback/apple")
	res, _ := postAppleCallback(t, testServer, appleCallbackForm(state))
	if assert.NotNil(t, res) {
		assert.Equal(t, http.StatusSeeOther, res.StatusCode)
		location, err := url.Parse(res.Header.Get("Location"))
		assert.NoError(t, err)
		assert.Equal(t, "other.example.com", location.Host)
		assert.Equal(t, "testcode", location.Query().Get("code"))
	}

	// The configured origin is rejected when the custom validator rejects it
	assertAppleCallbackRejected(t, testServer, appleCallbackForm(encodeAppleStateWithRedirect("https://app.example.com/auth/callback/apple")), "same origin rejected by custom validator")

	assert.Equal(t, []string{
		"https://other.example.com/auth/callback/apple",
		"https://app.example.com/auth/callback/apple",
	}, validatedURIs)

	// Disallowed schemes are still rejected before the custom validator is called
	assertAppleCallbackRejected(t, testServer, appleCallbackForm(encodeAppleStateWithRedirect("javascript:alert(1)")), "javascript URL")
	assert.Len(t, validatedURIs, 2)
}

func TestAppleCallbackCustomRedirectValidatorAllowsCustomSchemes(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, &tpmodels.TypeInput{
		IsValidFrontendRedirectURI: func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
			return true, nil
		},
	})
	defer testServer.Close()

	for _, redirect := range []string{
		"myapp://callback",
		"myapp:callback",
		"com.example.app:/oauth2redirect",
		"intent://callback#Intent;package=com.example.app;scheme=signinwithapple;end",
	} {
		state := encodeAppleStateWithRedirect(redirect)
		res, _ := postAppleCallback(t, testServer, appleCallbackForm(state))
		if !assert.NotNil(t, res, redirect) {
			continue
		}
		assert.Equal(t, http.StatusSeeOther, res.StatusCode, redirect)

		// The posted fields are added as the query, keeping the rest of the target as it was
		locationHeader := res.Header.Get("Location")
		targetWithoutFragment, fragment, hasFragment := strings.Cut(redirect, "#")
		assert.True(t, strings.HasPrefix(locationHeader, targetWithoutFragment+"?"), "%s -> %s", redirect, locationHeader)
		if hasFragment {
			assert.True(t, strings.HasSuffix(locationHeader, "#"+fragment), "%s -> %s", redirect, locationHeader)
		}

		location, err := url.Parse(locationHeader)
		if !assert.NoError(t, err, redirect) {
			continue
		}
		target, _ := url.Parse(redirect)
		assert.Equal(t, target.Scheme, location.Scheme, redirect)
		assert.Equal(t, target.Host, location.Host, redirect)
		assert.Equal(t, target.Path, location.Path, redirect)
		assert.Equal(t, target.Opaque, location.Opaque, redirect)
		assert.Equal(t, target.Fragment, location.Fragment, redirect)
		assert.Equal(t, "testcode", location.Query().Get("code"), redirect)
		assert.Equal(t, state, location.Query().Get("state"), redirect)
	}
}

func TestAppleCallbackCustomRedirectValidatorReceivesTheRedirectTarget(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	var validatedURI string
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, &tpmodels.TypeInput{
		IsValidFrontendRedirectURI: func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
			validatedURI = frontendRedirectURI
			return true, nil
		},
	})
	defer testServer.Close()

	postedFields := appleCallbackForm("")
	for _, redirect := range []string{
		"HTTPS://other.example.com/auth/callback/apple",
		"https://other.example.com/auth/callback/apple?existing=1",
		"https://other.example.com/a b",
		"MyApp://callback",
		"myapp:callback",
		"intent://callback#Intent;package=com.example.app;scheme=signinwithapple;end",
	} {
		validatedURI = ""
		res, _ := postAppleCallback(t, testServer, appleCallbackForm(encodeAppleStateWithRedirect(redirect)))
		if !assert.NotNil(t, res, redirect) {
			continue
		}
		assert.Equal(t, http.StatusSeeOther, res.StatusCode, redirect)

		// The function must see exactly where the browser is sent, without the posted fields
		location, err := url.Parse(res.Header.Get("Location"))
		if !assert.NoError(t, err, redirect) {
			continue
		}
		query := location.Query()
		for field := range postedFields {
			query.Del(field)
		}
		location.RawQuery = query.Encode()
		assert.Equal(t, location.String(), validatedURI, redirect)
	}
}

func TestAppleCallbackCustomRedirectValidatorIsNotCalledForDisallowedTargets(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	validatedURIs := []string{}
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, &tpmodels.TypeInput{
		IsValidFrontendRedirectURI: func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
			validatedURIs = append(validatedURIs, frontendRedirectURI)
			return true, nil
		},
	})
	defer testServer.Close()

	for _, redirect := range []string{
		"javascript:alert(1)",
		"JavaScript:alert(1)",
		"data:text/html,<script>alert(1)</script>",
		"vbscript:msgbox(1)",
		"file:///etc/passwd",
		"blob:https://app.example.com/x",
		"about:blank",
		"filesystem:https://app.example.com/temporary/x",
		"/auth/callback/apple",
		"//evil.test/cb",
		"https:evil.test",
	} {
		assertAppleCallbackRejected(t, testServer, appleCallbackForm(encodeAppleStateWithRedirect(redirect)), redirect)
	}
	assert.Empty(t, validatedURIs)
}

func TestAppleCallbackAcceptsURLSafeBase64State(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, nil)
	defer testServer.Close()

	stateJSON, _ := json.Marshal(map[string]interface{}{
		"frontendRedirectURI": "https://app.example.com/auth/callback/apple",
		"stateForProvider":    "??>>??>>",
	})
	// Make sure the test covers the characters that differ between the two alphabets
	assert.True(t, strings.ContainsAny(base64.StdEncoding.EncodeToString(stateJSON), "+/"))

	for description, state := range map[string]string{
		"URL-safe alphabet with padding":    base64.URLEncoding.EncodeToString(stateJSON),
		"URL-safe alphabet without padding": base64.RawURLEncoding.EncodeToString(stateJSON),
		"standard alphabet without padding": base64.RawStdEncoding.EncodeToString(stateJSON),
	} {
		res, _ := postAppleCallback(t, testServer, appleCallbackForm(state))
		if !assert.NotNil(t, res, description) {
			continue
		}
		assert.Equal(t, http.StatusSeeOther, res.StatusCode, description)
		assert.True(t, strings.HasPrefix(res.Header.Get("Location"), "https://app.example.com/auth/callback/apple?"), description)
	}
}

func TestAppleCallbackReturnsErrorFromCustomRedirectValidator(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, &tpmodels.TypeInput{
		IsValidFrontendRedirectURI: func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
			return true, errors.New("custom validator failure")
		},
	})
	defer testServer.Close()

	res, body := postAppleCallback(t, testServer, appleCallbackForm(encodeAppleStateWithRedirect("https://app.example.com/auth/callback/apple")))
	if assert.NotNil(t, res) {
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Empty(t, res.Header.Get("Location"))
		assert.Contains(t, body, "custom validator failure")
	}
}

func TestAppleCallbackReplacesQueryParamsWithThePostedFields(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpAppleRedirectTest(t, supertokens.AppInfo{WebsiteDomain: "https://app.example.com"}, nil)
	defer testServer.Close()

	state := encodeAppleStateWithRedirect("https://app.example.com/auth/callback/apple?code=stale-code&state=stale-state&existing=1")
	res, _ := postAppleCallback(t, testServer, appleCallbackForm(state))
	if !assert.NotNil(t, res) {
		return
	}
	assert.Equal(t, http.StatusSeeOther, res.StatusCode)

	location, err := url.Parse(res.Header.Get("Location"))
	if !assert.NoError(t, err) {
		return
	}
	query := location.Query()
	assert.Equal(t, []string{"testcode"}, query["code"])
	assert.Equal(t, []string{state}, query["state"])
	assert.Equal(t, []string{"testidtoken"}, query["id_token"])
	assert.Equal(t, []string{"1"}, query["existing"])
}
