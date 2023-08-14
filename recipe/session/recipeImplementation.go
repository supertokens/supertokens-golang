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
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var protectedProps = []string{
	"sub",
	"iat",
	"exp",
	"sessionHandle",
	"parentRefreshTokenHash1",
	"refreshTokenHash1",
	"antiCsrfToken",
	"tId",
}

var JWKCacheMaxAgeInMs int64 = 60000
var JWKRefreshRateLimit = 500
var jwksCache *sessmodels.GetJWKSResult = nil
var mutex sync.RWMutex

func getJWKSFromCacheIfPresent() *sessmodels.GetJWKSResult {
	mutex.RLock()
	defer mutex.RUnlock()
	if jwksCache != nil {
		// This means that we have valid JWKs for the given core path
		// We check if we need to refresh before returning
		currentTime := time.Now().UnixNano() / int64(time.Millisecond)

		// This means that the value in cache is not expired, in this case we return the cached value
		//
		// Note that this also means that the SDK will not try to query any other Core (if there are multiple)
		// if it has a valid cache entry from one of the core URLs. It will only attempt to fetch
		// from the cores again after the entry in the cache is expired
		if (currentTime - jwksCache.LastFetched) < JWKCacheMaxAgeInMs {
			if supertokens.IsRunningInTestMode() {
				if len(returnedFromCache) == cap(returnedFromCache) { // need to clear the channel if full because it's not being consumed in the test
					close(returnedFromCache)
					returnedFromCache = make(chan bool, 1000)
				}
				returnedFromCache <- true
			}

			return jwksCache
		}
	}

	return nil
}

func getJWKS() (*keyfunc.JWKS, error) {
	corePaths := supertokens.GetAllCoreUrlsForPath("/.well-known/jwks.json")

	if len(corePaths) == 0 {
		return nil, defaultErrors.New("No SuperTokens core available to query. Please pass supertokens > connectionURI to the init function, or override all the functions of the recipe you are using.")
	}

	resultFromCache := getJWKSFromCacheIfPresent()

	if resultFromCache != nil {
		return resultFromCache.JWKS, nil
	}

	var lastError error

	mutex.Lock()
	defer mutex.Unlock()
	for _, path := range corePaths {
		if supertokens.IsRunningInTestMode() {
			urlsAttemptedForJWKSFetch = append(urlsAttemptedForJWKSFetch, path)
		}

		// RefreshUnknownKID - Fetch JWKS again if the kid in the header of the JWT does not match any in
		// the keyfunc library's cache
		jwks, jwksError := keyfunc.Get(path, keyfunc.Options{
			RefreshUnknownKID: true,
		})

		if jwksError == nil {
			jwksResult := sessmodels.GetJWKSResult{
				JWKS:        jwks,
				Error:       jwksError,
				LastFetched: time.Now().UnixNano() / int64(time.Millisecond),
			}

			// Dont add to cache if there is an error to keep the logic of checking cache simple
			//
			// This also has the added benefit where if initially the request failed because the core
			// was down and then it comes back up, the next time it will try to request that core again
			// after the cache has expired
			jwksCache = &jwksResult

			if supertokens.IsRunningInTestMode() {
				if len(returnedFromCache) == cap(returnedFromCache) { // need to clear the channel if full because it's not being consumed in the test
					close(returnedFromCache)
					returnedFromCache = make(chan bool, 1000)
				}
				returnedFromCache <- false
			}

			return jwksResult.JWKS, nil
		}

		lastError = jwksError
	}

	// This means that fetching from all cores failed
	return nil, lastError
}

/*
*
This function fetches all JWKs from the first available core instance. This combines the other JWKS functions to become
error resistant.

Every core instance a backend is connected to is expected to connect to the same database and use the same key set for
token verification. Otherwise, the result of session verification would depend on which core is currently available.
*/
func GetCombinedJWKS() (*keyfunc.JWKS, error) {
	if supertokens.IsRunningInTestMode() {
		urlsAttemptedForJWKSFetch = []string{}
	}

	jwksResult, err := getJWKS()

	if err != nil {
		return nil, err
	}

	return jwksResult, nil
}

func MakeRecipeImplementation(querier supertokens.Querier, config sessmodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) sessmodels.RecipeInterface {
	var result sessmodels.RecipeInterface

	createNewSession := func(userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, disableAntiCsrf *bool, tenantId string, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		supertokens.LogDebugMessage("createNewSession: Started")

		sessionResponse, err := createNewSessionHelper(
			config, querier, userID, disableAntiCsrf != nil && *disableAntiCsrf == true, accessTokenPayload, sessionDataInDatabase, tenantId,
		)
		if err != nil {
			return nil, err
		}

		supertokens.LogDebugMessage("createNewSession: Finished")

		parsedJWT, parseErr := ParseJWTWithoutSignatureVerification(sessionResponse.AccessToken.Token)
		if parseErr != nil {
			return nil, parseErr
		}

		frontToken := BuildFrontToken(sessionResponse.Session.UserID, sessionResponse.AccessToken.Expiry, parsedJWT.Payload)
		session := sessionResponse.Session
		sessionContainerInput := makeSessionContainerInput(sessionResponse.AccessToken.Token, session.Handle, session.UserID, session.TenantId, parsedJWT.Payload, result, frontToken, sessionResponse.AntiCsrfToken, nil, &sessionResponse.RefreshToken, true)
		return newSessionContainer(config, &sessionContainerInput), nil
	}

	// In all cases if sIdRefreshToken token exists (so it's a legacy session) we return TRY_REFRESH_TOKEN. The refresh endpoint will clear this cookie and try to upgrade the session.
	// Check https://supertokens.com/docs/contribute/decisions/session/0007 for further details and a table of expected behaviours
	getSession := func(accessTokenString *string, antiCsrfToken *string, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		if options != nil && options.AntiCsrfCheck != nil && *options.AntiCsrfCheck != false && config.AntiCsrf == AntiCSRF_VIA_CUSTOM_HEADER {
			return nil, defaultErrors.New("Since the anti-csrf mode is VIA_CUSTOM_HEADER getSession can't check the CSRF token. Please either use VIA_TOKEN or set antiCsrfCheck to false")
		}

		supertokens.LogDebugMessage("getSession: Started")

		if accessTokenString == nil {
			if options != nil && options.SessionRequired != nil && *options.SessionRequired == false {
				supertokens.LogDebugMessage("getSession: returning nil because accessToken is nil and sessionRequired is false")
				return nil, nil
			}

			supertokens.LogDebugMessage("getSession: UNAUTHORISED because accessToken in request is nil")
			False := false
			return nil, errors.UnauthorizedError{
				Msg: "Session does not exist. Are you sending the session tokens in the request with the appropriate token transfer method?",
				// we do not clear the session here because of a
				// race condition mentioned here: https://github.com/supertokens/supertokens-node/issues/17
				ClearTokens: &False,
			}
		}

		var accessToken *sessmodels.ParsedJWTInfo

		accessTokenResponse, err := ParseJWTWithoutSignatureVerification(*accessTokenString)

		if err != nil {
			if options != nil && *options.SessionRequired == false {
				supertokens.LogDebugMessage("getSession: Returning nil because parsing failed and sessionRequired is false")
				return nil, nil
			}

			supertokens.LogDebugMessage("getSession: Returning UNAUTHORISED because parsing failed")
			return nil, errors.UnauthorizedError{
				Msg:         "Token parsing failed",
				ClearTokens: nil,
			}
		}

		accessToken = &accessTokenResponse
		err = ValidateAccessTokenStructure(accessTokenResponse.Payload, accessTokenResponse.Version)

		if err != nil {
			if options != nil && *options.SessionRequired == false {
				supertokens.LogDebugMessage("getSession: Returning nil because parsing failed and sessionRequired is false")
				return nil, nil
			}

			supertokens.LogDebugMessage("getSession: Returning UNAUTHORISED because parsing failed")
			return nil, errors.UnauthorizedError{
				Msg:         "Token parsing failed",
				ClearTokens: nil,
			}
		}

		alwaysCheckCore := false

		if options != nil && options.CheckDatabase != nil {
			alwaysCheckCore = *options.CheckDatabase == true
		}

		doAntiCsrfCheck := true

		if options != nil && options.AntiCsrfCheck != nil && *options.AntiCsrfCheck == false {
			doAntiCsrfCheck = false
		}

		response, err := getSessionHelper(config, querier, *accessToken, antiCsrfToken, doAntiCsrfCheck, alwaysCheckCore)
		if err != nil {
			return nil, err
		}

		supertokens.LogDebugMessage("getSession: Success!")
		var payload map[string]interface{}

		if accessToken.Version >= 3 {
			if !reflect.DeepEqual(response.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
				parsedToken, parseErr := ParseJWTWithoutSignatureVerification(response.AccessToken.Token)

				if parseErr != nil {
					return nil, parseErr
				}

				payload = parsedToken.Payload
			} else {
				payload = accessToken.Payload
			}
		} else {
			payload = response.Session.UserDataInAccessToken
		}

		accessTokenStringForSession := *accessTokenString

		accessTokenNil := reflect.DeepEqual(response.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{})

		if !accessTokenNil {
			accessTokenStringForSession = response.AccessToken.Token
		}

		frontToken := BuildFrontToken(response.Session.UserID, response.Session.ExpiryTime, payload)
		session := response.Session

		sessionContainerInput := makeSessionContainerInput(accessTokenStringForSession, session.Handle, session.UserID, session.TenantId, payload, result, frontToken, antiCsrfToken, nil, nil, !accessTokenNil)
		sessionContainer := newSessionContainer(config, &sessionContainerInput)

		return sessionContainer, nil
	}

	getSessionInformation := func(sessionHandle string, userContext supertokens.UserContext) (*sessmodels.SessionInformation, error) {
		return getSessionInformationHelper(querier, sessionHandle)
	}

	refreshSession := func(refreshToken string, antiCsrfToken *string, disableAntiCsrf bool, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		if disableAntiCsrf != true && config.AntiCsrf == AntiCSRF_VIA_CUSTOM_HEADER {
			return nil, defaultErrors.New("Since the anti-csrf mode is VIA_CUSTOM_HEADER getSession can't check the CSRF token. Please either use VIA_TOKEN or set antiCsrfCheck to false")
		}

		supertokens.LogDebugMessage("refreshSession: Started")

		response, err := refreshSessionHelper(config, querier, refreshToken, antiCsrfToken, disableAntiCsrf)
		if err != nil {
			return nil, err
		}
		supertokens.LogDebugMessage("refreshSession: Success!")

		responseToken, parseErr := ParseJWTWithoutSignatureVerification(response.AccessToken.Token)
		if parseErr != nil {
			return nil, err
		}

		session := response.Session
		frontToken := BuildFrontToken(session.UserID, response.AccessToken.Expiry, responseToken.Payload)

		sessionContainerInput := makeSessionContainerInput(response.AccessToken.Token, session.Handle, session.UserID, session.TenantId, responseToken.Payload, result, frontToken, response.AntiCsrfToken, nil, &response.RefreshToken, true)
		sessionContainer := newSessionContainer(config, &sessionContainerInput)

		return sessionContainer, nil
	}

	revokeAllSessionsForUser := func(userID string, tenantId string, revokeAcrossAllTenants *bool, userContext supertokens.UserContext) ([]string, error) {
		return revokeAllSessionsForUserHelper(querier, userID, tenantId, revokeAcrossAllTenants)
	}

	getAllSessionHandlesForUser := func(userID string, tenantId string, fetchAcrossAllTenants *bool, userContext supertokens.UserContext) ([]string, error) {
		return getAllSessionHandlesForUserHelper(querier, userID, tenantId, fetchAcrossAllTenants)
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
		for k, v := range sessionInfo.CustomClaimsInAccessTokenPayload {
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

	getGlobalClaimValidators := func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, tenantId string, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
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
					tenantId, ok := accessTokenPayload["tId"].(string)
					if !ok {
						tenantId = multitenancymodels.DefaultTenantId
					}
					value, err := claim.FetchValue(userId, tenantId, userContext)
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
		accessTokenPayloadUpdate, err := claim.Build(sessionInfo.UserId, sessionInfo.TenantId, nil, userContext)
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
				Value: claim.GetValueFromPayload(sessionInfo.CustomClaimsInAccessTokenPayload, userContext),
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
