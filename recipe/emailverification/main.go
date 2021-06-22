package emailverification

import "github.com/supertokens/supertokens-golang/recipe/emailverification/models"

func CreateEmailVerificationToken(userID, email string) (string, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	return instance.CreateEmailVerificationToken(userID, email)
}

func VerifyEmailUsingToken(token string) (*models.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.VerifyEmailUsingToken(token)
}

func IsEmailVerified(userID, email string) (bool, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return instance.RecipeImpl.IsEmailVerified(userID, email)
}
