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

	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeGetTenantIdForUserID func(userID string, userContext supertokens.UserContext) (TenantIdResult, error)

type TenantIdResult struct {
	OK *struct {
		TenantId *string
	}
	UnknownUserIDError *struct{}
}

type TypeInput struct {
	GetTenantIdForUserID  TypeGetTenantIdForUserID
	GetDomainsForTenantId func(tenantId *string, userContext supertokens.UserContext) ([]string, error)
	ErrorHandlers         *ErrorHandlers
	Override              *OverrideStruct
}

type ErrorHandlers struct {
	OnTenantDoesNotExistError      *func(err error, req *http.Request, res http.ResponseWriter) error
	OnRecipeDisabledForTenantError *func(err error, req *http.Request, res http.ResponseWriter) error
}

type TypeNormalisedInput struct {
	GetTenantIdForUserID  TypeGetTenantIdForUserID
	GetDomainsForTenantId func(tenantId *string, userContext supertokens.UserContext) ([]string, error)

	ErrorHandlers NormalisedErrorHandlers
	Override      OverrideStruct
}

type NormalisedErrorHandlers struct {
	OnTenantDoesNotExistError      func(err error, req *http.Request, res http.ResponseWriter) error
	OnRecipeDisabledForTenantError func(err error, req *http.Request, res http.ResponseWriter) error
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}
