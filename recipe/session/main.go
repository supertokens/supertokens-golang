/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package session

import (
	"context"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *sessmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func CreateNewSessionWithContext(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	config := instance.Config
	appInfo := instance.RecipeModule.GetAppInfo()

	return CreateNewSessionInRequest(req, res, config, appInfo, *instance, instance.RecipeImpl, userID, accessTokenPayload, sessionDataInDatabase, userContext)
}

func CreateNewSessionWithContextWithoutRequestResponse(userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, disableAntiCSRF *bool, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}

	claimsAddedByOtherRecipes := instance.GetClaimsAddedByOtherRecipes()
	finalAccessTokenPayload := accessTokenPayload
	if finalAccessTokenPayload == nil {
		finalAccessTokenPayload = map[string]interface{}{}
	}

	appInfo := instance.RecipeModule.GetAppInfo()
	issuer := appInfo.APIDomain.GetAsStringDangerous() + appInfo.APIBasePath.GetAsStringDangerous()

	finalAccessTokenPayload["iss"] = issuer

	for _, claim := range claimsAddedByOtherRecipes {
		finalAccessTokenPayload, err = claim.Build(userID, finalAccessTokenPayload, userContext)
		if err != nil {
			return nil, err
		}
	}

	_disableAntiCSRF := false

	if disableAntiCSRF != nil {
		_disableAntiCSRF = *disableAntiCSRF
	}

	return (*instance.RecipeImpl.CreateNewSession)(userID, accessTokenPayload, sessionDataInDatabase, &_disableAntiCSRF, userContext)
}

func GetSessionWithContext(req http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	config := instance.Config

	return GetSessionFromRequest(req, res, config, options, instance.RecipeImpl, userContext)
}

func GetSessionWithContextWithoutRequestResponse(accessToken string, antiCSRFToken *string, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}

	result, err := (*instance.RecipeImpl.GetSession)(accessToken, antiCSRFToken, options, userContext)

	if err != nil {
		return nil, err
	}

	if result != nil {
		var overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) = nil
		if options != nil {
			overrideGlobalClaimValidators = options.OverrideGlobalClaimValidators
		}

		if err != nil {
			return nil, err
		}
		claimValidators, err := GetRequiredClaimValidators(result, overrideGlobalClaimValidators, userContext)

		if err != nil {
			return nil, err
		}

		err = (*result).AssertClaimsWithContext(claimValidators, userContext)

		if err != nil {
			return nil, err
		}

	}

	return result, nil
}

func GetSessionInformationWithContext(sessionHandle string, userContext supertokens.UserContext) (*sessmodels.SessionInformation, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetSessionInformation)(sessionHandle, userContext)
}

func RefreshSessionWithContext(req *http.Request, res http.ResponseWriter, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return RefreshSessionInRequest(req, res, instance.Config, instance.RecipeImpl, userContext)
}

func RefreshSessionWithContextWithoutRequestResponse(refreshToken string, disableAntiCSRF *bool, antiCSRFToken *string, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}

	_disableAntiCSRF := false

	if disableAntiCSRF != nil {
		_disableAntiCSRF = *disableAntiCSRF
	}

	return (*instance.RecipeImpl.RefreshSession)(refreshToken, antiCSRFToken, _disableAntiCSRF, userContext)
}

func RevokeAllSessionsForUserWithContext(userID string, userContext supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.RevokeAllSessionsForUser)(userID, userContext)
}

func GetAllSessionHandlesForUserWithContext(userID string, userContext supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetAllSessionHandlesForUser)(userID, userContext)
}

func RevokeSessionWithContext(sessionHandle string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return (*instance.RecipeImpl.RevokeSession)(sessionHandle, userContext)
}

func RevokeMultipleSessionsWithContext(sessionHandles []string, userContext supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.RevokeMultipleSessions)(sessionHandles, userContext)
}

func UpdateSessionDataInDatabaseWithContext(sessionHandle string, newSessionData map[string]interface{}, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return (*instance.RecipeImpl.UpdateSessionDataInDatabase)(sessionHandle, newSessionData, userContext)
}

func CreateJWTWithContext(payload map[string]interface{}, validitySecondsPointer *uint64, useStaticSigningKey *bool, userContext supertokens.UserContext) (jwtmodels.CreateJWTResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return jwtmodels.CreateJWTResponse{}, err
	}

	return (*instance.OpenIdRecipe.RecipeImpl.CreateJWT)(payload, validitySecondsPointer, useStaticSigningKey, userContext)
}

func GetJWKSWithContext(userContext supertokens.UserContext) (jwtmodels.GetJWKSResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return jwtmodels.GetJWKSResponse{}, err
	}

	return (*instance.OpenIdRecipe.RecipeImpl.GetJWKS)(userContext)
}

func GetOpenIdDiscoveryConfigurationWithContext(userContext supertokens.UserContext) (openidmodels.GetOpenIdDiscoveryConfigurationResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return openidmodels.GetOpenIdDiscoveryConfigurationResponse{}, err
	}

	return (*instance.OpenIdRecipe.RecipeImpl.GetOpenIdDiscoveryConfiguration)(userContext)
}

func ValidateClaimsForSessionHandleWithContext(
	sessionHandle string,
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionInfo sessmodels.SessionInformation, userContext supertokens.UserContext) []claims.SessionClaimValidator,
	userContext supertokens.UserContext,
) (sessmodels.ValidateClaimsResponse, error) {

	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}

	sessionInfo, err := (*instance.RecipeImpl.GetSessionInformation)(sessionHandle, userContext)
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}

	if sessionInfo == nil {
		return sessmodels.ValidateClaimsResponse{
			SessionDoesNotExistError: &struct{}{},
		}, nil
	}

	claimValidatorsAddedByOtherRecipes := instance.getClaimValidatorsAddedByOtherRecipes()
	claimValidators, err := (*instance.RecipeImpl.GetGlobalClaimValidators)(sessionInfo.UserId, claimValidatorsAddedByOtherRecipes, userContext)
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}

	if overrideGlobalClaimValidators != nil {
		claimValidators = overrideGlobalClaimValidators(claimValidators, *sessionInfo, userContext)
	}

	claimValidationResponse, err := (*instance.RecipeImpl.ValidateClaims)(sessionInfo.UserId, sessionInfo.CustomClaimsInAccessTokenPayload, claimValidators, userContext)
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}
	if claimValidationResponse.AccessTokenPayloadUpdate != nil {
		ok, err := (*instance.RecipeImpl.MergeIntoAccessTokenPayload)(sessionHandle, claimValidationResponse.AccessTokenPayloadUpdate, userContext)
		if err != nil {
			return sessmodels.ValidateClaimsResponse{}, err
		}

		if !ok {
			return sessmodels.ValidateClaimsResponse{
				SessionDoesNotExistError: &struct{}{},
			}, nil
		}
	}
	return sessmodels.ValidateClaimsResponse{
		OK: &struct {
			InvalidClaims []claims.ClaimValidationError
		}{
			InvalidClaims: claimValidationResponse.InvalidClaims,
		},
	}, nil
}

func ValidateClaimsInJWTPayloadWithContext(
	userID string,
	jwtPayload map[string]interface{},
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, userID string, userContext supertokens.UserContext) []claims.SessionClaimValidator,
	userContext supertokens.UserContext,
) ([]claims.ClaimValidationError, error) {

	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}

	claimValidatorsAddedByOtherRecipes := instance.getClaimValidatorsAddedByOtherRecipes()
	claimValidators, err := (*instance.RecipeImpl.GetGlobalClaimValidators)(userID, claimValidatorsAddedByOtherRecipes, userContext)
	if err != nil {
		return nil, err
	}

	if overrideGlobalClaimValidators != nil {
		claimValidators = overrideGlobalClaimValidators(claimValidators, userID, userContext)
	}

	invalidClaims, err := (*instance.RecipeImpl.ValidateClaimsInJWTPayload)(userID, jwtPayload, claimValidators, userContext)
	if err != nil {
		return nil, err
	}

	return invalidClaims, nil
}

func MergeIntoAccessTokenPayloadWithContext(sessionHandle string, accessTokenPayloadUpdate map[string]interface{}, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}

	return (*instance.RecipeImpl.MergeIntoAccessTokenPayload)(sessionHandle, accessTokenPayloadUpdate, userContext)
}

func FetchAndSetClaimWithContext(sessionHandle string, claim *claims.TypeSessionClaim, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}

	return (*instance.RecipeImpl.FetchAndSetClaim)(sessionHandle, claim, userContext)
}

func SetClaimValueWithContext(sessionHandle string, claim *claims.TypeSessionClaim, value interface{}, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}

	return (*instance.RecipeImpl.SetClaimValue)(sessionHandle, claim, value, userContext)
}

func GetClaimValueWithContext(sessionHandle string, claim *claims.TypeSessionClaim, userContext supertokens.UserContext) (sessmodels.GetClaimValueResult, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.GetClaimValueResult{}, err
	}

	return (*instance.RecipeImpl.GetClaimValue)(sessionHandle, claim, userContext)
}

func RemoveClaimWithContext(sessionHandle string, claim *claims.TypeSessionClaim, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}

	return (*instance.RecipeImpl.RemoveClaim)(sessionHandle, claim, userContext)
}

func VerifySession(options *sessmodels.VerifySessionOptions, otherHandler http.HandlerFunc) http.HandlerFunc {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		panic("can't fetch supertokens instance. You should call the supertokens.Init function before using the VerifySession function.")
	}
	return VerifySessionHelper(*instance, options, otherHandler)
}

func GetSessionFromRequestContext(ctx context.Context) sessmodels.SessionContainer {
	value := ctx.Value(sessmodels.SessionContext)
	if value == nil {
		return nil
	}
	temp := value.(sessmodels.SessionContainer)
	return temp
}

func CreateNewSession(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}) (sessmodels.SessionContainer, error) {
	return CreateNewSessionWithContext(req, res, userID, accessTokenPayload, sessionDataInDatabase, &map[string]interface{}{})
}

func CreateNewSessionWithoutRequestResponse(userId string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, disableAntiCSRF *bool) (sessmodels.SessionContainer, error) {
	return CreateNewSessionWithContextWithoutRequestResponse(userId, accessTokenPayload, sessionDataInDatabase, disableAntiCSRF, nil)
}

func GetSession(req http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions) (sessmodels.SessionContainer, error) {
	return GetSessionWithContext(req, res, options, &map[string]interface{}{})
}

func GetSessionWithoutRequestResponse(accessToken string, antiCSRFToken *string, options *sessmodels.VerifySessionOptions) (sessmodels.SessionContainer, error) {
	return GetSessionWithContextWithoutRequestResponse(accessToken, antiCSRFToken, options, nil)
}

func GetSessionInformation(sessionHandle string) (*sessmodels.SessionInformation, error) {
	return GetSessionInformationWithContext(sessionHandle, &map[string]interface{}{})
}

func RefreshSession(req *http.Request, res http.ResponseWriter) (sessmodels.SessionContainer, error) {
	return RefreshSessionWithContext(req, res, &map[string]interface{}{})
}

func RefreshSessionWithoutRequestResponse(refreshToken string, disableAntiCSRF *bool, antiCSRFToken *string) (sessmodels.SessionContainer, error) {
	return RefreshSessionWithContextWithoutRequestResponse(refreshToken, disableAntiCSRF, antiCSRFToken, nil)
}

func RevokeAllSessionsForUser(userID string) ([]string, error) {
	return RevokeAllSessionsForUserWithContext(userID, &map[string]interface{}{})
}

func GetAllSessionHandlesForUser(userID string) ([]string, error) {
	return GetAllSessionHandlesForUserWithContext(userID, &map[string]interface{}{})
}

func RevokeSession(sessionHandle string) (bool, error) {
	return RevokeSessionWithContext(sessionHandle, &map[string]interface{}{})
}

func RevokeMultipleSessions(sessionHandles []string) ([]string, error) {
	return RevokeMultipleSessionsWithContext(sessionHandles, &map[string]interface{}{})
}

func UpdateSessionDataInDatabase(sessionHandle string, newSessionData map[string]interface{}) (bool, error) {
	return UpdateSessionDataInDatabaseWithContext(sessionHandle, newSessionData, &map[string]interface{}{})
}

func CreateJWT(payload map[string]interface{}, useStaticSigningKey *bool, validitySecondsPointer *uint64) (jwtmodels.CreateJWTResponse, error) {
	return CreateJWTWithContext(payload, validitySecondsPointer, useStaticSigningKey, &map[string]interface{}{})
}

func GetJWKS() (jwtmodels.GetJWKSResponse, error) {
	return GetJWKSWithContext(&map[string]interface{}{})
}

func GetOpenIdDiscoveryConfiguration() (openidmodels.GetOpenIdDiscoveryConfigurationResponse, error) {
	return GetOpenIdDiscoveryConfigurationWithContext(&map[string]interface{}{})
}

func ValidateClaimsForSessionHandle(
	sessionHandle string,
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionInfo sessmodels.SessionInformation, userContext supertokens.UserContext) []claims.SessionClaimValidator,
) (sessmodels.ValidateClaimsResponse, error) {
	return ValidateClaimsForSessionHandleWithContext(sessionHandle, overrideGlobalClaimValidators, &map[string]interface{}{})
}

func ValidateClaimsInJWTPayload(
	userID string,
	jwtPayload map[string]interface{},
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, userID string, userContext supertokens.UserContext) []claims.SessionClaimValidator,
) ([]claims.ClaimValidationError, error) {
	return ValidateClaimsInJWTPayloadWithContext(userID, jwtPayload, overrideGlobalClaimValidators, &map[string]interface{}{})
}

func MergeIntoAccessTokenPayload(sessionHandle string, accessTokenPayloadUpdate map[string]interface{}) (bool, error) {
	return MergeIntoAccessTokenPayloadWithContext(sessionHandle, accessTokenPayloadUpdate, &map[string]interface{}{})
}

func FetchAndSetClaim(sessionHandle string, claim *claims.TypeSessionClaim) (bool, error) {
	return FetchAndSetClaimWithContext(sessionHandle, claim, &map[string]interface{}{})
}

func SetClaimValue(sessionHandle string, claim *claims.TypeSessionClaim, value interface{}) (bool, error) {
	return SetClaimValueWithContext(sessionHandle, claim, value, &map[string]interface{}{})
}

func GetClaimValue(sessionHandle string, claim *claims.TypeSessionClaim) (sessmodels.GetClaimValueResult, error) {
	return GetClaimValueWithContext(sessionHandle, claim, &map[string]interface{}{})
}

func RemoveClaim(sessionHandle string, claim *claims.TypeSessionClaim) (bool, error) {
	return RemoveClaimWithContext(sessionHandle, claim, &map[string]interface{}{})
}
