package emailverification

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *evmodels.TypeInput) supertokens.RecipeListFunction {
	return recipeInit(config)
}

func CreateEmailVerificationToken(userID, email string) (evmodels.CreateEmailVerificationTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	return instance.RecipeImpl.CreateEmailVerificationToken(userID, email)
}

func VerifyEmailUsingToken(token string) (evmodels.VerifyEmailUsingTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.VerifyEmailUsingTokenResponse{}, err
	}
	return instance.RecipeImpl.VerifyEmailUsingToken(token)
}

func IsEmailVerified(userID, email string) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return instance.RecipeImpl.IsEmailVerified(userID, email)
}

func RevokeEmailVerificationTokens(userID, email string) (evmodels.RevokeEmailVerificationTokensResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.RevokeEmailVerificationTokensResponse{}, err
	}
	return instance.RecipeImpl.RevokeEmailVerificationTokens(userID, email)
}

func UnverifyEmail(userID, email string) (evmodels.UnverifyEmailResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	return instance.RecipeImpl.UnverifyEmail(userID, email)
}
