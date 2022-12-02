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

package multitenancy

import (
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *multitenancymodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func CreateOrUpdateThirdPartyConfigWithContext(tenantId *string, thirdPartyId string, config tpmodels.ProviderConfig, userContext supertokens.UserContext) (multitenancymodels.CreateOrUpdateTenantIdConfigResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.CreateOrUpdateTenantIdConfigResponse{}, err
	}
	return (*instance.RecipeImpl.CreateOrUpdateThirdPartyConfig)(tenantId, thirdPartyId, config, userContext)
}

func FetchThirdPartyConfigWithContext(tenantId *string, thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.FetchTenantIdConfigResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.FetchTenantIdConfigResponse{}, err
	}
	return (*instance.RecipeImpl.FetchThirdPartyConfig)(tenantId, thirdPartyId, userContext)
}

func DeleteThirdPartyConfigWithContext(tenantId *string, thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.DeleteTenantIdConfigResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.DeleteTenantIdConfigResponse{}, err
	}
	return (*instance.RecipeImpl.DeleteThirdPartyConfig)(tenantId, thirdPartyId, userContext)
}

func ListThirdPartyConfigsWithContext(tenantId *string, userContext supertokens.UserContext) (multitenancymodels.ListTenantConfigMappingsResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.ListTenantConfigMappingsResponse{}, err
	}
	return (*instance.RecipeImpl.ListThirdPartyConfigs)(tenantId, userContext)
}

func CreateOrUpdateThirdPartyConfig(tenantId *string, thirdPartyId string, config tpmodels.ProviderConfig) (multitenancymodels.CreateOrUpdateTenantIdConfigResponse, error) {
	return CreateOrUpdateThirdPartyConfigWithContext(tenantId, thirdPartyId, config, &map[string]interface{}{})
}

func FetchThirdPartyConfig(tenantId *string, thirdPartyId string) (multitenancymodels.FetchTenantIdConfigResponse, error) {
	return FetchThirdPartyConfigWithContext(tenantId, thirdPartyId, &map[string]interface{}{})
}

func DeleteThirdPartyConfig(tenantId *string, thirdPartyId string) (multitenancymodels.DeleteTenantIdConfigResponse, error) {
	return DeleteThirdPartyConfigWithContext(tenantId, thirdPartyId, &map[string]interface{}{})
}

func ListThirdPartyConfigs(tenantId *string) (multitenancymodels.ListTenantConfigMappingsResponse, error) {
	return ListThirdPartyConfigsWithContext(tenantId, &map[string]interface{}{})
}
