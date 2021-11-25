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
	"math"
	"net/http"

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
		// TODO: refresh session
	}

	{
		// TODO: update session
	}

	return originalImplementation
}

func addJWTToAccessTokenPayload(accessTokenPayload map[string]interface{}, jwtExpiry uint64, userId string, jwtPropertyName string, appInfo supertokens.NormalisedAppinfo, jwtRecipeImplementation jwtmodels.RecipeInterface) (map[string]interface{}, error) {
	// TODO:
	return map[string]interface{}{}, nil
}
