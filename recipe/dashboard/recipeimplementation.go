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
	"github.com/supertokens/supertokens-golang/recipe/dashboard/constants"
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
			}, userContext)

			if err != nil {
				return false, err
			}

			status, ok := verifyResponse["status"]

			if !ok || status != "OK" {
				return false, nil
			}

			// For all non GET requests we also want to check if the user is allowed to perform this operation
			if req.Method != http.MethodGet {
				// We dont want to block the analytics API
				if strings.HasSuffix(req.RequestURI, constants.DashboardAnalyticsAPI) {
					return true, nil
				}

				// We do not want to block the sign out request
				if strings.HasSuffix(req.RequestURI, constants.SignOutAPI) {
					return true, nil
				}

				admins := config.Admins

				if admins == nil {
					return true, nil
				}

				if len(*admins) == 0 {
					supertokens.LogDebugMessage("User Dashboard: Throwing OPERATION_NOT_ALLOWED because user is not an admin")
					return false, errors.ForbiddenAccessError{
						Msg: "You are not permitted to perform this operation",
					}
				}

				userEmail, emailOk := verifyResponse["email"]

				if !emailOk || userEmail.(string) == "" {
					supertokens.LogDebugMessage("User Dashboard: Returning Unauthorised because no email was returned from the core. Should never come here")
					return false, nil
				}

				if !supertokens.DoesSliceContainString(userEmail.(string), *admins) {
					supertokens.LogDebugMessage("User Dashboard: Throwing OPERATION_NOT_ALLOWED because user is not an admin")
					return false, errors.ForbiddenAccessError{
						Msg: "You are not permitted to perform this operation",
					}
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
