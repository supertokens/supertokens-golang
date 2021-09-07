package thirdparty

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier) models.RecipeInterface {
	return models.RecipeInterface{
		SignInUp: func(thirdPartyID, thirdPartyUserID string, email models.EmailStruct) (models.SignInUpResponse, error) {
			response, err := querier.SendPostRequest("/recipe/signinup", map[string]interface{}{
				"thirdPartyId":     thirdPartyID,
				"thirdPartyUserId": thirdPartyUserID,
				"email":            email,
			})
			if err != nil {
				return models.SignInUpResponse{}, err
			}
			user, err := parseUser(response["user"])
			if err != nil {
				return models.SignInUpResponse{}, err
			}
			return models.SignInUpResponse{
				OK: &struct {
					CreatedNewUser bool
					User           models.User
				}{
					CreatedNewUser: response["createdNewUser"].(bool),
					User:           *user,
				},
			}, nil
		},

		GetUserByID: func(userID string) (*models.User, error) {
			response, err := querier.SendGetRequest("/recipe/user", map[string]interface{}{
				"userId": userID,
			})
			if err != nil {
				return nil, err
			}
			if response["status"] == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return nil, err
				}
				return user, nil
			}
			return nil, nil
		},

		GetUserByThirdPartyInfo: func(thirdPartyID, thirdPartyUserID string) (*models.User, error) {
			response, err := querier.SendGetRequest("/recipe/user", map[string]interface{}{
				"thirdPartyId":     thirdPartyID,
				"thirdPartyUserId": thirdPartyUserID,
			})
			if err != nil {
				return nil, err
			}
			if response["status"] == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return nil, err
				}
				return user, nil
			}
			return nil, nil
		},

		GetUsersByEmail: func(email string) ([]models.User, error) {
			response, err := querier.SendGetRequest("/recipe/users/by-email", map[string]interface{}{
				"email": email,
			})
			if err != nil {
				return []models.User{}, err
			}
			users, err := parseUsers(response["users"])
			if err != nil {
				return []models.User{}, err
			}
			return users, nil
		},
	}
}
