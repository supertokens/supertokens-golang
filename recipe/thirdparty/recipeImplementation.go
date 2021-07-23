package thirdparty

import (
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier) models.RecipeImplementation {
	return models.RecipeImplementation{
		SignInUp: func(thirdPartyID, thirdPartyUserID string, email models.EmailStruct) models.SignInUpResponse {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/signinup")
			if err != nil {
				fmt.Println("here")
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
				fmt.Println("here 2")
				return models.SignInUpResponse{
					Status: "FIELD_ERROR",
					Error:  err,
				}
			}
			user, err := parseUser(response["user"])
			if err != nil {
				fmt.Println("here 3")
				return models.SignInUpResponse{
					Status: "FIELD_ERROR",
					Error:  err,
				}
			}
			return models.SignInUpResponse{
				Status:         "OK",
				User:           *user,
				CreatedNewUser: response["createdNewUser"].(bool),
			}
		},
		GetUserByID: func(userID string) *models.User {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user")
			if err != nil {
				return nil
			}
			response, err := querier.SendGetRequest(*normalisedURLPath, map[string]interface{}{
				"userId": userID,
			})
			if response["status"] == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return nil
				}
				return user
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
				user, err := parseUser(response["user"])
				if err != nil {
					return nil
				}
				return user
			}
			return nil
		},
	}
}
