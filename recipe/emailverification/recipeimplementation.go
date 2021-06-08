package emailverification

import (
	"fmt"
	"strconv"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
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

// Implementing RecipeInterface

func (r *RecipeImplementation) CreateEmailVerificationToken(userID, email string) schema.ReturnMap {
	normalisedURLPath, err := supertokens.NewNormalisedURLPath("/recipe/user/email/verify/token")
	if err != nil {
		fmt.Println(err) // todo handle error
		return nil
	}
	response, _ := r.Querier.SendPostRequest(*normalisedURLPath, map[string]string{
		"userId": userID,
		"email":  email,
	})
	if response["status"] == "OK" {
		resp := schema.CreateEmailVerificationTokenOk
		resp["token"] = response["token"]
		return resp
	}
	return schema.CreateEmailVerificationTokenError
}

func (r *RecipeImplementation) VerifyEmailUsingToken(token string) schema.ReturnMap {
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

func (r *RecipeImplementation) IsEmailVerified(userID, email string) bool {
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
