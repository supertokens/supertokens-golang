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
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *tpmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func ManuallyCreateOrUpdateUserWithContext(thirdPartyID string, thirdPartyUserID string, email string, userContext supertokens.UserContext) (tpmodels.ManuallyCreateOrUpdateUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.ManuallyCreateOrUpdateUserResponse{}, err
	}
	return (*instance.RecipeImpl.ManuallyCreateOrUpdateUser)(thirdPartyID, thirdPartyUserID, email, userContext)
}

func GetUserByIDWithContext(userID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUsersByEmailWithContext(email string, userContext supertokens.UserContext) ([]tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return []tpmodels.User{}, err
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, userContext)
}

func GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, userContext)
}

func CreateOrUpdateThirdPartyConfigWithContext(thirdPartyId string, tenantId *string, config tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.CreateOrUpdateTenantIdConfigResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.CreateOrUpdateTenantIdConfigResponse{}, err
	}
	return (*instance.RecipeImpl.CreateOrUpdateThirdPartyConfig)(thirdPartyId, tenantId, config, userContext)
}

func FetchThirdPartyConfigWithContext(thirdPartyId string, tenantId *string, userContext supertokens.UserContext) (tpmodels.FetchTenantIdConfigResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.FetchTenantIdConfigResponse{}, err
	}
	return (*instance.RecipeImpl.FetchThirdPartyConfig)(thirdPartyId, tenantId, userContext)
}

func DeleteThirdPartyConfigWithContext(thirdPartyId string, tenantId *string, userContext supertokens.UserContext) (tpmodels.DeleteTenantIdConfigResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.DeleteTenantIdConfigResponse{}, err
	}
	return (*instance.RecipeImpl.DeleteThirdPartyConfig)(thirdPartyId, tenantId, userContext)
}

func ListThirdPartyConfigsWithContext(tenantId *string, userContext supertokens.UserContext) (tpmodels.ListTenantConfigMappingsResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.ListTenantConfigMappingsResponse{}, err
	}
	return (*instance.RecipeImpl.ListThirdPartyConfigs)(tenantId, userContext)
}

func ManuallyCreateOrUpdateUser(thirdPartyID string, thirdPartyUserID string, email string) (tpmodels.ManuallyCreateOrUpdateUserResponse, error) {
	return ManuallyCreateOrUpdateUserWithContext(thirdPartyID, thirdPartyUserID, email, &map[string]interface{}{})
}

func GetUserByID(userID string) (*tpmodels.User, error) {
	return GetUserByIDWithContext(userID, &map[string]interface{}{})
}

func GetUsersByEmail(email string) ([]tpmodels.User, error) {
	return GetUsersByEmailWithContext(email, &map[string]interface{}{})
}

func GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID string) (*tpmodels.User, error) {
	return GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID, &map[string]interface{}{})
}

func CreateOrUpdateThirdPartyConfig(thirdPartyId string, tenantId *string, config tpmodels.ProviderConfig) (tpmodels.CreateOrUpdateTenantIdConfigResponse, error) {
	return CreateOrUpdateThirdPartyConfigWithContext(thirdPartyId, tenantId, config, &map[string]interface{}{})
}

func FetchThirdPartyConfig(thirdPartyId string, tenantId *string) (tpmodels.FetchTenantIdConfigResponse, error) {
	return FetchThirdPartyConfigWithContext(thirdPartyId, tenantId, &map[string]interface{}{})
}

func DeleteThirdPartyConfig(thirdPartyId string, tenantId *string) (tpmodels.DeleteTenantIdConfigResponse, error) {
	return DeleteThirdPartyConfigWithContext(thirdPartyId, tenantId, &map[string]interface{}{})
}

func ListThirdPartyConfigs(tenantId *string) (tpmodels.ListTenantConfigMappingsResponse, error) {
	return ListThirdPartyConfigsWithContext(tenantId, &map[string]interface{}{})
}

func ActiveDirectory(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.ActiveDirectory(input)
}

func Apple(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.Apple(input)
}

func BoxySaml(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.BoxySaml(input)
}

func Discord(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.Discord(input)
}

func Facebook(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.Facebook(input)
}

func Github(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.Github(input)
}

func Google(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.Google(input)
}

func GoogleWorkspaces(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.GoogleWorkspaces(input)
}

func Linkedin(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.Linkedin(input)
}

func Okta(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.Okta(input)
}

func CustomProvider(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	return providers.NewProvider(input)
}
