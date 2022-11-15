package providers

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const linkedinID = "linkedin"

func Linkedin(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	if input.ThirdPartyID == "" {
		input.ThirdPartyID = linkedinID
	}

	if input.Config.AuthorizationEndpoint == "" {
		input.Config.AuthorizationEndpoint = "https://www.linkedin.com/oauth/v2/authorization"
	}

	if input.Config.TokenEndpoint == "" {
		input.Config.TokenEndpoint = "https://www.linkedin.com/oauth/v2/accessToken"
	}

	if input.Config.UserInfoEndpoint == "" {
		input.Config.UserInfoEndpoint = "https://api.linkedin.com/v2/me"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "id"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.EmailVerified = "email_verified"
	}

	if input.Config.AuthorizationEndpointQueryParams == nil {
		input.Config.AuthorizationEndpointQueryParams = map[string]interface{}{}
	}

	if input.Config.AuthorizationEndpointQueryParams["include_granted_scopes"] == nil {
		input.Config.AuthorizationEndpointQueryParams["include_granted_scopes"] = "true"
	}
	if input.Config.AuthorizationEndpointQueryParams["access_type"] == nil {
		input.Config.AuthorizationEndpointQueryParams["access_type"] = "offline"
	}

	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfigForClientType
		provider.GetConfigForClientType = func(clientType *string, input tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, input, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"r_emailaddress", "r_liteprofile"}
			}

			return config, err
		}

		provider.GetUserInfo = func(config tpmodels.ProviderConfigForClientType, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			accessToken, accessTokenOk := oAuthTokens["access_token"].(string)
			if !accessTokenOk {
				return tpmodels.TypeUserInfo{}, errors.New("access token not found")
			}

			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			rawUserInfoFromProvider := tpmodels.TypeRawUserInfoFromProvider{}
			userInfoFromAccessToken, err := doGetRequest(config.UserInfoEndpoint, nil, headers)
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

			userInfoResult, err := oauth2_getSupertokensUserInfoResultFromRawUserInfo(config, rawUserInfoFromProvider)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}

			if config.TenantId != "" {
				userInfoResult.ThirdPartyUserId += "|" + config.TenantId
			}

			return tpmodels.TypeUserInfo{
				ThirdPartyUserId:        userInfoResult.ThirdPartyUserId,
				Email:                   userInfoResult.EmailInfo,
				RawUserInfoFromProvider: rawUserInfoFromProvider,
			}, nil
		}

		if oOverride != nil {
			provider = oOverride(provider)
		}
		return provider
	}

	return NewProvider(input)
}
