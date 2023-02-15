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

package thirdparty

import (
	"encoding/json"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance *Recipe, appInfo supertokens.NormalisedAppinfo, config *tpmodels.TypeInput) (tpmodels.TypeNormalisedInput, error) {
	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance)

	signInAndUpFeature, err := validateAndNormaliseSignInAndUpConfig(config.SignInAndUpFeature)
	if err != nil {
		return tpmodels.TypeNormalisedInput{}, err
	}
	typeNormalisedInput.SignInAndUpFeature = signInAndUpFeature

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
	}

	return typeNormalisedInput, nil
}

func makeTypeNormalisedInput(recipeInstance *Recipe) tpmodels.TypeNormalisedInput {
	return tpmodels.TypeNormalisedInput{
		SignInAndUpFeature: tpmodels.TypeNormalisedInputSignInAndUp{},
		Override: tpmodels.OverrideStruct{
			Functions: func(originalImplementation tpmodels.RecipeInterface) tpmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation tpmodels.APIInterface) tpmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}

func validateAndNormaliseSignInAndUpConfig(config tpmodels.TypeInputSignInAndUp) (tpmodels.TypeNormalisedInputSignInAndUp, error) {
	providers := config.Providers
	// With the feature of dynamic providers from core, we need not check if the provider array is empty

	thirdPartyIdSet := map[string]bool{}

	for _, provider := range providers {
		if thirdPartyIdSet[provider.Config.ThirdPartyId] {
			return tpmodels.TypeNormalisedInputSignInAndUp{}, supertokens.BadInputError{Msg: "The providers array has multiple entries for the same third party provider."}
		}
		thirdPartyIdSet[provider.Config.ThirdPartyId] = true
	}

	return tpmodels.TypeNormalisedInputSignInAndUp{
		Providers: providers,
	}, nil
}

func parseUser(value interface{}) (*tpmodels.User, error) {
	respJSON, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var user tpmodels.User
	err = json.Unmarshal(respJSON, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func parseUsers(value interface{}) ([]tpmodels.User, error) {
	respJSON, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var users []interface{}
	err = json.Unmarshal(respJSON, &users)

	usersResult := make([]tpmodels.User, len(users))
	for i, user := range users {
		userObj, err := parseUser(user)
		if err != nil {
			return nil, err
		}
		usersResult[i] = *userObj
	}

	if err != nil {
		return nil, err
	}
	return usersResult, nil
}
