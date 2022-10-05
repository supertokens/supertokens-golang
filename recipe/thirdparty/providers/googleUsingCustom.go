package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var googleBuilder = BuildCustomProvider(tpmodels.TypeCustomProviderBuilderInput{
	ID:                            googleID,
	AuthorizationURL:              "https://accounts.google.com/o/oauth2/v2/auth",
	ExchangeCodeForOAuthTokensURL: "https://accounts.google.com/o/oauth2/token",
	GetUserInfoURL:                "https://www.googleapis.com/oauth2/v1/userinfo?alt=json",
	DefaultScopes:                 []string{"https://www.googleapis.com/auth/userinfo.email"},
})

func GoogleUsingCustom(input tpmodels.TypeGoogleInput) (tpmodels.TypeProvider, error) {
	configs := make([]tpmodels.CustomProviderConfig, len(input.Config))
	for i, config := range input.Config {
		configs[i] = tpmodels.CustomProviderConfig{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Scope:        config.Scope,
		}
	}

	if input.Override == nil {
		return googleBuilder(tpmodels.TypeCustomProviderInput{
			Config: configs,
		})
	}

	googleProvider, err := googleBuilder(tpmodels.TypeCustomProviderInput{
		Config: configs,
		Override: func(provider *tpmodels.CustomProvider) *tpmodels.CustomProvider {
			overrideProvider := &tpmodels.GoogleProvider{}
			overrideProvider.GetConfig = func(clientID *string, userContext supertokens.UserContext) (tpmodels.GoogleConfig, error) {
				config, err := provider.GetConfig(clientID, userContext)
				if err != nil {
					return tpmodels.GoogleConfig{}, err
				}
				return tpmodels.GoogleConfig{
					ClientID:     config.ClientID,
					ClientSecret: config.ClientSecret,
					Scope:        config.Scope,
				}, nil
			}

			overrideProvider.GetAuthorisationRedirectURL = func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
				return provider.GetAuthorisationRedirectURL(clientID, redirectURIOnProviderDashboard, userContext)
			}

			overrideProvider.ExchangeAuthCodeForOAuthTokens = func(clientID *string, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
				return provider.ExchangeAuthCodeForOAuthTokens(clientID, redirectInfo, userContext)
			}

			overrideProvider.GetUserInfo = func(clientID *string, accessToken tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
				return provider.GetUserInfo(clientID, accessToken, userContext)
			}

			overriddenProvider := input.Override(overrideProvider)

			return &tpmodels.CustomProvider{
				GetConfig: func(clientID *string, userContext supertokens.UserContext) (tpmodels.CustomProviderConfig, error) {
					config, err := overriddenProvider.GetConfig(clientID, userContext)
					if err != nil {
						return tpmodels.CustomProviderConfig{}, err
					}
					return tpmodels.CustomProviderConfig{
						ClientID:     config.ClientID,
						ClientSecret: config.ClientSecret,
						Scope:        config.Scope,
					}, nil
				},
				TypeProvider: overriddenProvider.TypeProvider,
			}
		},
	})

	return googleProvider, err
}
