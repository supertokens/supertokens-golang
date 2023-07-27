package passwordless

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
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
		RecipeList: []supertokens.Recipe{Init(plessmodels.TypeInput{
			FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
			ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
				Enabled: true,
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

	if unittesting.MaxVersion(cdiVersion, "2.14") == "2.14" {
		return
	}

	signUpResponse, err := SignInUpByEmail("public", "test@example.com")
	assert.NoError(t, err)

	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.User.ID, externalUserId, &externalUserIdInfo, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	{ // Using supertokens ID
		userResp, err := GetUserByID(signUpResponse.User.ID)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using external ID
		userResp, err := GetUserByID(externalUserId)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using email
		userResp, err := GetUserByEmail("public", "test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using sign in
		codeResp, err := CreateCodeWithEmail("public", "test@example.com", nil)
		assert.NoError(t, err)
		assert.NotNil(t, codeResp.OK)

		resp, err := ConsumeCodeWithUserInputCode("public", codeResp.OK.DeviceID, codeResp.OK.UserInputCode, codeResp.OK.PreAuthSessionID)
		assert.NoError(t, err)
		assert.NotNil(t, resp.OK)

		assert.Equal(t, externalUserId, resp.OK.User.ID)
	}
}

func TestCreateUserIdMappingUsingPhone(t *testing.T) {
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

	signUpResponse, err := SignInUpByPhoneNumber("public", "+919876543210")
	assert.NoError(t, err)

	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.User.ID, externalUserId, &externalUserIdInfo, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	{ // Using supertokens ID
		userResp, err := GetUserByID(signUpResponse.User.ID)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using external ID
		userResp, err := GetUserByID(externalUserId)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using email
		userResp, err := GetUserByPhoneNumber("public", "+919876543210")
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using sign in
		codeResp, err := CreateCodeWithPhoneNumber("public", "+919876543210", nil)
		assert.NoError(t, err)
		assert.NotNil(t, codeResp.OK)

		resp, err := ConsumeCodeWithUserInputCode("public", codeResp.OK.DeviceID, codeResp.OK.UserInputCode, codeResp.OK.PreAuthSessionID)
		assert.NoError(t, err)
		assert.NotNil(t, resp.OK)

		assert.Equal(t, externalUserId, resp.OK.User.ID)
	}
}
