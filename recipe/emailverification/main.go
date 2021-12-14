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

package emailverification

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config evmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func CreateEmailVerificationToken(userID, email string, userContext supertokens.UserContext) (evmodels.CreateEmailVerificationTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	return (*instance.RecipeImpl.CreateEmailVerificationToken)(userID, email, userContext)
}

func VerifyEmailUsingToken(token string, userContext supertokens.UserContext) (evmodels.VerifyEmailUsingTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.VerifyEmailUsingTokenResponse{}, err
	}
	return (*instance.RecipeImpl.VerifyEmailUsingToken)(token, userContext)
}

func IsEmailVerified(userID, email string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return (*instance.RecipeImpl.IsEmailVerified)(userID, email, userContext)
}

func RevokeEmailVerificationTokens(userID, email string, userContext supertokens.UserContext) (evmodels.RevokeEmailVerificationTokensResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.RevokeEmailVerificationTokensResponse{}, err
	}
	return (*instance.RecipeImpl.RevokeEmailVerificationTokens)(userID, email, userContext)
}

func UnverifyEmail(userID, email string, userContext supertokens.UserContext) (evmodels.UnverifyEmailResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	return (*instance.RecipeImpl.UnverifyEmail)(userID, email, userContext)
}
