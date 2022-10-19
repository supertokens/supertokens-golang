package providers

import (
	"errors"
	"fmt"

	"github.com/derekstavis/go-qs"
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
	GetConfig func(ID *tpmodels.TypeID, userContext supertokens.UserContext) (GoogleConfig, error)
	*tpmodels.TypeProvider
}

func Google(input TypeGoogleInput) tpmodels.TypeProvider {
	googleProvider := &GoogleProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: googleID,
		},
	}

	googleProvider.GetConfig = func(ID *tpmodels.TypeID, userContext supertokens.UserContext) (GoogleConfig, error) {
		if ID == nil && len(input.Config) == 0 {
			return GoogleConfig{}, errors.New("please specify a config or override GetConfig")
		}

		if ID == nil && len(input.Config) > 1 {
			return GoogleConfig{}, errors.New("please specify a clientID as there are multiple configs")
		}

		if ID == nil {
			return input.Config[0], nil
		}

		if ID.Type == tpmodels.TypeClientID {
			for _, config := range input.Config {
				if config.ClientID == ID.ID {
					return config, nil
				}
			}
		} else {
			// TODO Multitenant
		}

		return GoogleConfig{}, errors.New("config for specified clientID not found")
	}

	customProvider := CustomProvider(TypeCustomProviderInput{
		ThirdPartyID: googleID,
		Override: func(provider *TypeCustomProvider) *TypeCustomProvider {
			provider.GetConfig = func(ID *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
				googleConfig, err := googleProvider.GetConfig(ID, userContext)
				if err != nil {
					return CustomProviderConfig{}, err
				}

				authURL := "https://accounts.google.com/o/oauth2/v2/auth"
				tokenURL := "https://oauth2.googleapis.com/token"
				userInfoURL := "https://www.googleapis.com/oauth2/v1/userinfo"

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
					UserInfoURL:    &userInfoURL,
					DefaultScope:   []string{"https://www.googleapis.com/auth/userinfo.email"},
					GetSupertokensUserFromRawResponse: func(rawResponse map[string]interface{}, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
						result := tpmodels.TypeUserInfo{}
						result.ThirdPartyUserId = fmt.Sprint(rawResponse["id"])
						result.EmailInfo = &tpmodels.EmailStruct{
							ID: fmt.Sprint(rawResponse["email"]),
						}
						emailVerified, emailVerifiedOk := rawResponse["email_verified"].(bool)
						result.EmailInfo.IsVerified = emailVerified && emailVerifiedOk

						return result, nil
					},
				}, nil
			}

			oGetAuthorisationRedirectURL := provider.GetAuthorisationRedirectURL
			provider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
				config, err := provider.GetConfig(id, userContext)
				if err != nil {
					return tpmodels.TypeAuthorisationRedirect{}, err
				}
				result, err := oGetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
				if err != nil {
					return tpmodels.TypeAuthorisationRedirect{}, err
				}

				if config.ClientSecret == "" {
					challenge, verifier, err := generateCodeChallengeS256(32)
					if err != nil {
						return tpmodels.TypeAuthorisationRedirect{}, err
					}
					extraQueryParams, err := qs.Marshal(map[string]interface{}{
						"code_challenge":        challenge,
						"code_challenge_method": "S256",
					})
					if err != nil {
						return tpmodels.TypeAuthorisationRedirect{}, err
					}
					result.URLWithQueryParams += "&" + extraQueryParams
					result.PKCECodeVerifier = &verifier
				}

				return result, nil
			}

			return provider
		},
	})

	googleProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		return customProvider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
	}

	googleProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		return customProvider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
	}

	googleProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		return customProvider.GetUserInfo(id, oAuthTokens, userContext)
	}

	if input.Override != nil {
		googleProvider = input.Override(googleProvider)
	}

	return customProvider
}
