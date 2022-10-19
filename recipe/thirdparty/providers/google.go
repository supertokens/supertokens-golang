package providers

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const googleID = "google"

type TypeGoogleInput struct {
	Config   []GoogleConfig
	Override func(provider *GoogleProvider) *GoogleProvider
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type GoogleProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (GoogleConfig, error)
	*tpmodels.TypeProvider
}

func Google(input TypeGoogleInput) tpmodels.TypeProvider {
	googleProvider := &GoogleProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: googleID,
		},
	}

	googleProvider.GetConfig = func(clientID *string, userContext supertokens.UserContext) (GoogleConfig, error) {
		if len(input.Config) == 0 {
			return GoogleConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if clientID == nil && len(input.Config) > 1 {
			return GoogleConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if clientID == nil {
			return input.Config[0], nil
		}

		for _, config := range input.Config {
			if config.ClientID == *clientID {
				return config, nil
			}
		}

		return GoogleConfig{}, errors.New("config for specified clientID not found")
	}

	customProvider := CustomProvider(TypeCustomProviderInput{
		ThirdPartyID: googleID,
		Override: func(provider *TypeCustomProvider) *TypeCustomProvider {
			provider.GetConfig = func(clientID *string, userContext supertokens.UserContext) (CustomProviderConfig, error) {
				googleConfig, err := googleProvider.GetConfig(clientID, userContext)
				if err != nil {
					return CustomProviderConfig{}, err
				}

				authURL := "https://accounts.google.com/o/oauth2/v2/auth"
				tokenURL := "https://oauth2.googleapis.com/token"
				userInfoURL := "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

				accessType := "offline"
				if googleConfig.ClientSecret == "" {
					accessType = "online"
				}

				return CustomProviderConfig{
					ClientID:     googleConfig.ClientID,
					ClientSecret: googleConfig.ClientSecret,
					Scope:        googleConfig.Scope,

					AuthorizationURL: &authURL,
					AuthorizationURLQueryParams: map[string]interface{}{
						"access_type":            accessType,
						"include_granted_scopes": "true",
						"response_type":          "code",
					},
					AccessTokenURL: &tokenURL,
					AccessTokenParams: map[string]interface{}{
						"grant_type": "authorization_code",
					},
					UserInfoURL: &userInfoURL,
				}, nil
			}

			return provider
		},
	})

	googleProvider.GetAuthorisationRedirectURL = func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		return customProvider.GetAuthorisationRedirectURL(clientID, redirectURIOnProviderDashboard, userContext)
	}

	googleProvider.ExchangeAuthCodeForOAuthTokens = func(clientID *string, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		return customProvider.ExchangeAuthCodeForOAuthTokens(clientID, redirectInfo, userContext)
	}

	googleProvider.GetUserInfo = func(clientID *string, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		return customProvider.GetUserInfo(clientID, oAuthTokens, userContext)
	}

	if input.Override != nil {
		googleProvider = input.Override(googleProvider)
	}

	return customProvider
}
