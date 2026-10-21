package thirdparty

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func setUpSignInUpWithUserInfoEndpoint(t *testing.T, userInfoHandler http.HandlerFunc) *httptest.Server {
	userInfoServer := httptest.NewServer(userInfoHandler)
	t.Cleanup(userInfoServer.Close)

	connectionURI := unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
			Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId:     "custom",
								UserInfoEndpoint: userInfoServer.URL + "/userinfo",
								Clients: []tpmodels.ProviderClientConfig{
									{ClientID: "client-id", ClientSecret: "client-secret"},
								},
							},
						},
					},
				},
			}),
		},
	})
	assert.NoError(t, err)

	testServer := httptest.NewServer(supertokens.Middleware(http.NewServeMux()))
	t.Cleanup(testServer.Close)
	return testServer
}

func assertSignInUpWithAccessTokenFails(t *testing.T, testServer *httptest.Server) {
	res, err := http.Post(testServer.URL+"/auth/signinup", "application/json", strings.NewReader(`{"thirdPartyId": "custom", "oAuthTokens": {"access_token": "some-access-token"}}`))
	// A panic in the handler would close the connection, resulting in an error here
	if assert.NoError(t, err) {
		defer res.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Empty(t, res.Header.Values("St-Access-Token"))
		assert.Empty(t, res.Header.Values("Set-Cookie"))
	}
}

func TestSignInUpWithAccessTokenRejectedByUserInfoEndpointReturnsAnError(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpSignInUpWithUserInfoEndpoint(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(`{"error": "invalid_token"}`))
	})

	assertSignInUpWithAccessTokenFails(t, testServer)
}

func TestSignInUpWithUserInfoEndpointReturningAnArrayReturnsAnError(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer := setUpSignInUpWithUserInfoEndpoint(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`[{"sub": "user-1", "email": "user@example.com"}]`))
	})

	assertSignInUpWithAccessTokenFails(t, testServer)
}
