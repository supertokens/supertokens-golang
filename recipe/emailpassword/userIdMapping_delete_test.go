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

	if unittesting.MaxVersion(cdiVersion, "2.14") == "2.14" {
		return
	}

	{
		supertokensType := supertokens.UserIdTypeSupertokens
		deleteResp, err := supertokens.DeleteUserIdMapping("unknownUserId", &supertokensType, nil)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.False(t, deleteResp.OK.DidMappingExist)
	}

	{
		externalType := supertokens.UserIdTypeExternal
		deleteResp, err := supertokens.DeleteUserIdMapping("unknownUserId", &externalType, nil)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
		assert.False(t, deleteResp.OK.DidMappingExist)
	}

	{
		anyType := supertokens.UserIdTypeAny
		deleteResp, err := supertokens.DeleteUserIdMapping("unknownUserId", &anyType, nil)
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

	if unittesting.MaxVersion(cdiVersion, "2.14") == "2.14" {
		return
	}

	signUpResponse, err := SignUp("test@example.com", "testpass123", nil)
	assert.NoError(t, err)

	assert.NotNil(t, signUpResponse.OK)

	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	supertokensType := supertokens.UserIdTypeSupertokens
	deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, &supertokensType, nil)
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

	if unittesting.MaxVersion(cdiVersion, "2.14") == "2.14" {
		return
	}

	signUpResponse, err := SignUp("test@example.com", "testpass123", nil)
	assert.NoError(t, err)

	assert.NotNil(t, signUpResponse.OK)

	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	externalType := supertokens.UserIdTypeExternal
	deleteResp, err := supertokens.DeleteUserIdMapping(externalUserId, &externalType, nil)
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

	if unittesting.MaxVersion(cdiVersion, "2.14") == "2.14" {
		return
	}

	signUpResponse, err := SignUp("test@example.com", "testpass123", nil)
	assert.NoError(t, err)
	assert.NotNil(t, signUpResponse.OK)

	anyType := supertokens.UserIdTypeAny

	{
		externalUserId := "externalId"
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.OK)

		deleteResp, err := supertokens.DeleteUserIdMapping(externalUserId, &anyType, nil)
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
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.OK)

		deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, &anyType, nil)
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

	if unittesting.MaxVersion(cdiVersion, "2.14") == "2.14" {
		return
	}

	signUpResponse, err := SignUp("test@example.com", "testpass123", nil)
	assert.NoError(t, err)
	assert.NotNil(t, signUpResponse.OK)

	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	userMetadata := map[string]interface{}{
		"role": "admin",
	}
	metadataResp, err := usermetadata.UpdateUserMetadata(externalUserId, userMetadata, nil)
	assert.NoError(t, err)
	assert.NotNil(t, metadataResp)
	{ // with force nil
		deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, nil, nil)
		assert.Contains(t, err.Error(), "UserId is already in use in UserMetadata recipe")
		assert.Nil(t, deleteResp.OK)
	}

	{ // without force
		False := false
		deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, nil, &False)
		assert.Contains(t, err.Error(), "UserId is already in use in UserMetadata recipe")
		assert.Nil(t, deleteResp.OK)
	}

	{ // with force
		True := true
		deleteResp, err := supertokens.DeleteUserIdMapping(signUpResponse.OK.User.ID, nil, &True)
		assert.NoError(t, err)
		assert.NotNil(t, deleteResp.OK)
	}
}
