package emailpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestGetUserIdMapping(t *testing.T) {
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

	{
		getResp, err := supertokens.GetUserIdMapping(signUpResponse.OK.User.ID, "SUPERTOKENS")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, signUpResponse.OK.User.ID, getResp.OK.SupertokensUserId)
		assert.Equal(t, externalUserId, getResp.OK.ExternalUserId)
		assert.Equal(t, externalUserIdInfo, *getResp.OK.ExternalUserIdInfo)
	}

	{
		getResp, err := supertokens.GetUserIdMapping(externalUserId, "EXTERNAL")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, signUpResponse.OK.User.ID, getResp.OK.SupertokensUserId)
		assert.Equal(t, externalUserId, getResp.OK.ExternalUserId)
		assert.Equal(t, externalUserIdInfo, *getResp.OK.ExternalUserIdInfo)
	}

	{
		getResp, err := supertokens.GetUserIdMapping(signUpResponse.OK.User.ID, "ANY")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, signUpResponse.OK.User.ID, getResp.OK.SupertokensUserId)
		assert.Equal(t, externalUserId, getResp.OK.ExternalUserId)
		assert.Equal(t, externalUserIdInfo, *getResp.OK.ExternalUserIdInfo)
	}

	{
		getResp, err := supertokens.GetUserIdMapping(externalUserId, "ANY")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, signUpResponse.OK.User.ID, getResp.OK.SupertokensUserId)
		assert.Equal(t, externalUserId, getResp.OK.ExternalUserId)
		assert.Equal(t, externalUserIdInfo, *getResp.OK.ExternalUserIdInfo)
	}
}

func TestGetUserIdMappingThatDoesNotExist(t *testing.T) {
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

	{
		getResp, err := supertokens.GetUserIdMapping("unknownId", "ANY")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.UnknownMappingError)
	}

	{
		getResp, err := supertokens.GetUserIdMapping("unknownId", "SUPERTOKENS")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.UnknownMappingError)
	}

	{
		getResp, err := supertokens.GetUserIdMapping("unknownId", "EXTERNAL")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.UnknownMappingError)
	}
}

func TestGetUserIdMappingWithNoExternalUserIdInfo(t *testing.T) {
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
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, nil, false)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	{
		getResp, err := supertokens.GetUserIdMapping(signUpResponse.OK.User.ID, "SUPERTOKENS")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, signUpResponse.OK.User.ID, getResp.OK.SupertokensUserId)
		assert.Equal(t, externalUserId, getResp.OK.ExternalUserId)
		assert.Nil(t, getResp.OK.ExternalUserIdInfo)
	}

	{
		getResp, err := supertokens.GetUserIdMapping(externalUserId, "EXTERNAL")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, signUpResponse.OK.User.ID, getResp.OK.SupertokensUserId)
		assert.Equal(t, externalUserId, getResp.OK.ExternalUserId)
		assert.Nil(t, getResp.OK.ExternalUserIdInfo)
	}

	{
		getResp, err := supertokens.GetUserIdMapping(signUpResponse.OK.User.ID, "ANY")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, signUpResponse.OK.User.ID, getResp.OK.SupertokensUserId)
		assert.Equal(t, externalUserId, getResp.OK.ExternalUserId)
		assert.Nil(t, getResp.OK.ExternalUserIdInfo)
	}

	{
		getResp, err := supertokens.GetUserIdMapping(externalUserId, "ANY")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, signUpResponse.OK.User.ID, getResp.OK.SupertokensUserId)
		assert.Equal(t, externalUserId, getResp.OK.ExternalUserId)
		assert.Nil(t, getResp.OK.ExternalUserIdInfo)
	}
}
