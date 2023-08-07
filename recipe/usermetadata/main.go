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

func Init(config *usermetadatamodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func GetUserMetadata(userID string, userContext ...supertokens.UserContext) (map[string]interface{}, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return map[string]interface{}{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserMetadata)(userID, userContext[0])
}

func UpdateUserMetadata(userID string, metadataUpdate map[string]interface{}, userContext ...supertokens.UserContext) (map[string]interface{}, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return map[string]interface{}{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.UpdateUserMetadata)(userID, metadataUpdate, userContext[0])
}

func ClearUserMetadata(userID string, userContext ...supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ClearUserMetadata)(userID, userContext[0])
}
