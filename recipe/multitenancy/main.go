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

func CreateOrUpdateTenant(tenantId string, config multitenancymodels.TenantConfig, userContext ...supertokens.UserContext) (multitenancymodels.CreateOrUpdateTenantResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.CreateOrUpdateTenantResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.CreateOrUpdateTenant)(tenantId, config, userContext[0])
}

func DeleteTenant(tenantId string, userContext ...supertokens.UserContext) (multitenancymodels.DeleteTenantResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.DeleteTenantResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.DeleteTenant)(tenantId, userContext[0])
}

func GetTenant(tenantId string, userContext ...supertokens.UserContext) (*multitenancymodels.Tenant, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetTenant)(tenantId, userContext[0])
}

func ListAllTenants(userContext ...supertokens.UserContext) (multitenancymodels.ListAllTenantsResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.ListAllTenantsResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ListAllTenants)(userContext[0])
}

// Third party provider management
func CreateOrUpdateThirdPartyConfig(tenantId string, config tpmodels.ProviderConfig, skipValidation *bool, userContext ...supertokens.UserContext) (multitenancymodels.CreateOrUpdateThirdPartyConfigResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.CreateOrUpdateThirdPartyConfigResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.CreateOrUpdateThirdPartyConfig)(tenantId, config, skipValidation, userContext[0])
}

func DeleteThirdPartyConfig(tenantId string, thirdPartyId string, userContext ...supertokens.UserContext) (multitenancymodels.DeleteThirdPartyConfigResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.DeleteThirdPartyConfigResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.DeleteThirdPartyConfig)(tenantId, thirdPartyId, userContext[0])
}

func AssociateUserToTenant(tenantId string, userId string, userContext ...supertokens.UserContext) (multitenancymodels.AssociateUserToTenantResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.AssociateUserToTenantResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.AssociateUserToTenant)(tenantId, userId, userContext[0])
}

func DisassociateUserFromTenant(tenantId string, userId string, userContext ...supertokens.UserContext) (multitenancymodels.DisassociateUserFromTenantResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return multitenancymodels.DisassociateUserFromTenantResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.DisassociateUserFromTenant)(tenantId, userId, userContext[0])
}
