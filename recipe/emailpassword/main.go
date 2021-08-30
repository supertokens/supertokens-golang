package emailpassword

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
)

func SignUp(email string, password string) (models.SignUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.SignUpResponse{}, err
	}
	return instance.RecipeImpl.SignUp(email, password)
}

func SignIn(email string, password string) (models.SignInResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.SignInResponse{}, err
	}
	return instance.RecipeImpl.SignIn(email, password)
}

func GetUserByID(userID string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByID(userID)
}

func GetUserByEmail(email string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByEmail(email)
}

func CreateResetPasswordToken(userID string) (models.CreateResetPasswordTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.CreateResetPasswordTokenResponse{}, err
	}
	return instance.RecipeImpl.CreateResetPasswordToken(userID)
}

func ResetPasswordUsingToken(token string, newPassword string) (models.ResetPasswordUsingTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.ResetPasswordUsingTokenResponse{}, nil
	}
	return instance.RecipeImpl.ResetPasswordUsingToken(token, newPassword)
}

func UpdateEmailOrPassword(userId string, email *string, password *string) (models.UpdateEmailOrPasswordResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.UpdateEmailOrPasswordResponse{}, nil
	}
	return instance.RecipeImpl.UpdateEmailOrPassword(userId, email, password)
}

// TODO: fix the functions below.
func CreateEmailVerificationToken(userID string) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	email, err := instance.getEmailForUserId(userID)
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
	email, err := instance.getEmailForUserId(userID)
	if err != nil {
		return false, err
	}
	return instance.EmailVerificationRecipe.RecipeImpl.IsEmailVerified(userID, email)
}
