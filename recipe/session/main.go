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

func CreateNewSession(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, userContext ...supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	config := instance.Config
	appInfo := instance.RecipeModule.GetAppInfo()

	return CreateNewSessionInRequest(req, res, config, appInfo, *instance, instance.RecipeImpl, userID, accessTokenPayload, sessionDataInDatabase, userContext[0])
}

func CreateNewSessionWithoutRequestResponse(tenantId string, userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, disableAntiCSRF *bool, userContext ...supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
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
		finalAccessTokenPayload, err = claim.Build(userID, tenantId, finalAccessTokenPayload, userContext[0])
		if err != nil {
			return nil, err
		}
	}

	_disableAntiCSRF := false

	if disableAntiCSRF != nil {
		_disableAntiCSRF = *disableAntiCSRF
	}

	return (*instance.RecipeImpl.CreateNewSession)(userID, accessTokenPayload, sessionDataInDatabase, &_disableAntiCSRF, userContext[0])
}

func GetSession(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions, userContext ...supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	config := instance.Config

	return GetSessionFromRequest(req, res, config, options, instance.RecipeImpl, userContext[0])
}

func GetSessionWithoutRequestResponse(accessToken string, antiCSRFToken *string, options *sessmodels.VerifySessionOptions, userContext ...supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}

	result, err := (*instance.RecipeImpl.GetSession)(&accessToken, antiCSRFToken, options, userContext[0])

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
		claimValidators, err := GetRequiredClaimValidators(result, overrideGlobalClaimValidators, userContext[0])

		if err != nil {
			return nil, err
		}

		err = (*result).AssertClaimsWithContext(claimValidators, userContext[0])

		if err != nil {
			return nil, err
		}

	}

	return result, nil
}

func GetSessionInformation(sessionHandle string, userContext ...supertokens.UserContext) (*sessmodels.SessionInformation, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetSessionInformation)(sessionHandle, userContext[0])
}

func RefreshSession(req *http.Request, res http.ResponseWriter, userContext ...supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return RefreshSessionInRequest(req, res, instance.Config, instance.RecipeImpl, userContext[0])
}

func RefreshSessionWithoutRequestResponse(refreshToken string, disableAntiCSRF *bool, antiCSRFToken *string, userContext ...supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}

	_disableAntiCSRF := false

	if disableAntiCSRF != nil {
		_disableAntiCSRF = *disableAntiCSRF
	}

	return (*instance.RecipeImpl.RefreshSession)(refreshToken, antiCSRFToken, _disableAntiCSRF, userContext[0])
}

func RevokeAllSessionsForUser(userID string, userContext ...supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.RevokeAllSessionsForUser)(userID, userContext[0])
}

func GetAllSessionHandlesForUser(userID string, userContext ...supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetAllSessionHandlesForUser)(userID, userContext[0])
}

func RevokeSession(sessionHandle string, userContext ...supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.RevokeSession)(sessionHandle, userContext[0])
}

func RevokeMultipleSessions(sessionHandles []string, userContext ...supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.RevokeMultipleSessions)(sessionHandles, userContext[0])
}

func UpdateSessionDataInDatabase(sessionHandle string, newSessionData map[string]interface{}, userContext ...supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.UpdateSessionDataInDatabase)(sessionHandle, newSessionData, userContext[0])
}

func CreateJWT(payload map[string]interface{}, validitySecondsPointer *uint64, useStaticSigningKey *bool, userContext ...supertokens.UserContext) (jwtmodels.CreateJWTResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return jwtmodels.CreateJWTResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.OpenIdRecipe.RecipeImpl.CreateJWT)(payload, validitySecondsPointer, useStaticSigningKey, userContext[0])
}

func GetJWKS(userContext ...supertokens.UserContext) (jwtmodels.GetJWKSResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return jwtmodels.GetJWKSResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.OpenIdRecipe.RecipeImpl.GetJWKS)(userContext[0])
}

func GetOpenIdDiscoveryConfiguration(userContext ...supertokens.UserContext) (openidmodels.GetOpenIdDiscoveryConfigurationResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return openidmodels.GetOpenIdDiscoveryConfigurationResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.OpenIdRecipe.RecipeImpl.GetOpenIdDiscoveryConfiguration)(userContext[0])
}

func ValidateClaimsForSessionHandle(
	sessionHandle string,
	overrideGlobalClaimValidators func([]claims.SessionClaimValidator, sessmodels.SessionInformation, supertokens.UserContext) []claims.SessionClaimValidator,
	userContext ...supertokens.UserContext,
) (sessmodels.ValidateClaimsResponse, error) {

	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	sessionInfo, err := (*instance.RecipeImpl.GetSessionInformation)(sessionHandle, userContext[0])
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}

	if sessionInfo == nil {
		return sessmodels.ValidateClaimsResponse{
			SessionDoesNotExistError: &struct{}{},
		}, nil
	}

	claimValidatorsAddedByOtherRecipes := instance.getClaimValidatorsAddedByOtherRecipes()
	claimValidators, err := (*instance.RecipeImpl.GetGlobalClaimValidators)(sessionInfo.UserId, claimValidatorsAddedByOtherRecipes, userContext[0])
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}

	if overrideGlobalClaimValidators != nil {
		claimValidators = overrideGlobalClaimValidators(claimValidators, *sessionInfo, userContext[0])
	}

	claimValidationResponse, err := (*instance.RecipeImpl.ValidateClaims)(sessionInfo.UserId, sessionInfo.CustomClaimsInAccessTokenPayload, claimValidators, userContext[0])
	if err != nil {
		return sessmodels.ValidateClaimsResponse{}, err
	}
	if claimValidationResponse.AccessTokenPayloadUpdate != nil {
		ok, err := (*instance.RecipeImpl.MergeIntoAccessTokenPayload)(sessionHandle, claimValidationResponse.AccessTokenPayloadUpdate, userContext[0])
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

func ValidateClaimsInJWTPayload(
	userID string,
	jwtPayload map[string]interface{},
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, userID string, userContext ...supertokens.UserContext) []claims.SessionClaimValidator,
	userContext ...supertokens.UserContext,
) ([]claims.ClaimValidationError, error) {

	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	claimValidatorsAddedByOtherRecipes := instance.getClaimValidatorsAddedByOtherRecipes()
	claimValidators, err := (*instance.RecipeImpl.GetGlobalClaimValidators)(userID, claimValidatorsAddedByOtherRecipes, userContext[0])
	if err != nil {
		return nil, err
	}

	if overrideGlobalClaimValidators != nil {
		claimValidators = overrideGlobalClaimValidators(claimValidators, userID, userContext[0])
	}

	invalidClaims, err := (*instance.RecipeImpl.ValidateClaimsInJWTPayload)(userID, jwtPayload, claimValidators, userContext[0])
	if err != nil {
		return nil, err
	}

	return invalidClaims, nil
}

func MergeIntoAccessTokenPayload(sessionHandle string, accessTokenPayloadUpdate map[string]interface{}, userContext ...supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.MergeIntoAccessTokenPayload)(sessionHandle, accessTokenPayloadUpdate, userContext[0])
}

func FetchAndSetClaim(sessionHandle string, claim *claims.TypeSessionClaim, userContext ...supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.FetchAndSetClaim)(sessionHandle, claim, userContext[0])
}

func SetClaimValue(sessionHandle string, claim *claims.TypeSessionClaim, value interface{}, userContext ...supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.SetClaimValue)(sessionHandle, claim, value, userContext[0])
}

func GetClaimValue(sessionHandle string, claim *claims.TypeSessionClaim, userContext ...supertokens.UserContext) (sessmodels.GetClaimValueResult, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.GetClaimValueResult{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetClaimValue)(sessionHandle, claim, userContext[0])
}

func RemoveClaim(sessionHandle string, claim *claims.TypeSessionClaim, userContext ...supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.RemoveClaim)(sessionHandle, claim, userContext[0])
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
