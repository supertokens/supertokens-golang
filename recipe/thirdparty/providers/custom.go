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

	impl.GetAllClientTypeConfigForTenant = func(tenantId *string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.ProviderConfig, error) {
		if tenantId == nil {
			return input.Config, nil
		}

		input.Config.TenantId = *tenantId

		configs, err := (*options.RecipeImplementation.FetchTenantIdConfigMapping)(input.ThirdPartyID, *tenantId, userContext)
		if err != nil {
			return tpmodels.ProviderConfig{}, err
		}

		if configs.UnknownMappingError != nil {
			// Return the static config as core doesn't have a mapping
			return input.Config, err
		}

		clientConfigs := []tpmodels.ProviderClientConfig{}
		copy(clientConfigs, input.Config.Clients)

		for _, config := range configs.OK.Config.Clients {
			found := false
			for i := range clientConfigs {
				if clientConfigs[i].ClientType == config.ClientType {
					clientConfigs[i].ClientID = config.ClientID
					clientConfigs[i].ClientSecret = config.ClientSecret
					clientConfigs[i].Scope = config.Scope
					clientConfigs[i].AdditionalConfig = config.AdditionalConfig
					found = true
					break
				}
			}
			if !found {
				clientConfigs = append(clientConfigs, tpmodels.ProviderClientConfig{
					ClientType:       config.ClientType,
					ClientID:         config.ClientID,
					ClientSecret:     config.ClientSecret,
					Scope:            config.Scope,
					AdditionalConfig: config.AdditionalConfig,
				})
			}
		}

		// Copy provider config
		config := tpmodels.ProviderConfig{
			Clients:                          clientConfigs,
			AuthorizationEndpoint:            configs.OK.Config.AuthorizationEndpoint,
			AuthorizationEndpointQueryParams: configs.OK.Config.AuthorizationEndpointQueryParams,
			TokenEndpoint:                    configs.OK.Config.TokenEndpoint,
			TokenParams:                      configs.OK.Config.TokenParams,
			ForcePKCE:                        configs.OK.Config.ForcePKCE,
			UserInfoEndpoint:                 configs.OK.Config.UserInfoEndpoint,
			UserInfoEndpointQueryParams:      configs.OK.Config.UserInfoEndpointQueryParams,
			UserInfoEndpointHeaders:          configs.OK.Config.UserInfoEndpointHeaders,
			JwksURI:                          configs.OK.Config.JwksURI,
			OIDCDiscoveryEndpoint:            configs.OK.Config.OIDCDiscoveryEndpoint,
			UserInfoMap:                      configs.OK.Config.UserInfoMap,
			TenantId:                         *tenantId,

			ValidateIdTokenPayload: input.Config.ValidateIdTokenPayload, // We may want to use this from static config
		}

		return config, nil
	}

	impl.GetConfigForClientType = func(clientType *string, inputConfig tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
		if clientType == nil {
			if len(inputConfig.Clients) != 1 {
				return tpmodels.ProviderConfigForClientType{}, errors.New("please provide exactly one client config or pass clientType or tenantId")
			}

			return getProviderConfigForClient(inputConfig, inputConfig.Clients[0]), nil
		}

		for _, client := range inputConfig.Clients {
			if client.ClientType == *clientType {
				config := getProviderConfigForClient(input.Config, client)
				return config, nil
			}
		}

		return tpmodels.ProviderConfigForClientType{}, errors.New("Could not find client config for clientType: " + *clientType)
	}

	impl.GetAuthorisationRedirectURL = func(config tpmodels.ProviderConfigForClientType, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		return oauth2_GetAuthorisationRedirectURL(config, redirectURIOnProviderDashboard, userContext)
	}

	impl.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.ProviderConfigForClientType, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		return oauth2_ExchangeAuthCodeForOAuthTokens(config, redirectURIInfo, userContext)
	}

	impl.GetUserInfo = func(config tpmodels.ProviderConfigForClientType, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		return oauth2_GetUserInfo(config, oAuthTokens, userContext)
	}

	if input.Override != nil {
		impl = input.Override(impl)
	}

	return *impl
}
