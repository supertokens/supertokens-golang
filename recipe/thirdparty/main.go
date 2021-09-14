package thirdparty

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type signInUpResponse struct {
	CreatedNewUser bool
	User           tpmodels.User
}

func Init(config *tpmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func SignInUp(thirdPartyID string, thirdPartyUserID string, email tpmodels.EmailStruct) (tpmodels.SignInUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.SignInUpResponse{}, err
	}
	return instance.RecipeImpl.SignInUp(thirdPartyID, thirdPartyUserID, email)
}

func GetUserByID(userID string) (*tpmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByID(userID)
}

func GetUsersByEmail(email string) ([]tpmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return []tpmodels.User{}, err
	}
	return instance.RecipeImpl.GetUsersByEmail(email)
}

func GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID string) (*tpmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID)
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

func VerifyEmailUsingToken(token string) (*tpmodels.User, error) {
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

// func Apple(config tpmodels.AppleConfig) tpmodels.TypeProvider {
// 	return providers.Apple(config)
// }

func Facebook(config tpmodels.FacebookConfig) tpmodels.TypeProvider {
	return providers.Facebook(config)
}

func Github(config tpmodels.GithubConfig) tpmodels.TypeProvider {
	return providers.Github(config)
}

func Google(config tpmodels.GoogleConfig) tpmodels.TypeProvider {
	return providers.Google(config)
}
