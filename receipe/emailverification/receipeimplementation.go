package emailverification

import (
	"fmt"
	"strconv"

	"github.com/supertokens/supertokens-golang/receipe/emailverification/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type RecipeImplementation schema.RecipeImplementation

func NewRecipeImplementation(querier supertokens.Querier) *schema.RecipeImplementation {
	return &schema.RecipeImplementation{
		Querier: querier,
	}
}

func (r *RecipeImplementation) createEmailVerificationToken(userId, email string) schema.ReturnMap {
	normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify/token")
	if err != nil {
		fmt.Println(err) // todo handle error
		return nil
	}
	response, _ := r.Querier.SendPostRequest(*normalisedURLPath, map[string]string{
		"userId": userId,
		"email":  email,
	})
	if response["status"] == "OK" {
		resp := schema.CreateEmailVerificationTokenOk
		resp["token"] = response["token"]
		return resp
	}
	return schema.CreateEmailVerificationTokenError
}

func (r *RecipeImplementation) verifyEmailUsingToken(token string) schema.ReturnMap {
	normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify")
	if err != nil {
		fmt.Println(err) // todo handle error
		return nil
	}

	response, _ := r.Querier.SendPostRequest(*normalisedURLPath, map[string]string{
		"method": "token",
		"token":  token,
	})

	if response["status"] == "OK" {
		resp := schema.VerifyEmailUsingTokenOk
		resp["user"] = schema.User{
			ID:    response["userId"],
			Email: response["email"],
		}
		return resp
	}
	return schema.VerifyEmailUsingTokenError
}

func (r *RecipeImplementation) isEmailVerified(userID, email string) bool {
	normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify")
	if err != nil {
		fmt.Println(err) // todo handle error
		return false
	}
	response, _ := r.Querier.SendGetRequest(*normalisedURLPath, map[string]string{
		"userId": userID,
		"email":  email,
	})

	isVerified, err := strconv.ParseBool(response["isVerified"])
	if err != nil {
		fmt.Println(err) // todo handle error
		return false
	}
	return isVerified
}
