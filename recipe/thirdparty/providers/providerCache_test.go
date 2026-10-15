package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testJWKSBody = `{"keys":[{"kty":"RSA","kid":"test","use":"sig","alg":"RS256","n":"sXchDaQebHnPiGvyDOAT4saGEUetSyo9MKLOoWFsueri23bOdgWp4Dy1WlUzewbgBHod5pcM9H95GQRV3JDXboIRROSBigeC5yjU1hGzHHyXss8UDprecbAYxknTcQkhslANGRUZmdTOQ5qTRsLAt6BTYuyvVRdhS8exSZEy_c4gs_7svlJJQ4H9_NxsiIoLwAEk7-Q3UXERGYw_75IDrGA84-lA_-Ct4eTlXHBIY2EaV7t7LjJaynVJCpkv4LKjTTAumiGUIuQhrNhZLuF_RJLqHpM2kgWFLU7-VTdL1VbC2tejvcI2BlMkEpk1BzBZI0KQB0GaDWFLN-aEAw3vRw","e":"AQAB"}]}`

// newSlowServer answers every request with body, but only after release is closed. It closes
// requested on the first request.
func newSlowServer(body string) (server *httptest.Server, requested chan struct{}, release chan struct{}) {
	requested = make(chan struct{})
	release = make(chan struct{})
	first := true
	server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if first {
			first = false
			close(requested)
		}
		<-release
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(body))
	}))
	return server, requested, release
}

func newFastServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(body))
	}))
}

// returnsWithin reports whether fn returns within the timeout.
func returnsWithin(timeout time.Duration, fn func()) bool {
	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

func TestCachedOIDCDiscoveryInfoIsNotBlockedByAnotherIssuersFetch(t *testing.T) {
	cachedServer := newFastServer(`{"issuer":"cached"}`)
	defer cachedServer.Close()
	_, err := getOIDCDiscoveryInfo(cachedServer.URL + "/cached")
	assert.NoError(t, err)

	slowServer, requested, release := newSlowServer(`{"issuer":"slow"}`)
	defer slowServer.Close()
	defer close(release)

	go getOIDCDiscoveryInfo(slowServer.URL + "/slow")
	<-requested

	assert.True(t, returnsWithin(2*time.Second, func() {
		info, err := getOIDCDiscoveryInfo(cachedServer.URL + "/cached")
		assert.NoError(t, err)
		assert.Equal(t, "cached", info["issuer"])
	}), "a cached lookup waited for another issuer's discovery request")
}

func TestCachedJWKSIsNotBlockedByAnotherJWKSFetch(t *testing.T) {
	cachedServer := newFastServer(testJWKSBody)
	defer cachedServer.Close()
	cachedJWKS, err := getJWKSFromURL(cachedServer.URL + "/cached")
	if !assert.NoError(t, err) {
		return
	}

	slowServer, requested, release := newSlowServer(testJWKSBody)
	defer slowServer.Close()
	defer close(release)

	go getJWKSFromURL(slowServer.URL + "/slow")
	<-requested

	assert.True(t, returnsWithin(2*time.Second, func() {
		jwks, err := getJWKSFromURL(cachedServer.URL + "/cached")
		assert.NoError(t, err)
		assert.Same(t, cachedJWKS, jwks)
	}), "a cached lookup waited for another JWKS request")
}
