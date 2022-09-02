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

	res := fakeRes{}
	sessionContainer, err := CreateNewSession(res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	sessionContainer.SetClaimValue(TrueClaim(), true)

	accessTokenPayload := sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 1, len(accessTokenPayload))
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
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						oCreateNewSession := *originalImplementation.CreateNewSession
						nCreateNewSession := func(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							accessTokenPayloadUpdate, err := TrueClaim().Build(userID, userContext)
							if err != nil {
								return sessmodels.SessionContainer{}, err
							}
							for k, v := range accessTokenPayloadUpdate {
								accessTokenPayload[k] = v
							}
							return oCreateNewSession(res, userID, accessTokenPayload, sessionData, userContext)
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
	sessionContainer, err := CreateNewSession(res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)
	accessTokenPayload := sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 1, len(accessTokenPayload))
	assert.Equal(t, true, accessTokenPayload["st-true"].(map[string]interface{})["v"])

	sessionContainer.SetClaimValue(TrueClaim(), false)

	accessTokenPayload = sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 1, len(accessTokenPayload))
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
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						oCreateNewSession := *originalImplementation.CreateNewSession
						nCreateNewSession := func(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							accessTokenPayloadUpdate, err := TrueClaim().Build(userID, userContext)
							if err != nil {
								return sessmodels.SessionContainer{}, err
							}
							for k, v := range accessTokenPayloadUpdate {
								accessTokenPayload[k] = v
							}
							return oCreateNewSession(res, userID, accessTokenPayload, sessionData, userContext)
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
	sessionContainer, err := CreateNewSession(res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)
	accessTokenPayload := sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 1, len(accessTokenPayload))
	assert.Equal(t, true, accessTokenPayload["st-true"].(map[string]interface{})["v"])

	ok, err := SetClaimValue(sessionContainer.GetHandle(), TrueClaim(), false)
	assert.NoError(t, err)
	assert.True(t, ok)

	sessInfo, err := GetSessionInformation(sessionContainer.GetHandle())
	assert.NoError(t, err)
	accessTokenPayload = sessInfo.AccessTokenPayload
	assert.Equal(t, 1, len(accessTokenPayload))
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

	ok, err := SetClaimValue("invalidHandle", TrueClaim(), false)
	assert.NoError(t, err)
	assert.False(t, ok)
}
