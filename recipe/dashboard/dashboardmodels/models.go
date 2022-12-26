/* Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package dashboardmodels

type TypeInput struct {
	ApiKey   string
	Override *OverrideStruct
}

type TypeNormalisedInput struct {
	ApiKey   string
	Override OverrideStruct
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}

type ThirdParty struct {
	Id     string `json:"id,omitempty"`
	UserId string `json:"userId,omitempty"`
}

type UserType struct {
	Id         string     `json:"id,omitempty"`
	TimeJoined uint64     `json:"timeJoined,omitempty"`
	FirstName  string     `json:"firstName,omitempty"`
	LastName   string     `json:"lastName,omitempty"`
	Email      string     `json:"email,omitempty"`
	ThirdParty ThirdParty `json:"thirdParty,omitempty"`
	Phone      string     `json:"phoneNumber,omitempty"`
}
