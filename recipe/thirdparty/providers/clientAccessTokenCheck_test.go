package providers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// fakeProviderAPI serves the given handlers (by path) for the given hosts, and counts the calls per path.
type fakeProviderAPI struct {
	mu       sync.Mutex
	calls    map[string]int
	requests map[string]*http.Request
	forms    map[string]map[string][]string
}

func (api *fakeProviderAPI) callCount(path string) int {
	api.mu.Lock()
	defer api.mu.Unlock()
	return api.calls[path]
}

func newFakeProviderAPI(t *testing.T, handlers map[string]func(r *http.Request) (int, interface{}), hosts ...string) *fakeProviderAPI {
	api := &fakeProviderAPI{calls: map[string]int{}, requests: map[string]*http.Request{}, forms: map[string]map[string][]string{}}
	mux := http.NewServeMux()
	for path, handler := range handlers {
		path, handler := path, handler
		mux.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
			_ = r.ParseForm()
			api.mu.Lock()
			api.calls[path]++
			api.requests[path] = r
			api.forms[path] = r.PostForm
			api.mu.Unlock()
			status, body := handler(r)
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(status)
			json.NewEncoder(rw).Encode(body)
		})
	}
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	redirectHostsToTestServer(t, server, hosts...)
	return api
}

func respond(status int, body interface{}) func(r *http.Request) (int, interface{}) {
	return func(r *http.Request) (int, interface{}) {
		return status, body
	}
}

func createProviderForTest(t *testing.T, config tpmodels.ProviderConfig, clientType *string) *tpmodels.TypeProvider {
	// Providers are created the same way the sign in up API creates them: a new instance for every request
	provider, err := FindAndCreateProviderInstance([]tpmodels.ProviderInput{{Config: config}}, config.ThirdPartyId, clientType, &map[string]interface{}{})
	assert.NoError(t, err)
	return provider
}

func googleConfigForTest() tpmodels.ProviderConfig {
	return tpmodels.ProviderConfig{
		ThirdPartyId: "google",
		Clients: []tpmodels.ProviderClientConfig{
			{ClientType: "web", ClientID: "web-client", ClientSecret: "web-secret"},
			{ClientType: "android", ClientID: "android-client"},
		},
	}
}

func newFakeGoogleAPI(t *testing.T, tokenInfo func(r *http.Request) (int, interface{})) *fakeProviderAPI {
	return newFakeProviderAPI(t, map[string]func(r *http.Request) (int, interface{}){
		"/.well-known/openid-configuration": respond(200, map[string]interface{}{
			"issuer":                 "https://accounts.google.com",
			"authorization_endpoint": "https://accounts.google.com/o/oauth2/v2/auth",
			"token_endpoint":         "https://oauth2.googleapis.com/token",
			"userinfo_endpoint":      "https://openidconnect.googleapis.com/v1/userinfo",
			"jwks_uri":               "https://www.googleapis.com/oauth2/v3/certs",
		}),
		"/tokeninfo":   tokenInfo,
		"/token":       respond(200, map[string]interface{}{"access_token": "exchanged-token"}),
		"/v1/userinfo": respond(200, map[string]interface{}{"sub": "google-user-1", "email": "user@example.com", "email_verified": true}),
	}, "accounts.google.com", "oauth2.googleapis.com", "openidconnect.googleapis.com")
}

func web() *string {
	clientType := "web"
	return &clientType
}

func getUserInfoWithToken(provider *tpmodels.TypeProvider, accessToken string) (tpmodels.TypeUserInfo, error) {
	return provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": accessToken}, &map[string]interface{}{})
}

func TestGoogleAcceptsAnAccessTokenOfAConfiguredClient(t *testing.T) {
	api := newFakeGoogleAPI(t, respond(200, map[string]interface{}{"aud": "web-client", "azp": "web-client"}))
	provider := createProviderForTest(t, googleConfigForTest(), web())

	userInfo, err := getUserInfoWithToken(provider, "token-1")
	assert.NoError(t, err)
	assert.Equal(t, "google-user-1", userInfo.ThirdPartyUserId)
	assert.Equal(t, 1, api.callCount("/tokeninfo"))
	// The token is sent in the body, not in the URL
	assert.Equal(t, []string{"token-1"}, api.forms["/tokeninfo"]["access_token"])
	assert.Equal(t, "", api.requests["/tokeninfo"].URL.RawQuery)
}

func TestGoogleAcceptsAnAccessTokenRequestedByAnotherConfiguredClient(t *testing.T) {
	api := newFakeGoogleAPI(t, respond(200, map[string]interface{}{"aud": "some-other-client", "azp": "android-client"}))
	provider := createProviderForTest(t, googleConfigForTest(), web())

	_, err := getUserInfoWithToken(provider, "token-1")
	assert.NoError(t, err)
	assert.Equal(t, 1, api.callCount("/v1/userinfo"))
}

func TestGoogleRejectsAnAccessTokenOfAClientThatIsNotConfigured(t *testing.T) {
	api := newFakeGoogleAPI(t, respond(200, map[string]interface{}{"aud": "some-other-client", "azp": "some-other-client"}))
	provider := createProviderForTest(t, googleConfigForTest(), web())

	_, err := getUserInfoWithToken(provider, "token-1")
	assert.EqualError(t, err, "Access token does not belong to your application")
	assert.Equal(t, 0, api.callCount("/v1/userinfo"))
}

func TestGoogleRejectsAnAccessTokenTheTokenInfoAPIRejects(t *testing.T) {
	api := newFakeGoogleAPI(t, respond(400, map[string]interface{}{"error": "invalid_token", "error_description": "Invalid Value"}))
	provider := createProviderForTest(t, googleConfigForTest(), web())

	_, err := getUserInfoWithToken(provider, "token-1")
	assert.EqualError(t, err, "Invalid access token")
	assert.Equal(t, 0, api.callCount("/v1/userinfo"))
}

func TestGoogleReportsOtherTokenInfoAPIErrorsWithTheirStatus(t *testing.T) {
	newFakeGoogleAPI(t, respond(503, map[string]interface{}{"error": "unavailable"}))
	provider := createProviderForTest(t, googleConfigForTest(), web())

	_, err := getUserInfoWithToken(provider, "token-1")
	assert.EqualError(t, err, "Could not check the access token with Google (status 503)")
}

func TestGoogleTokenCheckReturnsNetworkErrors(t *testing.T) {
	api := newFakeGoogleAPI(t, respond(200, map[string]interface{}{"aud": "web-client"}))
	provider := createProviderForTest(t, googleConfigForTest(), web())
	server := httptest.NewServer(http.NewServeMux())
	server.Close()
	redirectHostsToTestServer(t, server, "oauth2.googleapis.com")

	_, err := getUserInfoWithToken(provider, "token-1")
	if assert.Error(t, err) {
		assert.NotEqual(t, "Invalid access token", err.Error())
	}
	assert.Equal(t, 0, api.callCount("/v1/userinfo"))
}

func exchangeCodeForTest(t *testing.T, provider *tpmodels.TypeProvider) tpmodels.TypeOAuthTokens {
	oAuthTokens, err := provider.ExchangeAuthCodeForOAuthTokens(tpmodels.TypeRedirectURIInfo{
		RedirectURIOnProviderDashboard: "https://example.com/auth/callback/google",
		RedirectURIQueryParams:         map[string]interface{}{"code": "code"},
	}, &map[string]interface{}{})
	assert.NoError(t, err)
	return oAuthTokens
}

func TestGoogleDoesNotCheckTheAccessTokenFromTheCodeExchange(t *testing.T) {
	api := newFakeGoogleAPI(t, respond(200, map[string]interface{}{"aud": "some-other-client", "azp": "some-other-client"}))
	provider := createProviderForTest(t, googleConfigForTest(), web())

	oAuthTokens := exchangeCodeForTest(t, provider)
	_, err := provider.GetUserInfo(oAuthTokens, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, 0, api.callCount("/tokeninfo"))
	assert.Equal(t, 1, api.callCount("/v1/userinfo"))
}

func TestGoogleChecksAnAccessTokenThatIsNotTheOneFromTheCodeExchange(t *testing.T) {
	newFakeGoogleAPI(t, respond(200, map[string]interface{}{"aud": "some-other-client", "azp": "some-other-client"}))
	provider := createProviderForTest(t, googleConfigForTest(), web())

	exchangeCodeForTest(t, provider)
	_, err := getUserInfoWithToken(provider, "another-token")
	assert.EqualError(t, err, "Access token does not belong to your application")
}

func TestGoogleDoesNotCheckIfTheConfigSetsValidateAccessToken(t *testing.T) {
	api := newFakeGoogleAPI(t, respond(200, map[string]interface{}{"aud": "some-other-client", "azp": "some-other-client"}))
	config := googleConfigForTest()
	validatedTokens := []string{}
	config.ValidateAccessToken = func(accessToken string, clientConfig tpmodels.ProviderConfigForClientType, userContext supertokens.UserContext) error {
		validatedTokens = append(validatedTokens, accessToken)
		return nil
	}
	provider := createProviderForTest(t, config, web())

	_, err := getUserInfoWithToken(provider, "token-1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"token-1"}, validatedTokens)
	assert.Equal(t, 0, api.callCount("/tokeninfo"))
}

func TestGoogleWorkspacesChecksTheAccessToken(t *testing.T) {
	newFakeGoogleAPI(t, respond(200, map[string]interface{}{"aud": "some-other-client", "azp": "some-other-client"}))
	config := googleConfigForTest()
	config.ThirdPartyId = "google-workspaces"
	config.Clients = config.Clients[:1]
	provider := createProviderForTest(t, config, nil)

	_, err := getUserInfoWithToken(provider, "token-1")
	assert.EqualError(t, err, "Access token does not belong to your application")
}

func facebookConfigForTest(clientSecret string) tpmodels.ProviderConfig {
	return tpmodels.ProviderConfig{
		ThirdPartyId: "facebook",
		Clients:      []tpmodels.ProviderClientConfig{{ClientID: "fb-app", ClientSecret: clientSecret}},
	}
}

func newFakeFacebookAPI(t *testing.T, debugToken func(r *http.Request) (int, interface{})) *fakeProviderAPI {
	return newFakeProviderAPI(t, map[string]func(r *http.Request) (int, interface{}){
		"/debug_token": debugToken,
		"/me":          respond(200, map[string]interface{}{"id": "fb-user-1", "email": "user@example.com"}),
	}, "graph.facebook.com")
}

func TestFacebookAcceptsAValidAccessTokenOfTheConfiguredApp(t *testing.T) {
	api := newFakeFacebookAPI(t, respond(200, map[string]interface{}{"data": map[string]interface{}{"app_id": "fb-app", "is_valid": true, "user_id": "fb-user-1"}}))
	provider := createProviderForTest(t, facebookConfigForTest("fb-secret"), nil)

	userInfo, err := getUserInfoWithToken(provider, "fb-token")
	assert.NoError(t, err)
	assert.Equal(t, "fb-user-1", userInfo.ThirdPartyUserId)
	query := api.requests["/debug_token"].URL.Query()
	assert.Equal(t, "fb-token", query.Get("input_token"))
	assert.Equal(t, "fb-app|fb-secret", query.Get("access_token"))
}

func TestFacebookAcceptsANumericAppId(t *testing.T) {
	newFakeProviderAPI(t, map[string]func(r *http.Request) (int, interface{}){
		"/debug_token": respond(200, map[string]interface{}{"data": map[string]interface{}{"app_id": 138483919580948, "is_valid": true}}),
		"/me":          respond(200, map[string]interface{}{"id": "fb-user-1", "email": "user@example.com"}),
	}, "graph.facebook.com")
	config := facebookConfigForTest("fb-secret")
	config.Clients[0].ClientID = "138483919580948"
	provider := createProviderForTest(t, config, nil)

	_, err := getUserInfoWithToken(provider, "fb-token")
	assert.NoError(t, err)
}

func TestFacebookRejectsAnAccessTokenOfAnotherApp(t *testing.T) {
	api := newFakeFacebookAPI(t, respond(200, map[string]interface{}{"data": map[string]interface{}{"app_id": "another-app", "is_valid": true}}))
	provider := createProviderForTest(t, facebookConfigForTest("fb-secret"), nil)

	_, err := getUserInfoWithToken(provider, "fb-token")
	assert.EqualError(t, err, "Access token does not belong to your application")
	assert.Equal(t, 0, api.callCount("/me"))
}

func TestFacebookRejectsAnInvalidAccessToken(t *testing.T) {
	newFakeFacebookAPI(t, respond(200, map[string]interface{}{"data": map[string]interface{}{"app_id": "fb-app", "is_valid": false}}))
	provider := createProviderForTest(t, facebookConfigForTest("fb-secret"), nil)

	_, err := getUserInfoWithToken(provider, "fb-token")
	assert.EqualError(t, err, "Invalid access token")
}

func TestFacebookReportsAFailedDebugTokenCall(t *testing.T) {
	newFakeFacebookAPI(t, respond(400, map[string]interface{}{"error": map[string]interface{}{"message": "Invalid OAuth access token", "code": 190}}))
	provider := createProviderForTest(t, facebookConfigForTest("fb-secret"), nil)

	_, err := getUserInfoWithToken(provider, "fb-token")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Could not check the access token with Facebook (status 400)")
	}
}

func TestFacebookReportsAMissingClientSecretWithoutCallingFacebook(t *testing.T) {
	api := newFakeFacebookAPI(t, respond(200, map[string]interface{}{"data": map[string]interface{}{"app_id": "fb-app", "is_valid": true}}))
	provider := createProviderForTest(t, facebookConfigForTest(""), nil)

	_, err := getUserInfoWithToken(provider, "fb-token")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "ClientSecret")
	}
	assert.Equal(t, 0, api.callCount("/debug_token"))
}

func newFakeGitlabAPI(t *testing.T, tokenInfo func(r *http.Request) (int, interface{})) *fakeProviderAPI {
	return newFakeProviderAPI(t, map[string]func(r *http.Request) (int, interface{}){
		"/.well-known/openid-configuration": respond(200, map[string]interface{}{
			"issuer":                 "https://gitlab.example.com",
			"authorization_endpoint": "https://gitlab.example.com/oauth/authorize",
			"token_endpoint":         "https://gitlab.example.com/oauth/token",
			"userinfo_endpoint":      "https://gitlab.example.com/oauth/userinfo",
			"jwks_uri":               "https://gitlab.example.com/oauth/discovery/keys",
		}),
		"/oauth/token/info": tokenInfo,
		"/oauth/userinfo":   respond(200, map[string]interface{}{"sub": "gl-user-1", "email": "user@example.com", "email_verified": true}),
	}, "gitlab.example.com")
}

func gitlabConfigForTest() tpmodels.ProviderConfig {
	return tpmodels.ProviderConfig{
		ThirdPartyId: "gitlab",
		Clients: []tpmodels.ProviderClientConfig{{
			ClientID:         "gitlab-app",
			ClientSecret:     "gitlab-secret",
			AdditionalConfig: map[string]interface{}{"gitlabBaseUrl": "https://gitlab.example.com"},
		}},
	}
}

func TestGitlabAcceptsAnAccessTokenOfTheConfiguredApplication(t *testing.T) {
	api := newFakeGitlabAPI(t, respond(200, map[string]interface{}{"application": map[string]interface{}{"uid": "gitlab-app"}}))
	provider := createProviderForTest(t, gitlabConfigForTest(), nil)

	_, err := getUserInfoWithToken(provider, "gl-token")
	assert.NoError(t, err)
	assert.Equal(t, "Bearer gl-token", api.requests["/oauth/token/info"].Header.Get("Authorization"))
	assert.Equal(t, 1, api.callCount("/oauth/userinfo"))
}

func TestGitlabRejectsAnAccessTokenOfAnotherApplication(t *testing.T) {
	api := newFakeGitlabAPI(t, respond(200, map[string]interface{}{"application": map[string]interface{}{"uid": "another-app"}}))
	provider := createProviderForTest(t, gitlabConfigForTest(), nil)

	_, err := getUserInfoWithToken(provider, "gl-token")
	assert.EqualError(t, err, "Access token does not belong to your application")
	assert.Equal(t, 0, api.callCount("/oauth/userinfo"))
}

func TestGitlabRejectsAnAccessTokenGitlabDoesNotAccept(t *testing.T) {
	newFakeGitlabAPI(t, respond(401, map[string]interface{}{"error": "invalid_token"}))
	provider := createProviderForTest(t, gitlabConfigForTest(), nil)

	_, err := getUserInfoWithToken(provider, "gl-token")
	assert.EqualError(t, err, "Invalid access token")
}

func newFakeDiscordAPI(t *testing.T, authInfo func(r *http.Request) (int, interface{})) *fakeProviderAPI {
	return newFakeProviderAPI(t, map[string]func(r *http.Request) (int, interface{}){
		"/api/oauth2/@me": authInfo,
		"/api/users/@me":  respond(200, map[string]interface{}{"id": "dc-user-1", "email": "user@example.com", "verified": true}),
	}, "discord.com")
}

func discordConfigForTest() tpmodels.ProviderConfig {
	return tpmodels.ProviderConfig{
		ThirdPartyId: "discord",
		Clients:      []tpmodels.ProviderClientConfig{{ClientID: "discord-app", ClientSecret: "discord-secret"}},
	}
}

func TestDiscordAcceptsAnAccessTokenOfTheConfiguredApplication(t *testing.T) {
	api := newFakeDiscordAPI(t, respond(200, map[string]interface{}{"application": map[string]interface{}{"id": "discord-app"}}))
	provider := createProviderForTest(t, discordConfigForTest(), nil)

	_, err := getUserInfoWithToken(provider, "dc-token")
	assert.NoError(t, err)
	assert.Equal(t, "Bearer dc-token", api.requests["/api/oauth2/@me"].Header.Get("Authorization"))
	assert.Equal(t, 1, api.callCount("/api/users/@me"))
}

func TestDiscordRejectsAnAccessTokenOfAnotherApplication(t *testing.T) {
	api := newFakeDiscordAPI(t, respond(200, map[string]interface{}{"application": map[string]interface{}{"id": "another-app"}}))
	provider := createProviderForTest(t, discordConfigForTest(), nil)

	_, err := getUserInfoWithToken(provider, "dc-token")
	assert.EqualError(t, err, "Access token does not belong to your application")
	assert.Equal(t, 0, api.callCount("/api/users/@me"))
}

func TestDiscordRejectsAnAccessTokenDiscordDoesNotAccept(t *testing.T) {
	newFakeDiscordAPI(t, respond(401, map[string]interface{}{"message": "401: Unauthorized", "code": 0}))
	provider := createProviderForTest(t, discordConfigForTest(), nil)

	_, err := getUserInfoWithToken(provider, "dc-token")
	assert.True(t, errors.Is(err, errInvalidAccessToken))
}
