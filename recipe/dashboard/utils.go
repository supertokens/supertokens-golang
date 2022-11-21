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

package dashboard

import (
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config dashboardmodels.TypeInput) dashboardmodels.TypeNormalisedInput {
	typeNormalisedInput := makeTypeNormalisedInput(appInfo)

	if strings.Trim(config.ApiKey, " ") == "" {
		panic("ApiKey provided to Dashboard recipe cannot be empty")
	}

	typeNormalisedInput.ApiKey = config.ApiKey

	if config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
	}

	return typeNormalisedInput
}

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) dashboardmodels.TypeNormalisedInput {
	return dashboardmodels.TypeNormalisedInput{
		Override: dashboardmodels.OverrideStruct{
			Functions: func(originalImplementation dashboardmodels.RecipeInterface) dashboardmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation dashboardmodels.APIInterface) dashboardmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}

func isApiPath(path supertokens.NormalisedURLPath, appInfo supertokens.NormalisedAppinfo) (bool, error) {
	normalizedDashboardAPI, err := supertokens.NewNormalisedURLPath(dashboardAPI)
	if err != nil {
		return false, err
	}
	dashboardRecipeBasePath := appInfo.APIBasePath.AppendPath(normalizedDashboardAPI)
	if !path.StartsWith(dashboardRecipeBasePath) {
		return false, nil
	}

	pathWithoutDashboardPath := strings.Split(path.GetAsStringDangerous(), dashboardAPI)[1]
	if len(pathWithoutDashboardPath) > 0 && pathWithoutDashboardPath[0] == '/' {
		pathWithoutDashboardPath = pathWithoutDashboardPath[1:]
	}

	if strings.Split(pathWithoutDashboardPath, "/")[0] == "api" {
		return true, nil
	}
	return false, nil
}

func getApiIdIfMatched(path supertokens.NormalisedURLPath, method string) (*string, error) {
	if method == http.MethodPost && strings.HasSuffix(path.GetAsStringDangerous(), validateKeyAPI) {
		val := validateKeyAPI
		return &val, nil
	}

	if method == http.MethodGet && strings.HasSuffix(path.GetAsStringDangerous(), usersListGetAPI) {
		val := usersListGetAPI
		return &val, nil
	}

	if method == http.MethodGet && strings.HasSuffix(path.GetAsStringDangerous(), usersCountAPI) {
		val := usersCountAPI
		return &val, nil
	}

	if method == http.MethodGet && strings.HasSuffix(path.GetAsStringDangerous(), userAPI) {
		val := userAPI
		return &val, nil
	}

	return nil, nil
}

func IsValidRecipeId(recipeId string)(bool) {
	return recipeId == "emailpassword" || recipeId == "thirdparty" || recipeId == "passwordless"
}

func GetUserForRecipeId(userId string, recipeId string)(user dashboardmodels.UserType, recipe string) {
	var userToReturn dashboardmodels.UserType
	var recipeToReturn string

	if recipeId == emailpassword.RECIPE_ID {
		response, error := emailpassword.GetUserByID(userId)

		if error == nil {
			userToReturn.Id = response.ID
			userToReturn.TimeJoined = response.TimeJoined
			userToReturn.FirstName = ""
			userToReturn.LastName = ""
			userToReturn.Email = response.Email

			recipeToReturn = emailpassword.RECIPE_ID
		}

		if userToReturn == (dashboardmodels.UserType{}) {
			tpepResponse, tpepError := thirdpartyemailpassword.GetUserById(userId)

			if tpepError == nil {
				userToReturn.Id = tpepResponse.ID
				userToReturn.TimeJoined = tpepResponse.TimeJoined
				userToReturn.FirstName = ""
				userToReturn.LastName = ""
				userToReturn.Email = tpepResponse.Email

				recipeToReturn = thirdpartyemailpassword.RECIPE_ID
			}
		}
	} else if recipeId == thirdparty.RECIPE_ID {
		response, error := thirdparty.GetUserByID(userId)

		if error == nil {
			userToReturn.Id = response.ID
			userToReturn.TimeJoined = response.TimeJoined
			userToReturn.FirstName = ""
			userToReturn.LastName = ""
			userToReturn.Email = response.Email
			userToReturn.ThirdParty.Id = response.ThirdParty.ID
			userToReturn.ThirdParty.UserId = response.ThirdParty.UserID
		}

		if userToReturn == (dashboardmodels.UserType{}) {
			tpepResponse, tpepError := thirdpartyemailpassword.GetUserById(userId)

			if tpepError == nil {
				userToReturn.Id = tpepResponse.ID
				userToReturn.TimeJoined = tpepResponse.TimeJoined
				userToReturn.FirstName = ""
				userToReturn.LastName = ""
				userToReturn.Email = tpepResponse.Email
				userToReturn.ThirdParty.Id = tpepResponse.ThirdParty.ID
				userToReturn.ThirdParty.UserId = tpepResponse.ThirdParty.UserID
			}
		}
	} else if recipeId == passwordless.RECIPE_ID {
		response, error := passwordless.GetUserByID(userId)

		if error == nil {
			userToReturn.Id = response.ID
			userToReturn.TimeJoined = response.TimeJoined
			userToReturn.FirstName = ""
			userToReturn.LastName = ""

			if response.Email != nil {
				userToReturn.Email = *response.Email
			}

			if response.PhoneNumber != nil {
				userToReturn.Phone = *response.PhoneNumber
			}
		}

		if userToReturn == (dashboardmodels.UserType{}) {
			tppResponse, tppError := thirdpartypasswordless.GetUserByID(userId)

			if tppError == nil {
				userToReturn.Id = tppResponse.ID
				userToReturn.TimeJoined = tppResponse.TimeJoined
				userToReturn.FirstName = ""
				userToReturn.LastName = ""

				if tppResponse.Email != nil {
					userToReturn.Email = *tppResponse.Email
				}

				if tppResponse.PhoneNumber != nil {
					userToReturn.Phone = *tppResponse.PhoneNumber
				}
			}
		}
	}

	return userToReturn, recipeToReturn
}
