package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

const googleID = "google"

type GoogleConfig = CustomConfig
type GoogleClientConfig = CustomClientConfig

type TypeGoogle = TypeCustomProviderImplementation

type Google struct {
	Config   CustomConfig
	Override func(provider *TypeGoogle) *TypeGoogle
}

func (input Google) GetID() string {
	return googleID
}

func (input Google) Build() tpmodels.TypeProvider {
	googleImpl := input.buildInternal()
	if input.Override != nil {
		googleImpl = input.Override(googleImpl)
	}
	return *googleImpl.TypeProvider
}

func (input Google) buildInternal() *TypeGoogle {
	return (CustomProvider{
		ThirdPartyID: googleID,
		Config:       input.Config,

		oAuth2Normalize: normalizeOAuth2ConfigForGoogle,
	}).buildInternal()
}

func normalizeOAuth2ConfigForGoogle(config *typeCombinedOAuth2Config) *typeCombinedOAuth2Config {
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
		config.UserInfoMap.From = tpmodels.FromUserInfoAPI
	}

	if config.UserInfoMap.UserId == "" {
		config.UserInfoMap.UserId = "id"
	}

	if config.UserInfoMap.Email == "" {
		config.UserInfoMap.Email = "email"
	}

	if config.UserInfoMap.EmailVerified == "" {
		config.UserInfoMap.EmailVerified = "email_verified"
	}

	return config
}
