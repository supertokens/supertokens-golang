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
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier, config multitenancymodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) multitenancymodels.RecipeInterface {
	getTenantId := func(tenantIdFromFrontend *string, userContext supertokens.UserContext) (*string, error) {
		return tenantIdFromFrontend, nil
	}

	createOrUpdateTenant := func(tenantId *string, config multitenancymodels.TenantConfig, userContext supertokens.UserContext) (multitenancymodels.CreateOrUpdateTenantResponse, error) {
		// TODO impl
		return multitenancymodels.CreateOrUpdateTenantResponse{}, errors.New("not implemented")
	}

	deleteTenant := func(tenantId string, userContext supertokens.UserContext) (multitenancymodels.DeleteTenantResponse, error) {
		// TODO impl
		return multitenancymodels.DeleteTenantResponse{}, errors.New("not implemented")
	}

	getTenantConfig := func(tenantId *string, userContext supertokens.UserContext) (multitenancymodels.TenantConfigResponse, error) {
		// TODO impl
		// TODO may throw mterrors.TenantDoesNotExistError
		return multitenancymodels.TenantConfigResponse{
			OK: &struct {
				EmailPassword struct{ Enabled bool }
				Passwordless  struct{ Enabled bool }
				ThirdParty    struct {
					Enabled   bool
					Providers []tpmodels.ProviderConfig
				}
			}{
				EmailPassword: struct{ Enabled bool }{
					Enabled: true,
				},
				Passwordless: struct{ Enabled bool }{
					Enabled: true,
				},
				ThirdParty: struct {
					Enabled   bool
					Providers []tpmodels.ProviderConfig
				}{
					Enabled:   true,
					Providers: []tpmodels.ProviderConfig{},
				},
			},
		}, nil
	}

	listAllTenants := func(userContext supertokens.UserContext) (multitenancymodels.ListAllTenantsResponse, error) {
		// TODO impl
		return multitenancymodels.ListAllTenantsResponse{}, errors.New("not implemented")
	}

	createOrUpdateThirdPartyConfig := func(config tpmodels.ProviderConfig, skipValidation bool, userContext supertokens.UserContext) (multitenancymodels.CreateOrUpdateThirdPartyConfigResponse, error) {
		// TODO impl
		// TODO may throw mterrors.TenantDoesNotExistError
		return multitenancymodels.CreateOrUpdateThirdPartyConfigResponse{}, errors.New("not implemented")
	}

	deleteThirdPartyConfig := func(tenantId *string, thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.DeleteThirdPartyConfigResponse, error) {
		// TODO impl
		// TODO may throw mterrors.TenantDoesNotExistError
		return multitenancymodels.DeleteThirdPartyConfigResponse{}, errors.New("not implemented")
	}

	listThirdPartyConfigs := func(thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.ListThirdPartyConfigsForThirdPartyIdResponse, error) {
		// TODO impl
		return multitenancymodels.ListThirdPartyConfigsForThirdPartyIdResponse{}, errors.New("not implemented")
	}

	return multitenancymodels.RecipeInterface{
		GetTenantId: &getTenantId,

		CreateOrUpdateTenant: &createOrUpdateTenant,
		DeleteTenant:         &deleteTenant,
		GetTenantConfig:      &getTenantConfig,
		ListAllTenants:       &listAllTenants,

		CreateOrUpdateThirdPartyConfig:       &createOrUpdateThirdPartyConfig,
		DeleteThirdPartyConfig:               &deleteThirdPartyConfig,
		ListThirdPartyConfigsForThirdPartyId: &listThirdPartyConfigs,
	}
}
