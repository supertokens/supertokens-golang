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
	if len(providers) == 0 {
		return tpmodels.TypeNormalisedInputSignInAndUp{}, supertokens.BadInputError{Msg: "thirdparty recipe requires at least 1 provider to be passed in signInAndUpFeature.providers config"}
	}

	isDefaultProvidersSet := map[string]bool{}
	allProvidersSet := map[string]bool{}

	for i := 0; i < len(providers); i++ {
		id := providers[i].ID
		allProvidersSet[id] = true
		isDefault := providers[i].IsDefault

		// if this is the only provider with this ID, then we mark this as default
		var otherProvidersWithSameId []tpmodels.TypeProvider = []tpmodels.TypeProvider{}
		for y := 0; y < len(providers); y++ {
			if providers[y].ID == id && &providers[y] != &providers[i] {
				otherProvidersWithSameId = append(otherProvidersWithSameId, providers[y])
			}
		}
		if len(otherProvidersWithSameId) == 0 {
			isDefault = true
		}

		if isDefault {
			if isDefaultProvidersSet[id] {
				return tpmodels.TypeNormalisedInputSignInAndUp{}, supertokens.BadInputError{Msg: "You have provided multiple third party providers that have the id: " + providers[i].ID + " and are marked as 'IsDefault: true'. Please only mark one of them as isDefault"}
			}
			isDefaultProvidersSet[id] = true
		}
	}

	if len(isDefaultProvidersSet) != len(allProvidersSet) {
		return tpmodels.TypeNormalisedInputSignInAndUp{}, supertokens.BadInputError{Msg: "The providers array has multiple entries for the same third party provider. Please mark one of them as the default one by using 'IsDefault: true'"}
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
	var user []tpmodels.User
	err = json.Unmarshal(respJSON, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
