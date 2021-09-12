package thirdpartyemailpassword

import (
	"errors"

	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	envModels "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

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

func GetUserByThirdPartyInfo(thirdPartyID string, thirdPartyUserID string, email tpm.EmailStruct) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID)
}

func SignUp(email, password string) (models.SignUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.SignUpResponse{}, err
	}
	return instance.RecipeImpl.SignUp(email, password)
}

func SignIn(email, password string) (models.SignInResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return models.SignInResponse{}, err
	}
	return instance.RecipeImpl.SignIn(email, password)
}

func GetUserById(userID string) (*models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByID(userID)
}

func GetUsersByEmail(email string) ([]models.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUsersByEmail(email)
}

func CreateResetPasswordToken(userID string) (epm.CreateResetPasswordTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epm.CreateResetPasswordTokenResponse{}, err
	}
	return instance.RecipeImpl.CreateResetPasswordToken(userID)
}

func ResetPasswordUsingToken(token, newPassword string) (epm.ResetPasswordUsingTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epm.ResetPasswordUsingTokenResponse{}, err
	}
	return instance.RecipeImpl.ResetPasswordUsingToken(token, newPassword)
}

func UpdateEmailOrPassword(userId string, email *string, password *string) (epm.UpdateEmailOrPasswordResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epm.UpdateEmailOrPasswordResponse{}, err
	}
	return instance.RecipeImpl.UpdateEmailOrPassword(userId, email, password)
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

func Google(config providers.GoogleConfig) tpm.TypeProvider {
	return providers.Google(config)
}

func Github(config providers.GithubConfig) tpm.TypeProvider {
	return providers.Github(config)
}

func Facebook(config providers.FacebookConfig) tpm.TypeProvider {
	return providers.Facebook(config)
}

// func Apple(config providers.AppleConfig) tpm.TypeProvider {
// 	return providers.Apple(config)
// }
