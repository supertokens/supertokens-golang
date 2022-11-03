package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const googleID = "google"

type GoogleConfig = OAuth2ProviderConfig

type TypeGoogleInput struct {
	Config   GoogleConfig
	Override func(provider *GoogleProvider) *GoogleProvider
}

type GoogleProvider struct {
	*tpmodels.TypeProvider
}

func Google(input TypeGoogleInput) tpmodels.TypeProvider {
	googleProvider := &GoogleProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: googleID,
		},
	}

	oAuth2Provider := oAuth2Provider(TypeOAuth2ProviderInput{
		ThirdPartyID: googleID,
		Config:       input.Config,
	})

	{
		// Google provider APIs call into oAuth2 provider APIs

		googleProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
			return oAuth2Provider.GetConfig(clientType, tenantId, userContext)
		}

		googleProvider.GetAuthorisationRedirectURL = func(config tpmodels.TypeNormalisedProviderConfig, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return oAuth2Provider.GetAuthorisationRedirectURL(config, redirectURIOnProviderDashboard, userContext)
		}

		googleProvider.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.TypeNormalisedProviderConfig, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return oAuth2Provider.ExchangeAuthCodeForOAuthTokens(config, redirectInfo, userContext)
		}

		googleProvider.GetUserInfo = func(config tpmodels.TypeNormalisedProviderConfig, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return oAuth2Provider.GetUserInfo(config, oAuthTokens, userContext)
		}
	}

	if input.Override != nil {
		googleProvider = input.Override(googleProvider)
	}

	{
		// We want to always normalize (for google) the config before returning it
		oGetConfig := googleProvider.GetConfig
		googleProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
			config, err := oGetConfig(clientType, tenantId, userContext)
			if err != nil {
				return tpmodels.TypeNormalisedProviderConfig{}, err
			}
			return normalizeGoogleConfig(config), nil
		}
	}

	return *googleProvider.TypeProvider
}

func normalizeGoogleConfig(config tpmodels.TypeNormalisedProviderConfig) tpmodels.TypeNormalisedProviderConfig {
	if config.AuthorizationEndpoint == "" {
		config.AuthorizationEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	}

	if config.AuthorizationEndpointQueryParams == nil {
		accessType := "offline"
		if config.ClientSecret == "" {
			accessType = "online"
		}
		config.AuthorizationEndpointQueryParams = map[string]interface{}{
			"access_type":            accessType,
			"include_granted_scopes": "true",
			"response_type":          "code",
		}
	}

	if len(config.Scope) == 0 {
		config.Scope = []string{"https://www.googleapis.com/auth/userinfo.email"}
	}

	if config.TokenEndpoint == "" {
		config.TokenEndpoint = "https://oauth2.googleapis.com/token"
	}

	if config.UserInfoEndpoint == "" {
		config.UserInfoEndpoint = "https://www.googleapis.com/oauth2/v1/userinfo"
	}

	if config.UserInfoMap.From == "" {
		config.UserInfoMap.From = tpmodels.FromAccessTokenPayload
	}

	if config.UserInfoMap.IdField == "" {
		config.UserInfoMap.IdField = "id"
	}

	if config.UserInfoMap.EmailField == "" {
		config.UserInfoMap.EmailField = "email"
	}

	if config.UserInfoMap.EmailVerifiedField == "" {
		config.UserInfoMap.EmailVerifiedField = "email_verified"
	}

	return config
}
