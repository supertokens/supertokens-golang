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
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config *openidmodels.TypeInput) (openidmodels.TypeNormalisedInput, error) {
	result := makeTypeNormalisedInput(appInfo)

	if config != nil {
		if config.Issuer != nil {
			iD, err := supertokens.NewNormalisedURLDomain(*config.Issuer)
			if err != nil {
				return openidmodels.TypeNormalisedInput{}, err
			}
			iP, err := supertokens.NewNormalisedURLPath(*config.Issuer)
			if err != nil {
				return openidmodels.TypeNormalisedInput{}, err
			}
			result.IssuerDomain = iD
			result.IssuerPath = iP
		}

		if result.IssuerPath.GetAsStringDangerous() != appInfo.APIBasePath.GetAsStringDangerous() {
			return openidmodels.TypeNormalisedInput{}, errors.New("The path of the issuer URL must be equal to the apiBasePath. The default value is /auth")
		}
	}

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			result.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			result.Override.APIs = config.Override.APIs
		}
		result.Override.JwtFeature = config.Override.JwtFeature
	}

	return result, nil
}

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) openidmodels.TypeNormalisedInput {
	return openidmodels.TypeNormalisedInput{
		IssuerDomain: appInfo.APIDomain,
		IssuerPath:   appInfo.APIBasePath,
		Override: openidmodels.OverrideStruct{
			Functions: func(originalImplementation openidmodels.RecipeInterface) openidmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation openidmodels.APIInterface) openidmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}
