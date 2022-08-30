package emailpassword

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestCreateUserIdMappingAndDeleteUser(t *testing.T) {
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
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	err = supertokens.DeleteUser(externalUserId)
	assert.NoError(t, err)
}

func TestCreateUserIdMappingAndGetUsers(t *testing.T) {
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

	for i := 0; i < 4; i++ {
		signUpResponse, err := SignUp(fmt.Sprintf("test%d@example.com", i), "testpass123")
		assert.NoError(t, err)

		assert.NotNil(t, signUpResponse.OK)

		externalUserId := fmt.Sprintf("externalId%d", i)
		externalUserIdInfo := "externalIdInfo"
		createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
		assert.NoError(t, err)
		assert.NotNil(t, createResp.OK)
	}

	userResult, err := supertokens.GetUsersOldestFirst(nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, userResult.Users)

	for i, user := range userResult.Users {
		assert.Equal(t, user.User["id"], fmt.Sprintf("externalId%d", i))
	}
}
