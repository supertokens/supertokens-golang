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

package usermetadata

import (
	"github.com/supertokens/supertokens-golang/recipe/usermetadata/usermetadatamodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier, config usermetadatamodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) usermetadatamodels.RecipeInterface {
	getUserMetadata := func(userID string, tenantId *string, userContext supertokens.UserContext) (map[string]interface{}, error) {
		response, err := querier.SendGetRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/user/metadata", map[string]string{
			"userId": userID,
		})
		if err != nil {
			return map[string]interface{}{}, err
		}

		return response["metadata"].(map[string]interface{}), nil
	}

	updateUserMetadata := func(userID string, metadataUpdate map[string]interface{}, tenantId *string, userContext supertokens.UserContext) (map[string]interface{}, error) {
		response, err := querier.SendPutRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/user/metadata", map[string]interface{}{
			"userId":         userID,
			"metadataUpdate": metadataUpdate,
		})
		if err != nil {
			return map[string]interface{}{}, err
		}

		return response["metadata"].(map[string]interface{}), nil
	}

	clearUserMetadata := func(userID string, tenantId *string, userContext supertokens.UserContext) error {
		_, err := querier.SendPostRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/user/metadata/remove", map[string]interface{}{
			"userId": userID,
		})
		return err
	}

	return usermetadatamodels.RecipeInterface{
		GetUserMetadata:    &getUserMetadata,
		UpdateUserMetadata: &updateUserMetadata,
		ClearUserMetadata:  &clearUserMetadata,
	}
}
