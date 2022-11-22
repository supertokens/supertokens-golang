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

package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() tpmodels.APIInterface {

	providersForTenantGET := func(tenantId *string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.ProvidersForTenantGetResponse, error) {
		providers := []struct {
			ID   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{}

		configsFromCore, err := (*options.RecipeImplementation.ListConfigMappingsForTenant)(tenantId, userContext)
		if err != nil {
			return tpmodels.ProvidersForTenantGetResponse{}, err
		}

		// for default tenant = merge core and static config
		if tenantId == nil || *tenantId == tpmodels.DefaultTenantId {
			addedFromCore := map[string]bool{}

			for _, configFromCore := range configsFromCore.OK.Configs {
				providerResult := struct {
					ID   string `json:"id"`
					Name string `json:"name,omitempty"`
				}{
					ID:   configFromCore.ThirdPartyId,
					Name: configFromCore.Config.Name,
				}
				providers = append(providers, providerResult)
				addedFromCore[configFromCore.ThirdPartyId] = true
			}

			for _, staticProvider := range options.Providers {
				if staticProvider.UseForDefaultTenant {
					providerResult := struct {
						ID   string `json:"id"`
						Name string `json:"name,omitempty"`
					}{
						ID: staticProvider.ID,
					}

					if !addedFromCore[staticProvider.ID] {
						providers = append(providers, providerResult)
					}
				}
			}
		} else {
			// for other tenants = only core config if available else static config
			if len(configsFromCore.OK.Configs) > 0 {
				// Add from core
				for _, configFromCore := range configsFromCore.OK.Configs {
					providerResult := struct {
						ID   string `json:"id"`
						Name string `json:"name,omitempty"`
					}{
						ID:   configFromCore.ThirdPartyId,
						Name: configFromCore.Config.Name,
					}
					providers = append(providers, providerResult)
				}
			} else {
				// add from static
				for _, staticProvider := range options.Providers {
					providerResult := struct {
						ID   string `json:"id"`
						Name string `json:"name,omitempty"`
					}{
						ID: staticProvider.ID,
					}
					providers = append(providers, providerResult)
				}
			}
		}

		return tpmodels.ProvidersForTenantGetResponse{
			OK: &struct {
				Providers []struct {
					ID   string `json:"id"`
					Name string `json:"name,omitempty"`
				} `json:"providers"`
			}{
				Providers: providers,
			},
		}, nil
	}

	authorisationUrlGET := func(provider tpmodels.TypeProvider, config tpmodels.ProviderConfigForClientType, redirectURIOnProviderDashboard string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error) {
		authRedirect, err := provider.GetAuthorisationRedirectURL(config, redirectURIOnProviderDashboard, userContext)
		if err != nil {
			return tpmodels.AuthorisationUrlGETResponse{}, err
		}

		return tpmodels.AuthorisationUrlGETResponse{
			OK: &authRedirect,
		}, nil
	}

	signInUpPOST := func(provider tpmodels.TypeProvider, config tpmodels.ProviderConfigForClientType, input tpmodels.TypeSignInUpInput, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.SignInUpPOSTResponse, error) {
		var oAuthTokens map[string]interface{} = nil
		var err error

		if input.RedirectURIInfo != nil {
			oAuthTokens, err = provider.ExchangeAuthCodeForOAuthTokens(config, *input.RedirectURIInfo, userContext)
			if err != nil {
				return tpmodels.SignInUpPOSTResponse{}, err
			}
		} else {
			oAuthTokens = *input.OAuthTokens
		}

		userInfo, err := provider.GetUserInfo(config, oAuthTokens, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}

		emailInfo := userInfo.Email
		if emailInfo == nil {
			return tpmodels.SignInUpPOSTResponse{
				NoEmailGivenByProviderError: &struct{}{},
			}, nil
		}

		response, err := (*options.RecipeImplementation.SignInUp)(provider.ID, userInfo.ThirdPartyUserId, emailInfo.ID, oAuthTokens, userInfo.RawUserInfoFromProvider, config.TenantId, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}

		if emailInfo.IsVerified {
			evInstance := emailverification.GetRecipeInstance()
			if evInstance != nil {
				tokenResponse, err := (*evInstance.RecipeImpl.CreateEmailVerificationToken)(response.OK.User.ID, response.OK.User.Email, userContext)
				if err != nil {
					return tpmodels.SignInUpPOSTResponse{}, err
				}
				if tokenResponse.OK != nil {
					_, err := (*evInstance.RecipeImpl.VerifyEmailUsingToken)(tokenResponse.OK.Token, userContext)
					if err != nil {
						return tpmodels.SignInUpPOSTResponse{}, err
					}
				}
			}
		}

		session, err := session.CreateNewSessionWithContext(options.Res, response.OK.User.ID, nil, nil, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}
		return tpmodels.SignInUpPOSTResponse{
			OK: &struct {
				CreatedNewUser          bool
				User                    tpmodels.User
				Session                 sessmodels.SessionContainer
				OAuthTokens             tpmodels.TypeOAuthTokens
				RawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider
			}{
				CreatedNewUser:          response.OK.CreatedNewUser,
				User:                    response.OK.User,
				Session:                 session,
				OAuthTokens:             oAuthTokens,
				RawUserInfoFromProvider: userInfo.RawUserInfoFromProvider,
			},
		}, nil
	}

	appleRedirectHandlerPOST := func(formPostInfoFromProvider map[string]interface{}, options tpmodels.APIOptions, userContext supertokens.UserContext) error {
		state := formPostInfoFromProvider["state"].(string)
		stateBytes, err := base64.RawStdEncoding.DecodeString(state)

		if err != nil {
			return err
		}

		stateObj := map[string]interface{}{}
		err = json.Unmarshal(stateBytes, &stateObj)
		if err != nil {
			return err
		}

		redirectURL := stateObj["redirectURI"].(string)
		parsedRedirectURL, err := url.Parse(redirectURL)
		if err != nil {
			return err
		}

		query := parsedRedirectURL.Query()

		for k, v := range formPostInfoFromProvider {
			query.Add(k, fmt.Sprint(v))
		}

		parsedRedirectURL.RawQuery = query.Encode()

		options.Res.Header().Set("Location", parsedRedirectURL.String())
		options.Res.WriteHeader(http.StatusSeeOther)

		return nil
	}

	return tpmodels.APIInterface{
		ProvidersForTenantGET:    &providersForTenantGET,
		AuthorisationUrlGET:      &authorisationUrlGET,
		SignInUpPOST:             &signInUpPOST,
		AppleRedirectHandlerPOST: &appleRedirectHandlerPOST,
	}
}
