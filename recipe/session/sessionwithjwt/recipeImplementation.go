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
	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(originalImplementation sessmodels.RecipeInterface,
	jwtRecipeImplementation jwtmodels.RecipeInterface, config sessmodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) sessmodels.RecipeInterface {

	// Time difference between JWT expiry and access token expiry (JWT expiry = access token expiry + EXPIRY_OFFSET_SECONDS)
	var EXPIRY_OFFSET_SECONDS uint64 = 30

	{
		originalCreateNewSession := *originalImplementation.CreateNewSession

		(*originalImplementation.CreateNewSession) = func(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}) (sessmodels.SessionContainer, error) {
			if accessTokenPayload == nil {
				accessTokenPayload = map[string]interface{}{}
			}
			accessTokenValidityInSeconds, err := (*originalImplementation.GetAccessTokenLifeTimeMS)()
			if err != nil {
				return sessmodels.SessionContainer{}, err
			}
			accessTokenValidityInSeconds = uint64(math.Ceil(float64(accessTokenValidityInSeconds) / 1000))

			accessTokenPayload, err = addJWTToAccessTokenPayload(accessTokenPayload, accessTokenValidityInSeconds+EXPIRY_OFFSET_SECONDS, userID, config.Jwt.PropertyNameInAccessTokenPayload, appInfo, jwtRecipeImplementation)

			if err != nil {
				return sessmodels.SessionContainer{}, err
			}

			sessionContainer, err := originalCreateNewSession(res, userID, accessTokenPayload, sessionData)

			if err != nil {
				return sessionContainer, err
			}

			return newSessionWithJWTContainer(sessionContainer, jwtRecipeImplementation, appInfo), nil
		}
	}

	{
		originalGetSession := *originalImplementation.GetSession

		(*originalImplementation.GetSession) = func(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions) (*sessmodels.SessionContainer, error) {
			sessionContainer, err := originalGetSession(req, res, options)

			if err != nil {
				return nil, err
			}

			if sessionContainer == nil {
				return nil, nil
			}

			result := newSessionWithJWTContainer(*sessionContainer, jwtRecipeImplementation, appInfo)

			return &result, nil
		}
	}

	{
		originalRefreshSession := *originalImplementation.RefreshSession

		(*originalImplementation.RefreshSession) = func(req *http.Request, res http.ResponseWriter) (sessmodels.SessionContainer, error) {
			accessTokenValidityInSeconds, err := (*originalImplementation.GetAccessTokenLifeTimeMS)()
			if err != nil {
				return sessmodels.SessionContainer{}, err
			}
			accessTokenValidityInSeconds = uint64(math.Ceil(float64(accessTokenValidityInSeconds) / 1000))

			// Refresh session first because this will create a new access token
			newSession, err := originalRefreshSession(req, res)
			if err != nil {
				return sessmodels.SessionContainer{}, err
			}
			accessTokenPayload := newSession.GetAccessTokenPayload()

			accessTokenPayload, err = addJWTToAccessTokenPayload(accessTokenPayload, accessTokenValidityInSeconds+EXPIRY_OFFSET_SECONDS, newSession.GetUserID(), config.Jwt.PropertyNameInAccessTokenPayload, appInfo, jwtRecipeImplementation)

			if err != nil {
				return sessmodels.SessionContainer{}, err
			}

			err = newSession.UpdateAccessTokenPayload(accessTokenPayload)
			if err != nil {
				return sessmodels.SessionContainer{}, err
			}

			return newSessionWithJWTContainer(newSession, jwtRecipeImplementation, appInfo), nil
		}
	}

	{
		originalUpdateAccessTokenPayload := *originalImplementation.UpdateAccessTokenPayload

		(*originalImplementation.UpdateAccessTokenPayload) = func(sessionHandle string, newAccessTokenPayload map[string]interface{}) error {
			if newAccessTokenPayload == nil {
				newAccessTokenPayload = map[string]interface{}{}
			}
			sessionInformation, err := (*originalImplementation.GetSessionInformation)(sessionHandle)
			if err != nil {
				return err
			}
			accessTokenPayload := sessionInformation.AccessTokenPayload
			jwtPropertyName, ok := accessTokenPayload[ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY]

			if !ok {
				return originalUpdateAccessTokenPayload(sessionHandle, newAccessTokenPayload)
			}

			existingJWT := accessTokenPayload[jwtPropertyName.(string)].(string)

			currentTimeInSeconds := uint64(time.Now().UnixNano() / 1000000000) // time in seconds

			claims := jwt.MapClaims{}
			decodedPayload := map[string]interface{}{}
			_, _, err = new(jwt.Parser).ParseUnverified(existingJWT, claims)
			if err != nil {
				return err
			}
			for key, val := range claims {
				decodedPayload[key] = val
			}

			jwtExpiry := decodedPayload["exp"].(uint64) - currentTimeInSeconds

			if jwtExpiry <= 0 {
				// it can come here if someone calls this function well after
				// the access token and the jwt payload have expired (which can happen if an API takes a VERY long time). In this case, we still want the jwt payload to update, but the resulting JWT should
				// not be alive for too long (since it's expired already). So we set it to
				// 1 second lifetime.
				jwtExpiry = 1
			}

			newAccessTokenPayload, err = addJWTToAccessTokenPayload(newAccessTokenPayload, jwtExpiry, sessionInformation.UserId, jwtPropertyName.(string), appInfo, jwtRecipeImplementation)
			if err != nil {
				return err
			}

			return originalUpdateAccessTokenPayload(sessionHandle, newAccessTokenPayload)
		}
	}

	return originalImplementation
}

func addJWTToAccessTokenPayload(accessTokenPayload map[string]interface{}, jwtExpiry uint64, userId string, jwtPropertyName string, appInfo supertokens.NormalisedAppinfo, jwtRecipeImplementation jwtmodels.RecipeInterface) (map[string]interface{}, error) {

	// If jwtPropertyName is not undefined it means that the JWT was added to the access token payload already
	existingJwtPropertyName, ok := accessTokenPayload[ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY]

	if ok {
		delete(accessTokenPayload, existingJwtPropertyName.(string))
		delete(accessTokenPayload, ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY)
	}

	newAccessTokenPayload := map[string]interface{}{
		"sub": userId,
		"iss": appInfo.APIDomain.GetAsStringDangerous(),
	}
	for k, v := range accessTokenPayload {
		newAccessTokenPayload[k] = v
	}

	jwtResponse, err := (*jwtRecipeImplementation.CreateJWT)(newAccessTokenPayload, &jwtExpiry)
	if err != nil {
		return map[string]interface{}{}, err
	}

	if jwtResponse.UnsupportedAlgorithmError != nil {
		// Should never come here
		return map[string]interface{}{}, errors.New("JWT Signing key algorithm not supported")
	}

	newAccessTokenPayload[jwtPropertyName] = jwtResponse.OK.Jwt
	newAccessTokenPayload[ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY] = jwtPropertyName

	return newAccessTokenPayload, nil
}
