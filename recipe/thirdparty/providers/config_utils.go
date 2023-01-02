package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
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
