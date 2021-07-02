package emailpassword

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/constants"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier) models.RecipeImplementation {
	return models.RecipeImplementation{
		SignUp: func(email, password string) models.SignInUpResponse {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/signup")
			if err != nil {
				return models.SignInUpResponse{
					Status: constants.EmailAlreadyExistsError,
				}
			}
			response, err := querier.SendPostRequest(*normalisedURLPath, map[string]interface{}{
				"email":    email,
				"password": password,
			})
			if err != nil {
				return models.SignInUpResponse{
					Status: constants.EmailAlreadyExistsError,
				}
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				return models.SignInUpResponse{
					Status: status.(string),
					User:   response["user"].(models.User),
				}
			}
			return models.SignInUpResponse{
				Status: constants.EmailAlreadyExistsError,
			}
		},

		SignIn: func(email, password string) models.SignInUpResponse {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/signin")
			if err != nil {
				return models.SignInUpResponse{
					Status: constants.WrongCredentialsError,
				}
			}
			response, err := querier.SendPostRequest(*normalisedURLPath, map[string]interface{}{
				"email":    email,
				"password": password,
			})
			if err != nil {
				return models.SignInUpResponse{
					Status: constants.WrongCredentialsError,
				}
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				return models.SignInUpResponse{
					Status: status.(string),
					User:   response["user"].(models.User),
				}
			}
			return models.SignInUpResponse{
				Status: constants.WrongCredentialsError,
			}
		},

		GetUserById: func(userId string) *models.User {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user")
			if err != nil {
				return nil
			}
			response, err := querier.SendGetRequest(*normalisedURLPath, map[string]interface{}{
				"userId": userId,
			})
			if err != nil {
				return nil
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				user := response["user"].(models.User)
				return &user
			}
			return nil
		},

		GetUserByEmail: func(email string) *models.User {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user")
			if err != nil {
				return nil
			}
			response, err := querier.SendGetRequest(*normalisedURLPath, map[string]interface{}{
				"email": email,
			})
			if err != nil {
				return nil
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				user := response["user"].(models.User)
				return &user
			}
			return nil
		},

		CreateResetPasswordToken: func(userId string) models.CreateResetPasswordTokenResponse {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/password/reset/token")
			if err != nil {
				return models.CreateResetPasswordTokenResponse{
					Status: constants.UnknownUserID,
				}
			}
			response, err := querier.SendPostRequest(*normalisedURLPath, map[string]interface{}{
				"userId": userId,
			})
			if err != nil {
				return models.CreateResetPasswordTokenResponse{
					Status: constants.UnknownUserID,
				}
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				return models.CreateResetPasswordTokenResponse{
					Status: status.(string),
					Token:  response["token"].(string),
				}
			}
			return models.CreateResetPasswordTokenResponse{
				Status: constants.UnknownUserID,
			}
		},

		ResetPasswordUsingToken: func(token, newPassword string) models.ResetPasswordUsingTokenResponse {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/password/reset")
			if err != nil {
				return models.ResetPasswordUsingTokenResponse{
					Status: constants.ResetPasswordInvalidTokenError,
				}
			}
			response, err := querier.SendPostRequest(*normalisedURLPath, map[string]interface{}{
				"method":      "token",
				"token":       token,
				"newPassword": newPassword,
			})
			return models.ResetPasswordUsingTokenResponse{
				Status: response["status"].(string),
			}
		},
	}
}
