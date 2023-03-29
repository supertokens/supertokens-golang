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

type signInUpResponse struct {
	CreatedNewUser bool
	User           tpmodels.User
}

func Init(config *tpmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func SignInUpWithContext(thirdPartyID string, thirdPartyUserID string, email string, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.SignInUpResponse{}, err
	}
	return (*instance.RecipeImpl.SignInUp)(thirdPartyID, thirdPartyUserID, email, userContext)
}

func GetUserByIDWithContext(userID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUsersByEmailWithContext(email string, userContext supertokens.UserContext) ([]tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return []tpmodels.User{}, err
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, userContext)
}

func GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, userContext)
}

func SignInUp(thirdPartyID string, thirdPartyUserID string, email string) (tpmodels.SignInUpResponse, error) {
	return SignInUpWithContext(thirdPartyID, thirdPartyUserID, email, &map[string]interface{}{})
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

func Apple(config tpmodels.AppleConfig) tpmodels.TypeProvider {
	return providers.Apple(config)
}

func Facebook(config tpmodels.FacebookConfig) tpmodels.TypeProvider {
	return providers.Facebook(config)
}

func Github(config tpmodels.GithubConfig) tpmodels.TypeProvider {
	return providers.Github(config)
}

func Discord(config tpmodels.DiscordConfig) tpmodels.TypeProvider {
	return providers.Discord(config)
}

func GoogleWorkspaces(config tpmodels.GoogleWorkspacesConfig) tpmodels.TypeProvider {
	return providers.GoogleWorkspaces(config)
}

func Bitbucket(config tpmodels.BitbucketConfig) tpmodels.TypeProvider {
	return providers.Bitbucket(config)
}

func GitLab(config tpmodels.GitLabConfig) tpmodels.TypeProvider {
	return providers.GitLab(config)
}

func Google(config tpmodels.GoogleConfig) tpmodels.TypeProvider {
	return providers.Google(config)
}
