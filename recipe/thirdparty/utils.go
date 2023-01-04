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

package thirdparty

import (
	"encoding/json"
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance *Recipe, appInfo supertokens.NormalisedAppinfo, config *tpmodels.TypeInput) (tpmodels.TypeNormalisedInput, error) {
	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance)

	signInAndUpFeature, err := validateAndNormaliseSignInAndUpConfig(config.SignInAndUpFeature)
	if err != nil {
		return tpmodels.TypeNormalisedInput{}, err
	}
	typeNormalisedInput.SignInAndUpFeature = signInAndUpFeature

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
	}

	return typeNormalisedInput, nil
}

func makeTypeNormalisedInput(recipeInstance *Recipe) tpmodels.TypeNormalisedInput {
	return tpmodels.TypeNormalisedInput{
		SignInAndUpFeature: tpmodels.TypeNormalisedInputSignInAndUp{},
		Override: tpmodels.OverrideStruct{
			Functions: func(originalImplementation tpmodels.RecipeInterface) tpmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation tpmodels.APIInterface) tpmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}

func validateAndNormaliseSignInAndUpConfig(config tpmodels.TypeInputSignInAndUp) (tpmodels.TypeNormalisedInputSignInAndUp, error) {
	providers := config.Providers
	// With the feature of dynamic providers from core, we need not check if the provider array is empty

	thirdPartyIdSet := map[string]bool{}

	for _, provider := range providers {
		if thirdPartyIdSet[provider.Config.ThirdPartyId] {
			return tpmodels.TypeNormalisedInputSignInAndUp{}, supertokens.BadInputError{Msg: "The providers array has multiple entries for the same third party provider."}
		}
		thirdPartyIdSet[provider.Config.ThirdPartyId] = true
	}

	// TODO normalize provider names

	normalisedProviders := []tpmodels.ProviderInput{}
	for _, provider := range providers {
		provider, err := validateAndNormaliseProviderInput(provider)
		if err != nil {
			return tpmodels.TypeNormalisedInputSignInAndUp{}, err
		}
		normalisedProviders = append(normalisedProviders, provider)
	}

	return tpmodels.TypeNormalisedInputSignInAndUp{
		Providers: normalisedProviders,
	}, nil
}

func validateAndNormaliseProviderInput(providerInput tpmodels.ProviderInput) (tpmodels.ProviderInput, error) {
	switch providerInput.Config.ThirdPartyId {
	case "active-directory":
		return providers.ValidateAndNormaliseActiveDirectory(providerInput)
	case "apple":
		return providers.ValidateAndNormaliseApple(providerInput)
	case "discord":
		return providers.ValidateAndNormaliseDiscord(providerInput)
	case "facebook":
		return providers.ValidateAndNormaliseFacebook(providerInput)
	case "github":
		return providers.ValidateAndNormaliseGithub(providerInput)
	case "google":
		return providers.ValidateAndNormaliseGoogle(providerInput)
	case "google-workspaces":
		return providers.ValidateAndNormaliseGoogleWorkspaces(providerInput)
	case "okta":
		return providers.ValidateAndNormaliseOkta(providerInput)
	case "linkedin":
		return providers.ValidateAndNormaliseLinkedin(providerInput)
	case "boxy-saml":
		return providers.ValidateAndNormaliseBoxySaml(providerInput)
	}

	return providers.ValidateAndNormaliseNewProvider(providerInput)
}

func parseUser(value interface{}) (*tpmodels.User, error) {
	respJSON, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var user tpmodels.User
	err = json.Unmarshal(respJSON, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func parseUsers(value interface{}) ([]tpmodels.User, error) {
	respJSON, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var user []tpmodels.User
	err = json.Unmarshal(respJSON, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func findAndCreateProviderInstance(providers []tpmodels.ProviderInput, thirdPartyId string, tenantId *string) (tpmodels.TypeProvider, error) {
	for _, provider := range providers {
		if provider.Config.ThirdPartyId == thirdPartyId {
			providerInstance := createProvider(provider)
			return *providerInstance, nil
		}
	}
	return tpmodels.TypeProvider{}, fmt.Errorf("the provider %s could not be found in the configuration", thirdPartyId)
}

func createProvider(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	switch input.Config.ThirdPartyId {
	case "active-directory":
		return providers.ActiveDirectory(input)
	case "apple":
		return providers.Apple(input)
	case "discord":
		return providers.Discord(input)
	case "facebook":
		return providers.Facebook(input)
	case "github":
		return providers.Github(input)
	case "google":
		return providers.Google(input)
	case "google-workspaces":
		return providers.GoogleWorkspaces(input)
	case "okta":
		return providers.Okta(input)
	case "linkedin":
		return providers.Linkedin(input)
	case "boxy-saml":
		return providers.BoxySaml(input)
	}

	return providers.NewProvider(input)
}

func mergeConfig(staticConfig tpmodels.ProviderConfig, coreConfig tpmodels.ProviderConfig) tpmodels.ProviderConfig {
	result := staticConfig

	if coreConfig.AuthorizationEndpoint != "" {
		result.AuthorizationEndpoint = coreConfig.AuthorizationEndpoint
	}
	if coreConfig.AuthorizationEndpointQueryParams != nil {
		result.AuthorizationEndpointQueryParams = coreConfig.AuthorizationEndpointQueryParams
	}
	if coreConfig.TokenEndpoint != "" {
		result.TokenEndpoint = coreConfig.TokenEndpoint
	}
	if coreConfig.TokenEndpointBodyParams != nil {
		result.TokenEndpointBodyParams = coreConfig.TokenEndpointBodyParams
	}
	if coreConfig.UserInfoEndpoint != "" {
		result.UserInfoEndpoint = coreConfig.UserInfoEndpoint
	}
	if coreConfig.UserInfoEndpointHeaders != nil {
		result.UserInfoEndpointHeaders = coreConfig.UserInfoEndpointHeaders
	}
	if coreConfig.UserInfoEndpointQueryParams != nil {
		result.UserInfoEndpointQueryParams = coreConfig.UserInfoEndpointQueryParams
	}
	if coreConfig.JwksURI != "" {
		result.JwksURI = coreConfig.JwksURI
	}
	if coreConfig.Name != "" {
		result.Name = coreConfig.Name
	}
	if coreConfig.OIDCDiscoveryEndpoint != "" {
		result.OIDCDiscoveryEndpoint = coreConfig.OIDCDiscoveryEndpoint
	}
	if coreConfig.UserInfoMap.FromIdTokenPayload.Email != "" {
		result.UserInfoMap.FromIdTokenPayload.Email = coreConfig.UserInfoMap.FromIdTokenPayload.Email
	}
	if coreConfig.UserInfoMap.FromIdTokenPayload.EmailVerified != "" {
		result.UserInfoMap.FromIdTokenPayload.EmailVerified = coreConfig.UserInfoMap.FromIdTokenPayload.EmailVerified
	}
	if coreConfig.UserInfoMap.FromIdTokenPayload.UserId != "" {
		result.UserInfoMap.FromIdTokenPayload.UserId = coreConfig.UserInfoMap.FromIdTokenPayload.UserId
	}
	if coreConfig.UserInfoMap.FromUserInfoAPI.Email != "" {
		result.UserInfoMap.FromUserInfoAPI.Email = coreConfig.UserInfoMap.FromUserInfoAPI.Email
	}
	if coreConfig.UserInfoMap.FromUserInfoAPI.EmailVerified != "" {
		result.UserInfoMap.FromUserInfoAPI.EmailVerified = coreConfig.UserInfoMap.FromUserInfoAPI.EmailVerified
	}
	if coreConfig.UserInfoMap.FromUserInfoAPI.UserId != "" {
		result.UserInfoMap.FromUserInfoAPI.UserId = coreConfig.UserInfoMap.FromUserInfoAPI.UserId
	}

	// Merge the clients
	mergedClients := append([]tpmodels.ProviderClientConfig{}, staticConfig.Clients...)
	for _, client := range coreConfig.Clients {
		for i, staticClient := range mergedClients {
			if staticClient.ClientType == client.ClientType {
				mergedClients[i] = client
				break
			}
		}
	}
	result.Clients = mergedClients

	return result
}
