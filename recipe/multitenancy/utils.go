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

package multitenancy

import (
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
)

func validateAndNormaliseUserInput(config *multitenancymodels.TypeInput) multitenancymodels.TypeNormalisedInput {

	typeNormalisedInput := makeTypeNormalisedInput(config)

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

func makeTypeNormalisedInput(config *multitenancymodels.TypeInput) multitenancymodels.TypeNormalisedInput {
	if config == nil {
		config = &multitenancymodels.TypeInput{}
	}

	return multitenancymodels.TypeNormalisedInput{
		GetAllowedDomainsForTenantId: config.GetAllowedDomainsForTenantId,
		Override: multitenancymodels.OverrideStruct{
			Functions: func(originalImplementation multitenancymodels.RecipeInterface) multitenancymodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation multitenancymodels.APIInterface) multitenancymodels.APIInterface {
				return originalImplementation
			},
		},
	}
}
