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

package tplmodels

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type User struct {
	ID          string  `json:"id"`
	TimeJoined  uint64  `json:"timeJoined"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	ThirdParty  *struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
	} `json:"thirdParty"`
}

type TypeInput struct {
	ContactMethodPhone        plessmodels.ContactMethodPhoneConfig
	ContactMethodEmail        plessmodels.ContactMethodEmailConfig
	ContactMethodEmailOrPhone plessmodels.ContactMethodEmailOrPhoneConfig
	FlowType                  string
	GetCustomUserInputCode    func(userContext supertokens.UserContext) (string, error)
	Providers                 []tpmodels.TypeProvider
	Override                  *OverrideStruct
	EmailDelivery             *emaildelivery.TypeInput
	SmsDelivery               *smsdelivery.TypeInput
}

type TypeNormalisedInput struct {
	ContactMethodPhone        plessmodels.ContactMethodPhoneConfig
	ContactMethodEmail        plessmodels.ContactMethodEmailConfig
	ContactMethodEmailOrPhone plessmodels.ContactMethodEmailOrPhoneConfig
	FlowType                  string
	GetCustomUserInputCode    func(userContext supertokens.UserContext) (string, error)
	Providers                 []tpmodels.TypeProvider
	Override                  OverrideStruct
	GetEmailDeliveryConfig    func() emaildelivery.TypeInputWithService
	GetSmsDeliveryConfig      func() smsdelivery.TypeInputWithService
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}
