package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestShouldAllowWithoutClaimsRequiredOrPresent(t *testing.T) {
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
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

	app := getTestApp([]typeTestEndpoint{})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/default-claims", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
}

func TestShouldAllowClaimValidAfterRefetching(t *testing.T) {
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							_, trueValidators := TrueClaim()

							result = append(result, trueValidators.HasValue(true, nil, nil))
							return result, nil
						}
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

	app := getTestApp([]typeTestEndpoint{})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/default-claims", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
}

func TestShouldRejectClaimRequiredButNotAdded(t *testing.T) {
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							_, nilValidators := NilClaim()

							result = append(result, nilValidators.HasValue(true, nil, nil))
							return result, nil
						}
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

	app := getTestApp([]typeTestEndpoint{})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/default-claims", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 403, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
	assert.Equal(t, "invalid claim", res["message"].(string))
	assert.Equal(t, []interface{}{
		map[string]interface{}{
			"id": "st-nil",
			"reason": map[string]interface{}{
				"actualValue":   nil,
				"expectedValue": true,
				"message":       "value does not exist",
			},
		},
	}, res["claimValidationErrors"])
}

func TestShouldAllowCustomValidatorReturningTrue(t *testing.T) {
	customValidator := claims.SessionClaimValidator{
		ID: "testid",
		Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
			return claims.ClaimValidationResult{
				IsValid: true,
			}
		},
	}
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							result = append(result, customValidator)
							return result, nil
						}
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

	app := getTestApp([]typeTestEndpoint{})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/default-claims", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
}

func TestShouldRejectCustomValidatorReturningFalse(t *testing.T) {
	customValidator := claims.SessionClaimValidator{
		ID: "testid",
		Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
			return claims.ClaimValidationResult{
				IsValid: false,
			}
		},
	}
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							result = append(result, customValidator)
							return result, nil
						}
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

	app := getTestApp([]typeTestEndpoint{})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/default-claims", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 403, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
	assert.Equal(t, "invalid claim", res["message"].(string))
	assert.Equal(t, []interface{}{
		map[string]interface{}{
			"id": "testid",
		},
	}, res["claimValidationErrors"])
}

func TestShouldRejectCustomValidatorReturningFalseWithReason(t *testing.T) {
	customValidator := claims.SessionClaimValidator{
		ID: "testid",
		Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
			return claims.ClaimValidationResult{
				IsValid: false,
				Reason:  "custom reason",
			}
		},
	}
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							result = append(result, customValidator)
							return result, nil
						}
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

	app := getTestApp([]typeTestEndpoint{})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/default-claims", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 403, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
	assert.Equal(t, "invalid claim", res["message"].(string))
	assert.Equal(t, []interface{}{
		map[string]interface{}{
			"id":     "testid",
			"reason": "custom reason",
		},
	}, res["claimValidationErrors"])
}

func TestShouldRejectIfAssertClaimsReturnsError(t *testing.T) {
	customValidator := claims.SessionClaimValidator{
		ID: "testid",
		Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
			return claims.ClaimValidationResult{
				IsValid: true,
			}
		},
	}
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							result = append(result, customValidator)
							return result, nil
						}
						*originalImplementation.ValidateClaims = func(userId string, accessTokenPayload map[string]interface{}, claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) (sessmodels.ValidateClaimsResult, error) {
							return sessmodels.ValidateClaimsResult{
								InvalidClaims: []claims.ClaimValidationError{
									{
										ID:     "testid",
										Reason: "custom reason",
									},
								},
							}, nil
						}
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

	app := getTestApp([]typeTestEndpoint{})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/default-claims", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 403, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
	assert.Equal(t, "invalid claim", res["message"].(string))
	assert.Equal(t, []interface{}{
		map[string]interface{}{
			"id":     "testid",
			"reason": "custom reason",
		},
	}, res["claimValidationErrors"])
}

func TestShouldAllowIfAssertClaimsReturnsNoError(t *testing.T) {
	customValidator := claims.SessionClaimValidator{
		ID: "testid",
		Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
			return claims.ClaimValidationResult{
				IsValid: true,
			}
		},
	}
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							result = append(result, customValidator)
							return result, nil
						}
						*originalImplementation.ValidateClaims = func(userId string, accessTokenPayload map[string]interface{}, claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) (sessmodels.ValidateClaimsResult, error) {
							return sessmodels.ValidateClaimsResult{
								InvalidClaims: []claims.ClaimValidationError{},
							}, nil
						}
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

	app := getTestApp([]typeTestEndpoint{})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/default-claims", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
	assert.NotEqual(t, "invalid claim", res["message"].(string))
	assert.Nil(t, res["claimValidationErrors"])
}

func TestShouldAllowWithEmptyListAsOverride(t *testing.T) {
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							_, nilValidators := NilClaim()

							result = append(result, nilValidators.HasValue(true, nil, nil))
							return result, nil
						}
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

	app := getTestApp([]typeTestEndpoint{
		{
			path: "/no-claims",
			overrideGlobalClaimValidators: func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
				return []claims.SessionClaimValidator{}, nil
			},
		},
	})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/no-claims", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
	assert.NotEqual(t, "invalid claim", res["message"].(string))
	assert.Nil(t, res["claimValidationErrors"])
}

func TestShouldAllowClaimValidAfterRefetchingWithOverride(t *testing.T) {
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							_, trueValidators := TrueClaim()

							result = append(result, trueValidators.HasValue(false, nil, nil))
							return result, nil
						}
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

	app := getTestApp([]typeTestEndpoint{
		{
			path: "/refetched-claim",
			overrideGlobalClaimValidators: func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
				_, validators := TrueClaim()
				return []claims.SessionClaimValidator{
					validators.HasValue(true, nil, nil),
				}, nil
			},
		},
	})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/refetched-claim", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
}

func TestShouldRejectClaimInvalidAfterRefetchingWithOverride(t *testing.T) {
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							_, trueValidators := TrueClaim()

							result = append(result, trueValidators.HasValue(true, nil, nil))
							return result, nil
						}
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

	app := getTestApp([]typeTestEndpoint{
		{
			path: "/refetched-claim",
			overrideGlobalClaimValidators: func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
				_, validators := TrueClaim()
				return []claims.SessionClaimValidator{
					validators.HasValue(false, nil, nil),
				}, nil
			},
		},
	})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/refetched-claim", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 403, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
	assert.Equal(t, "invalid claim", res["message"].(string))
	assert.Equal(t, []interface{}{
		map[string]interface{}{
			"id": "st-true",
			"reason": map[string]interface{}{
				"actualValue":   true,
				"expectedValue": false,
				"message":       "wrong value",
			},
		},
	}, res["claimValidationErrors"])
}

func TestShouldRejectCustomValidatorReturningFalseWithOverride(t *testing.T) {
	customValidator := claims.SessionClaimValidator{
		ID: "testid",
		Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
			return claims.ClaimValidationResult{
				IsValid: false,
			}
		},
	}
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							_, validator := TrueClaim()
							result = append(result, validator.IsTrue(nil, nil))
							return result, nil
						}
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

	app := getTestApp([]typeTestEndpoint{
		{
			path: "/refetched-claim",
			overrideGlobalClaimValidators: func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
				return []claims.SessionClaimValidator{
					customValidator,
				}, nil
			},
		},
	})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/refetched-claim", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 403, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
	assert.Equal(t, "invalid claim", res["message"].(string))
	assert.Equal(t, []interface{}{
		map[string]interface{}{
			"id": "testid",
		},
	}, res["claimValidationErrors"])
}

func TestShouldAllowCustomValidatorReturningTrueWithOverride(t *testing.T) {
	customValidator := claims.SessionClaimValidator{
		ID: "testid",
		Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
			return claims.ClaimValidationResult{
				IsValid: true,
			}
		},
	}
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
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						*originalImplementation.GetGlobalClaimValidators = func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
							result := []claims.SessionClaimValidator{}
							result = append(result, claimValidatorsAddedByOtherRecipes...)
							_, validator := TrueClaim()

							result = append(result, validator.IsFalse(nil, nil))
							return result, nil
						}
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

	app := getTestApp([]typeTestEndpoint{
		{
			path: "/refetched-claim",
			overrideGlobalClaimValidators: func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
				return []claims.SessionClaimValidator{
					customValidator,
				}, nil
			},
		},
	})
	defer app.Close()

	cookies := createSession(app, nil)

	code, res, err := unittesting.GetRequestWithJSONResult(app.URL+"/refetched-claim", cookies)
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res["message"].(string))
}

/**
This test is to make sure that we dont always call the core for session verification by default.
1. Create a session
2. Call get session and expect to not call the core
3. Call refresh session
4. Call get session and expect the core to be called
*/
func TestThatVerifySessionDoesNotAlwaysCallCore(t *testing.T) {
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

	createSessionResp, err := CreateNewSessionWithoutRequestResponse("test-user-id", nil, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	createdTokens := createSessionResp.GetAllSessionTokensDangerously()
	assert.True(t, createdTokens.AccessToken != "")
	assert.True(t, *createdTokens.RefreshToken != "")
	assert.True(t, createdTokens.FrontToken != "")

	getSessionResp, err := GetSessionWithoutRequestResponse(createdTokens.AccessToken, createdTokens.AntiCsrfToken, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.False(t, didGetSessionCallCore)
	newTokens := getSessionResp.GetAllSessionTokensDangerously()
	assert.True(t, newTokens.AccessToken != "")
	assert.True(t, newTokens.RefreshToken == nil)
	assert.True(t, newTokens.FrontToken != "")

	refreshResp, err := RefreshSessionWithoutRequestResponse(*createdTokens.RefreshToken, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.False(t, didGetSessionCallCore)
	tokensAfterRefresh := refreshResp.GetAllSessionTokensDangerously()
	assert.True(t, tokensAfterRefresh.AccessToken != "")
	assert.True(t, *tokensAfterRefresh.RefreshToken != "")
	assert.True(t, tokensAfterRefresh.FrontToken != "")

	getSessionResp, err = GetSessionWithoutRequestResponse(tokensAfterRefresh.AccessToken, tokensAfterRefresh.AntiCsrfToken, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.True(t, didGetSessionCallCore)
}

func TestThatVerifySessionReturns401IfNoAccessTokenIsSentAndMiddlewareIsNotAdded(t *testing.T) {
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

	testServer := getTestServerWithoutMiddleware()
	bodyBytes := []byte("{}")
	res, err := http.Post(testServer.URL+"/verify", "application/json", bytes.NewBuffer(bodyBytes))

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, res.StatusCode, 401)
}

func TestThatAntiCSRFCheckIsSkippedIfSessionRequiredIsFalseAndNoAccessTokenIsPassed(t *testing.T) {
	AntiCsrf := AntiCSRF_VIA_CUSTOM_HEADER
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
				AntiCsrf: &AntiCsrf,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
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

	app := getTestApp([]typeTestEndpoint{})
	defer app.Close()

	session, err := CreateNewSessionWithoutRequestResponse("test-user", map[string]interface{}{}, map[string]interface{}{}, nil)
	assert.NoError(t, err)

	sessionTokens := session.GetAllSessionTokensDangerously()

	req, err := http.NewRequest(http.MethodGet, app.URL+"/verify", nil)
	assert.NoError(t, err)
	req.Header.Add("Cookie", "sAccessToken="+*sessionTokens.RefreshToken)

	res, err := http.DefaultClient.Do(req)
	assert.Equal(t, res.StatusCode, 401)

	req, err = http.NewRequest(http.MethodPost, app.URL+"/verify-optional", nil)
	assert.NoError(t, err)

	res, err = http.DefaultClient.Do(req)
	assert.Equal(t, res.StatusCode, 200)
}

type typeTestEndpoint struct {
	path                          string
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error)
}

func createSession(app *httptest.Server, body map[string]interface{}) []*http.Cookie {
	bodyBytes := []byte("{}")
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}
	res, err := http.Post(app.URL+"/create", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil
	}
	return res.Cookies()
}

func getTestApp(endpoints []typeTestEndpoint) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(400)
			return
		}
		body := map[string]interface{}{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		CreateNewSession(r, w, "testing-userId", body, map[string]interface{}{})
		resp := map[string]interface{}{
			"message": true,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", (len(respBytes))))
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	})

	False := false
	mux.HandleFunc("/verify-optional", VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &False,
	}, func(rw http.ResponseWriter, r *http.Request) {
		GetSession(r, rw, &sessmodels.VerifySessionOptions{
			SessionRequired: &False,
		})
	}))

	mux.HandleFunc("/verify", VerifySession(&sessmodels.VerifySessionOptions{}, func(rw http.ResponseWriter, r *http.Request) {
		GetSession(r, rw, &sessmodels.VerifySessionOptions{})
	}))

	mux.HandleFunc("/default-claims", VerifySession(nil, func(w http.ResponseWriter, r *http.Request) {
		sessionContainer := GetSessionFromRequestContext(r.Context())
		resp := map[string]interface{}{
			"message": sessionContainer.GetHandle(),
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", (len(respBytes))))
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))

	mux.HandleFunc("/logout", VerifySession(nil, func(w http.ResponseWriter, r *http.Request) {
		sessionContainer, err := GetSession(r, w, nil)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		RevokeSession((*sessionContainer).GetHandle())
		resp := map[string]interface{}{
			"message": true,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", (len(respBytes))))
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))

	for _, endpoint := range endpoints {
		mux.HandleFunc(endpoint.path, VerifySession(&sessmodels.VerifySessionOptions{
			OverrideGlobalClaimValidators: endpoint.overrideGlobalClaimValidators,
		}, func(w http.ResponseWriter, r *http.Request) {
			sessionContainer := GetSessionFromRequestContext(r.Context())
			resp := map[string]interface{}{
				"message": sessionContainer.GetHandle(),
			}
			respBytes, err := json.Marshal(resp)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", (len(respBytes))))
			w.WriteHeader(http.StatusOK)
			w.Write(respBytes)
		}))
	}

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	return testServer
}

func getTestServerWithoutMiddleware() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/verify", VerifySession(nil, func(w http.ResponseWriter, r *http.Request) {
		sessionContainer := GetSessionFromRequestContext(r.Context())
		resp := map[string]interface{}{
			"message": sessionContainer.GetHandle(),
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", (len(respBytes))))
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}))

	testServer := httptest.NewServer(mux)
	return testServer
}
