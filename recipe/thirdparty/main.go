package thirdparty

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
)

type signInUpResponse struct {
	CreatedNewUser bool
	User           models.User
}

func SignInUp(thirdPartyID string, thirdPartyUserID string, email models.EmailStruct) (*signInUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	result := instance.RecipeImpl.SignInUp(thirdPartyID, thirdPartyUserID, email)
	if result.Status == "OK" {
		return &signInUpResponse{
			CreatedNewUser: result.CreatedNewUser,
			User:           result.User,
		}, nil
	}
	return nil, result.Error
}

func GetUserById(userID string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserById(userID), nil
}

// TODO
func GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID), nil
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

func VerifyEmailUsingToken(token string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	user, err := instance.EmailVerificationRecipe.VerifyEmailUsingToken(token)
	if err != nil {
		return nil, err
	}
	userInThisRecipe := instance.RecipeImpl.GetUserById(user.ID)
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
	userInfo := instance.RecipeImpl.GetUserById(userID)
	if userInfo == nil {
		return "", errors.New("Unknown User ID provided")
	}
	return userInfo.Email, nil
}
