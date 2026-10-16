package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// These tests call the Apple callback handler directly, so they don't need a running core.

func appleStateWithRedirect(frontendRedirectURI string) string {
	stateJSON, _ := json.Marshal(map[string]interface{}{
		"frontendRedirectURI": frontendRedirectURI,
	})
	return base64.StdEncoding.EncodeToString(stateJSON)
}

// callAppleRedirectHandler calls the Apple callback handler with the given state and returns the
// Location header it set, or the error it returned.
func callAppleRedirectHandler(t *testing.T, websiteDomain string, isValidFrontendRedirectURI tpmodels.IsValidFrontendRedirectURIFunc, state interface{}) (string, error) {
	appInfo, err := supertokens.NormaliseInputAppInfoOrThrowError(supertokens.AppInfo{
		AppName:       "SuperTokens",
		APIDomain:     "https://api.example.com",
		WebsiteDomain: websiteDomain,
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	rec := httptest.NewRecorder()
	options := tpmodels.APIOptions{
		Config: tpmodels.TypeNormalisedInput{
			IsValidFrontendRedirectURI: isValidFrontendRedirectURI,
		},
		Req:     httptest.NewRequest(http.MethodPost, "/auth/callback/apple", nil),
		Res:     rec,
		AppInfo: appInfo,
	}
	formPostInfo := map[string]interface{}{
		"state":    state,
		"code":     "testcode",
		"id_token": "testidtoken",
	}

	apiImpl := MakeAPIImplementation()
	err = (*apiImpl.AppleRedirectHandlerPOST)(formPostInfo, options, &map[string]interface{}{})
	if err != nil {
		assert.Empty(t, rec.Header().Get("Location"))
		return "", err
	}
	assert.Equal(t, http.StatusSeeOther, rec.Code)
	return rec.Header().Get("Location"), nil
}

// withoutPostedFields removes the fields the handler adds to the redirect target's query.
func withoutPostedFields(t *testing.T, location string) string {
	parsedLocation, err := url.Parse(location)
	if !assert.NoError(t, err, location) {
		return ""
	}
	query := parsedLocation.Query()
	assert.Equal(t, "testcode", query.Get("code"), location)
	for _, field := range []string{"state", "code", "id_token"} {
		query.Del(field)
	}
	parsedLocation.RawQuery = query.Encode()
	return parsedLocation.String()
}

func TestRemoveDotSegments(t *testing.T) {
	for input, expected := range map[string]string{
		"https://app.example.com/a/../cb":                 "https://app.example.com/cb",
		"https://app.example.com/a/./b/../../cb":          "https://app.example.com/cb",
		"https://app.example.com/a/%2e%2e/cb":             "https://app.example.com/cb",
		"https://app.example.com/a/%2E%2e/cb":             "https://app.example.com/cb",
		"https://app.example.com/a/.%2e/cb":               "https://app.example.com/cb",
		"https://app.example.com/a/%2E./cb":               "https://app.example.com/cb",
		"https://app.example.com/a/%2e/cb":                "https://app.example.com/a/cb",
		"https://app.example.com/../../cb":                "https://app.example.com/cb",
		"https://app.example.com/a/b/..":                  "https://app.example.com/a/",
		"https://app.example.com/./cb?x=../y#../z":        "https://app.example.com/cb?x=../y#../z",
		"https://app.example.com/%41/%2e%2e/b%2fc":        "https://app.example.com/b%2fc",
		"https://app.example.com/a/%2e%2e%2e/b":           "https://app.example.com/a/%2e%2e%2e/b",
		"https://app.example.com/a/..b/c":                 "https://app.example.com/a/..b/c",
		"https://App.Example.com:443/a/../cb":             "https://App.Example.com:443/cb",
		"https://app.example.com":                         "https://app.example.com",
		"https://app.example.com/":                        "https://app.example.com/",
		"myapp://callback/a/../b":                         "myapp://callback/b",
		"myapp://callback/a/%2e%2e/b":                     "myapp://callback/b",
		"intent://callback/a/..#Intent;package=x.y;end":   "intent://callback/#Intent;package=x.y;end",
		"myapp:callback/../other":                         "myapp:callback/../other",
		"https://user@app.example.com/a/%2e%2e/b?c=d#e/f": "https://user@app.example.com/b?c=d#e/f",
		// Empty segments are kept, as in RFC 3986 section 5.2.4 and in browsers
		"https://app.example.com/..//evil.test/x":      "https://app.example.com//evil.test/x",
		"https://app.example.com/%2e%2e/%2e%2e//evil":  "https://app.example.com//evil",
		"https://app.example.com/a/b/../../../..//x":   "https://app.example.com//x",
		"https://app.example.com/%2e%2e//%2e%2e//evil": "https://app.example.com//evil",
		"https://app.example.com//a/../b":              "https://app.example.com//b",
		"https://app.example.com/a//../b":              "https://app.example.com/a/b",
		"https://app.example.com/a/b//..":              "https://app.example.com/a/b/",
		"myapp://host/..//evil/x":                      "myapp://host//evil/x",
		"myapp://user@host/a/..//b":                    "myapp://user@host//b",
		"https://app.example.com/a/b/.":                "https://app.example.com/a/b/",
		"https://app.example.com/a/b/%2E":              "https://app.example.com/a/b/",
		"https://app.example.com/.":                    "https://app.example.com/",
		"https://app.example.com/..":                   "https://app.example.com/",
		"https://app.example.com/.../x":                "https://app.example.com/.../x",
		"https://app.example.com/%C3%BC/../x":          "https://app.example.com/x",
		// A URL with an empty host keeps it
		"myapp:////evil.test/x": "myapp:///.//evil.test/x",
	} {
		checkRemoveDotSegments(t, input, expected)
	}
}

// urlStringOmitsMissingHost is false before Go 1.19, which serialises a URL without a host (e.g.
// com.example.app:/cb) with an empty one (com.example.app:///cb), because url.URL has no OmitHost
var urlStringOmitsMissingHost = func() bool {
	parsedURL, err := url.Parse("com.example.app:/cb")
	return err == nil && parsedURL.String() == "com.example.app:/cb"
}()

func TestRemoveDotSegmentsDoesNotAddAHost(t *testing.T) {
	for input, expected := range map[string]string{
		"com.example.app:/a/../oauth2redirect": "com.example.app:/oauth2redirect",
		"myapp:/.//evil.test/x":                "myapp:/.//evil.test/x",
		"myapp:/a/..//evil.test/x":             "myapp:/.//evil.test/x",
		"myapp:/%2e//evil.test/x":              "myapp:/.//evil.test/x",
		"myapp:/a/%2e%2e//evil.test/x":         "myapp:/.//evil.test/x",
		"myapp:/%2E%2E//evil.test/x":           "myapp:/.//evil.test/x",
		"myapp:/a/b/../..//evil.test/x":        "myapp:/.//evil.test/x",
		"myapp:/a/..//evil.test/x?y=z#f":       "myapp:/.//evil.test/x?y=z#f",
	} {
		parsedInput, err := url.Parse(input)
		if !assert.NoError(t, err, input) {
			continue
		}
		result := removeDotSegments(parsedInput)
		assert.NotContains(t, result.String(), "://evil.test", input)
		reparsed, err := url.Parse(result.String())
		if assert.NoError(t, err, input) {
			assert.Empty(t, reparsed.Host, input)
			assert.Nil(t, reparsed.User, input)
			assert.Equal(t, result.String(), removeDotSegments(reparsed).String(), input)
		}

		if !urlStringOmitsMissingHost {
			// Go 1.18 serialises these with an empty host (e.g. myapp:///.//evil.test/x), which is
			// checked above not to become a host either
			t.Logf("not comparing %s with %s: url.URL can't serialise a URL without a host before Go 1.19", input, expected)
			continue
		}
		checkRemoveDotSegments(t, input, expected)
	}
}

func checkRemoveDotSegments(t *testing.T, input string, expected string) {
	parsedInput, err := url.Parse(input)
	if !assert.NoError(t, err, input) {
		return
	}
	result := removeDotSegments(parsedInput)
	assert.Equal(t, expected, result.String(), input)
	// The input is not modified
	assert.Equal(t, input, parsedInput.String(), input)

	// Removing the dot segments again changes nothing, so the browser opens the URL as it is, and the
	// URL keeps its host
	reparsed, err := url.Parse(result.String())
	if assert.NoError(t, err, input) {
		assert.Equal(t, expected, removeDotSegments(reparsed).String(), input)
		assert.Equal(t, parsedInput.Host, reparsed.Host, input)
	}
}

func TestGetComparableOriginComparesHostsInIDNAASCIIForm(t *testing.T) {
	expected := comparableOrigin{scheme: "https", host: "xn--bcher-kva.example", port: "443"}
	for _, input := range []string{
		"https://bücher.example",
		"https://BÜCHER.example/cb",
		"HTTPS://Bücher.Example:443/cb",
		"https://xn--bcher-kva.example/cb",
		"https://XN--BCHER-KVA.example/cb",
		"https://b%C3%BCcher.example/cb",
		"https://ｂüｃｈｅｒ.example/cb",
	} {
		parsedInput, err := url.Parse(input)
		if !assert.NoError(t, err, input) {
			continue
		}
		origin, err := getComparableOrigin(parsedInput)
		if assert.NoError(t, err, input) {
			assert.Equal(t, expected, origin, input)
		}
	}

	for input, expectedOrigin := range map[string]comparableOrigin{
		"https://App.Example.com/cb":      {scheme: "https", host: "app.example.com", port: "443"},
		"http://localhost:3000/cb":        {scheme: "http", host: "localhost", port: "3000"},
		"http://127.0.0.1/cb":             {scheme: "http", host: "127.0.0.1", port: "80"},
		"http://[::1]:3000/cb":            {scheme: "http", host: "::1", port: "3000"},
		"https://my_app.example.com/cb":   {scheme: "https", host: "my_app.example.com", port: "443"},
		"https://r3---sn-x.example.com/":  {scheme: "https", host: "r3---sn-x.example.com", port: "443"},
		"https://bücher.example:8443/cb":  {scheme: "https", host: "xn--bcher-kva.example", port: "8443"},
		"https://app.example.com.evil.de": {scheme: "https", host: "app.example.com.evil.de", port: "443"},
	} {
		parsedInput, err := url.Parse(input)
		if !assert.NoError(t, err, input) {
			continue
		}
		origin, err := getComparableOrigin(parsedInput)
		if assert.NoError(t, err, input) {
			assert.Equal(t, expectedOrigin, origin, input)
		}
	}

	// Hosts that browsers can't convert to their ASCII form are rejected
	for _, input := range []string{
		"https://xn--zz.example/cb",
		"https://xn--a.example/cb",
	} {
		parsedInput, err := url.Parse(input)
		if !assert.NoError(t, err, input) {
			continue
		}
		_, err = getComparableOrigin(parsedInput)
		assert.Error(t, err, input)
	}
}

func TestAppleCallbackRemovesDotSegmentsBeforeTheOriginCheck(t *testing.T) {
	for redirect, expected := range map[string]string{
		"https://app.example.com/a/../cb":         "https://app.example.com/cb",
		"https://app.example.com/a/%2e%2e/cb":     "https://app.example.com/cb",
		"https://app.example.com/%2E%2e/%2e/cb":   "https://app.example.com/cb",
		"https://app.example.com/auth/./callback": "https://app.example.com/auth/callback",
	} {
		location, err := callAppleRedirectHandler(t, "https://app.example.com", nil, appleStateWithRedirect(redirect))
		if assert.NoError(t, err, redirect) {
			assert.Equal(t, expected, withoutPostedFields(t, location), redirect)
		}
	}

	for _, redirect := range []string{
		"https://evil.test/a/../cb",
		"https://evil.test/%2e%2e/app.example.com/cb",
		"https://app.example.com.evil.test/a/%2e%2e/cb",
	} {
		_, err := callAppleRedirectHandler(t, "https://app.example.com", nil, appleStateWithRedirect(redirect))
		assert.Equal(t, supertokens.BadInputError{Msg: "frontendRedirectURI in state does not match the configured origin"}, err, redirect)
	}
}

func TestAppleCallbackPassesTheURLWithoutDotSegmentsToIsValidFrontendRedirectURI(t *testing.T) {
	var validatedURIs []string
	isValidFrontendRedirectURI := func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
		validatedURIs = append(validatedURIs, frontendRedirectURI)
		return frontendRedirectURI == "https://other.example.com/cb" || frontendRedirectURI == "myapp://callback/cb", nil
	}

	for redirect, expected := range map[string]string{
		"https://other.example.com/a/../cb":     "https://other.example.com/cb",
		"https://other.example.com/a/%2e%2e/cb": "https://other.example.com/cb",
		"https://other.example.com/a/.%2E/cb":   "https://other.example.com/cb",
		"myapp://callback/a/../cb":              "myapp://callback/cb",
		"myapp://callback/a/b/%2e%2e/../cb":     "myapp://callback/cb",
	} {
		validatedURIs = nil
		location, err := callAppleRedirectHandler(t, "https://app.example.com", isValidFrontendRedirectURI, appleStateWithRedirect(redirect))
		if assert.NoError(t, err, redirect) {
			// The hook sees exactly the URL the browser is sent to, without the posted fields
			assert.Equal(t, []string{expected}, validatedURIs, redirect)
			assert.Equal(t, expected, withoutPostedFields(t, location), redirect)
		}
	}

	// A target that only matches before the dot segments are removed is rejected
	validatedURIs = nil
	_, err := callAppleRedirectHandler(t, "https://app.example.com", isValidFrontendRedirectURI, appleStateWithRedirect("https://other.example.com/cb/%2e%2e/admin"))
	assert.Equal(t, supertokens.BadInputError{Msg: "frontendRedirectURI in state is not allowed"}, err)
	assert.Equal(t, []string{"https://other.example.com/admin"}, validatedURIs)
}

func TestAppleCallbackMatchesAnIDNWebsiteDomain(t *testing.T) {
	for _, websiteDomain := range []string{"https://bücher.example", "https://xn--bcher-kva.example"} {
		for _, redirect := range []string{
			"https://bücher.example/auth/callback/apple",
			"https://BÜCHER.example/auth/callback/apple",
			"https://xn--bcher-kva.example/auth/callback/apple",
			"https://b%C3%BCcher.example:443/auth/callback/apple",
		} {
			location, err := callAppleRedirectHandler(t, websiteDomain, nil, appleStateWithRedirect(redirect))
			if !assert.NoError(t, err, websiteDomain+" "+redirect) {
				continue
			}
			parsedLocation, err := url.Parse(location)
			if assert.NoError(t, err, location) {
				assert.Equal(t, "xn--bcher-kva.example", parsedLocation.Hostname(), location)
				assert.Equal(t, "/auth/callback/apple", parsedLocation.Path, location)
				assert.Equal(t, "testcode", parsedLocation.Query().Get("code"), location)
			}
		}

		for _, redirect := range []string{
			"https://bucher.example/auth/callback/apple",
			"https://xn--bcher-kva.example.evil.test/auth/callback/apple",
			"http://bücher.example/auth/callback/apple",
		} {
			_, err := callAppleRedirectHandler(t, websiteDomain, nil, appleStateWithRedirect(redirect))
			assert.Equal(t, supertokens.BadInputError{Msg: "frontendRedirectURI in state does not match the configured origin"}, err, websiteDomain+" "+redirect)
		}

		_, err := callAppleRedirectHandler(t, websiteDomain, nil, appleStateWithRedirect("https://xn--zz.example/auth/callback/apple"))
		assert.Equal(t, supertokens.BadInputError{Msg: "frontendRedirectURI in state must be an absolute http or https URL"}, err)
	}
}

func TestAppleCallbackErrorMessages(t *testing.T) {
	alwaysValid := func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
		return true, nil
	}
	neverValid := func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
		return false, nil
	}

	encode := func(raw string) string {
		return base64.StdEncoding.EncodeToString([]byte(raw))
	}

	stateErrors := []struct {
		state    interface{}
		expected string
	}{
		{nil, "Please provide the state in the request body"},
		{"", "Please provide the state in the request body"},
		{1, "Please provide the state in the request body"},
		{"%%%not-base64%%%", "state must be base64 encoded"},
		{encode("not json"), "state must be a base64 encoded JSON object"},
		{encode("null"), "state must be a base64 encoded JSON object"},
		{encode(`["https://app.example.com/cb"]`), "state must be a base64 encoded JSON object"},
		{encode(`{}`), "state must contain frontendRedirectURI as a string"},
		{encode(`{"frontendRedirectURI": 1}`), "state must contain frontendRedirectURI as a string"},
	}
	for _, isValidFrontendRedirectURI := range []tpmodels.IsValidFrontendRedirectURIFunc{nil, alwaysValid} {
		for _, testCase := range stateErrors {
			_, err := callAppleRedirectHandler(t, "https://app.example.com", isValidFrontendRedirectURI, testCase.state)
			assert.Equal(t, supertokens.BadInputError{Msg: testCase.expected}, err, testCase.state)
		}
	}

	// Without IsValidFrontendRedirectURI
	for redirect, expected := range map[string]string{
		"https://app.example.com:port/cb": "frontendRedirectURI in state must be an absolute http or https URL",
		"/auth/callback/apple":            "frontendRedirectURI in state must be an absolute http or https URL",
		"myapp://callback":                "frontendRedirectURI in state must be an absolute http or https URL",
		"https:app.example.com":           "frontendRedirectURI in state must be an absolute http or https URL",
		"https://evil.test/cb":            "frontendRedirectURI in state does not match the configured origin",
		"http://app.example.com/cb":       "frontendRedirectURI in state does not match the configured origin",
	} {
		_, err := callAppleRedirectHandler(t, "https://app.example.com", nil, appleStateWithRedirect(redirect))
		assert.Equal(t, supertokens.BadInputError{Msg: expected}, err, redirect)
	}

	// With IsValidFrontendRedirectURI
	for redirect, expected := range map[string]string{
		"https://app.example.com:port/cb": "frontendRedirectURI in state must be an absolute URL with an allowed scheme",
		"/auth/callback/apple":            "frontendRedirectURI in state must be an absolute URL with an allowed scheme",
		"javascript:alert(1)":             "frontendRedirectURI in state must be an absolute URL with an allowed scheme",
		"https:app.example.com":           "frontendRedirectURI in state must be an absolute URL with an allowed scheme",
	} {
		_, err := callAppleRedirectHandler(t, "https://app.example.com", alwaysValid, appleStateWithRedirect(redirect))
		assert.Equal(t, supertokens.BadInputError{Msg: expected}, err, redirect)
	}
	_, err := callAppleRedirectHandler(t, "https://app.example.com", neverValid, appleStateWithRedirect("https://app.example.com/cb"))
	assert.Equal(t, supertokens.BadInputError{Msg: "frontendRedirectURI in state is not allowed"}, err)
}

func TestAppleCallbackDoesNotAddAHostToCustomSchemeTargets(t *testing.T) {
	var validatedURIs []string
	isValidFrontendRedirectURI := func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
		validatedURIs = append(validatedURIs, frontendRedirectURI)
		// A hook that only allows the evil.test host must not see it in these URLs
		parsedURI, err := url.Parse(frontendRedirectURI)
		return err == nil && parsedURI.Host != "evil.test", nil
	}

	for _, redirect := range []string{
		"myapp:/.//evil.test/x",
		"myapp:/a/..//evil.test/x",
		"myapp:/%2e//evil.test/x",
		"myapp:/a/%2e%2e//evil.test/x",
	} {
		validatedURIs = nil
		location, err := callAppleRedirectHandler(t, "https://app.example.com", isValidFrontendRedirectURI, appleStateWithRedirect(redirect))
		if !assert.NoError(t, err, redirect) || !assert.Len(t, validatedURIs, 1, redirect) {
			continue
		}
		for _, target := range []string{validatedURIs[0], location} {
			assert.NotContains(t, target, "://evil.test", redirect)
			parsedTarget, err := url.Parse(target)
			if assert.NoError(t, err, target) {
				assert.Empty(t, parsedTarget.Host, target)
				assert.Equal(t, "//evil.test/x", strings.TrimPrefix(parsedTarget.Path, "/."), target)
			}
		}
		if urlStringOmitsMissingHost {
			assert.Equal(t, "myapp:/.//evil.test/x", validatedURIs[0], redirect)
			assert.Equal(t, "myapp:/.//evil.test/x", withoutPostedFields(t, location), redirect)
		}
	}
}

func TestWithASCIIHost(t *testing.T) {
	for input, expected := range map[string]string{
		"https://bücher.example/cb":             "https://xn--bcher-kva.example/cb",
		"https://BÜCHER.example/cb":             "https://xn--bcher-kva.example/cb",
		"https://App.Bücher.example:8443/cb":    "https://app.xn--bcher-kva.example:8443/cb",
		"https://ｂüｃｈｅｒ.example/cb":             "https://xn--bcher-kva.example/cb",
		"https://b%C3%BCcher.example/cb":        "https://xn--bcher-kva.example/cb",
		"https://bücher.example:443/cb":         "https://xn--bcher-kva.example:443/cb",
		"https://bücher.example":                "https://xn--bcher-kva.example",
		"https://user@bücher.example/cb?q#f":    "https://user@xn--bcher-kva.example/cb?q#f",
		"http://ａｐｐ.example.com/cb":             "http://app.example.com/cb",
		"https://１２７.０.０.１/cb":                  "https://127.0.0.1/cb",
		"https://app。example。com/cb":            "https://app.example.com/cb",
		"https://xn--bcher-kva.example/cb":      "https://xn--bcher-kva.example/cb",
		"https://XN--BCHER-KVA.example/cb":      "https://XN--BCHER-KVA.example/cb",
		"https://App.Example.com/cb":            "https://App.Example.com/cb",
		"http://[::1]:3000/cb":                  "http://[::1]:3000/cb",
		"myapp://bücher/cb":                     "myapp://b%C3%BCcher/cb",
		"https://xn--zz.bücher.example/cb":      "https://xn--zz.b%C3%BCcher.example/cb",
		"https://evil.test／＠app.example.com/cb": "https://evil.test%EF%BC%8F%EF%BC%A0app.example.com/cb",
	} {
		parsedInput, err := url.Parse(input)
		if !assert.NoError(t, err, input) {
			continue
		}
		assert.Equal(t, expected, withASCIIHost(parsedInput).String(), input)
		// The input is not modified
		reparsedInput, _ := url.Parse(input)
		assert.Equal(t, reparsedInput, parsedInput, input)
	}
}

func TestAppleCallbackPassesTheASCIIFormOfInternationalisedHosts(t *testing.T) {
	var validatedURIs []string
	isValidFrontendRedirectURI := func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
		validatedURIs = append(validatedURIs, frontendRedirectURI)
		return true, nil
	}

	for redirect, expected := range map[string]string{
		"https://bücher.example/a/../cb":      "https://xn--bcher-kva.example/cb",
		"https://BÜCHER.example:8443/cb":      "https://xn--bcher-kva.example:8443/cb",
		"https://b%C3%BCcher.example/cb":      "https://xn--bcher-kva.example/cb",
		"https://ｂüｃｈｅｒ.example/cb":           "https://xn--bcher-kva.example/cb",
		"https://App.Example.com/cb":          "https://App.Example.com/cb",
		"https://app.example.com:443/cb":      "https://app.example.com:443/cb",
		"myapp://bücher/cb":                   "myapp://b%C3%BCcher/cb",
		"https://app.example.com./cb":         "https://app.example.com./cb",
		"https://my_app.example.com/cb":       "https://my_app.example.com/cb",
		"https://r3---sn-x.example.com/cb":    "https://r3---sn-x.example.com/cb",
		"https://bücher.example/%C3%BC/../cb": "https://xn--bcher-kva.example/cb",
	} {
		validatedURIs = nil
		location, err := callAppleRedirectHandler(t, "https://app.example.com", isValidFrontendRedirectURI, appleStateWithRedirect(redirect))
		if assert.NoError(t, err, redirect) {
			assert.Equal(t, []string{expected}, validatedURIs, redirect)
			assert.Equal(t, expected, withoutPostedFields(t, location), redirect)
		}
	}
}

func TestAppleCallbackRejectsHostsThatAreNotValidDomainNames(t *testing.T) {
	hookCalled := false
	alwaysValid := func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error) {
		hookCalled = true
		return true, nil
	}

	for _, redirect := range []string{
		// Full-width characters that are converted to characters not allowed in a domain name
		"https://evil.test／＠app.example.com/cb",
		"https://app.example.com／evil.test/cb",
		"https://evil.test＠app.example.com/cb",
		"https://app.example.com：443/cb",
		"https://app.example.com？x/cb",
		"https://app.example.com＃x/cb",
		"https://app.example.com％2f/cb",
		"https://app.example.com＼x/cb",
		"https://app.example.com＾x/cb",
		"https://app.example.com｜x/cb",
		"https://app.example.com＜x＞/cb",
		"https://app　example.com/cb",
		// ASCII characters that url.Parse allows in a host
		"https://app.example.com<x>/cb",
		// Empty labels
		"https://app..example.com/cb",
		"https://.app.example.com/cb",
		"https://app.example.com../cb",
		"https://app。。example.com/cb",
		"https://xn--zz.example/cb",
	} {
		_, err := callAppleRedirectHandler(t, "https://app.example.com", nil, appleStateWithRedirect(redirect))
		assert.Equal(t, supertokens.BadInputError{Msg: "frontendRedirectURI in state must be an absolute http or https URL"}, err, redirect)

		hookCalled = false
		_, err = callAppleRedirectHandler(t, "https://app.example.com", alwaysValid, appleStateWithRedirect(redirect))
		assert.Equal(t, supertokens.BadInputError{Msg: "frontendRedirectURI in state must be an absolute URL with an allowed scheme"}, err, redirect)
		assert.False(t, hookCalled, redirect)

		parsedRedirect, err := url.Parse(redirect)
		if assert.NoError(t, err, redirect) {
			_, err = getComparableOrigin(parsedRedirect)
			assert.Error(t, err, redirect)
		}
	}

	// A single trailing dot is allowed, but it is a different host for the default check
	for _, redirect := range []string{"https://app.example.com./cb", "https://app.example.com。/cb"} {
		parsedRedirect, err := url.Parse(redirect)
		if !assert.NoError(t, err, redirect) {
			continue
		}
		origin, err := getComparableOrigin(parsedRedirect)
		if assert.NoError(t, err, redirect) {
			assert.Equal(t, comparableOrigin{scheme: "https", host: "app.example.com.", port: "443"}, origin, redirect)
		}
		_, err = callAppleRedirectHandler(t, "https://app.example.com", nil, appleStateWithRedirect(redirect))
		assert.Equal(t, supertokens.BadInputError{Msg: "frontendRedirectURI in state does not match the configured origin"}, err, redirect)
	}
}
