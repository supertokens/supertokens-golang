package providers

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

const testIdTokenKid = "test-key"

type testIdTokenIssuer struct {
	server *httptest.Server
	key    *rsa.PrivateKey
}

func newTestIdTokenIssuer(t *testing.T) *testIdTokenIssuer {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	issuer := &testIdTokenIssuer{key: key}

	mux := http.NewServeMux()
	mux.HandleFunc("/jwks", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"keys": []map[string]interface{}{
				{
					"kty": "RSA",
					"kid": testIdTokenKid,
					"alg": "RS256",
					"use": "sig",
					"n":   base64.RawURLEncoding.EncodeToString(key.PublicKey.N.Bytes()),
					"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.PublicKey.E)).Bytes()),
				},
			},
		})
	})
	mux.HandleFunc("/.well-known/openid-configuration", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"issuer":   issuer.server.URL,
			"jwks_uri": issuer.server.URL + "/jwks",
		})
	})
	issuer.server = httptest.NewServer(mux)
	t.Cleanup(issuer.server.Close)
	return issuer
}

func (i *testIdTokenIssuer) jwksURI() string {
	return i.server.URL + "/jwks"
}

func (i *testIdTokenIssuer) sign(t *testing.T, aud interface{}) string {
	claims := jwt.MapClaims{
		"sub":   "user-1",
		"email": "user-1@example.com",
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	if aud != nil {
		claims["aud"] = aud
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = testIdTokenKid
	signed, err := token.SignedString(i.key)
	assert.NoError(t, err)
	return signed
}

func createCustomProviderForIdTokenTest(t *testing.T, issuer *testIdTokenIssuer, clientId string) *tpmodels.TypeProvider {
	provider, err := FindAndCreateProviderInstance([]tpmodels.ProviderInput{
		{
			Config: tpmodels.ProviderConfig{
				ThirdPartyId: "custom",
				JwksURI:      issuer.jwksURI(),
				Clients: []tpmodels.ProviderClientConfig{
					{
						ClientID:     clientId,
						ClientSecret: "secret",
					},
				},
			},
		},
	}, "custom", nil, &map[string]interface{}{})
	assert.NoError(t, err)
	return provider
}

func TestIdTokenWithAudienceOfAnotherClientIsRejected(t *testing.T) {
	issuer := newTestIdTokenIssuer(t)
	provider := createCustomProviderForIdTokenTest(t, issuer, "client-a")

	_, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{
		"id_token": issuer.sign(t, "client-b"),
	}, &map[string]interface{}{})
	assert.Error(t, err)
}

func TestIdTokenWithoutAudienceIsRejected(t *testing.T) {
	issuer := newTestIdTokenIssuer(t)
	provider := createCustomProviderForIdTokenTest(t, issuer, "client-a")

	_, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{
		"id_token": issuer.sign(t, nil),
	}, &map[string]interface{}{})
	assert.Error(t, err)
}

func TestIdTokenWithAudienceEqualToClientIdIsAccepted(t *testing.T) {
	issuer := newTestIdTokenIssuer(t)
	provider := createCustomProviderForIdTokenTest(t, issuer, "client-a")

	userInfo, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{
		"id_token": issuer.sign(t, "client-a"),
	}, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, "user-1", userInfo.ThirdPartyUserId)
}

func TestIdTokenWithAudienceArrayContainingClientIdIsAccepted(t *testing.T) {
	issuer := newTestIdTokenIssuer(t)
	provider := createCustomProviderForIdTokenTest(t, issuer, "client-a")

	userInfo, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{
		"id_token": issuer.sign(t, []string{"other-client", "client-a"}),
	}, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, "user-1", userInfo.ThirdPartyUserId)
}

func TestIdTokenAudienceIsCheckedAgainstActualClientIdForDevelopmentClientIds(t *testing.T) {
	issuer := newTestIdTokenIssuer(t)
	provider := createCustomProviderForIdTokenTest(t, issuer, DevKeyIdentifier+"client-a")

	userInfo, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{
		"id_token": issuer.sign(t, "client-a"),
	}, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, "user-1", userInfo.ThirdPartyUserId)

	_, err = provider.GetUserInfo(tpmodels.TypeOAuthTokens{
		"id_token": issuer.sign(t, "client-b"),
	}, &map[string]interface{}{})
	assert.Error(t, err)
}

func TestIdTokenIsRejectedWhenClientIdIsEmpty(t *testing.T) {
	issuer := newTestIdTokenIssuer(t)
	provider := createCustomProviderForIdTokenTest(t, issuer, "")

	_, err := provider.GetUserInfo(tpmodels.TypeOAuthTokens{
		"id_token": issuer.sign(t, "client-a"),
	}, &map[string]interface{}{})
	assert.Error(t, err)
}

func TestAppleIdTokenAudienceIsCheckedAgainstClientIdOfClientType(t *testing.T) {
	issuer := newTestIdTokenIssuer(t)

	providerInputs := []tpmodels.ProviderInput{
		{
			Config: tpmodels.ProviderConfig{
				ThirdPartyId:          "apple",
				OIDCDiscoveryEndpoint: issuer.server.URL,
				Clients: []tpmodels.ProviderClientConfig{
					{
						ClientType:   "web",
						ClientID:     "com.example.service",
						ClientSecret: "secret",
					},
					{
						ClientType:   "ios",
						ClientID:     "com.example.app",
						ClientSecret: "secret",
					},
				},
			},
		},
	}

	web := "web"
	webProvider, err := FindAndCreateProviderInstance(providerInputs, "apple", &web, &map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, issuer.jwksURI(), webProvider.Config.JwksURI)

	ios := "ios"
	iosProvider, err := FindAndCreateProviderInstance(providerInputs, "apple", &ios, &map[string]interface{}{})
	assert.NoError(t, err)

	_, err = webProvider.GetUserInfo(tpmodels.TypeOAuthTokens{"id_token": issuer.sign(t, "com.example.service")}, &map[string]interface{}{})
	assert.NoError(t, err)
	_, err = webProvider.GetUserInfo(tpmodels.TypeOAuthTokens{"id_token": issuer.sign(t, "com.example.app")}, &map[string]interface{}{})
	assert.Error(t, err)

	_, err = iosProvider.GetUserInfo(tpmodels.TypeOAuthTokens{"id_token": issuer.sign(t, "com.example.app")}, &map[string]interface{}{})
	assert.NoError(t, err)
	_, err = iosProvider.GetUserInfo(tpmodels.TypeOAuthTokens{"id_token": issuer.sign(t, "com.example.service")}, &map[string]interface{}{})
	assert.Error(t, err)
}
