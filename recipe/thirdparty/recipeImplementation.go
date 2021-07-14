package thirdparty

import (
	"encoding/json"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier) models.RecipeImplementation {
	return models.RecipeImplementation{
		SignInUp: func(thirdPartyID, thirdPartyUserID string, email struct {
			ID         string
			IsVerified bool
		}) models.SignInUpResponse {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/signinup")
			if err != nil {
				return models.SignInUpResponse{
					Status: "FIELD_ERROR",
					Error:  err,
				}
			}
			response, err := querier.SendPostRequest(*normalisedURLPath, map[string]interface{}{
				"thirdPartyId":     thirdPartyID,
				"thirdPartyUserId": thirdPartyUserID,
				"email":            email,
			})
			if err != nil {
				return models.SignInUpResponse{
					Status: "FIELD_ERROR",
					Error:  err,
				}
			}
			respJSON, err := json.Marshal(response["user"])
			if err != nil {
				return models.SignInUpResponse{
					Status: "FIELD_ERROR",
					Error:  err,
				}
			}
			var user models.User
			err = json.Unmarshal(respJSON, &user)
			if err != nil {
				return models.SignInUpResponse{
					Status: "FIELD_ERROR",
					Error:  err,
				}
			}
			return models.SignInUpResponse{
				Status:         "OK",
				User:           user,
				CreatedNewUser: response["createdNewUser"].(bool),
			}
		},
		GetUserById: func(userID string) *models.User {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user")
			if err != nil {
				return nil
			}
			response, err := querier.SendGetRequest(*normalisedURLPath, map[string]interface{}{
				"userId": userID,
			})
			if response["status"] == "OK" {
				respJSON, err := json.Marshal(response["user"])
				if err != nil {
					return nil
				}
				var user models.User
				err = json.Unmarshal(respJSON, &user)
				if err != nil {
					return nil
				}
				return &user
			}
			return nil
		},
		GetUserByThirdPartyInfo: func(thirdPartyID, thirdPartyUserID string) *models.User {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user")
			if err != nil {
				return nil
			}
			response, err := querier.SendGetRequest(*normalisedURLPath, map[string]interface{}{
				"thirdPartyId":     thirdPartyID,
				"thirdPartyUserId": thirdPartyUserID,
			})
			if response["status"] == "OK" {
				respJSON, err := json.Marshal(response["user"])
				if err != nil {
					return nil
				}
				var user models.User
				err = json.Unmarshal(respJSON, &user)
				if err != nil {
					return nil
				}
				return &user
			}
			return nil
		},
	}
}
