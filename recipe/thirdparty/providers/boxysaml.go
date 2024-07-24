package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func BoxySaml(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.Name == "" {
		input.Config.Name = "SAML"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "id"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}

	oOverride := input.Override

	input.Override = func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := originalImplementation.GetConfigForClientType
		originalImplementation.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			boxyURL, ok := config.AdditionalConfig["boxyURL"].(string)
			if ok {
				if config.AuthorizationEndpoint == "" {
					config.AuthorizationEndpoint = boxyURL + "/api/oauth/authorize"
				}

				if config.TokenEndpoint == "" {
					config.TokenEndpoint = boxyURL + "/api/oauth/token"
				}

				if config.UserInfoEndpoint == "" {
					config.UserInfoEndpoint = boxyURL + "/api/oauth/userinfo"
				}
			}

			return config, nil
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}
		return originalImplementation
	}

	return NewProvider(input)
}
