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

package openid

import (
	"github.com/supertokens/supertokens-golang/recipe/jwt"
	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(config openidmodels.TypeNormalisedInput, jwtRecipeImplementation jwtmodels.RecipeInterface) openidmodels.RecipeInterface {
	createJWT := func(payload map[string]interface{}, validitySecondsPointer *uint64, userContext supertokens.UserContext, useStaticSigningKey *bool) (jwtmodels.CreateJWTResponse, error) {
		issuer := config.IssuerDomain.GetAsStringDangerous() + config.IssuerPath.GetAsStringDangerous()
		if payload == nil {
			payload = map[string]interface{}{}
		}

		payload["iss"] = issuer
		return (*jwtRecipeImplementation.CreateJWT)(payload, validitySecondsPointer, userContext, useStaticSigningKey)
	}

	getJWKS := func(userContext supertokens.UserContext) (jwtmodels.GetJWKSResponse, error) {
		return (*jwtRecipeImplementation.GetJWKS)(userContext)
	}

	getOpenIdDiscoveryConfiguration := func(userContext supertokens.UserContext) (openidmodels.GetOpenIdDiscoveryConfigurationResponse, error) {
		issuer := config.IssuerDomain.GetAsStringDangerous() + config.IssuerPath.GetAsStringDangerous()
		jwksPath, err := supertokens.NewNormalisedURLPath(jwt.GetJWKSAPI)
		if err != nil {
			return openidmodels.GetOpenIdDiscoveryConfigurationResponse{}, err
		}
		jwks_uri := config.IssuerDomain.GetAsStringDangerous() + config.IssuerPath.AppendPath(jwksPath).GetAsStringDangerous()
		return openidmodels.GetOpenIdDiscoveryConfigurationResponse{
			OK: &struct {
				Issuer   string
				Jwks_uri string
			}{
				Issuer:   issuer,
				Jwks_uri: jwks_uri,
			},
		}, nil
	}

	return openidmodels.RecipeInterface{
		CreateJWT:                       &createJWT,
		GetJWKS:                         &getJWKS,
		GetOpenIdDiscoveryConfiguration: &getOpenIdDiscoveryConfiguration,
	}
}
