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

package multitenancymodels

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type RecipeInterface struct {
	// Tenant management
	CreateOrUpdateTenant *func(tenantId *string, config TenantConfig, userContext supertokens.UserContext) (CreateOrUpdateTenantResponse, error)
	DeleteTenant         *func(tenantId string, userContext supertokens.UserContext) (DeleteTenantResponse, error)
	GetTenantConfig      *func(tenantId *string, userContext supertokens.UserContext) (TenantConfigResponse, error)
	ListAllTenants       *func(userContext supertokens.UserContext) (ListAllTenantsResponse, error)

	// Third party provider management
	CreateOrUpdateThirdPartyConfig *func(config tpmodels.ProviderConfig, userContext supertokens.UserContext) (CreateOrUpdateThirdPartyConfigResponse, error)
	DeleteThirdPartyConfig         *func(tenantId *string, thirdPartyId string, userContext supertokens.UserContext) (DeleteThirdPartyConfigResponse, error)
	ListThirdPartyConfigs          *func(thirdPartyId string, userContext supertokens.UserContext) (ListTenantConfigMappingsResponse, error)
}

type TenantConfig struct {
	EmailpasswordEnabled *bool
	PasswordlessEnabled  *bool
	ThirdpartyEnabled    *bool
}

type CreateOrUpdateTenantResponse struct {
	OK *struct {
		CreatedNew bool
	}
}

type DeleteTenantResponse struct {
	OK *struct {
		TenantExisted bool
	}
}

type TenantConfigResponse struct {
	OK *struct {
		Emailpassword struct {
			Enabled bool
		}
		Passwordless struct {
			Enabled bool
		}
		ThirdParty struct {
			Enabled   bool
			Providers []tpmodels.ProviderConfig
		}
	}
	TenantDoesNotExistError *struct{}
}

type ListAllTenantsResponse struct {
	OK *struct {
		Tenants []string
	}
}

type CreateOrUpdateThirdPartyConfigResponse struct {
	OK *struct {
		CreatedNew bool
	}
	TenantDoesNotExistError *struct{}
}

type DeleteThirdPartyConfigResponse struct {
	OK *struct {
		DidConfigExist bool
	}
	TenantDoesNotExistError *struct{}
}

type ListTenantConfigMappingsResponse struct {
	OK *struct {
		Providers []tpmodels.ProviderConfig
	}
}
