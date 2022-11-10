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

const googleWorkspacesID = "google-workspaces"

type GoogleWorkspacesConfig struct {
	Clients []GoogleWorkspacesClientConfig

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	ForcePKCE                        bool
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      tpmodels.TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, config GoogleWorkspacesClientConfig) (bool, error)
}

func (config GoogleWorkspacesConfig) ToCustomConfig() CustomConfig {
	customConfig := CustomConfig{
		AuthorizationEndpoint:            config.AuthorizationEndpoint,
		AuthorizationEndpointQueryParams: config.AuthorizationEndpointQueryParams,
		TokenEndpoint:                    config.TokenEndpoint,
		TokenParams:                      config.TokenParams,
		ForcePKCE:                        config.ForcePKCE,
		UserInfoEndpoint:                 config.UserInfoEndpoint,
		JwksURI:                          config.JwksURI,
		OIDCDiscoveryEndpoint:            config.OIDCDiscoveryEndpoint,
		UserInfoMap:                      config.UserInfoMap,
	}

	if config.ValidateIdTokenPayload != nil {
		customConfig.ValidateIdTokenPayload = func(idTokenPayload map[string]interface{}, clientConfig CustomClientConfig) (bool, error) {
			googleWorkspacesClientConfig := GoogleWorkspacesClientConfig{}
			googleWorkspacesClientConfig.UpdateFromCustomClientConfig(clientConfig)
			return config.ValidateIdTokenPayload(idTokenPayload, googleWorkspacesClientConfig)
		}
	}

	customConfig.Clients = make([]CustomClientConfig, len(config.Clients))
	for i, client := range config.Clients {
		customConfig.Clients[i], _ = client.ToCustomClientConfig()
	}

	return customConfig
}

type GoogleWorkspacesClientConfig struct {
	ClientType       string
	ClientID         string
	ClientSecret     string
	Scope            []string
	Domain           string
	AdditionalConfig map[string]interface{}
}

func (config GoogleWorkspacesClientConfig) ToCustomClientConfig() (CustomClientConfig, error) {
	cConfig := CustomClientConfig{
		ClientID:         config.ClientID,
		ClientSecret:     config.ClientSecret,
		ClientType:       config.ClientType,
		Scope:            config.Scope,
		AdditionalConfig: config.AdditionalConfig,
	}
	if cConfig.AdditionalConfig == nil {
		cConfig.AdditionalConfig = map[string]interface{}{}
	}
	cConfig.AdditionalConfig["_domain"] = config.Domain
	return cConfig, nil
}

func (config *GoogleWorkspacesClientConfig) UpdateFromCustomClientConfig(input CustomClientConfig) {
	config.ClientType = input.ClientType
	config.ClientID = input.ClientID
	config.ClientSecret = input.ClientSecret
	config.Scope = input.Scope
	config.Domain = input.AdditionalConfig["_domain"].(string)
	config.AdditionalConfig = input.AdditionalConfig
}

type TypeGoogleWorkspaces struct {
	GetConfig func(clientType *string, tenantId *string, userContext supertokens.UserContext) (GoogleWorkspacesClientConfig, error)
	*tpmodels.TypeProvider
}

type GoogleWorkspaces struct {
	Config   GoogleWorkspacesConfig
	Override func(provider *TypeGoogleWorkspaces) *TypeGoogleWorkspaces
}

func (input GoogleWorkspaces) GetID() string {
	return googleWorkspacesID
}

func (input GoogleWorkspaces) Build() tpmodels.TypeProvider {
	googleWorkspacesImpl := input.buildInternal()
	if input.Override != nil {
		googleWorkspacesImpl = input.Override(googleWorkspacesImpl)
	}
	return *googleWorkspacesImpl.TypeProvider
}

func (input GoogleWorkspaces) buildInternal() *TypeGoogleWorkspaces {
	customProvider := (CustomProvider{
		ThirdPartyID: googleID,
		Config:       input.Config.ToCustomConfig(),

		oAuth2Normalize: normalizeOAuth2ConfigForGoogleWorkspaces,
	}).buildInternal()

	googleWorkspacesImpl := &TypeGoogleWorkspaces{
		TypeProvider: &tpmodels.TypeProvider{},
	}

	{
		oGetConfig := customProvider.GetConfig

		googleWorkspacesImpl.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (GoogleWorkspacesClientConfig, error) {
			customConfig, err := oGetConfig(clientType, tenantId, userContext)
			if err != nil {
				return GoogleWorkspacesClientConfig{}, err
			}
			googleWorkspacesConfig := &GoogleWorkspacesClientConfig{}
			googleWorkspacesConfig.UpdateFromCustomClientConfig(customConfig)

			return *googleWorkspacesConfig, nil
		}

		customProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (CustomClientConfig, error) {
			googleWorkspacesConfig, err := googleWorkspacesImpl.GetConfig(clientType, tenantId, userContext)
			if err != nil {
				return CustomClientConfig{}, err
			}
			return googleWorkspacesConfig.ToCustomClientConfig()
		}
	}

	googleWorkspacesImpl.GetAuthorisationRedirectURL = func(clientType, tenantId *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
		return customProvider.GetAuthorisationRedirectURL(clientType, tenantId, redirectURIOnProviderDashboard, userContext)
	}

	googleWorkspacesImpl.ExchangeAuthCodeForOAuthTokens = func(clientType, tenantId *string, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		return customProvider.ExchangeAuthCodeForOAuthTokens(clientType, tenantId, redirectInfo, userContext)
	}

	googleWorkspacesImpl.GetUserInfo = func(clientType, tenantId *string, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		return customProvider.GetUserInfo(clientType, tenantId, oAuthTokens, userContext)
	}

	return googleWorkspacesImpl
}

func normalizeOAuth2ConfigForGoogleWorkspaces(config *typeCombinedOAuth2Config) *typeCombinedOAuth2Config {
	if config.OIDCDiscoveryEndpoint == "" {
		config.OIDCDiscoveryEndpoint = "https://accounts.google.com"
	}

	if config.AuthorizationEndpointQueryParams == nil {
		accessType := "offline"
		if config.ClientSecret == "" {
			accessType = "online"
		}
		config.AuthorizationEndpointQueryParams = map[string]interface{}{
			"access_type":            accessType,
			"include_granted_scopes": "true",
			"response_type":          "code",
		}
		if config.AdditionalConfig["_domain"] == nil || config.AdditionalConfig["_domain"] == "" {
			config.AuthorizationEndpointQueryParams["hd"] = "*"
		} else {
			config.AuthorizationEndpointQueryParams["hd"] = config.AdditionalConfig["_domain"]
		}
	}

	if config.ValidateIdTokenPayload == nil {
		config.ValidateIdTokenPayload = func(idTokenPayload map[string]interface{}, config *typeCombinedOAuth2Config) (bool, error) {
			if config.AdditionalConfig["_domain"] == "*" {
				return true, nil
			}
			if idTokenPayload["hd"] == config.AdditionalConfig["_domain"] {
				return true, nil
			}
			return false, nil
		}
	}

	if len(config.Scope) == 0 {
		config.Scope = []string{"openid", "https://www.googleapis.com/auth/userinfo.email"}
	}

	if config.UserInfoMap.From == "" {
		config.UserInfoMap.From = tpmodels.FromIdTokenPayload
	}

	if config.UserInfoMap.UserId == "" {
		config.UserInfoMap.UserId = "sub"
	}

	if config.UserInfoMap.Email == "" {
		config.UserInfoMap.Email = "email"
	}

	if config.UserInfoMap.EmailVerified == "" {
		config.UserInfoMap.EmailVerified = "email_verified"
	}

	return config
}

// type GoogleWorkspacesConfig struct {
// 	Clients []GoogleWorkspacesClientConfig

// 	AuthorizationEndpoint            string
// 	AuthorizationEndpointQueryParams map[string]interface{}
// 	TokenEndpoint                    string
// 	TokenParams                      map[string]interface{}
// 	ForcePKCE                        *bool
// 	UserInfoEndpoint                 string
// 	JwksURI                          string
// 	OIDCDiscoveryEndpoint            string
// 	UserInfoMap                      tpmodels.TypeUserInfoMap
// 	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, config tpmodels.TypeNormalisedProviderConfig) (bool, error)
// }

// func (config GoogleWorkspacesConfig) ToOAuth2Config() OAuth2ProviderConfig {
// 	clients := make([]OAuth2ProviderClientConfig, len(config.Clients))
// 	for i, client := range config.Clients {
// 		clients[i] = OAuth2ProviderClientConfig{
// 			ClientType:       client.ClientType,
// 			ClientID:         client.ClientID,
// 			ClientSecret:     client.ClientSecret,
// 			Scope:            client.Scope,
// 			AdditionalConfig: client.AdditionalConfig,
// 		}
// 		if clients[i].AdditionalConfig == nil {
// 			clients[i].AdditionalConfig = map[string]interface{}{}
// 		}
// 		clients[i].AdditionalConfig["_domain"] = client.Domain
// 	}

// 	return OAuth2ProviderConfig{
// 		Clients:                          clients,
// 		AuthorizationEndpoint:            config.AuthorizationEndpoint,
// 		AuthorizationEndpointQueryParams: config.AuthorizationEndpointQueryParams,
// 		TokenEndpoint:                    config.TokenEndpoint,
// 		TokenParams:                      config.TokenParams,
// 		ForcePKCE:                        config.ForcePKCE,
// 		UserInfoEndpoint:                 config.UserInfoEndpoint,
// 		JwksURI:                          config.JwksURI,
// 		OIDCDiscoveryEndpoint:            config.OIDCDiscoveryEndpoint,
// 		UserInfoMap:                      config.UserInfoMap,
// 		ValidateIdTokenPayload:           config.ValidateIdTokenPayload,
// 	}
// }

// type GoogleWorkspacesClientConfig struct {
// 	ClientType       string
// 	ClientID         string
// 	ClientSecret     string
// 	Domain           string
// 	Scope            []string
// 	AdditionalConfig map[string]interface{}
// }

// type TypeGoogleWorkspacesInput struct {
// 	Config   GoogleWorkspacesConfig
// 	Override func(provider *GoogleWorkspacesProvider) *GoogleWorkspacesProvider
// }

// type GoogleWorkspacesProvider struct {
// 	*tpmodels.TypeProvider
// }

// func GoogleWorkspaces(input TypeGoogleWorkspacesInput) tpmodels.TypeProvider {
// 	googleWorkspacesProvider := &GoogleWorkspacesProvider{
// 		TypeProvider: &tpmodels.TypeProvider{
// 			ID: googleWorkspacesID,
// 		},
// 	}

// 	oAuth2Provider := oAuth2Provider(TypeOAuth2ProviderInput{
// 		ThirdPartyID: googleWorkspacesID,
// 		Config:       input.Config.ToOAuth2Config(),
// 	})

// 	{
// 		// GoogleWorkspaces provider APIs call into oAuth2 provider APIs

// 		googleWorkspacesProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
// 			return oAuth2Provider.GetConfig(clientType, tenantId, userContext)
// 		}

// 		googleWorkspacesProvider.GetAuthorisationRedirectURL = func(config tpmodels.TypeNormalisedProviderConfig, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
// 			return oAuth2Provider.GetAuthorisationRedirectURL(config, redirectURIOnProviderDashboard, userContext)
// 		}

// 		googleWorkspacesProvider.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.TypeNormalisedProviderConfig, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
// 			return oAuth2Provider.ExchangeAuthCodeForOAuthTokens(config, redirectInfo, userContext)
// 		}

// 		googleWorkspacesProvider.GetUserInfo = func(config tpmodels.TypeNormalisedProviderConfig, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
// 			return oAuth2Provider.GetUserInfo(config, oAuthTokens, userContext)
// 		}
// 	}

// 	if input.Override != nil {
// 		googleWorkspacesProvider = input.Override(googleWorkspacesProvider)
// 	}

// 	{
// 		// We want to always normalize (for googleWorkspaces) the config before returning it
// 		oGetConfig := googleWorkspacesProvider.GetConfig
// 		googleWorkspacesProvider.GetConfig = func(clientType, tenantId *string, userContext supertokens.UserContext) (tpmodels.TypeNormalisedProviderConfig, error) {
// 			config, err := oGetConfig(clientType, tenantId, userContext)
// 			if err != nil {
// 				return tpmodels.TypeNormalisedProviderConfig{}, err
// 			}
// 			return normalizeGoogleWorkspacesConfig(config), nil
// 		}
// 	}

// 	return *googleWorkspacesProvider.TypeProvider
// }

// func normalizeGoogleWorkspacesConfig(config tpmodels.TypeNormalisedProviderConfig) tpmodels.TypeNormalisedProviderConfig {

// 	if config.ValidateIdTokenPayload == nil {
// 		config.ValidateIdTokenPayload = func(idTokenPayload map[string]interface{}, config tpmodels.TypeNormalisedProviderConfig) (bool, error) {
// 			if config.AdditionalConfig["_domain"] == "*" {
// 				return true, nil
// 			}
// 			if idTokenPayload["hd"] == config.AdditionalConfig["_domain"] {
// 				return true, nil
// 			}
// 			return false, nil
// 		}
// 	}

// 	return config
// }
