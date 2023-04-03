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

package sessionwithjwt

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(originalImplementation sessmodels.RecipeInterface,
	openidRecipeImplementation openidmodels.RecipeInterface, config sessmodels.TypeNormalisedInput) sessmodels.RecipeInterface {

	// Time difference between JWT expiry and access token expiry (JWT expiry = access token expiry + EXPIRY_OFFSET_SECONDS)
	var EXPIRY_OFFSET_SECONDS uint64 = 30

	originalUpdateAccessTokenPayload := *originalImplementation.UpdateAccessTokenPayload

	jwtAwareUpdateAccessTokenPayload := func(sessionInformation *sessmodels.SessionInformation, newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) (bool, error) {
		accessTokenPayload := sessionInformation.AccessTokenPayload
		jwtPropertyName, ok := accessTokenPayload[ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY]

		if !ok {
			return originalUpdateAccessTokenPayload(sessionInformation.SessionHandle, newAccessTokenPayload, userContext)
		}

		existingJWT := accessTokenPayload[jwtPropertyName.(string)].(string)
		currentTimeInSeconds := uint64(time.Now().UnixNano() / 1000000000) // time in seconds
		claims := jwt.MapClaims{}
		decodedPayload := map[string]interface{}{}
		_, _, err := new(jwt.Parser).ParseUnverified(existingJWT, claims)
		if err != nil {
			return false, err
		}
		jwtExpiry := uint64(decodedPayload["exp"].(float64)) - currentTimeInSeconds

		if jwtExpiry <= 0 {
			// it can come here if someone calls this function well after
			// the access token and the jwt payload have expired (which can happen if an API takes a VERY long time). In this case, we still want the jwt payload to update, but the resulting JWT should
			// not be alive for too long (since it's expired already). So we set it to
			// 1 second lifetime.
			jwtExpiry = 1
		}
		newAccessTokenPayload, err = addJWTToAccessTokenPayload(newAccessTokenPayload, jwtExpiry, sessionInformation.UserId, jwtPropertyName.(string), openidRecipeImplementation, userContext)
		if err != nil {
			return false, err
		}
		return originalUpdateAccessTokenPayload(sessionInformation.SessionHandle, newAccessTokenPayload, userContext)
	}

	{
		originalCreateNewSession := *originalImplementation.CreateNewSession

		(*originalImplementation.CreateNewSession) = func(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
			if accessTokenPayload == nil {
				accessTokenPayload = map[string]interface{}{}
			}
			accessTokenValidityInSeconds, err := (*originalImplementation.GetAccessTokenLifeTimeMS)(userContext)
			if err != nil {
				return nil, err
			}
			accessTokenValidityInSeconds = uint64(math.Ceil(float64(accessTokenValidityInSeconds) / 1000))

			accessTokenPayload, err = addJWTToAccessTokenPayload(accessTokenPayload, accessTokenValidityInSeconds+EXPIRY_OFFSET_SECONDS, userID, config.Jwt.PropertyNameInAccessTokenPayload, openidRecipeImplementation, userContext)

			if err != nil {
				return nil, err
			}

			sessionContainer, err := originalCreateNewSession(req, res, userID, accessTokenPayload, sessionDataInDatabase, userContext)

			if err != nil {
				return sessionContainer, err
			}

			return newSessionWithJWTContainer(sessionContainer, openidRecipeImplementation), nil
		}
	}

	{
		originalGetSession := *originalImplementation.GetSession

		(*originalImplementation.GetSession) = func(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
			sessionContainer, err := originalGetSession(req, res, options, userContext)

			if err != nil {
				return nil, err
			}

			if sessionContainer == nil {
				return nil, nil
			}

			result := newSessionWithJWTContainer(sessionContainer, openidRecipeImplementation)

			return result, nil
		}
	}

	{
		originalRefreshSession := *originalImplementation.RefreshSession

		(*originalImplementation.RefreshSession) = func(req *http.Request, res http.ResponseWriter, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
			accessTokenValidityInSeconds, err := (*originalImplementation.GetAccessTokenLifeTimeMS)(userContext)
			if err != nil {
				return nil, err
			}
			accessTokenValidityInSeconds = uint64(math.Ceil(float64(accessTokenValidityInSeconds) / 1000))

			// Refresh session first because this will create a new access token
			newSession, err := originalRefreshSession(req, res, userContext)
			if err != nil {
				return nil, err
			}
			accessTokenPayload := newSession.GetAccessTokenPayloadWithContext(userContext)

			accessTokenPayload, err = addJWTToAccessTokenPayload(accessTokenPayload, accessTokenValidityInSeconds+EXPIRY_OFFSET_SECONDS, newSession.GetUserIDWithContext(userContext), config.Jwt.PropertyNameInAccessTokenPayload, openidRecipeImplementation, userContext)

			if err != nil {
				return nil, err
			}

			err = (newSession.UpdateAccessTokenPayloadWithContext)(accessTokenPayload, userContext)
			if err != nil {
				return nil, err
			}

			return newSessionWithJWTContainer(newSession, openidRecipeImplementation), nil
		}
	}

	{
		(*originalImplementation.UpdateAccessTokenPayload) = func(sessionHandle string, newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) (bool, error) {
			if newAccessTokenPayload == nil {
				newAccessTokenPayload = map[string]interface{}{}
			}
			sessionInformation, err := (*originalImplementation.GetSessionInformation)(sessionHandle, userContext)
			if err != nil {
				return false, err
			}
			if sessionInformation == nil {
				return false, nil
			}

			return jwtAwareUpdateAccessTokenPayload(sessionInformation, newAccessTokenPayload, userContext)
		}
	}

	{
		(*originalImplementation.MergeIntoAccessTokenPayload) = func(sessionHandle string, accessTokenPayloadUpdate map[string]interface{}, userContext supertokens.UserContext) (bool, error) {
			sessionInfo, err := (*originalImplementation.GetSessionInformation)(sessionHandle, userContext)
			if err != nil {
				return false, err
			}

			if sessionInfo == nil {
				return false, nil
			}

			newAccessTokenPayload := sessionInfo.AccessTokenPayload
			for k, v := range accessTokenPayloadUpdate {
				if v == nil {
					delete(newAccessTokenPayload, k)
				} else {
					newAccessTokenPayload[k] = v
				}
			}

			return jwtAwareUpdateAccessTokenPayload(sessionInfo, newAccessTokenPayload, userContext)
		}
	}

	return originalImplementation
}

func addJWTToAccessTokenPayload(accessTokenPayload map[string]interface{}, jwtExpiry uint64, userId string, jwtPropertyName string, openidRecipeImplementation openidmodels.RecipeInterface, userContext supertokens.UserContext) (map[string]interface{}, error) {

	// If jwtPropertyName is not undefined it means that the JWT was added to the access token payload already
	existingJwtPropertyName, ok := accessTokenPayload[ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY]

	if ok {
		delete(accessTokenPayload, existingJwtPropertyName.(string))
		delete(accessTokenPayload, ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY)
	}

	payloadInJWT := map[string]interface{}{
		"sub": userId,
	}
	for k, v := range accessTokenPayload {
		payloadInJWT[k] = v
	}

	jwtResponse, err := (*openidRecipeImplementation.CreateJWT)(payloadInJWT, &jwtExpiry, userContext, nil)
	if err != nil {
		return map[string]interface{}{}, err
	}

	if jwtResponse.UnsupportedAlgorithmError != nil {
		// Should never come here
		return map[string]interface{}{}, errors.New("JWT Signing key algorithm not supported")
	}

	accessTokenPayload[jwtPropertyName] = jwtResponse.OK.Jwt
	accessTokenPayload[ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY] = jwtPropertyName

	return accessTokenPayload, nil
}
