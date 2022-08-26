package emailpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		RecipeList: []supertokens.Recipe{Init(nil)},
	}

	err := supertokens.Init(config)
	assert.NoError(t, err)
}

func TestCreateUserIdMapping(t *testing.T) {
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

	signUpResponse, err := SignUp("test@example.com", "testpass123")
	assert.NoError(t, err)

	assert.NotNil(t, signUpResponse.OK)

	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, false)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	supertokensType := "SUPERTOKENS"
	getResp, err := supertokens.GetUserIdMapping(signUpResponse.OK.User.ID, &supertokensType)
	assert.NoError(t, err)
	assert.NotNil(t, getResp.OK)
	assert.Equal(t, signUpResponse.OK.User.ID, getResp.OK.SupertokensUserId)
	assert.Equal(t, externalUserId, getResp.OK.ExternalUserId)
	assert.Equal(t, externalUserIdInfo, *getResp.OK.ExternalUserIdInfo)
}

func TestCreateUserIdMappingWithUnknownSupertokensUserId(t *testing.T) {
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

	supertokensUserId := "unknownUserId"
	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"

	createResp, err := supertokens.CreateUserIdMapping(supertokensUserId, externalUserId, &externalUserIdInfo, false)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.UnknownSupertokensUserIdError)
}

func TestCreateUserIdMappingWhenAlreadyExists(t *testing.T) {
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

	signUpResponse, err := SignUp("test@example.com", "testpass123")
	assert.NoError(t, err)
	assert.NotNil(t, signUpResponse.OK)

	{
		externalUserId := "externalId"
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, false)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.OK)
	}

	{ // duplicate of both
		externalUserId := "externalId"
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, false)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.UserIdMappingAlreadyExistsError)
		assert.True(t, createResp.UserIdMappingAlreadyExistsError.DoesExternalUserIdExist)
		assert.True(t, createResp.UserIdMappingAlreadyExistsError.DoesSuperTokensUserIdExist)
	}

	{ // duplicate of supertokensUserId
		externalUserId := "differentId"
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, false)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.UserIdMappingAlreadyExistsError)
		assert.False(t, createResp.UserIdMappingAlreadyExistsError.DoesExternalUserIdExist)
		assert.True(t, createResp.UserIdMappingAlreadyExistsError.DoesSuperTokensUserIdExist)
	}

	{ // duplicate of externalUserId

		signUpResponse, err := SignUp("test2@example.com", "testpass123")
		assert.NoError(t, err)
		assert.NotNil(t, signUpResponse.OK)

		externalUserId := "externalId"
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, false)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.UserIdMappingAlreadyExistsError)
		assert.True(t, createResp.UserIdMappingAlreadyExistsError.DoesExternalUserIdExist)
		assert.False(t, createResp.UserIdMappingAlreadyExistsError.DoesSuperTokensUserIdExist)
	}
}
