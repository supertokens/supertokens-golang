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

package session

import (
	"context"
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *sessmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func CreateNewSession(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.SessionContainer{}, err
	}
	return (*instance.RecipeImpl.CreateNewSession)(res, userID, accessTokenPayload, sessionData, userContext)
}

func GetSession(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (*sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetSession)(req, res, options, userContext)
}

func GetSessionInformation(sessionHandle string, userContext supertokens.UserContext) (sessmodels.SessionInformation, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.SessionInformation{}, err
	}
	return (*instance.RecipeImpl.GetSessionInformation)(sessionHandle, userContext)
}

func RefreshSession(req *http.Request, res http.ResponseWriter, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.SessionContainer{}, err
	}
	return (*instance.RecipeImpl.RefreshSession)(req, res, userContext)
}

func RevokeAllSessionsForUser(userID string, userContext supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.RevokeAllSessionsForUser)(userID, userContext)
}

func GetAllSessionHandlesForUser(userID string, userContext supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetAllSessionHandlesForUser)(userID, userContext)
}

func RevokeSession(sessionHandle string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return (*instance.RecipeImpl.RevokeSession)(sessionHandle, userContext)
}

func RevokeMultipleSessions(sessionHandles []string, userContext supertokens.UserContext) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.RevokeMultipleSessions)(sessionHandles, userContext)
}

func UpdateSessionData(sessionHandle string, newSessionData map[string]interface{}, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.RecipeImpl.UpdateSessionData)(sessionHandle, newSessionData, userContext)
}

func UpdateAccessTokenPayload(sessionHandle string, newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.RecipeImpl.UpdateAccessTokenPayload)(sessionHandle, newAccessTokenPayload, userContext)
}

func VerifySession(options *sessmodels.VerifySessionOptions, otherHandler http.HandlerFunc) http.HandlerFunc {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		panic("can't fetch supertokens instance. You should call the supertokens.Init function before using the VerifySession function.")
	}
	return VerifySessionHelper(*instance, options, otherHandler)
}

func GetSessionFromRequestContext(ctx context.Context) *sessmodels.SessionContainer {
	value := ctx.Value(sessmodels.SessionContext)
	if value == nil {
		return nil
	}
	temp := value.(*sessmodels.SessionContainer)
	return temp
}

func CreateJWT(payload map[string]interface{}, validitySecondsPointer *uint64, userContext supertokens.UserContext) (jwtmodels.CreateJWTResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return jwtmodels.CreateJWTResponse{}, err
	}
	if instance.OpenIdRecipe == nil {
		return jwtmodels.CreateJWTResponse{}, errors.New("CreateJWT cannot be used without enabling the Jwt feature")
	}
	return (*instance.OpenIdRecipe.RecipeImpl.CreateJWT)(payload, validitySecondsPointer, userContext)
}

func GetJWKS(userContext supertokens.UserContext) (jwtmodels.GetJWKSResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return jwtmodels.GetJWKSResponse{}, err
	}
	if instance.OpenIdRecipe == nil {
		return jwtmodels.GetJWKSResponse{}, errors.New("GetJWKS cannot be used without enabling the Jwt feature")
	}
	return (*instance.OpenIdRecipe.RecipeImpl.GetJWKS)(userContext)
}

func GetOpenIdDiscoveryConfiguration(userContext supertokens.UserContext) (openidmodels.GetOpenIdDiscoveryConfigurationResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return openidmodels.GetOpenIdDiscoveryConfigurationResponse{}, err
	}
	if instance.OpenIdRecipe == nil {
		return openidmodels.GetOpenIdDiscoveryConfigurationResponse{}, errors.New("GetOpenIdDiscoveryConfiguration cannot be used without enabling the Jwt feature")
	}
	return (*instance.OpenIdRecipe.RecipeImpl.GetOpenIdDiscoveryConfiguration)(userContext)
}
