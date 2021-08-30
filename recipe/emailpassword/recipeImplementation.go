package emailpassword

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/constants"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier) models.RecipeImplementation {
	return models.RecipeImplementation{
		SignUp: func(email, password string) (models.SignInUpResponse, error) {
			response, err := querier.SendPostRequest("/recipe/signup", map[string]interface{}{
				"email":    email,
				"password": password,
			})
			if err != nil {
				return models.SignInUpResponse{}, err
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return models.SignInUpResponse{}, err
				}
				return models.SignInUpResponse{
					Status: status.(string),
					User:   *user,
				}, nil
			}
			return models.SignInUpResponse{
				Status: constants.EmailAlreadyExistsError,
			}, nil
		},

		SignIn: func(email, password string) (models.SignInUpResponse, error) {
			response, err := querier.SendPostRequest("/recipe/signin", map[string]interface{}{
				"email":    email,
				"password": password,
			})
			if err != nil {
				return models.SignInUpResponse{}, err
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return models.SignInUpResponse{}, err
				}
				return models.SignInUpResponse{
					Status: status.(string),
					User:   *user,
				}, nil
			}
			return models.SignInUpResponse{
				Status: constants.WrongCredentialsError,
			}, nil
		},

		GetUserByID: func(userID string) (*models.User, error) {
			response, err := querier.SendGetRequest("/recipe/user", map[string]interface{}{
				"userId": userID,
			})
			if err != nil {
				return nil, err
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return nil, err
				}
				return user, nil
			}
			return nil, nil
		},

		GetUserByEmail: func(email string) (*models.User, error) {
			response, err := querier.SendGetRequest("/recipe/user", map[string]interface{}{
				"email": email,
			})
			if err != nil {
				return nil, err
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return nil, err
				}
				return user, nil
			}
			return nil, nil
		},

		CreateResetPasswordToken: func(userID string) (models.CreateResetPasswordTokenResponse, error) {
			response, err := querier.SendPostRequest("/recipe/user/password/reset/token", map[string]interface{}{
				"userId": userID,
			})
			if err != nil {
				return models.CreateResetPasswordTokenResponse{}, err
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				return models.CreateResetPasswordTokenResponse{
					Status: status.(string),
					Token:  response["token"].(string),
				}, nil
			}
			return models.CreateResetPasswordTokenResponse{
				Status: constants.UnknownUserID,
			}, nil
		},

		ResetPasswordUsingToken: func(token, newPassword string) (models.ResetPasswordUsingTokenResponse, error) {
			response, err := querier.SendPostRequest("/recipe/user/password/reset", map[string]interface{}{
				"method":      "token",
				"token":       token,
				"newPassword": newPassword,
			})
			if err != nil {
				return models.ResetPasswordUsingTokenResponse{}, nil
			}
			return models.ResetPasswordUsingTokenResponse{
				Status: response["status"].(string),
			}, nil
		},

		UpdateEmailOrPassword: func(userId string, email, password *string) (models.UpdateEmailOrPasswordResponse, error) {
			requestBody := map[string]interface{}{
				"userId": userId,
			}
			if email != nil {
				requestBody["email"] = email
			}
			if password != nil {
				requestBody["password"] = password
			}
			response, err := querier.SendPutRequest("/recipe/user", requestBody)
			if err != nil {
				return models.UpdateEmailOrPasswordResponse{}, nil
			}
			return models.UpdateEmailOrPasswordResponse{
				Status: response["status"].(string),
			}, nil
		},
	}
}
