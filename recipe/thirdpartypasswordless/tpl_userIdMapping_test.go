package thirdpartypasswordless

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
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
		RecipeList: []supertokens.Recipe{Init(tplmodels.TypeInput{
			ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
				Enabled: true,
			},
			FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
			Providers: []tpmodels.ProviderInput{
				{
					Config: tpmodels.ProviderConfig{
						ThirdPartyId: "google",
						Clients: []tpmodels.ProviderClientConfig{
							{
								ClientID:     "clientID",
								ClientSecret: "clientSecret",
							},
						},
					},
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

	if unittesting.MaxVersion(cdiVersion, "2.14") == "2.14" {
		return
	}

	signUpResponse, err := ThirdPartyManuallyCreateOrUpdateUser("google", "googleID", "test@example.com", nil)
	assert.NoError(t, err)

	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.OK.User.ID, externalUserId, &externalUserIdInfo, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	{ // Using supertokens ID
		userResp, err := GetUserByID(signUpResponse.OK.User.ID, nil)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using external ID
		userResp, err := GetUserByID(externalUserId, nil)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using thirdparty info
		userResp, err := GetUserByThirdPartyInfo("google", "googleID", nil)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}
}

func TestPlessCreateUserIdMappingUsingEmail(t *testing.T) {
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

	signUpResponse, err := PasswordlessSignInUpByEmail("test@example.com", nil)
	assert.NoError(t, err)

	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.User.ID, externalUserId, &externalUserIdInfo, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	{ // Using supertokens ID
		userResp, err := GetUserByID(signUpResponse.User.ID, nil)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using external ID
		userResp, err := GetUserByID(externalUserId, nil)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using email
		userResp, err := GetUsersByEmail("test@example.com", nil)
		assert.NoError(t, err)
		assert.NotNil(t, userResp)
		assert.Equal(t, 1, len(userResp))
		for _, user := range userResp {
			assert.Equal(t, externalUserId, user.ID)
		}
	}

	{ // Using sign in
		codeResp, err := CreateCodeWithEmail("test@example.com", nil, nil)
		assert.NoError(t, err)
		assert.NotNil(t, codeResp.OK)

		resp, err := ConsumeCodeWithUserInputCode(codeResp.OK.DeviceID, codeResp.OK.UserInputCode, codeResp.OK.PreAuthSessionID, nil)
		assert.NoError(t, err)
		assert.NotNil(t, resp.OK)

		assert.Equal(t, externalUserId, resp.OK.User.ID)
	}
}

func TestPlessCreateUserIdMappingUsingPhone(t *testing.T) {
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

	signUpResponse, err := PasswordlessSignInUpByPhoneNumber("+919876543210", nil)
	assert.NoError(t, err)

	externalUserId := "externalId"
	externalUserIdInfo := "externalIdInfo"
	createResp, err := supertokens.CreateUserIdMapping(signUpResponse.User.ID, externalUserId, &externalUserIdInfo, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createResp.OK)

	{ // Using supertokens ID
		userResp, err := GetUserByID(signUpResponse.User.ID, nil)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using external ID
		userResp, err := GetUserByID(externalUserId, nil)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using email
		userResp, err := GetUserByPhoneNumber("+919876543210", nil)
		assert.NoError(t, err)
		assert.Equal(t, externalUserId, userResp.ID)
	}

	{ // Using sign in
		codeResp, err := CreateCodeWithPhoneNumber("+919876543210", nil, nil)
		assert.NoError(t, err)
		assert.NotNil(t, codeResp.OK)

		resp, err := ConsumeCodeWithUserInputCode(codeResp.OK.DeviceID, codeResp.OK.UserInputCode, codeResp.OK.PreAuthSessionID, nil)
		assert.NoError(t, err)
		assert.NotNil(t, resp.OK)

		assert.Equal(t, externalUserId, resp.OK.User.ID)
	}
}
