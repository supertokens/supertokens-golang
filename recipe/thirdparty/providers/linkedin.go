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
				config.Scope = []string{"r_emailaddress", "r_liteprofile"}
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
			userInfoFromAccessToken, err := doGetRequest("https://api.linkedin.com/v2/me", nil, headers)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
			rawUserInfoFromProvider.FromUserInfoAPI = userInfoFromAccessToken.(map[string]interface{})

			emailAPIURL := "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))"
			userInfoFromEmail, err := doGetRequest(emailAPIURL, nil, headers)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}

			elements := userInfoFromEmail.(map[string]interface{})["elements"].([]interface{})
			for _, elem := range elements {
				if elemMap, ok := elem.(map[string]interface{}); ok {
					for k, v := range elemMap {
						if k == "handle~" {
							emailMap := v.(map[string]interface{})
							rawUserInfoFromProvider.FromUserInfoAPI["email"] = emailMap["emailAddress"]
						}
					}
				}
			}

			for k, v := range userInfoFromEmail.(map[string]interface{}) {
				rawUserInfoFromProvider.FromUserInfoAPI[k] = v
			}

			userInfoResult := tpmodels.TypeUserInfo{
				ThirdPartyUserId: rawUserInfoFromProvider.FromUserInfoAPI["id"].(string),
				Email: &tpmodels.EmailStruct{
					ID: rawUserInfoFromProvider.FromUserInfoAPI["email"].(string),
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
