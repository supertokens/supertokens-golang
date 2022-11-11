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

package tpmodels

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type APIInterface struct {
	AuthorisationUrlGET *func(provider TypeProvider, config ProviderConfigForClient, redirectURIOnProviderDashboard string, options APIOptions, userContext supertokens.UserContext) (AuthorisationUrlGETResponse, error)
	SignInUpPOST        *func(provider TypeProvider, config ProviderConfigForClient, input TypeSignInUpInput, options APIOptions, userContext supertokens.UserContext) (SignInUpPOSTResponse, error)

	AppleRedirectHandlerPOST *func(formPostInfoFromProvider map[string]interface{}, options APIOptions, userContext supertokens.UserContext) error
	ProvidersForTenantGET    *func(tenantId string, userContext supertokens.UserContext) (ProvidersForTenantGetResponse, error)
}

type ProvidersForTenantGetResponse struct {
	OK *struct {
		Providers []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"providers"`
	}
	GeneralError *supertokens.GeneralErrorResponse
}

type AuthorisationUrlGETResponse struct {
	OK           *TypeAuthorisationRedirect
	GeneralError *supertokens.GeneralErrorResponse
}

type TypeSignInUpInput struct {
	// Either of the below
	RedirectURIInfo *TypeRedirectURIInfo `json:"redirectURIInfo"`
	OAuthTokens     *TypeOAuthTokens     `json:"oAuthTokens"`
}

type SignInUpPOSTResponse struct {
	OK *struct {
		CreatedNewUser          bool
		User                    User
		Session                 sessmodels.SessionContainer
		OAuthTokens             TypeOAuthTokens
		RawUserInfoFromProvider TypeRawUserInfoFromProvider
	}
	NoEmailGivenByProviderError *struct{}
	GeneralError                *supertokens.GeneralErrorResponse
}

type APIOptions struct {
	RecipeImplementation RecipeInterface
	Config               TypeNormalisedInput
	RecipeID             string
	Providers            []TypeProvider
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
	AppInfo              supertokens.NormalisedAppinfo
	EmailDelivery        emaildelivery.Ingredient
}
