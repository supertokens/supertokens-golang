package emailverification

import (
	"strconv"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier) schema.RecipeImplementation {
	return schema.RecipeImplementation{
		CreateEmailVerificationToken: func(userID, email string) (*schema.CreateEmailVerificationTokenResponse, error) {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify/token")
			if err != nil {
				return nil, err
			}
			response, _ := querier.SendPostRequest(*normalisedURLPath, map[string]string{
				"userId": userID,
				"email":  email,
			})
			if response["status"] == "OK" {
				resp := schema.CreateEmailVerificationTokenResponse{
					Ok: &struct{ Token string }{Token: response["token"]},
				}
				return &resp, nil
			}

			return &schema.CreateEmailVerificationTokenResponse{
				EmailAlreadyVerifiedError: true,
			}, nil
		},
		VerifyEmailUsingToken: func(token string) (*schema.VerifyEmailUsingTokenResponse, error) {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify")
			if err != nil {
				return nil, err
			}

			response, _ := querier.SendPostRequest(*normalisedURLPath, map[string]string{
				"method": "token",
				"token":  token,
			})

			if response["status"] == "OK" {
				return &schema.VerifyEmailUsingTokenResponse{
					Ok: &struct{ User schema.User }{User: schema.User{
						ID:    response["userId"],
						Email: response["email"],
					}},
				}, nil
			}
			return &schema.VerifyEmailUsingTokenResponse{InvalidTokenError: true}, nil
		},
		IsEmailVerified: func(userID, email string) (bool, error) {
			normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify")
			if err != nil {
				return false, err
			}
			response, _ := querier.SendGetRequest(*normalisedURLPath, map[string]string{
				"userId": userID,
				"email":  email,
			})

			isVerified, err := strconv.ParseBool(response["isVerified"])
			if err != nil {
				return false, err
			}
			return isVerified, nil
		},
	}
}
