package providers

import (
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const activeDirectoryID = "active-directory"

func ActiveDirectory(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	if input.ThirdPartyID == "" {
		input.ThirdPartyID = activeDirectoryID
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "sub"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}

	if input.Config.UserInfoMap.FromIdTokenPayload.UserId == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.UserId = "sub"
	}
	if input.Config.UserInfoMap.FromIdTokenPayload.Email == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.Email = "email"
	}

	if input.Config.AuthorizationEndpointQueryParams == nil {
		input.Config.AuthorizationEndpointQueryParams = map[string]interface{}{}
	}

	if input.Config.AuthorizationEndpointQueryParams["response_type"] == nil {
		input.Config.AuthorizationEndpointQueryParams["response_type"] = "code"
	}
	if input.Config.AuthorizationEndpointQueryParams["include_granted_scopes"] == nil {
		input.Config.AuthorizationEndpointQueryParams["include_granted_scopes"] = "true"
	}
	if input.Config.AuthorizationEndpointQueryParams["access_type"] == nil {
		input.Config.AuthorizationEndpointQueryParams["access_type"] = "offline"
	}

	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfig
		provider.GetConfig = func(clientType *string, input tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClient, error) {
			config, err := oGetConfig(clientType, input, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClient{}, err
			}

			if config.OIDCDiscoveryEndpoint == "" {
				config.OIDCDiscoveryEndpoint = fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0/", config.AdditionalConfig["directoryId"])
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"openid", "email"}
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
