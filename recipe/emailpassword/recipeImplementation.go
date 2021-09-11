package emailpassword

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier) models.RecipeInterface {
	return models.RecipeInterface{
		SignUp: func(email, password string) (models.SignUpResponse, error) {
			response, err := querier.SendPostRequest("/recipe/signup", map[string]interface{}{
				"email":    email,
				"password": password,
			})
			if err != nil {
				return models.SignUpResponse{}, err
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return models.SignUpResponse{}, err
				}
				return models.SignUpResponse{
					OK: &struct{ User models.User }{User: *user},
				}, nil
			}
			return models.SignUpResponse{
				EmailAlreadyExistsError: &struct{}{},
			}, nil
		},

		SignIn: func(email, password string) (models.SignInResponse, error) {
			response, err := querier.SendPostRequest("/recipe/signin", map[string]interface{}{
				"email":    email,
				"password": password,
			})
			if err != nil {
				return models.SignInResponse{}, err
			}
			status, ok := response["status"]
			if ok && status.(string) == "OK" {
				user, err := parseUser(response["user"])
				if err != nil {
					return models.SignInResponse{}, err
				}
				return models.SignInResponse{
					OK: &struct{ User models.User }{User: *user},
				}, nil
			}
			return models.SignInResponse{
				WrongCredentialsError: &struct{}{},
			}, nil
		},

		GetUserByID: func(userID string) (*models.User, error) {
			response, err := querier.SendGetRequest("/recipe/user", map[string]string{
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
			response, err := querier.SendGetRequest("/recipe/user", map[string]string{
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
					OK: &struct{ Token string }{Token: response["token"].(string)},
				}, nil
			}
			return models.CreateResetPasswordTokenResponse{
				UnknownUserIdError: &struct{}{},
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

			if response["status"].(string) == "OK" {
				return models.ResetPasswordUsingTokenResponse{
					OK: &struct{}{},
				}, nil
			} else {
				return models.ResetPasswordUsingTokenResponse{
					ResetPasswordInvalidTokenError: &struct{}{},
				}, nil
			}
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

			if response["status"].(string) == "OK" {
				return models.UpdateEmailOrPasswordResponse{
					OK: &struct{}{},
				}, nil
			} else if response["status"].(string) == "EMAIL_ALREADY_EXISTS_ERROR" {
				return models.UpdateEmailOrPasswordResponse{
					EmailAlreadyExistsError: &struct{}{},
				}, nil
			} else {
				return models.UpdateEmailOrPasswordResponse{
					UnknownUserIdError: &struct{}{},
				}, nil
			}
		},
	}
}
