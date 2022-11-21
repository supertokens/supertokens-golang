package userdetails

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userGetResponse struct {
	Status string  `json:"status"`
	RecipeId string `json:"recipeId,omitempty"`
	User dashboardmodels.UserType `json:"user,omitempty"`
}

func UserGet(apiImplementation dashboardmodels.APIInterface, options dashboardmodels.APIOptions)(userGetResponse, error) {
	req := options.Req
	userId := req.URL.Query().Get("userId")
	recipeId := req.URL.Query().Get("recipeId")

	if userId == "" {
		return userGetResponse{}, supertokens.BadInputError {
			Msg: "Missing required parameter 'userId'",
		}
	}

	if recipeId == "" {
		return userGetResponse{}, supertokens.BadInputError {
			Msg: "Missing required parameter 'recipeId'",
		}
	}

	if dashboard.IsValidRecipeId(recipeId) {
		return userGetResponse{}, supertokens.BadInputError {
			Msg: "Invalid recipe id",
		}
	}

	userForRecipeId, _ := dashboard.GetUserForRecipeId(userId, recipeId)

	if userForRecipeId == (dashboardmodels.UserType{}) {
		return userGetResponse{
			Status: "NO_USER_FOUND_ERROR",
		}, nil
	}

	_, err := usermetadata.GetRecipeInstanceOrThrowError()

	if err != nil {
		userForRecipeId.FirstName = "FEATURE_NOT_ENABLED"
		userForRecipeId.LastName = "FEATURE_NOT_ENABLED"

		return userGetResponse{
			Status: "OK",
			RecipeId: recipeId,
			User: userForRecipeId,
		}, nil
	}

	metadata, metadataerr := usermetadata.GetUserMetadata(userId)

	if metadataerr == nil {
		userForRecipeId.FirstName = metadata["first_name"].(string)
		userForRecipeId.LastName = metadata["last_name"].(string)
	}

	return userGetResponse{
		Status: "OK",
		RecipeId: recipeId,
		User: userForRecipeId,
	}, nil
}