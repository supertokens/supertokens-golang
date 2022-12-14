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

func makeRecipeImplementation(querier supertokens.Querier, config multitenancymodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) multitenancymodels.RecipeInterface {

	createOrUpdateTenant := func(tenantId *string, config multitenancymodels.TenantConfig, userContext supertokens.UserContext) (multitenancymodels.CreateOrUpdateTenantResponse, error) {
		// TODO impl
		return multitenancymodels.CreateOrUpdateTenantResponse{}, nil
	}

	deleteTenant := func(tenantId string, userContext supertokens.UserContext) (multitenancymodels.DeleteTenantResponse, error) {
		// TODO impl
		return multitenancymodels.DeleteTenantResponse{}, nil
	}

	getTenantConfig := func(tenantId *string, userContext supertokens.UserContext) (multitenancymodels.TenantConfigResponse, error) {
		// TODO impl
		return multitenancymodels.TenantConfigResponse{}, nil
	}

	listAllTenants := func(userContext supertokens.UserContext) (multitenancymodels.ListAllTenantsResponse, error) {
		// TODO impl
		return multitenancymodels.ListAllTenantsResponse{}, nil
	}

	createOrUpdateThirdPartyConfig := func(config tpmodels.ProviderConfig, userContext supertokens.UserContext) (multitenancymodels.CreateOrUpdateThirdPartyConfigResponse, error) {
		// TODO impl
		return multitenancymodels.CreateOrUpdateThirdPartyConfigResponse{}, nil
	}

	deleteThirdPartyConfig := func(tenantId *string, thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.DeleteThirdPartyConfigResponse, error) {
		// TODO impl
		return multitenancymodels.DeleteThirdPartyConfigResponse{}, nil
	}

	listThirdPartyConfigs := func(thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.ListThirdPartyConfigsResponse, error) {
		// TODO impl
		return multitenancymodels.ListThirdPartyConfigsResponse{}, nil
	}

	return multitenancymodels.RecipeInterface{
		CreateOrUpdateTenant: &createOrUpdateTenant,
		DeleteTenant:         &deleteTenant,
		GetTenantConfig:      &getTenantConfig,
		ListAllTenants:       &listAllTenants,

		CreateOrUpdateThirdPartyConfig: &createOrUpdateThirdPartyConfig,
		DeleteThirdPartyConfig:         &deleteThirdPartyConfig,
		ListThirdPartyConfigs:          &listThirdPartyConfigs,
	}
}
