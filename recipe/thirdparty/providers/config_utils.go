package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

func getProviderConfigForClient(config tpmodels.ProviderConfig, clientConfig tpmodels.ProviderClientConfig) tpmodels.ProviderConfigForClientType {
	return tpmodels.ProviderConfigForClientType{
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
		TenantId:                         config.TenantId,
	}
}
