package thirdparty

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
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

func GetUserByID(userID string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByID(userID), nil
}

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

func Apple(config providers.TypeThirdPartyProviderAppleConfig) models.TypeProvider {
	return providers.Apple(config)
}

func Facebook(config providers.TypeThirdPartyProviderFacebookConfig) models.TypeProvider {
	return providers.Facebook(config)
}

func Github(config providers.TypeThirdPartyProviderGithubConfig) models.TypeProvider {
	return providers.Github(config)
}

func Google(config providers.TypeThirdPartyProviderGoogleConfig) models.TypeProvider {
	return providers.Google(config)
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
