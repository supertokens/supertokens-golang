package emailpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
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
		supertokensType := "SUPERTOKENS"
		deleteResp, err := supertokens.DeleteUserIdMapping("unknownUserId", &supertokensType, false)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.False(t, deleteResp.OK.DidMappingExist)
	}

	{
		externalType := "EXTERNAL"
		deleteResp, err := supertokens.DeleteUserIdMapping("unknownUserId", &externalType, false)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.False(t, deleteResp.OK.DidMappingExist)
	}

	{
		anyType := "ANY"
		deleteResp, err := supertokens.DeleteUserIdMapping("unknownUserId", &anyType, false)
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

	supertokensType := "SUPERTOKENS"
	deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, &supertokensType, false)
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

	externalType := "EXTERNAL"
	deleteResp, err := supertokens.DeleteUserIdMapping(externalUserId, &externalType, false)
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

	anyType := "ANY"

	{
		externalUserId := "externalId"
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, false)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.OK)

		deleteResp, err := supertokens.DeleteUserIdMapping(externalUserId, &anyType, false)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.True(t, deleteResp.OK.DidMappingExist)

		getResp, err := supertokens.GetUserIdMapping(externalUserId, &anyType)
		assert.NoError(t, err)
		assert.NotNil(t, getResp.UnknownMappingError)
	}

	{
		externalUserId := "externalId"
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, false)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.OK)

		deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, &anyType, false)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.True(t, deleteResp.OK.DidMappingExist)

		getResp, err := supertokens.GetUserIdMapping(signUpResponse.OK.User.ID, &anyType)
		assert.NoError(t, err)
		assert.NotNil(t, getResp.UnknownMappingError)
	}
}

func TestDeleteUserIdMappingWithMetadataAndWithAndWithoutForce(t *testing.T) {
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

	userMetadata := map[string]interface{}{
		"role": "admin",
	}
	metadataResp, err := usermetadata.UpdateUserMetadata(externalUserId, userMetadata)
	assert.NoError(t, err)
	assert.NotNil(t, metadataResp)
	{ // without force
		deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, nil, false)
		assert.Contains(t, err.Error(), "UserId is already in use in UserMetadata recipe")
		assert.Nil(t, deleteResp.OK)
	}

	{ // without force
		deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, nil, true)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
	}
}
