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
	CreateOrUpdateThirdPartyConfig *func(tenantId *string, thirdPartyId string, config tpmodels.ProviderConfig, userContext supertokens.UserContext) (CreateOrUpdateTenantIdConfigResponse, error)
	FetchThirdPartyConfig          *func(tenantId *string, thirdPartyId string, userContext supertokens.UserContext) (FetchTenantIdConfigResponse, error)
	DeleteThirdPartyConfig         *func(tenantId *string, thirdPartyId string, userContext supertokens.UserContext) (DeleteTenantIdConfigResponse, error)
	ListThirdPartyConfigs          *func(tenantId *string, userContext supertokens.UserContext) (ListTenantConfigMappingsResponse, error)
}

type CreateOrUpdateTenantIdConfigResponse struct {
	OK *struct {
		CreatedNew bool
	}
}

type FetchTenantIdConfigResponse struct {
	OK *struct {
		Config tpmodels.ProviderConfig
	}
	UnknownMappingError *struct{}
}

type DeleteTenantIdConfigResponse struct {
	OK *struct {
		DidMappingExist bool
	}
}

type ListTenantConfigMappingsResponse struct {
	OK *struct {
		Configs []struct {
			ThirdPartyId string
			Config       tpmodels.ProviderConfig
		}
	}
}
