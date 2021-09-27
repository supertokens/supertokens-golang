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

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config *jwtmodels.TypeInput) jwtmodels.TypeNormalisedInput {

	typeNormalisedInput := makeTypeNormalisedInput(appInfo)

	if config != nil && config.JwtValiditySeconds != nil {
		typeNormalisedInput.JwtValiditySeconds = *config.JwtValiditySeconds
	}

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
	}

	return typeNormalisedInput
}

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) jwtmodels.TypeNormalisedInput {
	return jwtmodels.TypeNormalisedInput{
		JwtValiditySeconds: 3153600000, // 100 years in seconds
		Override: jwtmodels.OverrideStruct{
			Functions: func(originalImplementation jwtmodels.RecipeInterface) jwtmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation jwtmodels.APIInterface) jwtmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}
