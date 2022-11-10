package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

func getCombinedProviderConfig(config tpmodels.ProviderConfigInput, clientConfig tpmodels.ProviderClientConfigInput) tpmodels.ProviderClientConfig {
	return tpmodels.ProviderClientConfig{
		ClientType:       clientConfig.ClientType,
		ClientID:         clientConfig.ClientID,
		ClientSecret:     clientConfig.ClientSecret,
		Scope:            clientConfig.Scope,
		AdditionalConfig: clientConfig.AdditionalConfig,

		AuthorizationEndpoint:            config.AuthorizationEndpoint,
		AuthorizationEndpointQueryParams: config.AuthorizationEndpointQueryParams,
		TokenEndpoint:                    config.TokenEndpoint,
		TokenParams:                      config.TokenParams,
		ForcePKCE:                        config.ForcePKCE,
		UserInfoEndpoint:                 config.UserInfoEndpoint,
		JwksURI:                          config.JwksURI,
		OIDCDiscoveryEndpoint:            config.OIDCDiscoveryEndpoint,
		UserInfoMap:                      config.UserInfoMap,
		ValidateIdTokenPayload:           config.ValidateIdTokenPayload,
	}
}
