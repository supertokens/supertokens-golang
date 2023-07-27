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

package thirdparty

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *tpmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func ManuallyCreateOrUpdateUser(tenantId string, thirdPartyID string, thirdPartyUserID string, email string, userContext ...supertokens.UserContext) (tpmodels.ManuallyCreateOrUpdateUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.ManuallyCreateOrUpdateUserResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ManuallyCreateOrUpdateUser)(thirdPartyID, thirdPartyUserID, email, tenantId, userContext[0])
}

func GetUserByID(userID string, userContext ...supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext[0])
}

func GetUsersByEmail(tenantId string, email string, userContext ...supertokens.UserContext) ([]tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return []tpmodels.User{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, tenantId, userContext[0])
}

func GetUserByThirdPartyInfo(tenantId string, thirdPartyID, thirdPartyUserID string, userContext ...supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, tenantId, userContext[0])
}

func GetProvider(tenantId string, thirdPartyID string, clientType *string, userContext ...supertokens.UserContext) (*tpmodels.TypeProvider, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetProvider)(thirdPartyID, clientType, tenantId, userContext[0])
}
