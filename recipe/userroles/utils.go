/* Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package userroles

import (
	"github.com/supertokens/supertokens-golang/recipe/userroles/userrolesmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config *userrolesmodels.TypeInput) userrolesmodels.TypeNormalisedInput {

	typeNormalisedInput := makeTypeNormalisedInput(appInfo)

	if config != nil {
		typeNormalisedInput.SkipAddingRolesToAccessToken = config.SkipAddingRolesToAccessToken
		typeNormalisedInput.SkipAddingPermissionsToAccessToken = config.SkipAddingPermissionsToAccessToken
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

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) userrolesmodels.TypeNormalisedInput {
	return userrolesmodels.TypeNormalisedInput{
		Override: userrolesmodels.OverrideStruct{
			Functions: func(originalImplementation userrolesmodels.RecipeInterface) userrolesmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation userrolesmodels.APIInterface) userrolesmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}

func convertToStringArray(arr []interface{}) []string {
	result := make([]string, len(arr))
	for idx, v := range arr {
		result[idx] = v.(string)
	}
	return result
}
