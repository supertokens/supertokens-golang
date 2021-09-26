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

package jwt

import (
	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier, config jwtmodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) jwtmodels.RecipeInterface {
	return jwtmodels.RecipeInterface{
		CreateJWT: func(payload map[string]interface{}, validitySecondsPointer *uint64) (jwtmodels.CreateJWTResponse, error) {
			validitySeconds := config.JwtValiditySeconds
			if validitySecondsPointer != nil {
				validitySeconds = *validitySecondsPointer
			}
			if payload == nil {
				payload = map[string]interface{}{}
			}

			response, err := querier.SendPostRequest("/recipe/jwt", map[string]interface{}{
				"payload":    payload,
				"validity":   validitySeconds,
				"algorithm":  "RS256",
				"jwksDomain": appInfo.APIDomain.GetAsStringDangerous(),
			})
			if err != nil {
				return jwtmodels.CreateJWTResponse{}, err
			}

			status, ok := response["status"]
			if ok && status == "OK" {
				return jwtmodels.CreateJWTResponse{
					OK: &struct{ Jwt string }{
						Jwt: response["jwt"].(string),
					},
				}, nil
			} else {
				return jwtmodels.CreateJWTResponse{
					UnsupportedAlgorithmError: &struct{}{},
				}, nil
			}
		},
		GetJWKS: func() (jwtmodels.GetJWKSResponse, error) {
			response, err := querier.SendGetRequest("/recipe/jwt/jwks", map[string]string{})
			if err != nil {
				return jwtmodels.GetJWKSResponse{}, err
			}
			return jwtmodels.GetJWKSResponse{
				OK: &struct{ Keys []jwtmodels.JsonWebKeys }{
					Keys: response["keys"].([]jwtmodels.JsonWebKeys),
				},
			}, nil
		},
	}
}
