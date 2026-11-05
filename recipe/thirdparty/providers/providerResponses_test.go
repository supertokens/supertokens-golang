package providers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// newStaticResponseServer serves fixed responses per path.
func newStaticResponseServer(t *testing.T, responses map[string]struct {
	status int
	body   string
}) *httptest.Server {
	mux := http.NewServeMux()
	for path, response := range responses {
		response := response
		mux.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(response.status)
			rw.Write([]byte(response.body))
		})
	}
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

type staticResponse = struct {
	status int
	body   string
}

func createProviderForResponseTest(t *testing.T, config tpmodels.ProviderConfig) (*tpmodels.TypeProvider, error) {
	if len(config.Clients) == 0 {
		config.Clients = []tpmodels.ProviderClientConfig{{ClientID: "client-id", ClientSecret: "client-secret"}}
	}
	return FindAndCreateProviderInstance([]tpmodels.ProviderInput{{Config: config}}, config.ThirdPartyId, nil, &map[string]interface{}{})
}

func TestExchangeAuthCodeRejectsNonStringCode(t *testing.T) {
	provider, err := createProviderForResponseTest(t, tpmodels.ProviderConfig{
		ThirdPartyId:  "custom",
		TokenEndpoint: "http://127.0.0.1:1/token",
	})
	assert.NoError(t, err)

	for _, code := range []interface{}{123, map[string]interface{}{}, []interface{}{"a"}, true} {
		assert.NotPanics(t, func() {
			_, err = provider.ExchangeAuthCodeForOAuthTokens(tpmodels.TypeRedirectURIInfo{
				RedirectURIOnProviderDashboard: "https://app.example.com/callback",
				RedirectURIQueryParams:         map[string]interface{}{"code": code},
			}, &map[string]interface{}{})
		})
		assert.True(t, errors.As(err, &supertokens.BadInputError{}))
	}
}

func TestUserInfoMapPathThroughNonObjectValueReturnsAnError(t *testing.T) {
	server := newStaticResponseServer(t, map[string]staticResponse{
		"/userinfo": {status: 200, body: `{"data": "not-an-object"}`},
	})
	provider, err := createProviderForResponseTest(t, tpmodels.ProviderConfig{
		ThirdPartyId:     "custom",
		UserInfoEndpoint: server.URL + "/userinfo",
		UserInfoMap: tpmodels.TypeUserInfoMap{
			FromUserInfoAPI: tpmodels.TypeUserInfoMapFields{UserId: "data.id", Email: "data.email.address"},
		},
	})
	assert.NoError(t, err)

	assert.NotPanics(t, func() {
		_, err = provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "token"}, &map[string]interface{}{})
	})
	assert.Error(t, err)
}

func TestOIDCDiscoveryReturningNonObjectReturnsAnError(t *testing.T) {
	server := newStaticResponseServer(t, map[string]staticResponse{
		"/.well-known/openid-configuration": {status: 200, body: `["not", "an", "object"]`},
	})

	assert.NotPanics(t, func() {
		_, err := createProviderForResponseTest(t, tpmodels.ProviderConfig{
			ThirdPartyId:          "custom",
			OIDCDiscoveryEndpoint: server.URL + "/.well-known/openid-configuration",
		})
		assert.Error(t, err)
	})
}

func TestGithubResponsesWithUnexpectedShapeDoNotPanic(t *testing.T) {
	for description, responses := range map[string]map[string]staticResponse{
		"token check app is not an object": {
			"/applications/client-id/token": {status: 200, body: `{"app": "not-an-object"}`},
		},
		"user without id": {
			"/applications/client-id/token": {status: 200, body: `{"app": {"client_id": "client-id"}}`},
			"/user":                         {status: 200, body: `{"login": "someone"}`},
			"/user/emails":                  {status: 200, body: `[]`},
		},
		"emails is not an array": {
			"/applications/client-id/token": {status: 200, body: `{"app": {"client_id": "client-id"}}`},
			"/user":                         {status: 200, body: `{"id": 1}`},
			"/user/emails":                  {status: 200, body: `{"message": "not an array"}`},
		},
		"email entries with unexpected types": {
			"/applications/client-id/token": {status: 200, body: `{"app": {"client_id": "client-id"}}`},
			"/user":                         {status: 200, body: `{"id": 1}`},
			"/user/emails":                  {status: 200, body: `["a", {"primary": "yes"}, {"primary": true, "email": 1}]`},
		},
	} {
		t.Run(description, func(t *testing.T) {
			server := newStaticResponseServer(t, responses)
			redirectHostsToTestServer(t, server, "api.github.com")

			provider, err := createProviderForResponseTest(t, tpmodels.ProviderConfig{ThirdPartyId: "github"})
			assert.NoError(t, err)

			assert.NotPanics(t, func() {
				userInfo, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "token"}, &map[string]interface{}{})
				if err == nil {
					assert.Nil(t, userInfo.Email)
				}
			})
		})
	}
}

func TestBitbucketWithoutAccessToEmailsDoesNotPanic(t *testing.T) {
	server := newStaticResponseServer(t, map[string]staticResponse{
		"/2.0/user":        {status: 200, body: `{"uuid": "{user-uuid}"}`},
		"/2.0/user/emails": {status: 403, body: `{"type": "error"}`},
	})
	redirectHostsToTestServer(t, server, "api.bitbucket.org")

	provider, err := createProviderForResponseTest(t, tpmodels.ProviderConfig{ThirdPartyId: "bitbucket"})
	assert.NoError(t, err)

	assert.NotPanics(t, func() {
		userInfo, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "token"}, &map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, "{user-uuid}", userInfo.ThirdPartyUserId)
		assert.Nil(t, userInfo.Email)
	})
}

func TestLinkedinWithoutEmailDoesNotPanic(t *testing.T) {
	server := newStaticResponseServer(t, map[string]staticResponse{
		"/v2/userinfo": {status: 200, body: `{"sub": "linkedin-user"}`},
	})
	redirectHostsToTestServer(t, server, "api.linkedin.com")

	provider, err := createProviderForResponseTest(t, tpmodels.ProviderConfig{ThirdPartyId: "linkedin"})
	assert.NoError(t, err)

	assert.NotPanics(t, func() {
		userInfo, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "token"}, &map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, "linkedin-user", userInfo.ThirdPartyUserId)
		assert.Nil(t, userInfo.Email)
	})
}

func TestPostRequestToUnreachableServerReturnsAnError(t *testing.T) {
	server := httptest.NewServer(http.NewServeMux())
	unreachableURL := server.URL + "/token"
	server.Close()

	assert.NotPanics(t, func() {
		_, _, err := doPostRequest(unreachableURL, map[string]interface{}{"code": "abc"}, nil)
		assert.Error(t, err)
	})
	assert.NotPanics(t, func() {
		_, _, err := doPostRequestWithJSONBody(unreachableURL, map[string]interface{}{"access_token": "abc"}, nil)
		assert.Error(t, err)
	})
}
