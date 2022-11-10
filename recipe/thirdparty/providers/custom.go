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

func NewProvider(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	impl := &tpmodels.TypeProvider{
		ID: input.ThirdPartyID,
	}

	impl.GetConfig = func(clientType, tenantId *string, inputConfig tpmodels.ProviderConfigInput, userContext supertokens.UserContext) (tpmodels.ProviderClientConfig, error) {
		if tenantId == nil {
			if clientType == nil {

				if len(input.Config.Clients) == 0 || len(input.Config.Clients) > 1 {
					return tpmodels.ProviderClientConfig{}, errors.New("please provide exactly one client config or pass clientType or tenantId")
				}
				for _, client := range input.Config.Clients {
					if client.ClientType == *clientType {
						config := getCombinedProviderConfig(input.Config, client)
						// Discover the end points here
						return config, nil
					}
				}
			}

			return tpmodels.ProviderClientConfig{}, errors.New("Could not find client config for clientType: " + *clientType)
		}

		// TODO impl
		return tpmodels.ProviderClientConfig{}, errors.New("needs implementation")
	}

	impl.GetAuthorisationRedirectURL = func(clientType, tenantId *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		clientConfig, err := impl.GetConfig(clientType, tenantId, input.Config, userContext)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		return oauth2_GetAuthorisationRedirectURL(clientConfig, redirectURIOnProviderDashboard, userContext)
	}

	impl.ExchangeAuthCodeForOAuthTokens = func(clientType, tenantId *string, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		clientConfig, err := impl.GetConfig(clientType, tenantId, input.Config, userContext)
		if err != nil {
			return tpmodels.TypeOAuthTokens{}, err
		}

		return oauth2_ExchangeAuthCodeForOAuthTokens(clientConfig, redirectURIInfo, userContext)
	}

	impl.GetUserInfo = func(clientType, tenantId *string, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		clientConfig, err := impl.GetConfig(clientType, tenantId, input.Config, userContext)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		return oauth2_GetUserInfo(clientConfig, tenantId, oAuthTokens, userContext)
	}

	if input.Override != nil {
		impl = input.Override(impl)
	}

	return *impl
}
