package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
