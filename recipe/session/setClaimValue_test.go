package session

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestShouldSetClaimValueMergeTheRightValue(t *testing.T) {
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

	res := fakeRes{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	sessionContainer, err := CreateNewSession(req, res, "public", "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	trueClaim, _ := TrueClaim()
	sessionContainer.SetClaimValue(trueClaim, true)

	accessTokenPayload := sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 10, len(accessTokenPayload))
	assert.Equal(t, true, accessTokenPayload["st-true"].(map[string]interface{})["v"])
}

func TestShouldSetClaimValueOverwriteClaimValue(t *testing.T) {
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
						nCreateNewSession := func(userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, disableAntiCsrf *bool, tenantId string, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							trueClaim, _ := TrueClaim()
							accessTokenPayload, err := trueClaim.Build(userID, tenantId, accessTokenPayload, userContext)
							if err != nil {
								return nil, err
							}
							return oCreateNewSession(userID, accessTokenPayload, sessionDataInDatabase, disableAntiCsrf, tenantId, userContext)
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

	res := fakeRes{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	sessionContainer, err := CreateNewSession(req, res, "public", "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)
	accessTokenPayload := sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 10, len(accessTokenPayload))
	assert.Equal(t, true, accessTokenPayload["st-true"].(map[string]interface{})["v"])

	trueClaim, _ := TrueClaim()
	sessionContainer.SetClaimValue(trueClaim, false)

	accessTokenPayload = sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 10, len(accessTokenPayload))
	assert.Equal(t, false, accessTokenPayload["st-true"].(map[string]interface{})["v"])
}

func TestShouldSetClaimValueOverwriteClaimValueUsingHandle(t *testing.T) {
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
						nCreateNewSession := func(userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, disableAntiCsrf *bool, tenantId string, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							trueClaim, _ := TrueClaim()
							accessTokenPayload, err := trueClaim.Build(userID, tenantId, accessTokenPayload, userContext)
							if err != nil {
								return nil, err
							}
							return oCreateNewSession(userID, accessTokenPayload, sessionDataInDatabase, disableAntiCsrf, tenantId, userContext)
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

	res := fakeRes{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	sessionContainer, err := CreateNewSession(req, res, "public", "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)
	accessTokenPayload := sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 10, len(accessTokenPayload))
	assert.Equal(t, true, accessTokenPayload["st-true"].(map[string]interface{})["v"])

	trueClaim, _ := TrueClaim()
	ok, err := SetClaimValue(sessionContainer.GetHandle(), trueClaim, false)
	assert.NoError(t, err)
	assert.True(t, ok)

	sessInfo, err := GetSessionInformation(sessionContainer.GetHandle())
	assert.NoError(t, err)
	accessTokenPayload = sessInfo.CustomClaimsInAccessTokenPayload
	assert.Equal(t, 2, len(accessTokenPayload))
	assert.Equal(t, false, accessTokenPayload["st-true"].(map[string]interface{})["v"])
}

func TestShoulSetClaimValueWorkForNonExistantHandle(t *testing.T) {
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

	trueClaim, _ := TrueClaim()
	ok, err := SetClaimValue("invalidHandle", trueClaim, false)
	assert.NoError(t, err)
	assert.False(t, ok)
}
