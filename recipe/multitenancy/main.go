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

func CreateOrUpdateTenantWithContext(tenantId *string, config multitenancymodels.TenantConfig, userContext supertokens.UserContext) (multitenancymodels.CreateOrUpdateTenantResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.CreateOrUpdateTenantResponse{}, err
	}

	return (*instance.RecipeImpl.CreateOrUpdateTenant)(tenantId, config, userContext)
}

func DeleteTenantWithContext(tenantId string, userContext supertokens.UserContext) (multitenancymodels.DeleteTenantResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.DeleteTenantResponse{}, err
	}

	return (*instance.RecipeImpl.DeleteTenant)(tenantId, userContext)
}

func GetTenantConfigWithContext(tenantId *string, userContext supertokens.UserContext) (multitenancymodels.TenantConfigResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.TenantConfigResponse{}, err
	}

	return (*instance.RecipeImpl.GetTenantConfig)(tenantId, userContext)
}

func ListAllTenantsWithContext(userContext supertokens.UserContext) (multitenancymodels.ListAllTenantsResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.ListAllTenantsResponse{}, err
	}

	return (*instance.RecipeImpl.ListAllTenants)(userContext)
}

// Third party provider management
func CreateOrUpdateThirdPartyConfigWithContext(config tpmodels.ProviderConfig, skipValidation bool, userContext supertokens.UserContext) (multitenancymodels.CreateOrUpdateThirdPartyConfigResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.CreateOrUpdateThirdPartyConfigResponse{}, err
	}
	return (*instance.RecipeImpl.CreateOrUpdateThirdPartyConfig)(config, skipValidation, userContext)
}

func DeleteThirdPartyConfigWithContext(tenantId *string, thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.DeleteThirdPartyConfigResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.DeleteThirdPartyConfigResponse{}, err
	}
	return (*instance.RecipeImpl.DeleteThirdPartyConfig)(tenantId, thirdPartyId, userContext)
}

func ListThirdPartyConfigsForThirdPartyIdWithContext(thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.ListThirdPartyConfigsForThirdPartyIdResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.ListThirdPartyConfigsForThirdPartyIdResponse{}, err
	}
	return (*instance.RecipeImpl.ListThirdPartyConfigsForThirdPartyId)(thirdPartyId, userContext)
}

func CreateOrUpdateTenant(tenantId *string, config multitenancymodels.TenantConfig) (multitenancymodels.CreateOrUpdateTenantResponse, error) {
	return CreateOrUpdateTenantWithContext(tenantId, config, &map[string]interface{}{})
}

func DeleteTenant(tenantId string) (multitenancymodels.DeleteTenantResponse, error) {
	return DeleteTenantWithContext(tenantId, &map[string]interface{}{})
}

func GetTenantConfig(tenantId *string) (multitenancymodels.TenantConfigResponse, error) {
	return GetTenantConfigWithContext(tenantId, &map[string]interface{}{})
}

func ListAllTenants() (multitenancymodels.ListAllTenantsResponse, error) {
	return ListAllTenantsWithContext(&map[string]interface{}{})
}

// Third party provider management
func CreateOrUpdateThirdPartyConfig(config tpmodels.ProviderConfig, skipValidation bool) (multitenancymodels.CreateOrUpdateThirdPartyConfigResponse, error) {
	return CreateOrUpdateThirdPartyConfigWithContext(config, skipValidation, &map[string]interface{}{})
}

func DeleteThirdPartyConfig(tenantId *string, thirdPartyId string) (multitenancymodels.DeleteThirdPartyConfigResponse, error) {
	return DeleteThirdPartyConfigWithContext(tenantId, thirdPartyId, &map[string]interface{}{})
}

func ListThirdPartyConfigsForThirdPartyId(thirdPartyId string) (multitenancymodels.ListThirdPartyConfigsForThirdPartyIdResponse, error) {
	return ListThirdPartyConfigsForThirdPartyIdWithContext(thirdPartyId, &map[string]interface{}{})
}
