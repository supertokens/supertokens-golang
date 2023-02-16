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
	"errors"
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

func CreateNewSessionWithContext(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, tenantId *string, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}

	claimsAddedByOtherRecipes := instance.getClaimsAddedByOtherRecipes()
	finalAccessTokenPayload := accessTokenPayload
	if finalAccessTokenPayload == nil {
		finalAccessTokenPayload = map[string]interface{}{}
	}

	for _, claim := range claimsAddedByOtherRecipes {
		finalAccessTokenPayload, err = claim.Build(userID, finalAccessTokenPayload, tenantId, userContext)
		if err != nil {
			return nil, err
		}
	}

	return (*instance.RecipeImpl.CreateNewSession)(req, res, userID, finalAccessTokenPayload, sessionData, tenantId, userContext)
}

func GetSessionWithContext(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	sessionContainer, err := (*instance.RecipeImpl.GetSession)(req, res, options, userContext)
	if err != nil {
		return nil, err
	}

	if sessionContainer != nil {
		var overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) = nil
		if options != nil {
			overrideGlobalClaimValidators = options.OverrideGlobalClaimValidators
		}
		claimValidators, err := getRequiredClaimValidators(sessionContainer, overrideGlobalClaimValidators, userContext)
		if err != nil {
			return nil, err
		}
		err = sessionContainer.AssertClaimsWithContext(claimValidators, userContext)
		if err != nil {
			return nil, err
		}
	}
	return sessionContainer, nil
}

func GetSessionInformationWithContext(sessionHandle string, tenantId *string, userContext supertokens.UserContext) (*sessmodels.SessionInformation, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetSessionInformation)(sessionHandle, tenantId, userContext)
}

func RefreshSessionWithContext(req *http.Request, res http.ResponseWriter, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.RefreshSession)(req, res, userContext)
}

func RevokeAllSessionsForUserWithContext(userID string, tenantId *string, userContext supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.RevokeAllSessionsForUser)(userID, tenantId, userContext)
}

func GetAllSessionHandlesForUserWithContext(userID string, tenantId *string, userContext supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetAllSessionHandlesForUser)(userID, tenantId, userContext)
}

func RevokeSessionWithContext(sessionHandle string, tenantId *string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return (*instance.RecipeImpl.RevokeSession)(sessionHandle, tenantId, userContext)
}

func RevokeMultipleSessionsWithContext(sessionHandles []string, tenantId *string, userContext supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.RevokeMultipleSessions)(sessionHandles, tenantId, userContext)
}

func UpdateSessionDataWithContext(sessionHandle string, newSessionData map[string]interface{}, tenantId *string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return (*instance.RecipeImpl.UpdateSessionData)(sessionHandle, newSessionData, tenantId, userContext)
}

// Deprecated: use MergeIntoAccessTokenPayloadWithContext instead
func UpdateAccessTokenPayloadWithContext(sessionHandle string, newAccessTokenPayload map[string]interface{}, tenantId *string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return (*instance.RecipeImpl.UpdateAccessTokenPayload)(sessionHandle, newAccessTokenPayload, tenantId, userContext)
}

func CreateJWTWithContext(payload map[string]interface{}, validitySecondsPointer *uint64, userContext supertokens.UserContext) (jwtmodels.CreateJWTResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return jwtmodels.CreateJWTResponse{}, err
	}
	if instance.OpenIdRecipe == nil {
		return jwtmodels.CreateJWTResponse{}, errors.New("CreateJWT cannot be used without enabling the Jwt feature")
	}
	return (*instance.OpenIdRecipe.RecipeImpl.CreateJWT)(payload, validitySecondsPointer, userContext)
}

func GetJWKSWithContext(userContext supertokens.UserContext) (jwtmodels.GetJWKSResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return jwtmodels.GetJWKSResponse{}, err
	}
	if instance.OpenIdRecipe == nil {
		return jwtmodels.GetJWKSResponse{}, errors.New("GetJWKS cannot be used without enabling the Jwt feature")
	}
	return (*instance.OpenIdRecipe.RecipeImpl.GetJWKS)(userContext)
}

func GetOpenIdDiscoveryConfigurationWithContext(userContext supertokens.UserContext) (openidmodels.GetOpenIdDiscoveryConfigurationResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return openidmodels.GetOpenIdDiscoveryConfigurationResponse{}, err
	}
	if instance.OpenIdRecipe == nil {
		return openidmodels.GetOpenIdDiscoveryConfigurationResponse{}, errors.New("GetOpenIdDiscoveryConfiguration cannot be used without enabling the Jwt feature")
	}
	return (*instance.OpenIdRecipe.RecipeImpl.GetOpenIdDiscoveryConfiguration)(userContext)
}

func RegenerateAccessTokenWithContext(accessToken string, newAccessTokenPayload *map[string]interface{}, sessionHandle string, tenantId *string, userContext supertokens.UserContext) (*sessmodels.RegenerateAccessTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.RegenerateAccessToken)(accessToken, newAccessTokenPayload, tenantId, userContext)
}

func ValidateClaimsForSessionHandleWithContext(
	sessionHandle string,
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionInfo sessmodels.SessionInformation, userContext supertokens.UserContext) []claims.SessionClaimValidator, tenantId *string,
	userContext supertokens.UserContext,
) (sessmodels.ValidateClaimsResponse, error) {

	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}

	sessionInfo, err := (*instance.RecipeImpl.GetSessionInformation)(sessionHandle, tenantId, userContext)
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}

	if sessionInfo == nil {
		return sessmodels.ValidateClaimsResponse{
			SessionDoesNotExistError: &struct{}{},
		}, nil
	}

	claimValidatorsAddedByOtherRecipes := instance.getClaimValidatorsAddedByOtherRecipes()
	claimValidators, err := (*instance.RecipeImpl.GetGlobalClaimValidators)(sessionInfo.UserId, claimValidatorsAddedByOtherRecipes, tenantId, userContext)
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}

	if overrideGlobalClaimValidators != nil {
		claimValidators = overrideGlobalClaimValidators(claimValidators, *sessionInfo, userContext)
	}

	claimValidationResponse, err := (*instance.RecipeImpl.ValidateClaims)(sessionInfo.UserId, sessionInfo.AccessTokenPayload, claimValidators, tenantId, userContext)
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}
	if claimValidationResponse.AccessTokenPayloadUpdate != nil {
		ok, err := (*instance.RecipeImpl.MergeIntoAccessTokenPayload)(sessionHandle, claimValidationResponse.AccessTokenPayloadUpdate, tenantId, userContext)
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
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, userID string, userContext supertokens.UserContext) []claims.SessionClaimValidator, tenantId *string,
	userContext supertokens.UserContext,
) ([]claims.ClaimValidationError, error) {

	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}

	claimValidatorsAddedByOtherRecipes := instance.getClaimValidatorsAddedByOtherRecipes()
	claimValidators, err := (*instance.RecipeImpl.GetGlobalClaimValidators)(userID, claimValidatorsAddedByOtherRecipes, tenantId, userContext)
	if err != nil {
		return nil, err
	}

	if overrideGlobalClaimValidators != nil {
		claimValidators = overrideGlobalClaimValidators(claimValidators, userID, userContext)
	}

	invalidClaims, err := (*instance.RecipeImpl.ValidateClaimsInJWTPayload)(userID, jwtPayload, claimValidators, tenantId, userContext)
	if err != nil {
		return nil, err
	}

	return invalidClaims, nil
}

func MergeIntoAccessTokenPayloadWithContext(sessionHandle string, accessTokenPayloadUpdate map[string]interface{}, tenantId *string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}

	return (*instance.RecipeImpl.MergeIntoAccessTokenPayload)(sessionHandle, accessTokenPayloadUpdate, tenantId, userContext)
}

func FetchAndSetClaimWithContext(sessionHandle string, claim *claims.TypeSessionClaim, tenantId *string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}

	return (*instance.RecipeImpl.FetchAndSetClaim)(sessionHandle, claim, tenantId, userContext)
}

func SetClaimValueWithContext(sessionHandle string, claim *claims.TypeSessionClaim, value interface{}, tenantId *string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}

	return (*instance.RecipeImpl.SetClaimValue)(sessionHandle, claim, value, tenantId, userContext)
}

func GetClaimValueWithContext(sessionHandle string, claim *claims.TypeSessionClaim, tenantId *string, userContext supertokens.UserContext) (sessmodels.GetClaimValueResult, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.GetClaimValueResult{}, err
	}

	return (*instance.RecipeImpl.GetClaimValue)(sessionHandle, claim, tenantId, userContext)
}

func RemoveClaimWithContext(sessionHandle string, claim *claims.TypeSessionClaim, tenantId *string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}

	return (*instance.RecipeImpl.RemoveClaim)(sessionHandle, claim, tenantId, userContext)
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

func CreateNewSession(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, tenantId *string) (sessmodels.SessionContainer, error) {
	return CreateNewSessionWithContext(req, res, userID, accessTokenPayload, sessionData, tenantId, &map[string]interface{}{})
}

func GetSession(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions) (sessmodels.SessionContainer, error) {
	return GetSessionWithContext(req, res, options, &map[string]interface{}{})
}

func GetSessionInformation(sessionHandle string, tenantId *string) (*sessmodels.SessionInformation, error) {
	return GetSessionInformationWithContext(sessionHandle, tenantId, &map[string]interface{}{})
}

func RefreshSession(req *http.Request, res http.ResponseWriter) (sessmodels.SessionContainer, error) {
	return RefreshSessionWithContext(req, res, &map[string]interface{}{})
}

func RevokeAllSessionsForUser(userID string, tenantId *string) ([]string, error) {
	return RevokeAllSessionsForUserWithContext(userID, tenantId, &map[string]interface{}{})
}

func GetAllSessionHandlesForUser(userID string, tenantId *string) ([]string, error) {
	return GetAllSessionHandlesForUserWithContext(userID, tenantId, &map[string]interface{}{})
}

func RevokeSession(sessionHandle string, tenantId *string) (bool, error) {
	return RevokeSessionWithContext(sessionHandle, tenantId, &map[string]interface{}{})
}

func RevokeMultipleSessions(sessionHandles []string, tenantId *string) ([]string, error) {
	return RevokeMultipleSessionsWithContext(sessionHandles, tenantId, &map[string]interface{}{})
}

func UpdateSessionData(sessionHandle string, newSessionData map[string]interface{}, tenantId *string) (bool, error) {
	return UpdateSessionDataWithContext(sessionHandle, newSessionData, tenantId, &map[string]interface{}{})
}

// Deprecated: use MergeIntoAccessTokenPayload instead
func UpdateAccessTokenPayload(sessionHandle string, newAccessTokenPayload map[string]interface{}, tenantId *string) (bool, error) {
	return UpdateAccessTokenPayloadWithContext(sessionHandle, newAccessTokenPayload, tenantId, &map[string]interface{}{})
}

func CreateJWT(payload map[string]interface{}, validitySecondsPointer *uint64) (jwtmodels.CreateJWTResponse, error) {
	return CreateJWTWithContext(payload, validitySecondsPointer, &map[string]interface{}{})
}

func GetJWKS() (jwtmodels.GetJWKSResponse, error) {
	return GetJWKSWithContext(&map[string]interface{}{})
}

func GetOpenIdDiscoveryConfiguration() (openidmodels.GetOpenIdDiscoveryConfigurationResponse, error) {
	return GetOpenIdDiscoveryConfigurationWithContext(&map[string]interface{}{})
}

func RegenerateAccessToken(accessToken string, newAccessTokenPayload *map[string]interface{}, sessionHandle string, tenantId *string) (*sessmodels.RegenerateAccessTokenResponse, error) {
	return RegenerateAccessTokenWithContext(accessToken, newAccessTokenPayload, sessionHandle, tenantId, &map[string]interface{}{})
}

func ValidateClaimsForSessionHandle(
	sessionHandle string,
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionInfo sessmodels.SessionInformation, userContext supertokens.UserContext) []claims.SessionClaimValidator,
	tenantId *string,
) (sessmodels.ValidateClaimsResponse, error) {
	return ValidateClaimsForSessionHandleWithContext(sessionHandle, overrideGlobalClaimValidators, tenantId, &map[string]interface{}{})
}

func ValidateClaimsInJWTPayload(
	userID string,
	jwtPayload map[string]interface{},
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, userID string, userContext supertokens.UserContext) []claims.SessionClaimValidator,
	tenantId *string,
) ([]claims.ClaimValidationError, error) {
	return ValidateClaimsInJWTPayloadWithContext(userID, jwtPayload, overrideGlobalClaimValidators, tenantId, &map[string]interface{}{})
}

func MergeIntoAccessTokenPayload(sessionHandle string, accessTokenPayloadUpdate map[string]interface{}, tenantId *string) (bool, error) {
	return MergeIntoAccessTokenPayloadWithContext(sessionHandle, accessTokenPayloadUpdate, tenantId, &map[string]interface{}{})
}

func FetchAndSetClaim(sessionHandle string, claim *claims.TypeSessionClaim, tenantId *string) (bool, error) {
	return FetchAndSetClaimWithContext(sessionHandle, claim, tenantId, &map[string]interface{}{})
}

func SetClaimValue(sessionHandle string, claim *claims.TypeSessionClaim, value interface{}, tenantId *string) (bool, error) {
	return SetClaimValueWithContext(sessionHandle, claim, value, tenantId, &map[string]interface{}{})
}

func GetClaimValue(sessionHandle string, claim *claims.TypeSessionClaim, tenantId *string) (sessmodels.GetClaimValueResult, error) {
	return GetClaimValueWithContext(sessionHandle, claim, tenantId, &map[string]interface{}{})
}

func RemoveClaim(sessionHandle string, claim *claims.TypeSessionClaim, tenantId *string) (bool, error) {
	return RemoveClaimWithContext(sessionHandle, claim, tenantId, &map[string]interface{}{})
}
