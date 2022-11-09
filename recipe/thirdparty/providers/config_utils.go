package providers

import (
	"errors"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeToCustomProvider interface {
	ToCustomProviderClientConfig() (CustomClientConfig, error)
}

type TypeFromCustomProvider interface {
	UpdateFromCustomProviderClientConfig(config CustomClientConfig)
}

func findConfig(out TypeFromCustomProvider, clientType *string, tenantId *string, userContext supertokens.UserContext, clients []TypeToCustomProvider) error {
	if clientType == nil {
		if len(clients) == 0 || len(clients) > 1 {
			return errors.New("please provide exactly one client config or pass clientType or tenantId")
		}

		cConfig, err := clients[0].ToCustomProviderClientConfig()
		if err != nil {
			return err
		}
		out.UpdateFromCustomProviderClientConfig(cConfig)
		return nil
	}

	// (else) clientType is not nil
	if tenantId == nil {
		for _, config := range clients {
			cConfig, err := config.ToCustomProviderClientConfig()
			if err != nil {
				return err
			}
			if cConfig.ClientType == *clientType {
				out.UpdateFromCustomProviderClientConfig(cConfig)
				return nil
			}
		}

		return errors.New("config for specified clientType not found")
	} else {
		// TODO Multitenant
		return errors.New("needs implementation")
	}
}

func getCombinedOAuth2Config(config CustomConfig, clientConfig CustomClientConfig) *typeCombinedOAuth2Config {
	combinedConfig := &typeCombinedOAuth2Config{
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
	}
	combinedConfig.ValidateIdTokenPayload = func(idTokenPayload map[string]interface{}, config *typeCombinedOAuth2Config) (bool, error) {
		if config.ValidateIdTokenPayload != nil {
			return config.ValidateIdTokenPayload(idTokenPayload, combinedConfig)
		}
		return true, nil
	}
	return combinedConfig
}
