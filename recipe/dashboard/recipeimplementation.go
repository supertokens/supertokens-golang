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
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier, config dashboardmodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) dashboardmodels.RecipeInterface {

	getDashboardBundleLocation := func(userContext supertokens.UserContext) (string, error) {
		return fmt.Sprintf("https://cdn.jsdelivr.net/gh/supertokens/dashboard@v%s/build/", supertokens.DashboardVersion), nil
	}

	shouldAllowAccess := func(req *http.Request, userContext supertokens.UserContext) (bool, error) {
		apiKeyHeaderValue := req.Header.Get("authorization")

		// We receieve the api key as `Bearer API_KEY`, this retrieves just the key
		keyParts := strings.Split(apiKeyHeaderValue, " ")
		apiKeyHeaderValue = keyParts[len(keyParts)-1]

		if apiKeyHeaderValue == "" {
			return false, nil
		}

		return apiKeyHeaderValue == config.ApiKey, nil
	}

	return dashboardmodels.RecipeInterface{
		GetDashboardBundleLocation: &getDashboardBundleLocation,
		ShouldAllowAccess:          &shouldAllowAccess,
	}
}
