package emailpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDeleteUnknownUserIdMapping(t *testing.T) {
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
		deleteResp, err := supertokens.DeleteUserIdMapping("unknownUserId", "SUPERTOKENS", false)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.False(t, deleteResp.OK.DidMappingExist)
	}

	{
		deleteResp, err := supertokens.DeleteUserIdMapping("unknownUserId", "EXTERNAL", false)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.False(t, deleteResp.OK.DidMappingExist)
	}

	{
		deleteResp, err := supertokens.DeleteUserIdMapping("unknownUserId", "ANY", false)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.False(t, deleteResp.OK.DidMappingExist)
	}
}

func TestDeleteUserIdMappingOfSupertokensUserId(t *testing.T) {
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

	deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, "SUPERTOKENS", false)
	assert.NoError(t, err)
	assert.NotNil(t, deleteResp.OK)
	assert.True(t, deleteResp.OK.DidMappingExist)
}

func TestDeleteUserIdMappingOfExternalUserId(t *testing.T) {
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

	deleteResp, err := supertokens.DeleteUserIdMapping(externalUserId, "EXTERNAL", false)
	assert.NoError(t, err)
	assert.NotNil(t, deleteResp.OK)
	assert.True(t, deleteResp.OK.DidMappingExist)
}

func TestDeleteUserIdMappingOfAnyUsrId(t *testing.T) {
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

		deleteResp, err := supertokens.DeleteUserIdMapping(externalUserId, "ANY", false)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.True(t, deleteResp.OK.DidMappingExist)

		getResp, err := supertokens.GetUserIdMapping(externalUserId, "ANY")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.UnknownMappingError)
	}

	{
		externalUserId := "externalId"
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, false)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.OK)

		deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, "ANY", false)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.True(t, deleteResp.OK.DidMappingExist)

		getResp, err := supertokens.GetUserIdMapping(signUpResponse.OK.User.ID, "ANY")
		assert.NoError(t, err)
		assert.NotNil(t, getResp.UnknownMappingError)
	}
}
