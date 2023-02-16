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

func GetUserMetadata(userID string, tenantId *string) (map[string]interface{}, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return map[string]interface{}{}, err
	}
	return (*instance.RecipeImpl.GetUserMetadata)(userID, tenantId, &map[string]interface{}{})
}

func GetUserMetadataWithContext(userID string, tenantId *string, userContext supertokens.UserContext) (map[string]interface{}, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return map[string]interface{}{}, err
	}
	return (*instance.RecipeImpl.GetUserMetadata)(userID, tenantId, userContext)
}

func UpdateUserMetadata(userID string, metadataUpdate map[string]interface{}, tenantId *string) (map[string]interface{}, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return map[string]interface{}{}, err
	}
	return (*instance.RecipeImpl.UpdateUserMetadata)(userID, metadataUpdate, tenantId, &map[string]interface{}{})
}

func UpdateUserMetadataWithContext(userID string, metadataUpdate map[string]interface{}, tenantId *string, userContext supertokens.UserContext) (map[string]interface{}, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return map[string]interface{}{}, err
	}
	return (*instance.RecipeImpl.UpdateUserMetadata)(userID, metadataUpdate, tenantId, userContext)
}

func ClearUserMetadata(userID string, tenantId *string) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.RecipeImpl.ClearUserMetadata)(userID, tenantId, &map[string]interface{}{})
}

func ClearUserMetadataWithContext(userID string, tenantId *string, userContext supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.RecipeImpl.ClearUserMetadata)(userID, tenantId, userContext)
}
