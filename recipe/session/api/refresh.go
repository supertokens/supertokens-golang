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

package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func HandleRefreshAPI(apiImplementation sessmodels.APIInterface, options sessmodels.APIOptions) error {
	if apiImplementation.RefreshPOST == nil || (*apiImplementation.RefreshPOST) == nil {
		options.OtherHandler.ServeHTTP(options.Res, options.Req)
		return nil
	}

	userContext := supertokens.MakeDefaultUserContextFromAPI(options.Req)

	var tenantId *string = nil

	{
		// Fetch tenant id
		bodyBytes, err := ioutil.ReadAll(options.Req.Body)
		if err != nil {
			return err
		}
		if len(bodyBytes) > 0 {
			bodyObj := map[string]interface{}{}
			err = json.Unmarshal(bodyBytes, &bodyObj)
			if err != nil {
				return err
			}
			if tenantIdVal, ok := bodyObj["tenantId"].(string); ok {
				tenantId = &tenantIdVal
			}
		}
	}

	mtRecipe, err := multitenancy.GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	tenantId, err = (*mtRecipe.RecipeImpl.GetTenantId)(tenantId, userContext)
	if err != nil {
		return err
	}

	_, err = (*apiImplementation.RefreshPOST)(tenantId, options, userContext)
	if err != nil {
		return err
	}
	return supertokens.Send200Response(options.Res, nil)
}
