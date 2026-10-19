package providers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const concurrentSignIns = 50

// newFakeOAuthServer serves an OIDC discovery document, a token endpoint that returns an access token
// derived from the posted code, and a user info endpoint that returns a user derived from the access token.
func newFakeOAuthServer(t *testing.T) *httptest.Server {
	var server *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"issuer":                 server.URL,
			"authorization_endpoint": server.URL + "/authorize",
			"token_endpoint":         server.URL + "/token",
			"userinfo_endpoint":      server.URL + "/userinfo",
		})
	})
	mux.HandleFunc("/token", func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"access_token": "access-token-for-" + r.PostForm.Get("code"),
		})
	})
	mux.HandleFunc("/userinfo", func(rw http.ResponseWriter, r *http.Request) {
		accessToken := r.URL.Query().Get("access_token")
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"id":    "user-for-" + accessToken,
			"email": accessToken + "@example.com",
		})
	})
	server = httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

func TestFacebookSignInsWithDifferentAccessTokensResolveTheirOwnUsers(t *testing.T) {
	server := newFakeOAuthServer(t)

	providerInputs := []tpmodels.ProviderInput{
		{
			Config: tpmodels.ProviderConfig{
				ThirdPartyId:     "facebook",
				UserInfoEndpoint: server.URL + "/userinfo",
				UserInfoEndpointQueryParams: map[string]interface{}{
					"fields": "id,email",
				},
				Clients: []tpmodels.ProviderClientConfig{
					{ClientID: "facebook-client-id", ClientSecret: "facebook-client-secret"},
				},
				// This test is about the user info request, so it skips the check of the access token with Facebook
				ValidateAccessToken: func(accessToken string, clientConfig tpmodels.ProviderConfigForClientType, userContext supertokens.UserContext) error {
					return nil
				},
			},
		},
	}

	for _, accessToken := range []string{"token-1", "token-2"} {
		// Each sign in creates a new provider instance from the same static config, like the API does
		provider, err := FindAndCreateProviderInstance(providerInputs, "facebook", nil, &map[string]interface{}{})
		assert.NoError(t, err)

		userInfo, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{"access_token": accessToken}, &map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, "user-for-"+accessToken, userInfo.ThirdPartyUserId)
	}

	assert.Equal(t, map[string]interface{}{"fields": "id,email"}, providerInputs[0].Config.UserInfoEndpointQueryParams)
	assert.Nil(t, providerInputs[0].Config.UserInfoEndpointHeaders)
}

func generatePKCS8PrivateKeyPEM(t *testing.T) string {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	der, err := x509.MarshalPKCS8PrivateKey(key)
	assert.NoError(t, err)
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))
}

// runConcurrentSignIns creates a provider instance per sign in from the same static config and runs
// the authorisation URL and code exchange steps of 50 sign ins in parallel.
func runConcurrentSignIns(t *testing.T, providerInputs []tpmodels.ProviderInput, thirdPartyId string) {
	wg := sync.WaitGroup{}
	for i := 0; i < concurrentSignIns; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			userContext := &map[string]interface{}{}
			provider, err := FindAndCreateProviderInstance(providerInputs, thirdPartyId, nil, userContext)
			if !assert.NoError(t, err) {
				return
			}
			_, err = provider.GetAuthorisationRedirectURL("https://app.example.com/auth/callback", userContext)
			assert.NoError(t, err)

			code := fmt.Sprintf("code-%d", i)
			verifier := "verifier"
			tokens, err := provider.ExchangeAuthCodeForOAuthTokens(tpmodels.TypeRedirectURIInfo{
				RedirectURIOnProviderDashboard: "https://app.example.com/auth/callback",
				RedirectURIQueryParams:         map[string]interface{}{"code": code},
				PKCECodeVerifier:               &verifier,
			}, userContext)
			if assert.NoError(t, err) {
				assert.Equal(t, "access-token-for-"+code, tokens["access_token"])
			}
		}(i)
	}
	wg.Wait()
}

func TestConcurrentGoogleWorkspacesSignInsDoNotModifyTheStaticConfig(t *testing.T) {
	server := newFakeOAuthServer(t)

	providerInputs := []tpmodels.ProviderInput{
		{
			Config: tpmodels.ProviderConfig{
				ThirdPartyId:          "google-workspaces",
				OIDCDiscoveryEndpoint: server.URL,
				AuthorizationEndpointQueryParams: map[string]interface{}{
					"prompt": "select_account",
				},
				Clients: []tpmodels.ProviderClientConfig{
					{
						ClientID:         "gw-client-id",
						ClientSecret:     "gw-client-secret",
						AdditionalConfig: map[string]interface{}{"hd": "example.com"},
					},
				},
			},
		},
	}

	runConcurrentSignIns(t, providerInputs, "google-workspaces")

	assert.Equal(t, map[string]interface{}{"prompt": "select_account"}, providerInputs[0].Config.AuthorizationEndpointQueryParams)
	assert.Equal(t, map[string]interface{}{"hd": "example.com"}, providerInputs[0].Config.Clients[0].AdditionalConfig)
}

func TestConcurrentOktaSignInsDoNotModifyTheStaticConfig(t *testing.T) {
	server := newFakeOAuthServer(t)

	providerInputs := []tpmodels.ProviderInput{
		{
			Config: tpmodels.ProviderConfig{
				ThirdPartyId: "okta",
				TokenEndpointBodyParams: map[string]interface{}{
					"audience": "api://default",
				},
				Clients: []tpmodels.ProviderClientConfig{
					{
						ClientID: "okta-client-id",
						AdditionalConfig: map[string]interface{}{
							"oktaDomain": server.URL,
							"privateKey": generatePKCS8PrivateKeyPEM(t),
						},
					},
				},
			},
		},
	}

	runConcurrentSignIns(t, providerInputs, "okta")

	assert.Equal(t, map[string]interface{}{"audience": "api://default"}, providerInputs[0].Config.TokenEndpointBodyParams)
}

func TestConcurrentTwitterSignInsDoNotModifyTheStaticConfig(t *testing.T) {
	server := newFakeOAuthServer(t)

	providerInputs := []tpmodels.ProviderInput{
		{
			Config: tpmodels.ProviderConfig{
				ThirdPartyId:  "twitter",
				TokenEndpoint: server.URL + "/token",
				TokenEndpointBodyParams: map[string]interface{}{
					"custom": "value",
				},
				Clients: []tpmodels.ProviderClientConfig{
					{ClientID: "twitter-client-id", ClientSecret: "twitter-client-secret"},
				},
			},
		},
	}

	runConcurrentSignIns(t, providerInputs, "twitter")

	assert.Equal(t, map[string]interface{}{"custom": "value"}, providerInputs[0].Config.TokenEndpointBodyParams)
}
