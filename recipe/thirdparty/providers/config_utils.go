package providers

import (
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func getProviderConfigForClient(config tpmodels.ProviderConfig, clientConfig tpmodels.ProviderClientConfig) tpmodels.ProviderConfigForClientType {
	return tpmodels.ProviderConfigForClientType{
		ClientID:         clientConfig.ClientID,
		ClientSecret:     clientConfig.ClientSecret,
		Scope:            clientConfig.Scope,
		ForcePKCE:        clientConfig.ForcePKCE,
		AdditionalConfig: clientConfig.AdditionalConfig,

		AuthorizationEndpoint:            config.AuthorizationEndpoint,
		AuthorizationEndpointQueryParams: config.AuthorizationEndpointQueryParams,
		TokenEndpoint:                    config.TokenEndpoint,
		TokenEndpointBodyParams:          config.TokenEndpointBodyParams,
		UserInfoEndpoint:                 config.UserInfoEndpoint,
		UserInfoEndpointQueryParams:      config.UserInfoEndpointQueryParams,
		UserInfoEndpointHeaders:          config.UserInfoEndpointHeaders,
		JwksURI:                          config.JwksURI,
		OIDCDiscoveryEndpoint:            config.OIDCDiscoveryEndpoint,
		UserInfoMap:                      config.UserInfoMap,
		ValidateIdTokenPayload:           config.ValidateIdTokenPayload,
		TenantId:                         config.TenantId,
	}
}

func FindAndCreateProviderInstance(providers []tpmodels.ProviderInput, thirdPartyId string, tenantId *string, clientType *string, userContext supertokens.UserContext) (*tpmodels.TypeProvider, error) {
	for _, provider := range providers {
		if provider.Config.ThirdPartyId == thirdPartyId {
			providerInstance := createProvider(provider)

			err := fetchAndSetConfig(providerInstance, clientType, userContext)
			if err != nil {
				return nil, err
			}

			return providerInstance, nil
		}
	}
	return nil, fmt.Errorf("the provider %s could not be found in the configuration", thirdPartyId)
}

func createProvider(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	switch input.Config.ThirdPartyId {
	case "active-directory":
		return ActiveDirectory(input)
	case "apple":
		return Apple(input)
	case "discord":
		return Discord(input)
	case "facebook":
		return Facebook(input)
	case "github":
		return Github(input)
	case "google":
		return Google(input)
	case "google-workspaces":
		return GoogleWorkspaces(input)
	case "okta":
		return Okta(input)
	case "linkedin":
		return Linkedin(input)
	case "boxy-saml":
		return BoxySaml(input)
	}

	return NewProvider(input)
}

func fetchAndSetConfig(provider *tpmodels.TypeProvider, clientType *string, userContext supertokens.UserContext) error {
	config, err := provider.GetConfigForClientType(clientType, userContext)
	if err != nil {
		return err
	}

	config, err = discoverOIDCEndpoints(config)
	if err != nil {
		return err
	}

	provider.Config = config
	return nil
}
