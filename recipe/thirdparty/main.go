package thirdparty

import (
	"errors"

	envModels "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type signInUpResponse struct {
	CreatedNewUser bool
	User           models.User
}

func Init(config *models.TypeInput) supertokens.RecipeListFunction {
	return recipeInit(config)
}

func SignInUp(thirdPartyID string, thirdPartyUserID string, email models.EmailStruct) (models.SignInUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.SignInUpResponse{}, err
	}
	return instance.RecipeImpl.SignInUp(thirdPartyID, thirdPartyUserID, email)
}

func GetUserByID(userID string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByID(userID)
}

func GetUsersByEmail(email string) ([]models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return []models.User{}, err
	}
	return instance.RecipeImpl.GetUsersByEmail(email)
}

func GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID)
}

func CreateEmailVerificationToken(userID string) (envModels.CreateEmailVerificationTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return envModels.CreateEmailVerificationTokenResponse{}, err
	}
	email, err := instance.getEmailForUserId(userID)
	if err != nil {
		return envModels.CreateEmailVerificationTokenResponse{}, err
	}
	return instance.EmailVerificationRecipe.RecipeImpl.CreateEmailVerificationToken(userID, email)
}

func VerifyEmailUsingToken(token string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	response, err := instance.EmailVerificationRecipe.RecipeImpl.VerifyEmailUsingToken(token)
	if err != nil {
		return nil, err
	}
	if response.EmailVerificationInvalidTokenError != nil {
		return nil, errors.New("email verification token is invalid")
	}
	return instance.RecipeImpl.GetUserByID(response.OK.User.ID)
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

func RevokeEmailVerificationTokens(userID string) (envModels.RevokeEmailVerificationTokensResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return envModels.RevokeEmailVerificationTokensResponse{}, err
	}
	email, err := instance.getEmailForUserId(userID)
	if err != nil {
		return envModels.RevokeEmailVerificationTokensResponse{}, err
	}
	return instance.EmailVerificationRecipe.RecipeImpl.RevokeEmailVerificationTokens(userID, email)
}

func UnverifyEmail(userID string) (envModels.UnverifyEmailResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return envModels.UnverifyEmailResponse{}, err
	}
	email, err := instance.getEmailForUserId(userID)
	if err != nil {
		return envModels.UnverifyEmailResponse{}, err
	}
	return instance.EmailVerificationRecipe.RecipeImpl.UnverifyEmail(userID, email)
}

// func Apple(config providers.TypeThirdPartyProviderAppleConfig) models.TypeProvider {
// 	return providers.Apple(config)
// }

func Facebook(config providers.TypeThirdPartyProviderFacebookConfig) models.TypeProvider {
	return providers.Facebook(config)
}

func Github(config providers.TypeThirdPartyProviderGithubConfig) models.TypeProvider {
	return providers.Github(config)
}

func Google(config providers.TypeThirdPartyProviderGoogleConfig) models.TypeProvider {
	return providers.Google(config)
}
