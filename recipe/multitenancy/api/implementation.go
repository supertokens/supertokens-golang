package api

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	tpproviders "github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tperrors"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() multitenancymodels.APIInterface {

	loginMethodsAPI := func(tenantId *string, clientType *string, options multitenancymodels.APIOptions, userContext supertokens.UserContext) (multitenancymodels.LoginMethodsGETResponse, error) {

		tenantConfigResponse, err := (*options.RecipeImplementation.GetTenantConfig)(tenantId, userContext)
		if err != nil {
			return multitenancymodels.LoginMethodsGETResponse{}, err
		}

		providerInputsFromStatic := options.StaticThirdPartyProviders
		providerConfigsFromCore := tenantConfigResponse.OK.ThirdParty.Providers

		mergedProviders := tpproviders.MergeProvidersFromCoreAndStatic(tenantId, providerConfigsFromCore, providerInputsFromStatic)

		var finalProviderList []multitenancymodels.TypeThirdPartyProvider

		for _, providerInput := range mergedProviders {
			providerInstance, err := tpproviders.FindAndCreateProviderInstance(mergedProviders, providerInput.Config.ThirdPartyId, tenantId, clientType, userContext)
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

	return multitenancymodels.APIInterface{
		LoginMethodsGET: &loginMethodsAPI,
	}
}
