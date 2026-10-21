package api

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type tenantsListResponse struct {
	Status  string       `json:"status"`
	Tenants []tenantType `json:"tenants"`
}

type tenantType struct {
	TenantId      string `json:"tenantId"`
	EmailPassword struct {
		Enabled bool `json:"enabled"`
	} `json:"emailPassword"`
	Passwordless struct {
		Enabled bool `json:"enabled"`
	} `json:"passwordless"`
	ThirdParty struct {
		Enabled bool `json:"enabled"`
	} `json:"thirdParty"`
}

func TenantsListGet(apiInterface dashboardmodels.APIInterface, tenantId string, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (tenantsListResponse, error) {
	tenantsResponse, err := multitenancy.ListAllTenants(userContext)
	if err != nil {
		return tenantsListResponse{}, err
	}
	result := tenantsListResponse{
		Status:  "OK",
		Tenants: make([]tenantType, len(tenantsResponse.OK.Tenants)),
	}

	for i, tenant := range tenantsResponse.OK.Tenants {
		result.Tenants[i] = tenantType{
			TenantId: tenant.TenantId,
			EmailPassword: struct {
				Enabled bool `json:"enabled"`
			}{
				Enabled: tenant.EmailPassword.Enabled,
			},
			Passwordless: struct {
				Enabled bool `json:"enabled"`
			}{
				Enabled: tenant.Passwordless.Enabled,
			},
			ThirdParty: struct {
				Enabled bool `json:"enabled"`
			}{
				Enabled: tenant.ThirdParty.Enabled,
			},
		}
	}

	return result, nil
}
