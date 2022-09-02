package session

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestEmptyClaimsArray(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	var sessionContainer sessmodels.SessionContainer

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		var err error
		sessionContainer, err = CreateNewSession(rw, "rope", map[string]interface{}{}, map[string]interface{}{})
		assert.NoError(t, err)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	fmt.Println(sessionContainer)

	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{})
	assert.NoError(t, err)
}

func TestAssertClaimsWithPayload(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	var sessionContainer sessmodels.SessionContainer
	accessTokenPayload := map[string]interface{}{
		"hello": "world",
	}

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		var err error
		sessionContainer, err = CreateNewSession(rw, "rope", accessTokenPayload, map[string]interface{}{})
		assert.NoError(t, err)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	fmt.Println(sessionContainer)

	validateCallCount := 0
	var validationPayload map[string]interface{}

	validate := func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
		validateCallCount += 1
		validationPayload = payload
		return claims.ClaimValidationResult{
			IsValid: true,
		}
	}

	_, validators := StubClaim(validate)
	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{
		validators.Stub(),
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, validateCallCount)
	assert.Equal(t, accessTokenPayload, validationPayload)
}
