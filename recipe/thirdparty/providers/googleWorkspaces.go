package providers

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ValidateAndNormaliseGoogleWorkspaces(input tpmodels.ProviderInput) (tpmodels.ProviderInput, error) {
	if input.Config.ThirdPartyId == "" {
		input.Config.ThirdPartyId = "google-workspaces"
	}
	if input.Config.Name == "" {
		input.Config.Name = "Google Workspaces"
	}

	if input.Config.ValidateIdTokenPayload == nil {
		input.Config.ValidateIdTokenPayload = func(idTokenPayload map[string]interface{}, clientConfig tpmodels.ProviderConfigForClientType) error {
			if clientConfig.AdditionalConfig != nil && clientConfig.AdditionalConfig["hd"] != nil && clientConfig.AdditionalConfig["hd"] != "*" && idTokenPayload["hd"] != clientConfig.AdditionalConfig["hd"] {
				return errors.New("the value for hd claim in the id token does not match the value provided in the config")
			}
			return nil
		}
	}

	// TODO add validation

	return ValidateAndNormaliseNewProvider(input)
}

func GoogleWorkspaces(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfigForClientType
		provider.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if config.AdditionalConfig == nil || config.AdditionalConfig["hd"] == nil || config.AdditionalConfig["hd"] == "" {
				config.AuthorizationEndpointQueryParams["hd"] = "*"
			} else {
				config.AuthorizationEndpointQueryParams["hd"] = config.AdditionalConfig["hd"]
			}

			return discoverOIDCEndpoints(config)
		}

		if oOverride != nil {
			provider = oOverride(provider)
		}
		return provider
	}

	return Google(input)
}
