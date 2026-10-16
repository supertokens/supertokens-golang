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

		oGetUserInfo := originalImplementation.GetUserInfo
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			userInfo, err := oGetUserInfo(oAuthTokens, userContext)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
			err = checkGoogleWorkspacesHostedDomain(originalImplementation.Config, userInfo.RawUserInfoFromProvider)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
			return userInfo, nil
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}
		return originalImplementation
	}

	return Google(input)
}

// When a specific hosted domain is configured, the user has to come from a verified id_token (whose hd
// claim is checked by ValidateIdTokenPayload) or from a user info response with a matching hd claim.
// Only a missing hd and "*" allow any domain. Like in ValidateIdTokenPayload, an empty hd is compared
// as a value, so it only matches an empty hd claim.
func checkGoogleWorkspacesHostedDomain(config tpmodels.ProviderConfigForClientType, rawUserInfo tpmodels.TypeRawUserInfoFromProvider) error {
	hd := config.AdditionalConfig["hd"]
	if hd == nil || hd == "*" {
		return nil
	}
	if rawUserInfo.FromIdTokenPayload != nil {
		return nil
	}
	if rawUserInfo.FromUserInfoAPI != nil && rawUserInfo.FromUserInfoAPI["hd"] == hd {
		return nil
	}
	return errors.New("the value for hd claim in the user info does not match the value provided in the config")
}
