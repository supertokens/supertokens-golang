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
	createOrUpdateThirdPartyConfig := func(tenantId *string, thirdPartyId string, config tpmodels.ProviderConfig, userContext supertokens.UserContext) (multitenancymodels.CreateOrUpdateTenantIdConfigResponse, error) {
		// TODO impl
		return multitenancymodels.CreateOrUpdateTenantIdConfigResponse{}, nil
	}

	fetchThirdPartyConfig := func(tenantId *string, thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.FetchTenantIdConfigResponse, error) {
		// TODO impl
		return multitenancymodels.FetchTenantIdConfigResponse{}, nil
	}

	deleteThirdPartyConfig := func(tenantId *string, thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.DeleteTenantIdConfigResponse, error) {
		// TODO impl
		return multitenancymodels.DeleteTenantIdConfigResponse{}, nil
	}

	listThirdPartyConfigs := func(tenantId *string, userContext supertokens.UserContext) (multitenancymodels.ListTenantConfigMappingsResponse, error) {
		// TODO impl
		return multitenancymodels.ListTenantConfigMappingsResponse{}, nil
	}

	return multitenancymodels.RecipeInterface{
		CreateOrUpdateThirdPartyConfig: &createOrUpdateThirdPartyConfig,
		FetchThirdPartyConfig:          &fetchThirdPartyConfig,
		DeleteThirdPartyConfig:         &deleteThirdPartyConfig,
		ListThirdPartyConfigs:          &listThirdPartyConfigs,
	}
}
