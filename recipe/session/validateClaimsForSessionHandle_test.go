package session

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestValidateShouldReturnRightValidationErrors(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						oCreateNewSession := *originalImplementation.CreateNewSession
						nCreateNewSession := func(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							trueClaim, _ := TrueClaim()
							accessTokenPayload, err := trueClaim.Build(userID, accessTokenPayload, userContext)
							if err != nil {
								return sessmodels.SessionContainer{}, err
							}
							return oCreateNewSession(res, userID, accessTokenPayload, sessionData, userContext)
						}
						*originalImplementation.CreateNewSession = nCreateNewSession
						return originalImplementation
					},
				},
			}),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	res := fakeRes{}
	sessionContainer, err := CreateNewSession(res, "userId", map[string]interface{}{}, map[string]interface{}{})
	assert.NoError(t, err)

	_, nilValidator := NilClaim()
	failingValidator := nilValidator.HasValue(true, nil, nil)
	_, trueValidator := TrueClaim()
	passingValidator := trueValidator.IsTrue(nil, nil)

	validationRes, err := ValidateClaimsForSessionHandle(
		sessionContainer.GetHandle(),
		func(globalClaimValidators []claims.SessionClaimValidator, sessionInfo sessmodels.SessionInformation, userContext supertokens.UserContext) []claims.SessionClaimValidator {
			return []claims.SessionClaimValidator{
				passingValidator,
				failingValidator,
			}
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, validationRes.OK)
	assert.Equal(t, 1, len(validationRes.OK.InvalidClaims))
	assert.Equal(t, "st-nil", validationRes.OK.InvalidClaims[0].ID)
	assert.Equal(t, map[string]interface{}{
		"actualValue":   nil,
		"expectedValue": true,
		"message":       "value does not exist",
	}, validationRes.OK.InvalidClaims[0].Reason)
}

func TestValidateShouldWorkForNonExistantHandle(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	validationRes, err := ValidateClaimsForSessionHandle("nonExistantHandle", nil)
	assert.NoError(t, err)
	assert.Nil(t, validationRes.OK)
	assert.NotNil(t, validationRes.SessionDoesNotExistError)
}
