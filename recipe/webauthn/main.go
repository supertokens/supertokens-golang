/* Copyright (c) 2025, VRAI Labs and/or its affiliates. All rights reserved.
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

package webauthn

import (
	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ListCredentials(recipeUserId string, userContext ...supertokens.UserContext) (webauthnmodels.ListCredentialsResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return webauthnmodels.ListCredentialsResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ListCredentials)(recipeUserId, userContext[0])
}

func GetCredential(webauthnCredentialId string, recipeUserId string, tenantId string, userContext ...supertokens.UserContext) (webauthnmodels.GetCredentialResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return webauthnmodels.GetCredentialResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetCredential)(webauthnCredentialId, recipeUserId, tenantId, userContext[0])
}

func RemoveCredential(webauthnCredentialId string, recipeUserId string, userContext ...supertokens.UserContext) (webauthnmodels.RemoveCredentialResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return webauthnmodels.RemoveCredentialResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.RemoveCredential)(webauthnCredentialId, recipeUserId, userContext[0])
}
