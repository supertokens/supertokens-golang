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

var protectedProps = []string{
	"sub",
	"iat",
	"exp",
	"sessionHandle",
	"parentRefreshTokenHash1",
	"refreshTokenHash1",
	"antiCsrfToken",
}

func makeRecipeImplementation(querier supertokens.Querier, config sessmodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) sessmodels.RecipeInterface {

	// We are defining this here to reduce the scope of legacy code
	const LEGACY_ID_REFRESH_TOKEN_COOKIE_NAME = "sIdRefreshToken"

	var result sessmodels.RecipeInterface

	var recipeImplHandshakeInfo *sessmodels.HandshakeInfo = nil
	getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)

	createNewSession := func(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		supertokens.LogDebugMessage("createNewSession: Started")

		outputTokenTransferMethod := config.GetTokenTransferMethod(req, true, userContext)
		if outputTokenTransferMethod == sessmodels.AnyTransferMethod {
			outputTokenTransferMethod = sessmodels.HeaderTransferMethod
		}

		supertokens.LogDebugMessage(fmt.Sprintf("createNewSession: using transfer method %s", outputTokenTransferMethod))

		isTopLevelAPIDomainIPAddress, err := supertokens.IsAnIPAddress(appInfo.TopLevelAPIDomain)
		if err != nil {
			return nil, err
		}
		isTopLevelWebsiteDomainIPAddress, err := supertokens.IsAnIPAddress(appInfo.TopLevelWebsiteDomain)
		if err != nil {
			return nil, err
		}

		if outputTokenTransferMethod == sessmodels.CookieTransferMethod &&
			config.CookieSameSite == "none" &&
			!config.CookieSecure &&
			!((appInfo.TopLevelAPIDomain == "localhost" || isTopLevelAPIDomainIPAddress) &&
				(appInfo.TopLevelWebsiteDomain == "localhost" || isTopLevelWebsiteDomainIPAddress)) {
			// We can allow insecure cookie when both website & API domain are localhost or an IP
			// When either of them is a different domain, API domain needs to have https and a secure cookie to work
			return nil, defaultErrors.New("Since your API and website domain are different, for sessions to work, please use https on your apiDomain and dont set cookieSecure to false.")
		}

		disableAntiCSRF := outputTokenTransferMethod == sessmodels.HeaderTransferMethod
		sessionResponse, err := createNewSessionHelper(
			recipeImplHandshakeInfo, config, querier, userID, disableAntiCSRF, accessTokenPayload, sessionDataInDatabase,
		)
		if err != nil {
			return nil, err
		}

		for _, tokenTransferMethod := range availableTokenTransferMethods {
			if tokenTransferMethod != outputTokenTransferMethod {
				token, err := getToken(req, sessmodels.AccessToken, tokenTransferMethod)
				if err != nil {
					return nil, err
				}
				if token != nil {
					clearSession(config, res, tokenTransferMethod)
				}
			}
		}

		attachCreateOrRefreshSessionResponseToRes(config, res, sessionResponse, outputTokenTransferMethod)
		payload, parseErr := parseJWTWithoutSignatureVerification(sessionResponse.AccessToken.Token)
		if parseErr != nil {
			return nil, parseErr
		}

		sessionContainerInput := makeSessionContainerInput(sessionResponse.AccessToken.Token, sessionResponse.Session.Handle, sessionResponse.Session.UserID, payload.Payload, res, req, outputTokenTransferMethod, result)
		return newSessionContainer(config, &sessionContainerInput), nil
	}

	// In all cases if sIdRefreshToken token exists (so it's a legacy session) we return TRY_REFRESH_TOKEN. The refresh endpoint will clear this cookie and try to upgrade the session.
	// Check https://supertokens.com/docs/contribute/decisions/session/0007 for further details and a table of expected behaviours
	getSession := func(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		idRefreshToken := getCookieValue(req, LEGACY_ID_REFRESH_TOKEN_COOKIE_NAME)
		if idRefreshToken != nil {
			return nil, errors.TryRefreshTokenError{
				Msg: "using legacy session, please call the refresh API",
			}
		}

		sessionOptional := options != nil && options.SessionRequired != nil && !*options.SessionRequired
		supertokens.LogDebugMessage(fmt.Sprintf("getSession: optional validation %v", sessionOptional))

		accessTokens := map[sessmodels.TokenTransferMethod]*ParsedJWTInfo{}

		// We check all token transfer methods for available access tokens
		for _, tokenTransferMethod := range availableTokenTransferMethods {
			token, err := getToken(req, sessmodels.AccessToken, tokenTransferMethod)
			if err != nil {
				return nil, err
			}
			if token != nil {
				parsedToken, err := parseJWTWithoutSignatureVerification(*token)
				if err != nil {
					supertokens.LogDebugMessage(fmt.Sprintf("getSession: ignoring token in %s, because token parsing failed", tokenTransferMethod))
				} else {
					err := validateAccessTokenStructure(parsedToken.Payload, parsedToken.Version)
					if err != nil {
						supertokens.LogDebugMessage(fmt.Sprintf("getSession: ignoring token in %s, because it doesn't match our access token structure", tokenTransferMethod))
					} else {
						supertokens.LogDebugMessage(fmt.Sprintf("getSession: got access token from %s", tokenTransferMethod))
						accessTokens[tokenTransferMethod] = &parsedToken
					}
				}
			}
		}

		allowedTokenTransferMethod := config.GetTokenTransferMethod(req, false, userContext)

		var requestTokenTransferMethod sessmodels.TokenTransferMethod
		var accessToken *ParsedJWTInfo

		if (allowedTokenTransferMethod == sessmodels.AnyTransferMethod || allowedTokenTransferMethod == sessmodels.HeaderTransferMethod) && (accessTokens[sessmodels.HeaderTransferMethod] != nil) {
			supertokens.LogDebugMessage("getSession: using header transfer method")
			requestTokenTransferMethod = sessmodels.HeaderTransferMethod
			accessToken = accessTokens[sessmodels.HeaderTransferMethod]
		} else if (allowedTokenTransferMethod == sessmodels.AnyTransferMethod || allowedTokenTransferMethod == sessmodels.CookieTransferMethod) && (accessTokens[sessmodels.CookieTransferMethod] != nil) {
			supertokens.LogDebugMessage("getSession: using cookie transfer method")
			requestTokenTransferMethod = sessmodels.CookieTransferMethod
			accessToken = accessTokens[sessmodels.CookieTransferMethod]
		} else {
			if sessionOptional {
				supertokens.LogDebugMessage("getSession: returning undefined because accessToken is undefined and sessionRequired is false")
				return nil, nil
			}

			supertokens.LogDebugMessage("getSession: UNAUTHORISED because accessToken in request is undefined")
			False := false
			return nil, errors.UnauthorizedError{
				Msg: "Session does not exist. Are you sending the session tokens in the request as with the appropriate token transfer method?",
				// we do not clear the session here because of a
				// race condition mentioned here: https://github.com/supertokens/supertokens-node/issues/17
				ClearTokens: &False,
			}
		}

		antiCsrfToken := getAntiCsrfTokenFromHeaders(req)
		var doAntiCsrfCheck *bool

		if options != nil {
			doAntiCsrfCheck = options.AntiCsrfCheck
		}

		if doAntiCsrfCheck == nil {
			doAntiCsrfCheckBool := req.Method != http.MethodGet
			doAntiCsrfCheck = &doAntiCsrfCheckBool
		}

		if requestTokenTransferMethod == sessmodels.HeaderTransferMethod {
			False := false
			doAntiCsrfCheck = &False
		}

		supertokens.LogDebugMessage("getSession: Value of doAntiCsrfCheck is: " + strconv.FormatBool(*doAntiCsrfCheck))

		response, err := getSessionHelper(recipeImplHandshakeInfo, config, querier, *accessToken, antiCsrfToken, *doAntiCsrfCheck, getRidFromHeader(req) != nil)
		if err != nil {
			return nil, err
		}

		accessTokenStr := accessToken.RawTokenString

		if !reflect.DeepEqual(response.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
			tokenError := SetAccessTokenInResponse(config, res, response.AccessToken, response.Session, requestTokenTransferMethod)
			if tokenError != nil {
				return nil, tokenError
			}
			accessTokenStr = response.AccessToken.Token
		}

		supertokens.LogDebugMessage("getSession: Success!")
		var payload map[string]interface{}

		if accessToken.Version < 3 {
			payload = response.Session.UserDataInAccessToken
		} else {
			if reflect.DeepEqual(response.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
				payload = accessToken.Payload
			} else {
				parsedToken, parseErr := parseJWTWithoutSignatureVerification(response.AccessToken.Token)
				if parseErr != nil {
					return nil, parseErr
				}

				payload = parsedToken.Payload
			}
		}

		sessionContainerInput := makeSessionContainerInput(accessTokenStr, response.Session.Handle, response.Session.UserID, payload, res, req, requestTokenTransferMethod, result)
		sessionContainer := newSessionContainer(config, &sessionContainerInput)

		return sessionContainer, nil
	}

	getSessionInformation := func(sessionHandle string, userContext supertokens.UserContext) (*sessmodels.SessionInformation, error) {
		return getSessionInformationHelper(querier, sessionHandle)
	}

	refreshSession := func(req *http.Request, res http.ResponseWriter, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		supertokens.LogDebugMessage("refreshSession: Started")

		refreshTokens := map[sessmodels.TokenTransferMethod]*string{}
		// We check all token transfer methods for available refresh tokens
		// We do this so that we can later clear all we are not overwriting
		for _, tokenTransferMethod := range availableTokenTransferMethods {
			token, err := getToken(req, sessmodels.RefreshToken, tokenTransferMethod)
			if err != nil {
				return nil, err
			}
			refreshTokens[tokenTransferMethod] = token
			if token != nil {
				supertokens.LogDebugMessage("refreshSession: got refresh token from " + string(tokenTransferMethod))
			}
		}

		allowedTokenTransferMethod := config.GetTokenTransferMethod(req, false, userContext)
		supertokens.LogDebugMessage("refreshSession: getTokenTransferMethod returned " + string(allowedTokenTransferMethod))

		var requestTokenTransferMethod sessmodels.TokenTransferMethod
		var refreshToken *string

		if (allowedTokenTransferMethod == sessmodels.AnyTransferMethod || allowedTokenTransferMethod == sessmodels.HeaderTransferMethod) && refreshTokens[sessmodels.HeaderTransferMethod] != nil {
			supertokens.LogDebugMessage("refreshSession: using header transfer method")
			requestTokenTransferMethod = sessmodels.HeaderTransferMethod
			refreshToken = refreshTokens[sessmodels.HeaderTransferMethod]
		} else if (allowedTokenTransferMethod == sessmodels.AnyTransferMethod || allowedTokenTransferMethod == sessmodels.CookieTransferMethod) && refreshTokens[sessmodels.CookieTransferMethod] != nil {
			supertokens.LogDebugMessage("refreshSession: using cookie transfer method")
			requestTokenTransferMethod = sessmodels.CookieTransferMethod
			refreshToken = refreshTokens[sessmodels.CookieTransferMethod]
		} else {
			if getCookieValue(req, LEGACY_ID_REFRESH_TOKEN_COOKIE_NAME) != nil {
				supertokens.LogDebugMessage("refreshSession: cleared legacy id refresh token because refresh token was not found")
				setCookie(config, res, LEGACY_ID_REFRESH_TOKEN_COOKIE_NAME, "", 0, "accessTokenPath")
			}

			supertokens.LogDebugMessage("refreshSession: UNAUTHORISED because refresh token in request is undefined")
			False := false
			return nil, errors.UnauthorizedError{
				Msg:         "Refresh token not found. Are you sending the refresh token in the request as a cookie?",
				ClearTokens: &False,
			}
		}

		antiCsrfToken := getAntiCsrfTokenFromHeaders(req)
		response, err := refreshSessionHelper(recipeImplHandshakeInfo, config, querier, *refreshToken, antiCsrfToken, getRidFromHeader(req) != nil, requestTokenTransferMethod)
		if err != nil {
			unauthorisedErr := errors.UnauthorizedError{}
			isUnauthorisedErr := defaultErrors.As(err, &unauthorisedErr)
			isTokenTheftDetectedErr := defaultErrors.As(err, &errors.TokenTheftDetectedError{})

			// This token isn't handled by getToken/setToken to limit the scope of this legacy/migration code
			if (isTokenTheftDetectedErr) || (isUnauthorisedErr && unauthorisedErr.ClearTokens != nil && *unauthorisedErr.ClearTokens) {
				if getCookieValue(req, LEGACY_ID_REFRESH_TOKEN_COOKIE_NAME) != nil {
					supertokens.LogDebugMessage("refreshSession: cleared legacy id refresh token because refresh is clearing other tokens")
					setCookie(config, res, LEGACY_ID_REFRESH_TOKEN_COOKIE_NAME, "", 0, "accessTokenPath")
				}
			}
			return nil, err
		}

		supertokens.LogDebugMessage("refreshSession: Attaching refreshed session info as " + string(requestTokenTransferMethod))

		// We clear the tokens in all token transfer methods we are not going to overwrite
		for _, tokenTransferMethod := range availableTokenTransferMethods {
			if tokenTransferMethod != requestTokenTransferMethod && refreshTokens[tokenTransferMethod] != nil {
				clearSession(config, res, tokenTransferMethod)
			}
		}
		attachCreateOrRefreshSessionResponseToRes(config, res, response, requestTokenTransferMethod)
		supertokens.LogDebugMessage("refreshSession: Success!")

		// This token isn't handled by getToken/setToken to limit the scope of this legacy/migration code
		if getCookieValue(req, LEGACY_ID_REFRESH_TOKEN_COOKIE_NAME) != nil {
			supertokens.LogDebugMessage("refreshSession: cleared legacy id refresh token after successfull refresh")
			setCookie(config, res, LEGACY_ID_REFRESH_TOKEN_COOKIE_NAME, "", 0, "accessTokenPath")
		}

		responseToken, parseErr := parseJWTWithoutSignatureVerification(response.AccessToken.Token)
		if parseErr != nil {
			return nil, err
		}

		var payload map[string]interface{}

		if responseToken.Version < 3 {
			payload = response.Session.UserDataInAccessToken
		} else {
			payload = responseToken.Payload
		}

		sessionContainerInput := makeSessionContainerInput(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, payload, res, req, requestTokenTransferMethod, result)
		sessionContainer := newSessionContainer(config, &sessionContainerInput)

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

	updateSessionDataInDatabase := func(sessionHandle string, newSessionData map[string]interface{}, userContext supertokens.UserContext) (bool, error) {
		return updateSessionDataInDatabaseHelper(querier, sessionHandle, newSessionData)
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

		for k, _ := range newAccessTokenPayload {
			if supertokens.DoesSliceContainString(k, protectedProps) {
				delete(newAccessTokenPayload, k)
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

		invalidClaims := validateClaimsInPayload(claimValidators, accessTokenPayload, userContext)

		if len(accessTokenPayloadUpdate) == 0 {
			accessTokenPayloadUpdate = nil
		}

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
		UpdateSessionDataInDatabase: &updateSessionDataInDatabase,
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
