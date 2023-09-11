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

package dashboard

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"strings"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config *dashboardmodels.TypeInput) dashboardmodels.TypeNormalisedInput {
	typeNormalisedInput := makeTypeNormalisedInput(appInfo)

	_config := dashboardmodels.TypeInput{}
	if config != nil {
		_config = *config
	}

	if _config.ApiKey != "" {
		typeNormalisedInput.ApiKey = _config.ApiKey
		typeNormalisedInput.AuthMode = dashboardmodels.AuthModeAPIKey
	}

	if _config.Override != nil {
		if _config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = _config.Override.Functions
		}
		if _config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = _config.Override.APIs
		}
	}

	if _config.ApiKey != "" && config.Admins != nil {
		supertokens.LogDebugMessage("User Dashboard: Providing 'Admins' has no effect when using an apiKey.")
	}

	var admins *[]string
	if _config.Admins != nil {
		admins = _config.Admins
	}

	typeNormalisedInput.Admins = admins

	return typeNormalisedInput
}

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) dashboardmodels.TypeNormalisedInput {
	return dashboardmodels.TypeNormalisedInput{
		AuthMode: dashboardmodels.AuthModeEmailPassword,
		Override: dashboardmodels.OverrideStruct{
			Functions: func(originalImplementation dashboardmodels.RecipeInterface) dashboardmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation dashboardmodels.APIInterface) dashboardmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}

func normaliseEmail(email string) string {
	_email := strings.TrimSpace(email)
	_email = strings.ToLower(_email)
	return _email
}
