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
	GetTenantId *func(tenantIdFromFrontend string, userContext supertokens.UserContext) (string, error)

	// Tenant management
	CreateOrUpdateTenant *func(tenantId string, config TenantConfig, userContext supertokens.UserContext) (CreateOrUpdateTenantResponse, error)
	DeleteTenant         *func(tenantId string, userContext supertokens.UserContext) (DeleteTenantResponse, error)
	GetTenant            *func(tenantId string, userContext supertokens.UserContext) (*Tenant, error)
	ListAllTenants       *func(userContext supertokens.UserContext) (ListAllTenantsResponse, error)

	// Third party provider management
	CreateOrUpdateThirdPartyConfig *func(tenantId string, config tpmodels.ProviderConfig, skipValidation bool, userContext supertokens.UserContext) (CreateOrUpdateThirdPartyConfigResponse, error)
	DeleteThirdPartyConfig         *func(tenantId string, thirdPartyId string, userContext supertokens.UserContext) (DeleteThirdPartyConfigResponse, error)

	// User tenant association
	AssociateUserToTenant      *func(tenantId string, userId string, userContext supertokens.UserContext) (AssociateUserToTenantResponse, error)
	DisassociateUserFromTenant *func(tenantId string, userId string, userContext supertokens.UserContext) (DisassociateUserFromTenantResponse, error)
}

type TenantConfig struct {
	EmailPasswordEnabled *bool
	PasswordlessEnabled  *bool
	ThirdPartyEnabled    *bool
	CoreConfig           map[string]interface{}
}

type CreateOrUpdateTenantResponse struct {
	OK *struct {
		CreatedNew bool
	}
}

type DeleteTenantResponse struct {
	OK *struct {
		DidExist bool
	}
}

type Tenant struct {
	EmailPassword struct {
		Enabled bool `json:"enabled"`
	} `json:"emailPassword"`
	Passwordless struct {
		Enabled bool `json:"enabled"`
	} `json:"passwordless"`
	ThirdParty struct {
		Enabled   bool                      `json:"enabled"`
		Providers []tpmodels.ProviderConfig `json:"providers"`
	} `json:"thirdParty"`
	CoreConfig map[string]interface{} `json:"coreConfig"`
}

type ListAllTenantsResponse struct {
	OK *struct {
		Tenants []Tenant `json:"tenants"`
	}
}

type CreateOrUpdateThirdPartyConfigResponse struct {
	OK *struct {
		CreatedNew bool
	}
}

type DeleteThirdPartyConfigResponse struct {
	OK *struct {
		DidConfigExist bool
	}
}

type AssociateUserToTenantResponse struct {
	OK *struct {
		WasAlreadyAssociated bool
	}
	UnknownUserIdError               *struct{}
	EmailAlreadyExistsError          *struct{}
	PhoneNumberAlreadyExistsError    *struct{}
	ThirdPartyUserAlreadyExistsError *struct{}
}

type DisassociateUserFromTenantResponse struct {
	OK *struct {
		WasAssociated bool
	}
}
