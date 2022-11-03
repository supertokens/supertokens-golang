/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const facebookID = "facebook"

type FacebookConfig = OAuth2ProviderConfig

type TypeFacebookInput struct {
	Config   FacebookConfig
	Override func(provider *FacebookProvider) *FacebookProvider
}

type FacebookProvider = tpmodels.TypeProvider

func Facebook(input TypeFacebookInput) tpmodels.TypeProvider {
	facebookProvider := &FacebookProvider{
		ID: facebookID,
	}

	oAuth2Provider := oAuth2Provider(TypeOAuth2ProviderInput{
		ThirdPartyID: facebookID,
		Config:       input.Config,
	})

	{
		// Facebook provider APIs call into oAuth2 provider APIs

		facebookProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
			return oAuth2Provider.GetConfig(clientType, tenantId, userContext)
		}

		facebookProvider.GetAuthorisationRedirectURL = func(config tpmodels.TypeNormalisedProviderConfig, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return oAuth2Provider.GetAuthorisationRedirectURL(config, redirectURIOnProviderDashboard, userContext)
		}

		facebookProvider.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.TypeNormalisedProviderConfig, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return oAuth2Provider.ExchangeAuthCodeForOAuthTokens(config, redirectInfo, userContext)
		}
	}

	facebookProvider.GetUserInfo = func(config tpmodels.TypeNormalisedProviderConfig, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		accessToken := oAuthTokens["access_token"].(string)

		queryParams := map[string]interface{}{
			"access_token": accessToken,
			"fields":       "id,email",
			"format":       "json",
		}

		userInfoFromAccessToken, err := doGetRequest(config.UserInfoEndpoint, queryParams, nil)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		rawUserInfoFromProvider := tpmodels.TypeRawUserInfoFromProvider{
			FromAccessToken: userInfoFromAccessToken.(map[string]interface{}),
		}
		userInfoResult, err := getSupertokensUserInfoResultFromRawUserInfo(rawUserInfoFromProvider, config)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		return tpmodels.TypeUserInfo{
			ThirdPartyUserId:        userInfoResult.ThirdPartyUserId,
			EmailInfo:               userInfoResult.EmailInfo,
			RawUserInfoFromProvider: rawUserInfoFromProvider,
		}, nil
	}

	if input.Override != nil {
		facebookProvider = input.Override(facebookProvider)
	}

	{
		// We want to always normalize (for apple) the config before returning it
		oGetConfig := facebookProvider.GetConfig
		facebookProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
			config, err := oGetConfig(clientType, tenantId, userContext)
			if err != nil {
				return tpmodels.TypeNormalisedProviderConfig{}, err
			}
			return normalizeFacebookConfig(config), nil
		}
	}

	return *facebookProvider
}

func normalizeFacebookConfig(config tpmodels.TypeNormalisedProviderConfig) tpmodels.TypeNormalisedProviderConfig {
	if config.AuthorizationEndpoint == "" {
		config.AuthorizationEndpoint = "https://www.facebook.com/v12.0/dialog/oauth"
	}

	if config.TokenEndpoint == "" {
		config.TokenEndpoint = "https://graph.facebook.com/v12.0/oauth/access_token"
	}

	if config.UserInfoEndpoint == "" {
		config.UserInfoEndpoint = "https://graph.facebook.com/me"
	}

	if len(config.Scope) == 0 {
		config.Scope = []string{"email"}
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
