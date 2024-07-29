package api

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	tpproviders "github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tperrors"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() multitenancymodels.APIInterface {

	loginMethodsAPI := func(tenantId string, clientType *string, options multitenancymodels.APIOptions, userContext supertokens.UserContext) (multitenancymodels.LoginMethodsGETResponse, error) {

		tenantConfigResponse, err := (*options.RecipeImplementation.GetTenant)(tenantId, userContext)
		if err != nil {
			return multitenancymodels.LoginMethodsGETResponse{}, err
		}

		if tenantConfigResponse == nil {
			return multitenancymodels.LoginMethodsGETResponse{}, errors.New("tenant not found")
		}

		providerInputsFromStatic := options.StaticThirdPartyProviders
		providerConfigsFromCore := tenantConfigResponse.ThirdParty.Providers

		mergedProviders := tpproviders.MergeProvidersFromCoreAndStatic(providerConfigsFromCore, providerInputsFromStatic, tenantId == multitenancymodels.DefaultTenantId)

		var finalProviderList []multitenancymodels.TypeThirdPartyProvider = []multitenancymodels.TypeThirdPartyProvider{}

		for _, providerInput := range mergedProviders {
			providerInstance, err := tpproviders.FindAndCreateProviderInstance(mergedProviders, providerInput.Config.ThirdPartyId, clientType, userContext)
			if err != nil {
				if errors.As(err, &tperrors.ClientTypeNotFoundError{}) {
					continue // Skip as the clientType is missing for the particular provider
				}
				return multitenancymodels.LoginMethodsGETResponse{}, err
			}
			finalProviderList = append(finalProviderList, multitenancymodels.TypeThirdPartyProvider{
				Id:   providerInstance.ID,
				Name: providerInstance.Config.Name,
			})
		}

		result := multitenancymodels.LoginMethodsGETResponse{
			OK: &multitenancymodels.TypeLoginMethods{
				EmailPassword: multitenancymodels.TypeEmailPassword{
					Enabled: tenantConfigResponse.EmailPassword.Enabled,
				},
				Passwordless: multitenancymodels.TypePasswordless{
					Enabled: tenantConfigResponse.Passwordless.Enabled,
				},
				ThirdParty: multitenancymodels.TypeThirdParty{
					Enabled:   tenantConfigResponse.ThirdParty.Enabled,
					Providers: finalProviderList,
				},
			},
		}
		return result, nil
	}

	return multitenancymodels.APIInterface{
		LoginMethodsGET: &loginMethodsAPI,
	}
}
