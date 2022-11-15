package providers

import (
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const oktaID = "okta"

func Okta(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	if input.ThirdPartyID == "" {
		input.ThirdPartyID = oktaID
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "sub"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified = "email_verified"
	}

	if input.Config.UserInfoMap.FromIdTokenPayload.UserId == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.UserId = "sub"
	}
	if input.Config.UserInfoMap.FromIdTokenPayload.Email == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.Email = "email"
	}

	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfigForClientType
		provider.GetConfigForClientType = func(clientType *string, input tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, input, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if config.OIDCDiscoveryEndpoint == "" {
				config.OIDCDiscoveryEndpoint = fmt.Sprintf("https://%s.okta.com", config.AdditionalConfig["oktaDomain"])
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"openid", "email"}
			}

			if config.ClientSecret == "" && config.AdditionalConfig["certificate"] != nil {
				if config.TokenParams == nil {
					config.TokenParams = map[string]interface{}{}
				}
				config.TokenParams["client_assertion_type"] = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
				ca, err := getADClientAssertion(config)
				if err != nil {
					return tpmodels.ProviderConfigForClientType{}, err
				}
				config.TokenParams["client_assertion"] = ca
			}

			return config, err
		}

		if oOverride != nil {
			provider = oOverride(provider)
		}
		return provider
	}

	return NewProvider(input)
}
