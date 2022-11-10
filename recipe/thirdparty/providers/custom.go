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

type CustomConfig struct {
	Clients []CustomClientConfig

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	ForcePKCE                        bool // Providers like twitter expects PKCE to be used along with secret
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      tpmodels.TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig CustomClientConfig) (bool, error)
}

type CustomClientConfig struct {
	ClientType       string // optional
	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}
}

/*
ToCustomProviderClientConfig and UpdateFromCustomProviderClientConfig are used to convert between
CustomClientConfig and Provider specific client config. This helps us re-use a lot of code,
including the logic for GetConfig using the findConfig function.
*/
func (config CustomClientConfig) ToCustomProviderClientConfig() (CustomClientConfig, error) {
	return config, nil
}

func (config *CustomClientConfig) UpdateFromCustomProviderClientConfig(input CustomClientConfig) {
	config.ClientType = input.ClientType
	config.ClientID = input.ClientID
	config.ClientSecret = input.ClientSecret
	config.Scope = input.Scope
	config.AdditionalConfig = input.AdditionalConfig
}

type TypeCustomProviderImplementation struct {
	GetConfig func(clientType *string, tenantId *string, userContext supertokens.UserContext) (CustomClientConfig, error)
	*tpmodels.TypeProvider
}

type CustomProvider struct {
	ThirdPartyID string
	Config       CustomConfig
	Override     func(provider *TypeCustomProviderImplementation) *TypeCustomProviderImplementation

	// for internal use: to be used by built-in providers to populate provider specific config
	oAuth2Normalize func(config *typeCombinedOAuth2Config) *typeCombinedOAuth2Config
}

func (input CustomProvider) GetID() string {
	return input.ThirdPartyID
}

func (input CustomProvider) Build() tpmodels.TypeProvider {
	customImpl := input.buildInternal()
	if input.Override != nil {
		customImpl = input.Override(customImpl)
	}
	return *customImpl.TypeProvider
}

func (input CustomProvider) buildInternal() *TypeCustomProviderImplementation {
	if input.oAuth2Normalize == nil {
		input.oAuth2Normalize = func(config *typeCombinedOAuth2Config) *typeCombinedOAuth2Config {
			return config
		}
	}

	customProvider := &TypeCustomProviderImplementation{
		TypeProvider: &tpmodels.TypeProvider{
			ID: input.ThirdPartyID,
		},
	}

	customProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (CustomClientConfig, error) {
		clientConfig := CustomClientConfig{}
		clients := make([]TypeToCustomProvider, len(input.Config.Clients))
		for i, client := range input.Config.Clients {
			clients[i] = client
		}

		err := findConfig(&clientConfig, clientType, tenantId, userContext, clients)
		if err != nil {
			return CustomClientConfig{}, err
		}
		return clientConfig, nil
	}

	customProvider.GetAuthorisationRedirectURL = func(clientType, tenantId *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		clientConfig, err := customProvider.GetConfig(clientType, tenantId, userContext)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		return input.oAuth2Normalize(getCombinedOAuth2Config(input.Config, clientConfig)).GetAuthorisationRedirectURL(redirectURIOnProviderDashboard, userContext)
	}

	customProvider.ExchangeAuthCodeForOAuthTokens = func(clientType, tenantId *string, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		clientConfig, err := customProvider.GetConfig(clientType, tenantId, userContext)
		if err != nil {
			return tpmodels.TypeOAuthTokens{}, err
		}

		return input.oAuth2Normalize(getCombinedOAuth2Config(input.Config, clientConfig)).ExchangeAuthCodeForOAuthTokens(redirectURIInfo, userContext)
	}

	customProvider.GetUserInfo = func(clientType, tenantId *string, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		clientConfig, err := customProvider.GetConfig(clientType, tenantId, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		return input.oAuth2Normalize(getCombinedOAuth2Config(input.Config, clientConfig)).GetUserInfo(oAuthTokens, userContext)
	}

	if input.Override != nil {
		customProvider = input.Override(customProvider)
	}

	return customProvider
}
