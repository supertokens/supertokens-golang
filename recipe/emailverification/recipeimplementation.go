package emailverification

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier) evmodels.RecipeInterface {
	return evmodels.RecipeInterface{
		CreateEmailVerificationToken: func(userID, email string) (evmodels.CreateEmailVerificationTokenResponse, error) {
			response, err := querier.SendPostRequest("/recipe/user/email/verify/token", map[string]interface{}{
				"userId": userID,
				"email":  email,
			})
			if err != nil {
				return evmodels.CreateEmailVerificationTokenResponse{}, err
			}
			status, ok := response["status"]
			if ok && status == "OK" {
				return evmodels.CreateEmailVerificationTokenResponse{
					OK: &struct{ Token string }{Token: response["token"].(string)},
				}, nil
			}

			return evmodels.CreateEmailVerificationTokenResponse{
				EmailAlreadyVerifiedError: &struct{}{},
			}, nil
		},

		VerifyEmailUsingToken: func(token string) (evmodels.VerifyEmailUsingTokenResponse, error) {
			response, err := querier.SendPostRequest("/recipe/user/email/verify", map[string]interface{}{
				"method": "token",
				"token":  token,
			})
			if err != nil {
				return evmodels.VerifyEmailUsingTokenResponse{}, err
			}
			status, ok := response["status"]
			if ok && status == "OK" {
				return evmodels.VerifyEmailUsingTokenResponse{
					OK: &struct{ User evmodels.User }{User: evmodels.User{
						ID:    response["userId"].(string),
						Email: response["email"].(string),
					}},
				}, nil
			}
			return evmodels.VerifyEmailUsingTokenResponse{
				EmailVerificationInvalidTokenError: &struct{}{},
			}, nil
		},

		IsEmailVerified: func(userID, email string) (bool, error) {
			response, err := querier.SendGetRequest("/recipe/user/email/verify", map[string]string{
				"userId": userID,
				"email":  email,
			})
			if err != nil {
				return false, err
			}
			return response["isVerified"].(bool), nil
		},

		RevokeEmailVerificationTokens: func(userId string, email string) (evmodels.RevokeEmailVerificationTokensResponse, error) {
			_, err := querier.SendPostRequest("/recipe/user/email/verify/token/remove", map[string]interface{}{
				"userId": userId,
				"email":  email,
			})
			if err != nil {
				return evmodels.RevokeEmailVerificationTokensResponse{}, err
			}
			return evmodels.RevokeEmailVerificationTokensResponse{
				OK: &struct{}{},
			}, nil
		},

		UnverifyEmail: func(userId string, email string) (evmodels.UnverifyEmailResponse, error) {
			_, err := querier.SendPostRequest("/recipe/user/email/verify/remove", map[string]interface{}{
				"userId": userId,
				"email":  email,
			})
			if err != nil {
				return evmodels.UnverifyEmailResponse{}, err
			}
			return evmodels.UnverifyEmailResponse{
				OK: &struct{}{},
			}, nil
		},
	}
}
