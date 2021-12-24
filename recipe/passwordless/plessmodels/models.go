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

package plessmodels

import "github.com/supertokens/supertokens-golang/supertokens"

type User struct {
	ID          string  `json:"id"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	TimeJoined  uint64  `json:"timejoined"`
}

type TypeInput struct {
	ContactMethodPhone        ContactMethodPhoneConfig
	ContactMethodEmail        ContactMethodEmailConfig
	ContactMethodEmailOrPhone ContactMethodEmailOrPhoneConfig
	FlowType                  string
	GetLinkDomainAndPath      func(email *string, phoneNumber *string, userContext supertokens.UserContext) (string, error)
	GetCustomUserInputCode    func(userContext supertokens.UserContext) (string, error)
	Override                  *OverrideStruct
}

type TypeNormalisedInput struct {
	ContactMethodPhone        ContactMethodPhoneConfig
	ContactMethodEmail        ContactMethodEmailConfig
	ContactMethodEmailOrPhone ContactMethodEmailOrPhoneConfig
	FlowType                  string
	GetLinkDomainAndPath      func(email *string, phoneNumber *string, userContext supertokens.UserContext) (string, error)
	GetCustomUserInputCode    func(userContext supertokens.UserContext) (string, error)
	Override                  OverrideStruct
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}

type ContactMethodEmailConfig struct {
	Enabled                  bool
	ValidateEmailAddress     func(email interface{}) *string
	CreateAndSendCustomEmail func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext)
}

type ContactMethodEmailOrPhoneConfig struct {
	Enabled                        bool
	ValidateEmailAddress           func(email interface{}) *string
	CreateAndSendCustomEmail       func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext)
	ValidatePhoneNumber            func(phoneNumber interface{}) *string
	CreateAndSendCustomTextMessage func(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext)
}

type ContactMethodPhoneConfig struct {
	Enabled                        bool
	ValidatePhoneNumber            func(phoneNumber interface{}) *string
	CreateAndSendCustomTextMessage func(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext)
}
