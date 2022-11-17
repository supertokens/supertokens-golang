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
		input.Config.UserInfoMap.FromIdTokenPayload.EmailVerified = "email_verified"
	}

	if input.UseForDefaultTenant == nil {
		True := true
		input.UseForDefaultTenant = &True
	}

	impl.UseForDefaultTenant = *input.UseForDefaultTenant

	impl.GetAllClientTypeConfigForTenant = func(tenantId string, recipeImpl tpmodels.RecipeInterface, userContext supertokens.UserContext) (tpmodels.ProviderConfig, error) {
		input.Config.TenantId = tenantId

		configFromCore, err := (*recipeImpl.FetchTenantIdConfigMapping)(input.ThirdPartyID, tenantId, userContext)
		if err != nil {
			return tpmodels.ProviderConfig{}, err
		}

		if configFromCore.UnknownMappingError != nil {
			// Return the static config as core doesn't have a mapping
			return input.Config, err
		}

		configToReturn := input.Config

		if tenantId == tpmodels.DefaultTenantId {
			// if tenantId is default, merge the client configs
			staticClientConfigs := []tpmodels.ProviderClientConfig{}
			copy(staticClientConfigs, input.Config.Clients)

			for _, clientConfigFromCore := range configFromCore.OK.Config.Clients {
				found := false
				for i, staticClientConfig := range staticClientConfigs {
					if clientConfigFromCore.ClientType == staticClientConfig.ClientType {
						staticClientConfigs[i] = clientConfigFromCore
						found = true
						break
					}
				}
				if !found {
					staticClientConfigs = append(staticClientConfigs, clientConfigFromCore)
				}
			}

		} else {
			// We use the client configs from the db and ignore the static
			configToReturn.Clients = configFromCore.OK.Config.Clients
		}

		// Merge the config from Core
		if configFromCore.OK.Config.AuthorizationEndpoint != "" {
			configToReturn.AuthorizationEndpoint = configFromCore.OK.Config.AuthorizationEndpoint
		}
		if configFromCore.OK.Config.AuthorizationEndpointQueryParams != nil {
			configToReturn.AuthorizationEndpointQueryParams = configFromCore.OK.Config.AuthorizationEndpointQueryParams
		}
		if configFromCore.OK.Config.TokenEndpoint != "" {
			configToReturn.TokenEndpoint = configFromCore.OK.Config.TokenEndpoint
		}
		if configFromCore.OK.Config.TokenEndpointBodyParams != nil {
			configToReturn.TokenEndpointBodyParams = configFromCore.OK.Config.TokenEndpointBodyParams
		}
		if configFromCore.OK.Config.ForcePKCE != nil {
			configToReturn.ForcePKCE = configFromCore.OK.Config.ForcePKCE
		}
		if configFromCore.OK.Config.UserInfoEndpoint != "" {
			configToReturn.UserInfoEndpoint = configFromCore.OK.Config.UserInfoEndpoint
		}
		if configFromCore.OK.Config.UserInfoEndpointQueryParams != nil {
			configToReturn.UserInfoEndpointQueryParams = configFromCore.OK.Config.UserInfoEndpointQueryParams
		}
		if configFromCore.OK.Config.UserInfoEndpointHeaders != nil {
			configToReturn.UserInfoEndpointHeaders = configFromCore.OK.Config.UserInfoEndpointHeaders
		}
		if configFromCore.OK.Config.JwksURI != "" {
			configToReturn.JwksURI = configFromCore.OK.Config.JwksURI
		}
		if configFromCore.OK.Config.OIDCDiscoveryEndpoint != "" {
			configToReturn.OIDCDiscoveryEndpoint = configFromCore.OK.Config.OIDCDiscoveryEndpoint
		}
		if configFromCore.OK.Config.TenantId != "" {
			configToReturn.TenantId = configFromCore.OK.Config.TenantId
		}
		if configFromCore.OK.Config.UserInfoMap.FromIdTokenPayload.UserId != "" {
			configToReturn.UserInfoMap.FromIdTokenPayload.UserId = configFromCore.OK.Config.UserInfoMap.FromIdTokenPayload.UserId
		}
		if configFromCore.OK.Config.UserInfoMap.FromIdTokenPayload.Email != "" {
			configToReturn.UserInfoMap.FromIdTokenPayload.Email = configFromCore.OK.Config.UserInfoMap.FromIdTokenPayload.Email
		}
		if configFromCore.OK.Config.UserInfoMap.FromIdTokenPayload.EmailVerified != "" {
			configToReturn.UserInfoMap.FromIdTokenPayload.EmailVerified = configFromCore.OK.Config.UserInfoMap.FromIdTokenPayload.EmailVerified
		}

		if configFromCore.OK.Config.UserInfoMap.FromUserInfoAPI.UserId != "" {
			configToReturn.UserInfoMap.FromUserInfoAPI.UserId = configFromCore.OK.Config.UserInfoMap.FromUserInfoAPI.UserId
		}
		if configFromCore.OK.Config.UserInfoMap.FromUserInfoAPI.Email != "" {
			configToReturn.UserInfoMap.FromUserInfoAPI.Email = configFromCore.OK.Config.UserInfoMap.FromUserInfoAPI.Email
		}
		if configFromCore.OK.Config.UserInfoMap.FromUserInfoAPI.EmailVerified != "" {
			configToReturn.UserInfoMap.FromUserInfoAPI.EmailVerified = configFromCore.OK.Config.UserInfoMap.FromUserInfoAPI.EmailVerified
		}

		if configFromCore.OK.Config.Name != "" {
			configToReturn.Name = configFromCore.OK.Config.Name
		}

		return configToReturn, nil
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
