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

type TypeCustomProviderInput struct {
	ThirdPartyID string
	Config       CustomProviderConfig
	Override     func(provider *TypeCustomProvider) *TypeCustomProvider
}

type CustomProviderConfig struct {
	Clients []CustomProviderClientConfig

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	ForcePKCE                        bool
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      tpmodels.TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig CustomProviderClientConfig) (bool, error)
}

type CustomProviderClientConfig struct {
	ClientType       string
	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}
}

func (config CustomProviderClientConfig) ToCustomProviderClientConfig() (CustomProviderClientConfig, error) {
	return config, nil
}

func (config *CustomProviderClientConfig) UpdateFromCustomProviderClientConfig(input CustomProviderClientConfig) {
	config.ClientType = input.ClientType
	config.ClientID = input.ClientID
	config.ClientSecret = input.ClientSecret
	config.Scope = input.Scope
	config.AdditionalConfig = input.AdditionalConfig
}

type TypeSupertokensUserInfoMap struct {
	IdField            string
	EmailField         string
	EmailVerifiedField string
}

const scopeParameter = "scope"
const scopeSeparator = " "

type TypeCustomProvider struct {
	GetConfig func(clientType *string, tenantId *string, userContext supertokens.UserContext) (CustomProviderClientConfig, error)
	*tpmodels.TypeProvider
}

func customProvider(input TypeCustomProviderInput, oAuth2Normalize func(config *typeCombinedOAuth2Config) *typeCombinedOAuth2Config) *TypeCustomProvider {
	if oAuth2Normalize == nil {
		oAuth2Normalize = func(config *typeCombinedOAuth2Config) *typeCombinedOAuth2Config {
			return config
		}
	}

	customProvider := &TypeCustomProvider{
		TypeProvider: &tpmodels.TypeProvider{
			ID: input.ThirdPartyID,
		},
	}

	customProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (CustomProviderClientConfig, error) {
		clientConfig := CustomProviderClientConfig{}
		clients := make([]TypeToCustomProvider, len(input.Config.Clients))
		for i, client := range input.Config.Clients {
			clients[i] = client
		}

		err := findConfig(&clientConfig, clientType, tenantId, userContext, clients)
		if err != nil {
			return CustomProviderClientConfig{}, err
		}
		return clientConfig, nil
	}

	customProvider.GetAuthorisationRedirectURL = func(clientType *string, tenantId *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		clientConfig, err := customProvider.GetConfig(clientType, tenantId, userContext)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		return oAuth2Normalize(getCombinedOAuth2Config(input.Config, clientConfig)).GetAuthorisationRedirectURL(redirectURIOnProviderDashboard, userContext)
	}

	customProvider.ExchangeAuthCodeForOAuthTokens = func(clientType *string, tenantId *string, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		clientConfig, err := customProvider.GetConfig(clientType, tenantId, userContext)
		if err != nil {
			return tpmodels.TypeOAuthTokens{}, err
		}

		return oAuth2Normalize(getCombinedOAuth2Config(input.Config, clientConfig)).ExchangeAuthCodeForOAuthTokens(redirectURIInfo, userContext)

	}

	customProvider.GetUserInfo = func(clientType *string, tenantId *string, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		clientConfig, err := customProvider.GetConfig(clientType, tenantId, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		return oAuth2Normalize(getCombinedOAuth2Config(input.Config, clientConfig)).GetUserInfo(oAuthTokens, userContext)
	}

	if input.Override != nil {
		customProvider = input.Override(customProvider)
	}

	return customProvider
}

func CustomProvider(input TypeCustomProviderInput) tpmodels.TypeProvider {
	return *customProvider(input, nil).TypeProvider
}
