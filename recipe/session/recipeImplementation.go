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

	return sessmodels.RecipeInterface{
		CreateNewSession: func(res http.ResponseWriter, userID string, jwtPayload map[string]interface{}, sessionData map[string]interface{}) (sessmodels.SessionContainer, error) {
			response, err := createNewSessionHelper(recipeImplHandshakeInfo, config, querier, userID, jwtPayload, sessionData)
			if err != nil {
				return sessmodels.SessionContainer{}, err
			}
			attachCreateOrRefreshSessionResponseToRes(config, res, response)
			sessionContainerInput := makeSessionContainerInput(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInJWT, res)
			return newSessionContainer(querier, config, &sessionContainerInput), nil
		},

		GetSession: func(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions) (*sessmodels.SessionContainer, error) {
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
				setFrontTokenInHeaders(res, response.Session.UserID, response.AccessToken.Expiry, response.Session.UserDataInJWT)
				attachAccessTokenToCookie(config, res, response.AccessToken.Token, response.AccessToken.Expiry)
				accessToken = &response.AccessToken.Token
			}
			sessionContainerInput := makeSessionContainerInput(*accessToken, response.Session.Handle, response.Session.UserID, response.Session.UserDataInJWT, res)
			sessionContainer := newSessionContainer(querier, config, &sessionContainerInput)
			return &sessionContainer, nil
		},

		GetSessionInformation: func(sessionHandle string) (sessmodels.SessionInformation, error) {
			return getSessionInformationHelper(querier, sessionHandle)
		},

		RefreshSession: func(req *http.Request, res http.ResponseWriter) (sessmodels.SessionContainer, error) {
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
			sessionContainerInput := makeSessionContainerInput(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInJWT, res)
			sessionContainer := newSessionContainer(querier, config, &sessionContainerInput)
			return sessionContainer, nil
		},

		RevokeAllSessionsForUser: func(userID string) ([]string, error) {
			return revokeAllSessionsForUserHelper(querier, userID)
		},

		GetAllSessionHandlesForUser: func(userID string) ([]string, error) {
			return getAllSessionHandlesForUserHelper(querier, userID)
		},

		RevokeSession: func(sessionHandle string) (bool, error) {
			return revokeSessionHelper(querier, sessionHandle)
		},

		RevokeMultipleSessions: func(sessionHandles []string) ([]string, error) {
			return revokeMultipleSessionsHelper(querier, sessionHandles)
		},

		UpdateSessionData: func(sessionHandle string, newSessionData map[string]interface{}) error {
			return updateSessionDataHelper(querier, sessionHandle, newSessionData)
		},

		UpdateJWTPayload: func(sessionHandle string, newJWTPayload map[string]interface{}) error {
			return updateJWTPayloadHelper(querier, sessionHandle, newJWTPayload)
		},

		GetAccessTokenLifeTimeMS: func() (uint64, error) {
			err := getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)
			if err != nil {
				return 0, err
			}
			return recipeImplHandshakeInfo.AccessTokenValidity, nil
		},

		GetRefreshTokenLifeTimeMS: func() (uint64, error) {
			err := getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)
			if err != nil {
				return 0, err
			}
			return recipeImplHandshakeInfo.RefreshTokenValidity, nil
		},
	}
}

func getHandshakeInfo(recipeImplHandshakeInfo **sessmodels.HandshakeInfo, config sessmodels.TypeNormalisedInput, querier supertokens.Querier, forceFetch bool) error {
	handshakeInfoLock.Lock()
	defer handshakeInfoLock.Unlock()
	if *recipeImplHandshakeInfo == nil || forceFetch {
		response, err := querier.SendPostRequest("/recipe/handshake", nil)
		if err != nil {
			return err
		}
		signingKeyLastUpdated := getCurrTimeInMS()
		if *recipeImplHandshakeInfo != nil {
			if uint64(response["jwtSigningPublicKeyExpiryTime"].(float64)) == (**recipeImplHandshakeInfo).JWTSigningPublicKeyExpiryTime &&
				response["jwtSigningPublicKey"].(string) == (**recipeImplHandshakeInfo).JWTSigningPublicKey {
				signingKeyLastUpdated = (**recipeImplHandshakeInfo).SigningKeyLastUpdated
			}
		}

		*recipeImplHandshakeInfo = &sessmodels.HandshakeInfo{
			SigningKeyLastUpdated:          signingKeyLastUpdated,
			JWTSigningPublicKey:            response["jwtSigningPublicKey"].(string),
			AntiCsrf:                       config.AntiCsrf,
			AccessTokenBlacklistingEnabled: response["accessTokenBlacklistingEnabled"].(bool),
			JWTSigningPublicKeyExpiryTime:  uint64(response["jwtSigningPublicKeyExpiryTime"].(float64)),
			AccessTokenValidity:            uint64(response["accessTokenValidity"].(float64)),
			RefreshTokenValidity:           uint64(response["refreshTokenValidity"].(float64)),
		}
	}
	return nil
}

func updateJwtSigningPublicKeyInfo(recipeImplHandshakeInfo **sessmodels.HandshakeInfo, newKey string, newExpiry uint64) {
	handshakeInfoLock.Lock()
	defer handshakeInfoLock.Unlock()
	if *recipeImplHandshakeInfo != nil {
		if (**recipeImplHandshakeInfo).JWTSigningPublicKeyExpiryTime != newExpiry ||
			(**recipeImplHandshakeInfo).JWTSigningPublicKey != newKey {
			(**recipeImplHandshakeInfo).SigningKeyLastUpdated = getCurrTimeInMS()
		}
		(**recipeImplHandshakeInfo).JWTSigningPublicKey = newKey
		(**recipeImplHandshakeInfo).JWTSigningPublicKeyExpiryTime = newExpiry
	}

}
