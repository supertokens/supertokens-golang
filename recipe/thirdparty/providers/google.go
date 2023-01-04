package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ValidateAndNormaliseGoogle(input tpmodels.ProviderInput) (tpmodels.ProviderInput, error) {
	if input.Config.Name == "" {
		input.Config.Name = "Google"
	}

	if input.Config.OIDCDiscoveryEndpoint == "" {
		input.Config.OIDCDiscoveryEndpoint = "https://accounts.google.com/"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "id"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified = "email_verified"
	}

	if input.Config.AuthorizationEndpointQueryParams == nil {
		input.Config.AuthorizationEndpointQueryParams = map[string]interface{}{}
	}

	if input.Config.AuthorizationEndpointQueryParams["include_granted_scopes"] == nil {
		input.Config.AuthorizationEndpointQueryParams["include_granted_scopes"] = "true"
	}
	if input.Config.AuthorizationEndpointQueryParams["access_type"] == nil {
		input.Config.AuthorizationEndpointQueryParams["access_type"] = "offline"
	}

	// TODO add validation

	return ValidateAndNormaliseNewProvider(input)
}

func Google(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfigForClientType
		provider.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"openid", "email"}
			}

			return discoverOIDCEndpoints(config)
		}

		if oOverride != nil {
			provider = oOverride(provider)
		}
		return provider
	}

	return NewProvider(input)
}
