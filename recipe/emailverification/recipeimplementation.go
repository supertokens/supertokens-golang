package emailverification

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier) models.RecipeInterface {
	return models.RecipeInterface{
		CreateEmailVerificationToken: func(userID, email string) (models.CreateEmailVerificationTokenResponse, error) {
			response, err := querier.SendPostRequest("/recipe/user/email/verify/token", map[string]interface{}{
				"userId": userID,
				"email":  email,
			})
			if err != nil {
				return models.CreateEmailVerificationTokenResponse{}, err
			}
			status, ok := response["status"]
			if ok && status == "OK" {
				return models.CreateEmailVerificationTokenResponse{
					OK: &struct{ Token string }{Token: response["token"].(string)},
				}, nil
			}

			return models.CreateEmailVerificationTokenResponse{
				EmailAlreadyVerifiedError: &struct{}{},
			}, nil
		},

		VerifyEmailUsingToken: func(token string) (models.VerifyEmailUsingTokenResponse, error) {
			response, err := querier.SendPostRequest("/recipe/user/email/verify", map[string]interface{}{
				"method": "token",
				"token":  token,
			})
			if err != nil {
				return models.VerifyEmailUsingTokenResponse{}, err
			}
			status, ok := response["status"]
			if ok && status == "OK" {
				return models.VerifyEmailUsingTokenResponse{
					OK: &struct{ User models.User }{User: models.User{
						ID:    response["userId"].(string),
						Email: response["email"].(string),
					}},
				}, nil
			}
			return models.VerifyEmailUsingTokenResponse{
				EmailVerificationInvalidTokenError: &struct{}{},
			}, nil
		},

		IsEmailVerified: func(userID, email string) (bool, error) {
			response, err := querier.SendGetRequest("/recipe/user/email/verify", map[string]interface{}{
				"userId": userID,
				"email":  email,
			})
			if err != nil {
				return false, err
			}
			return response["isVerified"].(bool), nil
		},

		RevokeEmailVerificationTokens: func(userId string, email string) (models.RevokeEmailVerificationTokensResponse, error) {
			_, err := querier.SendPostRequest("/recipe/user/email/verify/token/remove", map[string]interface{}{
				"userId": userId,
				"email":  email,
			})
			if err != nil {
				return models.RevokeEmailVerificationTokensResponse{}, err
			}
			return models.RevokeEmailVerificationTokensResponse{
				OK: &struct{}{},
			}, nil
		},

		UnverifyEmail: func(userId string, email string) (models.UnverifyEmailResponse, error) {
			_, err := querier.SendPostRequest("/recipe/user/email/verify/remove", map[string]interface{}{
				"userId": userId,
				"email":  email,
			})
			if err != nil {
				return models.UnverifyEmailResponse{}, err
			}
			return models.UnverifyEmailResponse{
				OK: &struct{}{},
			}, nil
		},
	}
}
