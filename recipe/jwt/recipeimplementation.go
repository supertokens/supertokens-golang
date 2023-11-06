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
	"regexp"
	"strconv"
)

var defaultJWKSMaxAge = 60 // This corresponds to the dynamicSigningKeyOverlapMS in the core

func makeRecipeImplementation(querier supertokens.Querier, config jwtmodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) jwtmodels.RecipeInterface {
	createJWT := func(payload map[string]interface{}, validitySecondsPointer *uint64, useStaticSigningKey *bool, userContext supertokens.UserContext) (jwtmodels.CreateJWTResponse, error) {
		validitySeconds := config.JwtValiditySeconds
		if validitySecondsPointer != nil {
			validitySeconds = *validitySecondsPointer
		}
		if payload == nil {
			payload = map[string]interface{}{}
		}

		shouldUseStaticSigningKey := true

		if useStaticSigningKey != nil {
			shouldUseStaticSigningKey = *useStaticSigningKey
		}

		response, err := querier.SendPostRequest("/recipe/jwt", map[string]interface{}{
			"payload":             payload,
			"validity":            validitySeconds,
			"algorithm":           "RS256",
			"jwksDomain":          appInfo.APIDomain.GetAsStringDangerous(),
			"useStaticSigningKey": shouldUseStaticSigningKey,
		}, userContext)
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
	}
	getJWKS := func(userContext supertokens.UserContext) (jwtmodels.GetJWKSResponse, error) {
		response, headers, err := querier.SendGetRequestWithResponseHeaders("/.well-known/jwks.json", map[string]string{}, userContext)
		if err != nil {
			return jwtmodels.GetJWKSResponse{}, err
		}

		keys := []jwtmodels.JsonWebKeys{}

		for _, v := range response["keys"].([]interface{}) {
			keys = append(keys, jwtmodels.JsonWebKeys{
				Kty: (v.(map[string]interface{}))["kty"].(string),
				Kid: (v.(map[string]interface{}))["kid"].(string),
				N:   (v.(map[string]interface{}))["n"].(string),
				E:   (v.(map[string]interface{}))["e"].(string),
				Alg: (v.(map[string]interface{}))["alg"].(string),
				Use: (v.(map[string]interface{}))["use"].(string),
			})
		}

		validityInSeconds := defaultJWKSMaxAge
		cacheControlHeader := headers.Get("Cache-Control")

		if cacheControlHeader != "" {
			regex := regexp.MustCompile(`/,?\s*max-age=(\d+)(?:,|$)/`)
			maxAgeHeader := regex.FindAllString(cacheControlHeader, -1)

			if maxAgeHeader != nil && len(maxAgeHeader) > 0 {
				validityInSeconds, err = strconv.Atoi(maxAgeHeader[1])

				if err != nil {
					validityInSeconds = defaultJWKSMaxAge
				}
			}
		}

		return jwtmodels.GetJWKSResponse{
			OK: &struct {
				Keys              []jwtmodels.JsonWebKeys
				ValidityInSeconds int
			}{
				Keys:              keys,
				ValidityInSeconds: validityInSeconds,
			},
		}, nil
	}

	return jwtmodels.RecipeInterface{
		CreateJWT: &createJWT,
		GetJWKS:   &getJWKS,
	}
}
