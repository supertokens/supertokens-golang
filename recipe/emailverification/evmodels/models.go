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

package evmodels

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeGetEmailForUserID func(userID string, userContext supertokens.UserContext) (TypeEmailInfo, error)

type TypeMode string

const (
	ModeRequired TypeMode = "REQUIRED"
	ModeOptional TypeMode = "OPTIONAL"
)

type TypeEmailInfo struct {
	OK *struct {
		Email string
	}
	EmailDoesNotExistError *struct{}
	UnknownUserIDError     *struct{}
}

type TypeInput struct {
	Mode              TypeMode
	GetEmailForUserID TypeGetEmailForUserID
	Override          *OverrideStruct
	EmailDelivery     *emaildelivery.TypeInput
}

type TypeNormalisedInput struct {
	Mode                   TypeMode
	GetEmailForUserID      TypeGetEmailForUserID
	Override               OverrideStruct
	GetEmailDeliveryConfig func() emaildelivery.TypeInputWithService
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
