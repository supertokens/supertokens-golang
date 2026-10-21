package providers

import (
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func getProviderConfigForClient(config tpmodels.ProviderConfig, clientConfig tpmodels.ProviderClientConfig) tpmodels.ProviderConfigForClientType {
	return tpmodels.ProviderConfigForClientType{
		ClientID:         clientConfig.ClientID,
		ClientSecret:     clientConfig.ClientSecret,
		Scope:            clientConfig.Scope,
		ForcePKCE:        clientConfig.ForcePKCE,
		AdditionalConfig: copyConfigMap(clientConfig.AdditionalConfig),

		Name: config.Name,

		AuthorizationEndpoint:            config.AuthorizationEndpoint,
		AuthorizationEndpointQueryParams: copyConfigMap(config.AuthorizationEndpointQueryParams),
		TokenEndpoint:                    config.TokenEndpoint,
		TokenEndpointBodyParams:          copyConfigMap(config.TokenEndpointBodyParams),
		UserInfoEndpoint:                 config.UserInfoEndpoint,
		UserInfoEndpointQueryParams:      copyConfigMap(config.UserInfoEndpointQueryParams),
		UserInfoEndpointHeaders:          copyConfigMap(config.UserInfoEndpointHeaders),
		JwksURI:                          config.JwksURI,
		OIDCDiscoveryEndpoint:            config.OIDCDiscoveryEndpoint,
		UserInfoMap:                      config.UserInfoMap,
		ValidateIdTokenPayload:           config.ValidateIdTokenPayload,
		ValidateAccessToken:              config.ValidateAccessToken,
		RequireEmail:                     config.RequireEmail,
		GenerateFakeEmail:                config.GenerateFakeEmail,
	}
}

// copyConfigMap returns a deep copy of a config map (nil stays nil), so that the config used for a
// single request never shares a map with the static config or with other requests.
func copyConfigMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[k] = copyConfigValue(v)
	}
	return result
}

func copyConfigValue(v interface{}) interface{} {
	switch typed := v.(type) {
	case map[string]interface{}:
		return copyConfigMap(typed)
	case []interface{}:
		copied := make([]interface{}, len(typed))
		for i, item := range typed {
			copied[i] = copyConfigValue(item)
		}
		return copied
	default:
		return v
	}
}

// copyConfigMapForUpdate is like copyConfigMap, but returns an empty map instead of nil so that the
// result can be written to.
func copyConfigMapForUpdate(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return map[string]interface{}{}
	}
	return copyConfigMap(m)
}

func FindAndCreateProviderInstance(providers []tpmodels.ProviderInput, thirdPartyId string, clientType *string, userContext supertokens.UserContext) (*tpmodels.TypeProvider, error) {
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

	return nil, nil
}

func createProvider(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if strings.HasPrefix(input.Config.ThirdPartyId, "active-directory") {
		return ActiveDirectory(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "apple") {
		return Apple(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "bitbucket") {
		return Bitbucket(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "discord") {
		return Discord(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "facebook") {
		return Facebook(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "github") {
		return Github(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "gitlab") {
		return Gitlab(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "google-workspaces") {
		return GoogleWorkspaces(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "google") {
		return Google(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "okta") {
		return Okta(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "linkedin") {
		return Linkedin(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "boxy-saml") {
		return BoxySaml(input)
	} else if strings.HasPrefix(input.Config.ThirdPartyId, "twitter") {
		return Twitter(input)
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

func MergeProvidersFromCoreAndStatic(providerConfigsFromCore []tpmodels.ProviderConfig, providerInputsFromStatic []tpmodels.ProviderInput, includeAllProviders bool) []tpmodels.ProviderInput {
	mergedProviders := []tpmodels.ProviderInput{}

	if len(providerConfigsFromCore) == 0 {
		for _, config := range providerInputsFromStatic {
			if includeAllProviders || config.IncludeInNonPublicTenantsByDefault {
				mergedProviders = append(mergedProviders, config)
			}
		}
	} else {
		for _, providerConfigFromCore := range providerConfigsFromCore {
			mergedProviderInput := tpmodels.ProviderInput{
				Config: providerConfigFromCore,
			}

			for _, providerInputFromStatic := range providerInputsFromStatic {
				if providerInputFromStatic.Config.ThirdPartyId == providerConfigFromCore.ThirdPartyId {
					mergedProviderInput.Config = mergeConfig(providerInputFromStatic.Config, providerConfigFromCore)
					mergedProviderInput.Override = providerInputFromStatic.Override
					break
				}
			}

			mergedProviders = append(mergedProviders, mergedProviderInput)
		}
	}

	return mergedProviders
}

func mergeConfig(staticConfig tpmodels.ProviderConfig, coreConfig tpmodels.ProviderConfig) tpmodels.ProviderConfig {
	result := staticConfig

	if coreConfig.AuthorizationEndpoint != "" {
		result.AuthorizationEndpoint = coreConfig.AuthorizationEndpoint
	}
	if coreConfig.AuthorizationEndpointQueryParams != nil {
		result.AuthorizationEndpointQueryParams = coreConfig.AuthorizationEndpointQueryParams
	}
	if coreConfig.TokenEndpoint != "" {
		result.TokenEndpoint = coreConfig.TokenEndpoint
	}
	if coreConfig.TokenEndpointBodyParams != nil {
		result.TokenEndpointBodyParams = coreConfig.TokenEndpointBodyParams
	}
	if coreConfig.UserInfoEndpoint != "" {
		result.UserInfoEndpoint = coreConfig.UserInfoEndpoint
	}
	if coreConfig.UserInfoEndpointHeaders != nil {
		result.UserInfoEndpointHeaders = coreConfig.UserInfoEndpointHeaders
	}
	if coreConfig.UserInfoEndpointQueryParams != nil {
		result.UserInfoEndpointQueryParams = coreConfig.UserInfoEndpointQueryParams
	}
	if coreConfig.JwksURI != "" {
		result.JwksURI = coreConfig.JwksURI
	}
	if coreConfig.Name != "" {
		result.Name = coreConfig.Name
	}
	if coreConfig.OIDCDiscoveryEndpoint != "" {
		result.OIDCDiscoveryEndpoint = coreConfig.OIDCDiscoveryEndpoint
	}
	if coreConfig.UserInfoMap.FromIdTokenPayload.Email != "" {
		result.UserInfoMap.FromIdTokenPayload.Email = coreConfig.UserInfoMap.FromIdTokenPayload.Email
	}
	if coreConfig.UserInfoMap.FromIdTokenPayload.EmailVerified != "" {
		result.UserInfoMap.FromIdTokenPayload.EmailVerified = coreConfig.UserInfoMap.FromIdTokenPayload.EmailVerified
	}
	if coreConfig.UserInfoMap.FromIdTokenPayload.UserId != "" {
		result.UserInfoMap.FromIdTokenPayload.UserId = coreConfig.UserInfoMap.FromIdTokenPayload.UserId
	}
	if coreConfig.UserInfoMap.FromUserInfoAPI.Email != "" {
		result.UserInfoMap.FromUserInfoAPI.Email = coreConfig.UserInfoMap.FromUserInfoAPI.Email
	}
	if coreConfig.UserInfoMap.FromUserInfoAPI.EmailVerified != "" {
		result.UserInfoMap.FromUserInfoAPI.EmailVerified = coreConfig.UserInfoMap.FromUserInfoAPI.EmailVerified
	}
	if coreConfig.UserInfoMap.FromUserInfoAPI.UserId != "" {
		result.UserInfoMap.FromUserInfoAPI.UserId = coreConfig.UserInfoMap.FromUserInfoAPI.UserId
	}

	// Merge the clients
	mergedClients := append([]tpmodels.ProviderClientConfig{}, staticConfig.Clients...)
	for _, client := range coreConfig.Clients {
		for i, staticClient := range mergedClients {
			if staticClient.ClientType == client.ClientType {
				mergedClients[i] = client
				break
			}
		}
	}
	result.Clients = mergedClients

	return result
}
