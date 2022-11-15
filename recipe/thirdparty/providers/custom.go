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

	if input.Config.UserInfoMap.FromIdTokenPayload.UserId == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.UserId = "sub"
	}
	if input.Config.UserInfoMap.FromIdTokenPayload.Email == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.Email = "email"
	}
	if input.Config.UserInfoMap.FromIdTokenPayload.EmailVerified == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.Email = "email_verified"
	}

	impl.GetAllClientTypeConfigForTenant = func(tenantId *string, recipeImpl tpmodels.RecipeInterface, userContext supertokens.UserContext) (tpmodels.ProviderConfig, error) {
		if tenantId == nil {
			return input.Config, nil
		}

		input.Config.TenantId = *tenantId

		configs, err := (*recipeImpl.FetchTenantIdConfigMapping)(input.ThirdPartyID, *tenantId, userContext)
		if err != nil {
			return tpmodels.ProviderConfig{}, err
		}

		if configs.UnknownMappingError != nil {
			// Return the static config as core doesn't have a mapping
			return input.Config, err
		}

		// We use the client configs frpm the db and ignore the static
		clientConfigsFromDb := make([]tpmodels.ProviderClientConfig, len(configs.OK.Config.Clients))
		for i, config := range configs.OK.Config.Clients {
			clientConfigsFromDb[i] = tpmodels.ProviderClientConfig{
				ClientType:       config.ClientType,
				ClientID:         config.ClientID,
				ClientSecret:     config.ClientSecret,
				Scope:            config.Scope,
				AdditionalConfig: config.AdditionalConfig,
			}
		}

		// Merge the config from DB
		config := input.Config
		config.Clients = clientConfigsFromDb

		if configs.OK.Config.AuthorizationEndpoint != "" {
			config.AuthorizationEndpoint = configs.OK.Config.AuthorizationEndpoint
		}
		if configs.OK.Config.AuthorizationEndpointQueryParams != nil {
			config.AuthorizationEndpointQueryParams = configs.OK.Config.AuthorizationEndpointQueryParams
		}
		if configs.OK.Config.TokenEndpoint != "" {
			config.TokenEndpoint = configs.OK.Config.TokenEndpoint
		}
		if configs.OK.Config.TokenEndpointBodyParams != nil {
			config.TokenEndpointBodyParams = configs.OK.Config.TokenEndpointBodyParams
		}
		if configs.OK.Config.ForcePKCE != nil {
			config.ForcePKCE = configs.OK.Config.ForcePKCE
		}
		if configs.OK.Config.UserInfoEndpoint != "" {
			config.UserInfoEndpoint = configs.OK.Config.UserInfoEndpoint
		}
		if configs.OK.Config.UserInfoEndpointQueryParams != nil {
			config.UserInfoEndpointQueryParams = configs.OK.Config.UserInfoEndpointQueryParams
		}
		if configs.OK.Config.UserInfoEndpointHeaders != nil {
			config.UserInfoEndpointHeaders = configs.OK.Config.UserInfoEndpointHeaders
		}
		if configs.OK.Config.JwksURI != "" {
			config.JwksURI = configs.OK.Config.JwksURI
		}
		if configs.OK.Config.OIDCDiscoveryEndpoint != "" {
			config.OIDCDiscoveryEndpoint = configs.OK.Config.OIDCDiscoveryEndpoint
		}
		if configs.OK.Config.TenantId != "" {
			config.TenantId = configs.OK.Config.TenantId
		}
		if configs.OK.Config.UserInfoMap.FromIdTokenPayload.UserId != "" {
			config.UserInfoMap.FromIdTokenPayload.UserId = configs.OK.Config.UserInfoMap.FromIdTokenPayload.UserId
		}
		if configs.OK.Config.UserInfoMap.FromIdTokenPayload.Email != "" {
			config.UserInfoMap.FromIdTokenPayload.Email = configs.OK.Config.UserInfoMap.FromIdTokenPayload.Email
		}
		if configs.OK.Config.UserInfoMap.FromIdTokenPayload.EmailVerified != "" {
			config.UserInfoMap.FromIdTokenPayload.EmailVerified = configs.OK.Config.UserInfoMap.FromIdTokenPayload.EmailVerified
		}

		if configs.OK.Config.UserInfoMap.FromUserInfoAPI.UserId != "" {
			config.UserInfoMap.FromUserInfoAPI.UserId = configs.OK.Config.UserInfoMap.FromUserInfoAPI.UserId
		}
		if configs.OK.Config.UserInfoMap.FromUserInfoAPI.Email != "" {
			config.UserInfoMap.FromUserInfoAPI.Email = configs.OK.Config.UserInfoMap.FromUserInfoAPI.Email
		}
		if configs.OK.Config.UserInfoMap.FromUserInfoAPI.EmailVerified != "" {
			config.UserInfoMap.FromUserInfoAPI.EmailVerified = configs.OK.Config.UserInfoMap.FromUserInfoAPI.EmailVerified
		}

		if configs.OK.Config.Name != "" {
			config.Name = configs.OK.Config.Name
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
