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
	"net/http"
	"reflect"
	"strconv"
	"sync"

	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var handshakeInfoLock sync.Mutex

func makeRecipeImplementation(querier supertokens.Querier, config sessmodels.TypeNormalisedInput) sessmodels.RecipeInterface {

	var result sessmodels.RecipeInterface

	var recipeImplHandshakeInfo *sessmodels.HandshakeInfo = nil
	getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)

	createNewSession := func(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		response, err := createNewSessionHelper(recipeImplHandshakeInfo, config, querier, userID, accessTokenPayload, sessionData)
		if err != nil {
			return nil, err
		}
		attachCreateOrRefreshSessionResponseToRes(config, res, response)
		sessionContainerInput := makeSessionContainerInput(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInAccessToken, res, result)
		return newSessionContainer(config, &sessionContainerInput), nil
	}

	getSession := func(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		supertokens.LogDebugMessage("getSession: Started")
		supertokens.LogDebugMessage("getSession: rid in header: " + strconv.FormatBool(frontendHasInterceptor((req))))
		supertokens.LogDebugMessage("getSession: request method: " + req.Method)

		var doAntiCsrfCheck *bool = nil
		if options != nil {
			doAntiCsrfCheck = options.AntiCsrfCheck
		}

		idRefreshToken := getIDRefreshTokenFromCookie(req)
		if idRefreshToken == nil {
			if options != nil && options.SessionRequired != nil &&
				!(*options.SessionRequired) {
				supertokens.LogDebugMessage("getSession: returning nil because idRefreshToken is nil and sessionRequired is false")
				return nil, nil
			}
			supertokens.LogDebugMessage("getSession: UNAUTHORISED because idRefreshToken from cookies is nil")
			return nil, errors.UnauthorizedError{Msg: "Session does not exist. Are you sending the session tokens in the request as cookies?"}
		}

		accessToken := getAccessTokenFromCookie(req)
		if accessToken == nil {
			if options == nil || (options.SessionRequired != nil && *options.SessionRequired) || frontendHasInterceptor(req) || req.Method == http.MethodGet {
				supertokens.LogDebugMessage("getSession: Returning try refresh token because access token from cookies is nil")
				return nil, errors.TryRefreshTokenError{
					Msg: "Access token has expired. Please call the refresh API",
				}
			}
			return nil, nil
		}

		antiCsrfToken := getAntiCsrfTokenFromHeaders(req)
		if doAntiCsrfCheck == nil {
			doAntiCsrfCheckBool := req.Method != http.MethodGet
			doAntiCsrfCheck = &doAntiCsrfCheckBool
		}

		if doAntiCsrfCheck != nil {
			supertokens.LogDebugMessage("getSession: Value of doAntiCsrfCheck is: " + strconv.FormatBool(*doAntiCsrfCheck))
		} else {
			supertokens.LogDebugMessage("getSession: Value of doAntiCsrfCheck is: nil")
		}

		response, err := getSessionHelper(recipeImplHandshakeInfo, config, querier, *accessToken, antiCsrfToken, *doAntiCsrfCheck, getRidFromHeader(req) != nil)
		if err != nil {
			if defaultErrors.As(err, &errors.UnauthorizedError{}) {
				supertokens.LogDebugMessage("getSession: Clearing cookies because of UNAUTHORISED response")
				clearSessionFromCookie(config, res)
			}
			return nil, err
		}

		if !reflect.DeepEqual(response.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
			setFrontTokenInHeaders(res, response.Session.UserID, response.AccessToken.Expiry, response.Session.UserDataInAccessToken)
			attachAccessTokenToCookie(config, res, response.AccessToken.Token, response.AccessToken.Expiry)
			accessToken = &response.AccessToken.Token
		}
		sessionContainerInput := makeSessionContainerInput(*accessToken, response.Session.Handle, response.Session.UserID, response.Session.UserDataInAccessToken, res, result)
		sessionContainer := newSessionContainer(config, &sessionContainerInput)

		supertokens.LogDebugMessage("getSession: Success!")
		return sessionContainer, nil
	}

	getSessionInformation := func(sessionHandle string, userContext supertokens.UserContext) (*sessmodels.SessionInformation, error) {
		return getSessionInformationHelper(querier, sessionHandle)
	}

	refreshSession := func(req *http.Request, res http.ResponseWriter, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		supertokens.LogDebugMessage("refreshSession: Started")
		inputIdRefreshToken := getIDRefreshTokenFromCookie(req)
		if inputIdRefreshToken == nil {
			supertokens.LogDebugMessage("refreshSession: UNAUTHORISED because idRefreshToken from cookies is nil")
			return nil, errors.UnauthorizedError{Msg: "Session does not exist. Are you sending the session tokens in the request as cookies?"}
		}

		inputRefreshToken := getRefreshTokenFromCookie(req)
		if inputRefreshToken == nil {
			clearSessionFromCookie(config, res)
			supertokens.LogDebugMessage("refreshSession: UNAUTHORISED because refresh token from cookies is undefined")
			return nil, errors.UnauthorizedError{Msg: "Refresh token not found. Are you sending the refresh token in the request as a cookie?"}
		}

		antiCsrfToken := getAntiCsrfTokenFromHeaders(req)
		response, err := refreshSessionHelper(recipeImplHandshakeInfo, config, querier, *inputRefreshToken, antiCsrfToken, getRidFromHeader(req) != nil)
		if err != nil {
			// we clear cookies if it is UnauthorizedError & ClearCookies in it is nil or true
			// we clear cookies if it is TokenTheftDetectedError
			if (defaultErrors.As(err, &errors.UnauthorizedError{}) && (err.(errors.UnauthorizedError).ClearCookies == nil || *err.(errors.UnauthorizedError).ClearCookies)) || defaultErrors.As(err, &errors.TokenTheftDetectedError{}) {
				supertokens.LogDebugMessage("refreshSession: Clearing cookies because of UNAUTHORISED or TOKEN_THEFT_DETECTED response")
				clearSessionFromCookie(config, res)
			}
			return nil, err
		}
		attachCreateOrRefreshSessionResponseToRes(config, res, response)
		sessionContainerInput := makeSessionContainerInput(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInAccessToken, res, result)
		sessionContainer := newSessionContainer(config, &sessionContainerInput)

		supertokens.LogDebugMessage("refreshSession: Success!")
		return sessionContainer, nil
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

	updateSessionData := func(sessionHandle string, newSessionData map[string]interface{}, userContext supertokens.UserContext) (bool, error) {
		return updateSessionDataHelper(querier, sessionHandle, newSessionData)
	}

	updateAccessTokenPayload := func(sessionHandle string, newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) (bool, error) {
		return updateAccessTokenPayloadHelper(querier, sessionHandle, newAccessTokenPayload)
	}

	getAccessTokenLifeTimeMS := func(userContext supertokens.UserContext) (uint64, error) {
		err := getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)
		if err != nil {
			return 0, err
		}
		return recipeImplHandshakeInfo.AccessTokenValidity, nil
	}

	getRefreshTokenLifeTimeMS := func(userContext supertokens.UserContext) (uint64, error) {
		err := getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)
		if err != nil {
			return 0, err
		}
		return recipeImplHandshakeInfo.RefreshTokenValidity, nil
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
		for k, v := range accessTokenPayloadUpdate {
			if v == nil {
				delete(newAccessTokenPayload, k)
			} else {
				newAccessTokenPayload[k] = v
			}
		}
		return (*result.UpdateAccessTokenPayload)(sessionHandle, newAccessTokenPayload, userContext)
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

		if len(accessTokenPayload) == 0 {
			accessTokenPayload = nil
		}

		invalidClaims := validateClaimsInPayload(claimValidators, accessTokenPayload, userContext)
		return sessmodels.ValidateClaimsResult{
			InvalidClaims:            invalidClaims,
			AccessTokenPayloadUpdate: accessTokenPayloadUpdate,
		}, nil
	}

	validateClaimsInJWTPayload := func(userId string, jwtPayload map[string]interface{}, claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.ClaimValidationError, error) {
		invalidClaims := validateClaimsInPayload(claimValidators, jwtPayload, userContext)
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
		UpdateSessionData:           &updateSessionData,
		UpdateAccessTokenPayload:    &updateAccessTokenPayload,
		GetAccessTokenLifeTimeMS:    &getAccessTokenLifeTimeMS,
		GetRefreshTokenLifeTimeMS:   &getRefreshTokenLifeTimeMS,
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

// updates recipeImplHandshakeInfo in place.
func getHandshakeInfo(recipeImplHandshakeInfo **sessmodels.HandshakeInfo, config sessmodels.TypeNormalisedInput, querier supertokens.Querier, forceFetch bool) error {
	handshakeInfoLock.Lock()
	defer handshakeInfoLock.Unlock()
	if *recipeImplHandshakeInfo == nil ||
		len((*recipeImplHandshakeInfo).GetJwtSigningPublicKeyList()) == 0 ||
		forceFetch {
		response, err := querier.SendPostRequest("/recipe/handshake", nil)
		if err != nil {
			return err
		}

		*recipeImplHandshakeInfo = &sessmodels.HandshakeInfo{
			AntiCsrf:                       config.AntiCsrf,
			AccessTokenBlacklistingEnabled: response["accessTokenBlacklistingEnabled"].(bool),
			AccessTokenValidity:            uint64(response["accessTokenValidity"].(float64)),
			RefreshTokenValidity:           uint64(response["refreshTokenValidity"].(float64)),
		}

		updateJwtSigningPublicKeyInfoWithoutLock(recipeImplHandshakeInfo, getKeyInfoFromJson(response), response["jwtSigningPublicKey"].(string), uint64(response["jwtSigningPublicKeyExpiryTime"].(float64)))

	}
	return nil
}

func updateJwtSigningPublicKeyInfoWithoutLock(recipeImplHandshakeInfo **sessmodels.HandshakeInfo, keyList []sessmodels.KeyInfo, newKey string, newExpiry uint64) {
	if len(keyList) == 0 {
		// means we are using an older CDI version
		keyList = []sessmodels.KeyInfo{
			{
				PublicKey:  newKey,
				ExpiryTime: newExpiry,
				CreatedAt:  getCurrTimeInMS(),
			},
		}
	}

	if *recipeImplHandshakeInfo != nil {
		(*recipeImplHandshakeInfo).SetJwtSigningPublicKeyList(keyList)
	}

}

func updateJwtSigningPublicKeyInfo(recipeImplHandshakeInfo **sessmodels.HandshakeInfo, keyList []sessmodels.KeyInfo, newKey string, newExpiry uint64) {
	handshakeInfoLock.Lock()
	defer handshakeInfoLock.Unlock()
	updateJwtSigningPublicKeyInfoWithoutLock(recipeImplHandshakeInfo, keyList, newKey, newExpiry)
}
