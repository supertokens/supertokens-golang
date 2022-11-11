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

	impl.GetProviderConfig = func(tenantId *string, userContext supertokens.UserContext) (tpmodels.ProviderConfig, error) {
		if tenantId == nil {
			return input.Config, nil
		}

		// TODO impl
		return tpmodels.ProviderConfig{}, errors.New("not implemented")
	}

	impl.GetConfig = func(clientType *string, inputConfig tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClient, error) {
		if clientType == nil {
			if len(inputConfig.Clients) == 0 || len(inputConfig.Clients) > 1 {
				return tpmodels.ProviderConfigForClient{}, errors.New("please provide exactly one client config or pass clientType or tenantId")
			}

			return getProviderConfigForClient(inputConfig, inputConfig.Clients[0]), nil
		}

		for _, client := range inputConfig.Clients {
			if client.ClientType == *clientType {
				config := getProviderConfigForClient(input.Config, client)
				return config, nil
			}
		}

		return tpmodels.ProviderConfigForClient{}, errors.New("Could not find client config for clientType: " + *clientType)
	}

	impl.GetAuthorisationRedirectURL = func(config tpmodels.ProviderConfigForClient, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		return oauth2_GetAuthorisationRedirectURL(config, redirectURIOnProviderDashboard, userContext)
	}

	impl.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.ProviderConfigForClient, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		return oauth2_ExchangeAuthCodeForOAuthTokens(config, redirectURIInfo, userContext)
	}

	impl.GetUserInfo = func(config tpmodels.ProviderConfigForClient, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		return oauth2_GetUserInfo(config, oAuthTokens, userContext)
	}

	if input.Override != nil {
		impl = input.Override(impl)
	}

	return *impl
}
