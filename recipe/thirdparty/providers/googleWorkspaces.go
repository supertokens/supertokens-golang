package providers

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GoogleWorkspaces(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.Name == "" {
		input.Config.Name = "Google Workspaces"
	}

	if input.Config.ValidateIdTokenPayload == nil {
		input.Config.ValidateIdTokenPayload = func(idTokenPayload map[string]interface{}, clientConfig tpmodels.ProviderConfigForClientType, userContext supertokens.UserContext) error {
			if clientConfig.AdditionalConfig != nil && clientConfig.AdditionalConfig["hd"] != nil && clientConfig.AdditionalConfig["hd"] != "*" && idTokenPayload["hd"] != clientConfig.AdditionalConfig["hd"] {
				return errors.New("the value for hd claim in the id token does not match the value provided in the config")
			}
			return nil
		}
	}

	oOverride := input.Override

	input.Override = func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := originalImplementation.GetConfigForClientType
		originalImplementation.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if config.AdditionalConfig == nil || config.AdditionalConfig["hd"] == nil || config.AdditionalConfig["hd"] == "" {
				config.AuthorizationEndpointQueryParams["hd"] = "*"
			} else {
				config.AuthorizationEndpointQueryParams["hd"] = config.AdditionalConfig["hd"]
			}

			return config, nil
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}
		return originalImplementation
	}

	return Google(input)
}
