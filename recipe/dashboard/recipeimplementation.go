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
	"fmt"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/errors"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/validationUtils"
	"github.com/supertokens/supertokens-golang/supertokens"
	"net/http"
	"strings"
)

func makeRecipeImplementation(querier supertokens.Querier) dashboardmodels.RecipeInterface {

	getDashboardBundleLocation := func(userContext supertokens.UserContext) (string, error) {
		return fmt.Sprintf("https://cdn.jsdelivr.net/gh/supertokens/dashboard@v%s/build/", supertokens.DashboardVersion), nil
	}

	shouldAllowAccess := func(req *http.Request, config dashboardmodels.TypeNormalisedInput, userContext supertokens.UserContext) (bool, error) {
		if config.ApiKey == "" {
			authHeaderValue := req.Header.Get("authorization")
			// We receive the api key as `Bearer API_KEY`, this retrieves just the key
			keyParts := strings.Split(authHeaderValue, " ")
			authHeaderValue = keyParts[len(keyParts)-1]

			verifyResponse, err := querier.SendPostRequest("/recipe/dashboard/session/verify", map[string]interface{}{
				"sessionId": authHeaderValue,
			})

			if err != nil {
				return false, err
			}

			status, ok := verifyResponse["status"]

			if !ok || status != "OK" {
				return false, nil
			}

			// For all non GET requests we also want to check if the user is allowed to perform this operation
			if req.Method != "GET" {
				// We dont want to block the analytics API
				if strings.HasSuffix(req.RequestURI, dashboardAnalyticsAPI) {
					return true, nil
				}

				admins := config.Admins

				// If the user has provided no admins, allow
				if len(admins) == 0 {
					return true, nil
				}

				emailFromResponse, emailOk := verifyResponse["email"]

				if !emailOk || (emailFromResponse.(string)) == "" {
					supertokens.LogDebugMessage("User Dashboard: You are using an older version of SuperTokens core, to use the 'admins' property when initialising the Dashboard recipe please upgrade to a core version that is >= 6.0")
					return false, errors.OperationNotAllowedError{}
				}

				if !supertokens.DoesSliceContainString(emailFromResponse.(string), admins) {
					return false, errors.OperationNotAllowedError{}
				}
			}

			return true, nil
		}

		validateKeyResponse, err := validationUtils.ValidateApiKey(req, config, userContext)

		if err != nil {
			return false, err
		}

		return validateKeyResponse, nil
	}

	return dashboardmodels.RecipeInterface{
		GetDashboardBundleLocation: &getDashboardBundleLocation,
		ShouldAllowAccess:          &shouldAllowAccess,
	}
}
