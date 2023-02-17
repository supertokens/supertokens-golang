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

func ManuallyCreateOrUpdateUserWithContext(thirdPartyID string, thirdPartyUserID string, email string, tenantId *string, userContext supertokens.UserContext) (tpmodels.ManuallyCreateOrUpdateUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.ManuallyCreateOrUpdateUserResponse{}, err
	}
	return (*instance.RecipeImpl.ManuallyCreateOrUpdateUser)(thirdPartyID, thirdPartyUserID, email, tenantId, userContext)
}

func GetUserByIDWithContext(userID string, tenantId *string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, tenantId, userContext)
}

func GetUsersByEmailWithContext(email string, tenantId *string, userContext supertokens.UserContext) ([]tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return []tpmodels.User{}, err
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, tenantId, userContext)
}

func GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID string, tenantId *string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, tenantId, userContext)
}

func GetProviderWithContext(thirdPartyID string, tenantId *string, clientType *string, userContext supertokens.UserContext) (tpmodels.GetProviderResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.GetProviderResponse{}, err
	}
	return (*instance.RecipeImpl.GetProvider)(thirdPartyID, tenantId, clientType, userContext)
}

func ManuallyCreateOrUpdateUser(thirdPartyID string, thirdPartyUserID string, email string, tenantId *string) (tpmodels.ManuallyCreateOrUpdateUserResponse, error) {
	return ManuallyCreateOrUpdateUserWithContext(thirdPartyID, thirdPartyUserID, email, tenantId, &map[string]interface{}{})
}

func GetUserByID(userID string, tenantId *string) (*tpmodels.User, error) {
	return GetUserByIDWithContext(userID, tenantId, &map[string]interface{}{})
}

func GetUsersByEmail(email string, tenantId *string) ([]tpmodels.User, error) {
	return GetUsersByEmailWithContext(email, tenantId, &map[string]interface{}{})
}

func GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID string, tenantId *string) (*tpmodels.User, error) {
	return GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID, tenantId, &map[string]interface{}{})
}

func GetProvider(thirdPartyID string, tenantId *string, clientType *string) (tpmodels.GetProviderResponse, error) {
	return GetProviderWithContext(thirdPartyID, tenantId, clientType, &map[string]interface{}{})
}
