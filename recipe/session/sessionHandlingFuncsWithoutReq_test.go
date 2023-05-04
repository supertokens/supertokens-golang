package session

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	sessionError "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"testing"
)

func TestShouldCreateNewSession(t *testing.T) {
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

	res, err2 := CreateNewSessionWithoutRequestResponse("test-user-id", map[string]interface{}{
		"tokenProp": true,
	}, map[string]interface{}{
		"dbProp": true,
	}, nil)

	if err2 != nil {
		t.Error(err2.Error())
	}

	tokens := res.GetAllSessionTokensDangerously()
	assert.Equal(t, tokens.AccessAndFrontendTokenUpdated, true)
	assert.Nil(t, tokens.AntiCsrfToken)

	payload := res.GetAccessTokenPayload()
	assert.Equal(t, payload["sub"], "test-user-id")
	assert.Equal(t, payload["tokenProp"], true)
	assert.Equal(t, payload["iss"], "https://api.supertokens.io/auth")

	dataInDb, err3 := res.GetSessionDataInDatabase()

	if err3 != nil {
		t.Error(err2.Error())
	}

	assert.Equal(t, dataInDb["dbProp"], true)
}

func TestShouldCreateSessionWithAntiCSRF(t *testing.T) {
	antiCsrf := "VIA_TOKEN"
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
				AntiCsrf: &antiCsrf,
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

	res, err2 := CreateNewSessionWithoutRequestResponse("test-user-id", map[string]interface{}{
		"tokenProp": true,
	}, map[string]interface{}{
		"dbProp": true,
	}, nil)

	if err2 != nil {
		t.Error(err2.Error())
	}

	tokens := res.GetAllSessionTokensDangerously()
	assert.Equal(t, tokens.AccessAndFrontendTokenUpdated, true)
	assert.NotNil(t, tokens.AntiCsrfToken)

	payload := res.GetAccessTokenPayload()
	assert.Equal(t, payload["sub"], "test-user-id")
	assert.Equal(t, payload["tokenProp"], true)
	assert.Equal(t, payload["iss"], "https://api.supertokens.io/auth")

	dataInDb, err3 := res.GetSessionDataInDatabase()

	if err3 != nil {
		t.Error(err2.Error())
	}

	assert.Equal(t, dataInDb["dbProp"], true)
}

func TestShouldValidateBasicAccessToken(t *testing.T) {
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

	res, err2 := CreateNewSessionWithoutRequestResponse("test-user-id", nil, nil, nil)
	if err2 != nil {
		t.Error(err2.Error())
	}

	tokens := res.GetAllSessionTokensDangerously()
	session, err3 := GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, nil)

	if err3 != nil {
		t.Error(err3.Error())
	}

	tokenInfo := session.GetAllSessionTokensDangerously()
	assert.Equal(t, tokenInfo, sessmodels.SessionTokens{
		AccessToken:                   tokens.AccessToken,
		RefreshToken:                  nil,
		AntiCsrfToken:                 tokens.AntiCsrfToken,
		FrontToken:                    tokens.FrontToken,
		AccessAndFrontendTokenUpdated: false,
	})
}

func TestShouldValidateBasicAccessTokenWithAntiCSRF(t *testing.T) {
	antiCsrf := "VIA_TOKEN"
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
				AntiCsrf: &antiCsrf,
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

	res, err2 := CreateNewSessionWithoutRequestResponse("test-user-id", nil, nil, nil)
	if err2 != nil {
		t.Error(err2.Error())
	}

	tokens := res.GetAllSessionTokensDangerously()
	session, err3 := GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, nil)

	if err3 != nil {
		t.Error(err3.Error())
	}

	tokenInfo := session.GetAllSessionTokensDangerously()
	assert.Equal(t, tokenInfo, sessmodels.SessionTokens{
		AccessToken:                   tokens.AccessToken,
		RefreshToken:                  nil,
		AntiCsrfToken:                 tokens.AntiCsrfToken,
		FrontToken:                    tokens.FrontToken,
		AccessAndFrontendTokenUpdated: false,
	})

	_, sessionErr2 := GetSessionWithoutRequestResponse(tokens.AccessToken, nil, nil)
	assert.NotNil(t, sessionErr2)

	assert.True(t, errors.As(sessionErr2, &sessionError.TryRefreshTokenError{}))

	antiCsrfCheck := false
	sessionWithoutAnti, sessionErr3 := GetSessionWithoutRequestResponse(tokens.AccessToken, nil, &sessmodels.VerifySessionOptions{
		AntiCsrfCheck: &antiCsrfCheck,
	})

	assert.Nil(t, sessionErr3)
	assert.NotNil(t, sessionWithoutAnti)
}

func TestShouldErrorForNonTokens(t *testing.T) {
	antiCsrf := "VIA_TOKEN"
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
				AntiCsrf: &antiCsrf,
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

	_, err2 := GetSessionWithoutRequestResponse("nope", nil, nil)
	assert.NotNil(t, err2)

	assert.True(t, errors.As(err2, &sessionError.UnauthorizedError{}))
}

func TestShouldReturnNilForNonTokenWithoutSessionRequired(t *testing.T) {
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

	sessionRequired := false
	session, err2 := GetSessionWithoutRequestResponse("nope", nil, &sessmodels.VerifySessionOptions{
		SessionRequired: &sessionRequired,
	})

	assert.Nil(t, err2)
	assert.Nil(t, session)
}

func TestShouldReturnErrorForClaimValidationFailures(t *testing.T) {
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

	res, err2 := CreateNewSessionWithoutRequestResponse("test-user-id", nil, nil, nil)
	if err2 != nil {
		t.Error(err2.Error())
	}

	tokens := res.GetAllSessionTokensDangerously()

	_, sessionErr := GetSessionWithoutRequestResponse(tokens.AccessToken, nil, &sessmodels.VerifySessionOptions{
		OverrideGlobalClaimValidators: func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
			globalClaimValidators = append(globalClaimValidators, claims.SessionClaimValidator{
				ID: "test",
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
					return claims.ClaimValidationResult{
						IsValid: false,
						Reason:  "test",
					}
				},
			})

			return globalClaimValidators, nil
		},
	})

	assert.NotNil(t, sessionErr)
	assert.True(t, errors.As(err2, &sessionError.InvalidClaimError{}))
}

func TestShouldRefreshSession(t *testing.T) {
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

	res, err2 := CreateNewSessionWithoutRequestResponse("test-user-id", map[string]interface{}{
		"tokenProp": true,
	}, map[string]interface{}{
		"dbProp": true,
	}, nil)

	if err2 != nil {
		t.Error(err2.Error())
	}

	tokens := res.GetAllSessionTokensDangerously()

	disableAntiCsrf := false
	session, err3 := RefreshSessionWithoutRequestResponse(*tokens.RefreshToken, &disableAntiCsrf, tokens.AntiCsrfToken)
	assert.Nil(t, err3)
	assert.NotNil(t, session)

	tokensAfterRefresh := session.GetAllSessionTokensDangerously()

	assert.Equal(t, tokensAfterRefresh.AccessAndFrontendTokenUpdated, true)
	assert.True(t, tokensAfterRefresh.AntiCsrfToken == nil)

	payload := session.GetAccessTokenPayload()
	assert.Equal(t, payload["sub"], "test-user-id")
	assert.Equal(t, payload["tokenProp"], true)

	sessionData, err4 := session.GetSessionDataInDatabase()
	if err4 != nil {
		t.Error(err4.Error())
	}

	assert.Equal(t, sessionData, map[string]interface{}{
		"dbProp": true,
	})

}

func TestShouldWorkWithAntiCSRF(t *testing.T) {
	antiCsrf := "VIA_TOKEN"
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
				AntiCsrf: &antiCsrf,
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

	createRes, createErr := CreateNewSessionWithoutRequestResponse("test-user-id", nil, nil, nil)
	if createErr != nil {
		t.Error(createErr.Error())
	}

	tokens := createRes.GetAllSessionTokensDangerously()

	disableAntiCsrf := false
	session, err3 := RefreshSessionWithoutRequestResponse(*tokens.RefreshToken, &disableAntiCsrf, tokens.AntiCsrfToken)

	assert.Nil(t, err3)
	assert.NotNil(t, session)

	tokensAfterRefresh := session.GetAllSessionTokensDangerously()

	assert.Equal(t, tokensAfterRefresh.AccessAndFrontendTokenUpdated, true)

	_, sessionErr2 := RefreshSessionWithoutRequestResponse(*tokensAfterRefresh.RefreshToken, &disableAntiCsrf, nil)
	assert.NotNil(t, sessionErr2)

	assert.True(t, errors.As(sessionErr2, &sessionError.UnauthorizedError{}))

	disableAntiCsrf = true
	sessionAfterRefresh, sessionAfterRefreshErr := RefreshSessionWithoutRequestResponse(*tokensAfterRefresh.RefreshToken, &disableAntiCsrf, nil)
	assert.Nil(t, sessionAfterRefreshErr)
	assert.NotNil(t, sessionAfterRefresh)

	finalTokens := sessionAfterRefresh.GetAllSessionTokensDangerously()
	assert.True(t, finalTokens.AccessAndFrontendTokenUpdated == true)
}

func TestRefreshShouldReturnErrorForNonTokens(t *testing.T) {
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

	_, err2 := RefreshSessionWithoutRequestResponse("nope", nil, nil)

	assert.NotNil(t, err2)
	assert.True(t, errors.As(err2, &sessionError.UnauthorizedError{}))
}
