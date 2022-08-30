package thirdparty

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func initForUserIdMappingTest(t *testing.T) {

	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{Init(&tpmodels.TypeInput{
			SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
				Providers: []tpmodels.TypeProvider{
					Google(tpmodels.GoogleConfig{ClientID: "clientID", ClientSecret: "clientSecret"}),
				},
			},
		})},
	}

	err := supertokens.Init(config)
	assert.NoError(t, err)
}

func TestCreateUserIdMappingUsingEmail(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	initForUserIdMappingTest(t)

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	assert.NoError(t, err)

	cdiVersion, err := querier.GetQuerierAPIVersion()
	assert.NoError(t, err)

	if unittesting.MaxVersion(cdiVersion, "2.14") != cdiVersion {
		return
	}

	signUpResponse, err := SignInUp("google", "googleID", tpmodels.EmailStruct{ID: "test@example.com"})
	assert.NoError(t, err)

	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	{ // Using supertokens ID
		userResp, err := GetUserByID(signUpResponse.OK.User.ID)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using external ID
		userResp, err := GetUserByID(externalUserId)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using thirdparty info
		userResp, err := GetUserByThirdPartyInfo("google", "googleID")
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}
}
