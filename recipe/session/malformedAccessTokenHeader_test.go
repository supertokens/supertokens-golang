package session

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sessErrors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

// accessTokensWithMalformedHeaders returns unsigned tokens with a well formed payload, whose header
// has a missing or invalid typ, version or kid.
func accessTokensWithMalformedHeaders() map[string]string {
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"user","sessionHandle":"handle","refreshTokenHash1":"hash","exp":` +
		strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10) + `,"iat":1,"tId":"public"}`))
	headers := map[string]string{
		"missing typ":     `{"alg":"RS256","kid":"key","version":"5"}`,
		"numeric typ":     `{"alg":"RS256","kid":"key","version":"5","typ":1}`,
		"null typ":        `{"alg":"RS256","kid":"key","version":"5","typ":null}`,
		"null version":    `{"alg":"RS256","kid":"key","version":null,"typ":"JWT"}`,
		"null kid":        `{"alg":"RS256","kid":null,"version":"5","typ":"JWT"}`,
		"numeric kid":     `{"alg":"RS256","kid":1,"version":"5","typ":"JWT"}`,
		"numeric version": `{"alg":"RS256","kid":"key","version":5,"typ":"JWT"}`,
	}
	tokens := map[string]string{}
	for description, header := range headers {
		tokens[description] = base64.RawURLEncoding.EncodeToString([]byte(header)) + "." + payload + ".c2lnbmF0dXJl"
	}
	return tokens
}

func initSessionWithHeaderBasedTokens(t *testing.T) {
	connectionURI := unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.HeaderTransferMethod
				},
			}),
		},
	})
	assert.NoError(t, err)
}

func TestParseJWTWithoutSignatureVerificationRejectsMalformedHeaders(t *testing.T) {
	for description, token := range accessTokensWithMalformedHeaders() {
		assert.NotPanics(t, func() {
			_, err := ParseJWTWithoutSignatureVerification(token)
			assert.Error(t, err, description)
		}, description)
	}
}

func TestVerifySessionRejectsAccessTokensWithMalformedHeaders(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initSessionWithHeaderBasedTokens(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/protected", VerifySession(nil, func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	for description, token := range accessTokensWithMalformedHeaders() {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/protected", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := http.DefaultClient.Do(req)
		// A panic in the handler would close the connection, resulting in an error here
		if assert.NoError(t, err, description) {
			assert.Equal(t, http.StatusUnauthorized, res.StatusCode, description)
			res.Body.Close()
		}
	}
}

func TestGetSessionRejectsAccessTokensWithMalformedHeaders(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	initSessionWithHeaderBasedTokens(t)

	for description, token := range accessTokensWithMalformedHeaders() {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		assert.NotPanics(t, func() {
			sessionContainer, err := GetSession(req, httptest.NewRecorder(), nil)
			assert.Nil(t, sessionContainer, description)
			assert.True(t, errors.As(err, &sessErrors.UnauthorizedError{}), description)
		}, description)

		assert.NotPanics(t, func() {
			sessionContainer, err := GetSessionWithoutRequestResponse(token, nil, nil)
			assert.Nil(t, sessionContainer, description)
			isUnauthorised := errors.As(err, &sessErrors.UnauthorizedError{})
			isTryRefresh := errors.As(err, &sessErrors.TryRefreshTokenError{})
			assert.True(t, isUnauthorised || isTryRefresh, description)
		}, description)
	}
}
