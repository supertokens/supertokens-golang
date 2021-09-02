package thirdpartyemailpassword

import (
	"errors"

	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func SignInUp(thirdPartyID string, thirdPartyUserID string, email tpm.EmailStruct) (models.SignInUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.SignInUpResponse{}, err
	}
	return instance.RecipeImpl.SignInUp(thirdPartyID, thirdPartyUserID, email), nil
}

func GetUserByThirdPartyInfo(thirdPartyID string, thirdPartyUserID string, email tpm.EmailStruct) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID), nil
}

func SignUp(email, password string) (models.SignInUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.SignInUpResponse{}, err
	}
	return instance.RecipeImpl.SignUp(email, password), nil
}

func SignIn(email, password string) (models.SignInUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.SignInUpResponse{}, err
	}
	return instance.RecipeImpl.SignIn(email, password), nil
}

func GetUserById(userID string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByID(userID), nil
}

func GetUserByEmail(email string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByEmail(email), nil
}

func CreateResetPasswordToken(userID string) (epm.CreateResetPasswordTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epm.CreateResetPasswordTokenResponse{}, err
	}
	return instance.RecipeImpl.CreateResetPasswordToken(userID), nil
}

func ResetPasswordUsingToken(token, newPassword string) (epm.ResetPasswordUsingTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epm.ResetPasswordUsingTokenResponse{}, err
	}
	return instance.RecipeImpl.ResetPasswordUsingToken(token, newPassword), nil
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

func Google(config providers.TypeThirdPartyProviderGoogleConfig) tpm.TypeProvider {
	return providers.Google(config)
}

func Github(config providers.TypeThirdPartyProviderGithubConfig) tpm.TypeProvider {
	return providers.Github(config)
}

func Facebook(config providers.TypeThirdPartyProviderFacebookConfig) tpm.TypeProvider {
	return providers.Facebook(config)
}

func Apple(config providers.TypeThirdPartyProviderAppleConfig) tpm.TypeProvider {
	return providers.Apple(config)
}
