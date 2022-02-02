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
	defaultErrors "errors"
	"net/http"
	"reflect"
	"sync"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var handshakeInfoLock sync.Mutex

func makeRecipeImplementation(querier supertokens.Querier, config sessmodels.TypeNormalisedInput) sessmodels.RecipeInterface {

	var recipeImplHandshakeInfo *sessmodels.HandshakeInfo = nil
	getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)

	createNewSession := func(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		response, err := createNewSessionHelper(recipeImplHandshakeInfo, config, querier, userID, accessTokenPayload, sessionData)
		if err != nil {
			return sessmodels.SessionContainer{}, err
		}
		attachCreateOrRefreshSessionResponseToRes(config, res, response)
		sessionContainerInput := makeSessionContainerInput(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInAccessToken, res)
		return newSessionContainer(querier, config, &sessionContainerInput), nil
	}

	getSession := func(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (*sessmodels.SessionContainer, error) {
		var doAntiCsrfCheck *bool = nil
		if options != nil {
			doAntiCsrfCheck = options.AntiCsrfCheck
		}

		idRefreshToken := getIDRefreshTokenFromCookie(req)
		if idRefreshToken == nil {
			if options != nil && options.SessionRequired != nil &&
				!(*options.SessionRequired) {
				return nil, nil
			}
			return nil, errors.UnauthorizedError{Msg: "Session does not exist. Are you sending the session tokens in the request as cookies?"}
		}

		accessToken := getAccessTokenFromCookie(req)
		if accessToken == nil {
			if options == nil || (options.SessionRequired != nil && *options.SessionRequired) || frontendHasInterceptor(req) || req.Method == http.MethodGet {
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

		response, err := getSessionHelper(recipeImplHandshakeInfo, config, querier, *accessToken, antiCsrfToken, *doAntiCsrfCheck, getRidFromHeader(req) != nil)
		if err != nil {
			if defaultErrors.As(err, &errors.UnauthorizedError{}) {
				clearSessionFromCookie(config, res)
			}
			return nil, err
		}

		if !reflect.DeepEqual(response.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
			setFrontTokenInHeaders(res, response.Session.UserID, response.AccessToken.Expiry, response.Session.UserDataInAccessToken)
			attachAccessTokenToCookie(config, res, response.AccessToken.Token, response.AccessToken.Expiry)
			accessToken = &response.AccessToken.Token
		}
		sessionContainerInput := makeSessionContainerInput(*accessToken, response.Session.Handle, response.Session.UserID, response.Session.UserDataInAccessToken, res)
		sessionContainer := newSessionContainer(querier, config, &sessionContainerInput)
		return &sessionContainer, nil
	}

	getSessionInformation := func(sessionHandle string, userContext supertokens.UserContext) (sessmodels.SessionInformation, error) {
		return getSessionInformationHelper(querier, sessionHandle)
	}

	refreshSession := func(req *http.Request, res http.ResponseWriter, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
		inputIdRefreshToken := getIDRefreshTokenFromCookie(req)
		if inputIdRefreshToken == nil {
			return sessmodels.SessionContainer{}, errors.UnauthorizedError{Msg: "Session does not exist. Are you sending the session tokens in the request as cookies?"}
		}

		inputRefreshToken := getRefreshTokenFromCookie(req)
		if inputRefreshToken == nil {
			clearSessionFromCookie(config, res)
			return sessmodels.SessionContainer{}, errors.UnauthorizedError{Msg: "Refresh token not found. Are you sending the refresh token in the request as a cookie?"}
		}

		antiCsrfToken := getAntiCsrfTokenFromHeaders(req)
		response, err := refreshSessionHelper(recipeImplHandshakeInfo, config, querier, *inputRefreshToken, antiCsrfToken, getRidFromHeader(req) != nil)
		if err != nil {
			// we clear cookies if it is UnauthorizedError & ClearCookies in it is nil or true
			// we clear cookies if it is TokenTheftDetectedError
			if (defaultErrors.As(err, &errors.UnauthorizedError{}) && (err.(errors.UnauthorizedError).ClearCookies == nil || *err.(errors.UnauthorizedError).ClearCookies)) || defaultErrors.As(err, &errors.TokenTheftDetectedError{}) {
				clearSessionFromCookie(config, res)
			}
			return sessmodels.SessionContainer{}, err
		}
		attachCreateOrRefreshSessionResponseToRes(config, res, response)
		sessionContainerInput := makeSessionContainerInput(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInAccessToken, res)
		sessionContainer := newSessionContainer(querier, config, &sessionContainerInput)
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

	updateSessionData := func(sessionHandle string, newSessionData map[string]interface{}, userContext supertokens.UserContext) error {
		return updateSessionDataHelper(querier, sessionHandle, newSessionData)
	}

	updateAccessTokenPayload := func(sessionHandle string, newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) error {
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

	regenerateAccessToken := func(accessToken string, newAccessTokenPayload *map[string]interface{}, userContext supertokens.UserContext) (sessmodels.RegenerateAccessTokenResponse, error) {
		return regenerateAccessTokenHelper(querier, newAccessTokenPayload, accessToken)
	}

	return sessmodels.RecipeInterface{
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
	}

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
