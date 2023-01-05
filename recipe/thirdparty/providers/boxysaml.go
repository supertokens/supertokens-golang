package providers

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ValidateAndNormaliseBoxySaml(input tpmodels.ProviderInput) (tpmodels.ProviderInput, error) {
	if input.Config.Name == "" {
		input.Config.Name = "SAML"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "id"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}

	// TODO add validation

	return ValidateAndNormaliseNewProvider(input)
}

func BoxySaml(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfigForClientType
		provider.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			boxyURL, ok := config.AdditionalConfig["boxyURL"].(string)
			if !ok {
				return tpmodels.ProviderConfigForClientType{}, errors.New("boxyURL is missing or an invalid value in the additionalConfig")
			}

			if config.AuthorizationEndpoint == "" {
				config.AuthorizationEndpoint = boxyURL + "/api/oauth/authorize"
			}

			if config.TokenEndpoint == "" {
				config.TokenEndpoint = boxyURL + "/api/oauth/token"
			}

			if config.UserInfoEndpoint == "" {
				config.UserInfoEndpoint = boxyURL + "/api/oauth/userinfo"
			}

			return config, nil
		}

		if oOverride != nil {
			provider = oOverride(provider)
		}
		return provider
	}

	return NewProvider(input)
}
