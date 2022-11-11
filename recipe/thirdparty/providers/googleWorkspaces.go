package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GoogleWorkspaces(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	input.ThirdPartyID = "google-workspaces"

	if input.Config.ValidateIdTokenPayload == nil {
		input.Config.ValidateIdTokenPayload = func(idTokenPayload map[string]interface{}, clientConfig tpmodels.ProviderConfigForClient) (bool, error) {
			return idTokenPayload["hd"] == clientConfig.AdditionalConfig["domain"], nil
		}
	}

	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfig
		provider.GetConfig = func(clientType *string, input tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClient, error) {
			config, err := oGetConfig(clientType, input, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClient{}, err
			}

			if config.AdditionalConfig == nil || config.AdditionalConfig["domain"] == nil || config.AdditionalConfig["domain"] == "" {
				config.AuthorizationEndpointQueryParams["hd"] = "*"
			} else {
				config.AuthorizationEndpointQueryParams["hd"] = config.AdditionalConfig["domain"]
			}

			return config, err
		}

		if oOverride != nil {
			provider = oOverride(provider)
		}
		return provider
	}

	return Google(input)
}
