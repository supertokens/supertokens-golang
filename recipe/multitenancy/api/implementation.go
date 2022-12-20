package api

import (
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() multitenancymodels.APIInterface {

	loginMethodsAPI := func(tenantId *string, options multitenancymodels.APIOptions, userContext supertokens.UserContext) (multitenancymodels.LoginMethodsGETResponse, error) {
		tenantConfigResponse, err := (*options.RecipeImplementation.GetTenantConfig)(tenantId, userContext)
		if err != nil {
			return multitenancymodels.LoginMethodsGETResponse{}, err
		}

		if tenantConfigResponse.OK != nil {

			staticProviders := []multitenancymodels.TypeThirdPartyProvider{}

			for _, provider := range options.StaticThirdPartyProviders {
				if tenantId == nil || *tenantId == tpmodels.DefaultTenantId {
					if provider.UseForDefaultTenant != nil && !*provider.UseForDefaultTenant {
						continue
					}
				}
				staticProviders = append(staticProviders, multitenancymodels.TypeThirdPartyProvider{
					Id:   provider.Config.ThirdPartyId,
					Name: provider.Config.Name,
				})
			}

			providersFromCore := []multitenancymodels.TypeThirdPartyProvider{}

			for _, provider := range tenantConfigResponse.OK.ThirdParty.Providers {
				providersFromCore = append(providersFromCore, multitenancymodels.TypeThirdPartyProvider{
					Id:   provider.ThirdPartyId,
					Name: provider.Name,
				})
			}

			var finalProviderList []multitenancymodels.TypeThirdPartyProvider

			if len(providersFromCore) > 0 {
				finalProviderList = providersFromCore
			} else {
				finalProviderList = staticProviders
			}

			result := multitenancymodels.LoginMethodsGETResponse{
				OK: &multitenancymodels.TypeLoginMethods{
					Emailpassword: multitenancymodels.TypeEmailpassword{
						Enabled: tenantConfigResponse.OK.Emailpassword.Enabled,
					},
					Passwordless: multitenancymodels.TypePasswordless{
						Enabled: tenantConfigResponse.OK.Passwordless.Enabled,
					},
					ThirdParty: multitenancymodels.TypeThirdParty{
						Enabled:   tenantConfigResponse.OK.ThirdParty.Enabled,
						Providers: finalProviderList,
					},
				},
			}
			return result, nil
		} else if tenantConfigResponse.TenantDoesNotExistError != nil {
			return multitenancymodels.LoginMethodsGETResponse{
				TenantDoesNotExistError: &struct{}{},
			}, nil
		}

		return multitenancymodels.LoginMethodsGETResponse{}, nil
	}

	return multitenancymodels.APIInterface{
		LoginMethodsGET: &loginMethodsAPI,
	}
}
