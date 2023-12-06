package providers

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Linkedin(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.Name == "" {
		input.Config.Name = "LinkedIn"
	}

	if input.Config.AuthorizationEndpoint == "" {
		input.Config.AuthorizationEndpoint = "https://www.linkedin.com/oauth/v2/authorization"
	}

	if input.Config.TokenEndpoint == "" {
		input.Config.TokenEndpoint = "https://www.linkedin.com/oauth/v2/accessToken"
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
				// https://learn.microsoft.com/en-us/linkedin/consumer/integrations/self-serve/sign-in-with-linkedin-v2?context=linkedin%2Fconsumer%2Fcontext#authenticating-members
				config.Scope = []string{"openid", "profile", "email"}
			}

			return config, nil
		}

		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			accessToken, accessTokenOk := oAuthTokens["access_token"].(string)
			if !accessTokenOk {
				return tpmodels.TypeUserInfo{}, errors.New("access token not found")
			}

			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			rawUserInfoFromProvider := tpmodels.TypeRawUserInfoFromProvider{}
			// https://learn.microsoft.com/en-us/linkedin/consumer/integrations/self-serve/sign-in-with-linkedin-v2?context=linkedin%2Fconsumer%2Fcontext#sample-api-response
			userInfoFromAccessToken, err := doGetRequest("https://api.linkedin.com/v2/userinfo", nil, headers)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
			rawUserInfoFromProvider.FromUserInfoAPI = userInfoFromAccessToken.(map[string]interface{})

			userInfoResult := tpmodels.TypeUserInfo{
				ThirdPartyUserId: rawUserInfoFromProvider.FromUserInfoAPI["sub"].(string),
				Email: &tpmodels.EmailStruct{
					ID:         rawUserInfoFromProvider.FromUserInfoAPI["email"].(string),
					IsVerified: rawUserInfoFromProvider.FromUserInfoAPI["email_verified"].(bool),
				},
			}

			return tpmodels.TypeUserInfo{
				ThirdPartyUserId:        userInfoResult.ThirdPartyUserId,
				Email:                   userInfoResult.Email,
				RawUserInfoFromProvider: rawUserInfoFromProvider,
			}, nil
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}
		return originalImplementation
	}

	return NewProvider(input)
}
