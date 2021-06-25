package emailpassword

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// TODO:
func MakeRecipeImplementation(querier supertokens.Querier) models.RecipeImplementation {
	return models.RecipeImplementation{
		SignUp: func(email, password string) models.SignInUpResponse {
			return models.SignInUpResponse{}
		},

		SignIn: func(email, password string) models.SignInUpResponse {
			return models.SignInUpResponse{}
		},

		GetUserById: func(userId string) *models.User {
			return nil
		},

		GetUserByEmail: func(email string) *models.User {
			return nil
		},

		CreateResetPasswordToken: func(userId string) models.CreateResetPasswordTokenResponse {
			return models.CreateResetPasswordTokenResponse{}
		},

		ResetPasswordUsingToken: func(token, newPassword string) models.ResetPasswordUsingTokenResponse {
			return models.ResetPasswordUsingTokenResponse{}
		},
	}
}
