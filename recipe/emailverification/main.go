package emailverification

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func EmailVerificationInit(config *models.TypeInput) supertokens.RecipeListFunction {
	return recipeInit(config)
}

func CreateEmailVerificationToken(userID, email string) (models.CreateEmailVerificationTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.CreateEmailVerificationTokenResponse{}, err
	}
	return instance.RecipeImpl.CreateEmailVerificationToken(userID, email)
}

func VerifyEmailUsingToken(token string) (models.VerifyEmailUsingTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.VerifyEmailUsingTokenResponse{}, err
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

func RevokeEmailVerificationTokens(userID, email string) (models.RevokeEmailVerificationTokensResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.RevokeEmailVerificationTokensResponse{}, err
	}
	return instance.RecipeImpl.RevokeEmailVerificationTokens(userID, email)
}

func UnverifyEmail(userID, email string) (models.UnverifyEmailResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.UnverifyEmailResponse{}, err
	}
	return instance.RecipeImpl.UnverifyEmail(userID, email)
}
