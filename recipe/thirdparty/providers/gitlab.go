/* Copyright (c) 2023, VRAI Labs and/or its affiliates. All rights reserved.
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

func Gitlab(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.Name == "" {
		input.Config.Name = "Gitlab"
	}

	oOverride := input.Override

	input.Override = func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := originalImplementation.GetConfigForClientType
		originalImplementation.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"openid", "email"}
			}

			if config.AdditionalConfig != nil && config.AdditionalConfig["gitlabBaseUrl"] != nil {
				gitlabBaseUrl := config.AdditionalConfig["gitlabBaseUrl"].(string)
				oidcDomain, err := supertokens.NewNormalisedURLDomain(gitlabBaseUrl)
				if err != nil {
					return tpmodels.ProviderConfigForClientType{}, err
				}
				oidcPath, err := supertokens.NewNormalisedURLPath("/.well-known/openid-configuration")
				if err != nil {
					return tpmodels.ProviderConfigForClientType{}, err
				}
				config.OIDCDiscoveryEndpoint = oidcDomain.GetAsStringDangerous() + oidcPath.GetAsStringDangerous()
			} else if config.OIDCDiscoveryEndpoint == "" {
				config.OIDCDiscoveryEndpoint = "https://gitlab.com/.well-known/openid-configuration"
			}

			// The config could be coming from core where we didn't add the well-known previously
			config.OIDCDiscoveryEndpoint = normaliseOIDCEndpointToIncludeWellKnown(config.OIDCDiscoveryEndpoint)

			return config, nil
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}

		return originalImplementation
	}

	return NewProvider(input)
}
