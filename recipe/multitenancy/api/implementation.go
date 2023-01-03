package api

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() multitenancymodels.APIInterface {

	loginMethodsAPI := func(tenantId *string, options multitenancymodels.APIOptions, userContext supertokens.UserContext) (multitenancymodels.LoginMethodsGETResponse, error) {

		tenantId, err := (*options.RecipeImplementation.GetTenantId)(tenantId, userContext)
		if err != nil {
			return multitenancymodels.LoginMethodsGETResponse{}, err
		}

		tenantConfigResponse, err := (*options.RecipeImplementation.GetTenantConfig)(tenantId, userContext)
		if err != nil {
			return multitenancymodels.LoginMethodsGETResponse{}, err
		}

		if tenantConfigResponse.OK != nil {

			staticProviders := []multitenancymodels.TypeThirdPartyProvider{}

			for _, provider := range options.StaticThirdPartyProviders {
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

			/*
				With respect to https://supertokens.com/docs/contribute/decisions/multitenancy/0002,
				we are not merging the providers from core on top of static providers, because,
				we are just interested in the `Name` in the context of this API and core is expected
				to have the value for it.
			*/
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
		}

		return multitenancymodels.LoginMethodsGETResponse{}, errors.New("should never come here")
	}

	return multitenancymodels.APIInterface{
		LoginMethodsGET: &loginMethodsAPI,
	}
}
