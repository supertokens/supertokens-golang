package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestJWTShouldCreateRightAccessTokenPayloadWithClaims(t *testing.T) {
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
			Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						oCreateNewSession := *originalImplementation.CreateNewSession
						nCreateNewSession := func(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							if accessTokenPayload == nil {
								accessTokenPayload = map[string]interface{}{}
							}
							claim, _ := TrueClaim()
							accessTokenPayload, err := claim.Build(userID, accessTokenPayload, userContext)
							if err != nil {
								return nil, err
							}
							return oCreateNewSession(req, res, userID, accessTokenPayload, sessionDataInDatabase, userContext)
						}
						*originalImplementation.CreateNewSession = nCreateNewSession
						return originalImplementation
					},
				},
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.8", cdiVersion) == "2.8" {
		return
	}

	mux := http.NewServeMux()
	var sessionContainer sessmodels.SessionContainer

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		var err error
		sessionContainer, err = CreateNewSession(r, rw, "rope", map[string]interface{}{}, map[string]interface{}{})
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

	sessInfo, err := GetSessionInformation(sessionContainer.GetHandle())
	assert.NoError(t, err)
	jwtPayloadStr := sessInfo.AccessTokenPayload["jwt"].(string)
	jwtPayload := jwt.MapClaims{}

	_, _, err = (&jwt.Parser{}).ParseUnverified(jwtPayloadStr, jwtPayload)
	assert.NoError(t, err)

	assert.Equal(t, true, jwtPayload["st-true"].(map[string]interface{})["v"])
	assert.Equal(t, "rope", jwtPayload["sub"])
}

func TestAssertClaimsWithPayloadWithJWTAndCallRightUpdateAccessTokenPayload(t *testing.T) {
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
			Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.8", cdiVersion) == "2.8" {
		return
	}

	mux := http.NewServeMux()
	var sessionContainer sessmodels.SessionContainer
	accessTokenPayload := map[string]interface{}{
		"hello": "world",
	}

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		var err error
		sessionContainer, err = CreateNewSession(r, rw, "rope", accessTokenPayload, map[string]interface{}{})
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

	validateCallCount := 0
	var validationPayload map[string]interface{}

	validate := func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
		validateCallCount += 1

		validationPayload = payload

		return claims.ClaimValidationResult{
			IsValid: true,
		}
	}

	_, validators := StubClaimWithRefetch(validate)
	err = sessionContainer.AssertClaims([]claims.SessionClaimValidator{
		validators.Stub(),
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, validateCallCount)
	assert.Equal(t, "world", validationPayload["hello"])
	assert.NotNil(t, validationPayload, "st-stub")
	assert.Equal(t, "stub", validationPayload["st-stub"].(map[string]interface{})["v"])

	// Check if claim was updated in jwt
	sessInfo, err := GetSessionInformation(sessionContainer.GetHandle())
	assert.NoError(t, err)
	jwtPayloadStr := sessInfo.AccessTokenPayload["jwt"].(string)
	jwtPayload := jwt.MapClaims{}

	_, _, err = (&jwt.Parser{}).ParseUnverified(jwtPayloadStr, jwtPayload)
	assert.NoError(t, err)

	assert.Equal(t, "stub", jwtPayload["st-stub"].(map[string]interface{})["v"])
	assert.Equal(t, "rope", jwtPayload["sub"])
}

func TestMergeIntoAccessTokenPayloadForJWT(t *testing.T) {
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
			Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.8", cdiVersion) == "2.8" {
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		CreateNewSession(r, rw, "rope", nil, map[string]interface{}{})
	})

	mux.HandleFunc("/verifySession", VerifySession(nil, func(rw http.ResponseWriter, r *http.Request) {
		sessionContainer := GetSessionFromRequestContext(r.Context())
		assert.NotNil(t, sessionContainer)

		sessionContainer.MergeIntoAccessTokenPayload(map[string]interface{}{
			"testClaim": "newValue",
		})
		jwtPayloadStr := sessionContainer.GetAccessTokenPayload()["jwt"].(string)
		jwtPayload := jwt.MapClaims{}

		_, _, err = (&jwt.Parser{}).ParseUnverified(jwtPayloadStr, jwtPayload)
		assert.NoError(t, err)

		assert.Equal(t, "newValue", jwtPayload["testClaim"])
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	cookieData := unittesting.ExtractInfoFromResponse(res)

	reqV, err := http.NewRequest(http.MethodGet, testServer.URL+"/verifySession", nil)
	assert.NoError(t, err)
	reqV.Header.Add("Cookie", "sAccessToken="+cookieData["sAccessToken"])
	reqV.Header.Add("anti-csrf", cookieData["antiCsrf"])
	resv, err := http.DefaultClient.Do(reqV)
	assert.NoError(t, err)
	assert.Equal(t, resv.StatusCode, 200)
}
