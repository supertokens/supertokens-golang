package emailpassword

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
)

func SignUp(email string, password string) (models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.User{}, err
	}
	return instance.RecipeImpl.SignUp(email, password).User, nil
}

func SignIn(email string, password string) (models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.User{}, err
	}
	return instance.RecipeImpl.SignIn(email, password).User, nil
}

func GetUserById(userID string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserById(userID), nil
}

func GetUserByEmail(email string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByEmail(email), nil
}

func CreateResetPasswordToken(userID string) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	return instance.RecipeImpl.CreateResetPasswordToken(userID).Token, nil
}

func ResetPasswordUsingToken(token string, newPassword string) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	instance.RecipeImpl.ResetPasswordUsingToken(token, newPassword)
	return nil
}

func CreateEmailVerificationToken(userID string) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	email, err := getEmailForUserId(userID)
	if err != nil {
		return "", err
	}
	return instance.EmailVerificationRecipe.CreateEmailVerificationToken(userID, email)
}

func VerifyEmailUsingToken(token string) (models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.User{}, err
	}
	user, err := instance.EmailVerificationRecipe.VerifyEmailUsingToken(token)
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		ID:    user.ID,
		Email: user.Email,
	}, nil
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
	userInfo := instance.RecipeImpl.GetUserById(userID)
	if userInfo == nil {
		return "", errors.New("Unknown User ID provided")
	}
	return userInfo.Email, nil
}
