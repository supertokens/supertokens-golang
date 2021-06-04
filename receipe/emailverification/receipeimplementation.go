package emailverification

import (
	"fmt"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type RecipeImplementation struct {
	Querier supertokens.Querier
}

func NewRecipeImplementation(querier supertokens.Querier) *RecipeImplementation {
	return &RecipeImplementation{
		Querier: querier,
	}
}

func (r *RecipeImplementation) createEmailVerificationToken(userId, email string) map[string]interface{} {
	normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify/token")
	if err != nil {
		fmt.Println(err) // todo handle error
		return nil
	}
	response, _ := r.Querier.SendPostRequest(*normalisedURLPath, map[string]interface{}{
		"userId": userId,
		"email":  email,
	})
	if response["status"] == "OK" {
		return map[string]interface{}{
			"status": "OK",
			"token":  response["token"],
		}
	}
	return map[string]interface{}{
		"status": "EMAIL_ALREADY_VERIFIED_ERROR",
	}
}

func (r *RecipeImplementation) verifyEmailUsingToken(token string) map[string]interface{} {
	normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify")
	if err != nil {
		fmt.Println(err) // todo handle error
		return nil
	}
	response, _ := r.Querier.SendPostRequest(*normalisedURLPath, map[string]interface{}{
		"method": "token",
		"token":  token,
	})
	if response["status"] == "OK" {
		return map[string]interface{}{
			"status": "OK",
			"user": map[string]interface{}{
				"id":    response["userId"],
				"email": response["email"],
			},
		}
	}
	return map[string]interface{}{
		"status": "EMAIL_VERIFICATION_INVALID_TOKEN_ERROR",
	}
}

func (r *RecipeImplementation) isEmailVerified(userId, email string) bool {
	normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify")
	if err != nil {
		fmt.Println(err) // todo handle error
		return false
	}
	response, _ := r.Querier.SendGetRequest(*normalisedURLPath, map[string]string{
		"userId": userId,
		"email":  email,
	})

	return response["isVerified"].(bool)
}
