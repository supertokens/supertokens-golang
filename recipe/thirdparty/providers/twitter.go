package providers

import (
	"encoding/base64"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Twitter(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.Name == "" {
		input.Config.Name = "Twitter"
	}

	if input.Config.AuthorizationEndpoint == "" {
		input.Config.AuthorizationEndpoint = "https://twitter.com/i/oauth2/authorize"
	}

	if input.Config.TokenEndpoint == "" {
		input.Config.TokenEndpoint = "https://api.twitter.com/2/oauth2/token"
	}

	if input.Config.UserInfoEndpoint == "" {
		input.Config.UserInfoEndpoint = "https://api.twitter.com/2/users/me"
	}

	if input.Config.RequireEmail == nil {
		False := false
		input.Config.RequireEmail = &False
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "data.id"
	}

	oOverride := input.Override

	input.Override = func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := originalImplementation.GetConfigForClientType
		originalImplementation.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"users.read", "tweet.read"}
			}

			if config.ForcePKCE == nil {
				True := true
				config.ForcePKCE = &True
			}

			return config, nil
		}

		originalImplementation.ExchangeAuthCodeForOAuthTokens = func(redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			basicAuthToken := base64.StdEncoding.EncodeToString([]byte(originalImplementation.Config.ClientID + ":" + originalImplementation.Config.ClientSecret))
			twitterOauthParams := map[string]interface{}{}

			if originalImplementation.Config.TokenEndpointBodyParams != nil {
				twitterOauthParams = originalImplementation.Config.TokenEndpointBodyParams
			}

			twitterOauthParams["grant_type"] = "authorization_code"
			twitterOauthParams["client_id"] = originalImplementation.Config.ClientID
			twitterOauthParams["code_verifier"] = redirectURIInfo.PKCECodeVerifier
			twitterOauthParams["redirect_uri"] = redirectURIInfo.RedirectURIOnProviderDashboard
			twitterOauthParams["code"] = redirectURIInfo.RedirectURIQueryParams["code"]

			return doPostRequest(originalImplementation.Config.TokenEndpoint, twitterOauthParams, map[string]interface{}{
				"Authorization": "Basic " + basicAuthToken,
			})
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}

		return originalImplementation
	}

	return NewProvider(input)
}
