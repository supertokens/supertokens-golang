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

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
)

type APIInterface struct {
	AuthorisationUrlGET *func(provider TypeProvider, options APIOptions) (AuthorisationUrlGETResponse, error)
	SignInUpPOST        *func(provider TypeProvider, code string, authCodeResponse interface{}, redirectURI string, options APIOptions) (SignInUpPOSTResponse, error)
}

type AuthorisationUrlGETResponse struct {
	OK *struct{ Url string }
}

type SignInUpPOSTResponse struct {
	OK *struct {
		CreatedNewUser   bool
		User             User
		AuthCodeResponse interface{}
	}
	NoEmailGivenByProviderError *struct{}
	FieldError                  *struct{ Error string }
}

type APIOptions struct {
	RecipeImplementation                  RecipeInterface
	EmailVerificationRecipeImplementation evmodels.RecipeInterface
	Config                                TypeNormalisedInput
	RecipeID                              string
	Providers                             []TypeProvider
	Req                                   *http.Request
	Res                                   http.ResponseWriter
	OtherHandler                          http.HandlerFunc
}
