package emailverification

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier) models.RecipeImplementation {
	return models.RecipeImplementation{
		CreateEmailVerificationToken: func(userID, email string) (*models.CreateEmailVerificationTokenResponse, error) {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify/token")
			if err != nil {
				return nil, err
			}
			response, err := querier.SendPostRequest(*normalisedURLPath, map[string]interface{}{
				"userId": userID,
				"email":  email,
			})
			if err != nil {
				return nil, err
			}
			status, ok := response["status"]
			if ok && status == "OK" {
				resp := models.CreateEmailVerificationTokenResponse{
					Status: "OK",
					Token:  response["token"].(string),
				}
				return &resp, nil
			}

			return &models.CreateEmailVerificationTokenResponse{
				Status: "EMAIL_ALREADY_VERIFIED_ERROR",
			}, nil
		},
		VerifyEmailUsingToken: func(token string) (*models.VerifyEmailUsingTokenResponse, error) {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify")
			if err != nil {
				return nil, err
			}

			response, err := querier.SendPostRequest(*normalisedURLPath, map[string]interface{}{
				"method": "token",
				"token":  token,
			})
			if err != nil {
				return nil, err
			}
			status, ok := response["status"]
			if ok && status == "OK" {
				return &models.VerifyEmailUsingTokenResponse{
					Status: "OK",
					User: models.User{
						ID:    response["userId"].(string),
						Email: response["email"].(string),
					},
				}, nil
			}
			return &models.VerifyEmailUsingTokenResponse{
				Status: "EMAIL_VERIFICATION_INVALID_TOKEN_ERROR",
			}, nil
		},
		IsEmailVerified: func(userID, email string) (bool, error) {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify")
			if err != nil {
				return false, err
			}
			response, err := querier.SendGetRequest(*normalisedURLPath, map[string]interface{}{
				"userId": userID,
				"email":  email,
			})
			if err != nil {
				return false, err
			}
			isVerified := response["isVerified"].(bool)
			return isVerified, nil
		},
	}
}
