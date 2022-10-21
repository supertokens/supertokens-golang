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
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const facebookID = "facebook"

type FacebookConfig = CustomProviderConfig

type TypeFacebookInput struct {
	Config   []FacebookConfig
	Override func(provider *FacebookProvider) *FacebookProvider
}

type FacebookProvider struct {
	GetConfig func(id *tpmodels.TypeID, userContext supertokens.UserContext) (FacebookConfig, error)
	*tpmodels.TypeProvider
}

func Facebook(input TypeFacebookInput) tpmodels.TypeProvider {
	facebookProvider := &FacebookProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: facebookID,
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
		// Custom provider needs to use the config returned by facebook provider GetConfig
		// Also, facebook provider needs to use the default implementation of GetConfig provided by custom provider
		oGetConfig := customProvider.GetConfig
		customProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (CustomProviderConfig, error) {
			return facebookProvider.GetConfig(id, userContext)
		}
		facebookProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (FacebookConfig, error) {
			return oGetConfig(id, userContext)
		}
	}

	{
		// Facebook provider APIs call into custom provider APIs

		facebookProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
			return customProvider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
		}

		facebookProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
			return customProvider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
		}
	}

	facebookProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		config, err := facebookProvider.GetConfig(id, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		accessToken := oAuthTokens["access_token"].(string)

		queryParams := map[string]interface{}{
			"access_token": accessToken,
			"fields":       "id,email",
			"format":       "json",
		}

		status, userInfoFromAccessToken, err := doGetRequest(config.UserInfoEndpoint, queryParams, nil)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		if status >= 300 {
			return tpmodels.TypeUserInfo{}, errors.New("get user info returned a non 2xx response") // TODO return code and response
		}

		rawUserInfoFromProvider := tpmodels.TypeRawUserInfoFromProvider{
			FromAccessToken: userInfoFromAccessToken.(map[string]interface{}),
		}
		userInfoResult, err := config.GetSupertokensUserInfoFromRawUserInfoResponse(rawUserInfoFromProvider, userContext)
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
		facebookProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (FacebookConfig, error) {
			config, err := oGetConfig(id, userContext)
			if err != nil {
				return FacebookConfig{}, err
			}
			return normalizeFacebookConfig(config), nil
		}
	}

	return *facebookProvider.TypeProvider
}

func normalizeFacebookConfig(config FacebookConfig) FacebookConfig {
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

	if config.GetSupertokensUserInfoFromRawUserInfoResponse == nil {
		config.GetSupertokensUserInfoFromRawUserInfoResponse = getSupertokensUserInfoFromRawUserInfo("id", "email", "email_verified", "access_token")
	}

	return config
}
