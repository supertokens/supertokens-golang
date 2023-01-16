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

package multitenancymodels

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type APIInterface struct {
	LoginMethodsGET *func(tenantId *string, clientType *string, options APIOptions, userContext supertokens.UserContext) (LoginMethodsGETResponse, error)
}

type LoginMethodsGETResponse struct {
	OK           *TypeLoginMethods
	GeneralError *supertokens.GeneralErrorResponse
}

type TypeLoginMethods struct {
	EmailPassword TypeEmailPassword `json:"emailPassword"`
	Passwordless  TypePasswordless  `json:"passwordless"`
	ThirdParty    TypeThirdParty    `json:"thirdParty"`
}

type TypeEmailPassword struct {
	Enabled bool `json:"enabled"`
}

type TypePasswordless struct {
	Enabled bool `json:"enabled"`
}

type TypeThirdParty struct {
	Enabled   bool                     `json:"enabled"`
	Providers []TypeThirdPartyProvider `json:"providers"`
}

type TypeThirdPartyProvider struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type APIOptions struct {
	RecipeImplementation      RecipeInterface
	Config                    TypeNormalisedInput
	RecipeID                  string
	Req                       *http.Request
	Res                       http.ResponseWriter
	OtherHandler              http.HandlerFunc
	StaticThirdPartyProviders []tpmodels.ProviderInput
}
