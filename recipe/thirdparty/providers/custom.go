package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func BuildCustomProvider(builderInput tpmodels.TypeCustomProviderBuilderInput) tpmodels.CustomProviderFunc {
	return func(input tpmodels.TypeCustomProviderInput) (tpmodels.TypeProvider, error) {
		provider := &tpmodels.CustomProvider{
			TypeProvider: &tpmodels.TypeProvider{
				ID: builderInput.ID,
			},
		}

		provider.GetAuthorisationRedirectURL = func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			// ...
			return tpmodels.TypeAuthorisationRedirect{}, nil
		}

		provider.ExchangeAuthCodeForOAuthTokens = func(clientID *string, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			// ...
			return tpmodels.TypeOAuthTokens{}, nil
		}

		provider.GetUserInfo = func(clientID *string, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			// ...
			return tpmodels.TypeUserInfo{}, nil
		}

		if input.Override != nil {
			provider = input.Override(provider)
		}

		return *provider.TypeProvider, nil
	}
}
