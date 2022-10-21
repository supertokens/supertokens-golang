package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const googleID = "google"

type GoogleConfig = CustomProviderConfig

type TypeGoogleInput struct {
	Config   []GoogleConfig
	Override func(provider *GoogleProvider) *GoogleProvider
}

type GoogleProvider struct {
	GetConfig func(id *tpmodels.TypeID, userContext supertokens.UserContext) (GoogleConfig, error)
	*tpmodels.TypeProvider
}

func Google(input TypeGoogleInput) tpmodels.TypeProvider {
	googleProvider := &GoogleProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: googleID,
		},
	}

	var customProviderConfig []CustomProviderConfig
	if input.Config != nil {
		customProviderConfig = make([]CustomProviderConfig, len(input.Config))
		for idx, config := range input.Config {
			customProviderConfig[idx] = config
		}
	}

	customProvider := customProvider(TypeCustomProviderInput{
		ThirdPartyID: googleID,
		Config:       customProviderConfig,
	})

	{
		// Custom provider needs to use the config returned by google provider GetConfig
		// Also, google provider needs to use the default implementation of GetConfig provided by custom provider
		oGetConfig := customProvider.GetConfig
		customProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
			return googleProvider.GetConfig(id, userContext)
		}
		googleProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (GoogleConfig, error) {
			return oGetConfig(id, userContext)
		}
	}

	{
		// Google provider APIs call into custom provider APIs

		googleProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return customProvider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
		}

		googleProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return customProvider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
		}

		googleProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return customProvider.GetUserInfo(id, oAuthTokens, userContext)
		}
	}

	if input.Override != nil {
		googleProvider = input.Override(googleProvider)
	}

	{
		// We want to always normalize (for google) the config before returning it
		oGetConfig := googleProvider.GetConfig
		googleProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (GoogleConfig, error) {
			config, err := oGetConfig(id, userContext)
			if err != nil {
				return GoogleConfig{}, err
			}
			return normalizeGoogleConfig(config), nil
		}
	}

	return *googleProvider.TypeProvider
}

func normalizeGoogleConfig(config GoogleConfig) GoogleConfig {
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

	if config.GetSupertokensUserInfoFromRawUserInfoResponse == nil {
		config.GetSupertokensUserInfoFromRawUserInfoResponse = getSupertokensUserInfoFromRawUserInfo("id", "email", "email_verified", "access_token")
	}

	return config
}
