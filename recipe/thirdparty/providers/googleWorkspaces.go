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

// import (
// 	"errors"

// 	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
// 	"github.com/supertokens/supertokens-golang/supertokens"
// )

// const googleWorkspacesID = "google-workspaces"

// type GoogleWorkspacesConfig struct {
// 	ClientID     string
// 	ClientSecret string
// 	Scope        []string
// 	Domain       string

// 	AuthorizationEndpoint            string
// 	AuthorizationEndpointQueryParams map[string]interface{}

// 	TokenEndpoint string
// 	TokenParams   map[string]interface{}

// 	UserInfoEndpoint string

// 	JwksURI      string
// 	OIDCEndpoint string

// 	GetSupertokensUserInfoFromRawUserInfoResponse func(rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider, userContext supertokens.UserContext) (tpmodels.TypeSupertokensUserInfo, error)

// 	AdditionalConfig map[string]interface{}
// }

// type TypeGoogleWorkspacesInput struct {
// 	Config   []GoogleWorkspacesConfig
// 	Override func(provider *GoogleWorkspacesProvider) *GoogleWorkspacesProvider
// }

// type GoogleWorkspacesProvider struct {
// 	GetConfig func(id *tpmodels.TypeID, userContext supertokens.UserContext) (GoogleWorkspacesConfig, error)
// 	*tpmodels.TypeProvider
// }

// func GoogleWorkspaces(input TypeGoogleWorkspacesInput) tpmodels.TypeProvider {
// 	googleWorkspacesProvider := &GoogleWorkspacesProvider{
// 		TypeProvider: &tpmodels.TypeProvider{
// 			ID: googleWorkspacesID,
// 		},
// 	}

// 	var oAuth2ProviderConfig []OAuth2ProviderConfig
// 	if input.Config != nil {
// 		oAuth2ProviderConfig = make([]OAuth2ProviderConfig, len(input.Config))
// 		for idx, config := range input.Config {
// 			oAuth2ProviderConfig[idx] = googleWorkspacesConfigToOAuth2ProviderConfig(config)
// 		}
// 	}

// 	oAuth2Provider := oAuth2Provider(TypeOAuth2ProviderInput{
// 		ThirdPartyID: googleWorkspacesID,
// 		Config:       oAuth2ProviderConfig,
// 	})

// 	{
// 		// OAuth2 provider needs to use the config returned by google provider GetConfig
// 		// Also, google provider needs to use the default implementation of GetConfig provided by oAuth2 provider
// 		oGetConfig := oAuth2Provider.GetConfig
// 		oAuth2Provider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (OAuth2ProviderConfig, error) {
// 			config, err := googleWorkspacesProvider.GetConfig(id, userContext)
// 			if err != nil {
// 				return OAuth2ProviderConfig{}, err
// 			}
// 			return googleWorkspacesConfigToOAuth2ProviderConfig(config), nil
// 		}
// 		googleWorkspacesProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (GoogleWorkspacesConfig, error) {
// 			config, err := oGetConfig(id, userContext)
// 			if err != nil {
// 				return GoogleWorkspacesConfig{}, err
// 			}
// 			return googleWorkspacesConfigFromOAuth2ProviderConfig(config), nil
// 		}
// 	}

// 	{
// 		// GoogleWorkspaces provider APIs call into oAuth2 provider APIs

// 		googleWorkspacesProvider.GetAuthorisationRedirectURL = func(id *tpmodels.TypeID, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
// 			return oAuth2Provider.GetAuthorisationRedirectURL(id, redirectURIOnProviderDashboard, userContext)
// 		}

// 		googleWorkspacesProvider.ExchangeAuthCodeForOAuthTokens = func(id *tpmodels.TypeID, redirectInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
// 			return oAuth2Provider.ExchangeAuthCodeForOAuthTokens(id, redirectInfo, userContext)
// 		}

// 		googleWorkspacesProvider.GetUserInfo = func(id *tpmodels.TypeID, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
// 			return oAuth2Provider.GetUserInfo(id, oAuthTokens, userContext)
// 		}
// 	}

// 	if input.Override != nil {
// 		googleWorkspacesProvider = input.Override(googleWorkspacesProvider)
// 	}

// 	{
// 		// We want to always normalize (for google) the config before returning it
// 		oGetConfig := googleWorkspacesProvider.GetConfig
// 		googleWorkspacesProvider.GetConfig = func(id *tpmodels.TypeID, userContext supertokens.UserContext) (GoogleWorkspacesConfig, error) {
// 			config, err := oGetConfig(id, userContext)
// 			if err != nil {
// 				return GoogleWorkspacesConfig{}, err
// 			}
// 			return normalizeGoogleWorkspacesConfig(config), nil
// 		}
// 	}

// 	return *googleWorkspacesProvider.TypeProvider
// }

// func normalizeGoogleWorkspacesConfig(config GoogleWorkspacesConfig) GoogleWorkspacesConfig {
// 	if config.OIDCEndpoint == "" {
// 		config.OIDCEndpoint = "https://accounts.google.com"
// 	}

// 	if config.Domain == "" {
// 		config.Domain = "*"
// 	}

// 	if config.AuthorizationEndpointQueryParams == nil {
// 		accessType := "offline"
// 		if config.ClientSecret == "" {
// 			accessType = "online"
// 		}
// 		config.AuthorizationEndpointQueryParams = map[string]interface{}{
// 			"access_type":            accessType,
// 			"include_granted_scopes": "true",
// 			"response_type":          "code",
// 			"hd":                     config.Domain,
// 		}
// 	}

// 	if len(config.Scope) == 0 {
// 		config.Scope = []string{"openid", "https://www.googleapis.com/auth/userinfo.email"}
// 	}

// 	if config.GetSupertokensUserInfoFromRawUserInfoResponse == nil {
// 		config.GetSupertokensUserInfoFromRawUserInfoResponse = func(rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider, userContext supertokens.UserContext) (tpmodels.TypeSupertokensUserInfo, error) {
// 			result, err := getSupertokensUserInfoFromRawUserInfo("sub", "email", "email_verified", "id_token")(rawUserInfoResponse, userContext)
// 			if err != nil {
// 				return tpmodels.TypeSupertokensUserInfo{}, err
// 			}
// 			if config.Domain != "*" {
// 				if rawUserInfoResponse.FromIdToken["hd"] != config.Domain {
// 					return tpmodels.TypeSupertokensUserInfo{}, errors.New("Please use emails from " + config.Domain + " to login")
// 				}
// 			}
// 			return result, nil
// 		}
// 	}

// 	return config
// }

// func googleWorkspacesConfigToOAuth2ProviderConfig(googleWorkspacesConfig GoogleWorkspacesConfig) OAuth2ProviderConfig {
// 	additionalConfig := map[string]interface{}{}

// 	for k, v := range googleWorkspacesConfig.AdditionalConfig {
// 		additionalConfig[k] = v
// 	}
// 	additionalConfig["_domain"] = googleWorkspacesConfig.Domain

// 	return OAuth2ProviderConfig{
// 		ClientID:     googleWorkspacesConfig.ClientID,
// 		ClientSecret: googleWorkspacesConfig.ClientSecret,
// 		Scope:        googleWorkspacesConfig.Scope,

// 		AuthorizationEndpoint:            googleWorkspacesConfig.AuthorizationEndpoint,
// 		AuthorizationEndpointQueryParams: googleWorkspacesConfig.AuthorizationEndpointQueryParams,

// 		TokenEndpoint: googleWorkspacesConfig.TokenEndpoint,
// 		TokenParams:   googleWorkspacesConfig.TokenParams,

// 		UserInfoEndpoint: googleWorkspacesConfig.UserInfoEndpoint,

// 		JwksURI:      googleWorkspacesConfig.JwksURI,
// 		OIDCEndpoint: googleWorkspacesConfig.OIDCEndpoint,

// 		GetSupertokensUserInfoFromRawUserInfoResponse: googleWorkspacesConfig.GetSupertokensUserInfoFromRawUserInfoResponse,

// 		AdditionalConfig: additionalConfig,
// 	}
// }

// func googleWorkspacesConfigFromOAuth2ProviderConfig(config OAuth2ProviderConfig) GoogleWorkspacesConfig {
// 	return GoogleWorkspacesConfig{
// 		ClientID:     config.ClientID,
// 		ClientSecret: config.ClientSecret,
// 		Scope:        config.Scope,
// 		Domain:       config.AdditionalConfig["_domain"].(string),

// 		AuthorizationEndpoint:            config.AuthorizationEndpoint,
// 		AuthorizationEndpointQueryParams: config.AuthorizationEndpointQueryParams,

// 		TokenEndpoint: config.TokenEndpoint,
// 		TokenParams:   config.TokenParams,

// 		UserInfoEndpoint: config.UserInfoEndpoint,

// 		JwksURI:      config.JwksURI,
// 		OIDCEndpoint: config.OIDCEndpoint,

// 		GetSupertokensUserInfoFromRawUserInfoResponse: config.GetSupertokensUserInfoFromRawUserInfoResponse,

// 		AdditionalConfig: config.AdditionalConfig,
// 	}
// }
