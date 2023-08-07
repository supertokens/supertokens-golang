package emailpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestUpdateUserIdMappingForUnknownUserId(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	initForUserIdMappingTest(t)

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	assert.NoError(t, err)

	cdiVersion, err := querier.GetQuerierAPIVersion()
	assert.NoError(t, err)

	if unittesting.MaxVersion(cdiVersion, "2.14") == "2.14" {
		return
	}

	externalInfo := "someInfo"
	{
		supertokensType := supertokens.UserIdTypeSupertokens
		updateResp, err := supertokens.UpdateOrDeleteUserIdMappingInfo("unknownUserId", &supertokensType, &externalInfo)
		assert.NoError(t, err)
		assert.NotNil(t, updateResp.UnknownMappingError)
	}

	{
		externalType := supertokens.UserIdTypeExternal
		updateResp, err := supertokens.UpdateOrDeleteUserIdMappingInfo("unknownUserId", &externalType, &externalInfo)
		assert.NoError(t, err)
		assert.NotNil(t, updateResp.UnknownMappingError)
	}

	{
		anyType := supertokens.UserIdTypeAny
		updateResp, err := supertokens.UpdateOrDeleteUserIdMappingInfo("unknownUserId", &anyType, &externalInfo)
		assert.NoError(t, err)
		assert.NotNil(t, updateResp.UnknownMappingError)
	}
}

func TestUpdateInfoInUserIdMapping(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	initForUserIdMappingTest(t)

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	assert.NoError(t, err)

	cdiVersion, err := querier.GetQuerierAPIVersion()
	assert.NoError(t, err)

	if unittesting.MaxVersion(cdiVersion, "2.14") == "2.14" {
		return
	}

	signUpResponse, err := SignUp("public", "test@example.com", "testpass123")
	assert.NoError(t, err)
	assert.NotNil(t, signUpResponse.OK)

	{
		externalUserId := "externalId"
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.OK)
	}

	{
		infoA := "infoA"
		supertokensType := supertokens.UserIdTypeSupertokens
		updateResp, err := supertokens.UpdateOrDeleteUserIdMappingInfo(signUpResponse.OK.User.ID, &supertokensType, &infoA)
		assert.NoError(t, err)
		assert.NotNil(t, updateResp.OK)

		getResp, err := supertokens.GetUserIdMapping(signUpResponse.OK.User.ID, &supertokensType)
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, infoA, *getResp.OK.ExternalUserIdInfo)
	}

	{
		infoB := "infoB"
		externalType := supertokens.UserIdTypeExternal
		updateResp, err := supertokens.UpdateOrDeleteUserIdMappingInfo("externalId", &externalType, &infoB)
		assert.NoError(t, err)
		assert.NotNil(t, updateResp.OK)

		getResp, err := supertokens.GetUserIdMapping("externalId", &externalType)
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, infoB, *getResp.OK.ExternalUserIdInfo)
	}

	{
		infoC := "infoC"
		anyType := supertokens.UserIdTypeAny
		updateResp, err := supertokens.UpdateOrDeleteUserIdMappingInfo(signUpResponse.OK.User.ID, &anyType, &infoC)
		assert.NoError(t, err)
		assert.NotNil(t, updateResp.OK)

		getResp, err := supertokens.GetUserIdMapping(signUpResponse.OK.User.ID, &anyType)
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, infoC, *getResp.OK.ExternalUserIdInfo)
	}

	{
		infoD := "infoD"
		anyType := supertokens.UserIdTypeAny
		updateResp, err := supertokens.UpdateOrDeleteUserIdMappingInfo("externalId", &anyType, &infoD)
		assert.NoError(t, err)
		assert.NotNil(t, updateResp.OK)

		getResp, err := supertokens.GetUserIdMapping("externalId", &anyType)
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Equal(t, infoD, *getResp.OK.ExternalUserIdInfo)
	}
}

func TestDeleteInfoInUserIdMapping(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	initForUserIdMappingTest(t)

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	assert.NoError(t, err)

	cdiVersion, err := querier.GetQuerierAPIVersion()
	assert.NoError(t, err)

	if unittesting.MaxVersion(cdiVersion, "2.14") == "2.14" {
		return
	}

	signUpResponse, err := SignUp("public", "test@example.com", "testpass123")
	assert.NoError(t, err)
	assert.NotNil(t, signUpResponse.OK)

	{
		externalUserId := "externalId"
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.OK)
	}

	{ // set external Info to nil
		supertokensType := supertokens.UserIdTypeSupertokens
		updateResp, err := supertokens.UpdateOrDeleteUserIdMappingInfo(signUpResponse.OK.User.ID, &supertokensType, nil)
		assert.NoError(t, err)
		assert.NotNil(t, updateResp.OK)

		getResp, err := supertokens.GetUserIdMapping(signUpResponse.OK.User.ID, &supertokensType)
		assert.NoError(t, err)
		assert.NotNil(t, getResp.OK)
		assert.Nil(t, getResp.OK.ExternalUserIdInfo)
	}
}
