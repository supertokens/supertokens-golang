package session

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestGetClaimValueReturnsRightValue(t *testing.T) {
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
						nCreateNewSession := func(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, tenantId *string, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							trueClaim, _ := TrueClaim()
							accessTokenPayload, err := trueClaim.Build(userID, accessTokenPayload, tenantId, userContext)
							if err != nil {
								return nil, err
							}
							return oCreateNewSession(req, res, userID, accessTokenPayload, sessionData, tenantId, userContext)
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
	sessionContainer, err := CreateNewSession(req, res, "userId", map[string]interface{}{}, map[string]interface{}{}, nil)
	assert.NoError(t, err)

	trueClaim, _ := TrueClaim()
	assert.True(t, sessionContainer.GetClaimValue(trueClaim).(bool))
}

func TestGetClaimValueFromHandleReturnsRightValue(t *testing.T) {
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
						nCreateNewSession := func(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, tenantId *string, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							trueClaim, _ := TrueClaim()
							accessTokenPayload, err := trueClaim.Build(userID, accessTokenPayload, tenantId, userContext)
							if err != nil {
								return nil, err
							}
							return oCreateNewSession(req, res, userID, accessTokenPayload, sessionData, tenantId, userContext)
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
	sessionContainer, err := CreateNewSession(req, res, "userId", map[string]interface{}{}, map[string]interface{}{}, nil)
	assert.NoError(t, err)

	trueClaim, _ := TrueClaim()
	getRes, err := GetClaimValue(sessionContainer.GetHandle(), trueClaim, sessionContainer.GetTenantId())
	assert.NoError(t, err)
	assert.NotNil(t, getRes.OK)
	assert.True(t, getRes.OK.Value.(bool))
}

func TestGetClaimValueForNonExistantSessionhandle(t *testing.T) {
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
						nCreateNewSession := func(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, tenantId *string, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							trueClaim, _ := TrueClaim()
							accessTokenPayload, err := trueClaim.Build(userID, accessTokenPayload, tenantId, userContext)
							if err != nil {
								return nil, err
							}
							return oCreateNewSession(req, res, userID, accessTokenPayload, sessionData, tenantId, userContext)
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

	trueClaim, _ := TrueClaim()
	getRes, err := GetClaimValue("invalidSessionHandle", trueClaim, nil)
	assert.NoError(t, err)
	assert.Nil(t, getRes.OK)
	assert.NotNil(t, getRes.SessionDoesNotExistError)
}
