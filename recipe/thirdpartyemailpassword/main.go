package thirdpartyemailpassword

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func VerifyEmailUsingToken(token string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	user, err := instance.EmailVerificationRecipe.VerifyEmailUsingToken(token)
	if err != nil {
		return nil, err
	}
	userInThisRecipe := instance.RecipeImpl.GetUserByID(user.ID)
	if userInThisRecipe == nil {
		return nil, errors.New("Unknown User ID provided")
	}
	return userInThisRecipe, nil
}

func IsEmailVerified(userID string) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	email, err := getEmailForUserId(userID)
	if err != nil {
		return false, err
	}
	return instance.EmailVerificationRecipe.RecipeImpl.IsEmailVerified(userID, email)
}

func getEmailForUserId(userID string) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	userInfo := instance.RecipeImpl.GetUserByID(userID)
	if userInfo == nil {
		return "", errors.New("Unknown User ID provided")
	}
	return userInfo.Email, nil
}
