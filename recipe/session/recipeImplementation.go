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
	"bytes"
	"encoding/json"
	defaultErrors "errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var handshakeInfoLock sync.Mutex

var protectedProps = []string{
	"sub",
	"iat",
	"exp",
	"sessionHandle",
	"parentRefreshTokenHash1",
	"refreshTokenHash1",
	"antiCsrfToken",
}

func MakeRecipeImplementation(querier supertokens.Querier, config sessmodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) sessmodels.RecipeInterface {
	var result sessmodels.RecipeInterface

	createNewSession := func(userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, disableAntiCsrf *bool, userContext supertokens.UserContext) (sessmodels.CreateNewSessionResponse, error) {
		supertokens.LogDebugMessage("createNewSession: Started")

		sessionResponse, err := createNewSessionHelper(
			config, querier, userID, disableAntiCsrf != nil && *disableAntiCsrf == true, accessTokenPayload, sessionDataInDatabase,
		)
		if err != nil {
			return sessmodels.CreateNewSessionResponse{}, err
		}

		supertokens.LogDebugMessage("createNewSession: Finished")

		parsedJWT, parseErr := ParseJWTWithoutSignatureVerification(sessionResponse.AccessToken.Token)
		if parseErr != nil {
			return sessmodels.CreateNewSessionResponse{}, parseErr
		}

		frontToken := BuildFrontToken(sessionResponse.Session.UserID, sessionResponse.Session.ExpiryTime, sessionResponse.Session.UserDataInAccessToken)
		session := sessionResponse.Session
		sessionContainerInput := makeSessionContainerInput(sessionResponse.AccessToken.Token, session.Handle, session.UserID, parsedJWT.Payload, result, frontToken, sessionResponse.AntiCsrfToken, nil, &sessionResponse.RefreshToken, true)
		return sessmodels.CreateNewSessionResponse{
			Status:  "OK",
			Session: newSessionContainer(config, &sessionContainerInput),
		}, nil
	}

	// In all cases if sIdRefreshToken token exists (so it's a legacy session) we return TRY_REFRESH_TOKEN. The refresh endpoint will clear this cookie and try to upgrade the session.
	// Check https://supertokens.com/docs/contribute/decisions/session/0007 for further details and a table of expected behaviours
	getSession := func(accessTokenString string, antiCsrfToken *string, options *sessmodels.GetSessionOptions, userContext supertokens.UserContext) (sessmodels.GetSessionFunctionResponse, error) {
		if options != nil && *options.AntiCsrfCheck != false && config.AntiCsrf != AntiCSRF_VIA_CUSTOM_HEADER {
			return sessmodels.GetSessionFunctionResponse{}, defaultErrors.New("Since the anti-csrf mode is VIA_CUSTOM_HEADER getSession can't check the CSRF token. Please either use VIA_TOKEN or set antiCsrfCheck to false")
		}

		supertokens.LogDebugMessage("getSession: Started")
		var accessToken *sessmodels.ParsedJWTInfo

		accessTokenResponse, err := ParseJWTWithoutSignatureVerification(accessTokenString)

		if err != nil {
			supertokens.LogDebugMessage("getSession: Returning UNAUTHORISED because parsing failed")
			return sessmodels.GetSessionFunctionResponse{
				Status:  "UNAUTHORISED",
				Session: nil,
				Error:   &err,
			}, nil
		}

		err = ValidateAccessTokenStructure(accessTokenResponse.Payload, accessTokenResponse.Version)

		if err != nil {
			supertokens.LogDebugMessage("getSession: Returning UNAUTHORISED because parsing failed")
			return sessmodels.GetSessionFunctionResponse{
				Status:  "UNAUTHORISED",
				Session: nil,
				Error:   &err,
			}, nil
		}

		alwaysCheckCore := false

		if options.CheckDatabase != nil {
			alwaysCheckCore = *options.CheckDatabase == true
		}

		doAntiCsrfCheck := options != nil && *options.AntiCsrfCheck != false

		response, err := getSessionHelper(config, querier, *accessToken, antiCsrfToken, doAntiCsrfCheck, alwaysCheckCore)
		if err != nil {
			if defaultErrors.As(err, &errors.TryRefreshTokenError{}) {
				supertokens.LogDebugMessage("getSession: Returning TRY_REFRESH_TOKEN_ERROR because of an exception during getSession")
				return sessmodels.GetSessionFunctionResponse{
					Status:  "TRY_REFRESH_TOKEN_ERROR",
					Session: nil,
					Error:   &err,
				}, nil
			}

			supertokens.LogDebugMessage("getSession: Returning UNAUTHORISED because of an exception during getSession")
			return sessmodels.GetSessionFunctionResponse{
				Status:  "UNAUTHORISED",
				Session: nil,
				Error:   &err,
			}, nil
		}

		supertokens.LogDebugMessage("getSession: Success!")
		var payload map[string]interface{}

		if reflect.DeepEqual(response.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
			parsedToken, parseErr := ParseJWTWithoutSignatureVerification(response.AccessToken.Token)

			if parseErr != nil {
				return sessmodels.GetSessionFunctionResponse{}, parseErr
			}

			payload = parsedToken.Payload
		} else {
			payload = accessToken.Payload
		}

		accessTokenStringForSession := accessTokenString

		accessTokenNil := reflect.DeepEqual(response.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{})

		if !accessTokenNil {
			accessTokenStringForSession = response.AccessToken.Token
		}

		frontToken := BuildFrontToken(response.Session.UserID, response.Session.ExpiryTime, response.Session.UserDataInAccessToken)
		session := response.Session

		sessionContainerInput := makeSessionContainerInput(accessTokenStringForSession, session.Handle, session.UserID, payload, result, frontToken, antiCsrfToken, nil, nil, !accessTokenNil)
		sessionContainer := newSessionContainer(config, &sessionContainerInput)

		return sessmodels.GetSessionFunctionResponse{
			Status:  "OK",
			Session: &sessionContainer,
			Error:   nil,
		}, nil
	}

	getSessionInformation := func(sessionHandle string, userContext supertokens.UserContext) (*sessmodels.SessionInformation, error) {
		return getSessionInformationHelper(querier, sessionHandle)
	}

	refreshSession := func(refreshToken string, antiCsrfToken *string, disableAntiCsrf bool, userContext supertokens.UserContext) (sessmodels.GetSessionFunctionResponse, error) {
		if disableAntiCsrf != true && config.AntiCsrf != AntiCSRF_VIA_CUSTOM_HEADER {
			return sessmodels.GetSessionFunctionResponse{}, defaultErrors.New("Since the anti-csrf mode is VIA_CUSTOM_HEADER getSession can't check the CSRF token. Please either use VIA_TOKEN or set antiCsrfCheck to false")
		}

		supertokens.LogDebugMessage("refreshSession: Started")

		response, err := refreshSessionHelper(config, querier, refreshToken, antiCsrfToken, disableAntiCsrf)
		if err != nil {
			unauthorisedErr := errors.UnauthorizedError{}
			isUnauthorisedErr := defaultErrors.As(err, &unauthorisedErr)
			isTokenTheftDetectedErr := defaultErrors.As(err, &errors.TokenTheftDetectedError{})

			// This token isn't handled by GetToken/setToken to limit the scope of this legacy/migration code
			if (isTokenTheftDetectedErr) || (isUnauthorisedErr && unauthorisedErr.ClearTokens != nil && *unauthorisedErr.ClearTokens) {
				return sessmodels.GetSessionFunctionResponse{
					Status: "TOKEN_THEFT_DETECTED",
					Error:  &err,
				}, nil
			}

			return sessmodels.GetSessionFunctionResponse{
				Status: "UNAUTHORISED",
				Error:  &err,
			}, nil
		}
		supertokens.LogDebugMessage("refreshSession: Success!")

		responseToken, parseErr := ParseJWTWithoutSignatureVerification(response.AccessToken.Token)
		if parseErr != nil {
			return sessmodels.GetSessionFunctionResponse{}, err
		}

		session := response.Session
		frontToken := BuildFrontToken(session.UserID, session.ExpiryTime, responseToken.Payload)

		sessionContainerInput := makeSessionContainerInput(response.AccessToken.Token, session.Handle, session.UserID, responseToken.Payload, result, frontToken, antiCsrfToken, nil, &response.RefreshToken, true)
		sessionContainer := newSessionContainer(config, &sessionContainerInput)

		return sessmodels.GetSessionFunctionResponse{
			Status:  "OK",
			Session: &sessionContainer,
			Error:   nil,
		}, nil
	}

	revokeAllSessionsForUser := func(userID string, userContext supertokens.UserContext) ([]string, error) {
		return revokeAllSessionsForUserHelper(querier, userID)
	}

	getAllSessionHandlesForUser := func(userID string, userContext supertokens.UserContext) ([]string, error) {
		return getAllSessionHandlesForUserHelper(querier, userID)
	}

	revokeSession := func(sessionHandle string, userContext supertokens.UserContext) (bool, error) {
		return revokeSessionHelper(querier, sessionHandle)
	}

	revokeMultipleSessions := func(sessionHandles []string, userContext supertokens.UserContext) ([]string, error) {
		return revokeMultipleSessionsHelper(querier, sessionHandles)
	}

	updateSessionDataInDatabase := func(sessionHandle string, newSessionData map[string]interface{}, userContext supertokens.UserContext) (bool, error) {
		return updateSessionDataInDatabaseHelper(querier, sessionHandle, newSessionData)
	}

	regenerateAccessToken := func(accessToken string, newAccessTokenPayload *map[string]interface{}, userContext supertokens.UserContext) (*sessmodels.RegenerateAccessTokenResponse, error) {
		return regenerateAccessTokenHelper(querier, newAccessTokenPayload, accessToken)
	}

	mergeIntoAccessTokenPayload := func(sessionHandle string, accessTokenPayloadUpdate map[string]interface{}, userContext supertokens.UserContext) (bool, error) {
		sessionInfo, err := (*result.GetSessionInformation)(sessionHandle, userContext)
		if err != nil {
			return false, err
		}
		if sessionInfo == nil {
			return false, nil
		}
		newAccessTokenPayload := map[string]interface{}{}
		for k, v := range sessionInfo.AccessTokenPayload {
			newAccessTokenPayload[k] = v
		}

		for k, _ := range newAccessTokenPayload {
			if supertokens.DoesSliceContainString(k, protectedProps) {
				delete(newAccessTokenPayload, k)
			}
		}

		for k, v := range accessTokenPayloadUpdate {
			if v == nil {
				delete(newAccessTokenPayload, k)
			} else {
				newAccessTokenPayload[k] = v
			}
		}

		return updateAccessTokenPayloadHelper(querier, sessionHandle, newAccessTokenPayload)
	}

	getGlobalClaimValidators := func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
		return claimValidatorsAddedByOtherRecipes, nil
	}

	validateClaims := func(userId string, accessTokenPayload map[string]interface{}, claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) (sessmodels.ValidateClaimsResult, error) {
		accessTokenPayloadUpdate := map[string]interface{}{}
		origSessionClaimPayloadJSON, err := json.Marshal(accessTokenPayload)
		if err != nil {
			return sessmodels.ValidateClaimsResult{}, err
		}

		for _, validator := range claimValidators {
			supertokens.LogDebugMessage("updateClaimsInPayloadIfNeeded checking shouldRefetch for " + validator.ID)
			claim := validator.Claim
			if claim != nil && validator.ShouldRefetch != nil {
				if validator.ShouldRefetch(accessTokenPayload, userContext) {
					supertokens.LogDebugMessage("updateClaimsInPayloadIfNeeded refetching " + validator.ID)
					value, err := claim.FetchValue(userId, userContext)
					if err != nil {
						return sessmodels.ValidateClaimsResult{}, err
					}
					supertokens.LogDebugMessage(fmt.Sprint("updateClaimsInPayloadIfNeeded ", validator.ID, " refetch result ", value))
					if value != nil {
						accessTokenPayload = claim.AddToPayload_internal(accessTokenPayload, value, userContext)
					}
				}
			}
		}

		newSessionClaimPayloadJSON, err := json.Marshal(accessTokenPayload)
		if err != nil {
			return sessmodels.ValidateClaimsResult{}, err
		}
		if !bytes.Equal(origSessionClaimPayloadJSON, newSessionClaimPayloadJSON) {
			accessTokenPayloadUpdate = accessTokenPayload
		}

		invalidClaims := ValidateClaimsInPayload(claimValidators, accessTokenPayload, userContext)

		if len(accessTokenPayloadUpdate) == 0 {
			accessTokenPayloadUpdate = nil
		}

		return sessmodels.ValidateClaimsResult{
			InvalidClaims:            invalidClaims,
			AccessTokenPayloadUpdate: accessTokenPayloadUpdate,
		}, nil
	}

	validateClaimsInJWTPayload := func(userId string, jwtPayload map[string]interface{}, claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.ClaimValidationError, error) {
		invalidClaims := ValidateClaimsInPayload(claimValidators, jwtPayload, userContext)
		return invalidClaims, nil
	}

	fetchAndSetClaim := func(sessionHandle string, claim *claims.TypeSessionClaim, userContext supertokens.UserContext) (bool, error) {
		sessionInfo, err := (*result.GetSessionInformation)(sessionHandle, userContext)
		if err != nil {
			return false, err
		}
		if sessionInfo == nil {
			return false, nil
		}
		accessTokenPayloadUpdate, err := claim.Build(sessionInfo.UserId, nil, userContext)
		if err != nil {
			return false, err
		}
		return (*result.MergeIntoAccessTokenPayload)(sessionHandle, accessTokenPayloadUpdate, userContext)
	}

	setClaimValue := func(sessionHandle string, claim *claims.TypeSessionClaim, value interface{}, userContext supertokens.UserContext) (bool, error) {
		accessTokenPayloadUpdate := claim.AddToPayload_internal(map[string]interface{}{}, value, userContext)
		return (*result.MergeIntoAccessTokenPayload)(sessionHandle, accessTokenPayloadUpdate, userContext)
	}

	getClaimValue := func(sessionHandle string, claim *claims.TypeSessionClaim, userContext supertokens.UserContext) (sessmodels.GetClaimValueResult, error) {
		sessionInfo, err := (*result.GetSessionInformation)(sessionHandle, userContext)
		if err != nil {
			return sessmodels.GetClaimValueResult{}, err
		}
		if sessionInfo == nil {
			return sessmodels.GetClaimValueResult{
				SessionDoesNotExistError: &struct{}{},
			}, nil
		}

		return sessmodels.GetClaimValueResult{
			OK: &struct{ Value interface{} }{
				Value: claim.GetValueFromPayload(sessionInfo.AccessTokenPayload, userContext),
			},
		}, nil
	}

	removeClaim := func(sessionHandle string, claim *claims.TypeSessionClaim, userContext supertokens.UserContext) (bool, error) {
		accessTokenPayloadUpdate := claim.RemoveFromPayloadByMerge_internal(map[string]interface{}{}, userContext)
		return (*result.MergeIntoAccessTokenPayload)(sessionHandle, accessTokenPayloadUpdate, userContext)
	}
	result = sessmodels.RecipeInterface{
		CreateNewSession:            &createNewSession,
		GetSession:                  &getSession,
		RefreshSession:              &refreshSession,
		GetSessionInformation:       &getSessionInformation,
		RevokeAllSessionsForUser:    &revokeAllSessionsForUser,
		GetAllSessionHandlesForUser: &getAllSessionHandlesForUser,
		RevokeSession:               &revokeSession,
		RevokeMultipleSessions:      &revokeMultipleSessions,
		UpdateSessionDataInDatabase: &updateSessionDataInDatabase,
		RegenerateAccessToken:       &regenerateAccessToken,

		MergeIntoAccessTokenPayload: &mergeIntoAccessTokenPayload,
		GetGlobalClaimValidators:    &getGlobalClaimValidators,
		ValidateClaims:              &validateClaims,
		ValidateClaimsInJWTPayload:  &validateClaimsInJWTPayload,
		FetchAndSetClaim:            &fetchAndSetClaim,
		SetClaimValue:               &setClaimValue,
		GetClaimValue:               &getClaimValue,
		RemoveClaim:                 &removeClaim,
	}

	return result
}
