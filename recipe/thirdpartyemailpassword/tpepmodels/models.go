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

package tpepmodels

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

type User struct {
	ID         string `json:"id"`
	TimeJoined uint64 `json:"timeJoined"`
	Email      string `json:"email"`
	ThirdParty *struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
	} `json:"thirdParty"`
	TenantIds []string `json:"tenantIds"`
}

type TypeContext struct {
	FormFields                 []epmodels.TypeFormField
	ThirdPartyAuthCodeResponse interface{}
}

type TypeInput struct {
	SignUpFeature *epmodels.TypeInputSignUp
	Providers     []tpmodels.ProviderInput
	Override      *OverrideStruct
	EmailDelivery *emaildelivery.TypeInput
}

type TypeNormalisedInput struct {
	SignUpFeature          *epmodels.TypeInputSignUp
	Providers              []tpmodels.ProviderInput
	Override               OverrideStruct
	GetEmailDeliveryConfig func(recipeImpl RecipeInterface, epRecipeImpl epmodels.RecipeInterface) emaildelivery.TypeInputWithService
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}
