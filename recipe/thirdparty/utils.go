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
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/emaildelivery/backwardCompatibilityService"
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

	typeNormalisedInput.EmailVerificationFeature = validateAndNormaliseEmailVerificationConfig(recipeInstance, config)

	typeNormalisedInput.GetEmailDeliveryConfig = func(recipeImpl tpmodels.RecipeInterface) emaildelivery.TypeInputWithService {
		sendEmailVerificationEmail := emailverification.DefaultCreateAndSendCustomEmail(appInfo)
		if typeNormalisedInput.EmailVerificationFeature.CreateAndSendCustomEmail != nil {
			sendEmailVerificationEmail = typeNormalisedInput.EmailVerificationFeature.CreateAndSendCustomEmail
		}

		emailService := backwardCompatibilityService.MakeBackwardCompatibilityService(recipeImpl, appInfo, sendEmailVerificationEmail)
		if config != nil && config.EmailDelivery != nil && config.EmailDelivery.Service != nil {
			emailService = *config.EmailDelivery.Service
		}
		result := emaildelivery.TypeInputWithService{
			Service: emailService,
		}
		if config != nil && config.EmailDelivery != nil && config.EmailDelivery.Override != nil {
			result.Override = config.EmailDelivery.Override
		}

		return result
	}

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
		if config.Override.EmailVerificationFeature != nil {
			typeNormalisedInput.Override.EmailVerificationFeature = config.Override.EmailVerificationFeature
		}
	}

	return typeNormalisedInput, nil
}

func makeTypeNormalisedInput(recipeInstance *Recipe) tpmodels.TypeNormalisedInput {
	return tpmodels.TypeNormalisedInput{
		SignInAndUpFeature:       tpmodels.TypeNormalisedInputSignInAndUp{},
		EmailVerificationFeature: validateAndNormaliseEmailVerificationConfig(recipeInstance, nil),
		Override: tpmodels.OverrideStruct{
			Functions: func(originalImplementation tpmodels.RecipeInterface) tpmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation tpmodels.APIInterface) tpmodels.APIInterface {
				return originalImplementation
			},
			EmailVerificationFeature: nil,
		},
	}
}

func validateAndNormaliseEmailVerificationConfig(recipeInstance *Recipe, config *tpmodels.TypeInput) evmodels.TypeInput {
	emailverificationTypeInput := evmodels.TypeInput{
		GetEmailForUserID: recipeInstance.getEmailForUserId,
		Override:          nil,
	}

	if config != nil {
		if config.Override != nil {
			emailverificationTypeInput.Override = config.Override.EmailVerificationFeature
		}
		if config.EmailVerificationFeature != nil {
			if config.EmailVerificationFeature.CreateAndSendCustomEmail != nil {
				emailverificationTypeInput.CreateAndSendCustomEmail = func(user evmodels.User, link string, userContext supertokens.UserContext) {
					userInfo, err := (*recipeInstance.RecipeImpl.GetUserByID)(user.ID, userContext)
					if err != nil {
						return
					}
					if userInfo == nil {
						return
					}
					config.EmailVerificationFeature.CreateAndSendCustomEmail(*userInfo, link, userContext)
				}
			}

			if config.EmailVerificationFeature.GetEmailVerificationURL != nil {
				emailverificationTypeInput.GetEmailVerificationURL = func(user evmodels.User, userContext supertokens.UserContext) (string, error) {
					userInfo, err := (*recipeInstance.RecipeImpl.GetUserByID)(user.ID, userContext)
					if err != nil {
						return "", err
					}
					if userInfo == nil {
						return "", errors.New("unknown User ID provided")
					}
					return config.EmailVerificationFeature.GetEmailVerificationURL(*userInfo, userContext)
				}
			}
		}
	}

	return emailverificationTypeInput
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
