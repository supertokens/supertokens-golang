package session

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestUseDynamicAccessTokenSigningKeySwitchingFromTrueToFalse(t *testing.T) {
	True := true
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
				UseDynamicAccessTokenSigningKey: &True,
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

	res, err := CreateNewSessionWithoutRequestResponse("public", "test-user-id", map[string]interface{}{
		"tokenProp": true,
	}, map[string]interface{}{
		"dbProp": true,
	}, nil)
	assert.Nil(t, err)

	tokens := res.GetAllSessionTokensDangerously()
	checkAccessTokenSigningKeyType(t, tokens, true)

	resetAll()
	False := false
	configValue = supertokens.TypeInput{
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
				UseDynamicAccessTokenSigningKey: &False,
			}),
		},
	}

	err = supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, nil)
	assert.NotNil(t, err)

	assert.Equal(t, "The access token doesn't match the useDynamicAccessTokenSigningKey setting", err.Error())
}

func TestUseDynamicAccessTokenSigningKeySwitchingFromTrueToFalseShouldWorkAfterRefresh(t *testing.T) {
	True := true
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
				UseDynamicAccessTokenSigningKey: &True,
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

	res, err := CreateNewSessionWithoutRequestResponse("public", "test-user-id", map[string]interface{}{
		"tokenProp": true,
	}, map[string]interface{}{
		"dbProp": true,
	}, nil)
	assert.Nil(t, err)

	tokens := res.GetAllSessionTokensDangerously()
	checkAccessTokenSigningKeyType(t, tokens, true)

	resetAll()
	False := false
	configValue = supertokens.TypeInput{
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
				UseDynamicAccessTokenSigningKey: &False,
			}),
		},
	}

	err = supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	refreshedSession, err := RefreshSessionWithoutRequestResponse(*tokens.RefreshToken, nil, tokens.AntiCsrfToken)
	assert.Nil(t, err)

	tokens = refreshedSession.GetAllSessionTokensDangerously()
	checkAccessTokenSigningKeyType(t, tokens, false)

	verifiedSession, err := GetSessionWithoutRequestResponse(tokens.AccessToken, tokens.AntiCsrfToken, nil)
	assert.Nil(t, err)

	tokensAfterVerify := verifiedSession.GetAllSessionTokensDangerously()
	assert.True(t, tokensAfterVerify.AccessAndFrontendTokenUpdated)

	verified2Session, err := GetSessionWithoutRequestResponse(tokensAfterVerify.AccessToken, tokensAfterVerify.AntiCsrfToken, nil)
	assert.Nil(t, err)

	tokensAfterVerify2 := verified2Session.GetAllSessionTokensDangerously()
	assert.False(t, tokensAfterVerify2.AccessAndFrontendTokenUpdated)
}

func checkAccessTokenSigningKeyType(t *testing.T, tokens sessmodels.SessionTokens, isDynamic bool) {
	info, err := ParseJWTWithoutSignatureVerification(tokens.AccessToken)
	assert.Nil(t, err)
	if isDynamic {
		assert.True(t, strings.HasPrefix(*info.KID, "d-"))
	} else {
		assert.True(t, strings.HasPrefix(*info.KID, "s-"))
	}
}
