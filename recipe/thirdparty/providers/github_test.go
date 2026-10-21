package providers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// hostRewritingTransport sends requests for the given hosts to a test server instead.
type hostRewritingTransport struct {
	hosts  map[string]bool
	target *url.URL
	base   http.RoundTripper
}

func (t *hostRewritingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.hosts[req.URL.Host] {
		req = req.Clone(req.Context())
		req.URL.Scheme = t.target.Scheme
		req.URL.Host = t.target.Host
		req.Host = t.target.Host
	}
	return t.base.RoundTrip(req)
}

// redirectHostsToTestServer routes requests made through http.DefaultTransport for the given hosts
// to the test server for the duration of the test.
func redirectHostsToTestServer(t *testing.T, server *httptest.Server, hosts ...string) {
	target, err := url.Parse(server.URL)
	assert.NoError(t, err)
	hostSet := map[string]bool{}
	for _, h := range hosts {
		hostSet[h] = true
	}
	original := http.DefaultTransport
	http.DefaultTransport = &hostRewritingTransport{hosts: hostSet, target: target, base: original}
	t.Cleanup(func() {
		http.DefaultTransport = original
	})
}

type fakeGithubAPI struct {
	mu                 sync.Mutex
	tokenCheckCalls    int
	userInfoCalls      int
	validAccessToken   string
	appClientIdOfToken string
}

func newFakeGithubAPI(t *testing.T, validAccessToken string, appClientIdOfToken string) *fakeGithubAPI {
	api := &fakeGithubAPI{validAccessToken: validAccessToken, appClientIdOfToken: appClientIdOfToken}

	mux := http.NewServeMux()
	mux.HandleFunc("/applications/", func(rw http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		api.tokenCheckCalls++
		api.mu.Unlock()
		rw.Header().Set("Content-Type", "application/json")
		// Like GitHub, the token check needs the app's client credentials
		clientId, clientSecret, ok := r.BasicAuth()
		if !ok || clientId != "github-client-id" || clientSecret != "github-client-secret" {
			rw.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(rw).Encode(map[string]interface{}{"message": "Bad credentials"})
			return
		}
		// Like GitHub, the token check only accepts a JSON body
		var body struct {
			AccessToken string `json:"access_token"`
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") || json.NewDecoder(r.Body).Decode(&body) != nil {
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(map[string]interface{}{"message": "Problems parsing JSON"})
			return
		}
		if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, "/token") || body.AccessToken != api.validAccessToken {
			rw.WriteHeader(http.StatusNotFound)
			json.NewEncoder(rw).Encode(map[string]interface{}{"message": "Not Found"})
			return
		}
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"app": map[string]interface{}{"client_id": api.appClientIdOfToken},
		})
	})
	mux.HandleFunc("/user/emails", func(rw http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		api.userInfoCalls++
		api.mu.Unlock()
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode([]interface{}{
			map[string]interface{}{"email": "user@example.com", "primary": true, "verified": true},
		})
	})
	mux.HandleFunc("/user", func(rw http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		api.userInfoCalls++
		api.mu.Unlock()
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{"id": 123})
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	redirectHostsToTestServer(t, server, "api.github.com")
	return api
}

func createGithubProviderForTest(t *testing.T, validateAccessToken func(accessToken string, clientConfig tpmodels.ProviderConfigForClientType, userContext supertokens.UserContext) error) *tpmodels.TypeProvider {
	provider, err := FindAndCreateProviderInstance([]tpmodels.ProviderInput{
		{
			Config: tpmodels.ProviderConfig{
				ThirdPartyId: "github",
				Clients: []tpmodels.ProviderClientConfig{
					{
						ClientID:     "github-client-id",
						ClientSecret: "github-client-secret",
					},
				},
				ValidateAccessToken: validateAccessToken,
			},
		},
	}, "github", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	return provider
}

func TestGithubSignInFailsWhenGithubRejectsTheAccessToken(t *testing.T) {
	fakeAPI := newFakeGithubAPI(t, "valid-token", "github-client-id")
	provider := createGithubProviderForTest(t, nil)

	_, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "token-unknown-to-github"}, &map[string]interface{}{})
	assert.Error(t, err)
	assert.Equal(t, 1, fakeAPI.tokenCheckCalls)
	assert.Equal(t, 0, fakeAPI.userInfoCalls)
}

func TestGithubSignInFailsWhenTheAccessTokenBelongsToAnotherApp(t *testing.T) {
	fakeAPI := newFakeGithubAPI(t, "valid-token", "another-client-id")
	provider := createGithubProviderForTest(t, nil)

	_, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "valid-token"}, &map[string]interface{}{})
	assert.Error(t, err)
	assert.Equal(t, 1, fakeAPI.tokenCheckCalls)
	assert.Equal(t, 0, fakeAPI.userInfoCalls)
}

func TestGithubSignInWorksWhenGithubAcceptsTheAccessToken(t *testing.T) {
	fakeAPI := newFakeGithubAPI(t, "valid-token", "github-client-id")
	provider := createGithubProviderForTest(t, nil)

	userInfo, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "valid-token"}, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, "123", userInfo.ThirdPartyUserId)
	assert.Equal(t, "user@example.com", userInfo.Email.ID)
	assert.True(t, userInfo.Email.IsVerified)
	assert.Equal(t, 1, fakeAPI.tokenCheckCalls)
}

func TestGithubSignInUsesCustomValidateAccessToken(t *testing.T) {
	fakeAPI := newFakeGithubAPI(t, "valid-token", "github-client-id")
	validatedTokens := []string{}
	provider := createGithubProviderForTest(t, func(accessToken string, clientConfig tpmodels.ProviderConfigForClientType, userContext supertokens.UserContext) error {
		validatedTokens = append(validatedTokens, accessToken)
		assert.Equal(t, "github-client-id", clientConfig.ClientID)
		if accessToken != "custom-valid-token" {
			return errors.New("rejected by custom validator")
		}
		return nil
	})

	_, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "valid-token"}, &map[string]interface{}{})
	assert.EqualError(t, err, "rejected by custom validator")
	assert.Equal(t, 0, fakeAPI.userInfoCalls)

	userInfo, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "custom-valid-token"}, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, "123", userInfo.ThirdPartyUserId)

	assert.Equal(t, []string{"valid-token", "custom-valid-token"}, validatedTokens)
	assert.Equal(t, 0, fakeAPI.tokenCheckCalls)
}

func createGithubProviderWithClientSecretForTest(t *testing.T, clientSecret string) *tpmodels.TypeProvider {
	provider, err := FindAndCreateProviderInstance([]tpmodels.ProviderInput{
		{
			Config: tpmodels.ProviderConfig{
				ThirdPartyId: "github",
				Clients: []tpmodels.ProviderClientConfig{
					{
						ClientID:     "github-client-id",
						ClientSecret: clientSecret,
					},
				},
			},
		},
	}, "github", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	return provider
}

func TestGithubTokenCheckReportsAMissingClientSecret(t *testing.T) {
	fakeAPI := newFakeGithubAPI(t, "valid-token", "github-client-id")
	provider := createGithubProviderWithClientSecretForTest(t, "")

	_, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "valid-token"}, &map[string]interface{}{})
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "ClientSecret")
	}
	assert.Equal(t, 0, fakeAPI.tokenCheckCalls)
	assert.Equal(t, 0, fakeAPI.userInfoCalls)
}

func TestGithubTokenCheckReportsRejectedClientCredentials(t *testing.T) {
	fakeAPI := newFakeGithubAPI(t, "valid-token", "github-client-id")
	provider := createGithubProviderWithClientSecretForTest(t, "wrong-client-secret")

	_, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "valid-token"}, &map[string]interface{}{})
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "client credentials")
	}
	assert.Equal(t, 1, fakeAPI.tokenCheckCalls)
	assert.Equal(t, 0, fakeAPI.userInfoCalls)
}

func TestGithubTokenCheckReturnsNetworkErrors(t *testing.T) {
	server := httptest.NewServer(http.NewServeMux())
	server.Close()
	redirectHostsToTestServer(t, server, "api.github.com")
	provider := createGithubProviderForTest(t, nil)

	_, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": "valid-token"}, &map[string]interface{}{})
	if assert.Error(t, err) {
		assert.NotEqual(t, "Invalid access token", err.Error())
	}
}
