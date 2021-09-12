package thirdpartyemailpassword

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *tpepmodels.TypeInput) supertokens.RecipeListFunction {
	return recipeInit(config)
}

func SignInUp(thirdPartyID string, thirdPartyUserID string, email tpepmodels.EmailStruct) (tpepmodels.SignInUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpepmodels.SignInUpResponse{}, err
	}
	return instance.RecipeImpl.SignInUp(thirdPartyID, thirdPartyUserID, email)
}

func GetUserByThirdPartyInfo(thirdPartyID string, thirdPartyUserID string, email tpmodels.EmailStruct) (*tpepmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID)
}

func SignUp(email, password string) (tpepmodels.SignUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpepmodels.SignUpResponse{}, err
	}
	return instance.RecipeImpl.SignUp(email, password)
}

func SignIn(email, password string) (tpepmodels.SignInResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpepmodels.SignInResponse{}, err
	}
	return instance.RecipeImpl.SignIn(email, password)
}

func GetUserById(userID string) (*tpepmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByID(userID)
}

func GetUsersByEmail(email string) ([]tpepmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUsersByEmail(email)
}

func CreateResetPasswordToken(userID string) (epmodels.CreateResetPasswordTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.CreateResetPasswordTokenResponse{}, err
	}
	return instance.RecipeImpl.CreateResetPasswordToken(userID)
}

func ResetPasswordUsingToken(token, newPassword string) (epmodels.ResetPasswordUsingTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.ResetPasswordUsingTokenResponse{}, err
	}
	return instance.RecipeImpl.ResetPasswordUsingToken(token, newPassword)
}

func UpdateEmailOrPassword(userId string, email *string, password *string) (epmodels.UpdateEmailOrPasswordResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return epmodels.UpdateEmailOrPasswordResponse{}, err
	}
	return instance.RecipeImpl.UpdateEmailOrPassword(userId, email, password)
}

func CreateEmailVerificationToken(userID string) (evmodels.CreateEmailVerificationTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	email, err := instance.getEmailForUserId(userID)
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	return instance.EmailVerificationRecipe.RecipeImpl.CreateEmailVerificationToken(userID, email)
}

func VerifyEmailUsingToken(token string) (*tpepmodels.User, error) {
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

func RevokeEmailVerificationTokens(userID string) (evmodels.RevokeEmailVerificationTokensResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.RevokeEmailVerificationTokensResponse{}, err
	}
	email, err := instance.getEmailForUserId(userID)
	if err != nil {
		return evmodels.RevokeEmailVerificationTokensResponse{}, err
	}
	return instance.EmailVerificationRecipe.RecipeImpl.RevokeEmailVerificationTokens(userID, email)
}

func UnverifyEmail(userID string) (evmodels.UnverifyEmailResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	email, err := instance.getEmailForUserId(userID)
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	return instance.EmailVerificationRecipe.RecipeImpl.UnverifyEmail(userID, email)
}

func Google(config tpmodels.GoogleConfig) tpmodels.TypeProvider {
	return providers.Google(config)
}

func Github(config tpmodels.GithubConfig) tpmodels.TypeProvider {
	return providers.Github(config)
}

func Facebook(config tpmodels.FacebookConfig) tpmodels.TypeProvider {
	return providers.Facebook(config)
}

// func Apple(config providers.AppleConfig) tpm.TypeProvider {
// 	return providers.Apple(config)
// }
