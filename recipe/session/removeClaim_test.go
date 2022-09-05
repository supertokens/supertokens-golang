package session

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestShouldRemoveNonExistantClaim(t *testing.T) {
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

	trueClaim, _ := TrueClaim()
	err = sessionContainer.RemoveClaim(trueClaim)
	assert.NoError(t, err)
}

func TestShouldClearPreviouslySetClaim(t *testing.T) {
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
							trueClaim, _ := TrueClaim()
							err := trueClaim.Build(userID, accessTokenPayload, userContext)
							if err != nil {
								return sessmodels.SessionContainer{}, err
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

	trueClaim, _ := TrueClaim()
	err = sessionContainer.RemoveClaim(trueClaim)
	assert.NoError(t, err)

	accessTokenPayload = sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 0, len(accessTokenPayload))
}

func TestShouldClearPreviouslySetClaimUsingHandle(t *testing.T) {
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
							trueClaim, _ := TrueClaim()
							err := trueClaim.Build(userID, accessTokenPayload, userContext)
							if err != nil {
								return sessmodels.SessionContainer{}, err
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

	trueClaim, _ := TrueClaim()
	ok, err := RemoveClaim(sessionContainer.GetHandle(), trueClaim)
	assert.True(t, ok)
	assert.NoError(t, err)

	sessInfo, err := GetSessionInformation(sessionContainer.GetHandle())
	assert.NoError(t, err)
	accessTokenPayload = sessInfo.AccessTokenPayload
	assert.Equal(t, 0, len(accessTokenPayload))
}

func TestShouldRemoveWorkForNonExistantHandle(t *testing.T) {
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

	trueClaim, _ := TrueClaim()
	ok, err := RemoveClaim("invalidHandle", trueClaim)
	assert.False(t, ok)
	assert.NoError(t, err)
}
