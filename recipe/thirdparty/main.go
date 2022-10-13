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
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *tpmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func CreateUserWithContext(thirdPartyID string, thirdPartyUserID string, email string, userContext supertokens.UserContext) (tpmodels.CreateUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.CreateUserResponse{}, err
	}
	return (*instance.RecipeImpl.CreateUser)(thirdPartyID, thirdPartyUserID, email, userContext)
}

func GetUserByIDWithContext(userID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUsersByEmailWithContext(email string, userContext supertokens.UserContext) ([]tpmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return []tpmodels.User{}, err
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, userContext)
}

func GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, userContext)
}

func CreateUser(thirdPartyID string, thirdPartyUserID string, email string) (tpmodels.CreateUserResponse, error) {
	return CreateUserWithContext(thirdPartyID, thirdPartyUserID, email, &map[string]interface{}{})
}

func GetUserByID(userID string) (*tpmodels.User, error) {
	return GetUserByIDWithContext(userID, &map[string]interface{}{})
}

func GetUsersByEmail(email string) ([]tpmodels.User, error) {
	return GetUsersByEmailWithContext(email, &map[string]interface{}{})
}

func GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID string) (*tpmodels.User, error) {
	return GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID, &map[string]interface{}{})
}

// TODO: enable below once implemented
// func Apple(input tpmodels.TypeAppleInput) tpmodels.TypeProvider {
// 	return providers.Apple(input)
// }

// func Facebook(input tpmodels.TypeFacebookInput) tpmodels.TypeProvider {
// 	return providers.Facebook(input)
// }

// func Github(input tpmodels.TypeGithubInput) tpmodels.TypeProvider {
// 	return providers.Github(input)
// }

// func Discord(input tpmodels.TypeDiscordInput) tpmodels.TypeProvider {
// 	return providers.Discord(input)
// }

// func GoogleWorkspaces(input tpmodels.TypeGoogleWorkspacesInput) tpmodels.TypeProvider {
// 	return providers.GoogleWorkspaces(input)
// }

func Google(input tpmodels.TypeGoogleInput) tpmodels.TypeProvider {
	return providers.Google(input)
}
