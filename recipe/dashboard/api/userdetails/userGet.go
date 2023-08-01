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

package userdetails

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/api"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userGetResponse struct {
	Status   string                   `json:"status"`
	RecipeId string                   `json:"recipeId,omitempty"`
	User     dashboardmodels.UserType `json:"user,omitempty"`
}

func UserGet(apiImplementation dashboardmodels.APIInterface, tenantId string, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (userGetResponse, error) {
	req := options.Req
	userId := req.URL.Query().Get("userId")
	recipeId := req.URL.Query().Get("recipeId")

	if userId == "" {
		return userGetResponse{}, supertokens.BadInputError{
			Msg: "Missing required parameter 'userId'",
		}
	}

	if recipeId == "" {
		return userGetResponse{}, supertokens.BadInputError{
			Msg: "Missing required parameter 'recipeId'",
		}
	}

	if !api.IsValidRecipeId(recipeId) {
		return userGetResponse{}, supertokens.BadInputError{
			Msg: "Invalid recipe id",
		}
	}

	if !api.IsRecipeInitialised(recipeId) {
		return userGetResponse{
			Status: "RECIPE_NOT_INITIALISED",
		}, nil
	}

	userForRecipeId, _ := api.GetUserForRecipeId(userId, recipeId, userContext)

	if userForRecipeId == (dashboardmodels.UserType{}) {
		return userGetResponse{
			Status: "NO_USER_FOUND_ERROR",
		}, nil
	}

	_, err := usermetadata.GetRecipeInstanceOrThrowError()

	if err != nil {
		// If metadata is not enabled then the frontend will show this as the name
		userForRecipeId.FirstName = "FEATURE_NOT_ENABLED"
		userForRecipeId.LastName = "FEATURE_NOT_ENABLED"

		return userGetResponse{
			Status:   "OK",
			RecipeId: recipeId,
			User:     userForRecipeId,
		}, nil
	}

	metadata, metadataerr := usermetadata.GetUserMetadata(userId, userContext)

	if metadataerr != nil {
		return userGetResponse{}, metadataerr
	}

	// first and last name should be an empty string if they dont exist in metadata
	userForRecipeId.FirstName = ""
	userForRecipeId.LastName = ""

	if metadata["first_name"] != nil {
		userForRecipeId.FirstName = metadata["first_name"].(string)
	}

	if metadata["last_name"] != nil {
		userForRecipeId.LastName = metadata["last_name"].(string)
	}

	return userGetResponse{
		Status:   "OK",
		RecipeId: recipeId,
		User:     userForRecipeId,
	}, nil
}
